package cli

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/sd3/linuxlab/internal/challenge"
	"github.com/sd3/linuxlab/internal/progress"
	"github.com/sd3/linuxlab/internal/reference"
)

var categoryLabels = map[string]string{
	"linux-basics":    "Linux 基础命令",
	"vim":             "Vim 操作",
	"shell-scripting": "Shell 脚本",
	"ops":             "运维实战",
	"containers":      "容器与部署",
}

type dataDump struct {
	Categories []categoryDTO         `json:"categories"`
	Progress   progress.ProgressData `json:"progress"`
	References referenceDTO          `json:"references"`
}

type categoryDTO struct {
	ID         string         `json:"id"`
	Label      string         `json:"label"`
	Total      int            `json:"total"`
	Passed     int            `json:"passed"`
	Challenges []challengeDTO `json:"challenges"`
}

type challengeDTO struct {
	ID             string         `json:"id"`
	Title          string         `json:"title"`
	Difficulty     int            `json:"difficulty"`
	Category       string         `json:"category"`
	Subcategory    string         `json:"subcategory"`
	Tags           []string       `json:"tags"`
	Description    string         `json:"description"`
	Hints          []hintDTO      `json:"hints"`
	Verify         []verifyDTO    `json:"verify"`
	SetupFiles     []setupFileDTO `json:"setup_files,omitempty"`
	ComposeFile    string         `json:"compose_file,omitempty"`
	RequiresDocker bool           `json:"requires_docker,omitempty"`
	Dir            string         `json:"dir,omitempty"`
}

type hintDTO struct {
	Level int    `json:"level"`
	Text  string `json:"text"`
}

type verifyDTO struct {
	Type    string `json:"type"`
	Path    string `json:"path,omitempty"`
	Expect  string `json:"expect,omitempty"`
	Command string `json:"command,omitempty"`
}

type setupFileDTO struct {
	Path    string `json:"path"`
	Content string `json:"content"`
}

type referenceDTO struct {
	Commands []commandDTO `json:"commands"`
}

type commandDTO struct {
	Name              string       `json:"name"`
	Brief             string       `json:"brief"`
	Examples          []exampleDTO `json:"examples"`
	RelatedChallenges []string     `json:"related_challenges,omitempty"`
}

type exampleDTO struct {
	Desc string `json:"desc"`
	Cmd  string `json:"cmd"`
}

func runData(_ context.Context, env Env) int {
	if len(env.Args) < 3 || env.Args[1] != "dump" || !hasFlag(env.Args[2:], "--json") {
		return writeError(env.Stderr, "用法: linuxlab data dump --json")
	}

	categories, err := challenge.LoadAllByCategory(env.ChallengesDir)
	if err != nil {
		return writeError(env.Stderr, fmt.Sprintf("加载题库失败: %v", err))
	}
	store, err := progress.NewStore(env.ProgressPath)
	if err != nil {
		return writeError(env.Stderr, fmt.Sprintf("加载进度失败: %v", err))
	}
	refs := reference.ReferenceData{}
	if loaded, err := reference.LoadReferences(env.RefsPath); err == nil {
		refs = *loaded
	}

	dump := dataDump{
		Categories: toCategoryDTOs(categories, store.Data.Challenges),
		Progress:   store.Data,
		References: toReferenceDTO(refs),
	}
	if err := json.NewEncoder(env.Stdout).Encode(dump); err != nil {
		return writeError(env.Stderr, fmt.Sprintf("输出数据失败: %v", err))
	}
	return 0
}

func toCategoryDTOs(categories map[string][]*challenge.Challenge, entries map[string]*progress.ChallengeEntry) []categoryDTO {
	order := []string{"containers", "linux-basics", "ops", "shell-scripting", "vim"}
	seen := map[string]bool{}
	result := make([]categoryDTO, 0, len(categories))
	for _, id := range order {
		if challenges, ok := categories[id]; ok {
			result = append(result, toCategoryDTO(id, challenges, entries))
			seen[id] = true
		}
	}
	for id, challenges := range categories {
		if !seen[id] {
			result = append(result, toCategoryDTO(id, challenges, entries))
		}
	}
	return result
}

func toCategoryDTO(id string, challenges []*challenge.Challenge, entries map[string]*progress.ChallengeEntry) categoryDTO {
	dto := categoryDTO{
		ID:         id,
		Label:      categoryLabels[id],
		Total:      len(challenges),
		Challenges: make([]challengeDTO, 0, len(challenges)),
	}
	if dto.Label == "" {
		dto.Label = id
	}
	for _, ch := range challenges {
		if entries[ch.ID] != nil && entries[ch.ID].Status == "passed" {
			dto.Passed++
		}
		dto.Challenges = append(dto.Challenges, challengeDTO{
			ID:             ch.ID,
			Title:          ch.Title,
			Difficulty:     ch.Difficulty,
			Category:       ch.Category,
			Subcategory:    ch.Subcategory,
			Tags:           ch.Tags,
			Description:    ch.Description,
			Hints:          toHintDTOs(ch.Hints),
			Verify:         toVerifyDTOs(ch.Verify),
			SetupFiles:     toSetupFileDTOs(ch.SetupFiles),
			ComposeFile:    ch.ComposeFile,
			RequiresDocker: ch.RequiresDocker,
			Dir:            ch.Dir,
		})
	}
	return dto
}

func toHintDTOs(hints []challenge.Hint) []hintDTO {
	out := make([]hintDTO, 0, len(hints))
	for _, hint := range hints {
		out = append(out, hintDTO{Level: hint.Level, Text: hint.Text})
	}
	return out
}

func toVerifyDTOs(rules []challenge.VerifyRule) []verifyDTO {
	out := make([]verifyDTO, 0, len(rules))
	for _, rule := range rules {
		out = append(out, verifyDTO{Type: rule.Type, Path: rule.Path, Expect: rule.Expect, Command: rule.Command})
	}
	return out
}

func toSetupFileDTOs(files []challenge.SetupFile) []setupFileDTO {
	out := make([]setupFileDTO, 0, len(files))
	for _, file := range files {
		out = append(out, setupFileDTO{Path: file.Path, Content: file.Content})
	}
	return out
}

func toReferenceDTO(refs reference.ReferenceData) referenceDTO {
	out := referenceDTO{Commands: make([]commandDTO, 0, len(refs.Commands))}
	for _, command := range refs.Commands {
		examples := make([]exampleDTO, 0, len(command.Examples))
		for _, example := range command.Examples {
			examples = append(examples, exampleDTO{Desc: example.Desc, Cmd: example.Cmd})
		}
		out.Commands = append(out.Commands, commandDTO{
			Name:              command.Name,
			Brief:             command.Brief,
			Examples:          examples,
			RelatedChallenges: command.RelatedChallenges,
		})
	}
	return out
}
