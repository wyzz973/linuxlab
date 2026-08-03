package sandbox

import (
	"context"
	"os/exec"
	"sync"
	"time"

	"github.com/sd3/linuxlab/internal/challenge"
)

// Sandbox is the common interface for all challenge execution backends.
type Sandbox interface {
	Exec(ctx context.Context, command string) (string, int, error)
	Destroy(ctx context.Context) error
	InteractiveShellArgs() []string
}

// cleanupTimeout bounds container/compose teardown. Cleanup deliberately runs
// on a context independent from the caller's, so that an expired or cancelled
// business context cannot prevent resources from being released (see the
// Destroy implementations).
const cleanupTimeout = 15 * time.Second

// Docker availability probing. The result rarely changes within a session,
// but the probe costs a subprocess round-trip (and hangs when the daemon is
// wedged), so it is cached in-process. A successful probe is trusted for a
// long time; a failed one only briefly, so that starting the Docker daemon
// mid-session is picked up quickly.
const (
	dockerProbeTimeout    = 3 * time.Second
	dockerProbeSuccessTTL = 5 * time.Minute
	dockerProbeFailureTTL = 5 * time.Second
)

var (
	dockerProbeMu     sync.Mutex
	dockerProbeResult bool
	dockerProbeTime   time.Time
)

// DockerAvailable returns true if the Docker daemon is reachable.
// The probe runs with a short timeout so a wedged daemon cannot block the
// caller indefinitely (the TUI calls this from its event loop), and the
// result is cached in-process (see the TTL constants above).
// `docker version --format` is used instead of `docker info` because it only
// needs a single lightweight API call (measured ~4x faster).
func DockerAvailable() bool {
	dockerProbeMu.Lock()
	defer dockerProbeMu.Unlock()

	ttl := dockerProbeFailureTTL
	if dockerProbeResult {
		ttl = dockerProbeSuccessTTL
	}
	if !dockerProbeTime.IsZero() && time.Since(dockerProbeTime) < ttl {
		return dockerProbeResult
	}

	ctx, cancel := context.WithTimeout(context.Background(), dockerProbeTimeout)
	defer cancel()
	err := exec.CommandContext(ctx, "docker", "version", "--format", "{{.Server.Version}}").Run()
	dockerProbeResult = err == nil
	dockerProbeTime = time.Now()
	return dockerProbeResult
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
// It picks ComposeSandbox if a compose file is specified and Docker is available,
// DockerSandbox if Docker is available, or LocalSandbox as a degraded-mode fallback.
func NewSandbox(ctx context.Context, ch *challenge.Challenge) (Sandbox, error) {
	// containers-category challenges without a compose file exercise the
	// host `docker` CLI itself (their check.sh runs `docker ps`,
	// `docker images`, ...). Inside the ubuntu sandbox container there is
	// neither a docker binary nor a daemon socket, so both the exercise
	// and its verification would inevitably fail there — route these
	// challenges to LocalSandbox, where the host Docker daemon is reachable.
	if ch.Category == "containers" && ch.ComposeFile == "" {
		return NewLocalSandbox()
	}
	hasDocker := DockerAvailable()
	if ch.ComposeFile != "" && hasDocker {
		return NewComposeSandbox(ctx, ch.Dir, ch.ComposeFile)
	}
	if hasDocker {
		return NewDockerSandbox(ctx, "")
	}
	return NewLocalSandbox()
}
