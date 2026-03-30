package sandbox

import (
	"context"
	"os/exec"
	"testing"
	"time"
)

func dockerAvailable() bool {
	cmd := exec.Command("docker", "info")
	return cmd.Run() == nil
}

func TestDockerSandbox_CreateAndDestroy(t *testing.T) {
	if !dockerAvailable() {
		t.Skip("docker not available")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	sb, err := NewDockerSandbox(ctx, "")
	if err != nil {
		t.Fatalf("NewDockerSandbox failed: %v", err)
	}
	defer sb.Destroy(ctx)

	if sb.ContainerID() == "" {
		t.Error("container ID is empty")
	}
}

func TestDockerSandbox_ExecCommand(t *testing.T) {
	if !dockerAvailable() {
		t.Skip("docker not available")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	sb, err := NewDockerSandbox(ctx, "")
	if err != nil {
		t.Fatalf("NewDockerSandbox failed: %v", err)
	}
	defer sb.Destroy(ctx)

	output, exitCode, err := sb.Exec(ctx, "echo hello")
	if err != nil {
		t.Fatalf("Exec failed: %v", err)
	}
	if exitCode != 0 {
		t.Errorf("exit code = %d, want 0", exitCode)
	}
	if output != "hello\n" {
		t.Errorf("output = %q, want %q", output, "hello\n")
	}
}

func TestDockerSandbox_RunInitScript(t *testing.T) {
	if !dockerAvailable() {
		t.Skip("docker not available")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	sb, err := NewDockerSandbox(ctx, "")
	if err != nil {
		t.Fatalf("NewDockerSandbox failed: %v", err)
	}
	defer sb.Destroy(ctx)

	_, exitCode, err := sb.Exec(ctx, "mkdir -p /tmp/testdir && echo done > /tmp/testdir/result.txt")
	if err != nil {
		t.Fatalf("Exec failed: %v", err)
	}
	if exitCode != 0 {
		t.Errorf("init exit code = %d", exitCode)
	}

	output, _, _ := sb.Exec(ctx, "cat /tmp/testdir/result.txt")
	if output != "done\n" {
		t.Errorf("result = %q, want %q", output, "done\n")
	}
}
