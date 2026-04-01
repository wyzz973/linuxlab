package sandbox

import (
	"context"
	"os"
	"os/exec"
)

// LocalSandbox executes commands locally in a temporary directory.
// It serves as a degraded-mode fallback when Docker is unavailable.
type LocalSandbox struct {
	workDir string
}

// NewLocalSandbox creates a new local sandbox with a temporary working directory.
func NewLocalSandbox() (*LocalSandbox, error) {
	dir, err := os.MkdirTemp("", "linuxlab-*")
	if err != nil {
		return nil, err
	}
	return &LocalSandbox{workDir: dir}, nil
}

// WorkDir returns the temporary working directory path.
func (s *LocalSandbox) WorkDir() string { return s.workDir }

// Exec runs a command in the sandbox working directory via sh -c.
func (s *LocalSandbox) Exec(ctx context.Context, command string) (string, int, error) {
	cmd := exec.CommandContext(ctx, "sh", "-c", command)
	cmd.Dir = s.workDir
	return runAndCapture(cmd)
}

// Destroy removes the temporary working directory.
func (s *LocalSandbox) Destroy(_ context.Context) error {
	return os.RemoveAll(s.workDir)
}

// InteractiveShellArgs returns arguments to open an interactive shell in the working directory.
func (s *LocalSandbox) InteractiveShellArgs() []string {
	return []string{"bash", "--norc", "-c", "cd '" + s.workDir + "' && exec bash"}
}
