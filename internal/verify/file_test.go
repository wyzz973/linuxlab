package verify

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/sd3/linuxlab/internal/challenge"
)

func TestFileContentVerifier_Pass(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "test.txt")
	os.WriteFile(p, []byte("hello\nworld\n"), 0644)

	v := &FileContentVerifier{}
	r := v.Verify(challenge.VerifyRule{Type: "file_content", Path: p, Expect: "hello\nworld\n"})
	if !r.Passed {
		t.Fatalf("expected pass, got fail: %s", r.Message)
	}
}

func TestFileContentVerifier_Fail_WrongContent(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "test.txt")
	os.WriteFile(p, []byte("wrong\n"), 0644)

	v := &FileContentVerifier{}
	r := v.Verify(challenge.VerifyRule{Type: "file_content", Path: p, Expect: "hello\nworld\n"})
	if r.Passed {
		t.Fatal("expected fail, got pass")
	}
}

func TestFileContentVerifier_Fail_NoFile(t *testing.T) {
	v := &FileContentVerifier{}
	r := v.Verify(challenge.VerifyRule{Type: "file_content", Path: "/nonexistent/path/file.txt", Expect: "hello"})
	if r.Passed {
		t.Fatal("expected fail, got pass")
	}
}

func TestFileContentVerifier_TrimWhitespace(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "test.txt")
	os.WriteFile(p, []byte("  hello  \n"), 0644)

	v := &FileContentVerifier{}
	r := v.Verify(challenge.VerifyRule{Type: "file_content", Path: p, Expect: "hello"})
	if !r.Passed {
		t.Fatalf("expected pass with trimmed whitespace, got fail: %s", r.Message)
	}
}

func TestFileExistsVerifier_Pass(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "test.txt")
	os.WriteFile(p, []byte(""), 0644)

	v := &FileExistsVerifier{}
	r := v.Verify(challenge.VerifyRule{Type: "file_exists", Path: p})
	if !r.Passed {
		t.Fatalf("expected pass, got fail: %s", r.Message)
	}
}

func TestFileExistsVerifier_Fail(t *testing.T) {
	v := &FileExistsVerifier{}
	r := v.Verify(challenge.VerifyRule{Type: "file_exists", Path: "/nonexistent/path/file.txt"})
	if r.Passed {
		t.Fatal("expected fail, got pass")
	}
}

func TestFileExistsVerifier_Directory(t *testing.T) {
	dir := t.TempDir()

	v := &FileExistsVerifier{}
	r := v.Verify(challenge.VerifyRule{Type: "file_exists", Path: dir})
	if !r.Passed {
		t.Fatalf("expected pass for directory, got fail: %s", r.Message)
	}
}
