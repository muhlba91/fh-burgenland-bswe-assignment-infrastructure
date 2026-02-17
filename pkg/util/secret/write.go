package secret

import (
	"fmt"

	"github.com/muhlba91/fh-burgenland-bswe-assignment-infrastructure/pkg/model/config/repository"
	"github.com/muhlba91/fh-burgenland-bswe-assignment-infrastructure/pkg/util/provider"
	"github.com/muhlba91/pulumi-shared-library/pkg/lib/github/actions/secret"
	"github.com/pulumi/pulumi-github/sdk/v6/go/github"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Write creates or updates a secret for the specified repository.
// ctx: The Pulumi context for resource management.
// repository: The configuration of the repository for which the secret should be created or updated.
// githubRepositories: A map of repository names to their corresponding GitHub repository resources.
// key: The name of the secret to create or update.
// value: The value of the secret, which can be a Pulumi string input.
func Write(
	ctx *pulumi.Context,
	repository *repository.Config,
	githubRepositories map[string]*github.Repository,
	key string,
	value pulumi.StringInput,
) error {
	if provider.GitHub(repository) {
		return writeGitHub(ctx, repository.Name, githubRepositories[repository.Name], key, value)
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

	secret.Create(ctx, &secret.CreateOptions{
		Repository: repository,
		Key:        key,
		Value:      value,
	})

	return nil
}
