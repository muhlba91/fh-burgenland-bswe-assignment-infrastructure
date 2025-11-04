import * as github from '@pulumi/github';
import { Output } from '@pulumi/pulumi';

/**
 * Stores a value in GitHub Actions secrets.
 *
 * @param {string} repository the repository
 * @param {string} key the key
 * @param {Output<string>} value the value
 */
export const writeToGitHubActionsSecret = (
  repository: string,
  key: string,
  value: Output<string>,
) => {
  new github.ActionsSecret(
    `github-actions-secret-${repository}-${key}`,
    {
      repository: repository,
      secretName: key,
      plaintextValue: value,
      destroyOnDrift: false,
    },
    {},
  );
};
