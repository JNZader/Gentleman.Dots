# Tasks: Project Init & Skill Manager

## Overview

This change adds two new top-level TUI flows to the Javi.Dots installer: "Initialize Project" (a multi-step wizard collecting path, stack, memory, CI choices and executing `init-project.sh`) and "Skill Manager" (browse/install/remove skills via `add-skill.sh`). Both follow the existing Bubbletea Model-Update-View pattern with drill-down navigation. Implementation is ordered so each task builds on the previous: model layer first, then handlers, then views, then execution, then CLI, then tests.

## Task List

### Task 1: Add screen constants, new types, and model fields

- **File(s)**: `installer/internal/tui/model.go`
- **Description**: Append 13 new screen constants to the existing `Screen` iota block (after `ScreenTrainerBossResult`): `ScreenProjectPath`, `ScreenProjectStack`, `ScreenProjectMemory`, `ScreenProjectEngram`, `ScreenProjectCI`, `ScreenProjectConfirm`, `ScreenProjectInstalling`, `ScreenProjectResult`, `ScreenSkillMenu`, `ScreenSkillBrowse`, `ScreenSkillInstall`, `ScreenSkillRemove`, `ScreenSkillResult`. Add new fields to the `Model` struct for project init state (`ProjectPathInput string`, `ProjectPathError string`, `ProjectStack string`, `ProjectMemory string`, `ProjectEngram bool`, `ProjectCI string`, `ProjectLogLines []string`) and skill manager state (`SkillList []string`, `InstalledSkills []string`, `SkillSelected []bool`, `SkillScroll int`, `SkillLoading bool`, `SkillLoadError string`, `SkillResultLog []string`). Add new fields to `UserChoices` for non-interactive passthrough (`InitProject bool`, `ProjectPath string`, `ProjectStack string`, `ProjectMemory string`, `ProjectCI string`, `ProjectEngram bool`, `SkillAction string`, `SkillNames []string`). Initialize all new fields in `NewModel()`.
- **Depends on**: none
- **Acceptance criteria**:
  - [ ] 13 new `Screen*` constants compile and follow iota ordering after `ScreenTrainerBossResult`
  - [ ] `Model` struct has all 14 new fields (7 project + 7 skill) with correct types
  - [ ] `UserChoices` struct has all 8 new fields for non-interactive mode
  - [ ] `NewModel()` initializes all new fields to zero/empty values
  - [ ] `go build ./...` compiles without errors

### Task 2: Wire GetCurrentOptions, GetScreenTitle, and GetScreenDescription for all new screens

- **File(s)**: `installer/internal/tui/model.go`
- **Description**: Add switch cases in `GetCurrentOptions()` for `ScreenProjectStack` (stack list with auto-detect label), `ScreenProjectMemory` (5 memory options with emoji+description), `ScreenProjectEngram` (Yes/No), `ScreenProjectCI` (4 CI options), `ScreenProjectConfirm` (Confirm & Run / Cancel), `ScreenSkillMenu` (Browse/Install/Remove/separator/Back), `ScreenSkillInstall` (SkillList + separator + Confirm), `ScreenSkillRemove` (InstalledSkills + separator + Confirm). Add the two new main menu items `"📦 Initialize Project"` and `"🎯 Skill Manager"` before `"❌ Exit"` in the `ScreenMainMenu` case. Add all 13 new screen titles to `GetScreenTitle()` and all applicable screen descriptions to `GetScreenDescription()`. The stack screen description should show "Auto-detected: X" when `ProjectStack` is set and not "unknown".
- **Depends on**: Task 1
- **Acceptance criteria**:
  - [ ] Main menu shows "Initialize Project" and "Skill Manager" before "Exit"
  - [ ] Each new screen returns non-empty options (where applicable), title, and description
  - [ ] `ScreenProjectStack` description dynamically includes auto-detected stack name
  - [ ] `ScreenSkillInstall` options length = `len(SkillList) + 2` (separator + confirm)
  - [ ] `ScreenSkillRemove` options length = `len(InstalledSkills) + 2`

### Task 3: Add helper functions (expandPath, detectStack) and new message types

