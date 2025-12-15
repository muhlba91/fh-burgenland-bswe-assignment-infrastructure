package aws

import (
	"fmt"

	"github.com/muhlba91/fh-burgenland-bswe-assignment-infrastructure/pkg/model/config/repository"
	"github.com/muhlba91/pulumi-shared-library/pkg/util/defaults"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/iam"
	"github.com/pulumi/pulumi-github/sdk/v6/go/github"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Configure sets up AWS resources based on the provided configuration.
// ctx: Pulumi context for resource management.
// repositories: List of repository configurations.
// githubRepositories: list of created GitHub repositories.
func Configure(ctx *pulumi.Context,
	repositories []*repository.Config,
	githubRepositories map[string]*github.Repository,
) (map[string]*pulumi.StringOutput, error) {
	accounts := make(map[string]*pulumi.StringOutput)

	githubOidcURL := "https://token.actions.githubusercontent.com"
	identityProvider, ipErr := iam.LookupOpenIdConnectProvider(ctx, &iam.LookupOpenIdConnectProviderArgs{
		Url: &githubOidcURL,
	})
	if ipErr != nil {
		return nil, ipErr
	}

	for _, repo := range repositories {
		if !defaults.GetOrDefault(repo.AWS, false) {
			continue
		}

		ghRepo, exists := githubRepositories[repo.Name]
		if !exists {
			return nil, fmt.Errorf("repository %s not found in created GitHub repositories", repo.Name)
		}

		roleArn := createAccountIAM(ctx, ghRepo, identityProvider.Arn)

		accounts[repo.Name] = &roleArn
	}

	return accounts, nil
}
