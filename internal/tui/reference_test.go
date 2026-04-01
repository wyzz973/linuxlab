package tui

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/sd3/linuxlab/internal/reference"
)

func sampleRefs() *reference.ReferenceData {
	return &reference.ReferenceData{
		Commands: []reference.CommandRef{
			{
				Name:  "ls",
				Brief: "列出目录内容",
				Examples: []reference.Example{
					{Desc: "列出所有文件", Cmd: "ls -la {{目录}}"},
					{Desc: "按大小排序", Cmd: "ls -lS"},
				},
				RelatedChallenges: []string{"ls-basic"},
			},
			{
				Name:  "grep",
				Brief: "搜索文本模式",
				Examples: []reference.Example{
					{Desc: "搜索模式", Cmd: "grep {{模式}} {{文件}}"},
				},
			},
			{
				Name:  "find",
				Brief: "搜索文件和目录",
				Examples: []reference.Example{
					{Desc: "按名称搜索", Cmd: "find {{路径}} -name {{模式}}"},
				},
			},
		},
	}
}

func TestReferenceModel_Init(t *testing.T) {
	refs := sampleRefs()
	m := NewReferenceModel(refs)
	rm := m.(ReferenceModel)
	if len(rm.filtered) != len(refs.Commands) {
		t.Errorf("expected %d filtered commands, got %d", len(refs.Commands), len(rm.filtered))
	}
	if rm.showDetail {
		t.Error("should not show detail initially")
	}
}

func TestReferenceModel_TypeToSearch(t *testing.T) {
	refs := sampleRefs()
	m := NewReferenceModel(refs)

	// Type "ls"
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'l'}})
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}})

	rm := m.(ReferenceModel)
	if rm.query != "ls" {
		t.Errorf("expected query 'ls', got %q", rm.query)
	}
	if len(rm.filtered) != 1 {
		t.Errorf("expected 1 filtered result, got %d", len(rm.filtered))
	}
	if rm.filtered[0].Name != "ls" {
		t.Errorf("expected 'ls' in filtered, got %q", rm.filtered[0].Name)
	}
}

func TestReferenceModel_Enter_ShowsDetail(t *testing.T) {
	refs := sampleRefs()
	m := NewReferenceModel(refs)

	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	rm := m.(ReferenceModel)
	if !rm.showDetail {
		t.Error("expected showDetail to be true after Enter")
	}

	view := m.View()
	if !strings.Contains(view, "ls") {
		t.Error("detail view should contain command name")
	}
}

func TestReferenceModel_Escape_FromDetail(t *testing.T) {
	refs := sampleRefs()
	m := NewReferenceModel(refs)

	// Enter detail
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	rm := m.(ReferenceModel)
	if !rm.showDetail {
		t.Fatal("expected detail mode")
	}

	// Escape back to list
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyEsc})
	rm = m.(ReferenceModel)
	if rm.showDetail {
		t.Error("expected to be back in list mode after Esc")
	}
}

func TestReferenceModel_Escape_FromList(t *testing.T) {
	refs := sampleRefs()
	m := NewReferenceModel(refs)

	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEsc})
	if cmd == nil {
		t.Fatal("expected a command from Esc in list mode")
	}
	msg := cmd()
	if _, ok := msg.(GoBackMsg); !ok {
		t.Fatalf("expected GoBackMsg, got %T", msg)
	}
}
