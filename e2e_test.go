package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/sd3/linuxlab/internal/challenge"
	"github.com/sd3/linuxlab/internal/progress"
	"github.com/sd3/linuxlab/internal/tui"
)

func TestE2E_LoadChallengesAndCreateApp(t *testing.T) {
	byCategory, err := challenge.LoadAllByCategory("challenges")
	if err != nil {
		t.Fatalf("LoadAllByCategory failed: %v", err)
	}

	if len(byCategory) < 4 {
		t.Errorf("expected at least 4 categories, got %d", len(byCategory))
	}

	for _, tc := range []struct {
		cat string
		min int
	}{
		{"linux-basics", 50},
		{"vim", 30},
		{"shell-scripting", 40},
		{"ops", 40},
	} {
		if len(byCategory[tc.cat]) < tc.min {
			t.Errorf("%s has %d challenges, want >= %d", tc.cat, len(byCategory[tc.cat]), tc.min)
		}
	}

	// Verify each challenge has required fields
	for cat, challenges := range byCategory {
		for _, ch := range challenges {
			if ch.ID == "" {
				t.Errorf("[%s] challenge has empty ID", cat)
			}
			if ch.Title == "" {
				t.Errorf("[%s/%s] has empty Title", cat, ch.ID)
			}
			if ch.Difficulty < 1 || ch.Difficulty > 5 {
				t.Errorf("[%s/%s] difficulty = %d, want 1-5", cat, ch.ID, ch.Difficulty)
			}
			if len(ch.Verify) == 0 {
				t.Errorf("[%s/%s] has no verify rules", cat, ch.ID)
			}
			if ch.Description == "" {
				t.Errorf("[%s/%s] has empty description", cat, ch.ID)
			}
		}
	}
}

func TestE2E_ProgressRoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "progress.json")

	store, _ := progress.NewStore(path)
	store.RecordAttempt("list-hidden-files", "linux-basics", "file-operations", true, 1)
	store.RecordAttempt("find-large-files", "linux-basics", "file-operations", false, 0)
	store.Save()

	store2, _ := progress.NewStore(path)
	sm := progress.BuildSkillMap(store2)

	if sm.TotalPassed != 1 {
		t.Errorf("TotalPassed = %d, want 1", sm.TotalPassed)
	}
	if sm.TotalCount != 2 {
		t.Errorf("TotalCount = %d, want 2", sm.TotalCount)
	}
}

func TestE2E_AppModelCreation(t *testing.T) {
	byCategory, err := challenge.LoadAllByCategory("challenges")
	if err != nil {
		t.Skip("challenges directory not available")
	}

	dir := t.TempDir()
	store, _ := progress.NewStore(filepath.Join(dir, "progress.json"))

	app := tui.NewAppModel(byCategory, store, nil)
	view := app.View()
	if view == "" {
		t.Error("app view is empty")
	}
}

func TestE2E_AllSolutionScriptsExist(t *testing.T) {
	err := filepath.Walk("challenges", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.Name() == "challenge.yaml" {
			dir := filepath.Dir(path)
			checkPath := filepath.Join(dir, "check.sh")
			if _, err := os.Stat(checkPath); os.IsNotExist(err) {
				t.Errorf("missing check.sh in %s", dir)
			}
			solutionPath := filepath.Join(dir, "solution.sh")
			if _, err := os.Stat(solutionPath); os.IsNotExist(err) {
				t.Errorf("missing solution.sh in %s", dir)
			}
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}
