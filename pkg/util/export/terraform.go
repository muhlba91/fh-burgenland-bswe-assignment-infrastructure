package export

import (
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Terraform exports the created Terraform buckets.
// ctx: pulumi.Context.
// terraform: map of created Terraform buckets.
func Terraform(ctx *pulumi.Context, terraform map[string]*pulumi.StringOutput) {
	buckets := pulumi.Map{}
	for repository, bucket := range terraform {
		buckets[repository] = bucket
	}

	ctx.Export("terraform", buckets)
}
