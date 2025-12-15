package project

import (
	"fmt"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumiverse/pulumi-harbor/sdk/v3/go/harbor"
)

// defaultPullPushHistory defines the default number of most recently pulled/pushed images to retain.
const defaultPullPushHistory = 3

// createPolicies creates retention policies for the given Harbor project.
// ctx: pulumi.Context.
// project: Harbor project to create policies for.
func createPolicies(ctx *pulumi.Context, project *harbor.Project) {
	project.Name.ApplyT(func(name string) error {
		_, err := harbor.NewRetentionPolicy(
			ctx,
			fmt.Sprintf("harbor-retention-policy-%s", name),
			&harbor.RetentionPolicyArgs{
				Scope:    project.ID(),
				Schedule: pulumi.String("Daily"),
				Rules: harbor.RetentionPolicyRuleArray{
					&harbor.RetentionPolicyRuleArgs{
						Disabled:           pulumi.Bool(false),
						RepoMatching:       pulumi.String("**"),
						TagMatching:        pulumi.String("**"),
						MostRecentlyPulled: pulumi.Int(defaultPullPushHistory),
					},
					&harbor.RetentionPolicyRuleArgs{
						Disabled:           pulumi.Bool(false),
						RepoMatching:       pulumi.String("**"),
						TagMatching:        pulumi.String("**"),
						MostRecentlyPushed: pulumi.Int(defaultPullPushHistory),
					},
				},
			},
		)
		return err
	})
}
