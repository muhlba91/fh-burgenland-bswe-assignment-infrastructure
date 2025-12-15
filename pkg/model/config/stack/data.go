package stack

import (
	"github.com/muhlba91/fh-burgenland-bswe-assignment-infrastructure/pkg/model/config/repository"
	"github.com/muhlba91/fh-burgenland-bswe-assignment-infrastructure/pkg/model/config/team"
)

// Config defines stack-related configuration.
type Config struct {
	// Repositories defines the repositories to be created in the stack.
	Repositories []*repository.Config `yaml:"repositories,omitempty"`
	// Teams defines the teams to be created in the stack.
	Teams []*team.Config `yaml:"teams,omitempty"`
}
