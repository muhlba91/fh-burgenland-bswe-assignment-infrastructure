package main

import (
	"maps"
	"slices"

	"github.com/muhlba91/fh-burgenland-bswe-assignment-infrastructure/pkg/lib/aws"
	"github.com/muhlba91/fh-burgenland-bswe-assignment-infrastructure/pkg/lib/config"
	"github.com/muhlba91/fh-burgenland-bswe-assignment-infrastructure/pkg/lib/github/repositories"
	"github.com/muhlba91/fh-burgenland-bswe-assignment-infrastructure/pkg/lib/github/team"
	"github.com/muhlba91/fh-burgenland-bswe-assignment-infrastructure/pkg/lib/harbor/auth"
	"github.com/muhlba91/fh-burgenland-bswe-assignment-infrastructure/pkg/lib/harbor/group"
	"github.com/muhlba91/fh-burgenland-bswe-assignment-infrastructure/pkg/lib/harbor/project"
	"github.com/muhlba91/fh-burgenland-bswe-assignment-infrastructure/pkg/lib/harbor/robot"
	"github.com/muhlba91/fh-burgenland-bswe-assignment-infrastructure/pkg/lib/terraform"
	"github.com/pulumi/pulumi-github/sdk/v6/go/github"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumiverse/pulumi-harbor/sdk/v3/go/harbor"
)

// main is the entry point of the Pulumi program.
func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		stackConfig, err := config.LoadConfig(ctx)
		if err != nil {
			return err
		}

		// github teams and repositories
		githubTeams, ghtErr := team.Create(ctx, stackConfig.Teams)
		if ghtErr != nil {
			return ghtErr
		}
		githubRepositories, ghrErr := repositories.Create(ctx, stackConfig.Repositories, githubTeams)
		if ghrErr != nil {
			return ghrErr
		}

		// harbor configuration
		haErr := auth.Configure(ctx)
		if haErr != nil {
			return haErr
		}
		harborGroups, hgErr := group.Create(ctx, githubTeams)
		if hgErr != nil {
			return hgErr
		}
		harborProjects := project.Create(ctx, stackConfig.Repositories, githubRepositories, harborGroups)
		harborRobotAccounts, hrErr := robot.Create(ctx, stackConfig.Repositories, githubRepositories, harborProjects)
		if hrErr != nil {
			return hrErr
		}

		// terraform and aws integrations
		terraform, tfErr := terraform.Configure(ctx, stackConfig.Repositories)
		if tfErr != nil {
			return tfErr
		}
		aws, awsErr := aws.Configure(ctx, stackConfig.Repositories, githubRepositories)
		if awsErr != nil {
			return awsErr
		}

		exportTeams(ctx, githubTeams)
		exportRepositories(ctx, githubRepositories)
		exportHarbor(ctx, harborProjects, harborRobotAccounts)
		exportTerraform(ctx, terraform)
		exportAWS(ctx, aws)

		return nil
	})
}

// exportTeams exports the created GitHub teams.
// ctx: pulumi.Context.
// teams: map of created GitHub teams.
func exportTeams(ctx *pulumi.Context, teams map[string]*github.Team) {
	teamNames := pulumi.Array{}
	teamKeys := slices.Collect(maps.Keys(teams))
	slices.Sort(teamKeys)
	for _, team := range teamKeys {
		teamNames = append(teamNames, teams[team].Name)
	}

	ctx.Export("teams", teamNames)
}

// exportRepositories exports the created GitHub repositories.
// ctx: pulumi.Context.
// repositories: map of created GitHub repositories.
func exportRepositories(ctx *pulumi.Context, repositories map[string]*github.Repository) {
	repositoryNames := pulumi.Array{}
	repositoryKeys := slices.Collect(maps.Keys(repositories))
	slices.Sort(repositoryKeys)
	for _, repository := range repositoryKeys {
		repositoryNames = append(repositoryNames, repositories[repository].Name)
	}

	ctx.Export("repositories", repositoryNames)
}

// exportHarbor exports the created Harbor details.
// ctx: pulumi.Context.
// projects: map of created Harbor projects.
func exportHarbor(
	ctx *pulumi.Context,
	projects map[string]*harbor.ProjectOutput,
	robotAccounts map[string]*pulumi.StringOutput,
) {
	projectNames := pulumi.Array{}
	projectKeys := slices.Collect(maps.Keys(projects))
	slices.Sort(projectKeys)
	for _, project := range projectKeys {
		project, _ := projects[project].ApplyT(func(p *harbor.Project) pulumi.StringOutput { return p.Name }).(pulumi.StringOutput)
		projectNames = append(projectNames, project)
	}

	accounts := pulumi.Array{}
	robotAccountKeys := slices.Collect(maps.Keys(robotAccounts))
	slices.Sort(robotAccountKeys)
	for _, account := range robotAccountKeys {
		accounts = append(accounts, robotAccounts[account])
	}

	ctx.Export("harbor", pulumi.ToMap(map[string]any{
		"projects":      projectNames,
		"robotAccounts": accounts,
	}))
}

// exportTerraform exports the created Terraform buckets.
// ctx: pulumi.Context.
// terraform: map of created Terraform buckets.
func exportTerraform(ctx *pulumi.Context, terraform map[string]*pulumi.StringOutput) {
	buckets := pulumi.Map{}
	for repository, bucket := range terraform {
		buckets[repository] = bucket
	}

	ctx.Export("terraform", buckets)
}

// exportAWS exports the created AWS resources.
// ctx: pulumi.Context.
// aws: map of created AWS resources.
func exportAWS(ctx *pulumi.Context, aws map[string]*pulumi.StringOutput) {
	roles := pulumi.Map{}
	for repository, role := range aws {
		roles[repository] = role
	}

	ctx.Export("aws", roles)
}
