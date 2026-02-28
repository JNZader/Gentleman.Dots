# Design: zed-editor-support

## Architecture

Zed follows the exact same integration pattern as Neovim. It touches the same 7 files with identical structural changes:

```
model.go     -> Screen constant + UserChoices field + 3 switch cases
update.go    -> Forward navigation + backward navigation + escape handler + cursor clamp
view.go      -> Render case + progress bar
installer.go -> stepInstallZed() + step registration
main.go      -> CLI flag + non-interactive wiring + help text
non_interactive.go -> buildStepsForChoices() entry
system/exec.go     -> ConfigPaths() entry
```

The screen sits between `ScreenNvimSelect` and `ScreenAIToolsSelect` in the wizard flow. All existing AI-related screens shift their step numbers by +1 in titles.

## Changes by File

### 1. `installer/internal/tui/model.go`

#### 1a. Add `ScreenZedSelect` constant (line 26, after `ScreenNvimSelect`)

Current code at lines 25-26:
```go
	ScreenNvimSelect
	ScreenInstalling
```

Insert `ScreenZedSelect` between them:
```go
	ScreenNvimSelect
	ScreenZedSelect  // NEW: Zed editor selection
	ScreenInstalling
```

> **Note:** All subsequent `iota` values shift by +1. This is safe because all references use symbolic names (`ScreenInstalling`, `ScreenComplete`, etc.), never raw integers. The iota reordering affects `ScreenInstalling` (was 9, becomes 10) and everything after it.

#### 1b. Add `InstallZed` field to `UserChoices` struct (line 122, after `InstallNvim`)

Current code at lines 121-122:
```go
	InstallNvim  bool
	CreateBackup bool // Whether to backup existing configs
```

Add:
```go
	InstallNvim  bool
	InstallZed   bool // Whether to install Zed editor with config
	CreateBackup bool // Whether to backup existing configs
```

#### 1c. Add `ScreenZedSelect` to `GetCurrentOptions()` (line 397, after `ScreenNvimSelect` case)

Current code at lines 397-399:
```go
	case ScreenNvimSelect:
		return []string{"Yes, install Neovim with config", "No, skip Neovim", "─────────────", "ℹ️  Learn about Neovim", "⌨️  View Keymaps", "📖 LazyVim Guide"}
	case ScreenAIToolsSelect:
```

Insert after the `ScreenNvimSelect` case:
```go
	case ScreenZedSelect:
		return []string{"Yes, install Zed with config", "No, skip Zed"}
```

#### 1d. Add `ScreenZedSelect` to `GetScreenTitle()` (line 566, after `ScreenNvimSelect` case)

Current code at lines 564-567:
```go
	case ScreenNvimSelect:
		return "Step 6: Neovim Configuration"
	case ScreenAIToolsSelect:
		return "Step 7: AI Coding Tools"
```

Insert after `ScreenNvimSelect` case and renumber AI screens:
```go
	case ScreenNvimSelect:
		return "Step 6: Neovim Configuration"
	case ScreenZedSelect:
		return "Step 7: Zed Editor"
	case ScreenAIToolsSelect:
		return "Step 8: AI Coding Tools"
	case ScreenAIFrameworkConfirm:
		return "Step 9: AI Framework"
	case ScreenAIFrameworkPreset:
		return "Step 9: Choose Framework Preset"
	case ScreenAIFrameworkCategories:
		return "Step 9: Select Module Categories"
	case ScreenAIFrameworkCategoryItems:
		if m.SelectedModuleCategory >= 0 && m.SelectedModuleCategory < len(moduleCategories) {
			cat := moduleCategories[m.SelectedModuleCategory]
			return fmt.Sprintf("Step 9: %s %s", cat.Icon, cat.Label)
		}
		return "Step 9: Select Modules"
```

> This replaces lines 564-579 (the existing `ScreenNvimSelect` through `ScreenAIFrameworkCategoryItems` block).

#### 1e. Add `ScreenZedSelect` to `GetScreenDescription()` (line 711, after `ScreenNvimSelect` case)

