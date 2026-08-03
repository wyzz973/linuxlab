package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/sd3/linuxlab/internal/runner"
	"github.com/sd3/linuxlab/internal/verify"
)

func TestUnknownCommandReturnsChineseError(t *testing.T) {
	var stderr bytes.Buffer
	code := Run(context.Background(), Env{Args: []string{"unknown"}, Stderr: &stderr})
	if code == 0 {
		t.Fatal("expected non-zero exit")
	}
	if !strings.Contains(stderr.String(), "未知命令") {
		t.Fatalf("expected Chinese error, got %q", stderr.String())
	}
}

func TestDoctorJSONReturnsParseableJSON(t *testing.T) {
	root := t.TempDir()
	writeChallenge(t, root, "linux-basics", "ls-basic")
	refsPath := filepath.Join(root, "commands.yaml")
	if err := os.WriteFile(refsPath, []byte("commands: []\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	var stdout bytes.Buffer
	code := Run(context.Background(), Env{
		Args:          []string{"doctor", "--json"},
		Stdout:        &stdout,
		Stderr:        ioDiscard{},
		ChallengesDir: filepath.Join(root, "challenges"),
		RefsPath:      refsPath,
	})
	if code != 0 {
		t.Fatalf("doctor returned %d", code)
	}
	var parsed map[string]any
	if err := json.Unmarshal(stdout.Bytes(), &parsed); err != nil {
		t.Fatalf("invalid json: %v: %s", err, stdout.String())
	}
}

func TestDoctorJSONReportsLoadErrorsWithNonZeroExit(t *testing.T) {
	root := t.TempDir()
	var stdout bytes.Buffer
	code := Run(context.Background(), Env{
		Args:          []string{"doctor", "--json"},
		Stdout:        &stdout,
		Stderr:        ioDiscard{},
		ChallengesDir: filepath.Join(root, "missing-challenges"),
		RefsPath:      filepath.Join(root, "missing-refs.yaml"),
	})
	if code == 0 {
		t.Fatal("doctor should return non-zero when loading fails")
	}
	var parsed map[string]any
	if err := json.Unmarshal(stdout.Bytes(), &parsed); err != nil {
		t.Fatalf("invalid json: %v: %s", err, stdout.String())
	}
	if s, _ := parsed["challenges_error"].(string); s == "" {
		t.Errorf("expected challenges_error in JSON output, got %v", parsed)
	}
	if s, _ := parsed["references_error"].(string); s == "" {
		t.Errorf("expected references_error in JSON output, got %v", parsed)
	}
}

func TestDoctorTextModeShowsLoadErrors(t *testing.T) {
	root := t.TempDir()
	var stdout bytes.Buffer
	code := Run(context.Background(), Env{
		Args:          []string{"doctor"},
		Stdout:        &stdout,
		Stderr:        ioDiscard{},
		ChallengesDir: filepath.Join(root, "missing-challenges"),
		RefsPath:      filepath.Join(root, "missing-refs.yaml"),
	})
	if code == 0 {
		t.Fatal("doctor should return non-zero when loading fails")
	}
	if !strings.Contains(stdout.String(), "加载失败") {
		t.Fatalf("expected load error details in text output, got %q", stdout.String())
	}
}

func TestRenderEventResultOutput(t *testing.T) {
	failed := false
	var buf bytes.Buffer
	renderEvent(&buf, runner.Event{
		Type:   "result",
		Passed: &failed,
		Results: []verify.Result{
			{Passed: true, Message: "文件内容匹配"},
			{Passed: false, Message: "命令输出不匹配"},
		},
		HintsUsed: 1,
	})
	out := buf.String()
	if !strings.Contains(out, "未通过") {
		t.Errorf("expected 未通过 in output, got %q", out)
	}
	if !strings.Contains(out, "文件内容匹配") || !strings.Contains(out, "命令输出不匹配") {
		t.Errorf("expected every check message in output, got %q", out)
	}
	if !strings.Contains(out, "使用提示: 1") {
		t.Errorf("expected hint count in output, got %q", out)
	}

	passed := true
	buf.Reset()
	renderEvent(&buf, runner.Event{Type: "result", Passed: &passed})
	if !strings.Contains(buf.String(), "通过") || strings.Contains(buf.String(), "未通过") {
		t.Errorf("expected 通过 for passed result, got %q", buf.String())
	}

	buf.Reset()
	renderEvent(&buf, runner.Event{Type: "setup", Message: "准备挑战"})
	if !strings.Contains(buf.String(), "准备挑战") {
		t.Errorf("expected plain message passthrough, got %q", buf.String())
	}
}

func TestVerifyFailedExitCodeIsDistinct(t *testing.T) {
	if exitVerifyFailed == 0 || exitVerifyFailed == 1 {
		t.Fatalf("exitVerifyFailed = %d, must differ from 0 (pass) and 1 (execution error)", exitVerifyFailed)
	}
}

func TestMissingChallengeIDReturnsNonZero(t *testing.T) {
	root := t.TempDir()
	var stderr bytes.Buffer
	code := Run(context.Background(), Env{
		Args:          []string{"challenge", "run", "missing", "--json"},
		Stdout:        ioDiscard{},
		Stderr:        &stderr,
		ChallengesDir: filepath.Join(root, "challenges"),
	})
	if code == 0 {
		t.Fatal("expected non-zero exit")
	}
}

func TestDataDumpJSON(t *testing.T) {
	root := t.TempDir()
	writeChallenge(t, root, "linux-basics", "ls-basic")
	refsPath := filepath.Join(root, "commands.yaml")
	progressPath := filepath.Join(root, "progress.json")
	if err := os.WriteFile(refsPath, []byte("commands: []\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	var stdout bytes.Buffer
	code := Run(context.Background(), Env{
		Args:          []string{"data", "dump", "--json"},
		Stdout:        &stdout,
		Stderr:        ioDiscard{},
		ChallengesDir: filepath.Join(root, "challenges"),
		RefsPath:      refsPath,
		ProgressPath:  progressPath,
	})
	if code != 0 {
		t.Fatalf("data dump returned %d", code)
	}
	if !strings.Contains(stdout.String(), `"categories"`) || !strings.Contains(stdout.String(), `"ls-basic"`) {
		t.Fatalf("unexpected data dump: %s", stdout.String())
	}
}

func writeChallenge(t *testing.T, root, category, id string) {
	t.Helper()
	dir := filepath.Join(root, "challenges", category, id)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	body := []byte(`id: ` + id + `
title: "ls 基础"
difficulty: 1
category: ` + category + `
subcategory: files
tags: [ls]
description: "列出文件"
verify: []
`)
	if err := os.WriteFile(filepath.Join(dir, "challenge.yaml"), body, 0o644); err != nil {
		t.Fatal(err)
	}
}

type ioDiscard struct{}

func (ioDiscard) Write(p []byte) (int, error) {
	return len(p), nil
}
