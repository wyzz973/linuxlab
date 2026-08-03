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

func TestNewStore_CorruptFileBackedUpAndStartsEmpty(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "progress.json")
	corrupt := `{"skills": {"linux.files": {"total":`
	if err := os.WriteFile(path, []byte(corrupt), 0o644); err != nil {
		t.Fatal(err)
	}

	store, err := NewStore(path)
	if err != nil {
		t.Fatalf("NewStore() on corrupt file should recover, got error: %v", err)
	}
	if len(store.Data.Skills) != 0 || len(store.Data.Challenges) != 0 {
		t.Error("store should start empty after corrupt file recovery")
	}

	backup, err := os.ReadFile(path + ".corrupt")
	if err != nil {
		t.Fatalf("corrupt file was not backed up: %v", err)
	}
	if string(backup) != corrupt {
		t.Errorf("backup content = %q, want original corrupt bytes", string(backup))
	}
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Error("corrupt progress.json should have been renamed away")
	}

	// Recovered store must be fully usable.
	store.RecordAttempt("ch1", "linux", "files", true, 0)
	if err := store.Save(); err != nil {
		t.Fatalf("Save() after recovery error = %v", err)
	}
	store2, err := NewStore(path)
	if err != nil {
		t.Fatalf("NewStore() reload after recovery error = %v", err)
	}
	if store2.Data.Challenges["ch1"] == nil {
		t.Error("progress recorded after recovery was not persisted")
	}
}

func TestNewStore_PrunesNullEntries(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "progress.json")
	body := `{"skills":{"a.b":null,"linux.files":{"total":1,"passed":1,"score":1}},"challenges":{"c1":null}}`
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}

	store, err := NewStore(path)
	if err != nil {
		t.Fatalf("NewStore() error = %v", err)
	}
	if _, exists := store.Data.Skills["a.b"]; exists {
		t.Error("nil skill entry should be pruned")
	}
	if _, exists := store.Data.Challenges["c1"]; exists {
		t.Error("nil challenge entry should be pruned")
	}
	if store.Data.Skills["linux.files"] == nil {
		t.Error("valid skill entry should be kept")
	}

	// These used to panic on nil map values (BuildSkillMap / RecordAttempt).
	sm := BuildSkillMap(store)
	if sm.TotalCount != 1 {
		t.Errorf("TotalCount = %d, want 1", sm.TotalCount)
	}
	store.RecordAttempt("c1", "a", "b", true, 0)
	entry := store.Data.Challenges["c1"]
	if entry == nil || entry.Status != "passed" {
		t.Errorf("RecordAttempt after prune failed, entry = %+v", entry)
	}
}

func TestStore_SaveAtomicLeavesNoTempFiles(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "progress.json")

	store, err := NewStore(path)
	if err != nil {
		t.Fatalf("NewStore() error = %v", err)
	}
	store.RecordAttempt("ch1", "linux", "files", true, 0)
	if err := store.Save(); err != nil {
		t.Fatalf("Save() error = %v", err)
	}
	// Overwrite an existing file — the rename path must replace it cleanly.
	store.RecordAttempt("ch2", "linux", "files", false, 0)
	if err := store.Save(); err != nil {
		t.Fatalf("Save() overwrite error = %v", err)
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 1 || entries[0].Name() != "progress.json" {
		var names []string
		for _, e := range entries {
			names = append(names, e.Name())
		}
		t.Errorf("expected only progress.json in dir, got %v", names)
	}

	store2, err := NewStore(path)
	if err != nil {
		t.Fatalf("NewStore() reload error = %v", err)
	}
	if len(store2.Data.Challenges) != 2 {
		t.Errorf("challenges after reload = %d, want 2", len(store2.Data.Challenges))
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
