package main

import (
	"github.com/muhlba91/fh-burgenland-bswe-assignment-infrastructure/pkg/lib/aws"
	"github.com/muhlba91/fh-burgenland-bswe-assignment-infrastructure/pkg/lib/config"
	"github.com/muhlba91/fh-burgenland-bswe-assignment-infrastructure/pkg/lib/github/repositories"
	ghTeam "github.com/muhlba91/fh-burgenland-bswe-assignment-infrastructure/pkg/lib/github/team"
	glTeam "github.com/muhlba91/fh-burgenland-bswe-assignment-infrastructure/pkg/lib/gitlab/team"
	harborCfg "github.com/muhlba91/fh-burgenland-bswe-assignment-infrastructure/pkg/lib/harbor"
	"github.com/muhlba91/fh-burgenland-bswe-assignment-infrastructure/pkg/lib/terraform"
	"github.com/muhlba91/fh-burgenland-bswe-assignment-infrastructure/pkg/util/export"
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
		githubRepositories, ghrErr := repositories.Create(ctx, stackConfig.Repositories, githubTeams)
		if ghrErr != nil {
			return ghrErr
		}
		unknownRepositories := make([]string, 0)
		for _, repo := range stackConfig.Repositories {
			_, ghExists := githubRepositories[repo.Name]
			if !ghExists {
				log.Warn().
					Msgf("[repository] repository %s does not have a known provider and was not created", repo.Name)
				unknownRepositories = append(unknownRepositories, repo.Name)
			}
		}

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
		export.GitLab(ctx, gitlabTeams, githubRepositories)
		export.Harbor(ctx, harborProjects, harborRobotAccounts)
		export.Terraform(ctx, terraform)
		export.AWS(ctx, aws)
		ctx.Export("unknownRepositories", pulumi.ToStringArray(unknownRepositories))

		return nil
	})
}
