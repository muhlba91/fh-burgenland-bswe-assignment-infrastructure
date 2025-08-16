import * as github from '@pulumi/github';
import { interpolate } from '@pulumi/pulumi';
import * as harbor from '@pulumiverse/harbor';

import { StringMap } from '../../model/map';
import { githubOrganisation, teams } from '../configuration';

/**
 * Creates all Harbor groups.
 *
 * @param {StringMap<github.Team>} githubTeams the GitHub teams
 * @returns {StringMap<harbor.Group>} the configured groups
 */
export const createGroups = (
  githubTeams: StringMap<github.Team>,
): StringMap<harbor.Group> =>
  Object.fromEntries(
    teams.map((team) => [
      team.name,
      new harbor.Group(`harbor-group-${team.name}`, {
        groupName: interpolate`${githubOrganisation}:${githubTeams[team.name].name}`,
        groupType: 3,
      }),
    ]),
  );
