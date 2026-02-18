package repositories

import (
	"fmt"

	"github.com/pulumi/pulumi-github/sdk/v6/go/github"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"

	"github.com/muhlba91/fh-burgenland-bswe-assignment-infrastructure/pkg/lib/config"
	"github.com/muhlba91/fh-burgenland-bswe-assignment-infrastructure/pkg/model/config/repository"
)

// createAccess creates access permissions for the given repository based on the provided configuration.
// ctx: The Pulumi context for resource creation.
// repository: The configuration for the repository.
// repo: The Pulumi GitHub repository resource.
// githubTeams: A map of GitHub teams for potential repository team associations.
func createAccess(
	ctx *pulumi.Context,
	repository *repository.Config,
	repo *github.Repository,
	githubTeams map[string]*github.Team,
) error {
	_, collErr := github.NewRepositoryCollaborator(
		ctx,
		fmt.Sprintf("github-repository-admin-%s", repository.Name),
		&github.RepositoryCollaboratorArgs{
			Repository: repo.Name,
			Username:   pulumi.String(config.OwnerHandle),
			Permission: pulumi.String("admin"),
		},
	)
	if collErr != nil {
		return collErr
	}

	for _, team := range repository.Teams {
		role := repositoryRoleToGitHubRole(team.Role)
		if len(role) == 0 {
			continue
		}

		ghTeam, exists := githubTeams[team.Name]
		if !exists {
			return fmt.Errorf("team %s not found for repository %s", team.Name, repository.Name)
		}

		_ = repo.Name.ApplyT(func(name string) error {
			_, tcErr := github.NewTeamRepository(
				ctx,
				fmt.Sprintf("github-team-repository-%s-%s", name, team.Name),
				&github.TeamRepositoryArgs{
					Repository: repo.Name,
					TeamId:     ghTeam.ID(),
					Permission: pulumi.String(role),
				},
				pulumi.DependsOn([]pulumi.Resource{repo}),
				pulumi.RetainOnDelete(true),
			)
			return tcErr
		})
	}

	return nil
}

// repositoryRoleToGitHubRole maps a custom repository role to a GitHub permission string.
// role: The custom repository role.
func repositoryRoleToGitHubRole(role string) string {
	switch role {
	case "owner":
		return "maintain"
	case "developer":
		return "maintain"
	default:
		return ""
	}
}
