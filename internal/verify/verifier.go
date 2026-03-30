package verify

import "github.com/sd3/linuxlab/internal/challenge"

// Result holds the outcome of a single verification check.
type Result struct {
	Passed  bool
	Message string
}

// Verifier is the interface for all verification types.
type Verifier interface {
	Verify(rule challenge.VerifyRule) Result
}
