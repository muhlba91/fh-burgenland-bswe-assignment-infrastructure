package harbor

import (
	"github.com/muhlba91/fh-burgenland-bswe-assignment-infrastructure/pkg/lib/harbor/auth"
	"github.com/muhlba91/fh-burgenland-bswe-assignment-infrastructure/pkg/lib/harbor/group"
	"github.com/muhlba91/fh-burgenland-bswe-assignment-infrastructure/pkg/lib/harbor/project"
	"github.com/muhlba91/fh-burgenland-bswe-assignment-infrastructure/pkg/lib/harbor/robot"
	"github.com/muhlba91/fh-burgenland-bswe-assignment-infrastructure/pkg/model/config/stack"
	"github.com/muhlba91/fh-burgenland-bswe-assignment-infrastructure/pkg/util/feature"
	"github.com/pulumi/pulumi-github/sdk/v6/go/github"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumiverse/pulumi-harbor/sdk/v3/go/harbor"
	"github.com/rs/zerolog/log"
)

// Configure sets up the Harbor configuration, including authentication, groups, projects, and robot accounts.
// ctx: The Pulumi context for resource management.
// stackConfig: The configuration for the stack, containing repository information.
// githubRepositories: A map of GitHub repositories that may be needed for project creation.
// teams: A map of GitHub teams that may be needed for group creation.
func Configure(
	ctx *pulumi.Context,
	stackConfig *stack.Config,
	githubRepositories map[string]*github.Repository,
	githubTeams map[string]*github.Team,
) (map[string]*harbor.ProjectOutput, map[string]*pulumi.StringOutput, error) {
	if !feature.Harbor() {
		log.Info().Msg("[harbor] feature is disabled, skipping harbor configuration")
		return map[string]*harbor.ProjectOutput{}, map[string]*pulumi.StringOutput{}, nil
	}

	haErr := auth.Configure(ctx)
	if haErr != nil {
		return nil, nil, haErr
	}

	harborGroups, hgErr := group.Create(ctx, githubTeams)
	if hgErr != nil {
		return nil, nil, hgErr
	}

	harborProjects := project.Create(ctx, stackConfig.Repositories, githubRepositories, harborGroups)

	harborRobotAccounts, hrErr := robot.Create(ctx, stackConfig.Repositories, githubRepositories, harborProjects)
	if hrErr != nil {
		return nil, nil, hrErr
	}

	return harborProjects, harborRobotAccounts, nil
}
