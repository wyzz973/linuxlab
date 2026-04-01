package sandbox

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func composeAvailable() bool {
	// Check both that docker compose CLI exists and that Docker daemon is running
	return exec.Command("docker", "compose", "version").Run() == nil &&
		exec.Command("docker", "info").Run() == nil
}

func TestComposeSandbox_ImplementsSandbox(t *testing.T) {
	var _ Sandbox = (*ComposeSandbox)(nil)
}

func TestComposeSandbox_UpAndDown(t *testing.T) {
	if !composeAvailable() {
		t.Skip("docker compose not available")
	}

	dir := t.TempDir()
	composeFile := filepath.Join(dir, "docker-compose.yaml")
	content := `services:
  web:
    image: alpine:3.18
    command: ["sleep", "300"]
`
	if err := os.WriteFile(composeFile, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	sb, err := NewComposeSandbox(ctx, dir)
	if err != nil {
		t.Fatalf("NewComposeSandbox: %v", err)
	}

	// Verify we can destroy without error
	if err := sb.Destroy(ctx); err != nil {
		t.Errorf("Destroy: %v", err)
	}
}

func TestComposeSandbox_Exec(t *testing.T) {
	if !composeAvailable() {
		t.Skip("docker compose not available")
	}

	dir := t.TempDir()
	composeFile := filepath.Join(dir, "docker-compose.yaml")
	content := `services:
  web:
    image: alpine:3.18
    command: ["sleep", "300"]
`
	if err := os.WriteFile(composeFile, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	sb, err := NewComposeSandbox(ctx, dir)
	if err != nil {
		t.Fatalf("NewComposeSandbox: %v", err)
	}
	defer sb.Destroy(ctx)

	out, code, err := sb.Exec(ctx, "echo hello")
	if err != nil {
		t.Fatalf("Exec: %v", err)
	}
	if code != 0 {
		t.Errorf("exit code = %d, want 0", code)
	}
	if strings.TrimSpace(out) != "hello" {
		t.Errorf("output = %q, want %q", out, "hello")
	}
}
