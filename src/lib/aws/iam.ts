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
                  actions: [
                    's3:AbortMultipartUpload',
                    's3:CompleteMultipartUpload',
                    's3:CopyObject',
                    's3:CreateBucket',
                    's3:CreateBucketMetadataTableConfiguration',
                    's3:CreateMultipartUpload',
                    's3:CreateSession',
                    's3:DeleteBucket',
                    's3:DeleteBucketAnalyticsConfiguration',
                    's3:DeleteBucketCors',
                    's3:DeleteBucketEncryption',
                    's3:DeleteBucketIntelligentTieringConfiguration',
                    's3:DeleteBucketInventoryConfiguration',
                    's3:DeleteBucketLifecycle',
                    's3:DeleteBucketMetadataTableConfiguration',
                    's3:DeleteBucketMetricsConfiguration',
                    's3:DeleteBucketOwnershipControls',
                    's3:DeleteBucketPolicy',
                    's3:DeleteBucketReplication',
                    's3:DeleteBucketTagging',
                    's3:DeleteBucketWebsite',
                    's3:DeleteObject',
                    's3:DeleteObjects',
                    's3:DeleteObjectTagging',
                    's3:DeletePublicAccessBlock',
                    's3:GetBucketAccelerateConfiguration',
                    's3:GetBucketAcl',
                    's3:GetBucketAnalyticsConfiguration',
                    's3:GetBucketCors',
                    's3:GetBucketEncryption',
                    's3:GetBucketIntelligentTieringConfiguration',
                    's3:GetBucketInventoryConfiguration',
                    's3:GetBucketLifecycle',
                    's3:GetBucketLifecycleConfiguration',
                    's3:GetBucketLocation',
                    's3:GetBucketLogging',
                    's3:GetBucketMetadataTableConfiguration',
                    's3:GetBucketMetricsConfiguration',
                    's3:GetBucketNotification',
                    's3:GetBucketNotificationConfiguration',
                    's3:GetBucketOwnershipControls',
                    's3:GetBucketPolicy',
                    's3:GetBucketPolicyStatus',
                    's3:GetBucketReplication',
                    's3:GetBucketRequestPayment',
                    's3:GetBucketTagging',
                    's3:GetBucketVersioning',
                    's3:GetBucketWebsite',
                    's3:GetObject',
                    's3:GetObjectAcl',
                    's3:GetObjectAttributes',
                    's3:GetObjectLegalHold',
                    's3:GetObjectLockConfiguration',
                    's3:GetObjectRetention',
                    's3:GetObjectTagging',
                    's3:GetObjectTorrent',
                    's3:GetPublicAccessBlock',
                    's3:HeadBucket',
                    's3:HeadObject',
                    's3:ListBucketAnalyticsConfigurations',
                    's3:ListBucketIntelligentTieringConfigurations',
                    's3:ListBucketInventoryConfigurations',
                    's3:ListBucketMetricsConfigurations',
                    's3:ListBuckets',
                    's3:ListBucket',
                    's3:ListDirectoryBuckets',
                    's3:ListMultipartUploads',
                    's3:ListObjects',
                    's3:ListObjectsV2',
                    's3:ListObjectVersions',
                    's3:ListParts',
                    's3:PutBucketAccelerateConfiguration',
                    's3:PutBucketAcl',
                    's3:PutBucketAnalyticsConfiguration',
                    's3:PutBucketCors',
                    's3:PutBucketEncryption',
                    's3:PutBucketIntelligentTieringConfiguration',
                    's3:PutBucketInventoryConfiguration',
                    's3:PutBucketLifecycle',
                    's3:PutBucketLifecycleConfiguration',
                    's3:PutBucketLogging',
                    's3:PutBucketMetricsConfiguration',
                    's3:PutBucketNotification',
                    's3:PutBucketNotificationConfiguration',
                    's3:PutBucketOwnershipControls',
                    's3:PutBucketPolicy',
                    's3:PutBucketReplication',
                    's3:PutBucketRequestPayment',
                    's3:PutBucketTagging',
                    's3:PutBucketVersioning',
                    's3:PutBucketWebsite',
                    's3:PutObject',
                    's3:PutObjectAcl',
                    's3:PutObjectLegalHold',
                    's3:PutObjectLockConfiguration',
                    's3:PutObjectRetention',
                    's3:PutObjectTagging',
                    's3:PutPublicAccessBlock',
                    's3:RestoreObject',
                    's3:SelectObjectContent',
                    's3:UploadPart',
                    's3:UploadPartCopy',
                    's3:WriteGetObjectResponse',
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
                    `arn:aws:cloudfront::${awsAccountId}:distribution/bswe-${globalName}-${environment}-*`,
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
