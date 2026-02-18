package repositories

import (
	"sort"
	"strconv"

	"github.com/pulumi/pulumi-gitlab/sdk/v9/go/gitlab"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"

	"github.com/muhlba91/fh-burgenland-bswe-assignment-infrastructure/pkg/lib/config"
	"github.com/muhlba91/fh-burgenland-bswe-assignment-infrastructure/pkg/model/config/repository"
	"github.com/muhlba91/fh-burgenland-bswe-assignment-infrastructure/pkg/util/provider"
	libRepo "github.com/muhlba91/pulumi-shared-library/pkg/lib/gitlab/repository"
	"github.com/muhlba91/pulumi-shared-library/pkg/util/defaults"
)

// defaultVisibility is the default visibility for GitLab repositories.
const defaultVisibility = "private"

// defaultDeletePipelinesInSeconds is the default time in seconds after which pipelines will be automatically deleted.
const defaultDeletePipelinesInSeconds = 1209600 // 14 days

// Create creates multiple GitLab repositories based on the provided configuration.
// ctx: The Pulumi context for resource creation.
// repositories: A slice of repository configurations to create.
// gitlabTeams: A map of GitLab teams (groups) for potential repository team associations.
func Create(
	ctx *pulumi.Context,
	repositories []*repository.Config,
	gitlabTeams map[string]*gitlab.Group,
) (map[string]*gitlab.Project, error) {
	repos := make(map[string]*gitlab.Project)

	for _, repo := range repositories {
		if !provider.GitLab(repo) {
			continue
		}

		glRepo, err := create(ctx, repo, gitlabTeams)
		if err != nil {
			return nil, err
		}
		repos[repo.Name] = glRepo
	}

	return repos, nil
}

// create creates a single GitLab repository based on the provided configuration.
// ctx: The Pulumi context for resource creation.
// repository: The configuration for the repository to create.
// gitlabTeams: A map of GitLab teams (groups) for potential repository team associations.
func create(
	ctx *pulumi.Context,
	repository *repository.Config,
	gitlabTeams map[string]*gitlab.Group,
) (*gitlab.Project, error) {
	topics := []string{
		config.Classroom.Tag,
		config.Environment,
		repository.Service,
	}
	sort.Strings(topics)

	var owner pulumi.IntOutput
	for _, team := range repository.Teams {
		if team.Role == "owner" {
			owner, _ = gitlabTeams[team.Name].ID().ToStringOutput().ApplyT(func(id string) int {
				gid, _ := strconv.Atoi(id)
				return gid
			}).(pulumi.IntOutput)
			break
		}
	}

	defVis := defaultVisibility
	trueValue := true
	deletePipelinesInSeconds := defaultDeletePipelinesInSeconds
	repo, err := libRepo.Create(ctx, repository.Name, &libRepo.CreateOptions{
		Name: pulumi.Sprintf("%s-%s-%s", config.Classroom.Tag, config.Environment, repository.Name),
		Description: pulumi.Sprintf(
			"%s %s: %s repository",
			config.Classroom.Name,
			config.Environment,
			repository.Service,
		),
		NamespaceID:              owner,
		ConversationResolution:   &trueValue,
		Topics:                   topics,
		Visibility:               &defVis,
		DeletePipelinesInSeconds: &deletePipelinesInSeconds,
		AllowRepositoryDeletion:  false,
		RetainOnDelete:           &trueValue,
		Protected:                !defaults.GetOrDefault(repository.DeleteOnDestroy, false),
	})
	if err != nil {
		return nil, err
	}

	rsErr := createRuleset(ctx, repository, repo)
	if rsErr != nil {
		return nil, rsErr
	}

	raErr := createAccess(ctx, repository, repo, gitlabTeams)
	if raErr != nil {
		return nil, raErr
	}

	return repo, nil
}
