package tui

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestExpandPath(t *testing.T) {
	t.Run("should expand ~/foo to home dir + /foo", func(t *testing.T) {
		home, err := os.UserHomeDir()
		if err != nil {
			t.Fatalf("could not get home dir: %v", err)
		}
		result := expandPath("~/foo")
		expected := filepath.Join(home, "foo")
		if result != expected {
			t.Errorf("expected %q, got %q", expected, result)
		}
	})

	t.Run("should return absolute path unchanged", func(t *testing.T) {
		result := expandPath("/abs/path")
		if result != "/abs/path" {
			t.Errorf("expected /abs/path, got %q", result)
		}
	})

	t.Run("should return relative path unchanged", func(t *testing.T) {
		result := expandPath("relative")
		if result != "relative" {
			t.Errorf("expected 'relative', got %q", result)
		}
	})
}

func TestDetectStack(t *testing.T) {
	tests := []struct {
		name     string
		file     string
		expected string
	}{
		{"go.mod → go", "go.mod", "go"},
		{"package.json → node", "package.json", "node"},
		{"angular.json → angular", "angular.json", "angular"},
		{"Cargo.toml → rust", "Cargo.toml", "rust"},
		{"pyproject.toml → python", "pyproject.toml", "python"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			f, err := os.Create(filepath.Join(dir, tt.file))
			if err != nil {
				t.Fatalf("failed to create indicator file: %v", err)
			}
			f.Close()

			result := detectStack(dir)
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}

	t.Run("empty dir → unknown", func(t *testing.T) {
		dir := t.TempDir()
		result := detectStack(dir)
		if result != "unknown" {
			t.Errorf("expected 'unknown', got %q", result)
		}
	})
}

func TestProjectPathValidation(t *testing.T) {
	t.Run("empty input + Enter sets error, screen stays", func(t *testing.T) {
		m := NewModel()
		m.Screen = ScreenProjectPath
		m.ProjectPathInput = ""

		result, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
		nm := result.(Model)

		if nm.ProjectPathError == "" {
			t.Error("expected ProjectPathError to be set for empty input")
		}
		if nm.Screen != ScreenProjectPath {
			t.Errorf("expected screen to stay at ScreenProjectPath, got %d", nm.Screen)
		}
	})

	t.Run("non-existent path + Enter sets error", func(t *testing.T) {
		m := NewModel()
		m.Screen = ScreenProjectPath
		m.ProjectPathInput = "/tmp/absolutely-does-not-exist-12345"

		result, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
		nm := result.(Model)

		if nm.ProjectPathError == "" {
			t.Error("expected ProjectPathError to be set for non-existent path")
		}
		if nm.Screen != ScreenProjectPath {
			t.Errorf("expected screen to stay at ScreenProjectPath, got %d", nm.Screen)
		}
	})

	t.Run("valid directory path + Enter advances to ScreenProjectStack", func(t *testing.T) {
		dir := t.TempDir()
		// Create a go.mod to test stack detection
		f, _ := os.Create(filepath.Join(dir, "go.mod"))
		f.Close()

		m := NewModel()
		m.Screen = ScreenProjectPath
		m.ProjectPathInput = dir

		result, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
		nm := result.(Model)

		if nm.Screen != ScreenProjectStack {
			t.Errorf("expected ScreenProjectStack, got %d", nm.Screen)
		}
		if nm.ProjectPathError != "" {
			t.Errorf("expected no error, got %q", nm.ProjectPathError)
		}
		if nm.ProjectStack != "go" {
			t.Errorf("expected ProjectStack='go', got %q", nm.ProjectStack)
		}
	})
}

func TestProjectPathBackspace(t *testing.T) {
	t.Run("backspace removes last character", func(t *testing.T) {
		m := NewModel()
		m.Screen = ScreenProjectPath
		m.ProjectPathInput = "hello"

		result, _ := m.Update(tea.KeyMsg{Type: tea.KeyBackspace})
		nm := result.(Model)

		if nm.ProjectPathInput != "hell" {
			t.Errorf("expected 'hell', got %q", nm.ProjectPathInput)
		}
	})

	t.Run("multiple backspaces eventually reach empty string without panic", func(t *testing.T) {
		m := NewModel()
		m.Screen = ScreenProjectPath
		m.ProjectPathInput = "ab"

		for i := 0; i < 5; i++ {
			result, _ := m.Update(tea.KeyMsg{Type: tea.KeyBackspace})
			m = result.(Model)
		}

		if m.ProjectPathInput != "" {
			t.Errorf("expected empty string, got %q", m.ProjectPathInput)
		}
	})
}

func TestProjectPathTyping(t *testing.T) {
	t.Run("character keys accumulate in ProjectPathInput", func(t *testing.T) {
		m := NewModel()
		m.Screen = ScreenProjectPath
		m.ProjectPathInput = ""

		chars := []rune{'/', 't', 'm', 'p'}
		for _, c := range chars {
			result, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{c}})
			m = result.(Model)
		}

		if m.ProjectPathInput != "/tmp" {
			t.Errorf("expected '/tmp', got %q", m.ProjectPathInput)
		}
	})

	t.Run("space is included in input, does not activate leader mode", func(t *testing.T) {
		m := NewModel()
		m.Screen = ScreenProjectPath
		m.ProjectPathInput = "/my"

		// Send space
		result, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{' '}})
		nm := result.(Model)

		if !strings.Contains(nm.ProjectPathInput, " ") {
			t.Errorf("expected space in input, got %q", nm.ProjectPathInput)
		}
		if nm.LeaderMode {
			t.Error("space in ScreenProjectPath should NOT activate leader mode")
		}
	})
}

