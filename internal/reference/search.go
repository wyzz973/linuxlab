package reference

import (
	"sort"
	"strings"
)

// relevance scores for ranking search results.
const (
	scoreExactName    = 100
	scorePrefixName   = 75
	scoreContainsName = 50
	scoreContainsBrief = 25
)

type scoredRef struct {
	ref   CommandRef
	score int
}

// Search finds commands matching the query using case-insensitive substring
// matching on Name and Brief. Results are sorted by relevance.
// An empty query returns all commands.
func Search(query string, commands []CommandRef) []CommandRef {
	if query == "" {
		result := make([]CommandRef, len(commands))
		copy(result, commands)
		return result
	}

	q := strings.ToLower(query)
	var scored []scoredRef

	for _, cmd := range commands {
		name := strings.ToLower(cmd.Name)
		brief := strings.ToLower(cmd.Brief)

		score := 0
		if name == q {
			score = scoreExactName
		} else if strings.HasPrefix(name, q) {
			score = scorePrefixName
		} else if strings.Contains(name, q) {
			score = scoreContainsName
		} else if strings.Contains(brief, q) {
			score = scoreContainsBrief
		}

		if score > 0 {
			scored = append(scored, scoredRef{ref: cmd, score: score})
		}
	}

	sort.SliceStable(scored, func(i, j int) bool {
		return scored[i].score > scored[j].score
	})

	result := make([]CommandRef, len(scored))
	for i, s := range scored {
		result[i] = s.ref
	}
	return result
}