Current code at lines 711-713:
```go
	case ScreenNvimSelect:
		return "Includes LSP, TreeSitter, and Gentleman config"
	case ScreenAIToolsSelect:
```

Insert after:
```go
	case ScreenZedSelect:
		return "High-performance editor with Vim mode and AI agent support"
```

#### 1f. Add Zed step to `SetupInstallSteps()` (line 1053, after nvim step block)

Current code at lines 1045-1064:
```go
	// Neovim (not interactive - brew doesn't need password)
	if m.Choices.InstallNvim {
		m.Steps = append(m.Steps, InstallStep{
			ID:          "nvim",
			Name:        "Install Neovim",
			Description: "Editor with config",
			Status:      StatusPending,
		})
	}

	// AI Tools: Claude Code + OpenCode (not interactive)
	if len(m.Choices.AITools) > 0 {
```

Insert after the nvim block (after line 1053):
```go
	// Zed editor (not interactive)
	if m.Choices.InstallZed {
		m.Steps = append(m.Steps, InstallStep{
			ID:          "zed",
			Name:        "Install Zed",
			Description: "Editor with Vim mode",
			Status:      StatusPending,
		})
	}
```

---

### 2. `installer/internal/tui/update.go`

#### 2a. Add `ScreenZedSelect` to the `handleSelectionKeys` dispatch (line 786)

Current code at lines 786-788:
```go
	case ScreenOSSelect, ScreenTerminalSelect, ScreenFontSelect, ScreenShellSelect, ScreenWMSelect, ScreenNvimSelect, ScreenAIFrameworkConfirm, ScreenAIFrameworkPreset, ScreenGhosttyWarning,
		ScreenProjectStack, ScreenProjectMemory, ScreenProjectObsidianInstall, ScreenProjectEngram, ScreenProjectCI, ScreenProjectConfirm, ScreenSkillMenu, ScreenLearnMenu:
		return m.handleSelectionKeys(key)
```

Add `ScreenZedSelect` to the list:
```go
	case ScreenOSSelect, ScreenTerminalSelect, ScreenFontSelect, ScreenShellSelect, ScreenWMSelect, ScreenNvimSelect, ScreenZedSelect, ScreenAIFrameworkConfirm, ScreenAIFrameworkPreset, ScreenGhosttyWarning,
		ScreenProjectStack, ScreenProjectMemory, ScreenProjectObsidianInstall, ScreenProjectEngram, ScreenProjectCI, ScreenProjectConfirm, ScreenSkillMenu, ScreenLearnMenu:
		return m.handleSelectionKeys(key)
```

#### 2b. Add `ScreenZedSelect` to the escape handler list (line 911)

Current code at line 911:
```go
	case ScreenOSSelect, ScreenTerminalSelect, ScreenFontSelect, ScreenShellSelect, ScreenWMSelect, ScreenNvimSelect, ScreenAIToolsSelect, ScreenAIFrameworkConfirm, ScreenAIFrameworkPreset, ScreenAIFrameworkCategories, ScreenAIFrameworkCategoryItems:
		return m.goBackInstallStep()
```

Add `ScreenZedSelect`:
```go
	case ScreenOSSelect, ScreenTerminalSelect, ScreenFontSelect, ScreenShellSelect, ScreenWMSelect, ScreenNvimSelect, ScreenZedSelect, ScreenAIToolsSelect, ScreenAIFrameworkConfirm, ScreenAIFrameworkPreset, ScreenAIFrameworkCategories, ScreenAIFrameworkCategoryItems:
		return m.goBackInstallStep()
```

#### 2c. Modify `ScreenNvimSelect` forward navigation (lines 1356-1365)

Current code at lines 1356-1365:
```go
	case ScreenNvimSelect:
		m.Choices.InstallNvim = m.Cursor == 0
		// Proceed to AI tools selection (skip on Termux)
		if m.SystemInfo.IsTermux {
			// Termux doesn't support AI tools, skip to backup/install
			return m.proceedToBackupOrInstall()
		}
		m.Screen = ScreenAIToolsSelect
		m.Cursor = 0
		m.AIToolSelected = make([]bool, len(aiToolIDMap))
```

