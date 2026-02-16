package export

import (
	"maps"
	"slices"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumiverse/pulumi-harbor/sdk/v3/go/harbor"
)

// Harbor exports the created Harbor details.
// ctx: pulumi.Context.
// projects: map of created Harbor projects.
func Harbor(
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
