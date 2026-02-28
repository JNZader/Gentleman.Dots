# Spec: Project Init & Skill Manager

## Requirements

### FR-1: Initialize Project Flow

- FR-1.1: Main Menu shows "📦 Initialize Project" item
  - Inserted between "🎮 Vim Trainer" and "🔄 Restore from Backup" (if present) / "❌ Exit"
  - Uses `strings.Contains(selected, "Initialize Project")` dispatch in `handleMainMenuKeys`

- FR-1.2: Text input screen for project directory path (`ScreenProjectPath`)
  - Uses Bubbletea's `textinput` component (charmbracelet/bubbles/textinput)
  - Placeholder: `/home/user/my-project`
  - On submit (Enter):
    - Expand leading `~` to `os.Getenv("HOME")`
    - Call `filepath.Abs()` to resolve relative paths
    - Validate with `os.Stat()`: path must exist and `IsDir()` must be true
    - If invalid: show inline error message below the input field; do NOT advance
    - If valid: proceed to `ScreenProjectStack`, store expanded absolute path in `Model.ProjectPath`
  - Esc: return to `ScreenMainMenu`, reset `ProjectPath`

- FR-1.3: Stack detection screen (`ScreenProjectStack`)
  - Auto-detects stack from project directory by checking for indicator files:
    - `build.gradle` or `build.gradle.kts` → `java-gradle`
    - `pom.xml` → `java-maven`
    - `package.json` → `node`
    - `requirements.txt`, `pyproject.toml`, or `setup.py` → `python`
    - `go.mod` → `go`
    - `Cargo.toml` → `rust`
    - None found → no pre-selection, user must choose
  - Detection runs as a `tea.Cmd` (non-blocking) when entering the screen; spinner shown while detecting
  - Screen shows single-select list of all supported stacks; auto-detected entry is highlighted with `(detected)` suffix
  - Supported options: `java-gradle`, `java-maven`, `node`, `python`, `go`, `rust`, `─────────────`, `← Back`
  - Enter confirms selection; store in `Model.ProjectStack`
  - Esc: return to `ScreenProjectPath`

- FR-1.4: Memory module selection screen (`ScreenProjectMemory`)
  - Single-select list with brief description shown as subtitle per item
  - Options (label → `--memory=N` value → description):
    1. `Obsidian Brain` → `--memory=1` → "Knowledge base via Obsidian vault"
    2. `VibeKanban` → `--memory=2` → "Kanban board task tracking"
    3. `Engram` → `--memory=3` → "Lightweight AI agent memory"
    4. `Simple` → `--memory=4` → "Plain markdown notes"
    5. `None` → `--memory=5` → "No memory module"
  - Enter confirms; store choice label in `Model.ProjectMemory` and numeric value in `Model.ProjectMemoryN`
  - If `Obsidian Brain` selected → advance to `ScreenProjectEngram`
  - Otherwise → advance to `ScreenProjectCI`
  - Esc: return to `ScreenProjectStack`

- FR-1.5: Conditional Engram prompt (`ScreenProjectEngram`)
  - Shown ONLY if `Obsidian Brain` was selected in FR-1.4
  - Yes/No single-select: "Also add Engram for AI agent memory?"
  - Options: `Yes, add Engram`, `No, skip Engram`
  - Store boolean in `Model.ProjectAddEngram`
  - Advance to `ScreenProjectCI`
  - Esc: return to `ScreenProjectMemory`

- FR-1.6: CI provider selection screen (`ScreenProjectCI`)
  - Single-select list
  - Options (label → `--ci=N` value):
    1. `GitHub Actions` → `--ci=1`
    2. `GitLab CI` → `--ci=2`
    3. `Woodpecker` → `--ci=3`
    4. `None` → `--ci=4`
  - Enter confirms; store in `Model.ProjectCI` and `Model.ProjectCIN`
  - Advance to `ScreenProjectConfirm`
  - Esc: return to `ScreenProjectEngram` (if Obsidian Brain was chosen) or `ScreenProjectMemory`

