package progress

import (
	"encoding/json"
	"os"
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
	Status      string `json:"status"`
	Attempts    int    `json:"attempts"`
	HintsUsed   int    `json:"hints_used"`
	LastAttempt string `json:"last_attempt"`
}

// Store manages reading and writing progress data.
type Store struct {
	Data ProgressData
	path string
}

// NewStore loads progress from the given path, or creates an empty store if the file does not exist.
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
		return nil, err
	}

	if s.Data.Skills == nil {
		s.Data.Skills = make(map[string]*SkillEntry)
	}
	if s.Data.Challenges == nil {
		s.Data.Challenges = make(map[string]*ChallengeEntry)
	}

	return s, nil
}

// RecordAttempt records a challenge attempt, updating both challenge and skill entries.
func (s *Store) RecordAttempt(challengeID, category, subcategory string, passed bool, hintsUsed int) {
	entry, exists := s.Data.Challenges[challengeID]
	wasAlreadyPassed := exists && entry.Status == "passed"
	isFirstAttempt := !exists

	if !exists {
		entry = &ChallengeEntry{}
		s.Data.Challenges[challengeID] = entry
	}

	entry.Attempts++
	entry.LastAttempt = time.Now().Format("2006-01-02")

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
	if !ok {
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
func (s *Store) Save() error {
	data, err := json.MarshalIndent(s.Data, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, data, 0644)
}
