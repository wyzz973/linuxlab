package progress

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewStore_CreatesFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "progress.json")

	store, err := NewStore(path)
	if err != nil {
		t.Fatalf("NewStore() error = %v", err)
	}
	if store.Data.Skills == nil {
		t.Fatal("Skills map should not be nil")
	}
	if store.Data.Challenges == nil {
		t.Fatal("Challenges map should not be nil")
	}
	if len(store.Data.Skills) != 0 {
		t.Errorf("Skills map should be empty, got %d", len(store.Data.Skills))
	}
	if len(store.Data.Challenges) != 0 {
		t.Errorf("Challenges map should be empty, got %d", len(store.Data.Challenges))
	}
}

func TestStore_RecordPass(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "progress.json")

	store, err := NewStore(path)
	if err != nil {
		t.Fatalf("NewStore() error = %v", err)
	}

	store.RecordAttempt("ch1", "linux", "files", true, 0)

	entry, ok := store.Data.Challenges["ch1"]
	if !ok {
		t.Fatal("challenge entry not found")
	}
	if entry.Status != "passed" {
		t.Errorf("status = %q, want %q", entry.Status, "passed")
	}
	if entry.Attempts != 1 {
		t.Errorf("attempts = %d, want 1", entry.Attempts)
	}
	if entry.HintsUsed != 0 {
		t.Errorf("hintsUsed = %d, want 0", entry.HintsUsed)
	}
}

func TestStore_RecordFail(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "progress.json")

	store, err := NewStore(path)
	if err != nil {
		t.Fatalf("NewStore() error = %v", err)
	}

	store.RecordAttempt("ch1", "linux", "files", false, 2)

	entry, ok := store.Data.Challenges["ch1"]
	if !ok {
		t.Fatal("challenge entry not found")
	}
	if entry.Status != "failed" {
		t.Errorf("status = %q, want %q", entry.Status, "failed")
	}
	if entry.HintsUsed != 2 {
		t.Errorf("hintsUsed = %d, want 2", entry.HintsUsed)
	}
}

func TestStore_MultipleAttempts(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "progress.json")

	store, err := NewStore(path)
	if err != nil {
		t.Fatalf("NewStore() error = %v", err)
	}

	store.RecordAttempt("ch1", "linux", "files", false, 1)
	store.RecordAttempt("ch1", "linux", "files", true, 0)

	entry := store.Data.Challenges["ch1"]
	if entry.Attempts != 2 {
		t.Errorf("attempts = %d, want 2", entry.Attempts)
	}
	if entry.Status != "passed" {
		t.Errorf("status = %q, want %q", entry.Status, "passed")
	}
	if entry.HintsUsed != 1 {
		t.Errorf("hintsUsed = %d, want 1 (max of all attempts)", entry.HintsUsed)
	}
}

func TestStore_SaveAndLoad(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "progress.json")

	store, err := NewStore(path)
	if err != nil {
		t.Fatalf("NewStore() error = %v", err)
	}

	store.RecordAttempt("ch1", "linux", "files", true, 1)
	if err := store.Save(); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Fatal("progress file was not created")
	}

	// Load from same path
	store2, err := NewStore(path)
	if err != nil {
		t.Fatalf("NewStore() reload error = %v", err)
	}

	entry, ok := store2.Data.Challenges["ch1"]
	if !ok {
		t.Fatal("challenge entry not found after reload")
	}
	if entry.Status != "passed" {
		t.Errorf("status after reload = %q, want %q", entry.Status, "passed")
	}
	if entry.Attempts != 1 {
		t.Errorf("attempts after reload = %d, want 1", entry.Attempts)
	}
}

func TestStore_SkillUpdatedOnRecord(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "progress.json")

	store, err := NewStore(path)
	if err != nil {
		t.Fatalf("NewStore() error = %v", err)
	}

	store.RecordAttempt("ch1", "linux", "files", true, 0)

	skill, ok := store.Data.Skills["linux.files"]
	if !ok {
		t.Fatal("skill entry not found for linux.files")
	}
	if skill.Total != 1 {
		t.Errorf("total = %d, want 1", skill.Total)
	}
	if skill.Passed != 1 {
		t.Errorf("passed = %d, want 1", skill.Passed)
	}
	if skill.Score != 1.0 {
		t.Errorf("score = %f, want 1.0", skill.Score)
	}

	// Add a failed challenge
	store.RecordAttempt("ch2", "linux", "files", false, 0)

	if skill.Total != 2 {
		t.Errorf("total = %d, want 2", skill.Total)
	}
	if skill.Passed != 1 {
		t.Errorf("passed = %d, want 1", skill.Passed)
	}
	if skill.Score != 0.5 {
		t.Errorf("score = %f, want 0.5", skill.Score)
	}
}