Change to navigate to `ScreenZedSelect` instead:
```go
	case ScreenNvimSelect:
		m.Choices.InstallNvim = m.Cursor == 0
		// Proceed to Zed selection (skip on Termux — Zed needs GUI)
		if m.SystemInfo.IsTermux {
			// Termux doesn't support Zed or AI tools, skip to backup/install
			return m.proceedToBackupOrInstall()
		}
		m.Screen = ScreenZedSelect
		m.Cursor = 0
```

#### 2d. Add `ScreenZedSelect` forward navigation handler (after the ScreenNvimSelect case, ~line 1366)

Insert new case:
```go
	case ScreenZedSelect:
		m.Choices.InstallZed = m.Cursor == 0
		m.Screen = ScreenAIToolsSelect
		m.Cursor = 0
		m.AIToolSelected = make([]bool, len(aiToolIDMap))
```

#### 2e. Add `ScreenZedSelect` backward navigation in `goBackInstallStep()` (line 1169)

Current code at lines 1169-1173:
```go
	case ScreenAIToolsSelect:
		m.Screen = ScreenNvimSelect
		m.Cursor = 0
		m.Choices.AITools = nil
		m.AIToolSelected = nil
```

Change `ScreenAIToolsSelect` to go back to `ScreenZedSelect`:
```go
	case ScreenAIToolsSelect:
		m.Screen = ScreenZedSelect
		m.Cursor = 0
		m.Choices.AITools = nil
		m.AIToolSelected = nil
```

And add `ScreenZedSelect` backward navigation (insert after `ScreenNvimSelect` case, ~line 1167):
```go
	case ScreenZedSelect:
		m.Screen = ScreenNvimSelect
		m.Cursor = 0
		m.Choices.InstallZed = false
```

---

### 3. `installer/internal/tui/view.go`

#### 3a. Add `ScreenZedSelect` to render dispatch (line 54)

Current code at line 54:
```go
	case ScreenOSSelect, ScreenTerminalSelect, ScreenFontSelect, ScreenShellSelect, ScreenWMSelect, ScreenNvimSelect, ScreenAIFrameworkConfirm, ScreenAIFrameworkPreset, ScreenGhosttyWarning:
		s.WriteString(m.renderSelection())
```

Add `ScreenZedSelect`:
```go
	case ScreenOSSelect, ScreenTerminalSelect, ScreenFontSelect, ScreenShellSelect, ScreenWMSelect, ScreenNvimSelect, ScreenZedSelect, ScreenAIFrameworkConfirm, ScreenAIFrameworkPreset, ScreenGhosttyWarning:
		s.WriteString(m.renderSelection())
```

#### 3b. Update `renderStepProgress()` steps array and index mapping (lines 253-274)

Current code at lines 253-274:
```go
func (m Model) renderStepProgress() string {
	steps := []string{"OS", "Terminal", "Font", "Shell", "WM", "Nvim", "AI Tools", "Framework"}
	currentIdx := 0

	switch m.Screen {
	case ScreenOSSelect:
		currentIdx = 0
	case ScreenTerminalSelect:
		currentIdx = 1
	case ScreenFontSelect:
		currentIdx = 2
	case ScreenShellSelect:
		currentIdx = 3
	case ScreenWMSelect:
		currentIdx = 4
	case ScreenNvimSelect:
		currentIdx = 5
	case ScreenAIToolsSelect:
		currentIdx = 6
	case ScreenAIFrameworkConfirm, ScreenAIFrameworkPreset, ScreenAIFrameworkCategories, ScreenAIFrameworkCategoryItems:
		currentIdx = 7
	}
```

