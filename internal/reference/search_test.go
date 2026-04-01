package reference

import (
	"testing"
)

var testCommands = []CommandRef{
	{Name: "find", Brief: "搜索文件和目录"},
	{Name: "grep", Brief: "搜索文本模式"},
	{Name: "ls", Brief: "列出目录内容"},
	{Name: "cat", Brief: "查看文件内容"},
	{Name: "findmnt", Brief: "查找挂载点"},
}

func TestSearch_ExactName(t *testing.T) {
	results := Search("find", testCommands)
	if len(results) == 0 {
		t.Fatal("expected results")
	}
	if results[0].Name != "find" {
		t.Errorf("expected 'find' first, got %q", results[0].Name)
	}
}

func TestSearch_Prefix(t *testing.T) {
	results := Search("gr", testCommands)
	if len(results) == 0 {
		t.Fatal("expected results")
	}
	if results[0].Name != "grep" {
		t.Errorf("expected 'grep' first, got %q", results[0].Name)
	}
}

func TestSearch_BriefMatch(t *testing.T) {
	results := Search("搜索", testCommands)
	if len(results) < 2 {
		t.Fatalf("expected at least 2 results, got %d", len(results))
	}
	// Both find and grep have "搜索" in brief
	names := map[string]bool{}
	for _, r := range results {
		names[r.Name] = true
	}
	if !names["find"] {
		t.Error("expected find in results")
	}
	if !names["grep"] {
		t.Error("expected grep in results")
	}
}

func TestSearch_CaseInsensitive(t *testing.T) {
	results := Search("FIND", testCommands)
	if len(results) == 0 {
		t.Fatal("expected results")
	}
	found := false
	for _, r := range results {
		if r.Name == "find" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected 'find' in results for query 'FIND'")
	}
}

func TestSearch_EmptyQuery(t *testing.T) {
	results := Search("", testCommands)
	if len(results) != len(testCommands) {
		t.Errorf("expected %d results, got %d", len(testCommands), len(results))
	}
}

func TestSearch_NoMatch(t *testing.T) {
	results := Search("zzz", testCommands)
	if len(results) != 0 {
		t.Errorf("expected 0 results, got %d", len(results))
	}
}