- FR-1.7: Confirmation screen (`ScreenProjectConfirm`)
  - Info screen (no cursor navigation) showing summary of all collected choices:
    - Project path
    - Stack
    - Memory module
    - Engram add-on (if applicable)
    - CI provider
    - The exact `init-project.sh` command that will be executed (derived from choices)
  - Two options: `✅ Run init-project.sh`, `❌ Cancel`
  - Enter on confirm → advance to `ScreenProjectInstalling`, trigger `stepInitProject()` as `tea.Cmd`
  - Enter on cancel / Esc → return to `ScreenMainMenu`, reset all `Project*` model fields

- FR-1.8: Execution screen (`ScreenProjectInstalling`)
  - Reuses existing `ScreenInstalling` rendering pattern (progress steps, log lines, spinner)
  - Single `InstallStep` with ID `"init-project"`
  - Step function `stepInitProject()` in `installer.go`:
    1. If `project-starter-framework` repo not cloned yet (check `/tmp/project-starter-framework/`), clone it: `git clone https://github.com/JNZader/project-starter-framework.git /tmp/project-starter-framework`
    2. Build command: `bash /tmp/project-starter-framework/init-project.sh --non-interactive --memory=N --ci=N [--engram]`
    3. Execute in the selected project directory (`Dir` field of `exec.Cmd`)
    4. Stream stdout/stderr via `SendLog("init-project", line)`
    5. On exit code 0 → send `stepCompleteMsg{stepID: "init-project"}`; on error → send with `err` set
  - No Esc available while running (same as existing `ScreenInstalling`)

- FR-1.9: Completion screen (`ScreenProjectComplete`)
  - Info screen showing:
    - Success message with project path
    - Summary list of what was created (parsed from script stdout, or static based on memory+CI choices if parsing is not feasible)
    - Hint: "Press Enter to return to Main Menu"
  - Enter / Esc → `ScreenMainMenu`, reset all `Project*` model fields and clear `Steps`

---

### FR-2: Skill Manager Flow

- FR-2.1: Main Menu shows "🎯 Skill Manager" item
  - Inserted after "📦 Initialize Project" in the main menu options list
  - Uses `strings.Contains(selected, "Skill Manager")` dispatch in `handleMainMenuKeys`

- FR-2.2: Skill sub-menu screen (`ScreenSkillMenu`)
  - Single-select list with 3 actions plus back:
    - `📖 Browse Skills`
    - `⬇️  Install Skill`
    - `🗑️  Remove Skill`
    - `─────────────`
    - `← Back`
  - Esc / "← Back" → `ScreenMainMenu`

- FR-2.3: Browse screen (`ScreenSkillBrowse`)
  - Scrollable read-only list of all skills available in the Gentleman-Skills repo
  - Skills are grouped by source tag: `curated` (first) and `community` (second)
  - Group headers rendered as non-selectable separators (same visual pattern as `─────────────`)
  - Each skill entry shows: skill name + one-line description (read from skill manifest or README first line)
  - Uses viewport scrolling (same pattern as `CategoryItemsScroll` / `LazyVimScroll`):
    - Model field: `SkillBrowseScroll int`
    - j/↓ scrolls down, k/↑ scrolls up, g goes to top, G goes to bottom
  - Loading state: spinner shown while Gentleman-Skills repo is cloned/cached (see FR-2.7)
  - Esc → `ScreenSkillMenu`

- FR-2.4: Install screen (`ScreenSkillInstall`)
  - Multi-select list (same toggle pattern as `ScreenAIToolsSelect` / `AIToolSelected []bool`):
    - Model field: `SkillInstallSelected []bool`
  - Lists all available skills from cached Gentleman-Skills repo (same source as Browse)
  - Supports viewport scrolling for long lists: model field `SkillInstallScroll int`
  - j/k/↑/↓ move cursor, Enter toggles selection, Space toggles selection
  - Bottom option: `✅ Confirm installation` (always visible, not scrolled away)
  - On confirm with at least one skill selected → advance to `ScreenSkillResult`, trigger `stepSkillInstall()`
  - On confirm with zero selected → show inline message "Select at least one skill"
  - Esc → `ScreenSkillMenu`

