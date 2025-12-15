package project

import (
	"fmt"

	"github.com/muhlba91/fh-burgenland-bswe-assignment-infrastructure/pkg/model/config/repository"
	"github.com/muhlba91/pulumi-shared-library/pkg/util/defaults"
	"github.com/pulumi/pulumi-github/sdk/v6/go/github"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumiverse/pulumi-harbor/sdk/v3/go/harbor"
)

// defaultStorageQuota defines the default storage quota for Harbor projects in GB.
const defaultStorageQuota = 10

// Create creates Harbor projects based on the provided repositories and groups.
// ctx: pulumi.Context.
// repositories: list of repository configurations.
// githubRepositories: map of created GitHub repositories.
// groups: map of created Harbor groups.
func Create(
	ctx *pulumi.Context,
	repositories []*repository.Config,
	githubRepositories map[string]*github.Repository,
	groups map[string]*harbor.Group,
) map[string]*harbor.ProjectOutput {
	harborProjects := make(map[string]*harbor.ProjectOutput)

	for _, repoConfig := range repositories {
		if !defaults.GetOrDefault(repoConfig.Harbor, false) {
			continue
		}

		githubRepo, ok := githubRepositories[repoConfig.Name]
		if !ok {
			return nil
		}

		project, _ := githubRepo.Name.ApplyT(func(name string) *harbor.Project {
			project, _ := harbor.NewProject(ctx, fmt.Sprintf("harbor-project-%s", name), &harbor.ProjectArgs{
				Name:                  pulumi.String(name),
				Public:                pulumi.Bool(false),
				StorageQuota:          pulumi.Int(defaultStorageQuota),
				AutoSbomGeneration:    pulumi.Bool(true),
				VulnerabilityScanning: pulumi.Bool(true),
				ForceDestroy:          pulumi.Bool(true),
			})

			createPolicies(ctx, project)

			for _, team := range repoConfig.Teams {
				_, _ = harbor.NewProjectMemberGroup(
					ctx,
					fmt.Sprintf("harbor-project-member-group-%s-%s", name, team.Name),
					&harbor.ProjectMemberGroupArgs{
						ProjectId: project.ID(),
						Type:      pulumi.String("oidc"),
						Role:      pulumi.String(repositoryRoleToHarborRole(team.Role)),
						GroupName: groups[team.Name].GroupName,
					},
				)
			}

			return project
		}).(harbor.ProjectOutput)
		harborProjects[repoConfig.Name] = &project
	}

	return harborProjects
}

// repositoryRoleToHarborRole maps repository roles to Harbor roles.
// role: repository role as string.
func repositoryRoleToHarborRole(role string) string {
	switch role {
	case "developer":
		return "maintainer"
	default:
		return "limitedguest"
	}
}
