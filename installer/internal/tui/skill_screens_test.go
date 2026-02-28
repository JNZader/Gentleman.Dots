package tui

import (
	"fmt"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestSkillMenuOptions(t *testing.T) {
	t.Run("ScreenSkillMenu returns 6 items", func(t *testing.T) {
		m := NewModel()
		m.Screen = ScreenSkillMenu
		opts := m.GetCurrentOptions()

		// Browse, Install, Remove, Update, separator, Back = 6
		if len(opts) != 6 {
			t.Errorf("expected 6 options (Browse, Install, Remove, Update, separator, Back), got %d: %v", len(opts), opts)
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

	t.Run("Update (cursor 3) → Enter → ScreenSkillUpdate", func(t *testing.T) {
		m := NewModel()
		m.Screen = ScreenSkillMenu
		m.Cursor = 3

		result, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
		nm := result.(Model)

		if nm.Screen != ScreenSkillUpdate {
			t.Errorf("expected ScreenSkillUpdate, got %d", nm.Screen)
		}
		if !nm.SkillLoading {
			t.Error("expected SkillLoading=true after navigating to Update")
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
		m.SkillCatalog = []SkillInfo{
			{Name: "react-19", Category: "curated", Installed: false},
			{Name: "typescript", Category: "curated", Installed: false},
			{Name: "tailwind-4", Category: "curated", Installed: false},
		}
		m.SkillSelected = []bool{false, false, false}
		// Options: [0] Select All, [1] 📦 Curated, [2] react-19, [3] typescript, [4] tailwind-4, [5] sep, [6] Confirm
		// Cursor at 2 = first skill item (index 0 in SkillSelected)
		m.Cursor = 2

		// Toggle on
		result, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
		nm := result.(Model)

		if !nm.SkillSelected[0] {
			t.Error("expected SkillSelected[0]=true after first toggle")
		}

		// Toggle off
		nm.Cursor = 2
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
		m.SkillCatalog = []SkillInfo{
			{Name: "react-19", Category: "curated", Installed: false},
			{Name: "typescript", Category: "curated", Installed: false},
		}
		m.SkillSelected = []bool{false, false}
		// Options: [0] Select All, [1] 📦 Curated, [2] react-19, [3] typescript, [4] sep, [5] Confirm
		m.Cursor = 5

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
	t.Run("successful load sets SkillCatalog and SkillSelected", func(t *testing.T) {
		m := NewModel()
		m.SkillLoading = true
		m.Screen = ScreenSkillInstall

		msg := skillsLoadedMsg{
			skills: []SkillInfo{
				{Name: "a", Category: "curated", Installed: false},
				{Name: "b", Category: "curated", Installed: false},
				{Name: "c", Category: "community", Installed: false},
			},
			err: nil,
		}

		result, _ := m.Update(msg)
		nm := result.(Model)

		if nm.SkillLoading {
			t.Error("expected SkillLoading=false after skillsLoadedMsg")
		}
		if len(nm.SkillCatalog) != 3 {
			t.Errorf("expected 3 skills in catalog, got %d", len(nm.SkillCatalog))
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
	t.Run("error sets SkillLoadError, catalog stays empty", func(t *testing.T) {
		m := NewModel()
		m.SkillLoading = true
		m.Screen = ScreenSkillInstall

		msg := skillsLoadedMsg{
			skills: nil,
			err:    fmt.Errorf("network timeout"),
		}

		result, _ := m.Update(msg)
		nm := result.(Model)

		if nm.SkillLoading {
			t.Error("expected SkillLoading=false after error")
		}
		if nm.SkillLoadError == "" {
			t.Error("expected SkillLoadError to be set")
		}
		if len(nm.SkillCatalog) != 0 {
			t.Errorf("expected empty SkillCatalog, got %d items", len(nm.SkillCatalog))
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
	t.Run("SkillInstall options include Select All, group headers, skills, separator, confirm", func(t *testing.T) {
		m := NewModel()
		m.Screen = ScreenSkillInstall
		m.SkillCatalog = []SkillInfo{
			{Name: "react-19", Description: "React 19 patterns", Category: "curated", Installed: false},
			{Name: "typescript", Description: "TypeScript types", Category: "curated", Installed: false},
		}

		opts := m.GetCurrentOptions()

		// Select All + 📦 Curated + 2 skills + separator + confirm = 6
		if len(opts) != 6 {
			t.Errorf("expected 6 options, got %d: %v", len(opts), opts)
		}
		if !strings.Contains(opts[len(opts)-1], "Confirm") {
			t.Errorf("last option should contain 'Confirm', got %q", opts[len(opts)-1])
		}
	})
}

func TestGetCurrentOptionsSkillRemove(t *testing.T) {
	t.Run("SkillRemove options = Select All + category header + installed skills + separator + confirm", func(t *testing.T) {
		m := NewModel()
		m.Screen = ScreenSkillRemove
		m.SkillCatalog = []SkillInfo{
			{Name: "react-19", Description: "React 19 patterns", Category: "curated", Installed: true},
		}

		opts := m.GetCurrentOptions()

		// Select All + header + 1 skill + separator + confirm = 5
		if len(opts) != 5 {
			t.Errorf("expected 5 options, got %d: %v", len(opts), opts)
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
		ScreenSkillUpdate,
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

func TestSkillOptionToIndex(t *testing.T) {
	t.Run("maps cursor to correct skill index skipping headers", func(t *testing.T) {
		options := []string{
			"✅ Select All",
			"📦 Curated",
			"react-19 — React 19 patterns",
			"typescript — TypeScript types",
			"🌐 Community",
			"electron — Electron patterns",
			"─────────────",
			"✅ Confirm installation",
		}

		// Select All → -1
		if idx := skillOptionToIndex(options, 0); idx != -1 {
			t.Errorf("expected -1 for Select All, got %d", idx)
		}
		// Group header → -1
		if idx := skillOptionToIndex(options, 1); idx != -1 {
			t.Errorf("expected -1 for group header, got %d", idx)
		}
		// react-19 → 0
		if idx := skillOptionToIndex(options, 2); idx != 0 {
			t.Errorf("expected 0 for react-19, got %d", idx)
		}
		// typescript → 1
		if idx := skillOptionToIndex(options, 3); idx != 1 {
			t.Errorf("expected 1 for typescript, got %d", idx)
		}
		// community header → -1
		if idx := skillOptionToIndex(options, 4); idx != -1 {
			t.Errorf("expected -1 for community header, got %d", idx)
		}
		// electron → 2
		if idx := skillOptionToIndex(options, 5); idx != 2 {
			t.Errorf("expected 2 for electron, got %d", idx)
		}
	})
}

func TestParseSkillFrontmatter(t *testing.T) {
	t.Run("returns empty for non-existent file", func(t *testing.T) {
		name, desc := parseSkillFrontmatter("/tmp/nonexistent-skill-test-file.md")
		if name != "" || desc != "" {
			t.Errorf("expected empty name and desc for missing file, got %q %q", name, desc)
		}
	})
}

func TestTruncateDesc(t *testing.T) {
	t.Run("short string unchanged", func(t *testing.T) {
		result := truncateDesc("short", 60)
		if result != "short" {
			t.Errorf("expected 'short', got %q", result)
		}
	})

	t.Run("long string truncated with ellipsis", func(t *testing.T) {
		long := strings.Repeat("a", 100)
		result := truncateDesc(long, 60)
		// Should be maxLen-1 chars (59 'a') + ellipsis
		if !strings.HasSuffix(result, "…") {
			t.Error("expected truncated string to end with …")
		}
		// Check that the non-ellipsis part is exactly 59 chars
		withoutEllipsis := strings.TrimSuffix(result, "…")
		if len(withoutEllipsis) != 59 {
			t.Errorf("expected 59 chars before ellipsis, got %d", len(withoutEllipsis))
		}
	})
}

func TestSkillRemoveCategoryToggleWithLocalSkills(t *testing.T) {
	t.Run("Backend header toggles category instead of going back", func(t *testing.T) {
		m := NewModel()
		m.Screen = ScreenSkillRemove
		m.SkillCatalog = []SkillInfo{
			{Name: "api-gateway", Description: "API Gateway", Category: "local:backend", Installed: true},
			{Name: "bff-concepts", Description: "BFF pattern", Category: "local:backend", Installed: true},
		}
		m.SkillSelected = make([]bool, len(m.getInstalledSkills()))

		opts := m.GetCurrentOptions()
		// Find the "🏠 Backend" header
		headerIdx := -1
		for i, o := range opts {
			if strings.Contains(o, "Backend") && strings.HasPrefix(o, "🏠") {
				headerIdx = i
				break
			}
		}
		if headerIdx == -1 {
			t.Fatal("Backend header not found in options")
		}

		// Press Enter on header — should NOT go back to SkillMenu
		m.Cursor = headerIdx
		result, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
		nm := result.(Model)

		if nm.Screen != ScreenSkillRemove {
			t.Errorf("expected to stay on ScreenSkillRemove, got screen %d (went back!)", nm.Screen)
		}
		// All backend skills should be selected
		for i := range nm.SkillSelected {
			if !nm.SkillSelected[i] {
				t.Errorf("expected SkillSelected[%d] to be true", i)
			}
		}
	})

	t.Run("toggling bff-concepts in local:backend does not panic", func(t *testing.T) {
		m := NewModel()
		m.Screen = ScreenSkillRemove
		m.SkillCatalog = []SkillInfo{
			{Name: "react-19", Description: "React 19 patterns", Category: "curated", Installed: true},
			{Name: "typescript", Description: "TypeScript types", Category: "curated", Installed: true},
			{Name: "api-gateway", Description: "API Gateway", Category: "local:backend", Installed: true},
			{Name: "bff-concepts", Description: "BFF pattern", Category: "local:backend", Installed: true},
			{Name: "chi-router", Description: "Chi router", Category: "local:backend", Installed: true},
		}
		installed := m.getInstalledSkills()
		m.SkillSelected = make([]bool, len(installed))

		opts := m.GetCurrentOptions()
		t.Logf("options: %v", opts)
		t.Logf("SkillSelected len: %d, installed len: %d", len(m.SkillSelected), len(installed))

		// Find bff-concepts in options
		bffIdx := -1
		for i, o := range opts {
			if strings.Contains(o, "bff-concepts") {
				bffIdx = i
				break
			}
		}
		if bffIdx == -1 {
			t.Fatal("bff-concepts not found in options")
		}
		t.Logf("bff-concepts at option index %d", bffIdx)

		// Position cursor on bff-concepts and press Enter
		m.Cursor = bffIdx
		result, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
		nm := result.(Model)

		// Verify it toggled correctly
		skillIdx := skillOptionToIndex(opts, bffIdx)
		t.Logf("skillOptionToIndex for bff-concepts: %d", skillIdx)
		if skillIdx < 0 || skillIdx >= len(nm.SkillSelected) {
			t.Fatalf("skillOptionToIndex returned %d, SkillSelected len %d", skillIdx, len(nm.SkillSelected))
		}
		if !nm.SkillSelected[skillIdx] {
			t.Error("expected bff-concepts to be selected after toggle")
		}
	})
}