func TestProjectMemoryConditionalEngram(t *testing.T) {
	t.Run("obsidian-brain without obsidian installed goes to ScreenProjectObsidianInstall", func(t *testing.T) {
		// In test env, "obsidian" binary is NOT in PATH, so it goes to install screen
		m := NewModel()
		m.Screen = ScreenProjectMemory
		m.Cursor = 0 // obsidian-brain

		result, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
		nm := result.(Model)

		if nm.ProjectMemory != "obsidian-brain" {
			t.Errorf("expected ProjectMemory='obsidian-brain', got %q", nm.ProjectMemory)
		}
		// Obsidian binary not in PATH → goes to install screen
		if nm.Screen != ScreenProjectObsidianInstall {
			t.Errorf("expected ScreenProjectObsidianInstall, got %d", nm.Screen)
		}
	})

	t.Run("vibekanban skips Engram, goes to ScreenProjectCI", func(t *testing.T) {
		m := NewModel()
		m.Screen = ScreenProjectMemory
		m.Cursor = 1 // vibekanban

		result, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
		nm := result.(Model)

		if nm.Screen != ScreenProjectCI {
			t.Errorf("expected ScreenProjectCI, got %d", nm.Screen)
		}
		if nm.ProjectMemory != "vibekanban" {
			t.Errorf("expected ProjectMemory='vibekanban', got %q", nm.ProjectMemory)
		}
	})
}

func TestObsidianInstallSelection(t *testing.T) {
	t.Run("Yes sets InstallObsidian=true and goes to Engram", func(t *testing.T) {
		m := NewModel()
		m.Screen = ScreenProjectObsidianInstall
		m.Cursor = 0 // Yes

		result, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
		nm := result.(Model)

		if !nm.Choices.InstallObsidian {
			t.Error("expected InstallObsidian=true")
		}
		if nm.Screen != ScreenProjectEngram {
			t.Errorf("expected ScreenProjectEngram, got %d", nm.Screen)
		}
	})

	t.Run("No sets InstallObsidian=false and goes to Engram", func(t *testing.T) {
		m := NewModel()
		m.Screen = ScreenProjectObsidianInstall
		m.Cursor = 1 // No

		result, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
		nm := result.(Model)

		if nm.Choices.InstallObsidian {
			t.Error("expected InstallObsidian=false")
		}
		if nm.Screen != ScreenProjectEngram {
			t.Errorf("expected ScreenProjectEngram, got %d", nm.Screen)
		}
	})
}

func TestObsidianInstallBackNav(t *testing.T) {
	t.Run("backspace on ObsidianInstall goes back to Memory", func(t *testing.T) {
		m := NewModel()
		m.Screen = ScreenProjectObsidianInstall

		result, _ := m.Update(tea.KeyMsg{Type: tea.KeyBackspace})
		nm := result.(Model)

		if nm.Screen != ScreenProjectMemory {
			t.Errorf("expected ScreenProjectMemory, got %d", nm.Screen)
		}
	})

	t.Run("backspace on Engram goes to ObsidianInstall when obsidian not in PATH", func(t *testing.T) {
		m := NewModel()
		m.Screen = ScreenProjectEngram

		result, _ := m.Update(tea.KeyMsg{Type: tea.KeyBackspace})
		nm := result.(Model)

		// Obsidian not in PATH → back goes to ObsidianInstall
		if nm.Screen != ScreenProjectObsidianInstall {
			t.Errorf("expected ScreenProjectObsidianInstall, got %d", nm.Screen)
		}
	})
}

