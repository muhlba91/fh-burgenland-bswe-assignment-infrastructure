package repository

// Config defines Repository-related configuration.
type Config struct {
	// Name is the name of the repository.
	Name string `yaml:"name"`
	// Service is the intended purpose of the repository.
	Service string `yaml:"service"`
	// Teams contains configuration for teams associated with the repository.
	Teams []*TeamConfig `yaml:"teams,omitempty"`
	// Approvers is the number of required approvers for pull requests.
	Approvers *int `yaml:"approvers,omitempty"`
	// DeleteOnDestroy indicates whether the repository should be deleted when destroyed.
	DeleteOnDestroy *bool `yaml:"deleteOnDestroy,omitempty"`
	// AWS indicates whether AWS-related configurations are enabled.
	AWS *bool `yaml:"aws,omitempty"`
	// Terraform indicates whether Terraform-related configurations are enabled.
	Terraform *bool `yaml:"terraform,omitempty"`
	// Pulumi indicates whether Pulumi-related configurations are enabled.
	Pulumi *bool `yaml:"pulumi,omitempty"`
	// Harbor indicates whether Harbor-related configurations are enabled.
	Harbor *bool `yaml:"harbor,omitempty"`
	// RequiredChecks contains a list of required checks for pull requests.
	RequiredChecks []string `yaml:"requiredChecks,omitempty"`
}