Replace with:
```go
func (m Model) renderStepProgress() string {
	steps := []string{"OS", "Terminal", "Font", "Shell", "WM", "Nvim", "Zed", "AI Tools", "Framework"}
	currentIdx := 0

	switch m.Screen {
	case ScreenOSSelect:
		currentIdx = 0
	case ScreenTerminalSelect:
		currentIdx = 1
	case ScreenFontSelect:
		currentIdx = 2
	case ScreenShellSelect:
		currentIdx = 3
	case ScreenWMSelect:
		currentIdx = 4
	case ScreenNvimSelect:
		currentIdx = 5
	case ScreenZedSelect:
		currentIdx = 6
	case ScreenAIToolsSelect:
		currentIdx = 7
	case ScreenAIFrameworkConfirm, ScreenAIFrameworkPreset, ScreenAIFrameworkCategories, ScreenAIFrameworkCategoryItems:
		currentIdx = 8
	}
```

---

### 4. `installer/internal/tui/installer.go`

#### 4a. Add `"zed"` case to `executeStep()` dispatch (line 69, after `"nvim"`)

Current code at lines 68-71:
```go
	case "nvim":
		return stepInstallNvim(m)
	case "aitools":
		return stepInstallAITools(m)
```

Insert:
```go
	case "nvim":
		return stepInstallNvim(m)
	case "zed":
		return stepInstallZed(m)
	case "aitools":
		return stepInstallAITools(m)
```

#### 4b. Add `stepInstallZed()` function (after `stepInstallNvim()`, after line 1080)

The function follows the same pattern as `stepInstallNvim()` (lines 980-1080):

```go
func stepInstallZed(m *Model) error {
	homeDir := os.Getenv("HOME")
	repoDir := "Gentleman.Dots"
	stepID := "zed"

	// Skip on Termux — Zed requires GUI with Vulkan
	if m.SystemInfo.IsTermux {
		SendLog(stepID, "Skipping Zed on Termux (requires GUI with Vulkan)")
		return nil
	}

	// Install Zed binary
	SendLog(stepID, "Installing Zed editor...")
	var result *system.ExecResult
	switch m.SystemInfo.OS {
	case system.OSMac:
		result = system.RunBrewWithLogs("install --cask zed", nil, func(line string) {
			SendLog(stepID, line)
		})
	case system.OSArch:
		result = system.RunSudoWithLogs("pacman -S --noconfirm zed", nil, func(line string) {
			SendLog(stepID, line)
		})
	case system.OSDebian, system.OSLinux:
		result = system.RunWithLogs("bash -c 'curl -f https://zed.dev/install.sh | sh'", nil, func(line string) {
			SendLog(stepID, line)
		})
	case system.OSFedora:
		result = system.RunWithLogs("bash -c 'curl -f https://zed.dev/install.sh | sh'", nil, func(line string) {
			SendLog(stepID, line)
		})
	default:
		SendLog(stepID, "Unknown OS, attempting curl install script...")
		result = system.RunWithLogs("bash -c 'curl -f https://zed.dev/install.sh | sh'", nil, func(line string) {
			SendLog(stepID, line)
		})
	}
	if result != nil && result.Error != nil {
		SendLog(stepID, "Warning: Zed install failed: "+result.Error.Error())
		SendLog(stepID, "You can install Zed manually from https://zed.dev/download")
		// Non-fatal: continue with config copy in case user has Zed already
	} else {
		SendLog(stepID, "Zed binary installed")
	}

	// Copy config
	SendLog(stepID, "Copying Zed configuration...")
	zedDir := filepath.Join(homeDir, ".config", "zed")
	if err := system.EnsureDir(zedDir); err != nil {
		return wrapStepError("zed", "Install Zed",
			"Failed to create Zed config directory",
			err)
	}

	srcZed := filepath.Join(repoDir, "GentlemanZed", "zed")
	if err := system.CopyDir(srcZed, zedDir); err != nil {
		return wrapStepError("zed", "Install Zed",
			"Failed to copy Zed configuration",
			err)
	}

	SendLog(stepID, "Zed configured with Gentleman setup")
	return nil
}
```

---

### 5. `installer/cmd/gentleman-installer/main.go`

#### 5a. Add `zed` field to `cliFlags` struct (line 27, after `nvim`)

