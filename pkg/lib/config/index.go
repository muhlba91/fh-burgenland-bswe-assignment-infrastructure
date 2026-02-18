package config

import (
	"fmt"
	"os"
	"strings"

	ghConfig "github.com/pulumi/pulumi-github/sdk/v6/go/github/config"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"

	"github.com/muhlba91/fh-burgenland-bswe-assignment-infrastructure/pkg/model/config/classroom"
	"github.com/muhlba91/fh-burgenland-bswe-assignment-infrastructure/pkg/model/config/stack"
	"github.com/muhlba91/fh-burgenland-bswe-assignment-infrastructure/pkg/util/data"
	"github.com/muhlba91/pulumi-shared-library/pkg/util/defaults"
)

//nolint:gochecknoglobals // global configuration is acceptable here
var (
	// Environment holds the current deployment environment (e.g., dev, staging, prod).
	Environment string
	// Classroom is a constant name used across resources.
	Classroom *classroom.Config
	// OwnerHandle is the owner handle used for resources.
	OwnerHandle = "muhlba91"
	// AWSDefaultRegion is the default AWS region for deployments.
	AWSDefaultRegion = "eu-west-1"
	// AWSAccountID is the AWS account ID used for deployments.
	AWSAccountID = "061039787254"
	// AllowRepositoryDeletion indicates whether repository deletion is permitted.
	AllowRepositoryDeletion = false
	// FeatureGates holds the feature gates for conditional resource creation.
	FeatureGates = []string{}
)

// LoadConfig loads the configuration for the given Pulumi context.
// ctx: The Pulumi context.
func LoadConfig(
	ctx *pulumi.Context,
) (*stack.Config, error) {
	Environment = ctx.Stack()

	repoDelEnv := strings.ToLower(os.Getenv("ALLOW_REPOSITORY_DELETION"))
	AllowRepositoryDeletion = defaults.GetOrDefault(&repoDelEnv, "false") == "true"

	stackConfig, rErr := data.ParseDataFromFiles(fmt.Sprintf("./assets/data_%s.yaml", Environment))
	if rErr != nil {
		return nil, rErr
	}

	FeatureGates = stackConfig.Features
	Classroom = stackConfig.Classroom

	if ghConfig.GetOwner(ctx) != Classroom.Github.Owner {
		return nil, fmt.Errorf(
			"GitHub owner mismatch: expected %s, got %s",
			Classroom.Github.Owner,
			ghConfig.GetOwner(ctx),
		)
	}

	return stackConfig, nil
}

// CommonLabels returns a map of common labels to be used across resources.
func CommonLabels() map[string]string {
	return map[string]string{
		"environment": Environment,
		"purpose":     Classroom.Tag,
	}
}
