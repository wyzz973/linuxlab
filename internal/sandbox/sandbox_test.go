package sandbox

import (
	"context"
	"testing"

	"github.com/sd3/linuxlab/internal/challenge"
)

func TestDockerSandbox_ImplementsSandbox(t *testing.T) {
	var _ Sandbox = (*DockerSandbox)(nil) // compile-time check
}

func TestDockerAvailable(t *testing.T) {
	// Just verify it returns a bool without panicking
	result := DockerAvailable()
	t.Logf("DockerAvailable() = %v", result)
}

func TestNewSandbox_LocalFallback(t *testing.T) {
	if DockerAvailable() {
		t.Skip("Docker is available; cannot test local fallback")
	}
	ch := &challenge.Challenge{
		ID:       "test-challenge",
		Category: "linux-basics",
	}
	sb, err := NewSandbox(context.Background(), ch)
	if err != nil {
		t.Fatalf("NewSandbox: %v", err)
	}
	defer sb.Destroy(context.Background())

	if _, ok := sb.(*LocalSandbox); !ok {
		t.Errorf("expected *LocalSandbox, got %T", sb)
	}
}
