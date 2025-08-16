/**
 * Defines repository config.
 */
export interface RepositoryConfig {
  readonly name: string;
  readonly service: string;
  readonly teams: readonly RepositoryTeamConfig[];
  readonly approvers?: number;
  readonly deleteOnDestroy?: boolean;
  readonly aws: boolean;
  readonly terraform: boolean;
  readonly pulumi: boolean;
  readonly harbor: boolean;
  readonly requiredChecks: readonly string[];
}

/**
 * Defines ra epository team config.
 */
export interface RepositoryTeamConfig {
  readonly name: string;
  readonly role: string;
}
