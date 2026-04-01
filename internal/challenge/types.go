package challenge

// Challenge represents a single learning challenge (Linux or Vim).
type Challenge struct {
	ID          string       `yaml:"id"`
	Title       string       `yaml:"title"`
	Difficulty  int          `yaml:"difficulty"`
	Category    string       `yaml:"category"`
	Subcategory string       `yaml:"subcategory"`
	Tags        []string     `yaml:"tags"`
	Description string       `yaml:"description"`
	Hints       []Hint       `yaml:"hints,omitempty"`
	Verify      []VerifyRule `yaml:"verify"`
	SetupFiles  []SetupFile  `yaml:"setup_files,omitempty"`
	ComposeFile    string       `yaml:"compose_file,omitempty"`
	RequiresDocker bool         `yaml:"requires_docker,omitempty"`
	Dir            string       `yaml:"-"`
}

// Hint provides a progressive hint for solving a challenge.
type Hint struct {
	Level int    `yaml:"level"`
	Text  string `yaml:"text"`
}

// VerifyRule defines how to check if a challenge is completed.
type VerifyRule struct {
	Type    string `yaml:"type"`
	Path    string `yaml:"path,omitempty"`
	Expect  string `yaml:"expect,omitempty"`
	Command string `yaml:"command,omitempty"`
}

// SetupFile defines a file to create before starting a challenge.
type SetupFile struct {
	Path    string `yaml:"path"`
	Content string `yaml:"content"`
}
