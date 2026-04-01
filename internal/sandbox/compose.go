package sandbox

import (
	"context"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
)

// ComposeSandbox manages a multi-service environment via docker compose.
type ComposeSandbox struct {
	dir     string // directory containing docker-compose.yaml
	service string // primary service name for exec
}

// NewComposeSandbox starts services defined in a docker-compose.yaml in the given directory.
// It auto-detects the first service name for exec commands.
func NewComposeSandbox(ctx context.Context, dir string) (*ComposeSandbox, error) {
	composeFile := filepath.Join(dir, "docker-compose.yaml")

	// Bring up services
	cmd := exec.CommandContext(ctx, "docker", "compose", "-f", composeFile, "up", "-d")
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("docker compose up: %s: %w", string(out), err)
	}

	// Detect the first service name
	svcCmd := exec.CommandContext(ctx, "docker", "compose", "-f", composeFile, "config", "--services")
	svcCmd.Dir = dir
	svcOut, err := svcCmd.Output()
	if err != nil {
		// Clean up on failure
		exec.CommandContext(ctx, "docker", "compose", "-f", composeFile, "down").Run()
		return nil, fmt.Errorf("docker compose config --services: %w", err)
	}

	services := strings.Split(strings.TrimSpace(string(svcOut)), "\n")
	if len(services) == 0 || services[0] == "" {
		exec.CommandContext(ctx, "docker", "compose", "-f", composeFile, "down").Run()
		return nil, fmt.Errorf("no services found in compose file")
	}

	return &ComposeSandbox{dir: dir, service: services[0]}, nil
}

// Exec runs a command in the primary service container.
func (s *ComposeSandbox) Exec(ctx context.Context, command string) (string, int, error) {
	composeFile := filepath.Join(s.dir, "docker-compose.yaml")
	cmd := exec.CommandContext(ctx, "docker", "compose", "-f", composeFile, "exec", "-T", s.service, "sh", "-c", command)
	cmd.Dir = s.dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return string(out), exitErr.ExitCode(), nil
		}
		return string(out), -1, err
	}
	return string(out), 0, nil
}

// Destroy tears down all compose services.
func (s *ComposeSandbox) Destroy(ctx context.Context) error {
	composeFile := filepath.Join(s.dir, "docker-compose.yaml")
	cmd := exec.CommandContext(ctx, "docker", "compose", "-f", composeFile, "down")
	cmd.Dir = s.dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("docker compose down: %s: %w", string(out), err)
	}
	return nil
}

// InteractiveShellArgs returns arguments to open an interactive shell in the primary service.
func (s *ComposeSandbox) InteractiveShellArgs() []string {
	composeFile := filepath.Join(s.dir, "docker-compose.yaml")
	return []string{"docker", "compose", "-f", composeFile, "exec", s.service, "/bin/sh"}
}
