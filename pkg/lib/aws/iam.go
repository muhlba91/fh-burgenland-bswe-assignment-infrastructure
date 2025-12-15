package aws

import (
	"encoding/json"
	"fmt"

	"github.com/muhlba91/fh-burgenland-bswe-assignment-infrastructure/pkg/lib/config"
	"github.com/muhlba91/pulumi-shared-library/pkg/lib/aws/iam/policy"
	"github.com/muhlba91/pulumi-shared-library/pkg/lib/aws/iam/role"
	"github.com/muhlba91/pulumi-shared-library/pkg/lib/github/actions/secret"
	"github.com/muhlba91/pulumi-shared-library/pkg/lib/random"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/iam"
	"github.com/pulumi/pulumi-github/sdk/v6/go/github"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// createAccountIAM creates AWS IAM roles for Continuous Integration for the specified repository account.
// ctx: Pulumi context for resource management.
// repository: The GitHub repository.
// identityProviderArn: ARN of the AWS IAM Identity Provider for GitHub OIDC.
func createAccountIAM(ctx *pulumi.Context,
	repository *github.Repository,
	identityProviderArn string,
) pulumi.StringOutput {
	tags := config.CommonLabels()
	tags["organization"] = "fh-burgenland-bswe"

	roleArn, _ := repository.Name.ApplyT(func(name string) pulumi.StringOutput {
		truncatedRepository := name[:min(maxRepositoryLength, len(name))]

		postfix, _ := random.CreateString(
			ctx,
			fmt.Sprintf("random-string-aws-iam-role-ci-%s", name),
			&random.StringOptions{
				Length:  postfixLength,
				Special: false,
			},
		)

		ciRole, _ := createRole(ctx, name, identityProviderArn, tags, truncatedRepository, postfix.Text)
		_ = createPolicy(ctx, name, ciRole, tags, truncatedRepository, postfix.Text)

		secret.Write(ctx, &secret.WriteArgs{
			Repository: repository,
			Key:        "AWS_IDENTITY_ROLE_ARN",
			Value:      ciRole.Arn,
		})
		secret.Write(ctx, &secret.WriteArgs{
			Repository: repository,
			Key:        "AWS_REGION",
			Value:      pulumi.String(config.AWSDefaultRegion),
		})

		return ciRole.Arn
	}).(pulumi.StringOutput)

	return roleArn
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
							config.GitHubOrganization,
							repository,
						),
					},
				},
			},
		},
	})

	ciRole, cirErr := role.Create(ctx, repository, &role.CreateOptions{
		Name:             pulumi.Sprintf("ci-%s-%s", truncatedRepository, ciPostfix),
		Description:      pulumi.Sprintf("FH Burgenland Softwaremanagement II GitHub Repository: %s", repository),
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
					fmt.Sprintf("arn:aws:s3:::bswe-%s-%s-*", config.GlobalName, config.Environment),
					fmt.Sprintf("arn:aws:s3:::bswe-%s-%s-*/*", config.GlobalName, config.Environment),
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
		Name:        pulumi.Sprintf("ci-%s-%s", truncatedRepository, ciPostfix),
		Description: pulumi.Sprintf("FH Burgenland Softwaremanagement II GitHub Repository: %s", repository),
		Policy:      pulumi.String(policyDoc.Json),
		Labels:      tags,
	})
	if cipolErr != nil {
		return cipolErr
	}

	_, paErr := role.CreatePolicyAttachment(ctx, fmt.Sprintf("ci-%s", repository), &role.CreatePolicyAttachmentOptions{
		Roles:     pulumi.StringArray{ciRole.Name},
		PolicyArn: ciPolicy.Arn,
		PulumiOptions: []pulumi.ResourceOption{
			pulumi.DependsOn([]pulumi.Resource{ciRole, ciPolicy}),
		},
	})
	return paErr
}