- FR-2.5: Remove screen (`ScreenSkillRemove`)
  - Multi-select list of currently installed skills
  - Installed skills detected by scanning `~/.claude/skills/` directory for subdirectories (or equivalent per-tool skill paths if applicable)
  - Model field: `InstalledSkills []string`, `SkillRemoveSelected []bool`
  - If no skills installed → show "No skills installed" message with single option `← Back`; Esc/Enter → `ScreenSkillMenu`
  - Same scrolling and toggle pattern as FR-2.4
  - On confirm → advance to `ScreenSkillResult`, trigger `stepSkillRemove()`
  - Esc → `ScreenSkillMenu`

- FR-2.6: Result screen (`ScreenSkillResult`)
  - Info screen showing success/error for each skill operation, one line per skill:
    - `✅ react-19 — installed successfully`
    - `❌ typescript — error: <message>`
  - Error details truncated to one line; full log available in `LogLines`
  - Options: `← Back to Skill Manager`
  - Enter / Esc → `ScreenSkillMenu`, clear `SkillInstallSelected` / `SkillRemoveSelected`

- FR-2.7: Loading state for Gentleman-Skills clone
  - When entering `ScreenSkillBrowse`, `ScreenSkillInstall`, or `ScreenSkillRemove`, a `tea.Cmd` checks for cached clone:
    - Cache path: `/tmp/gentleman-skills/`
    - Cache valid if clone exists AND `os.Stat("/tmp/gentleman-skills/.git")` succeeds AND mtime of `.git/FETCH_HEAD` is less than 1 hour old
    - If cache invalid: `git clone https://github.com/Gentleman-Programming/Gentleman-Skills.git /tmp/gentleman-skills` (or `git -C /tmp/gentleman-skills pull` if dir exists)
  - Spinner shown on screen with message "Fetching skill catalog..." via `SpinnerFrame` / `tickMsg`
  - Model field: `SkillListLoading bool`, `SkillList []SkillEntry` (where `SkillEntry` has `Name string`, `Description string`, `Source string`)
  - New message type: `skillsLoadedMsg struct { skills []SkillEntry; err error }`
  - On `skillsLoadedMsg`: set `SkillList`, clear `SkillListLoading`; if `err != nil` show error inline

---

### FR-3: Non-Interactive CLI

- FR-3.1: `--init-project` flag (bool) triggers project init flow without TUI
  - Parses remaining flags, runs `init-project.sh` directly, prints output to stdout, exits
  - Requires `--project-path`; exits with error if not provided

- FR-3.2: `--project-path=<path>` sets project directory
  - Same `~` expansion and `os.Stat()` validation as TUI flow
  - Exits with non-zero code and error message if path invalid

- FR-3.3: `--stack=<stack>` overrides auto-detection
  - Valid values: `java-gradle`, `java-maven`, `node`, `python`, `go`, `rust`
  - Passed as `--stack=<value>` to `init-project.sh`
  - If omitted, auto-detection runs (same file-check logic as TUI)

- FR-3.4: `--memory=<N>` sets memory module
  - Valid values: `1` (obsidian-brain), `2` (vibekanban), `3` (engram), `4` (simple), `5` (none)
  - Passed as `--memory=N` to `init-project.sh`
  - Default: `5` (none) if not specified

- FR-3.5: `--ci=<N>` sets CI provider
  - Valid values: `1` (github), `2` (gitlab), `3` (woodpecker), `4` (none)
  - Passed as `--ci=N` to `init-project.sh`
  - Default: `4` (none) if not specified

- FR-3.6: `--engram` flag (bool) adds Engram; passed as `--engram` to `init-project.sh`
  - Ignored if `--memory` is not `1` (engram only meaningful alongside Obsidian Brain)

- FR-3.7: `--skill-install=<name>` installs a named skill
  - Runs: `bash /tmp/gentleman-skills/add-skill.sh gentleman <name>`
  - Clones Gentleman-Skills repo if not cached (same logic as FR-2.7)
  - Can be repeated multiple times for multiple skills
  - Exits 0 on success, non-zero on first failure