- **File(s)**: `installer/internal/tui/update.go`
- **Description**: Add the `expandPath(p string) string` helper that expands leading `~/` to `os.UserHomeDir()`. Add `detectStack(path string) string` that checks for indicator files (`package.json`->node, `angular.json`->angular, `go.mod`->go, `Cargo.toml`->rust, `pom.xml`->java, `pyproject.toml`->python, etc.) using `filepath.Glob`. Add new message types: `projectInstallStartMsg struct{}`, `projectInstallLogMsg struct{ line string }`, `projectInstallCompleteMsg struct{ err error }`, `skillsLoadedMsg struct{ available []string; installed []string; err error }`, `skillActionCompleteMsg struct{ logLines []string; err error }`. These are pure data/utility additions with no behavioral changes yet.
- **Depends on**: Task 1
- **Acceptance criteria**:
  - [ ] `expandPath("~/foo")` returns `$HOME/foo`; `expandPath("/abs/path")` returns unchanged
  - [ ] `detectStack` returns correct stack for directories containing known indicator files
  - [ ] `detectStack` returns `"unknown"` when no indicator files are found
  - [ ] All 5 message types are declared and compile correctly
  - [ ] No existing tests break

### Task 4: Implement project init screen handlers and key dispatch

- **File(s)**: `installer/internal/tui/update.go`
- **Description**: Implement all project init handler functions: `handleProjectPathKeys` (text input via manual key accumulation -- same pattern as `TrainerInput` -- with backspace, enter for path validation using `expandPath` + `os.Stat` + `filepath.Abs`, esc to return to main menu), `handleProjectSelectionKeys` (single dispatch for `ScreenProjectStack`, `ScreenProjectMemory`, `ScreenProjectEngram`, `ScreenProjectCI`, `ScreenProjectConfirm` -- follows `handleSelectionKeys` pattern with per-screen enter/esc logic from the design doc), `handleProjectResultKeys` (enter/esc returns to main menu). Add dispatch cases in `handleKeyPress` switch for all new project screens. Add `ScreenProjectPath` to the space-key exclusion list so space appends to input instead of activating leader mode. Handle `projectInstallStartMsg`, `projectInstallLogMsg`, and `projectInstallCompleteMsg` in the main `Update()` message switch. Wire main menu dispatch for "Initialize Project" in `handleMainMenuKeys`. The `ScreenProjectInstalling` screen has no key handler (keys ignored, same as existing `ScreenInstalling`). Conditional routing: memory="obsidian-brain" -> Engram screen, otherwise skip to CI; CI esc goes back to Engram if obsidian-brain was selected, otherwise back to Memory.
- **Depends on**: Tasks 1, 2, 3
- **Acceptance criteria**:
  - [ ] Selecting "Initialize Project" from main menu navigates to `ScreenProjectPath` with cleared state
  - [ ] Typing characters appends to `ProjectPathInput`; backspace removes last rune
  - [ ] Enter on invalid path shows inline error without advancing; valid dir advances to `ScreenProjectStack`
  - [ ] Stack/Memory/Engram/CI/Confirm screens navigate forward and backward correctly per design flow diagram
  - [ ] Obsidian Brain memory selection routes to Engram screen; other selections skip to CI
  - [ ] Confirm screen cursor=0 triggers `projectInstallStartMsg`; cursor=1 cancels to main menu
  - [ ] `projectInstallLogMsg` appends to `ProjectLogLines` (capped at 30)
  - [ ] `projectInstallCompleteMsg` transitions to `ScreenProjectResult`
  - [ ] Esc from each screen returns to the correct predecessor screen

### Task 5: Implement skill manager screen handlers and key dispatch

