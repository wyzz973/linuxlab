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

// NewSandbox creates the appropriate sandbox backend for a challenge.
// It picks ComposeSandbox if a compose file is specified, DockerSandbox
// if Docker is available, or LocalSandbox as a degraded-mode fallback.
func NewSandbox(ctx context.Context, ch *challenge.Challenge) (Sandbox, error) {
	if ch.ComposeFile != "" && DockerAvailable() {
		return NewComposeSandbox(ctx, ch.Dir)
	}
	if DockerAvailable() {
		return NewDockerSandbox(ctx, "")
	}
	return NewLocalSandbox()
}