- FR-3.8: `--skill-remove=<name>` removes a named skill
  - Runs: `bash /tmp/gentleman-skills/add-skill.sh remove <name>`
  - Can be repeated multiple times for multiple skills

---

### NFR: Non-Functional Requirements

- NFR-1: All new screens implement the Bubbletea Model-Update-View pattern; no goroutine-based state mutation outside `tea.Cmd`
- NFR-2: Navigation is consistent with existing screens:
  - `j` / `↓` → move cursor down
  - `k` / `↑` → move cursor up
  - `Enter` → confirm / select
  - `Esc` → go back one screen (not jump to MainMenu)
  - `g` → top of list (for scrollable screens)
  - `G` → bottom of list (for scrollable screens)
- NFR-3: Esc from any project-init or skill-manager screen returns to the immediately preceding screen (the back-link defined in Screen Mapping), never directly to Main Menu except from `ScreenSkillMenu` and the completion/result screens
- NFR-4: Path text input uses `charmbracelet/bubbles/textinput`; model stores the `textinput.Model` instance in `Model.ProjectPathInput textinput.Model`
- NFR-5: Clone operations (`stepInitProject` repo clone, Gentleman-Skills clone) stream progress lines via `SendLog` so the TUI log panel updates live
- NFR-6: All new screens have unit tests in:
  - `installer/internal/tui/project_screens_test.go` (project init flow)
  - `installer/internal/tui/skill_screens_test.go` (skill manager flow)
  - Tests must cover: initial state, cursor navigation, Enter/Esc transitions, conditional screen branching (Engram prompt), empty-state handling (no installed skills), inline error display on invalid path

---

## Scenarios

### S-1: Happy path — Initialize Node project with Obsidian Brain + GitHub Actions

Given: User selects "📦 Initialize Project" from Main Menu

When:
1. Enters path `/home/user/my-app` (valid directory with `package.json`)
2. `ScreenProjectStack` auto-detects `node (detected)`, user presses Enter to confirm
3. `ScreenProjectMemory` — user selects `Obsidian Brain`
4. `ScreenProjectEngram` — user selects `Yes, add Engram`
5. `ScreenProjectCI` — user selects `GitHub Actions`
6. `ScreenProjectConfirm` — user selects `✅ Run init-project.sh`

Then:
- `stepInitProject()` executes: `bash /tmp/project-starter-framework/init-project.sh --non-interactive --memory=1 --ci=1 --engram` with `Dir = "/home/user/my-app"`
- `ScreenProjectInstalling` shows progress/log output
- On exit code 0: transitions to `ScreenProjectComplete`
- `ScreenProjectComplete` shows created files summary and "Press Enter to return to Main Menu"

### S-2: Happy path — Initialize Go project with Simple memory, no CI

Given: User selects "📦 Initialize Project"

When:
1. Enters path `/home/user/go-service` (valid directory with `go.mod`)
2. Stack screen auto-detects `go (detected)`, user confirms
3. Memory screen — user selects `Simple`
4. (No Engram prompt — Obsidian Brain was NOT selected)
5. CI screen — user selects `None`
6. Confirmation screen — user confirms

Then:
- Executes: `bash /tmp/project-starter-framework/init-project.sh --non-interactive --memory=4 --ci=4`
- No `--engram` flag in command
- Completes normally

### S-3: Invalid path

Given: User is on `ScreenProjectPath`

When: User types `/home/user/nonexistent-dir` and presses Enter

Then:
- `os.Stat()` fails (or succeeds but `IsDir()` is false)
- Inline error displayed below input field: `Directory not found`
- Cursor remains on text input; screen does NOT advance

When: User clears input, types `/home/user/real-dir` (a valid directory), presses Enter

Then:
- Validation passes; screen advances to `ScreenProjectStack`

### S-4: Path is not a git repo

Given: User enters a valid directory that has no `.git/` subdirectory

Then:
- `ScreenProjectStack` shows info note: `No git repo detected — init-project.sh will run git init`
- Flow continues normally through all screens
- `init-project.sh` handles `git init` internally; no TUI-level blocking

### S-5: Skill Manager — Browse and Install

Given: User selects "🎯 Skill Manager" → "⬇️  Install Skill"

