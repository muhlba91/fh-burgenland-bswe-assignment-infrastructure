package repository

// TeamConfig defines configuration for a team associated with a repository.
type TeamConfig struct {
	// Name is the name of the team.
	Name string `yaml:"name"`
	// Role is the role of the team in the repository.
	Role string `yaml:"role"`
}
