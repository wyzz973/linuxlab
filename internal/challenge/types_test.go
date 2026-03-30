package challenge

import (
	"testing"

	"gopkg.in/yaml.v3"
)

func TestChallengeUnmarshalYAML(t *testing.T) {
	data := []byte(`
id: find-large-files
title: "查找大文件"
difficulty: 2
category: linux-basics
subcategory: file-operations
tags: [find, du, sort]
description: |
  找出最大的文件。
hints:
  - level: 1
    text: "试试 du 命令"
  - level: 2
    text: "du -a 配合 sort"
verify:
  - type: file_content
    path: /tmp/result.txt
    expect: "/var/log/auth.log"
`)
	var c Challenge
	if err := yaml.Unmarshal(data, &c); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}
	if c.ID != "find-large-files" {
		t.Errorf("ID = %q, want %q", c.ID, "find-large-files")
	}
	if c.Title != "查找大文件" {
		t.Errorf("Title = %q", c.Title)
	}
	if c.Difficulty != 2 {
		t.Errorf("Difficulty = %d", c.Difficulty)
	}
	if c.Category != "linux-basics" {
		t.Errorf("Category = %q", c.Category)
	}
	if c.Subcategory != "file-operations" {
		t.Errorf("Subcategory = %q", c.Subcategory)
	}
	if len(c.Tags) != 3 {
		t.Fatalf("Tags len = %d", len(c.Tags))
	}
	if len(c.Hints) != 2 {
		t.Fatalf("Hints len = %d", len(c.Hints))
	}
	if c.Hints[0].Level != 1 {
		t.Errorf("Hints[0].Level = %d", c.Hints[0].Level)
	}
	if len(c.Verify) != 1 {
		t.Fatalf("Verify len = %d", len(c.Verify))
	}
	if c.Verify[0].Type != "file_content" {
		t.Errorf("Verify[0].Type = %q", c.Verify[0].Type)
	}
}

func TestVimChallengeUnmarshalYAML(t *testing.T) {
	data := []byte(`
id: vim-batch-replace
title: "批量替换"
difficulty: 2
category: vim
subcategory: editing
tags: [substitute]
description: "将 foo 替换为 bar"
setup_files:
  - path: /tmp/vim-challenge.txt
    content: "foo is great"
verify:
  - type: file_content
    path: /tmp/vim-challenge.txt
    expect: "bar is great"
`)
	var c Challenge
	if err := yaml.Unmarshal(data, &c); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}
	if c.Category != "vim" {
		t.Errorf("Category = %q", c.Category)
	}
	if len(c.SetupFiles) != 1 {
		t.Fatalf("SetupFiles len = %d", len(c.SetupFiles))
	}
	if c.SetupFiles[0].Path != "/tmp/vim-challenge.txt" {
		t.Errorf("Path = %q", c.SetupFiles[0].Path)
	}
}
