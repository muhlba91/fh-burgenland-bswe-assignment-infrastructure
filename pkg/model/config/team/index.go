package team

// Config defines team-related configuration.
type Config struct {
	// Name is the name of the team.
	Name string `yaml:"name"`
	// DeleteOnDestroy indicates whether the team should be deleted when the stack is destroyed.
	DeleteOnDestroy *bool `yaml:"deleteOnDestroy,omitempty"`
	// Members defines the members of the team.
	Members []string `yaml:"members,omitempty"`
}
