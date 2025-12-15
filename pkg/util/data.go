package util //nolint:revive // package name is util

import (
	"os"

	"github.com/muhlba91/fh-burgenland-bswe-assignment-infrastructure/pkg/model/config/stack"

	"gopkg.in/yaml.v3"
)

// ParseDataFromFiles reads the configuration file from the specified path.
func ParseDataFromFiles(path string) (*stack.Config, error) {
	b, rErr := os.ReadFile(path)
	if rErr != nil {
		return nil, rErr
	}

	var config stack.Config
	if yErr := yaml.Unmarshal(b, &config); yErr != nil {
		return nil, yErr
	}

	return &config, nil
}
