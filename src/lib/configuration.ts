import * as github from '@pulumi/github';
import { getStack } from '@pulumi/pulumi';

import { parseDataFromFile } from './util/data';

export const environment = getStack();

export const githubOrganisation = github.config.owner;
const data = parseDataFromFile(`./assets/data_${environment}.yaml`);
export const repositories = data.repositories;
export const teams = data.teams;

export const githubHandle = 'muhlba91';

export const awsDefaultRegion = 'eu-west-1';
export const awsAccountId = '061039787254';

export const globalName = 'swm2';
export const commonLabels = {
  environment: environment,
  purpose: globalName,
};
