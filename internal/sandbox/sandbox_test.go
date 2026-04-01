package sandbox

import (
	"testing"
)

func TestDockerSandbox_ImplementsSandbox(t *testing.T) {
	var _ Sandbox = (*DockerSandbox)(nil) // compile-time check
}

func TestLocalSandboxPlaceholder(t *testing.T) {
	// Will be implemented in Task 3
	t.Skip("LocalSandbox not yet implemented")
}