- **File(s)**: `installer/internal/tui/update.go`
- **Description**: Implement all skill manager handler functions: `handleSkillMenuKeys` (Browse/Install/Remove dispatch with async `loadSkillsCmd()` trigger for Browse and Install; `loadInstalledSkillsCmd()` for Remove), `handleSkillBrowseKeys` (read-only scroll with j/k/up/down, esc back to menu), `handleSkillInstallKeys` (multi-select toggle with scroll, enter on skill item toggles `SkillSelected`, enter on Confirm collects selected names and triggers `runSkillInstallCmd`), `handleSkillRemoveKeys` (identical to install but uses `InstalledSkills` and `runSkillRemoveCmd`), `handleSkillResultKeys` (enter/esc returns to skill menu). Add `loadSkillsCmd() tea.Cmd` and `fetchSkillCatalog()` (clones Gentleman-Skills repo to `/tmp/gentleman-skills-cache` if stale or missing, reads skill directories, reads installed skills from `$CWD/.claude/skills/`). Handle `skillsLoadedMsg` and `skillActionCompleteMsg` in `Update()`. Wire main menu dispatch for "Skill Manager" in `handleMainMenuKeys`. Add dispatch cases in `handleKeyPress` for all skill screens.
- **Depends on**: Tasks 1, 2, 3
- **Acceptance criteria**:
  - [ ] Selecting "Skill Manager" from main menu navigates to `ScreenSkillMenu`
  - [ ] Browse/Install/Remove dispatch from skill menu sets `SkillLoading = true` and returns appropriate `tea.Cmd`
  - [ ] `skillsLoadedMsg` populates `SkillList` and `InstalledSkills`, clears `SkillLoading`
  - [ ] `skillsLoadedMsg` with error sets `SkillLoadError`
  - [ ] Browse screen scrolls with j/k, esc returns to menu and resets scroll
  - [ ] Install screen toggles `SkillSelected[cursor]` on enter; confirm with zero selected is no-op
  - [ ] Install confirm with selected skills transitions to `ScreenSkillResult` and triggers async install
  - [ ] Remove screen works identically but on `InstalledSkills`
  - [ ] Result screen enter/esc returns to `ScreenSkillMenu`
  - [ ] `fetchSkillCatalog` reuses cache if less than 1 hour old

### Task 6: Implement view/render functions for all new screens

- **File(s)**: `installer/internal/tui/view.go`
- **Description**: Add render functions and wire them in the main `View()` switch. `renderProjectPath()`: centered layout with title, input line (`> input█`), inline error if `ProjectPathError != ""`, and help text. `ScreenProjectStack`, `ScreenProjectMemory`, `ScreenProjectEngram`, `ScreenProjectCI`: reuse existing `renderSelection()` (already handles title/description/options from `GetCurrentOptions`/`GetScreenTitle`/`GetScreenDescription`). `renderProjectConfirm()`: summary block showing path/stack/memory/engram/CI, then action buttons with cursor. `renderProjectInstalling()`: spinner + log lines from `ProjectLogLines`. `renderProjectResult()`: success or error based on `ErrorMsg`. `renderSkillMenu()`: simple option list (reuse `renderSelection`). `renderSkillBrowse()`: scrollable read-only list with `SkillScroll`, loading state, and error state. `renderSkillInstall()` / `renderSkillRemove()`: multi-select with checkboxes and scroll, loading/error states. `renderSkillResult()`: result log lines with success/error per skill. Add all 13 new cases to the `View()` switch statement.
- **Depends on**: Tasks 1, 2
- **Acceptance criteria**:
  - [ ] `renderProjectPath` shows input with cursor block and inline error when set
  - [ ] Project selection screens (stack/memory/engram/CI) render via `renderSelection` pattern
  - [ ] `renderProjectConfirm` displays all choice fields as summary
  - [ ] `renderProjectInstalling` shows spinner animation and log lines
  - [ ] `renderProjectResult` renders success or error conditionally
  - [ ] `renderSkillBrowse` shows "Fetching..." when `SkillLoading`, error when `SkillLoadError`, list with scroll indicators otherwise
  - [ ] `renderSkillInstall`/`renderSkillRemove` show checkboxes `[✓]`/`[ ]` for toggled items
  - [ ] `renderSkillResult` shows per-skill operation log
  - [ ] All 13 new screen constants have a case in `View()` -- no missing screens

### Task 7: Implement installer step functions (stepInitProject, stepSkillInstall, stepSkillRemove)

