import * as github from '@pulumi/github';
import { Output } from '@pulumi/pulumi';

import {
  RepositoryConfig,
  RepositoryTeamConfig,
} from '../../model/config/repository';
import { StringMap } from '../../model/map';
import { githubHandle } from '../configuration';

/**
 * Creates all access permissions for a repository.
 *
 * @param {RepositoryConfig} repository the repository configuration
 * @param {github.Repository} githubRepository the GitHub repository
 * @param {StringMap<github.Team>} githubTeams the GitHub teams
 */
export const createRepositoryAccess = (
  repository: RepositoryConfig,
  githubRepository: github.Repository,
  githubTeams: StringMap<github.Team>,
) => {
  new github.RepositoryCollaborator(
    `github-repository-admin-${repository.name}`,
    {
      repository: githubRepository.name,
      username: githubHandle,
      permission: 'admin',
    },
    {
      dependsOn: [githubRepository],
      retainOnDelete: true,
    },
  );

  repository.teams.forEach(async (team) => {
    createTeamAccess(githubRepository, team, githubTeams[team.name]?.id);
  });
};

/**
 * Creates access permissions for the team to the repository.
 *
 * @param {github.Repository} repository the repository
 * @param {RepositoryTeamConfig} team the team configuration
 * @param {number} teamId the team id
 */
const createTeamAccess = (
  repository: github.Repository,
  team: RepositoryTeamConfig,
  teamId: Output<string>,
) =>
  repository.name.apply(
    (repositoryName) =>
      new github.TeamRepository(
        `github-team-repository-${repositoryName}-${team.name}`,
        {
          repository: repositoryName,
          teamId: teamId,
          permission: repositoryRoleToGitHubRole(team.role),
        },
        {
          dependsOn: [repository],
          retainOnDelete: true,
        },
      ),
  );

/**
 * Maps a repository role to a Harbor role.
 *
 * @param role the repository role
 * @returns the corresponding Harbor role
 */
const repositoryRoleToGitHubRole = (role: string): string => {
  switch (role) {
    case 'developer':
      return 'maintain';
    default:
      return '';
  }
};
