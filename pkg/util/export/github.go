package export

import (
	"maps"
	"slices"

	"github.com/pulumi/pulumi-github/sdk/v6/go/github"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// GitHub exports the created GitHub resources.
// ctx: pulumi.Context.
// teams: map of created GitHub teams.
// repositories: map of created GitHub repositories.
func GitHub(ctx *pulumi.Context, teams map[string]*github.Team, repositories map[string]*github.Repository) {
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

	ctx.Export("github", pulumi.ToMap(map[string]any{
		"repositories": repositoryNames,
		"teams":        teamNames,
	}))
}
