package group

import (
	"fmt"

	"github.com/muhlba91/fh-burgenland-bswe-assignment-infrastructure/pkg/lib/config"
	"github.com/pulumi/pulumi-github/sdk/v6/go/github"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumiverse/pulumi-harbor/sdk/v3/go/harbor"
)

// defaultGroupType defines the default group type for Harbor groups.
const defaultGroupType = 3

// Create creates Harbor groups based on GitHub teams.
// ctx: pulumi.Context.
// teams: map of created GitHub teams.
func Create(ctx *pulumi.Context, teams map[string]*github.Team) (map[string]*harbor.Group, error) {
	groups := make(map[string]*harbor.Group)

	for teamName, team := range teams {
		group, err := harbor.NewGroup(ctx, fmt.Sprintf("harbor-group-%s", teamName), &harbor.GroupArgs{
			GroupName: pulumi.Sprintf("%s:%s", config.GitHubOrganization, team.Name),
			GroupType: pulumi.Int(defaultGroupType),
		})
		if err != nil {
			return nil, err
		}
		groups[teamName] = group
	}

	return groups, nil
}
