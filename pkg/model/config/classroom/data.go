package classroom

// Config defines stack-related configuration.
type Config struct {
	// Name is the name of the classroom.
	Name string `yaml:"name,omitempty"`
	// Tag is the tag of the classroom.
	Tag string `yaml:"tag,omitempty"`
	// Github defines the GitHub configuration for the classroom.
	Github *GithubConfig `yaml:"github,omitempty"`
	// Gitlab defines the GitLab configuration for the classroom.
	Gitlab *GitlabConfig `yaml:"gitlab,omitempty"`
}

// GithubConfig defines the GitHub configuration for the classroom.
type GithubConfig struct {
	// Owner is the owner of the GitHub organization.
	Owner string `yaml:"owner,omitempty"`
}

// GitlabConfig defines the GitLab configuration for the classroom.
type GitlabConfig struct {
	// Group is the group ID of the GitLab instance.
	Group int `yaml:"group,omitempty"`
}
