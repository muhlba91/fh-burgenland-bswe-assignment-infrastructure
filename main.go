package main

import (
	"github.com/muhlba91/fh-burgenland-bswe-assignment-infrastructure/pkg/lib/aws"
	"github.com/muhlba91/fh-burgenland-bswe-assignment-infrastructure/pkg/lib/config"
	ghRepositories "github.com/muhlba91/fh-burgenland-bswe-assignment-infrastructure/pkg/lib/github/repositories"
	ghTeam "github.com/muhlba91/fh-burgenland-bswe-assignment-infrastructure/pkg/lib/github/team"
	glRepositories "github.com/muhlba91/fh-burgenland-bswe-assignment-infrastructure/pkg/lib/gitlab/repositories"
	glTeam "github.com/muhlba91/fh-burgenland-bswe-assignment-infrastructure/pkg/lib/gitlab/team"
	harborCfg "github.com/muhlba91/fh-burgenland-bswe-assignment-infrastructure/pkg/lib/harbor"
	"github.com/muhlba91/fh-burgenland-bswe-assignment-infrastructure/pkg/lib/terraform"
	"github.com/muhlba91/fh-burgenland-bswe-assignment-infrastructure/pkg/model/config/repository"
	"github.com/muhlba91/fh-burgenland-bswe-assignment-infrastructure/pkg/util/export"
	"github.com/pulumi/pulumi-github/sdk/v6/go/github"
	"github.com/pulumi/pulumi-gitlab/sdk/v9/go/gitlab"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/rs/zerolog/log"
)

// main is the entry point of the Pulumi program.
func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		stackConfig, err := config.LoadConfig(ctx)
		if err != nil {
			return err
		}

		// teams
		githubTeams, ghtErr := ghTeam.Create(ctx, stackConfig.Teams)
		if ghtErr != nil {
			return ghtErr
		}
		gitlabTeams, gltErr := glTeam.Create(ctx, stackConfig.Teams)
		if gltErr != nil {
			return gltErr
		}

		// repositories
		githubRepositories, ghrErr := ghRepositories.Create(ctx, stackConfig.Repositories, githubTeams)
		if ghrErr != nil {
			return ghrErr
		}
		gitlabRepositories, glrErr := glRepositories.Create(ctx, stackConfig.Repositories, gitlabTeams)
		if glrErr != nil {
			return glrErr
		}
		unknownRepositories := detectUnknownRepositories(
			stackConfig.Repositories,
			githubRepositories,
			gitlabRepositories,
		)

		// harbor
		harborProjects, harborRobotAccounts, haErr := harborCfg.Configure(
			ctx,
			stackConfig,
			githubRepositories,
			githubTeams,
		)
		if haErr != nil {
			return haErr
		}

		// terraform integration
		terraform, tfErr := terraform.Configure(ctx, stackConfig.Repositories, githubRepositories)
		if tfErr != nil {
			return tfErr
		}

		// aws integration
		aws, awsErr := aws.Configure(ctx, stackConfig.Repositories, githubRepositories)
		if awsErr != nil {
			return awsErr
		}

		export.GitHub(ctx, githubTeams, githubRepositories)
		export.GitLab(ctx, gitlabTeams, gitlabRepositories)
		export.Harbor(ctx, harborProjects, harborRobotAccounts)
		export.Terraform(ctx, terraform)
		export.AWS(ctx, aws)
		ctx.Export("unknownRepositories", pulumi.ToStringArray(unknownRepositories))

		return nil
	})
}

// detectUnknownRepositories checks if there are any repositories in the configuration that do not have a known provider and returns their names.
// repositories: the list of repositories defined in the configuration.
// githubRepositories: the map of repositories created in GitHub, keyed by repository name.
// gitlabRepositories: the map of repositories created in GitLab, keyed by repository name.
func detectUnknownRepositories(
	repositories []*repository.Config,
	githubRepositories map[string]*github.Repository,
	gitlabRepositories map[string]*gitlab.Project,
) []string {
	unknownRepositories := make([]string, 0)

	for _, repo := range repositories {
		_, ghExists := githubRepositories[repo.Name]
		_, glExists := gitlabRepositories[repo.Name]

		if !ghExists && !glExists {
			log.Warn().
				Msgf("[repository] repository %s does not have a known provider and was not created", repo.Name)
			unknownRepositories = append(unknownRepositories, repo.Name)
		}
	}

	return unknownRepositories
}
