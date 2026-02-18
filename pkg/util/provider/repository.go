package provider

import (
	"github.com/muhlba91/fh-burgenland-bswe-assignment-infrastructure/pkg/model/config/repository"
	"github.com/muhlba91/pulumi-shared-library/pkg/util/defaults"
)

const defaultProvider = "none"

// GitHub checks if the given repository configuration is intended for GitHub based on the provider field.
// repository: The configuration of the repository to check.
func GitHub(repository *repository.Config) bool {
	return defaults.GetOrDefault(repository.Provider, defaultProvider) == "github"
}

// GitLab checks if the given repository configuration is intended for GitLab based on the provider field.
// repository: The configuration of the repository to check.
func GitLab(repository *repository.Config) bool {
	return defaults.GetOrDefault(repository.Provider, defaultProvider) == "gitlab"
}
