import * as aws from '@pulumi/aws';
import { Output } from '@pulumi/pulumi';

import { StringMap } from '../../model/map';
import {
  commonLabels,
  environment,
  globalName,
  repositories,
} from '../configuration';

/**
 * Creates all Terraform related infrastructure.
 *
 * @returns {StringMap<Output<string>>} the repositories and their Terraform backend buckets
 */
export const configureTerraform = (): StringMap<Output<string>> => {
  const buckets = Object.fromEntries(
    repositories
      .filter((repo) => repo.terraform)
      .map((repo) => [repo.name, configureRepository(repo.name)]),
  );

  return buckets;
};

/**
 * Configures a repository for Terraform.
 *
 * @param {string} repository the repository
 * @returns {Output<string>} the Terraform backend bucket
 */
const configureRepository = (repository: string): Output<string> => {
  const bucket = new aws.s3.Bucket(
    `aws-s3-bucket-terraform-${environment}-${repository}`,
    {
      bucketPrefix: `bswe-${globalName}-${environment}-${repository}`,
      tags: commonLabels,
    },
    {},
  );

  return bucket.bucket;
};
