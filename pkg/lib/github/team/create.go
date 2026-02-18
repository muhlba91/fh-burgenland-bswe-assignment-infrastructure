package team

import (
	"fmt"

	"github.com/muhlba91/fh-burgenland-bswe-assignment-infrastructure/pkg/lib/config"
	teamConf "github.com/muhlba91/fh-burgenland-bswe-assignment-infrastructure/pkg/model/config/team"
	"github.com/muhlba91/pulumi-shared-library/pkg/util/defaults"
	"github.com/pulumi/pulumi-github/sdk/v6/go/github"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Create creates GitHub teams based on the provided configurations.
// ctx: The Pulumi context for resource creation.
// teams: A slice of team configuration objects.
func Create(ctx *pulumi.Context, teams []*teamConf.Config) (map[string]*github.Team, error) {
	createdTeams := make(map[string]*github.Team)

	for _, team := range teams {
		githubTeam, err := createTeam(ctx, team)
		if err != nil {
			return nil, err
		}
		createdTeams[team.Name] = githubTeam
	}

	return createdTeams, nil
}

// createTeam creates a single GitHub team based on the provided configuration.
// ctx: The Pulumi context for resource creation.
// team: The configuration object for the team to be created.
func createTeam(ctx *pulumi.Context, team *teamConf.Config) (*github.Team, error) {
	retainOnDelete := pulumi.RetainOnDelete(!defaults.GetOrDefault(team.DeleteOnDestroy, false))

	githubTeam, ghtErr := github.NewTeam(ctx, fmt.Sprintf("github-team-%s", team.Name), &github.TeamArgs{
		Name:        pulumi.Sprintf("%s-%s-%s", config.Classroom.Tag, config.Environment, team.Name),
		Description: pulumi.Sprintf("%s %s: %s", config.Classroom.Name, config.Environment, team.Name),
		Privacy:     pulumi.String("secret"),
	}, retainOnDelete)
	if ghtErr != nil {
		return nil, ghtErr
	}

	for _, member := range team.Members {
		if member == config.OwnerHandle {
			continue
		}

		_, gtmErr := github.NewTeamMembership(
			ctx,
			fmt.Sprintf("github-team-membership-%s-%s", team.Name, member),
			&github.TeamMembershipArgs{
				TeamId:   githubTeam.ID(),
				Username: pulumi.String(member),
				Role:     pulumi.String("member"),
			},
			retainOnDelete,
			pulumi.DependsOn([]pulumi.Resource{githubTeam}),
		)
		if gtmErr != nil {
			return nil, gtmErr
		}
	}

	return githubTeam, nil
}
