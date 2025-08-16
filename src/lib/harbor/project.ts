import * as github from '@pulumi/github';
import { Output } from '@pulumi/pulumi';
import * as harbor from '@pulumiverse/harbor';

import { RepositoryTeamConfig } from '../../model/config/repository';
import { StringMap } from '../../model/map';
import { repositories } from '../configuration';

/**
 * Creates all Harbor projects.
 *
 * @param {StringMap<github.Repository>} githubRepositories the GitHub repositories
 * @param {StringMap<harbor.Group>} groups the Harbor groups
 * @returns {StringMap<Output<harbor.Project>>} the configured projects
 */
export const createProjects = (
  githubRepositories: StringMap<github.Repository>,
  groups: StringMap<harbor.Group>,
): StringMap<Output<harbor.Project>> =>
  Object.fromEntries(
    repositories
      .filter((repo) => repo.harbor)
      .map((repo) => [
        repo.name,
        createProject(githubRepositories[repo.name].name, repo.teams, groups),
      ]),
  );

/**
 * Creates a Harbor project.
 *
 * @param {Output<string>} repository the GitHub repository
 * @param {readonly RepositoryTeamConfig[]} teams the team names
 * @param {StringMap<harbor.Group>} groups the Harbor groups
 * @returns {Output<harbor.Project>} the project
 */
const createProject = (
  repository: Output<string>,
  teams: readonly RepositoryTeamConfig[],
  groups: StringMap<harbor.Group>,
): Output<harbor.Project> => {
  const harborProject = repository.apply(
    (repo) =>
      new harbor.Project(`harbor-project-${repo}`, {
        name: repo,
        public: false,
        storageQuota: 10,
        autoSbomGeneration: true,
        vulnerabilityScanning: true,
        forceDestroy: true,
      }),
  );

  repository.apply((repo) =>
    teams.forEach(
      (team) =>
        new harbor.ProjectMemberGroup(
          `harbor-project-member-group-${repo}-${team}`,
          {
            projectId: harborProject.id,
            type: 'oidc',
            role: repositoryRoleToHarborRole(team.role),
            groupName: groups[team.name].groupName,
          },
        ),
    ),
  );

  // TODO: pull-only access for other teams

  return harborProject;
};

/**
 * Maps a repository role to a Harbor role.
 *
 * @param role the repository role
 * @returns the corresponding Harbor role
 */
const repositoryRoleToHarborRole = (role: string): string => {
  switch (role) {
    case 'developer':
      return 'maintainer';
    default:
      return 'limitedGuest';
  }
};
