package export

import (
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// AWS exports the created AWS resources.
// ctx: pulumi.Context.
// aws: map of created AWS resources.
func AWS(ctx *pulumi.Context, aws map[string]*pulumi.StringOutput) {
	roles := pulumi.Map{}
	for repository, role := range aws {
		roles[repository] = role
	}

	ctx.Export("aws", roles)
}
