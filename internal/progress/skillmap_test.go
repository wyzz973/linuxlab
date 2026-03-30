package progress

import (
	"math"
	"testing"

	"github.com/sd3/linuxlab/internal/challenge"
)

func TestBuildSkillMap(t *testing.T) {
	dir := t.TempDir()
	store, err := NewStore(dir + "/progress.json")
	if err != nil {
		t.Fatal(err)
	}

	// 3 skills across 2 categories
	store.RecordAttempt("ch1", "linux", "files", true, 0)
	store.RecordAttempt("ch2", "linux", "permissions", true, 0)
	store.RecordAttempt("ch3", "linux", "permissions", false, 0)
	store.RecordAttempt("ch4", "vim", "navigation", true, 0)

	sm := BuildSkillMap(store)

	// Should have 2 categories: linux, vim
	if len(sm.Categories) != 2 {
		t.Fatalf("categories count = %d, want 2", len(sm.Categories))
	}

	// Find linux category
	var linux, vim *CategoryScore
	for i := range sm.Categories {
		switch sm.Categories[i].Name {
		case "linux":
			linux = &sm.Categories[i]
		case "vim":
			vim = &sm.Categories[i]
		}
	}

	if linux == nil {
		t.Fatal("linux category not found")
	}
	if vim == nil {
		t.Fatal("vim category not found")
	}

	// Linux has 2 subcategories: files (1/1=1.0) and permissions (1/2=0.5)
	if len(linux.Subcategories) != 2 {
		t.Fatalf("linux subcategories = %d, want 2", len(linux.Subcategories))
	}

	// Category score = avg of subcategory scores = (1.0 + 0.5) / 2 = 0.75
	if math.Abs(linux.Score-0.75) > 0.001 {
		t.Errorf("linux score = %f, want 0.75", linux.Score)
	}

	// Vim has 1 subcategory: navigation (1/1=1.0)
	if len(vim.Subcategories) != 1 {
		t.Fatalf("vim subcategories = %d, want 1", len(vim.Subcategories))
	}
	if math.Abs(vim.Score-1.0) > 0.001 {
		t.Errorf("vim score = %f, want 1.0", vim.Score)
	}

	// Overall counts
	if sm.TotalPassed != 3 {
		t.Errorf("totalPassed = %d, want 3", sm.TotalPassed)
	}
	if sm.TotalCount != 4 {
		t.Errorf("totalCount = %d, want 4", sm.TotalCount)
	}
}

func TestRecommendWeakest(t *testing.T) {
	dir := t.TempDir()
	store, err := NewStore(dir + "/progress.json")
	if err != nil {
		t.Fatal(err)
	}

	// linux.files: 1 passed (score 1.0)
	store.RecordAttempt("ch1", "linux", "files", true, 0)
	// linux.permissions: 0 passed, 1 failed (score 0.0) - weakest
	store.RecordAttempt("ch2", "linux", "permissions", false, 0)

	challenges := []*challenge.Challenge{
		{ID: "ch1", Category: "linux", Subcategory: "files"},
		{ID: "ch2", Category: "linux", Subcategory: "permissions"},
		{ID: "ch3", Category: "linux", Subcategory: "permissions"},
	}

	rec := RecommendWeakest(store, challenges)
	if rec == nil {
		t.Fatal("expected a recommendation, got nil")
	}
	// Should recommend from weakest subcategory "permissions".
	// ch2 is failed (unpassed) and comes first in the list.
	if rec.Subcategory != "permissions" {
		t.Errorf("recommended subcategory = %q, want %q", rec.Subcategory, "permissions")
	}
	if rec.ID != "ch2" && rec.ID != "ch3" {
		t.Errorf("recommended = %q, want ch2 or ch3 (unpassed in permissions)", rec.ID)
	}
}

func TestScoreWithHints(t *testing.T) {
	tests := []struct {
		hints int
		want  float64
	}{
		{0, 1.0},
		{1, 0.8},
		{2, 0.64},
	}

	for _, tt := range tests {
		got := ScoreWithHints(tt.hints)
		if math.Abs(got-tt.want) > 0.001 {
			t.Errorf("ScoreWithHints(%d) = %f, want %f", tt.hints, got, tt.want)
		}
	}
}
