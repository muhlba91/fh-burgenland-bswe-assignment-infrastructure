package export

import (
	"maps"
	"slices"

	"github.com/pulumi/pulumi-gitlab/sdk/v9/go/gitlab"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// GitLab exports the created GitLab resources.
// ctx: pulumi.Context.
// teams: map of created GitLab teams.
// repositories: map of created GitLab repositories.
func GitLab(ctx *pulumi.Context, teams map[string]*gitlab.Group, repositories map[string]*gitlab.Project) {
	teamNames := pulumi.Array{}
	teamKeys := slices.Collect(maps.Keys(teams))
	slices.Sort(teamKeys)
	for _, team := range teamKeys {
		teamNames = append(teamNames, teams[team].Name)
	}

	repositoryNames := pulumi.Array{}
	repositoryKeys := slices.Collect(maps.Keys(repositories))
	slices.Sort(repositoryKeys)
	for _, repository := range repositoryKeys {
		repositoryNames = append(repositoryNames, repositories[repository].Name)
	}

	ctx.Export("gitlab", pulumi.ToMap(map[string]any{
		"repositories": repositoryNames,
		"teams":        teamNames,
	}))
}