- **File(s)**: `installer/internal/tui/installer.go`
- **Description**: Add three new installer functions. `stepInitProject(m *Model) error`: checks for cached clone at `/tmp/project-starter-framework-install` (reuse if < 1 hour old, else re-clone with `git clone --depth 1`), makes `init-project.sh` executable, builds command string with `--non-interactive --memory=X --ci=X [--engram]` flags, executes via `system.RunWithLogs` with `ExecOptions{WorkDir: m.Choices.ProjectPath}` to run in the target project directory, streams output via `SendLog`. `stepSkillInstall(m *Model) error`: ensures framework repo clone exists (same cache logic), runs `bash add-skill.sh gentleman <name>` for each skill in `m.Choices.SkillNames`, collects per-skill success/failure. `stepSkillRemove(m *Model) error`: same pattern but runs `add-skill.sh remove <name>`. Add `runProjectInitCmd(m *Model) tea.Cmd` and `runSkillInstallCmd`/`runSkillRemoveCmd` wrappers that capture model state and return async commands. Note: use `system.RunWithLogs` with `ExecOptions{WorkDir: path}` since that field already exists on `ExecOptions` -- no need for a new `RunWithLogsInDir` function.
- **Depends on**: Tasks 1, 3
- **Acceptance criteria**:
  - [ ] `stepInitProject` clones framework repo if not cached or stale
  - [ ] `stepInitProject` constructs correct `init-project.sh` command with all flags
  - [ ] `stepInitProject` runs command in the target project directory via `WorkDir`
  - [ ] `stepSkillInstall` iterates selected skills and runs `add-skill.sh gentleman <name>` for each
  - [ ] `stepSkillRemove` runs `add-skill.sh remove <name>` for each selected skill
  - [ ] Failed skills are collected and returned as a combined error without stopping other installations
  - [ ] `runProjectInitCmd` returns a `tea.Cmd` that sends `projectInstallCompleteMsg`
  - [ ] All functions use `SendLog` for streaming output to the TUI

### Task 8: Add non-interactive CLI flags for project init and skill management

- **File(s)**: `installer/cmd/gentleman-installer/main.go`
- **Description**: Add new fields to `cliFlags` struct: `initProject bool`, `projectPath string`, `projectStack string`, `projectMemory string` (default "simple"), `projectCI string` (default "none"), `projectEngram bool`, `skillAction string`, `skillNames string`. Register them in `parseFlags()` with `flag.BoolVar`/`flag.StringVar`. In `runNonInteractive()`, add validation: if `--init-project` is set, require `--project-path` (validate with `expandPath` + `os.Stat`), validate `--project-memory` and `--project-ci` against allowed values, populate `UserChoices` project fields. If `--skill` is set, validate action is "install" or "remove", require `--skill-names`, parse comma-separated names. Add dispatch in `RunNonInteractive` (or `runNonInteractive`) to call `stepInitProject` and/or skill actions when those flags are present. Update `printHelp()` with new flag documentation and examples.
- **Depends on**: Tasks 1, 7
- **Acceptance criteria**:
  - [ ] `--init-project --project-path=/valid/dir` runs `stepInitProject` non-interactively
  - [ ] `--init-project` without `--project-path` exits with error message
  - [ ] `--project-memory=invalid` exits with validation error
  - [ ] `--project-ci=invalid` exits with validation error
  - [ ] `--skill=install --skill-names=react-19,typescript` runs skill install for each name
  - [ ] `--skill=remove --skill-names=react-19` runs skill remove
  - [ ] `--skill=invalid` exits with validation error
  - [ ] `--skill=install` without `--skill-names` exits with error
  - [ ] `--help` output includes all new flags with descriptions
  - [ ] Exit codes mirror the underlying script exit codes

### Task 9: Unit tests for project init flow

