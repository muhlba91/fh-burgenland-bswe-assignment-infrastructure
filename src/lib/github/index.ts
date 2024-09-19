import * as github from '@pulumi/github';

import { RepositoryConfig } from '../../model/config/repository';
import { StringMap } from '../../model/map';
import { environment, globalName, repositories } from '../configuration';

import { createRepositoryAccess } from './access';
import { createRepositoryRulesets } from './ruleset';

/**
 * Creates all GitHub repositories.
 *
 * @param {StringMap<github.Team>} githubTeams the GitHub teams
 * @returns {StringMap<github.Repository>} the configured repositories
 */
export const createRepositories = (
  githubTeams: StringMap<github.Team>,
): StringMap<github.Repository> =>
  Object.fromEntries(
    repositories.map((repository) => [
      repository.name,
      createRepository(repository, githubTeams),
    ]),
  );

/**
 * Creates a GitHub repository.
 *
 * @param {RepositoryConfig} repository the repository configuration
 * @param {StringMap<github.Team>} githubTeams the GitHub teams
 * @returns {github.Repository} the repository
 */
const createRepository = (
  repository: RepositoryConfig,
  githubTeams: StringMap<github.Team>,
): github.Repository => {
  const githubRepository = new github.Repository(
    `github-repo-${repository.name}`,
    {
      name: `${globalName}-${environment}-${repository.name}`,
      description: `Softwaremanagement II ${environment}: ${repository.service} repository`,
      hasDiscussions: false,
      hasWiki: true,
      topics: [globalName, environment, repository.service].sort(),
      visibility: 'private',
      allowAutoMerge: false,
      allowMergeCommit: false,
      allowRebaseMerge: true,
      allowSquashMerge: false,
      allowUpdateBranch: true,
      archived: false,
      archiveOnDestroy: !repository.deleteOnDestroy,
      autoInit: true,
      deleteBranchOnMerge: true,
      hasDownloads: true,
      hasIssues: true,
      hasProjects: true,
      mergeCommitMessage: 'PR_TITLE',
      mergeCommitTitle: 'MERGE_MESSAGE',
      squashMergeCommitMessage: 'COMMIT_MESSAGES',
      squashMergeCommitTitle: 'COMMIT_OR_PR_TITLE',
      vulnerabilityAlerts: true,
    },
    {
      retainOnDelete: !repository.deleteOnDestroy,
      ignoreChanges: ['securityAndAnalysis', 'template'],
    },
  );

  createRepositoryRulesets(repository, githubRepository);
  createRepositoryAccess(repository, githubRepository, githubTeams);

  return githubRepository;
};
