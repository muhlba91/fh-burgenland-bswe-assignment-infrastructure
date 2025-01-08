import * as aws from '@pulumi/aws';
import { interpolate, Output } from '@pulumi/pulumi';

import {
  awsAccountId,
  awsDefaultRegion,
  commonLabels,
  environment,
  githubOrganisation,
  globalName,
} from '../configuration';
import { writeToGitHubActionsSecret } from '../util/github/secret';
import { createRandomString } from '../util/random';

/**
 * Creates IAM for an AWS account.
 *
 * @param {string} repository the repository
 * @param {Output<string>} identityProviderArn the identity provider ARN
 * @returns {Output<string>} the IAM role ARN
 */
export const createAccountIam = (
  repository: Output<string>,
  identityProviderArn: string,
): Output<string> => {
  const labels = {
    ...commonLabels,
    organization: 'fh-burgenland-bswe',
    repository: repository,
  };

  const ciPostfix = repository.apply((repo) =>
    createRandomString(`aws-iam-role-ci-${repo}`, {}),
  );
  const truncatedRepository = repository.apply((repo) => repo.substring(0, 18));

  const ciRole = repository.apply(
    (repo) =>
      new aws.iam.Role(
        `aws-iam-role-ci-${repo}`,
        {
          name: interpolate`ci-${truncatedRepository}-${ciPostfix.result}`,
          description: `FH Burgenland Softwaremanagement II GitHub Repository: ${repo}`,
          assumeRolePolicy: JSON.stringify({
            Version: '2012-10-17',
            Statement: [
              {
                Action: 'sts:AssumeRoleWithWebIdentity',
                Effect: 'Allow',
                Principal: {
                  Federated: identityProviderArn,
                },
                Condition: {
                  StringEquals: {
                    'token.actions.githubusercontent.com:aud':
                      'sts.amazonaws.com',
                  },
                  StringLike: {
                    'token.actions.githubusercontent.com:sub': `repo:${githubOrganisation}/${repo}:*`,
                  },
                },
              },
            ],
          }),
          tags: labels,
        },
        {},
      ),
  );

  const policy = repository.apply(
    (repo) =>
      new aws.iam.Policy(
        `aws-iam-role-ci-policy-${repo}`,
        {
          name: interpolate`ci-${truncatedRepository}-${ciPostfix.result}`,
          description: `FH Burgenland Softwaremanagement II GitHub Repository: ${repo}`,
          policy: aws.iam
            .getPolicyDocument({
              statements: [
                {
                  effect: 'Allow',
                  actions: ['s3:*'],
                  resources: [
                    `arn:aws:s3:::bswe-${globalName}-${environment}-*`,
                    `arn:aws:s3:::bswe-${globalName}-${environment}-*/*`,
                  ],
                },
                {
                  effect: 'Allow',
                  actions: ['cloudfront:*'],
                  resources: [
                    `arn:aws:cloudfront::${awsAccountId}:distribution/*`,
                    `arn:aws:cloudfront::${awsAccountId}:origin-access-identity/*`,
                    `arn:aws:cloudfront::${awsAccountId}:origin-request-policy/*`,
                    `arn:aws:cloudfront::${awsAccountId}:response-headers-policy/*`,
                    `arn:aws:cloudfront::${awsAccountId}:origin-access-control/*`,
                  ],
                },
              ],
            })
            .then((doc) => doc.json),
          tags: labels,
        },
        {},
      ),
  );

  repository.apply(
    (repo) =>
      new aws.iam.RolePolicyAttachment(
        `aws-iam-role-ci-policy-attachment-${repo}`,
        {
          role: ciRole.name,
          policyArn: policy.arn,
        },
        {
          dependsOn: [ciRole, policy],
        },
      ),
  );

  repository.apply((repo) => {
    writeToGitHubActionsSecret(repo, 'AWS_IDENTITY_ROLE_ARN', ciRole.arn);
    writeToGitHubActionsSecret(
      repo,
      'AWS_REGION',
      Output.create(awsDefaultRegion),
    );
  });

  return ciRole.arn;
};
