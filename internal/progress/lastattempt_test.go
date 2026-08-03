package progress

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/sd3/linuxlab/internal/challenge"
)

func lastStore(t *testing.T) *Store {
	t.Helper()
	s, err := NewStore(filepath.Join(t.TempDir(), "progress.json"))
	if err != nil {
		t.Fatal(err)
	}
	return s
}

func lastCats() []*challenge.Challenge {
	return []*challenge.Challenge{
		{ID: "a", Title: "题 A", Category: "c", Subcategory: "s"},
		{ID: "b", Title: "题 B", Category: "c", Subcategory: "s"},
		{ID: "c", Title: "题 C", Category: "c", Subcategory: "s"},
	}
}

// Several attempts in one session share a date, so ordering has to come from
// the precise timestamp.
func TestLastAttempted_UsesPreciseTimestampWithinOneDay(t *testing.T) {
	store := lastStore(t)
	chs := lastCats()
	for _, id := range []string{"a", "b", "c"} {
		store.RecordAttempt(id, "c", "s", false, 0)
		time.Sleep(2 * time.Millisecond)
	}

	got, entry := LastAttempted(store, chs)
	if got == nil || got.ID != "c" {
		t.Fatalf("LastAttempted = %v, want the challenge attempted last (c)", got)
	}
	if entry == nil || entry.Attempts != 1 {
		t.Fatalf("entry = %+v, want the matching progress entry", entry)
	}
}

// Progress files written before LastAttemptAt existed only carry a day, and
// must still resolve deterministically.
func TestLastAttempted_FallsBackToDayForLegacyEntries(t *testing.T) {
	store := lastStore(t)
	store.Data.Challenges["a"] = &ChallengeEntry{Status: "failed", Attempts: 1, LastAttempt: "2026-07-20"}
	store.Data.Challenges["b"] = &ChallengeEntry{Status: "passed", Attempts: 3, LastAttempt: "2026-07-28"}
	store.Data.Challenges["c"] = &ChallengeEntry{Status: "failed", Attempts: 2, LastAttempt: "2026-07-28"}

	got, _ := LastAttempted(store, lastCats())
	if got == nil || got.ID != "b" {
		t.Fatalf("LastAttempted = %v, want b (latest day, ID tie-break)", got)
	}

	// Stability: repeated calls agree.
	again, _ := LastAttempted(store, lastCats())
	if again.ID != got.ID {
		t.Fatalf("LastAttempted is not deterministic: %s vs %s", again.ID, got.ID)
	}
}

// A timestamped entry is more recent information than a day-only one.
func TestLastAttempted_TimestampBeatsLegacyDay(t *testing.T) {
	store := lastStore(t)
	store.Data.Challenges["a"] = &ChallengeEntry{Status: "failed", Attempts: 1, LastAttempt: "2026-07-28"}
	store.RecordAttempt("b", "c", "s", false, 0)

	got, _ := LastAttempted(store, lastCats())
	if got == nil || got.ID != "b" {
		t.Fatalf("LastAttempted = %v, want the timestamped entry b", got)
	}
}

func TestLastAttempted_EmptyAndNil(t *testing.T) {
	if ch, entry := LastAttempted(nil, lastCats()); ch != nil || entry != nil {
		t.Fatal("nil store should yield no result")
	}
	if ch, _ := LastAttempted(lastStore(t), lastCats()); ch != nil {
		t.Fatal("a store with no attempts should yield no result")
	}
}
