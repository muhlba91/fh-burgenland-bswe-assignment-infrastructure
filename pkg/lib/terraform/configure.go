package terraform

import (
	"fmt"

	"github.com/muhlba91/fh-burgenland-bswe-assignment-infrastructure/pkg/lib/config"
	"github.com/muhlba91/fh-burgenland-bswe-assignment-infrastructure/pkg/model/config/repository"
	"github.com/muhlba91/pulumi-shared-library/pkg/lib/aws/s3/bucket"
	"github.com/muhlba91/pulumi-shared-library/pkg/lib/github/actions/secret"
	"github.com/muhlba91/pulumi-shared-library/pkg/util/defaults"
	"github.com/pulumi/pulumi-github/sdk/v6/go/github"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Configure configures Terraform resources.
// ctx: pulumi.Context.
// repositories: list of repository configurations.
// githubRepositories: list of created GitHub repositories.
func Configure(
	ctx *pulumi.Context,
	repositories []*repository.Config,
	githubRepositories map[string]*github.Repository,
) (map[string]*pulumi.StringOutput, error) {
	buckets := make(map[string]*pulumi.StringOutput)

	for _, repo := range repositories {
		if !defaults.GetOrDefault(repo.Terraform, false) {
			continue
		}

		prefix := pulumi.StringPtr(fmt.Sprintf("bswe-%s-%s-%s", config.GlobalName, config.Environment, repo.Name))
		bucket, err := bucket.Create(ctx, &bucket.CreateOptions{
			Name:   fmt.Sprintf("terraform-%s-%s", config.Environment, repo.Name),
			Prefix: &prefix,
			Labels: config.CommonLabels(),
		})
		if err != nil {
			return nil, err
		}

		ghRepo, exists := githubRepositories[repo.Name]
		if !exists {
			return nil, fmt.Errorf("repository %s not found in created GitHub repositories", repo.Name)
		}

		secret.Write(ctx, &secret.WriteArgs{
			Repository: ghRepo,
			Key:        "TERRAFORM_BACKEND_BUCKET",
			Value:      bucket.Bucket,
		})

		buckets[repo.Name] = &bucket.Bucket
	}

	return buckets, nil
}
