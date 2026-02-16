package aws

import (
	"github.com/muhlba91/fh-burgenland-bswe-assignment-infrastructure/pkg/model/config/repository"
	"github.com/muhlba91/fh-burgenland-bswe-assignment-infrastructure/pkg/util/feature"
	"github.com/muhlba91/pulumi-shared-library/pkg/util/defaults"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/iam"
	"github.com/pulumi/pulumi-github/sdk/v6/go/github"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/rs/zerolog/log"
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

	if !feature.AWS() {
		log.Info().Msg("[aws] feature is disabled, skipping aws configuration")
		return accounts, nil
	}

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

		roleArn := createAccountIAM(ctx, repo, identityProvider.Arn, githubRepositories)
		accounts[repo.Name] = &roleArn
	}

	return accounts, nil
}
