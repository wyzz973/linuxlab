package sandbox

import (
	"os"
	"path/filepath"

	"github.com/sd3/linuxlab/internal/challenge"
	"github.com/sd3/linuxlab/internal/verify"
)

// VimRunner manages file setup and verification for Vim challenges.
type VimRunner struct {
	WorkDir string
}

// PrepareFiles creates the setup files on disk and returns their absolute paths.
func (r *VimRunner) PrepareFiles(files []challenge.SetupFile) ([]string, error) {
	var paths []string
	for _, f := range files {
		fullPath := filepath.Join(r.WorkDir, f.Path)
		if err := os.MkdirAll(filepath.Dir(fullPath), 0o755); err != nil {
			return nil, err
		}
		if err := os.WriteFile(fullPath, []byte(f.Content), 0o644); err != nil {
			return nil, err
		}
		paths = append(paths, fullPath)
	}
	return paths, nil
}

// VimArgs returns the command arguments to open a file in Vim.
func (r *VimRunner) VimArgs(filePath string) []string {
	return []string{"vim", filePath}
}

// CheckResult runs verification rules and returns whether all passed.
func (r *VimRunner) CheckResult(rules []challenge.VerifyRule) bool {
	results := verify.RunAll(rules)
	return verify.AllPassed(results)
}
