package robot

import (
	"fmt"
	"os"
	"strings"

	"github.com/muhlba91/fh-burgenland-bswe-assignment-infrastructure/pkg/model/config/repository"
	"github.com/muhlba91/fh-burgenland-bswe-assignment-infrastructure/pkg/util/secret"
	"github.com/muhlba91/pulumi-shared-library/pkg/util/defaults"
	"github.com/pulumi/pulumi-github/sdk/v6/go/github"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumiverse/pulumi-harbor/sdk/v3/go/harbor"
)

// Create configures Harbor robot accounts for the given repositories.
// ctx: The Pulumi context.
// repositories: The list of repository configurations.
// githubRepositories: A map of GitHub repository names to their corresponding Pulumi GitHub Repository resources.
// projects: A map of Harbor project names to their corresponding Pulumi Harbor Project resources.
func Create(
	ctx *pulumi.Context,
	repositories []*repository.Config,
	githubRepositories map[string]*github.Repository,
	projects map[string]*harbor.ProjectOutput,
) (map[string]*pulumi.StringOutput, error) {
	harborRobotAccounts := make(map[string]*pulumi.StringOutput)

	for _, repoConfig := range repositories {
		if !defaults.GetOrDefault(repoConfig.Harbor, false) {
			continue
		}

		project, exists := projects[repoConfig.Name]
		if !exists {
			return nil, fmt.Errorf("harbor project %s not found for harbor robot account creation", repoConfig.Name)
		}

		ghRepo, exists := githubRepositories[repoConfig.Name]
		if !exists {
			return nil, fmt.Errorf("github repository %s not found for harbor robot account creation", repoConfig.Name)
		}

		robot := ghRepo.Name.ApplyT(func(name string) *harbor.RobotAccount {
			projectName, _ := project.ApplyT(func(p *harbor.Project) pulumi.StringOutput {
				return p.Name
			}).(pulumi.StringOutput)
			ra, _ := harbor.NewRobotAccount(ctx, fmt.Sprintf("harbor-robot-%s", name), &harbor.RobotAccountArgs{
				Level:       pulumi.String("project"),
				Description: pulumi.String(fmt.Sprintf("Robot account for %s", name)),
				Permissions: &harbor.RobotAccountPermissionArray{
					&harbor.RobotAccountPermissionArgs{
						Kind:      pulumi.String("project"),
						Namespace: projectName,
						Accesses: &harbor.RobotAccountPermissionAccessArray{
							&harbor.RobotAccountPermissionAccessArgs{
								Action:   pulumi.String("push"),
								Resource: pulumi.String("repository"),
								Effect:   pulumi.String("allow"),
							},
							&harbor.RobotAccountPermissionAccessArgs{
								Action:   pulumi.String("pull"),
								Resource: pulumi.String("repository"),
								Effect:   pulumi.String("allow"),
							},
						},
					},
				},
			})
			return ra
		})

		robotAccount, _ := robot.ApplyT(func(r *harbor.RobotAccount) pulumi.StringOutput {
			return r.FullName
		}).(pulumi.StringOutput)
		harborRobotAccounts[repoConfig.Name] = &robotAccount

		robotName, _ := robot.ApplyT(func(r *harbor.RobotAccount) pulumi.StringOutput { return r.FullName }).(pulumi.StringOutput)
		robotSecret, _ := robot.ApplyT(func(r *harbor.RobotAccount) pulumi.StringOutput { return r.Secret }).(pulumi.StringOutput)

		_ = secret.Write(
			ctx,
			repoConfig,
			githubRepositories,
			"HARBOR_REGISTRY_URL",
			pulumi.String(strings.ReplaceAll(os.Getenv("HARBOR_URL"), "https://", "")),
		)
		_ = secret.Write(ctx, repoConfig, githubRepositories, "HARBOR_ROBOT_NAME", robotName)
		_ = secret.Write(ctx, repoConfig, githubRepositories, "HARBOR_ROBOT_SECRET", robotSecret)
	}

	return harborRobotAccounts, nil
}
