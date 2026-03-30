package tui

import (
	"testing"

	"github.com/sd3/linuxlab/internal/challenge"
)

func TestAppModel_StartsAtMenu(t *testing.T) {
	store := tmpStore(t)
	cats := sampleCategories()
	m := NewAppModel(cats, store)
	app := m.(AppModel)
	if app.screen != screenMenu {
		t.Errorf("expected screenMenu, got %d", app.screen)
	}
}

func TestAppModel_MenuToPractice(t *testing.T) {
	store := tmpStore(t)
	cats := sampleCategories()
	m := NewAppModel(cats, store)
	updated, _ := m.Update(MenuChoiceMsg{Choice: "practice"})
	app := updated.(AppModel)
	if app.screen != screenModules {
		t.Errorf("expected screenModules, got %d", app.screen)
	}
}

func TestAppModel_MenuToSkillMap(t *testing.T) {
	store := tmpStore(t)
	cats := sampleCategories()
	m := NewAppModel(cats, store)
	updated, _ := m.Update(MenuChoiceMsg{Choice: "skillmap"})
	app := updated.(AppModel)
	if app.screen != screenSkillMap {
		t.Errorf("expected screenSkillMap, got %d", app.screen)
	}
}

func TestAppModel_GoBack_FromModules(t *testing.T) {
	store := tmpStore(t)
	cats := sampleCategories()
	m := NewAppModel(cats, store)
	// Navigate to modules first
	m, _ = m.Update(MenuChoiceMsg{Choice: "practice"})
	// Go back
	m, _ = m.Update(GoBackMsg{})
	app := m.(AppModel)
	if app.screen != screenMenu {
		t.Errorf("expected screenMenu, got %d", app.screen)
	}
}

func TestAppModel_GoBack_FromChallenges(t *testing.T) {
	store := tmpStore(t)
	cats := sampleCategories()
	m := NewAppModel(cats, store)
	// Navigate to modules, then challenges
	m, _ = m.Update(MenuChoiceMsg{Choice: "practice"})
	m, _ = m.Update(ModuleSelectedMsg{Category: "linux-basics"})
	app := m.(AppModel)
	if app.screen != screenChallenges {
		t.Fatalf("expected screenChallenges, got %d", app.screen)
	}
	// Go back
	m, _ = m.Update(GoBackMsg{})
	app = m.(AppModel)
	if app.screen != screenModules {
		t.Errorf("expected screenModules, got %d", app.screen)
	}
}

func TestAppModel_GoBack_FromDetail(t *testing.T) {
	store := tmpStore(t)
	cats := sampleCategories()
	m := NewAppModel(cats, store)
	// Navigate to modules, then challenges, then detail
	m, _ = m.Update(MenuChoiceMsg{Choice: "practice"})
	m, _ = m.Update(ModuleSelectedMsg{Category: "linux-basics"})
	ch := &challenge.Challenge{ID: "ls-basic", Title: "ls 基础", Category: "linux-basics"}
	m, _ = m.Update(ChallengeSelectedMsg{Challenge: ch})
	app := m.(AppModel)
	if app.screen != screenDetail {
		t.Fatalf("expected screenDetail, got %d", app.screen)
	}
	// Go back
	m, _ = m.Update(GoBackMsg{})
	app = m.(AppModel)
	if app.screen != screenChallenges {
		t.Errorf("expected screenChallenges, got %d", app.screen)
	}
}