func TestObsidianInstallScreenOptions(t *testing.T) {
	t.Run("has 2 options", func(t *testing.T) {
		m := NewModel()
		m.Screen = ScreenProjectObsidianInstall
		opts := m.GetCurrentOptions()

		if len(opts) != 2 {
			t.Errorf("expected 2 options, got %d: %v", len(opts), opts)
		}
	})

	t.Run("title is non-empty", func(t *testing.T) {
		m := NewModel()
		m.Screen = ScreenProjectObsidianInstall
		title := m.GetScreenTitle()
		if title == "" {
			t.Error("expected non-empty title")
		}
	})

	t.Run("description mentions obsidian", func(t *testing.T) {
		m := NewModel()
		m.Screen = ScreenProjectObsidianInstall
		desc := m.GetScreenDescription()
		if !strings.Contains(strings.ToLower(desc), "obsidian") {
			t.Errorf("expected description to mention obsidian, got %q", desc)
		}
	})
}

func TestProjectCIAdvancesToConfirm(t *testing.T) {
	t.Run("selecting CI advances to ScreenProjectConfirm", func(t *testing.T) {
		m := NewModel()
		m.Screen = ScreenProjectCI
		m.Cursor = 0 // GitHub Actions

		result, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
		nm := result.(Model)

		if nm.Screen != ScreenProjectConfirm {
			t.Errorf("expected ScreenProjectConfirm, got %d", nm.Screen)
		}
		if nm.ProjectCI != "github" {
			t.Errorf("expected ProjectCI='github', got %q", nm.ProjectCI)
		}
	})
}

func TestProjectEscapeBackNavigation(t *testing.T) {
	// Selection-based project screens use backspace (via handleSelectionKeys)
	// for back navigation, since ESC is intercepted by handleEscape() first.
	// ScreenProjectPath uses ESC via handleEscape() directly.

	t.Run("ScreenProjectStack → Backspace → ScreenProjectPath", func(t *testing.T) {
		m := NewModel()
		m.Screen = ScreenProjectStack

		result, _ := m.Update(tea.KeyMsg{Type: tea.KeyBackspace})
		nm := result.(Model)

		if nm.Screen != ScreenProjectPath {
			t.Errorf("expected ScreenProjectPath, got %d", nm.Screen)
		}
	})

	t.Run("ScreenProjectMemory → Backspace → ScreenProjectStack", func(t *testing.T) {
		m := NewModel()
		m.Screen = ScreenProjectMemory

		result, _ := m.Update(tea.KeyMsg{Type: tea.KeyBackspace})
		nm := result.(Model)

		if nm.Screen != ScreenProjectStack {
			t.Errorf("expected ScreenProjectStack, got %d", nm.Screen)
		}
	})

	t.Run("ScreenProjectEngram → Backspace → ScreenProjectObsidianInstall (obsidian not in PATH)", func(t *testing.T) {
		m := NewModel()
		m.Screen = ScreenProjectEngram

		result, _ := m.Update(tea.KeyMsg{Type: tea.KeyBackspace})
		nm := result.(Model)

		// Obsidian not in PATH → back goes to ObsidianInstall
		if nm.Screen != ScreenProjectObsidianInstall {
			t.Errorf("expected ScreenProjectObsidianInstall, got %d", nm.Screen)
		}
	})

	t.Run("ScreenProjectCI with memory=obsidian-brain → Backspace → ScreenProjectEngram", func(t *testing.T) {
		m := NewModel()
		m.Screen = ScreenProjectCI
		m.ProjectMemory = "obsidian-brain"

		result, _ := m.Update(tea.KeyMsg{Type: tea.KeyBackspace})
		nm := result.(Model)

		if nm.Screen != ScreenProjectEngram {
			t.Errorf("expected ScreenProjectEngram, got %d", nm.Screen)
		}
	})

	t.Run("ScreenProjectCI with memory=simple → Backspace → ScreenProjectMemory", func(t *testing.T) {
		m := NewModel()
		m.Screen = ScreenProjectCI
		m.ProjectMemory = "simple"

		result, _ := m.Update(tea.KeyMsg{Type: tea.KeyBackspace})
		nm := result.(Model)

		if nm.Screen != ScreenProjectMemory {
			t.Errorf("expected ScreenProjectMemory, got %d", nm.Screen)
		}
	})

	t.Run("ScreenProjectConfirm → Backspace → ScreenProjectCI", func(t *testing.T) {
		m := NewModel()
		m.Screen = ScreenProjectConfirm

		result, _ := m.Update(tea.KeyMsg{Type: tea.KeyBackspace})
		nm := result.(Model)

		if nm.Screen != ScreenProjectCI {
			t.Errorf("expected ScreenProjectCI, got %d", nm.Screen)
		}
	})

	t.Run("ScreenProjectPath → Esc → ScreenMainMenu", func(t *testing.T) {
		m := NewModel()
		m.Screen = ScreenProjectPath

		result, _ := m.Update(tea.KeyMsg{Type: tea.KeyEsc})
		nm := result.(Model)

		if nm.Screen != ScreenMainMenu {
			t.Errorf("expected ScreenMainMenu, got %d", nm.Screen)
		}
	})
}

