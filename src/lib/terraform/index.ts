import { Output } from '@pulumi/pulumi';
import * as pulumiservice from '@pulumi/pulumiservice';

import { StringMap } from '../../model/map';
import { environment, repositories } from '../configuration';

/**
 * Creates all Pulumi related infrastructure.
 *
 * @returns {StringMap<Output<string>>} the repositories and their Pulumi access tokens
 */
export const configurePulumi = (): StringMap<Output<string>> => {
  const repos = repositories
    .filter((repo) => repo.pulumi)
    .map((repo) => repo.name);

  const accessTokens = Object.fromEntries(
    repos.map((repository) => [repository, configureRepository(repository)]),
  );

  return accessTokens;
};

/**
 * Configures a repository for Pulumi.
 *
 * @param {string} repository the repository
 * @returns {Output<string>} the Pulumi access token
 */
const configureRepository = (repository: string): Output<string> => {
  const accessToken = new pulumiservice.AccessToken(
    `pulumi-access-token-${environment}-${repository}`,
    {
      description: `FH Burgenland: BSWE assignment ${environment} repository: ${repository}`,
    },
    {},
  );

  return accessToken.value;
};
