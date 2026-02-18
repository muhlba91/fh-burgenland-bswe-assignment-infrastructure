package repositories

import (
	"fmt"

	"github.com/pulumi/pulumi-gitlab/sdk/v9/go/gitlab"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"

	"github.com/muhlba91/fh-burgenland-bswe-assignment-infrastructure/pkg/lib/config"
	"github.com/muhlba91/fh-burgenland-bswe-assignment-infrastructure/pkg/model/config/repository"
	libRuleset "github.com/muhlba91/pulumi-shared-library/pkg/lib/gitlab/ruleset"
)

// createRuleset creates a branch ruleset for the given repository based on the provided configuration.
// ctx: The Pulumi context for resource creation.
// repository: The configuration for the repository.
// repo: The Pulumi GitLab repository resource.
func createRuleset(
	ctx *pulumi.Context,
	repository *repository.Config,
	repo *gitlab.Project,
) error {
	_, err := libRuleset.Create(
		ctx,
		fmt.Sprintf("%s-%s", config.Environment, repository.Name),
		&libRuleset.CreateOptions{
			Repository:    repo,
			Branch:        libRuleset.DefaultBranch,
			ReviewerCount: repository.Approvers,
		},
	)
	return err
}
