package tui

import (
	"fmt"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestSkillMenuOptions(t *testing.T) {
	t.Run("ScreenSkillMenu returns 5 items", func(t *testing.T) {
		m := NewModel()
		m.Screen = ScreenSkillMenu
		opts := m.GetCurrentOptions()

		if len(opts) != 5 {
			t.Errorf("expected 5 options (Browse, Install, Remove, separator, Back), got %d: %v", len(opts), opts)
		}
	})
}

func TestSkillMenuNavigation(t *testing.T) {
	t.Run("Browse (cursor 0) → Enter → ScreenSkillBrowse", func(t *testing.T) {
		m := NewModel()
		m.Screen = ScreenSkillMenu
		m.Cursor = 0

		result, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
		nm := result.(Model)

		if nm.Screen != ScreenSkillBrowse {
			t.Errorf("expected ScreenSkillBrowse, got %d", nm.Screen)
		}
		if !nm.SkillLoading {
			t.Error("expected SkillLoading=true after navigating to Browse")
		}
	})

	t.Run("Install (cursor 1) → Enter → ScreenSkillInstall", func(t *testing.T) {
		m := NewModel()
		m.Screen = ScreenSkillMenu
		m.Cursor = 1

		result, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
		nm := result.(Model)

		if nm.Screen != ScreenSkillInstall {
			t.Errorf("expected ScreenSkillInstall, got %d", nm.Screen)
		}
		if !nm.SkillLoading {
			t.Error("expected SkillLoading=true after navigating to Install")
		}
	})

	t.Run("Remove (cursor 2) → Enter → ScreenSkillRemove", func(t *testing.T) {
		m := NewModel()
		m.Screen = ScreenSkillMenu
		m.Cursor = 2

		result, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
		nm := result.(Model)

		if nm.Screen != ScreenSkillRemove {
			t.Errorf("expected ScreenSkillRemove, got %d", nm.Screen)
		}
		if !nm.SkillLoading {
			t.Error("expected SkillLoading=true after navigating to Remove")
		}
	})
}

func TestSkillMenuEscape(t *testing.T) {
	t.Run("Esc from ScreenSkillMenu → ScreenMainMenu", func(t *testing.T) {
		m := NewModel()
		m.Screen = ScreenSkillMenu

		result, _ := m.Update(tea.KeyMsg{Type: tea.KeyEsc})
		nm := result.(Model)

		if nm.Screen != ScreenMainMenu {
			t.Errorf("expected ScreenMainMenu, got %d", nm.Screen)
		}
	})
}

func TestSkillBrowseEscape(t *testing.T) {
	t.Run("Esc from ScreenSkillBrowse → ScreenSkillMenu", func(t *testing.T) {
		m := NewModel()
		m.Screen = ScreenSkillBrowse

		result, _ := m.Update(tea.KeyMsg{Type: tea.KeyEsc})
		nm := result.(Model)

		if nm.Screen != ScreenSkillMenu {
			t.Errorf("expected ScreenSkillMenu, got %d", nm.Screen)
		}
	})
}

func TestSkillInstallToggle(t *testing.T) {
	t.Run("Enter toggles skill selection on and off", func(t *testing.T) {
		m := NewModel()
		m.Screen = ScreenSkillInstall
		m.SkillList = []string{"react-19", "typescript", "tailwind-4"}
		m.SkillSelected = []bool{false, false, false}
		m.Cursor = 0

		// Toggle on
		result, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
		nm := result.(Model)

		if !nm.SkillSelected[0] {
			t.Error("expected SkillSelected[0]=true after first toggle")
		}

		// Toggle off
		nm.Cursor = 0
		result, _ = nm.Update(tea.KeyMsg{Type: tea.KeyEnter})
		nm = result.(Model)

		if nm.SkillSelected[0] {
			t.Error("expected SkillSelected[0]=false after second toggle")
		}
	})
}

func TestSkillInstallConfirmNoSelection(t *testing.T) {
	t.Run("Confirm with no selection is a no-op", func(t *testing.T) {
		m := NewModel()
		m.Screen = ScreenSkillInstall
		m.SkillList = []string{"react-19", "typescript"}
		m.SkillSelected = []bool{false, false}
		// Confirm option is at: len(SkillList) + 1 (separator at len, confirm at len+1)
		m.Cursor = len(m.SkillList) + 1

		result, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
		nm := result.(Model)

		// Screen should stay on SkillInstall (no-op)
		if nm.Screen != ScreenSkillInstall {
			t.Errorf("expected to stay on ScreenSkillInstall, got %d", nm.Screen)
		}
	})
}

func TestSkillRemoveEscape(t *testing.T) {
	t.Run("Esc from ScreenSkillRemove → ScreenSkillMenu", func(t *testing.T) {
		m := NewModel()
		m.Screen = ScreenSkillRemove

		result, _ := m.Update(tea.KeyMsg{Type: tea.KeyEsc})
		nm := result.(Model)

		if nm.Screen != ScreenSkillMenu {
			t.Errorf("expected ScreenSkillMenu, got %d", nm.Screen)
		}
	})
}

func TestSkillsLoadedMsg(t *testing.T) {
	t.Run("successful load sets SkillList and SkillSelected", func(t *testing.T) {
		m := NewModel()
		m.SkillLoading = true
		m.Screen = ScreenSkillInstall

		msg := skillsLoadedMsg{
			available: []string{"a", "b", "c"},
			installed: []string{},
			err:       nil,
		}

		result, _ := m.Update(msg)
		nm := result.(Model)

		if nm.SkillLoading {
			t.Error("expected SkillLoading=false after skillsLoadedMsg")
		}
		if len(nm.SkillList) != 3 {
			t.Errorf("expected 3 skills, got %d", len(nm.SkillList))
		}
		if len(nm.SkillSelected) != 3 {
			t.Errorf("expected 3 selection booleans, got %d", len(nm.SkillSelected))
		}
		for i, sel := range nm.SkillSelected {
			if sel {
				t.Errorf("expected SkillSelected[%d]=false, got true", i)
			}
		}
	})
}

func TestSkillsLoadedMsgError(t *testing.T) {
	t.Run("error sets SkillLoadError, list stays empty", func(t *testing.T) {
		m := NewModel()
		m.SkillLoading = true
		m.Screen = ScreenSkillInstall

		msg := skillsLoadedMsg{
			available: nil,
			installed: nil,
			err:       fmt.Errorf("network timeout"),
		}

		result, _ := m.Update(msg)
		nm := result.(Model)

		if nm.SkillLoading {
			t.Error("expected SkillLoading=false after error")
		}
		if nm.SkillLoadError == "" {
			t.Error("expected SkillLoadError to be set")
		}
		if len(nm.SkillList) != 0 {
			t.Errorf("expected empty SkillList, got %d items", len(nm.SkillList))
		}
	})
}

func TestSkillResultEnter(t *testing.T) {
	t.Run("Enter on ScreenSkillResult → ScreenSkillMenu", func(t *testing.T) {
		m := NewModel()
		m.Screen = ScreenSkillResult

		result, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
		nm := result.(Model)

		if nm.Screen != ScreenSkillMenu {
			t.Errorf("expected ScreenSkillMenu, got %d", nm.Screen)
		}
	})
}

func TestGetCurrentOptionsSkillInstall(t *testing.T) {
	t.Run("SkillInstall options = skills + separator + confirm", func(t *testing.T) {
		m := NewModel()
		m.Screen = ScreenSkillInstall
		m.SkillList = []string{"react-19", "typescript"}

		opts := m.GetCurrentOptions()

		// 2 skills + separator + confirm = 4
		if len(opts) != 4 {
			t.Errorf("expected 4 options, got %d: %v", len(opts), opts)
		}
		if !strings.Contains(opts[len(opts)-1], "Confirm") {
			t.Errorf("last option should contain 'Confirm', got %q", opts[len(opts)-1])
		}
	})
}

func TestGetCurrentOptionsSkillRemove(t *testing.T) {
	t.Run("SkillRemove options = installed skills + separator + confirm", func(t *testing.T) {
		m := NewModel()
		m.Screen = ScreenSkillRemove
		m.InstalledSkills = []string{"react-19"}

		opts := m.GetCurrentOptions()

		// 1 skill + separator + confirm = 3
		if len(opts) != 3 {
			t.Errorf("expected 3 options, got %d: %v", len(opts), opts)
		}
		if !strings.Contains(opts[len(opts)-1], "Confirm") {
			t.Errorf("last option should contain 'Confirm', got %q", opts[len(opts)-1])
		}
	})
}

func TestGetScreenTitleSkillScreens(t *testing.T) {
	screens := []Screen{
		ScreenSkillMenu,
		ScreenSkillBrowse,
		ScreenSkillInstall,
		ScreenSkillRemove,
		ScreenSkillResult,
	}

	m := NewModel()
	for _, s := range screens {
		t.Run(fmt.Sprintf("screen %d title non-empty", s), func(t *testing.T) {
			m.Screen = s
			title := m.GetScreenTitle()
			if title == "" {
				t.Errorf("expected non-empty title for screen %d", s)
			}
		})
	}
}

func TestMainMenuHasNewItems(t *testing.T) {
	t.Run("main menu contains Initialize Project and Skill Manager", func(t *testing.T) {
		m := NewModel()
		m.Screen = ScreenMainMenu
		opts := m.GetCurrentOptions()

		hasInitProject := false
		hasSkillManager := false
		for _, opt := range opts {
			if strings.Contains(opt, "Initialize Project") {
				hasInitProject = true
			}
			if strings.Contains(opt, "Skill Manager") {
				hasSkillManager = true
			}
		}

		if !hasInitProject {
			t.Error("main menu should contain 'Initialize Project'")
		}
		if !hasSkillManager {
			t.Error("main menu should contain 'Skill Manager'")
		}
	})
}