When:
1. `SkillListLoading = true`; spinner shown; Gentleman-Skills repo cloned to `/tmp/gentleman-skills/`
2. `skillsLoadedMsg` received; `SkillList` populated
3. `ScreenSkillInstall` renders multi-select list
4. User navigates to `react-19`, presses Enter to toggle; navigates to `typescript`, presses Enter to toggle
5. User navigates to `✅ Confirm installation`, presses Enter

Then:
- `stepSkillInstall()` runs in sequence:
  - `bash /tmp/gentleman-skills/add-skill.sh gentleman react-19`
  - `bash /tmp/gentleman-skills/add-skill.sh gentleman typescript`
- `ScreenSkillResult` shows:
  - `✅ react-19 — installed successfully`
  - `✅ typescript — installed successfully`

### S-6: Skill Manager — Remove

Given: User selects "🎯 Skill Manager" → "🗑️  Remove Skill"

When:
1. `InstalledSkills` populated from `~/.claude/skills/` directory scan
2. `ScreenSkillRemove` shows multi-select list of installed skills
3. User selects skills to remove, confirms

Then:
- `stepSkillRemove()` runs `bash /tmp/gentleman-skills/add-skill.sh remove <name>` for each selected skill
- `ScreenSkillResult` shows per-skill success/error

### S-7: Skill Manager — No skills installed (Remove)

Given: User selects "🎯 Skill Manager" → "🗑️  Remove Skill"

When: `~/.claude/skills/` directory is empty or does not exist

Then:
- `ScreenSkillRemove` shows "No skills installed" message
- Only option shown: `← Back`
- Enter or Esc → returns to `ScreenSkillMenu`

### S-8: Non-interactive project init

Given: CLI invoked with flags `--init-project --project-path=/home/user/app --memory=1 --ci=1 --engram`

Then:
- TUI is NOT launched (`nonInteractiveMode = true`)
- Path validated; stack auto-detected
- `init-project.sh --non-interactive --memory=1 --ci=1 --engram` executed in `/home/user/app`
- stdout/stderr streamed to terminal
- Exit code mirrors `init-project.sh` exit code

Note: `--ghagga` is not a supported flag (Ghagga integration is deferred).

---

## Screen Mapping

| Screen Constant | Type | Follows | Leads To |
|---|---|---|---|
| `ScreenProjectPath` | TextInput | MainMenu | `ScreenProjectStack` |
| `ScreenProjectStack` | SingleSelect + Spinner | `ScreenProjectPath` | `ScreenProjectMemory` |
| `ScreenProjectMemory` | SingleSelect | `ScreenProjectStack` | `ScreenProjectEngram` (if Obsidian Brain) or `ScreenProjectCI` |
| `ScreenProjectEngram` | YesNo | `ScreenProjectMemory` | `ScreenProjectCI` |
| `ScreenProjectCI` | SingleSelect | `ScreenProjectMemory` or `ScreenProjectEngram` | `ScreenProjectConfirm` |
| `ScreenProjectConfirm` | Confirm (Summary) | `ScreenProjectCI` | `ScreenProjectInstalling` |
| `ScreenProjectInstalling` | Progress (reuses Installing pattern) | `ScreenProjectConfirm` | `ScreenProjectComplete` |
| `ScreenProjectComplete` | Info | `ScreenProjectInstalling` | MainMenu |
| `ScreenSkillMenu` | SingleSelect | MainMenu | `ScreenSkillBrowse`, `ScreenSkillInstall`, or `ScreenSkillRemove` |
| `ScreenSkillBrowse` | ScrollList + Spinner | `ScreenSkillMenu` | `ScreenSkillMenu` |
| `ScreenSkillInstall` | MultiSelect + Scroll + Spinner | `ScreenSkillMenu` | `ScreenSkillResult` |
| `ScreenSkillRemove` | MultiSelect + Scroll | `ScreenSkillMenu` | `ScreenSkillResult` |
| `ScreenSkillResult` | Info | `ScreenSkillInstall` or `ScreenSkillRemove` | `ScreenSkillMenu` |

---

## Model Changes

### New fields on `UserChoices`