Current code at lines 26-27:
```go
	nvim           bool
	font           bool
```

Add:
```go
	nvim           bool
	zed            bool
	font           bool
```

#### 5b. Add `--zed` flag to `parseFlags()` (line 57, after `--nvim`)

Current code at lines 57-58:
```go
	flag.BoolVar(&flags.nvim, "nvim", false, "Install Neovim configuration")
	flag.BoolVar(&flags.font, "font", false, "Install Nerd Font")
```

Add:
```go
	flag.BoolVar(&flags.nvim, "nvim", false, "Install Neovim configuration")
	flag.BoolVar(&flags.zed, "zed", false, "Install Zed editor with config")
	flag.BoolVar(&flags.font, "font", false, "Install Nerd Font")
```

#### 5c. Add `InstallZed` to choices in `runNonInteractive()` (line 323, after `InstallNvim`)

Current code at lines 319-324:
```go
	choices := tui.UserChoices{
		Terminal:              terminal,
		Shell:                 shell,
		WindowMgr:             wm,
		InstallNvim:           flags.nvim,
		InstallFont:           flags.font,
```

Add:
```go
	choices := tui.UserChoices{
		Terminal:              terminal,
		Shell:                 shell,
		WindowMgr:             wm,
		InstallNvim:           flags.nvim,
		InstallZed:            flags.zed,
		InstallFont:           flags.font,
```

#### 5d. Add Zed line to summary output (line 338, after Neovim)

Current code at lines 338-339:
```go
	fmt.Printf("  Neovim:      %v\n", choices.InstallNvim)
	fmt.Printf("  Font:        %v\n", choices.InstallFont)
```

Add:
```go
	fmt.Printf("  Neovim:      %v\n", choices.InstallNvim)
	fmt.Printf("  Zed:         %v\n", choices.InstallZed)
	fmt.Printf("  Font:        %v\n", choices.InstallFont)
```

#### 5e. Add `--zed` to `printHelp()` (line 408, after `--nvim`)

Current code at lines 408-409:
```go
  --nvim               Install Neovim configuration
  --font               Install Nerd Font
```

Add:
```go
  --nvim               Install Neovim configuration
  --zed                Install Zed editor with config
  --font               Install Nerd Font
```

---

### 6. `installer/internal/tui/non_interactive.go`

#### 6a. Add Zed step to `buildStepsForChoices()` (line 108, after nvim block)

Current code at lines 105-113:
```go
	// Neovim
	if m.Choices.InstallNvim {
		steps = append(steps, InstallStep{ID: "nvim", Name: "Install Neovim configuration"})
	}

	// AI Tools
	if len(m.Choices.AITools) > 0 {
```

Add after the nvim block:
```go
	// Zed editor
	if m.Choices.InstallZed {
		steps = append(steps, InstallStep{ID: "zed", Name: "Install Zed editor"})
	}
```

---

### 7. `installer/internal/system/exec.go`

#### 7a. Add `"zed"` to `ConfigPaths()` map (line 348, after ghostty)

Current code at lines 347-349:
```go
		"ghostty":   home + "/.config/ghostty",
		"starship":  home + "/.config/starship.toml",
	}
```

Add:
```go
		"ghostty":   home + "/.config/ghostty",
		"zed":       home + "/.config/zed",
		"starship":  home + "/.config/starship.toml",
	}
```

---

## New Files

### `GentlemanZed/zed/settings.json`

