package sandbox

import (
	"context"
	"errors"
	"os/exec"
	"strings"
	"testing"
	"time"
)

func dockerAvailable() bool {
	cmd := exec.Command("docker", "info")
	return cmd.Run() == nil
}

func TestDrainPullStream(t *testing.T) {
	tests := []struct {
		name    string
		stream  string
		wantErr string
	}{
		{
			name:   "clean stream",
			stream: `{"status":"Pulling from library/ubuntu"}{"status":"Download complete"}`,
		},
		{
			name:   "empty stream",
			stream: ``,
		},
		{
			name:    "in-stream error is surfaced",
			stream:  `{"status":"Pulling"}{"error":"write /var/lib/docker: no space left on device"}`,
			wantErr: "no space left on device",
		},
		{
			name:    "corrupted stream is surfaced",
			stream:  `{"status":"Pul`,
			wantErr: "decode pull progress",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := drainPullStream(strings.NewReader(tt.stream))
			if tt.wantErr == "" {
				if err != nil {
					t.Fatalf("drainPullStream: %v", err)
				}
				return
			}
			if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
				t.Fatalf("err = %v, want error containing %q", err, tt.wantErr)
			}
		})
	}
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

	// Regression (sleep 3600): PID 1 must be an indefinite sleep so the
	// container cannot self-terminate mid-session.
	out, _, err := sb.Exec(ctx, "tr '\\0' ' ' < /proc/1/cmdline")
	if err != nil {
		t.Fatalf("Exec failed: %v", err)
	}
	if !strings.Contains(out, "infinity") {
		t.Errorf("PID 1 cmdline = %q, want sleep infinity", out)
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

func TestDockerSandbox_ExecMergesStderr(t *testing.T) {
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

	// Regression: stderr used to be discarded, diverging from the
	// CombinedOutput semantics of LocalSandbox/ComposeSandbox.
	output, exitCode, err := sb.Exec(ctx, "echo to-stdout; echo to-stderr 1>&2")
	if err != nil {
		t.Fatalf("Exec failed: %v", err)
	}
	if exitCode != 0 {
		t.Errorf("exit code = %d, want 0", exitCode)
	}
	if !strings.Contains(output, "to-stdout") || !strings.Contains(output, "to-stderr") {
		t.Errorf("output = %q, want both stdout and stderr content", output)
	}
}

func TestDockerSandbox_ExecContextTimeout(t *testing.T) {
	if !dockerAvailable() {
		t.Skip("docker not available")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	sb, err := NewDockerSandbox(ctx, "")
	if err != nil {
		t.Fatalf("NewDockerSandbox failed: %v", err)
	}
	defer sb.Destroy(context.Background())

	// Regression: StdCopy on the hijacked connection used to ignore ctx,
	// so a blocking foreground command hung Exec for its full duration.
	execCtx, execCancel := context.WithTimeout(ctx, 2*time.Second)
	defer execCancel()

	start := time.Now()
	_, _, err = sb.Exec(execCtx, "sleep 30")
	elapsed := time.Since(start)

	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("err = %v, want context.DeadlineExceeded", err)
	}
	if elapsed > 10*time.Second {
		t.Errorf("Exec returned after %v, want prompt return on context timeout", elapsed)
	}
}

func TestDockerSandbox_DestroyWithDeadContext(t *testing.T) {
	if !dockerAvailable() {
		t.Skip("docker not available")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	sb, err := NewDockerSandbox(ctx, "")
	if err != nil {
		t.Fatalf("NewDockerSandbox failed: %v", err)
	}

	// Regression: Destroy used to run ContainerRemove on the caller's ctx,
	// so an expired ctx leaked the container.
	deadCtx, deadCancel := context.WithCancel(context.Background())
	deadCancel()
	if err := sb.Destroy(deadCtx); err != nil {
		t.Errorf("Destroy with cancelled context should still remove the container: %v", err)
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
