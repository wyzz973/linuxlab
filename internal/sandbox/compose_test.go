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

func writeComposeFile(t *testing.T, content string) (dir, name string) {
	t.Helper()
	dir = t.TempDir()
	name = "docker-compose.yaml"
	if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	return dir, name
}

func TestComposeSandbox_ImplementsSandbox(t *testing.T) {
	var _ Sandbox = (*ComposeSandbox)(nil)
}

func TestFirstDeclaredService(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    string
		wantErr bool
	}{
		{
			// Regression: `docker compose config --services` sorts
			// alphabetically and would return "cache" here.
			name: "declaration order beats alphabetical order",
			content: `services:
  web:
    image: nginx:alpine
  cache:
    image: redis:7-alpine
`,
			want: "web",
		},
		{
			name: "single service",
			content: `services:
  app:
    image: alpine:3.18
`,
			want: "app",
		},
		{
			name: "services after other top-level keys",
			content: `version: "3"
services:
  zeta:
    image: alpine:3.18
  alpha:
    image: alpine:3.18
`,
			want: "zeta",
		},
		{
			name:    "no services key",
			content: "volumes: {}\n",
			wantErr: true,
		},
		{
			name: "empty services mapping",
			content: `services: {}
`,
			wantErr: true,
		},
		{
			name:    "invalid yaml",
			content: "services: [unclosed\n",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir, name := writeComposeFile(t, tt.content)
			got, err := firstDeclaredService(filepath.Join(dir, name))
			if tt.wantErr {
				if err == nil {
					t.Fatalf("firstDeclaredService = %q, want error", got)
				}
				return
			}
			if err != nil {
				t.Fatalf("firstDeclaredService: %v", err)
			}
			if got != tt.want {
				t.Errorf("firstDeclaredService = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestFirstDeclaredService_MissingFile(t *testing.T) {
	if _, err := firstDeclaredService(filepath.Join(t.TempDir(), "nope.yaml")); err == nil {
		t.Fatal("expected error for missing compose file")
	}
}

func TestComposeSandbox_InteractiveShellArgs_FallsBackToSh(t *testing.T) {
	sb := &ComposeSandbox{dir: "/x", composeFile: "/x/docker-compose.yaml", service: "web"}
	joined := strings.Join(sb.InteractiveShellArgs(), " ")
	// Regression: /bin/bash was hardcoded, which fails outright on
	// bash-less images (nginx:alpine, redis:7-alpine, ...).
	if strings.Contains(joined, "/bin/bash") {
		t.Errorf("interactive shell must not hardcode /bin/bash: %q", joined)
	}
	if !strings.Contains(joined, "exec sh") {
		t.Errorf("interactive shell should fall back to sh: %q", joined)
	}
}

func TestComposeSandbox_UpAndDown(t *testing.T) {
	if !composeAvailable() {
		t.Skip("docker not available")
	}

	dir, name := writeComposeFile(t, `services:
  web:
    image: alpine:3.18
    command: ["sleep", "300"]
`)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	sb, err := NewComposeSandbox(ctx, dir, name)
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
		t.Skip("docker not available")
	}

	dir, name := writeComposeFile(t, `services:
  web:
    image: alpine:3.18
    command: ["sleep", "300"]
`)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	sb, err := NewComposeSandbox(ctx, dir, name)
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

func TestComposeSandbox_PrimaryServiceIsFirstDeclared(t *testing.T) {
	if !composeAvailable() {
		t.Skip("docker not available")
	}

	// "web" is declared first but sorts after "cache" — the old
	// `config --services`-based detection picked "cache" here.
	dir, name := writeComposeFile(t, `services:
  web:
    image: alpine:3.18
    command: ["sleep", "300"]
    environment:
      - LINUXLAB_ROLE=web
  cache:
    image: alpine:3.18
    command: ["sleep", "300"]
    environment:
      - LINUXLAB_ROLE=cache
`)

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	sb, err := NewComposeSandbox(ctx, dir, name)
	if err != nil {
		t.Fatalf("NewComposeSandbox: %v", err)
	}
	defer sb.Destroy(ctx)

	out, code, err := sb.Exec(ctx, "echo $LINUXLAB_ROLE")
	if err != nil {
		t.Fatalf("Exec: %v", err)
	}
	if code != 0 {
		t.Errorf("exit code = %d, want 0", code)
	}
	if strings.TrimSpace(out) != "web" {
		t.Errorf("primary service role = %q, want %q (first declared)", strings.TrimSpace(out), "web")
	}
}
