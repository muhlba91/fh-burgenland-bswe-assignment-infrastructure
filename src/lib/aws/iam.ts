import * as aws from '@pulumi/aws';
import { interpolate, Output } from '@pulumi/pulumi';

import {
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
                  actions: [
                    's3:AbortMultipartUpload',
                    's3:CreateBucket',
                    's3:DeleteBucket',
                    's3:DeleteBucketPolicy',
                    's3:DeleteBucketWebsite',
                    's3:DeleteObject',
                    's3:DeleteObjectTagging',
                    's3:DeleteObjectVersion',
                    's3:DeleteObjectVersionTagging',
                    's3:GetAccountPublicAccessBlock',
                    's3:GetBucketAcl',
                    's3:GetBucketCORS',
                    's3:GetBucketLocation',
                    's3:GetBucketObjectLockConfiguration',
                    's3:GetBucketPolicy',
                    's3:GetBucketPolicyStatus',
                    's3:GetBucketPublicAccessBlock',
                    's3:GetBucketTagging',
                    's3:GetBucketVersioning',
                    's3:GetBucketWebsite',
                    's3:GetEncryptionConfiguration',
                    's3:GetLifecycleConfiguration',
                    's3:GetObject',
                    's3:GetObjectAcl',
                    's3:GetObjectAttributes',
                    's3:GetObjectRetention',
                    's3:GetObjectTagging',
                    's3:GetObjectVersion',
                    's3:GetObjectVersionAcl',
                    's3:GetObjectVersionAttributes',
                    's3:GetObjectVersionTagging',
                    's3:ListAllMyBuckets',
                    's3:ListBucket',
                    's3:ListBucketMultipartUploads',
                    's3:ListBucketVersions',
                    's3:ListTagsForResource',
                    's3:ObjectOwnerOverrideToBucketOwner',
                    's3:PutBucketAcl',
                    's3:PutBucketCORS',
                    's3:PutBucketObjectLockConfiguration',
                    's3:PutBucketPolicy',
                    's3:PutBucketPublicAccessBlock',
                    's3:PutBucketTagging',
                    's3:PutBucketVersioning',
                    's3:PutBucketWebsite',
                    's3:PutEncryptionConfiguration',
                    's3:PutObject',
                    's3:PutObjectAcl',
                    's3:PutObjectRetention',
                    's3:PutObjectTagging',
                    's3:PutObjectVersionAcl',
                    's3:PutObjectVersionTagging',
                    's3:TagResource',
                    's3:UntagResource',
                  ],
                  resources: [
                    `arn:aws:s3:::bswe-${globalName}-${environment}-*`,
                    `arn:aws:s3:::bswe-${globalName}-${environment}-*/*`,
                  ],
                },
                {
                  effect: 'Allow',
                  actions: [
                    'cloudfront:CreateDistribution',
                    'cloudfront:CreateInvalidation',
                    'cloudfront:CreateOriginAccessControl',
                    'cloudfront:CreateOriginRequestPolicy',
                    'cloudfront:DeleteDistribution',
                    'cloudfront:DeleteOriginAccessControl',
                    'cloudfront:DeleteOriginRequestPolicy',
                    'cloudfront:GetCloudFrontOriginAccessIdentity',
                    'cloudfront:GetCloudFrontOriginAccessIdentityConfig',
                    'cloudfront:GetDistribution',
                    'cloudfront:GetDistributionConfig',
                    'cloudfront:GetInvalidation',
                    'cloudfront:GetOriginAccessControl',
                    'cloudfront:GetOriginAccessControlConfig',
                    'cloudfront:GetOriginRequestPolicy',
                    'cloudfront:GetOriginRequestPolicyConfig',
                    'cloudfront:GetResponseHeadersPolicy',
                    'cloudfront:GetResponseHeadersPolicyConfig',
                    'cloudfront:ListCloudFrontOriginAccessIdentities',
                    'cloudfront:ListDistributions',
                    'cloudfront:ListDistributionsByOriginRequestPolicyId',
                    'cloudfront:ListDistributionsByResponseHeadersPolicyId',
                    'cloudfront:ListInvalidations',
                    'cloudfront:ListOriginAccessControls',
                    'cloudfront:ListOriginRequestPolicies',
                    'cloudfront:ListResponseHeadersPolicies',
                    'cloudfront:ListTagsForResource',
                    'cloudfront:TagResource',
                    'cloudfront:UntagResource',
                    'cloudfront:UpdateCloudFrontOriginAccessIdentity',
                    'cloudfront:UpdateDistribution',
                    'cloudfront:UpdateDistributionWithStagingConfig',
                    'cloudfront:UpdateOriginAccessControl',
                    'cloudfront:UpdateOriginRequestPolicy',
                    'cloudfront:UpdateResponseHeadersPolicy',
                  ],
                  resources: [
                    `arn:aws:cloudfront:${awsDefaultRegion}::distribution/bswe-${globalName}-${environment}-*`,
                    `arn:aws:cloudfront:${awsDefaultRegion}::origin-access-identity/*`,
                    `arn:aws:cloudfront:${awsDefaultRegion}::origin-request-policy/*`,
                    `arn:aws:cloudfront:${awsDefaultRegion}::response-headers-policy/*`,
                    `arn:aws:cloudfront:${awsDefaultRegion}::origin-access-control/*`,
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
