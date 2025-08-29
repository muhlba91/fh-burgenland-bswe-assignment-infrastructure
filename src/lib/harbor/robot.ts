import * as github from '@pulumi/github';
import { Output } from '@pulumi/pulumi';
import * as harbor from '@pulumiverse/harbor';

import { StringMap } from '../../model/map';
import { repositories } from '../configuration';
import { writeToGitHubActionsSecret } from '../util/github/secret';

/**
 * Configures Harbor robot accounts for CI/CD pipelines.
 *
 * @param {StringMap<Output<harbor.Project>>} projects the Harbor projects
 * @param {StringMap<github.Repository>} githubRepositories the GitHub repositories
 * @return {StringMap<Output<string>>} the robot accounts for each project
 */
export const configureHarborRobotAccounts = (
  projects: StringMap<Output<harbor.Project>>,
  githubRepositories: StringMap<github.Repository>,
): StringMap<Output<string>> => {
  const robots = Object.fromEntries(
    repositories
      .filter((repo) => repo.harbor)
      .map((repo) => [
        repo.name,
        createRobotAccount(
          githubRepositories[repo.name].name,
          projects[repo.name],
        ),
      ]),
  );

  return robots;
};

/**
 * Creates the Harbor robot accounts for a GitHub repository.
 *
 * @param {Output<string>} githubRepository the GitHub repository name
 * @param {Output<harbor.Project>} project the Harbor projects
 * @returns {Output<string>} the robot account name
 */
export const createRobotAccount = (
  githubRepository: Output<string>,
  project: Output<harbor.Project>,
): Output<string> => {
  const robot = githubRepository.apply(
    (repo) =>
      new harbor.RobotAccount(`harbor-robot-${repo}`, {
        level: 'project',
        description: `Robot account for ${repo}`,
        permissions: [
          {
            kind: 'project',
            namespace: project.name,
            accesses: [
              {
                action: 'push',
                resource: 'repository',
                effect: 'allow',
              },
              {
                action: 'pull',
                resource: 'repository',
                effect: 'allow',
              },
            ],
          },
        ],
      }),
  );

  githubRepository.apply((repo) => {
    writeToGitHubActionsSecret(
      repo,
      'HARBOR_REGISTRY_URL',
      Output.create(process.env.HARBOR_URL?.replace('https://', '') || ''),
    );
    writeToGitHubActionsSecret(repo, 'HARBOR_ROBOT_NAME', robot.fullName);
    writeToGitHubActionsSecret(
      repo,
      'HARBOR_ROBOT_SECRET',
      Output.create(robot.secret),
    );
  });

  return robot.fullName;
};
