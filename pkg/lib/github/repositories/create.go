package repositories

import (
	"sort"

	"github.com/pulumi/pulumi-github/sdk/v6/go/github"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"

	"github.com/muhlba91/fh-burgenland-bswe-assignment-infrastructure/pkg/lib/config"
	"github.com/muhlba91/fh-burgenland-bswe-assignment-infrastructure/pkg/model/config/repository"
	libRepo "github.com/muhlba91/pulumi-shared-library/pkg/lib/github/repository"
	"github.com/muhlba91/pulumi-shared-library/pkg/util/defaults"
)

// defaultVisibility is the default visibility for GitHub repositories.
const defaultVisibility = "private"

// Create creates multiple GitHub repositories based on the provided configuration.
// ctx: The Pulumi context for resource creation.
// repositories: A slice of repository configurations to create.
// githubTeams: A map of GitHub teams for potential repository team associations.
func Create(
	ctx *pulumi.Context,
	repositories []*repository.Config,
	githubTeams map[string]*github.Team,
) (map[string]*github.Repository, error) {
	repos := make(map[string]*github.Repository)

	for _, repo := range repositories {
		ghRepo, err := create(ctx, repo, githubTeams)
		if err != nil {
			return nil, err
		}
		repos[repo.Name] = ghRepo
	}

	return repos, nil
}

// create creates a single GitHub repository based on the provided configuration.
// ctx: The Pulumi context for resource creation.
// repository: The configuration for the repository to create.
// githubTeams: A map of GitHub teams for potential repository team associations.
func create(
	ctx *pulumi.Context,
	repository *repository.Config,
	githubTeams map[string]*github.Team,
) (*github.Repository, error) {
	topics := []string{
		config.GlobalName,
		config.Environment,
		repository.Service,
	}
	sort.Strings(topics)

	defVis := defaultVisibility
	repo, err := libRepo.Create(ctx, repository.Name, &libRepo.CreateOptions{
		Name: pulumi.Sprintf("%s-%s-%s", config.GlobalName, config.Environment, repository.Name),
		Description: pulumi.Sprintf(
			"Softwaremanagement II %s: %s repository",
			config.Environment,
			repository.Service,
		),
		EnableDiscussions: pulumi.Bool(false),
		EnableWiki:        pulumi.Bool(true),
		Topics:            topics,
		Visibility:        &defVis,
		Protected:         !defaults.GetOrDefault(repository.DeleteOnDestroy, false),
	})
	if err != nil {
		return nil, err
	}

	rsErr := createRuleset(ctx, repository, repo)
	if rsErr != nil {
		return nil, rsErr
	}

	raErr := createAccess(ctx, repository, repo, githubTeams)
	if raErr != nil {
		return nil, raErr
	}

	return repo, nil
}