- **File(s)**: `installer/internal/tui/project_screens_test.go` (new file)
- **Description**: Write unit tests following the existing pattern in `update_test.go` (create model with `NewModel()`, set fields, send `tea.KeyMsg`, assert resulting screen and state). Tests to include: `TestProjectPathValidation` (empty input error, non-existent path error, file-not-dir error, valid dir advances), `TestProjectPathTildeExpansion` (verify `expandPath` expands `~` correctly), `TestProjectPathBackspace` (verify rune-safe backspace), `TestProjectStackNavigation` (cursor movement, enter saves stack), `TestProjectMemoryConditionalEngram` (obsidian-brain routes to Engram screen, others skip to CI), `TestProjectCIAdvancesToConfirm`, `TestProjectEscapeBackNavigation` (esc from each screen returns to correct predecessor), `TestProjectConfirmCancel` (cursor=1 returns to main menu), `TestGetCurrentOptionsProjectScreens` (correct option counts and labels), `TestGetScreenTitleProjectScreens` (non-empty titles), `TestGetScreenDescriptionProjectStack` (auto-detected prefix when stack is set).
- **Depends on**: Tasks 1, 2, 3, 4
- **Acceptance criteria**:
  - [ ] All tests pass with `go test ./installer/internal/tui/ -run TestProject`
  - [ ] Path validation tests cover: empty, non-existent, file (not dir), valid directory
  - [ ] Conditional Engram routing is tested for both branches
  - [ ] Escape back-navigation is tested for every project screen
  - [ ] `GetCurrentOptions`/`GetScreenTitle`/`GetScreenDescription` coverage for all project screens
  - [ ] No existing tests are broken

### Task 10: Unit tests for skill manager flow

- **File(s)**: `installer/internal/tui/skill_screens_test.go` (new file)
- **Description**: Write unit tests following the same pattern. Tests to include: `TestSkillMenuNavigation` (Browse/Install/Remove options exist and dispatch to correct screens), `TestSkillMenuEscape` (returns to main menu), `TestSkillBrowseScrolling` (j/k adjusts `SkillScroll`, esc resets and returns to menu), `TestSkillInstallToggle` (enter toggles `SkillSelected[cursor]`), `TestSkillInstallConfirmNoSelection` (confirm with nothing selected is no-op), `TestSkillInstallConfirm` (with selections transitions to result screen), `TestSkillRemoveEmpty` (empty `InstalledSkills` shows no options, no panic), `TestSkillsLoadedMsg` (populates `SkillList`, clears `SkillLoading`, initializes `SkillSelected`), `TestSkillsLoadedMsgError` (sets `SkillLoadError`, keeps list empty), `TestSkillResultEnterReturnsToMenu`, `TestGetCurrentOptionsSkillScreens` (correct option lengths for install/remove), `TestRenderSkillInstallLoading` (loading state renders without panic on empty list).
- **Depends on**: Tasks 1, 2, 3, 5
- **Acceptance criteria**:
  - [ ] All tests pass with `go test ./installer/internal/tui/ -run TestSkill`
  - [ ] Skill menu navigation tested for all 3 sub-actions
  - [ ] Multi-select toggle tested for both select and deselect
  - [ ] Empty installed skills list renders safely (no index out of bounds)
  - [ ] `skillsLoadedMsg` handling tested for both success and error paths
  - [ ] Loading state render does not panic on empty `SkillList`
  - [ ] No existing tests are broken

### Task 11: Integration verification and full build

- **File(s)**: All modified files
- **Description**: Run `go build ./...` and `go test ./...` across the entire project. Verify that: (1) all new screens are reachable via navigation from main menu, (2) the project init flow chains correctly through all 8 screens with proper conditional branching at the Engram step, (3) the skill manager flow chains correctly through all sub-actions, (4) escape navigation from every new screen leads to the correct predecessor, (5) non-interactive CLI flags parse and validate correctly, (6) no existing tests are broken by the additions. Fix any compilation errors, missing imports, or test failures discovered. Ensure the `View()` switch has no missing cases for the 13 new screen constants.
- **Depends on**: Tasks 1-10
- **Acceptance criteria**:
  - [ ] `go build ./...` completes with zero errors
  - [ ] `go test ./...` passes with zero failures
  - [ ] `go vet ./...` reports no issues
  - [ ] Every new `Screen*` constant has a handler case in `handleKeyPress` (or is explicitly ignored like `ScreenProjectInstalling`)
  - [ ] Every new `Screen*` constant has a case in `View()`
  - [ ] Every new `Screen*` constant has an entry in `GetScreenTitle()`
  - [ ] Main menu shows both new items in correct position
  - [ ] Non-interactive `--init-project` and `--skill` flags are functional