```json
{
  "vim_mode": true,
  "theme": {
    "mode": "dark",
    "dark": "Kanagawa Wave",
    "light": "Kanagawa Lotus"
  },
  "ui_font_size": 14,
  "buffer_font_family": "Iosevka Term",
  "buffer_font_size": 14,
  "relative_line_numbers": true,
  "tab_size": 2,
  "hard_tabs": false,
  "format_on_save": "on",
  "autosave": {
    "after_delay": {
      "milliseconds": 1000
    }
  },
  "cursor_blink": false,
  "scrollbar": {
    "show": "never"
  },
  "vertical_scroll_margin": 8,
  "gutter": {
    "line_numbers": true,
    "code_actions": true,
    "folds": true
  },
  "indent_guides": {
    "enabled": true,
    "coloring": "indent_aware"
  },
  "inlay_hints": {
    "enabled": true,
    "show_type_hints": true,
    "show_parameter_hints": true
  },
  "terminal": {
    "font_family": "Iosevka Term",
    "font_size": 14,
    "shell": "system"
  },
  "telemetry": {
    "diagnostics": false,
    "metrics": false
  },
  "features": {
    "copilot": false
  },
  "assistant": {
    "default_model": {
      "provider": "anthropic",
      "model": "claude-sonnet-4-20250514"
    },
    "version": "2"
  },
  "language_models": {
    "anthropic": {
      "version": "1"
    }
  }
}
```

### `GentlemanZed/zed/keymap.json`

```json
[
  {
    "context": "Workspace",
    "bindings": {
      "ctrl-h": "workspace::ActivatePaneLeft",
      "ctrl-l": "workspace::ActivatePaneRight",
      "ctrl-k": "workspace::ActivatePaneUp",
      "ctrl-j": "workspace::ActivatePaneDown"
    }
  },
  {
    "context": "Editor && vim_mode == normal",
    "bindings": {
      "space f f": "file_finder::Toggle",
      "space f g": "workspace::NewSearch",
      "space e": "workspace::ToggleLeftDock",
      "space b d": "pane::CloseActiveItem",
      "space b n": "pane::ActivateNextItem",
      "space b p": "pane::ActivatePrevItem",
      "space /": "editor::ToggleComments",
      "space w": "workspace::Save",
      "space q": "pane::CloseActiveItem",
      "g d": "editor::GoToDefinition",
      "g r": "editor::FindAllReferences",
      "g i": "editor::GoToImplementation",
      "K": "editor::Hover",
      "space c a": "editor::ToggleCodeActions",
      "space r n": "editor::Rename",
      "[ d": "editor::GoToPrevDiagnostic",
      "] d": "editor::GoToNextDiagnostic",
      "space s v": "pane::SplitRight",
      "space s h": "pane::SplitDown"
    }
  },
  {
    "context": "Editor && vim_mode == visual",
    "bindings": {
      "space /": "editor::ToggleComments",
      "J": "editor::MoveLineDown",
      "K": "editor::MoveLineUp"
    }
  }
]
```

---

## Test Updates

### Files requiring changes:

1. **`installer/internal/tui/model_test.go`**
   - `TestGetCurrentOptions`: Add test case for `ScreenZedSelect` — verify it returns 2 options
   - No existing assertions break from iota reorder (tests use symbolic constants)

2. **`installer/internal/tui/update_test.go`**
   - Add test: `ScreenNvimSelect` enter navigates to `ScreenZedSelect` (not `ScreenAIToolsSelect`)
   - Add test: `ScreenZedSelect` enter with cursor=0 sets `InstallZed=true`, advances to `ScreenAIToolsSelect`
   - Add test: `ScreenZedSelect` enter with cursor=1 sets `InstallZed=false`, advances to `ScreenAIToolsSelect`
   - Add test: `ScreenZedSelect` escape goes back to `ScreenNvimSelect`
   - Add test: `ScreenAIToolsSelect` escape goes back to `ScreenZedSelect`
   - Update any existing test that asserts `ScreenNvimSelect` -> `ScreenAIToolsSelect` transition

3. **`installer/internal/tui/comprehensive_test.go`**
   - Update any wizard flow tests that assert screen sequence
   - Verify progress bar test includes "Zed" in the steps array

4. **`installer/internal/tui/installation_steps_test.go`**
   - Add test: `SetupInstallSteps` includes `"zed"` step when `InstallZed=true`
   - Add test: `SetupInstallSteps` excludes `"zed"` step when `InstallZed=false`

5. **`installer/internal/tui/integration_test.go`**
   - Update full-flow integration tests to include `ScreenZedSelect` in the expected sequence