func TestProjectConfirmCancel(t *testing.T) {
	t.Run("Cursor=1 (Cancel) → ScreenMainMenu", func(t *testing.T) {
		m := NewModel()
		m.Screen = ScreenProjectConfirm
		m.Cursor = 1 // Cancel

		result, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
		nm := result.(Model)

		if nm.Screen != ScreenMainMenu {
			t.Errorf("expected ScreenMainMenu, got %d", nm.Screen)
		}
	})
}

func TestGetCurrentOptionsProjectScreens(t *testing.T) {
	t.Run("ScreenProjectMemory → 5 options", func(t *testing.T) {
		m := NewModel()
		m.Screen = ScreenProjectMemory
		opts := m.GetCurrentOptions()

		if len(opts) != 5 {
			t.Errorf("expected 5 options, got %d: %v", len(opts), opts)
		}
	})

	t.Run("ScreenProjectEngram → 2 options", func(t *testing.T) {
		m := NewModel()
		m.Screen = ScreenProjectEngram
		opts := m.GetCurrentOptions()

		if len(opts) != 2 {
			t.Errorf("expected 2 options, got %d: %v", len(opts), opts)
		}
	})

	t.Run("ScreenProjectCI → 4 options", func(t *testing.T) {
		m := NewModel()
		m.Screen = ScreenProjectCI
		opts := m.GetCurrentOptions()

		if len(opts) != 4 {
			t.Errorf("expected 4 options, got %d: %v", len(opts), opts)
		}
	})

	t.Run("ScreenProjectConfirm → 2 options", func(t *testing.T) {
		m := NewModel()
		m.Screen = ScreenProjectConfirm
		opts := m.GetCurrentOptions()

		if len(opts) != 2 {
			t.Errorf("expected 2 options, got %d: %v", len(opts), opts)
		}
	})
}

func TestGetScreenTitleProjectScreens(t *testing.T) {
	screens := []Screen{
		ScreenProjectPath,
		ScreenProjectStack,
		ScreenProjectMemory,
		ScreenProjectObsidianInstall,
		ScreenProjectEngram,
		ScreenProjectCI,
		ScreenProjectConfirm,
		ScreenProjectInstalling,
		ScreenProjectResult,
	}

	m := NewModel()
	for _, s := range screens {
		t.Run("screen title non-empty", func(t *testing.T) {
			m.Screen = s
			title := m.GetScreenTitle()
			if title == "" {
				t.Errorf("expected non-empty title for screen %d", s)
			}
		})
	}
}

func TestGetScreenDescriptionProjectStack(t *testing.T) {
	t.Run("with ProjectStack=go, description contains 'go'", func(t *testing.T) {
		m := NewModel()
		m.Screen = ScreenProjectStack
		m.ProjectStack = "go"

		desc := m.GetScreenDescription()
		if !strings.Contains(strings.ToLower(desc), "go") {
			t.Errorf("expected description to contain 'go', got %q", desc)
		}
	})

	t.Run("with ProjectStack empty, description is generic", func(t *testing.T) {
		m := NewModel()
		m.Screen = ScreenProjectStack
		m.ProjectStack = ""

		desc := m.GetScreenDescription()
		if desc == "" {
			t.Error("expected non-empty generic description")
		}
		// Should NOT contain "Auto-detected"
		if strings.Contains(desc, "Auto-detected") {
			t.Errorf("empty stack should give generic description, got %q", desc)
		}
	})
}

func TestProjectResultEnter(t *testing.T) {
	t.Run("Enter on ScreenProjectResult → ScreenMainMenu", func(t *testing.T) {
		m := NewModel()
		m.Screen = ScreenProjectResult

		result, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
		nm := result.(Model)

		if nm.Screen != ScreenMainMenu {
			t.Errorf("expected ScreenMainMenu, got %d", nm.Screen)
		}
	})
}
