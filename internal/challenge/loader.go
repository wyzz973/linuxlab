package challenge

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"gopkg.in/yaml.v3"
)

// LoadChallenge reads a challenge.yaml file from the given directory.
func LoadChallenge(dir string) (*Challenge, error) {
	data, err := os.ReadFile(filepath.Join(dir, "challenge.yaml"))
	if err != nil {
		return nil, fmt.Errorf("read challenge.yaml in %s: %w", dir, err)
	}

	var c Challenge
	if err := yaml.Unmarshal(data, &c); err != nil {
		return nil, fmt.Errorf("parse challenge.yaml in %s: %w", dir, err)
	}
	c.Dir = dir
	return &c, nil
}

// LoadAll walks root/{category}/{challenge}/ directories and returns all
// challenges sorted by category then ID.
func LoadAll(root string) ([]*Challenge, error) {
	entries, err := os.ReadDir(root)
	if err != nil {
		return nil, fmt.Errorf("read root dir %s: %w", root, err)
	}

	var challenges []*Challenge
	for _, catEntry := range entries {
		if !catEntry.IsDir() {
			continue
		}
		catDir := filepath.Join(root, catEntry.Name())
		chEntries, err := os.ReadDir(catDir)
		if err != nil {
			return nil, fmt.Errorf("read category dir %s: %w", catDir, err)
		}
		for _, chEntry := range chEntries {
			if !chEntry.IsDir() {
				continue
			}
			chDir := filepath.Join(catDir, chEntry.Name())
			// Skip directories without challenge.yaml
			if _, err := os.Stat(filepath.Join(chDir, "challenge.yaml")); os.IsNotExist(err) {
				continue
			}
			c, err := LoadChallenge(chDir)
			if err != nil {
				return nil, err
			}
			challenges = append(challenges, c)
		}
	}

	sort.Slice(challenges, func(i, j int) bool {
		if challenges[i].Category != challenges[j].Category {
			return challenges[i].Category < challenges[j].Category
		}
		return challenges[i].ID < challenges[j].ID
	})

	return challenges, nil
}

// LoadAllByCategory loads all challenges and groups them by category.
func LoadAllByCategory(root string) (map[string][]*Challenge, error) {
	challenges, err := LoadAll(root)
	if err != nil {
		return nil, err
	}

	result := make(map[string][]*Challenge)
	for _, c := range challenges {
		result[c.Category] = append(result[c.Category], c)
	}
	return result, nil
}
