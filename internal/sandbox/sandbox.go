package sandbox

import (
	"context"
	"os/exec"

	"github.com/sd3/linuxlab/internal/challenge"
)

// Sandbox is the common interface for all challenge execution backends.
type Sandbox interface {
	Exec(ctx context.Context, command string) (string, int, error)
	Destroy(ctx context.Context) error
	InteractiveShellArgs() []string
}

// DockerAvailable returns true if the Docker daemon is reachable.
func DockerAvailable() bool {
	return exec.Command("docker", "info").Run() == nil
}

// runAndCapture runs a command and returns output, exit code, and error.
// Shared by LocalSandbox and ComposeSandbox to avoid duplicating ExitError handling.
func runAndCapture(cmd *exec.Cmd) (string, int, error) {
	out, err := cmd.CombinedOutput()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return string(out), exitErr.ExitCode(), nil
		}
		return string(out), -1, err
	}
	return string(out), 0, nil
}

// NewSandbox creates the appropriate sandbox backend for a challenge.
// It picks ComposeSandbox if a compose file is specified, DockerSandbox
// if Docker is available, or LocalSandbox as a degraded-mode fallback.
func NewSandbox(ctx context.Context, ch *challenge.Challenge) (Sandbox, error) {
	hasDocker := DockerAvailable()
	if ch.ComposeFile != "" && hasDocker {
		return NewComposeSandbox(ctx, ch.Dir)
	}
	if hasDocker {
		return NewDockerSandbox(ctx, "")
	}
	return NewLocalSandbox()
}
