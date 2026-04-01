package reference

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// LoadReferences reads and parses a command reference YAML file.
func LoadReferences(path string) (*ReferenceData, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read references %s: %w", path, err)
	}

	var refs ReferenceData
	if err := yaml.Unmarshal(data, &refs); err != nil {
		return nil, fmt.Errorf("parse references %s: %w", path, err)
	}

	return &refs, nil
}
