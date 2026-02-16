package feature

import (
	"slices"

	"github.com/muhlba91/fh-burgenland-bswe-assignment-infrastructure/pkg/lib/config"
)

// Enabled checks if a feature gate is enabled in the stack.
// featureGate: The name of the feature gate to check.
func Enabled(featureGate string) bool {
	return slices.Contains(config.FeatureGates, featureGate)
}

// Harbor checks if the "harbor" feature gate is enabled in the stack configuration.
func Harbor() bool {
	return Enabled("harbor")
}

// Terraform checks if the "terraform" feature gate is enabled in the stack configuration.
func Terraform() bool {
	return Enabled("terraform")
}

// AWS checks if the "aws" feature gate is enabled in the stack configuration.
func AWS() bool {
	return Enabled("aws")
}
