import * as harbor from '@pulumiverse/harbor';

import { TeamConfig } from '../../model/config/team';
import { StringMap } from '../../model/map';
import { environment, globalName, teams } from '../configuration';

/**
 * Creates all Harbor projects.
 *
 * @param {StringMap<harbor.Group>} groups the Harbor groups
 * @returns {StringMap<harbor.Project>} the configured projects
 */
export const createProjects = (
  groups: StringMap<harbor.Group>,
): StringMap<harbor.Project> =>
  Object.fromEntries(
    teams.map((team) => [team.name, createProject(team, groups)]),
  );

/**
 * Creates a Harbor project.
 *
 * @param {TeamConfig} team the team configuration
 * @param {StringMap<harbor.Group>} groups the Harbor groups
 * @returns {harbor.Project} the project
 */
const createProject = (
  team: TeamConfig,
  groups: StringMap<harbor.Group>,
): harbor.Project => {
  const harborProject = new harbor.Project(`harbor-project-${team.name}`, {
    name: `${globalName}-${environment}-${team.name}`,
    public: false,
    storageQuota: 15,
    autoSbomGeneration: true,
    vulnerabilityScanning: true,
    forceDestroy: true,
  });

  new harbor.ProjectMemberGroup(`harbor-project-member-group-${team.name}`, {
    projectId: harborProject.id,
    type: 'oidc',
    role: 'maintainer',
    groupName: groups[team.name].groupName,
  });
  // TODO: add joining group

  return harborProject;
};
