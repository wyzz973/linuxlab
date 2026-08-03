package runner

import "github.com/sd3/linuxlab/internal/verify"

type Event struct {
	Type      string          `json:"type"`
	Message   string          `json:"message,omitempty"`
	Mode      string          `json:"mode,omitempty"`
	Passed    *bool           `json:"passed,omitempty"`
	HintsUsed int             `json:"hintsUsed,omitempty"`
	Results   []verify.Result `json:"results,omitempty"`
}
