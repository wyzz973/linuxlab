package cli

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/sd3/linuxlab/internal/challenge"
	"github.com/sd3/linuxlab/internal/reference"
	"github.com/sd3/linuxlab/internal/sandbox"
)

type doctorResult struct {
	Docker          bool   `json:"docker"`
	Challenges      int    `json:"challenges"`
	References      int    `json:"references"`
	ChallengesError string `json:"challenges_error,omitempty"`
	ReferencesError string `json:"references_error,omitempty"`
}

func runDoctor(_ context.Context, env Env) int {
	result := doctorResult{Docker: sandbox.DockerAvailable()}

	challenges, chErr := challenge.LoadAll(env.ChallengesDir)
	if chErr != nil {
		result.ChallengesError = chErr.Error()
	} else {
		result.Challenges = len(challenges)
	}

	refs, refErr := reference.LoadReferences(env.RefsPath)
	if refErr != nil {
		result.ReferencesError = refErr.Error()
	} else {
		result.References = len(refs.Commands)
	}

	exitCode := 0
	if chErr != nil || refErr != nil {
		exitCode = 1
	}

	if hasFlag(env.Args[1:], "--json") {
		if err := json.NewEncoder(env.Stdout).Encode(result); err != nil {
			return writeError(env.Stderr, fmt.Sprintf("输出 doctor JSON 失败: %v", err))
		}
		return exitCode
	}

	fmt.Fprintf(env.Stdout, "Docker: %v\n", result.Docker)
	if chErr != nil {
		fmt.Fprintf(env.Stdout, "题目: 加载失败: %v\n", chErr)
	} else {
		fmt.Fprintf(env.Stdout, "题目: %d\n", result.Challenges)
	}
	if refErr != nil {
		fmt.Fprintf(env.Stdout, "命令速查: 加载失败: %v\n", refErr)
	} else {
		fmt.Fprintf(env.Stdout, "命令速查: %d\n", result.References)
	}
	return exitCode
}
