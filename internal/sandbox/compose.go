package sandbox

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// ComposeSandbox manages a multi-service environment via docker compose.
type ComposeSandbox struct {
	dir         string // directory containing docker-compose.yaml
	composeFile string // full path to docker-compose.yaml
	service     string // primary service name for exec
}

// firstDeclaredService parses the compose file and returns the first service
// in declaration order. `docker compose config --services` cannot be used
// for this: compose v2 sorts its output alphabetically, so its first line is
// not the first declared service (e.g. a file declaring web before cache
// would yield "cache").
func firstDeclaredService(composeFile string) (string, error) {
	data, err := os.ReadFile(composeFile)
	if err != nil {
		return "", fmt.Errorf("read compose file: %w", err)
	}
	var doc yaml.Node
	if err := yaml.Unmarshal(data, &doc); err != nil {
		return "", fmt.Errorf("parse compose file: %w", err)
	}
	if len(doc.Content) == 0 || doc.Content[0].Kind != yaml.MappingNode {
		return "", fmt.Errorf("no services found in compose file")
	}
	root := doc.Content[0]
	// Mapping nodes store key/value pairs as alternating entries in Content,
	// preserving document order.
	for i := 0; i+1 < len(root.Content); i += 2 {
		if root.Content[i].Value != "services" {
			continue
		}
		services := root.Content[i+1]
		if services.Kind == yaml.MappingNode && len(services.Content) > 0 {
			return services.Content[0].Value, nil
		}
		break
	}
	return "", fmt.Errorf("no services found in compose file")
}

// NewComposeSandbox starts services defined in a compose file in the given directory.
// The first service declared in the compose file becomes the primary service
// used for exec commands and the interactive shell.
func NewComposeSandbox(ctx context.Context, dir, composeFileName string) (*ComposeSandbox, error) {
	if composeFileName == "" {
		composeFileName = "docker-compose.yaml"
	}
	composeFile := filepath.Join(dir, composeFileName)

	// Detect the primary service before starting anything, so a bad compose
	// file fails fast with nothing to clean up.
	service, err := firstDeclaredService(composeFile)
	if err != nil {
		return nil, err
	}

	// Bring up services
	cmd := exec.CommandContext(ctx, "docker", "compose", "-f", composeFile, "up", "-d")
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("docker compose up: %s: %w", string(out), err)
	}

	return &ComposeSandbox{dir: dir, composeFile: composeFile, service: service}, nil
}

// Exec runs a command in the primary service container.
func (s *ComposeSandbox) Exec(ctx context.Context, command string) (string, int, error) {
	cmd := exec.CommandContext(ctx, "docker", "compose", "-f", s.composeFile, "exec", "-T", s.service, "sh", "-c", command)
	cmd.Dir = s.dir
	return runAndCapture(cmd)
}

// Destroy tears down all compose services.
// Teardown runs on an independent timeout context so it still succeeds when
// the caller's context has already expired or been cancelled — otherwise a
// timed-out challenge would leak its services.
func (s *ComposeSandbox) Destroy(_ context.Context) error {
	cleanupCtx, cancel := context.WithTimeout(context.Background(), cleanupTimeout)
	defer cancel()
	cmd := exec.CommandContext(cleanupCtx, "docker", "compose", "-f", s.composeFile, "down")
	cmd.Dir = s.dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("docker compose down: %s: %w", string(out), err)
	}
	return nil
}

// InteractiveShellArgs returns arguments to open an interactive shell in the primary service.
// bash is preferred, but images without it (e.g. alpine-based ones) fall back to sh.
func (s *ComposeSandbox) InteractiveShellArgs() []string {
	return []string{
		"docker", "compose", "-f", s.composeFile, "exec", "-it", s.service,
		"sh", "-c", "command -v bash >/dev/null 2>&1 && exec bash --rcfile /tmp/.linuxlab_bashrc || exec sh",
	}
}
