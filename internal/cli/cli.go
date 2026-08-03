package cli

import (
	"context"
	"fmt"
	"io"
)

type Env struct {
	Args          []string
	Stdout        io.Writer
	Stderr        io.Writer
	Stdin         io.Reader
	ChallengesDir string
	RefsPath      string
	ProgressPath  string
}

func Run(ctx context.Context, env Env) int {
	if env.Stdout == nil {
		env.Stdout = io.Discard
	}
	if env.Stderr == nil {
		env.Stderr = io.Discard
	}
	if len(env.Args) == 0 {
		return -1
	}
	switch env.Args[0] {
	case "challenge":
		return runChallenge(ctx, env)
	case "doctor":
		return runDoctor(ctx, env)
	case "data":
		return runData(ctx, env)
	default:
		return writeError(env.Stderr, "未知命令: "+env.Args[0])
	}
}

func writeError(stderr io.Writer, message string) int {
	if stderr != nil {
		fmt.Fprintln(stderr, message)
	}
	return 1
}

func hasFlag(args []string, flag string) bool {
	for _, arg := range args {
		if arg == flag {
			return true
		}
	}
	return false
}
