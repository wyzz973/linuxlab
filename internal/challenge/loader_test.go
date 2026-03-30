package challenge

import (
	"os"
	"path/filepath"
	"testing"
)

func writeChallenge(t *testing.T, dir string, yaml string) {
	t.Helper()
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "challenge.yaml"), []byte(yaml), 0o644); err != nil {
		t.Fatal(err)
	}
}

const sampleYAML = `id: test-challenge
title: "测试挑战"
difficulty: 1
category: linux-basics
subcategory: files
tags: [ls, cat]
description: "一个测试挑战"
verify:
  - type: file_content
    path: /tmp/out.txt
    expect: "hello"
`

func TestLoadChallenge(t *testing.T) {
	dir := t.TempDir()
	challengeDir := filepath.Join(dir, "linux-basics", "test-challenge")
	writeChallenge(t, challengeDir, sampleYAML)

	c, err := LoadChallenge(challengeDir)
	if err != nil {
		t.Fatalf("LoadChallenge failed: %v", err)
	}
	if c.ID != "test-challenge" {
		t.Errorf("ID = %q, want %q", c.ID, "test-challenge")
	}
	if c.Title != "测试挑战" {
		t.Errorf("Title = %q", c.Title)
	}
	if c.Dir != challengeDir {
		t.Errorf("Dir = %q, want %q", c.Dir, challengeDir)
	}
	if c.Difficulty != 1 {
		t.Errorf("Difficulty = %d, want 1", c.Difficulty)
	}
}

func TestLoadAllChallenges(t *testing.T) {
	root := t.TempDir()

	writeChallenge(t, filepath.Join(root, "linux-basics", "challenge-a"), `id: challenge-a
title: "挑战A"
difficulty: 1
category: linux-basics
subcategory: files
tags: [ls]
description: "A"
verify:
  - type: file_content
    path: /tmp/a.txt
    expect: "a"
`)
	writeChallenge(t, filepath.Join(root, "linux-basics", "challenge-b"), `id: challenge-b
title: "挑战B"
difficulty: 2
category: linux-basics
subcategory: files
tags: [cat]
description: "B"
verify:
  - type: file_content
    path: /tmp/b.txt
    expect: "b"
`)
	writeChallenge(t, filepath.Join(root, "vim", "vim-a"), `id: vim-a
title: "Vim挑战A"
difficulty: 1
category: vim
subcategory: navigation
tags: [hjkl]
description: "VA"
verify:
  - type: file_content
    path: /tmp/va.txt
    expect: "va"
`)
	writeChallenge(t, filepath.Join(root, "vim", "vim-b"), `id: vim-b
title: "Vim挑战B"
difficulty: 2
category: vim
subcategory: editing
tags: [insert]
description: "VB"
verify:
  - type: file_content
    path: /tmp/vb.txt
    expect: "vb"
`)

	challenges, err := LoadAll(root)
	if err != nil {
		t.Fatalf("LoadAll failed: %v", err)
	}
	if len(challenges) != 4 {
		t.Fatalf("LoadAll returned %d challenges, want 4", len(challenges))
	}

	// Should be sorted by category then ID
	if challenges[0].ID != "challenge-a" {
		t.Errorf("challenges[0].ID = %q, want %q", challenges[0].ID, "challenge-a")
	}
	if challenges[1].ID != "challenge-b" {
		t.Errorf("challenges[1].ID = %q, want %q", challenges[1].ID, "challenge-b")
	}
	if challenges[2].ID != "vim-a" {
		t.Errorf("challenges[2].ID = %q, want %q", challenges[2].ID, "vim-a")
	}
}

func TestLoadAllByCategory(t *testing.T) {
	root := t.TempDir()

	writeChallenge(t, filepath.Join(root, "linux-basics", "ch1"), `id: ch1
title: "挑战1"
difficulty: 1
category: linux-basics
subcategory: files
tags: [ls]
description: "1"
verify:
  - type: file_content
    path: /tmp/1.txt
    expect: "1"
`)
	writeChallenge(t, filepath.Join(root, "vim", "ch2"), `id: ch2
title: "挑战2"
difficulty: 1
category: vim
subcategory: nav
tags: [hjkl]
description: "2"
verify:
  - type: file_content
    path: /tmp/2.txt
    expect: "2"
`)

	byCategory, err := LoadAllByCategory(root)
	if err != nil {
		t.Fatalf("LoadAllByCategory failed: %v", err)
	}
	if len(byCategory) != 2 {
		t.Fatalf("got %d categories, want 2", len(byCategory))
	}
	if len(byCategory["linux-basics"]) != 1 {
		t.Errorf("linux-basics has %d challenges, want 1", len(byCategory["linux-basics"]))
	}
	if len(byCategory["vim"]) != 1 {
		t.Errorf("vim has %d challenges, want 1", len(byCategory["vim"]))
	}
}

func TestLoadChallengeNotFound(t *testing.T) {
	_, err := LoadChallenge("/nonexistent/path/that/does/not/exist")
	if err == nil {
		t.Fatal("expected error for nonexistent path, got nil")
	}
}
