package reference

// CommandRef represents a single command reference entry.
type CommandRef struct {
	Name              string    `yaml:"name"`
	Brief             string    `yaml:"brief"`
	Examples          []Example `yaml:"examples"`
	RelatedChallenges []string  `yaml:"related_challenges,omitempty"`
}

// Example represents a usage example for a command.
type Example struct {
	Desc string `yaml:"desc"`
	Cmd  string `yaml:"cmd"`
}

// ReferenceData holds all command references loaded from YAML.
type ReferenceData struct {
	Commands []CommandRef `yaml:"commands"`
}
