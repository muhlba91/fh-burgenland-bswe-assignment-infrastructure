package repositories

import (
	"fmt"

	"github.com/pulumi/pulumi-github/sdk/v6/go/github"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"

	"github.com/muhlba91/fh-burgenland-bswe-assignment-infrastructure/pkg/lib/config"
	"github.com/muhlba91/fh-burgenland-bswe-assignment-infrastructure/pkg/model/config/repository"
	libRuleset "github.com/muhlba91/pulumi-shared-library/pkg/lib/github/ruleset"
)

// createRuleset creates a branch ruleset for the given repository based on the provided configuration.
// ctx: The Pulumi context for resource creation.
// repository: The configuration for the repository.
// repo: The Pulumi GitHub repository resource.
func createRuleset(
	ctx *pulumi.Context,
	repository *repository.Config,
	repo *github.Repository,
) error {
	wipIntegration := false
	_, err := libRuleset.Create(
		ctx,
		fmt.Sprintf("%s-%s", config.Environment, repository.Name),
		&libRuleset.CreateOptions{
			Repository:     repo,
			Patterns:       []string{libRuleset.DefaultBranchRulesetPattern},
			ReviewerCount:  repository.Approvers,
			RequiredChecks: repository.RequiredChecks,
			WIPIntegration: &wipIntegration,
		},
	)
	return err
}
