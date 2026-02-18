package aws

import (
	"encoding/json"
	"fmt"

	"github.com/muhlba91/fh-burgenland-bswe-assignment-infrastructure/pkg/lib/config"
	"github.com/muhlba91/fh-burgenland-bswe-assignment-infrastructure/pkg/model/config/repository"
	"github.com/muhlba91/fh-burgenland-bswe-assignment-infrastructure/pkg/util/secret"
	"github.com/muhlba91/pulumi-shared-library/pkg/lib/aws/iam/policy"
	"github.com/muhlba91/pulumi-shared-library/pkg/lib/aws/iam/role"
	"github.com/muhlba91/pulumi-shared-library/pkg/lib/random"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/iam"
	"github.com/pulumi/pulumi-github/sdk/v6/go/github"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// createAccountIAM creates AWS IAM roles for Continuous Integration for the specified repository account.
// ctx: Pulumi context for resource management.
// repository: The name of the GitHub repository for which to create the IAM role.
// identityProviderArn: ARN of the AWS IAM Identity Provider for GitHub OIDC.
// githubRepositories: Map of created GitHub repositories to ensure the repository exists before creating IAM roles.
func createAccountIAM(ctx *pulumi.Context,
	repository *repository.Config,
	identityProviderArn string,
	githubRepositories map[string]*github.Repository,
) pulumi.StringOutput {
	tags := config.CommonLabels()
	tags["organization"] = "fh-burgenland-bswe"

	truncatedRepository := repository.Name[:min(maxRepositoryLength, len(repository.Name))]

	postfix, _ := random.CreateString(
		ctx,
		fmt.Sprintf("random-string-aws-iam-role-%s", repository.Name),
		&random.StringOptions{
			Length:  postfixLength,
			Special: false,
		},
	)

	ciRole, _ := createRole(ctx, repository.Name, identityProviderArn, tags, truncatedRepository, postfix.Text)
	_ = createPolicy(ctx, repository.Name, ciRole, tags, truncatedRepository, postfix.Text)

	_ = secret.Write(ctx, repository, githubRepositories, "AWS_IDENTITY_ROLE_ARN", ciRole.Arn)
	_ = secret.Write(ctx, repository, githubRepositories, "AWS_REGION", pulumi.String(config.AWSDefaultRegion))

	return ciRole.Arn
}

// createRole creates an AWS IAM role for Continuous Integration for the specified repository account.
// ctx: Pulumi context for resource management.
// repository: The name of the GitHub repository.
// identityProviderArn: ARN of the AWS IAM Identity Provider for GitHub OIDC.
// tags: Tags to be applied to the IAM role.
// truncatedRepository: Truncated name of the repository for naming purposes.
// ciPostfix: Random postfix for ensuring unique role names.
func createRole(ctx *pulumi.Context,
	repository string,
	identityProviderArn string,
	tags map[string]string,
	truncatedRepository string,
	ciPostfix pulumi.StringOutput,
) (*iam.Role, error) {
	//nolint:gosec // the token url is required for the trust relationship and does not contain any sensitive information
	policyDoc, _ := json.Marshal(map[string]any{
		"Version": "2012-10-17",
		"Statement": []map[string]any{
			{
				"Effect": "Allow",
				"Action": "sts:AssumeRoleWithWebIdentity",
				"Principal": map[string]any{
					"Federated": identityProviderArn,
				},
				"Condition": map[string]any{
					"StringEquals": map[string]any{
						"token.actions.githubusercontent.com:aud": "sts.amazonaws.com",
					},
					"StringLike": map[string]any{
						"token.actions.githubusercontent.com:sub": fmt.Sprintf(
							"repo:%s/%s:*",
							config.Classroom.Github.Owner,
							repository,
						),
					},
				},
			},
		},
	})

	ciRole, cirErr := role.Create(ctx, repository, &role.CreateOptions{
		Name:             pulumi.Sprintf("%s-%s", truncatedRepository, ciPostfix),
		Description:      pulumi.Sprintf("%s GitHub Repository: %s", config.Classroom.Name, repository),
		AssumeRolePolicy: pulumi.String(policyDoc),
		Labels:           tags,
	})
	if cirErr != nil {
		return nil, cirErr
	}

	return ciRole, nil
}

// createPolicy creates an AWS IAM policy for Continuous Integration for the specified repository account.
// ctx: Pulumi context for resource management.
// repository: The name of the GitHub repository.
// ciRole: The IAM role to attach the policy to.
// tags: Tags to be applied to the IAM role.
// truncatedRepository: Truncated name of the repository for naming purposes.
// ciPostfix: Random postfix for ensuring unique role names.
func createPolicy(ctx *pulumi.Context,
	repository string,
	ciRole *iam.Role,
	tags map[string]string,
	truncatedRepository string,
	ciPostfix pulumi.StringOutput,
) error {
	allow := "Allow"
	policyDoc, polErr := iam.GetPolicyDocument(ctx, &iam.GetPolicyDocumentArgs{
		Statements: []iam.GetPolicyDocumentStatement{
			{
				Effect:  &allow,
				Actions: []string{"s3:*"},
				Resources: []string{
					fmt.Sprintf("arn:aws:s3:::bswe-%s-%s-*", config.Classroom.Tag, config.Environment),
					fmt.Sprintf("arn:aws:s3:::bswe-%s-%s-*/*", config.Classroom.Tag, config.Environment),
				},
			},
			{
				Effect:  &allow,
				Actions: []string{"cloudfront:*"},
				Resources: []string{
					fmt.Sprintf("arn:aws:cloudfront::%s:distribution/*", config.AWSAccountID),
					fmt.Sprintf("arn:aws:cloudfront::%s:origin-access-identity/*", config.AWSAccountID),
					fmt.Sprintf("arn:aws:cloudfront::%s:origin-request-policy/*", config.AWSAccountID),
					fmt.Sprintf("arn:aws:cloudfront::%s:response-headers-policy/*", config.AWSAccountID),
					fmt.Sprintf("arn:aws:cloudfront::%s:origin-access-control/*", config.AWSAccountID),
				},
			},
		},
	})
	if polErr != nil {
		return polErr
	}

	ciPolicy, cipolErr := policy.Create(ctx, repository, &policy.CreateOptions{
		Name:        pulumi.Sprintf("%s-%s", truncatedRepository, ciPostfix),
		Description: pulumi.Sprintf("%s GitHub Repository: %s", config.Classroom.Name, repository),
		Policy:      pulumi.String(policyDoc.Json),
		Labels:      tags,
	})
	if cipolErr != nil {
		return cipolErr
	}

	_, paErr := role.CreatePolicyAttachment(ctx, repository, &role.CreatePolicyAttachmentOptions{
		Roles:     pulumi.StringArray{ciRole.Name},
		PolicyArn: ciPolicy.Arn,
		PulumiOptions: []pulumi.ResourceOption{
			pulumi.DependsOn([]pulumi.Resource{ciRole, ciPolicy}),
		},
	})
	return paErr
}
