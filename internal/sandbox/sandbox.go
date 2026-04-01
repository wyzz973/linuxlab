package sandbox

import "context"

// Sandbox is the common interface for all challenge execution backends.
type Sandbox interface {
	Exec(ctx context.Context, command string) (string, int, error)
	Destroy(ctx context.Context) error
	InteractiveShellArgs() []string
}
