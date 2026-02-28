# Tasks: zed-editor-support

## Pre-flight

- [ ] **TASK-01**: Verify Zed install commands work on test platforms
  - **Description**: Manually verify that `brew install --cask zed` (macOS), `pacman -S zed` (Arch), and `curl -f https://zed.dev/install.sh | sh` (Debian/Fedora) work. Confirm the Zed settings.json schema is stable and accepts the keys we use (`vim_mode`, `theme`, `buffer_font_family`, `relative_line_numbers`, `assistant`).
  - **Files**: None (manual verification)
  - **Estimated lines**: 0
  - **Dependencies**: None

## Implementation

- [ ] **TASK-02**: Add `ScreenZedSelect` constant and `UserChoices.InstallZed` field
  - **Description**: In `model.go`, insert `ScreenZedSelect` at line 26 (between `ScreenNvimSelect` and `ScreenInstalling` in the `Screen` iota). Add `InstallZed bool` field to `UserChoices` struct at line 122 (after `InstallNvim`).
  - **Files**: `installer/internal/tui/model.go`
  - **Estimated lines**: +2 lines
  - **Dependencies**: None

- [ ] **TASK-03**: Add screen options, title, and description for `ScreenZedSelect`
  - **Description**: In `model.go`, add three switch cases:
    1. `GetCurrentOptions()` (~line 399): return `["Yes, install Zed with config", "No, skip Zed"]`
    2. `GetScreenTitle()` (~line 566): return `"Step 7: Zed Editor"` and renumber AI screens to Step 8/9
    3. `GetScreenDescription()` (~line 713): return `"High-performance editor with Vim mode and AI agent support"`
  - **Files**: `installer/internal/tui/model.go`
  - **Estimated lines**: +8 lines (net, after renumbering edits)
  - **Dependencies**: TASK-02

- [ ] **TASK-04**: Add navigation (forward, backward, escape) for `ScreenZedSelect`
  - **Description**: In `update.go`:
    1. Add `ScreenZedSelect` to the `handleSelectionKeys` dispatch list (line 786)
    2. Add `ScreenZedSelect` to the escape handler list (line 911)
    3. Change `ScreenNvimSelect` forward navigation (lines 1356-1365): go to `ScreenZedSelect` instead of `ScreenAIToolsSelect`
    4. Add `ScreenZedSelect` forward handler: set `InstallZed`, advance to `ScreenAIToolsSelect` with `AIToolSelected` init
    5. Change `ScreenAIToolsSelect` backward navigation (line 1169-1173): go back to `ScreenZedSelect` instead of `ScreenNvimSelect`
    6. Add `ScreenZedSelect` backward handler: go back to `ScreenNvimSelect`, reset `InstallZed = false`
  - **Files**: `installer/internal/tui/update.go`
  - **Estimated lines**: +12 lines
  - **Dependencies**: TASK-02

- [ ] **TASK-05**: Add `ScreenZedSelect` to view render and update progress bar
  - **Description**: In `view.go`:
    1. Add `ScreenZedSelect` to the `renderSelection()` dispatch at line 54
    2. Update `renderStepProgress()`: add `"Zed"` to steps array at index 6, add `ScreenZedSelect` case mapping to index 6, shift `ScreenAIToolsSelect` to 7 and `ScreenAIFramework*` to 8
  - **Files**: `installer/internal/tui/view.go`
  - **Estimated lines**: +4 lines
  - **Dependencies**: TASK-02

- [ ] **TASK-06**: Add Zed install step registration in `SetupInstallSteps()` and `executeStep()`
  - **Description**: In `installer.go`:
    1. Add `"zed"` case to `executeStep()` dispatch (after `"nvim"` case, line 69): call `stepInstallZed(m)`
    2. Add Zed step block in `SetupInstallSteps()` after the nvim block (after line 1053): register `InstallStep{ID: "zed", Name: "Install Zed", Description: "Editor with Vim mode"}`
  - **Files**: `installer/internal/tui/installer.go`
  - **Estimated lines**: +12 lines
  - **Dependencies**: TASK-02

- [ ] **TASK-07**: Implement `stepInstallZed()` function
  - **Description**: In `installer.go`, add the `stepInstallZed(m *Model) error` function after `stepInstallNvim()` (after line 1080). The function must:
    1. Early-return on Termux with a skip log
    2. Install Zed binary via platform-specific commands (brew cask for macOS, pacman for Arch, curl install script for Debian/Fedora/other)
    3. Non-fatal on install failure (log warning + manual install URL)
    4. Create `~/.config/zed/` with `system.EnsureDir()`
    5. Copy config from `GentlemanZed/zed/` to `~/.config/zed/` with `system.CopyDir()`
    6. Fatal on config copy failure (return `wrapStepError`)
  - **Files**: `installer/internal/tui/installer.go`
  - **Estimated lines**: +55 lines
  - **Dependencies**: TASK-06, TASK-08

