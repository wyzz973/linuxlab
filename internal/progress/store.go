package progress

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// ProgressData holds all progress tracking information.
type ProgressData struct {
	Skills     map[string]*SkillEntry     `json:"skills"`
	Challenges map[string]*ChallengeEntry `json:"challenges"`
}

// SkillEntry tracks progress for a category.subcategory skill.
type SkillEntry struct {
	Total  int     `json:"total"`
	Passed int     `json:"passed"`
	Score  float64 `json:"score"`
}

// ChallengeEntry tracks progress for an individual challenge.
type ChallengeEntry struct {
	Status    string `json:"status"`
	Attempts  int    `json:"attempts"`
	HintsUsed int    `json:"hints_used"`
	// LastAttempt is the day of the last attempt, kept for display and for
	// files written before LastAttemptAt existed.
	LastAttempt string `json:"last_attempt"`
	// LastAttemptAt is the precise RFC3339 timestamp of the last attempt.
	// Day granularity cannot order several attempts made in one session, which
	// is what "where did I leave off" needs. Empty in older progress files.
	LastAttemptAt string `json:"last_attempt_at,omitempty"`
}

// Store manages reading and writing progress data.
type Store struct {
	Data ProgressData
	path string
}

// NewStore loads progress from the given path, or creates an empty store if the file does not exist.
// A corrupt progress file is renamed to <path>.corrupt and the store starts empty,
// so the application can still launch.
func NewStore(path string) (*Store, error) {
	s := &Store{
		Data: ProgressData{
			Skills:     make(map[string]*SkillEntry),
			Challenges: make(map[string]*ChallengeEntry),
		},
		path: path,
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return s, nil
		}
		return nil, err
	}

	if err := json.Unmarshal(data, &s.Data); err != nil {
		backupPath := path + ".corrupt"
		if renameErr := os.Rename(path, backupPath); renameErr != nil {
			return nil, fmt.Errorf("进度文件损坏 (%v) 且备份失败: %w", err, renameErr)
		}
		fmt.Fprintf(os.Stderr, "警告: 进度文件损坏 (%v)，已备份到 %s，以空进度启动\n", err, backupPath)
		s.Data = ProgressData{
			Skills:     make(map[string]*SkillEntry),
			Challenges: make(map[string]*ChallengeEntry),
		}
		return s, nil
	}

	if s.Data.Skills == nil {
		s.Data.Skills = make(map[string]*SkillEntry)
	}
	if s.Data.Challenges == nil {
		s.Data.Challenges = make(map[string]*ChallengeEntry)
	}

	// JSON null values deserialize to entries with nil pointer values;
	// drop them so later reads (BuildSkillMap/RecordAttempt) never dereference nil.
	for key, skill := range s.Data.Skills {
		if skill == nil {
			delete(s.Data.Skills, key)
		}
	}
	for key, entry := range s.Data.Challenges {
		if entry == nil {
			delete(s.Data.Challenges, key)
		}
	}

	return s, nil
}

// RecordAttempt records a challenge attempt, updating both challenge and skill entries.
func (s *Store) RecordAttempt(challengeID, category, subcategory string, passed bool, hintsUsed int) {
	entry, exists := s.Data.Challenges[challengeID]
	if entry == nil {
		exists = false
	}
	wasAlreadyPassed := exists && entry.Status == "passed"
	isFirstAttempt := !exists

	if !exists {
		entry = &ChallengeEntry{}
		s.Data.Challenges[challengeID] = entry
	}

	entry.Attempts++
	now := time.Now()
	entry.LastAttempt = now.Format("2006-01-02")
	entry.LastAttemptAt = now.Format(time.RFC3339Nano)

	// Status: passed overrides failed, never downgrade.
	if passed {
		entry.Status = "passed"
	} else if entry.Status != "passed" {
		entry.Status = "failed"
	}

	// HintsUsed = max of current and new.
	if hintsUsed > entry.HintsUsed {
		entry.HintsUsed = hintsUsed
	}

	// Update skill entry.
	skillKey := category + "." + subcategory
	skill, ok := s.Data.Skills[skillKey]
	if !ok || skill == nil {
		skill = &SkillEntry{}
		s.Data.Skills[skillKey] = skill
	}

	if isFirstAttempt {
		skill.Total++
	}

	// Increment passed only on first-time pass for this challenge.
	if passed && !wasAlreadyPassed {
		skill.Passed++
	}

	if skill.Total > 0 {
		skill.Score = float64(skill.Passed) / float64(skill.Total)
	}
}

// Save writes the progress data to disk as JSON.
// The write is atomic: data goes to a temp file in the same directory,
// is fsynced, then renamed over the target, so a crash mid-write can
// never leave a truncated progress.json behind.
// Path returns the file the progress data is persisted to.
func (s *Store) Path() string { return s.path }

func (s *Store) Save() error {
	data, err := json.MarshalIndent(s.Data, "", "  ")
	if err != nil {
		return err
	}

	dir := filepath.Dir(s.path)
	tmp, err := os.CreateTemp(dir, ".progress-*.tmp")
	if err != nil {
		return err
	}
	tmpPath := tmp.Name()
	cleanup := func() { _ = os.Remove(tmpPath) }

	if _, err := tmp.Write(data); err != nil {
		_ = tmp.Close()
		cleanup()
		return err
	}
	if err := tmp.Sync(); err != nil {
		_ = tmp.Close()
		cleanup()
		return err
	}
	if err := tmp.Close(); err != nil {
		cleanup()
		return err
	}
	if err := os.Chmod(tmpPath, 0o644); err != nil {
		cleanup()
		return err
	}
	if err := os.Rename(tmpPath, s.path); err != nil {
		cleanup()
		return err
	}
	return nil
}
