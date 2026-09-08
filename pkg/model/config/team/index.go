package team

// Config defines team-related configuration.
type Config struct {
	// Enabled indicates whether the team is enabled.
	Enabled *bool `yaml:"enabled"`
	// Name is the name of the team.
	Name string `yaml:"name"`
	// DeleteOnDestroy indicates whether the team should be deleted when the stack is destroyed.
	DeleteOnDestroy *bool `yaml:"deleteOnDestroy,omitempty"`
	// Members defines the members of the team.
	Members *Members `yaml:"members,omitempty"`
}

// Members defines team-member-related configuration.
type Members struct {
	// GitHub defines the GitHub usernames of the team members.
	GitHub []string `yaml:"github,omitempty"`
	// GitLab defines the GitLab usernames of the team members.
	GitLab []string `yaml:"gitlab,omitempty"`
}