```go
// Project Init
ProjectPath      string // Absolute, ~ expanded
ProjectStack     string // e.g. "node", "go", "java-gradle"
ProjectMemory    string // Label e.g. "Obsidian Brain"
ProjectMemoryN   int    // 1-5 for --memory=N
ProjectAddEngram bool   // --engram flag
ProjectCI        string // Label e.g. "GitHub Actions"
ProjectCIN       int    // 1-4 for --ci=N
```

### New fields on `Model`

```go
// Project Init
ProjectPathInput    textinput.Model // bubbles/textinput instance
ProjectPathError    string          // Inline validation error; empty if none

// Skill Manager
SkillList           []SkillEntry // Available skills from Gentleman-Skills repo
SkillListLoading    bool         // True while cloning/fetching
InstalledSkills     []string     // Skills found in ~/.claude/skills/
SkillInstallSelected []bool      // Toggle state for install multi-select
SkillRemoveSelected  []bool      // Toggle state for remove multi-select
SkillBrowseScroll   int          // Viewport scroll offset for Browse screen
SkillInstallScroll  int          // Viewport scroll offset for Install screen
SkillRemoveScroll   int          // Viewport scroll offset for Remove screen
```

### New type

```go
type SkillEntry struct {
    Name        string // e.g. "react-19"
    Description string // One-line description
    Source      string // "curated" or "community"
}
```

### New screen constants (appended to existing iota block)

```go
// Project Init screens
ScreenProjectPath
ScreenProjectStack
ScreenProjectMemory
ScreenProjectEngram
ScreenProjectCI
ScreenProjectConfirm
ScreenProjectInstalling
ScreenProjectComplete
// Skill Manager screens
ScreenSkillMenu
ScreenSkillBrowse
ScreenSkillInstall
ScreenSkillRemove
ScreenSkillResult
```

---

## Affected Files

| File | Change Type | Description |
|---|---|---|
| `installer/internal/tui/model.go` | Modified | 13 new screen constants; new `SkillEntry` type; new `Model` fields for project path input, project choices, skill list state; new entries in `GetCurrentOptions`, `GetScreenTitle`, `GetScreenDescription`; `NewModel()` initializes new fields |
| `installer/internal/tui/update.go` | Modified | New message types (`skillsLoadedMsg`, `stackDetectedMsg`); handler functions for each new screen; dispatch in `handleMainMenuKeys` for "Initialize Project" and "Skill Manager"; `handleKeyPress` switch cases for all new `Screen*` constants |
| `installer/internal/tui/view.go` | Modified | New render functions: `viewProjectPath`, `viewProjectStack`, `viewProjectMemory`, `viewProjectEngram`, `viewProjectCI`, `viewProjectConfirm`, `viewProjectComplete`, `viewSkillMenu`, `viewSkillBrowse`, `viewSkillInstall`, `viewSkillRemove`, `viewSkillResult`; dispatch in main `View()` switch |
| `installer/internal/tui/installer.go` | Modified | `stepInitProject()` — clones framework repo if needed, builds and executes `init-project.sh` command, streams logs; `stepSkillInstall()` — iterates selected skills, runs `add-skill.sh gentleman <name>` for each; `stepSkillRemove()` — runs `add-skill.sh remove <name>` for each selected skill; `loadGentlemanSkills()` — clones/pulls Gentleman-Skills repo, parses skill list |
| `installer/cmd/gentleman-installer/main.go` | Modified | New flags: `--init-project` (bool), `--project-path` (string), `--stack` (string), `--memory` (int), `--ci` (int), `--engram` (bool), `--skill-install` (string, repeatable), `--skill-remove` (string, repeatable); non-interactive execution path branching |
| `installer/internal/tui/project_screens_test.go` | New | Tests for all project init screens: initial render, cursor movement, Enter/Esc transitions, `~` expansion, invalid path error display, conditional Engram screen branching, confirmation summary content, stack auto-detection logic |
| `installer/internal/tui/skill_screens_test.go` | New | Tests for skill manager flow: sub-menu navigation, browse scroll, install multi-select toggle, remove with empty installed list, result screen per-skill output, cache check logic |
