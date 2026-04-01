package verify

import (
	"fmt"
	"strings"

	"github.com/sd3/linuxlab/internal/challenge"
)

// Result holds the outcome of a single verification check.
type Result struct {
	Passed  bool
	Message string
}

// Verifier is the interface for all verification types.
type Verifier interface {
	Verify(rule challenge.VerifyRule) Result
}

// registry maps verify rule type names to their Verifier implementations.
var registry = map[string]Verifier{
	"file_content":   &FileContentVerifier{},
	"file_exists":    &FileExistsVerifier{},
	"command_output": &CommandOutputVerifier{},
	"exit_code":      &ExitCodeVerifier{},
	"permissions":    &PermissionsVerifier{},
	"script":         &ScriptVerifier{},
}

// RunAll executes all verification rules and returns their results.
func RunAll(rules []challenge.VerifyRule) []Result {
	results := make([]Result, len(rules))
	for i, rule := range rules {
		v, ok := registry[rule.Type]
		if !ok {
			results[i] = Result{Passed: false, Message: fmt.Sprintf("未知验证类型: %s", rule.Type)}
			continue
		}
		results[i] = v.Verify(rule)
	}
	return results
}

// AllPassed returns true if every result in the slice passed.
func AllPassed(results []Result) bool {
	for _, r := range results {
		if !r.Passed {
			return false
		}
	}
	return true
}

// compareStrings trims and compares actual vs expected, returning a Result.
func compareStrings(actual, expected, context string) Result {
	a := strings.TrimSpace(actual)
	e := strings.TrimSpace(expected)
	if a != e {
		return Result{
			Passed:  false,
			Message: fmt.Sprintf("%s不匹配\n期望: %q\n实际: %q", context, e, a),
		}
	}
	return Result{Passed: true, Message: context + "匹配"}
}
