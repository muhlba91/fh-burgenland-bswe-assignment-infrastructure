import * as github from '@pulumi/github';
import { Output } from '@pulumi/pulumi';

import { RepositoryConfig } from '../../model/config/repository';
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
    createTeamAccess(githubRepository, team, githubTeams[team]?.id);
  });
};

/**
 * Creates access permissions for the team to the repository.
 *
 * @param {github.Repository} repository the repository
 * @param {string} team the team name
 * @param {number} teamId the team id
 */
const createTeamAccess = (
  repository: github.Repository,
  team: string,
  teamId: Output<string>,
) =>
  repository.name.apply(
    (repositoryName) =>
      new github.TeamRepository(
        `github-team-repository-${repositoryName}-${team}`,
        {
          repository: repositoryName,
          teamId: teamId,
          permission: 'maintain',
        },
        {
          dependsOn: [repository],
          retainOnDelete: true,
        },
      ),
  );
