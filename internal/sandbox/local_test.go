package sandbox

import (
	"context"
	"os"
	"strings"
	"testing"
)

func TestLocalSandbox_ImplementsSandbox(t *testing.T) {
	var _ Sandbox = (*LocalSandbox)(nil)
}

func TestLocalSandbox_CreateAndDestroy(t *testing.T) {
	sb := newLocalSandboxWithDir(t.TempDir())
	if _, err := os.Stat(sb.WorkDir()); err != nil {
		t.Errorf("workDir doesn't exist: %v", err)
	}
}

func TestLocalSandbox_Exec(t *testing.T) {
	sb := newLocalSandboxWithDir(t.TempDir())
	out, code, err := sb.Exec(context.Background(), "echo hello")
	if err != nil {
		t.Fatal(err)
	}
	if code != 0 {
		t.Errorf("code = %d", code)
	}
	if strings.TrimSpace(out) != "hello" {
		t.Errorf("out = %q", out)
	}
}

func TestLocalSandbox_ExecInWorkDir(t *testing.T) {
	sb := newLocalSandboxWithDir(t.TempDir())
	sb.Exec(context.Background(), "touch testfile.txt")
	out, _, _ := sb.Exec(context.Background(), "ls testfile.txt")
	if strings.TrimSpace(out) != "testfile.txt" {
		t.Errorf("out = %q", out)
	}
}

func TestLocalSandbox_Destroy(t *testing.T) {
	// Use NewLocalSandbox to test the real constructor + Destroy cleanup
	sb, err := NewLocalSandbox()
	if err != nil {
		t.Fatal(err)
	}
	dir := sb.WorkDir()
	sb.Destroy(context.Background())
	if _, err := os.Stat(dir); !os.IsNotExist(err) {
		t.Error("workDir should be removed after Destroy")
	}
}
