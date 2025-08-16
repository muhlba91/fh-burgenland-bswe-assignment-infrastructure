/**
 * Defines repository config.
 */
export interface RepositoryConfig {
  readonly name: string;
  readonly service: string;
  readonly teams: readonly string[];
  readonly approvers?: number;
  readonly deleteOnDestroy?: boolean;
  readonly aws: boolean;
  readonly terraform: boolean;
  readonly pulumi: boolean;
  readonly harbor: boolean;
  readonly requiredChecks: readonly string[];
}
