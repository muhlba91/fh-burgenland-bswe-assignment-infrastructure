/* eslint-disable functional/no-let */

import { Output } from '@pulumi/pulumi';

import { configureAwsAccounts } from './lib/aws';
import { createRepositories } from './lib/github';
import { createTeams } from './lib/github/team';
import { createGroups } from './lib/harbor/group';
import { configureHarborAuth } from './lib/harbor/oidc';
import { createProjects } from './lib/harbor/project';
import { configureTerraform } from './lib/pulumi';
import { StringMap } from './model/map';

export = async () => {
  // github teams and repositories
  const githubTeams = createTeams();
  const githubRepositories = createRepositories(githubTeams);

  // harbor configuration
  configureHarborAuth();
  const harborGroups = createGroups(githubTeams);
  const harborProjects = createProjects(harborGroups);

  // terraform and aws integrations
  let terraform: StringMap<Output<string>> = {};
  let aws: StringMap<Output<string>> = {};
  terraform = configureTerraform();
  aws = await configureAwsAccounts(githubRepositories);

  return {
    aws: aws,
    terraform: terraform,
    teams: Object.values(githubTeams).map((team) => team.name),
    repositories: Object.values(githubRepositories).map((repo) => repo.name),
    harbor: {
      projects: Object.values(harborProjects).map((project) => project.name),
    },
  };
};
