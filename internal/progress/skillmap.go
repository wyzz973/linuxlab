package progress

import (
	"math"
	"sort"
	"strings"

	"github.com/sd3/linuxlab/internal/challenge"
)

// SkillMap provides an overview of progress across all categories.
type SkillMap struct {
	Categories   []CategoryScore
	TotalPassed  int
	TotalCount   int
	OverallScore float64
}

// CategoryScore holds aggregate scores for a category.
type CategoryScore struct {
	Name          string
	Score         float64
	Subcategories []SubcategoryScore
}

// SubcategoryScore holds progress for a single subcategory.
type SubcategoryScore struct {
	Name   string
	Total  int
	Passed int
	Score  float64
}

// BuildSkillMap parses skill keys "category.subcategory" from the store,
// groups them into categories, and calculates averages.
func BuildSkillMap(store *Store) *SkillMap {
	catMap := make(map[string][]SubcategoryScore)

	totalPassed := 0
	totalCount := 0

	for key, skill := range store.Data.Skills {
		parts := strings.SplitN(key, ".", 2)
		if len(parts) != 2 {
			continue
		}
		cat, sub := parts[0], parts[1]

		sc := SubcategoryScore{
			Name:   sub,
			Total:  skill.Total,
			Passed: skill.Passed,
			Score:  skill.Score,
		}
		catMap[cat] = append(catMap[cat], sc)

		totalPassed += skill.Passed
		totalCount += skill.Total
	}

	var categories []CategoryScore
	for name, subs := range catMap {
		sort.Slice(subs, func(i, j int) bool {
			return subs[i].Name < subs[j].Name
		})

		var sum float64
		for _, s := range subs {
			sum += s.Score
		}
		avg := sum / float64(len(subs))

		categories = append(categories, CategoryScore{
			Name:          name,
			Score:         avg,
			Subcategories: subs,
		})
	}

	sort.Slice(categories, func(i, j int) bool {
		return categories[i].Name < categories[j].Name
	})

	var overall float64
	if totalCount > 0 {
		overall = float64(totalPassed) / float64(totalCount)
	}

	return &SkillMap{
		Categories:   categories,
		TotalPassed:  totalPassed,
		TotalCount:   totalCount,
		OverallScore: overall,
	}
}

// RecommendWeakest finds the lowest-scoring subcategory and returns the first
// unpassed challenge from it.
func RecommendWeakest(store *Store, challenges []*challenge.Challenge) *challenge.Challenge {
	sm := BuildSkillMap(store)

	if len(sm.Categories) == 0 {
		return nil
	}

	// Find the subcategory with the lowest score.
	var weakestCat, weakestSub string
	lowestScore := math.MaxFloat64

	for _, cat := range sm.Categories {
		for _, sub := range cat.Subcategories {
			if sub.Score < lowestScore {
				lowestScore = sub.Score
				weakestCat = cat.Name
				weakestSub = sub.Name
			}
		}
	}

	// Find first unpassed challenge in that subcategory.
	for _, ch := range challenges {
		if ch.Category != weakestCat || ch.Subcategory != weakestSub {
			continue
		}
		entry, exists := store.Data.Challenges[ch.ID]
		if !exists || entry.Status != "passed" {
			return ch
		}
	}

	return nil
}

// RecommendMultiple returns up to limit unpassed challenges from the weakest subcategories.
func RecommendMultiple(store *Store, challenges []*challenge.Challenge, limit int) []*challenge.Challenge {
	sm := BuildSkillMap(store)
	if len(sm.Categories) == 0 || len(challenges) == 0 {
		return nil
	}

	// Collect all subcategories sorted by score (ascending = weakest first)
	type subKey struct{ cat, sub string }
	type scored struct {
		key   subKey
		score float64
	}
	var subs []scored
	for _, cat := range sm.Categories {
		for _, sub := range cat.Subcategories {
			subs = append(subs, scored{key: subKey{cat.Name, sub.Name}, score: sub.Score})
		}
	}
	sort.Slice(subs, func(i, j int) bool { return subs[i].score < subs[j].score })

	var result []*challenge.Challenge
	seen := make(map[string]bool)

	for _, s := range subs {
		if len(result) >= limit {
			break
		}
		for _, ch := range challenges {
			if len(result) >= limit {
				break
			}
			if ch.Category != s.key.cat || ch.Subcategory != s.key.sub {
				continue
			}
			if seen[ch.ID] {
				continue
			}
			entry := store.Data.Challenges[ch.ID]
			if entry != nil && entry.Status == "passed" {
				continue
			}
			result = append(result, ch)
			seen[ch.ID] = true
		}
	}
	return result
}

// ScoreWithHints calculates a score penalty for hint usage.
// Each hint multiplies the score by 0.8.
func ScoreWithHints(hintsUsed int) float64 {
	return math.Pow(0.8, float64(hintsUsed))
}
