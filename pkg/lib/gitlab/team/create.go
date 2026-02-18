package team

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/muhlba91/fh-burgenland-bswe-assignment-infrastructure/pkg/lib/config"
	teamConf "github.com/muhlba91/fh-burgenland-bswe-assignment-infrastructure/pkg/model/config/team"
	"github.com/muhlba91/pulumi-shared-library/pkg/util/defaults"
	"github.com/pulumi/pulumi-gitlab/sdk/v9/go/gitlab"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// defaultVisibility is the default visibility for GitLab groups.
const defaultVisibility = "private"

// Create creates GitLab teams (groups) based on the provided configurations.
// ctx: The Pulumi context for resource creation.
// teams: A slice of team configuration objects.
func Create(ctx *pulumi.Context, teams []*teamConf.Config) (map[string]*gitlab.Group, error) {
	createdTeams := make(map[string]*gitlab.Group)

	envGroup, ghtErr := createGroup(
		ctx,
		pulumi.Int(config.Classroom.Gitlab.Group).ToIntOutput(),
		config.Environment,
		config.Environment,
		fmt.Sprintf("%s %s", config.Classroom.Name, strings.ToUpper(config.Environment)),
		nil,
	)
	if ghtErr != nil {
		return nil, ghtErr
	}
	envGroupID, _ := envGroup.ID().ToStringOutput().ApplyT(func(id string) int {
		gid, _ := strconv.Atoi(id)
		return gid
	}).(pulumi.IntOutput)

	for _, team := range teams {
		gitlabTeam, err := createTeam(ctx, team, envGroupID)
		if err != nil {
			return nil, err
		}
		createdTeams[team.Name] = gitlabTeam
	}

	return createdTeams, nil
}

// createTeam creates a single GitLab team (group) based on the provided configuration.
// ctx: The Pulumi context for resource creation.
// team: The configuration object for the team to be created.
// parentID: The ID of the parent group under which the new team group will be created.
func createTeam(ctx *pulumi.Context, team *teamConf.Config, parentID pulumi.IntOutput) (*gitlab.Group, error) {
	retainOnDelete := pulumi.RetainOnDelete(!defaults.GetOrDefault(team.DeleteOnDestroy, false))

	gitlabTeam, ghtErr := createGroup(
		ctx,
		parentID,
		team.Name,
		team.Name,
		fmt.Sprintf("%s %s: %s", config.Classroom.Name, strings.ToUpper(config.Environment), team.Name),
		team.DeleteOnDestroy,
	)
	if ghtErr != nil {
		return nil, ghtErr
	}

	groupID, _ := gitlabTeam.ID().ToStringOutput().ApplyT(func(id string) int {
		gid, _ := strconv.Atoi(id)
		return gid
	}).(pulumi.IntOutput)

	for _, member := range team.Members {
		if member == config.OwnerHandle {
			continue
		}

		user, uErr := gitlab.LookupUser(ctx, &gitlab.LookupUserArgs{
			Username: &member,
		})
		if uErr != nil {
			return nil, uErr
		}

		_, gtmErr := gitlab.NewGroupMembership(
			ctx,
			fmt.Sprintf("gitlab-group-member-%s-%s", team.Name, member),
			&gitlab.GroupMembershipArgs{
				GroupId:                    groupID,
				UserId:                     pulumi.Int(user.UserId),
				AccessLevel:                pulumi.String("developer"),
				SkipSubresourcesOnDestroy:  pulumi.Bool(true),
				UnassignIssuablesOnDestroy: pulumi.Bool(true),
			},
			retainOnDelete,
			pulumi.DependsOn([]pulumi.Resource{gitlabTeam}),
		)
		if gtmErr != nil {
			return nil, gtmErr
		}
	}

	return gitlabTeam, nil
}

// createGroup creates a GitLab group for the given team configuration.
// ctx: The Pulumi context for resource creation.
// parentID: The ID of the parent group under which the new group will be created.
// name: The name of the group to be created.
// path: The path of the group to be created.
// description: The description of the group to be created.
// deleteOnDestroy: Whether to delete the group on destroy.
func createGroup(
	ctx *pulumi.Context,
	parentID pulumi.IntOutput,
	name string,
	path string,
	description string,
	deleteOnDestroy *bool,
) (*gitlab.Group, error) {
	retainOnDelete := pulumi.RetainOnDelete(!defaults.GetOrDefault(deleteOnDestroy, false))

	return gitlab.NewGroup(ctx, fmt.Sprintf("gitlab-group-%s", path), &gitlab.GroupArgs{
		ParentId:    parentID,
		Path:        pulumi.String(path),
		Name:        pulumi.String(name),
		Description: pulumi.String(description),
		OnlyAllowMergeIfAllDiscussionsAreResolved: pulumi.Bool(true),
		OnlyAllowMergeIfPipelineSucceeds:          pulumi.Bool(true),
		AllowMergeOnSkippedPipeline:               pulumi.Bool(true),
		AutoDevopsEnabled:                         pulumi.Bool(false),
		ProjectCreationLevel:                      pulumi.String("maintainer"),
		RequestAccessEnabled:                      pulumi.Bool(false),
		SubgroupCreationLevel:                     pulumi.String("maintainer"),
		VisibilityLevel:                           pulumi.String(defaultVisibility),
		WikiAccessLevel:                           pulumi.String("private"),
	}, retainOnDelete)
}