- [ ] **TASK-08**: Create `GentlemanZed/zed/` config files
  - **Description**: Create two new files in the repository:
    1. `GentlemanZed/zed/settings.json` — Zed settings with vim_mode, Kanagawa theme, Iosevka font, relative line numbers, Claude assistant config, telemetry disabled
    2. `GentlemanZed/zed/keymap.json` — Vim-mode keybindings: window navigation (ctrl-h/j/k/l), leader bindings (space-based: file finder, search, sidebar, buffer ops, LSP actions), visual mode helpers
  - **Files**: `GentlemanZed/zed/settings.json` (NEW), `GentlemanZed/zed/keymap.json` (NEW)
  - **Estimated lines**: ~120 lines total
  - **Dependencies**: None

- [ ] **TASK-09**: Add `--zed` CLI flag
  - **Description**: In `main.go`:
    1. Add `zed bool` field to `cliFlags` struct (line 27)
    2. Add `flag.BoolVar(&flags.zed, "zed", false, ...)` in `parseFlags()` (line 58)
    3. Add `InstallZed: flags.zed` to choices in `runNonInteractive()` (line 324)
    4. Add `Zed:` line to summary output (line 339)
    5. Add `--zed` to `printHelp()` (line 409)
  - **Files**: `installer/cmd/gentleman-installer/main.go`
  - **Estimated lines**: +5 lines
  - **Dependencies**: TASK-02

- [ ] **TASK-10**: Add non-interactive support for Zed
  - **Description**: In `non_interactive.go`, add a Zed step block in `buildStepsForChoices()` after the nvim block (after line 108): if `m.Choices.InstallZed` is true, append `InstallStep{ID: "zed", Name: "Install Zed editor"}`.
  - **Files**: `installer/internal/tui/non_interactive.go`
  - **Estimated lines**: +4 lines
  - **Dependencies**: TASK-02

- [ ] **TASK-11**: Add `ConfigPaths` entry for backup detection
  - **Description**: In `installer/internal/system/exec.go`, add `"zed": home + "/.config/zed"` to the `ConfigPaths()` map (line 348, after the `ghostty` entry). This enables the backup system to detect and backup existing Zed configs before overwrite.
  - **Files**: `installer/internal/system/exec.go`
  - **Estimated lines**: +1 line
  - **Dependencies**: None

- [ ] **TASK-12**: Update documentation
  - **Description**: Update any relevant examples in `printHelp()` to include `--zed` flag in example commands. Verify the help text shows `--zed` in the correct section.
  - **Files**: `installer/cmd/gentleman-installer/main.go`
  - **Estimated lines**: +2 lines (example additions)
  - **Dependencies**: TASK-09

## Verification

- [ ] **TASK-13**: Update existing tests
  - **Description**: Update test files to account for the new `ScreenZedSelect` screen:
    1. `model_test.go`: Add `ScreenZedSelect` options test (2 options expected)
    2. `update_test.go`: Fix any test asserting `ScreenNvimSelect` -> `ScreenAIToolsSelect` transition (now goes via `ScreenZedSelect`). Add tests for Zed forward/backward navigation.
    3. `comprehensive_test.go`: Update wizard flow sequence assertions to include `ScreenZedSelect`
    4. `installation_steps_test.go`: Add tests verifying `SetupInstallSteps` includes/excludes `"zed"` step based on `InstallZed` flag
    5. `integration_test.go`: Update full-flow tests
  - **Files**: `installer/internal/tui/model_test.go`, `installer/internal/tui/update_test.go`, `installer/internal/tui/comprehensive_test.go`, `installer/internal/tui/installation_steps_test.go`, `installer/internal/tui/integration_test.go`
  - **Estimated lines**: +40 lines (across files)
  - **Dependencies**: TASK-02 through TASK-10

- [ ] **TASK-14**: Run `go test ./...` and fix any failures
  - **Description**: Run the full test suite from `installer/` directory. Fix any failures caused by iota reordering, navigation changes, or missing switch cases. Ensure zero failures before marking complete.
  - **Files**: All test files
  - **Estimated lines**: Variable
  - **Dependencies**: TASK-13

- [ ] **TASK-15**: Manual testing on macOS and Linux
  - **Description**: Run the installer in TUI mode and verify:
    1. `ScreenZedSelect` appears as Step 7 after Neovim
    2. Progress bar shows 9 steps with "Zed" at position 6
    3. Selecting "Yes" leads to AI Tools; selecting "No" also leads to AI Tools
    4. Esc from Zed goes back to Nvim; Esc from AI Tools goes back to Zed
    5. Non-interactive mode with `--zed` triggers Zed installation
    6. Backup system detects existing `~/.config/zed/` directory
    7. Config files are correctly copied to `~/.config/zed/`
    8. Zed launches with vim_mode, Kanagawa theme, and correct font
  - **Files**: None (manual testing)
  - **Estimated lines**: 0
  - **Dependencies**: TASK-14
