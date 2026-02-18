package repositories

import (
	"fmt"
	"strconv"

	"github.com/pulumi/pulumi-gitlab/sdk/v9/go/gitlab"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"

	"github.com/muhlba91/fh-burgenland-bswe-assignment-infrastructure/pkg/model/config/repository"
)

// createAccess creates access permissions for the given repository based on the provided configuration.
// ctx: The Pulumi context for resource creation.
// repository: The configuration for the repository.
// repo: The Pulumi GitLab repository resource.
// gitlabTeams: A map of GitLab teams for potential repository team associations.
func createAccess(
	ctx *pulumi.Context,
	repository *repository.Config,
	repo *gitlab.Project,
	gitlabTeams map[string]*gitlab.Group,
) error {
	for _, team := range repository.Teams {
		role := repositoryRoleToGitLabRole(team.Role)
		if len(role) == 0 {
			continue
		}

		glTeam, exists := gitlabTeams[team.Name]
		if !exists {
			return fmt.Errorf("team %s not found for repository %s", team.Name, repository.Name)
		}
		teamID, _ := glTeam.ID().ToStringOutput().ApplyT(func(id string) int {
			gid, _ := strconv.Atoi(id)
			return gid
		}).(pulumi.IntOutput)

		_ = repo.Name.ApplyT(func(name string) error {
			_, tcErr := gitlab.NewProjectShareGroup(
				ctx,
				fmt.Sprintf("gitlab-team-repository-%s-%s", name, team.Name),
				&gitlab.ProjectShareGroupArgs{
					Project:     repo.ID(),
					GroupId:     teamID,
					GroupAccess: pulumi.String(role),
				},
				pulumi.DependsOn([]pulumi.Resource{repo}),
				pulumi.RetainOnDelete(true),
			)
			return tcErr
		})
	}

	return nil
}

// repositoryRoleToGitLabRole maps a custom repository role to a GitLab permission string.
// role: The custom repository role.
func repositoryRoleToGitLabRole(role string) string {
	switch role {
	case "developer":
		return "developer"
	default:
		return ""
	}
}
