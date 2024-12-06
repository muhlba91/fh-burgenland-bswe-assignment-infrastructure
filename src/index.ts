/* eslint-disable functional/no-let */

import { Output } from '@pulumi/pulumi';

import { configureAwsAccounts } from './lib/aws';
import { createRepositories } from './lib/github';
import { createTeams } from './lib/github/team';
import { configureTerraform } from './lib/pulumi';
import { StringMap } from './model/map';

export = async () => {
  const githubTeams = createTeams();
  const githubRepositories = createRepositories(githubTeams);

  let terraform: StringMap<Output<string>> = {};
  let aws: StringMap<Output<string>> = {};

  terraform = configureTerraform();
  aws = await configureAwsAccounts(githubRepositories);

  return {
    aws: aws,
    terraform: terraform,
    teams: Object.values(githubTeams).map((team) => team.name),
    repositories: Object.values(githubRepositories).map((repo) => repo.name),
  };
};
