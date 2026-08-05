package secret

import (
	"fmt"

	"github.com/muhlba91/fh-burgenland-bswe-assignment-infrastructure/pkg/model/config/repository"
	"github.com/muhlba91/fh-burgenland-bswe-assignment-infrastructure/pkg/util/provider"
	ghSecret "github.com/muhlba91/pulumi-shared-library/pkg/lib/github/actions/secret"
	glSecret "github.com/muhlba91/pulumi-shared-library/pkg/lib/gitlab/actions/secret"
	"github.com/pulumi/pulumi-github/sdk/v6/go/github"
	"github.com/pulumi/pulumi-gitlab/sdk/v10/go/gitlab"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Write creates or updates a secret for the specified repository.
// ctx: The Pulumi context for resource management.
// repository: The configuration of the repository for which the secret should be created or updated.
// githubRepositories: A map of repository names to their corresponding GitHub repository resources.
// gitlabRepositories: A map of repository names to their corresponding GitLab project resources.
// key: The name of the secret to create or update.
// value: The value of the secret, which can be a Pulumi string input.
func Write(
	ctx *pulumi.Context,
	repository *repository.Config,
	githubRepositories map[string]*github.Repository,
	gitlabRepositories map[string]*gitlab.Project,
	key string,
	value pulumi.StringInput,
) error {
	if provider.GitHub(repository) {
		return writeGitHub(ctx, repository.Name, githubRepositories[repository.Name], key, value)
	}
	if provider.GitLab(repository) {
		return writeGitlab(ctx, repository.Name, gitlabRepositories[repository.Name], key, value, true)
	}

	return nil
}

// WriteUnmasked creates or updates an unmasked secret for the specified repository.
// ctx: The Pulumi context for resource management.
// repository: The configuration of the repository for which the secret should be created or updated.
// githubRepositories: A map of repository names to their corresponding GitHub repository resources.
// gitlabRepositories: A map of repository names to their corresponding GitLab project resources.
// key: The name of the secret to create or update.
// value: The value of the secret, which can be a Pulumi string input.
func WriteUnmasked(
	ctx *pulumi.Context,
	repository *repository.Config,
	githubRepositories map[string]*github.Repository,
	gitlabRepositories map[string]*gitlab.Project,
	key string,
	value pulumi.StringInput,
) error {
	if provider.GitHub(repository) {
		return writeGitHub(ctx, repository.Name, githubRepositories[repository.Name], key, value)
	}
	if provider.GitLab(repository) {
		return writeGitlab(ctx, repository.Name, gitlabRepositories[repository.Name], key, value, false)
	}

	return nil
}

// writeGitHub creates or updates a secret for the specified GitHub repository.
// ctx: The Pulumi context for resource management.
// name: The name of the repository for which the secret should be created or updated.
// repository: The configuration of the repository for which the secret should be created or updated.
// key: The name of the secret to create or update.
// value: The value of the secret, which can be a Pulumi string input.
func writeGitHub(
	ctx *pulumi.Context,
	name string,
	repository *github.Repository,
	key string,
	value pulumi.StringInput,
) error {
	if repository == nil {
		return fmt.Errorf("[secret]repository %s not found in created GitHub repositories", name)
	}

	ghSecret.Create(ctx, &ghSecret.CreateOptions{
		Repository: repository,
		Key:        key,
		Value:      value,
	})

	return nil
}

// writeGitlab creates or updates a secret for the specified GitLab project.
// ctx: The Pulumi context for resource management.
// name: The name of the project for which the secret should be created or updated.
// repository: The configuration of the project for which the secret should be created or updated.
// key: The name of the secret to create or update.
// value: The value of the secret, which can be a Pulumi string input.
// masked: A boolean indicating whether the secret should be masked (true) or not (false).
func writeGitlab(
	ctx *pulumi.Context,
	name string,
	repository *gitlab.Project,
	key string,
	value pulumi.StringInput,
	masked bool,
) error {
	if repository == nil {
		return fmt.Errorf("[secret]repository %s not found in created GitLab projects", name)
	}

	trueValue := true
	glSecret.Create(ctx, &glSecret.CreateOptions{
		Repository:               repository,
		Key:                      key,
		Value:                    value,
		Masked:                   &masked,
		DisableVariableExpansion: &trueValue,
	})

	return nil
}
