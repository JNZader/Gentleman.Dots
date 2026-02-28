# Design: Project Init & Skill Manager

## Architecture Decision

Hybrid pattern: TUI screens collect user choices with drill-down navigation, then delegate execution to framework scripts via `system.RunWithLogs`. Identical to the AI Framework integration (`stepInstallAIFramework` in `installer.go`). No new external dependencies — path input uses the existing manual key-accumulation pattern from `TrainerInput`, not `charmbracelet/bubbles/textinput` (which is not in `go.mod`).

Skill list loading uses a dedicated message type (`skillsLoadedMsg`) and a Bubbletea `Cmd` that clones/reads the repo in the background, populating `Model` fields before rendering.

---

## New Screen Constants

Add to `model.go` after `ScreenTrainerBossResult` (currently the last constant, value 40):

```go
// Project Init screens
ScreenProjectPath        // 41 — text input: project directory
ScreenProjectStack       // 42 — single-select: detected stack confirmation/override
ScreenProjectMemory      // 43 — single-select: memory module
ScreenProjectEngram      // 44 — yes/no: add Engram alongside Obsidian Brain
ScreenProjectCI          // 45 — single-select: CI provider
ScreenProjectConfirm     // 46 — summary before execution
ScreenProjectInstalling  // 47 — progress log (reuses existing installing pattern)
ScreenProjectResult      // 48 — success/error

// Skill Manager screens
ScreenSkillMenu          // 49 — Browse / Install / Remove
ScreenSkillBrowse        // 50 — scrollable read-only list
ScreenSkillInstall       // 51 — multi-select from available skills
ScreenSkillRemove        // 52 — multi-select from installed skills
ScreenSkillResult        // 53 — success/error output
```

Total new constants: 13 (6 project init + 7 skill manager). Existing iota block continues unbroken — no renumbering needed.

---

## New Model Fields

Add to the `Model` struct in `model.go`:

```go
// Project Init state
ProjectPathInput    string   // Raw text being typed (grows/shrinks via backspace)
ProjectPathError    string   // Validation error message shown inline ("" = valid)
ProjectStack        string   // Detected or user-selected stack ("react", "angular", "python", etc.)
ProjectMemory       string   // Selected memory module ID: "obsidian-brain", "vibekanban", "engram", "simple", "none"
ProjectEngram       bool     // Add Engram alongside Obsidian Brain
ProjectCI           string   // Selected CI provider ID: "github", "gitlab", "woodpecker", "none"
ProjectLogLines     []string // Log output from init-project.sh (separate from main LogLines)

// Skill Manager state
SkillList           []string        // Available skills from Gentleman-Skills repo
InstalledSkills     []string        // Skills already installed in current project
SkillSelected       []bool          // Toggle state for install/remove multi-select
SkillScroll         int             // Scroll offset for skill list screens
SkillLoading        bool            // True while async skill fetch is running
SkillLoadError      string          // Error message if skill fetch fails
SkillResultLog      []string        // Output lines from add-skill.sh
```

Initialize in `NewModel()`:

```go
ProjectPathInput:  "",
ProjectPathError:  "",
ProjectStack:      "",
ProjectMemory:     "",
ProjectEngram:     false,
ProjectCI:         "",
ProjectLogLines:   []string{},
SkillList:         []string{},
InstalledSkills:   []string{},
SkillSelected:     []bool{},
SkillScroll:       0,
SkillLoading:      false,
SkillLoadError:    "",
SkillResultLog:    []string{},
```

---

## New UserChoices Fields

Add to `UserChoices` struct in `model.go` for non-interactive CLI passthrough:

```go
// Project Init (non-interactive)
InitProject     bool   // Whether to run project init flow
ProjectPath     string // Target project directory path
ProjectStack    string // Stack override (empty = auto-detect)
ProjectMemory   string // Memory module: obsidian-brain, vibekanban, engram, simple, none
ProjectCI       string // CI provider: github, gitlab, woodpecker, none
ProjectEngram   bool   // Add Engram alongside Obsidian Brain

// Skill Manager (non-interactive)
SkillAction     string   // "install" or "remove"
SkillNames      []string // Skill names to install or remove
```

---

## Screen Flow Diagrams

### Initialize Project Flow

```
ScreenMainMenu
  │
  └─ "📦 Initialize Project"
       │
       ▼
  ScreenProjectPath  ──[esc]──► ScreenMainMenu
       │ [enter, valid path]
       ▼
  ScreenProjectStack  ──[esc]──► ScreenProjectPath
       │ [enter]
       ▼
  ScreenProjectMemory  ──[esc]──► ScreenProjectStack
       │
       ├─ "obsidian-brain" ──► ScreenProjectEngram ──[esc]──► ScreenProjectMemory
       │                             │ [enter yes/no]
       │                             ▼
       │                       ScreenProjectCI ──[esc]──► ScreenProjectEngram
       │                             │ [enter]
       │                             ▼
       │                       ScreenProjectConfirm ──[esc]──► ScreenProjectCI
       │
       └─ other memory ──► ScreenProjectCI ──[esc]──► ScreenProjectMemory
                                 │ [enter]
                                 ▼
                           ScreenProjectConfirm ──[esc]──► ScreenProjectCI
                                 │ [enter: "Confirm & Run"]
                                 ▼
                           ScreenProjectInstalling
                                 │ (stepInitProject runs async)
                                 ▼
                           ScreenProjectResult
                                 │ [enter or esc]
                                 ▼
                           ScreenMainMenu
```

### Skill Manager Flow

```
ScreenMainMenu
  │
  └─ "🎯 Skill Manager"
       │
       ▼
  ScreenSkillMenu  ──[esc]──► ScreenMainMenu
       │
       ├─ "Browse"
       │     │
       │     ▼
       │  ScreenSkillBrowse (async load ─► skillsLoadedMsg populates SkillList)
       │     │ [esc]
       │     ▼
       │  ScreenSkillMenu
       │
       ├─ "Install"
       │     │
       │     ▼
       │  ScreenSkillInstall (async load if needed)
       │     │ [enter: "✅ Confirm"]
       │     ▼
       │  ScreenSkillResult
       │     │ [enter or esc]
       │     ▼
       │  ScreenSkillMenu
       │
       └─ "Remove"
             │
             ▼
          ScreenSkillRemove (loads InstalledSkills)
             │ [enter: "✅ Confirm"]
             ▼
          ScreenSkillResult
             │ [enter or esc]
             ▼
          ScreenSkillMenu
```

---

## Handler Functions

### `handleMainMenuKeys` additions (update.go)

In the existing `handleMainMenuKeys` `switch` block, add two new `case` branches inside the `"enter", " "` handler:

```go
case strings.Contains(selected, "Initialize Project"):
    m.Screen = ScreenProjectPath
    m.ProjectPathInput = ""
    m.ProjectPathError = ""
    m.Cursor = 0

case strings.Contains(selected, "Skill Manager"):
    m.Screen = ScreenSkillMenu
    m.Cursor = 0
```

`GetCurrentOptions()` for `ScreenMainMenu` gains two entries before `"❌ Exit"`:

```go
"📦 Initialize Project",
"🎯 Skill Manager",
```

---

### `handleProjectPathKeys(key string) (tea.Model, tea.Cmd)` — update.go

**Signature:** `func (m Model) handleProjectPathKeys(key string) (tea.Model, tea.Cmd)`

**Key bindings:**

| Key | Action |
|-----|--------|
| `esc` | Return to `ScreenMainMenu`, reset `ProjectPathInput` |
| `backspace` | Remove last rune from `ProjectPathInput` |
| `enter` | Validate path → advance or show inline error |
| any printable char (len == 1) | Append to `ProjectPathInput` |

**Path validation on `enter`:**

```go
case "enter":
    path := expandPath(m.ProjectPathInput) // see below
    if path == "" {
        m.ProjectPathError = "Path cannot be empty"
        return m, nil
    }
    info, err := os.Stat(path)
    if err != nil {
        m.ProjectPathError = fmt.Sprintf("Path does not exist: %s", path)
        return m, nil
    }
    if !info.IsDir() {
        m.ProjectPathError = "Path must be a directory"
        return m, nil
    }
    // Valid — store canonical path and advance
    absPath, _ := filepath.Abs(path)
    m.ProjectPathInput = absPath
    m.ProjectPathError = ""
    m.ProjectStack = detectStack(absPath) // see Skill List Loading section
    m.Screen = ScreenProjectStack
    m.Cursor = 0
```

**Helper `expandPath`:**

```go
func expandPath(p string) string {
    if strings.HasPrefix(p, "~/") {
        home, _ := os.UserHomeDir()
        return filepath.Join(home, p[2:])
    }
    return p
}
```

**Helper `detectStack`:** runs `ls` on the path and returns the best guess:

```go
func detectStack(path string) string {
    checks := []struct{ file, stack string }{
        {"package.json", "node"},
        {"angular.json", "angular"},
        {"next.config.*", "nextjs"},
        {"pyproject.toml", "python"},
        {"go.mod", "go"},
        {"Cargo.toml", "rust"},
        {"pom.xml", "java"},
    }
    for _, c := range checks {
        matches, _ := filepath.Glob(filepath.Join(path, c.file))
        if len(matches) > 0 {
            return c.stack
        }
    }
    return "unknown"
}
```

**Space key:** `ScreenProjectPath` must be added to the exclusion list in `handleKeyPress` so that space appends a literal space character instead of activating leader mode:

```go
case ScreenProjectPath:
    // space is part of path input — fall through to handleProjectPathKeys
```

---

### `handleProjectStackKeys(key string) (tea.Model, tea.Cmd)` — update.go

**Follows `handleSelectionKeys` pattern exactly.**

| Key | Action |
|-----|--------|
| `up`/`k`, `down`/`j` | Move cursor, skip separators |
| `esc`, `backspace` | `m.Screen = ScreenProjectPath; m.Cursor = 0` |
| `enter` | Save selected stack label → `m.ProjectStack` → `m.Screen = ScreenProjectMemory` |

Stack options are defined in `GetCurrentOptions` (see below). User can override the auto-detected stack.

---

### `handleProjectMemoryKeys(key string) (tea.Model, tea.Cmd)` — update.go

**Follows `handleSelectionKeys` pattern.**

| Key | Action |
|-----|--------|
| `esc` | Back to `ScreenProjectStack` |
| `enter` | Save memory ID → advance conditionally |

```go
case "enter":
    ids := []string{"obsidian-brain", "vibekanban", "engram", "simple", "none"}
    // cursor 0-4 maps to ids, cursor 5 is separator (skip), cursor 6 is never reachable
    if m.Cursor < len(ids) {
        m.ProjectMemory = ids[m.Cursor]
    }
    if m.ProjectMemory == "obsidian-brain" {
        m.Screen = ScreenProjectEngram
    } else {
        m.Screen = ScreenProjectCI
    }
    m.Cursor = 0
```

---

### `handleProjectEngramKeys(key string) (tea.Model, tea.Cmd)` — update.go

**Yes/No screen. Follows `handleSelectionKeys` pattern with 2 options.**

| Key | Action |
|-----|--------|
| `esc` | Back to `ScreenProjectMemory` |
| `enter` | `m.ProjectEngram = (m.Cursor == 0)` → `m.Screen = ScreenProjectCI; m.Cursor = 0` |

---

### `handleProjectCIKeys(key string) (tea.Model, tea.Cmd)` — update.go

**Follows `handleSelectionKeys` pattern.**

| Key | Action |
|-----|--------|
| `esc` | Back to `ScreenProjectEngram` if `m.ProjectMemory == "obsidian-brain"`, else back to `ScreenProjectMemory` |
| `enter` | Save CI ID → advance conditionally |

```go
case "enter":
    ids := []string{"github", "gitlab", "woodpecker", "none"}
    if m.Cursor < len(ids) {
        m.ProjectCI = ids[m.Cursor]
    }
    m.Screen = ScreenProjectConfirm
    m.Cursor = 0
```

---

### `handleProjectConfirmKeys(key string) (tea.Model, tea.Cmd)` — update.go

**Single-select: "Confirm & Run" or "Cancel".**

| Key | Action |
|-----|--------|
| `esc` | Back to `ScreenProjectCI` |
| `enter`, cursor=0 | Start execution |
| `enter`, cursor=1 | `m.Screen = ScreenMainMenu; m.Cursor = 0` |

Execution on confirm:

```go
m.ProjectLogLines = []string{}
m.Screen = ScreenProjectInstalling
m.SpinnerFrame = 0
return m, func() tea.Msg { return projectInstallStartMsg{} }
```

---

### `handleProjectInstallingKeys` — update.go

Only `ctrl+c` is handled (quit). All other keys are ignored. The screen exits automatically when `projectInstallCompleteMsg` or `projectInstallErrorMsg` arrives.

Add to `Update()` message switch:

```go
case projectInstallStartMsg:
    return m, runProjectInitCmd(&m)

case projectInstallLogMsg:
    m.ProjectLogLines = append(m.ProjectLogLines, msg.line)
    if len(m.ProjectLogLines) > 30 {
        m.ProjectLogLines = m.ProjectLogLines[len(m.ProjectLogLines)-30:]
    }
    m.SpinnerFrame++
    return m, nil

case projectInstallCompleteMsg:
    m.Screen = ScreenProjectResult
    if msg.err != nil {
        m.ErrorMsg = msg.err.Error()
    } else {
        m.ErrorMsg = ""
    }
    return m, nil
```

---

### `handleProjectResultKeys(key string) (tea.Model, tea.Cmd)` — update.go

| Key | Action |
|-----|--------|
| `enter`, `esc` | `m.Screen = ScreenMainMenu; m.Cursor = 0` |

---

### `handleSkillMenuKeys(key string) (tea.Model, tea.Cmd)` — update.go

**Follows `handleSelectionKeys` pattern exactly.**

| Key | Action |
|-----|--------|
| `esc` | `m.Screen = ScreenMainMenu; m.Cursor = 0` |
| `enter` on "Browse" | `m.Screen = ScreenSkillBrowse; m.SkillLoading = true; return m, loadSkillsCmd()` |
| `enter` on "Install" | `m.Screen = ScreenSkillInstall; m.SkillLoading = true; return m, loadSkillsCmd()` |
| `enter` on "Remove" | `m.Screen = ScreenSkillRemove; return m, loadInstalledSkillsCmd()` |

Add to `Update()`:

```go
case skillsLoadedMsg:
    m.SkillLoading = false
    if msg.err != nil {
        m.SkillLoadError = msg.err.Error()
    } else {
        m.SkillLoadError = ""
        m.SkillList = msg.available
        m.InstalledSkills = msg.installed
        m.SkillSelected = make([]bool, len(msg.available))
    }
    return m, nil
```

---

### `handleSkillBrowseKeys(key string) (tea.Model, tea.Cmd)` — update.go

Read-only list. Follows keymap scroll pattern.

| Key | Action |
|-----|--------|
| `up`/`k` | `m.SkillScroll = max(0, m.SkillScroll-1)` |
| `down`/`j` | `m.SkillScroll = min(len(m.SkillList)-1, m.SkillScroll+1)` |
| `esc` | `m.Screen = ScreenSkillMenu; m.Cursor = 0; m.SkillScroll = 0` |

---

### `handleSkillInstallKeys(key string) (tea.Model, tea.Cmd)` — update.go

**Multi-select, follows `handleAIToolsKeys` pattern.**

| Key | Action |
|-----|--------|
| `up`/`k` | Move cursor up, update `SkillScroll` if needed |
| `down`/`j` | Move cursor down, update `SkillScroll` if needed |
| `esc` | `m.Screen = ScreenSkillMenu; m.Cursor = 0` |
| `enter` on skill item | Toggle `m.SkillSelected[m.Cursor]` |
| `enter` on "✅ Confirm" | Collect selected names → run install → `m.Screen = ScreenSkillResult` |

Execution on confirm:

```go
var names []string
for i, sel := range m.SkillSelected {
    if sel && i < len(m.SkillList) {
        names = append(names, m.SkillList[i])
    }
}
if len(names) == 0 {
    return m, nil // nothing selected
}
m.SkillResultLog = []string{}
m.Screen = ScreenSkillResult
m.SkillLoading = true
return m, runSkillInstallCmd(names)
```

Add to `Update()`:

```go
case skillActionCompleteMsg:
    m.SkillLoading = false
    m.SkillResultLog = msg.logLines
    if msg.err != nil {
        m.ErrorMsg = msg.err.Error()
    } else {
        m.ErrorMsg = ""
    }
    return m, nil
```

---

### `handleSkillRemoveKeys(key string) (tea.Model, tea.Cmd)` — update.go

**Identical structure to `handleSkillInstallKeys` but operates on `InstalledSkills`.**

On confirm, run `runSkillRemoveCmd(names)` which calls `add-skill.sh remove <name>` for each selected skill.

---

### `handleSkillResultKeys(key string) (tea.Model, tea.Cmd)` — update.go

| Key | Action |
|-----|--------|
| `enter`, `esc` | `m.Screen = ScreenSkillMenu; m.Cursor = 0; m.SkillLoading = false; m.SkillLoadError = ""` |

---

### Dispatch additions in `handleKeyPress` switch — update.go

```go
case ScreenProjectPath:
    return m.handleProjectPathKeys(key)

case ScreenProjectStack, ScreenProjectMemory, ScreenProjectEngram,
     ScreenProjectCI, ScreenProjectConfirm:
    return m.handleProjectSelectionKeys(key)
    // (single dispatch function that branches internally by m.Screen)

case ScreenProjectResult:
    return m.handleProjectResultKeys(key)

case ScreenSkillMenu:
    return m.handleSkillMenuKeys(key)

case ScreenSkillBrowse:
    return m.handleSkillBrowseKeys(key)

case ScreenSkillInstall:
    return m.handleSkillInstallKeys(key)

case ScreenSkillRemove:
    return m.handleSkillRemoveKeys(key)

case ScreenSkillResult:
    return m.handleSkillResultKeys(key)
```

`ScreenProjectInstalling` needs NO key handler — ignore all keys (same as `ScreenInstalling`).

---

### `handleEscape` additions — update.go

```go
case ScreenProjectPath:
    m.Screen = ScreenMainMenu
    m.ProjectPathInput = ""
    m.ProjectPathError = ""
    m.Cursor = 0

case ScreenProjectStack:
    m.Screen = ScreenProjectPath
    m.Cursor = 0

case ScreenProjectMemory:
    m.Screen = ScreenProjectStack
    m.Cursor = 0

case ScreenProjectEngram:
    m.Screen = ScreenProjectMemory
    m.Cursor = 0

case ScreenProjectCI:
    if m.ProjectMemory == "obsidian-brain" {
        m.Screen = ScreenProjectEngram
    } else {
        m.Screen = ScreenProjectMemory
    }
    m.Cursor = 0

case ScreenProjectConfirm:
    m.Screen = ScreenProjectCI
    m.Cursor = 0

case ScreenProjectResult, ScreenSkillMenu:
    m.Screen = ScreenMainMenu
    m.Cursor = 0

case ScreenSkillBrowse, ScreenSkillInstall, ScreenSkillRemove:
    m.Screen = ScreenSkillMenu
    m.Cursor = 0
    m.SkillScroll = 0

case ScreenSkillResult:
    m.Screen = ScreenSkillMenu
    m.Cursor = 0
```

---

## View Functions

### `renderProjectPath()` — view.go

**What it displays:** Full-screen text input with inline error. No step-progress bar (this is a separate top-level flow, not the install wizard).

**Pattern:** New — closest analog is `renderWelcome` for the centered layout but with an input field box.

```go
func (m Model) renderProjectPath() string {
    var s strings.Builder

    s.WriteString(TitleStyle.Render("📦 Initialize Project"))
    s.WriteString("\n")
    s.WriteString(MutedStyle.Render("Enter the path to the project directory"))
    s.WriteString("\n\n")

    // Input box
    inputDisplay := m.ProjectPathInput
    if inputDisplay == "" {
        inputDisplay = " "
    }
    inputLine := fmt.Sprintf("> %s█", inputDisplay)
    s.WriteString(BoxStyle.Render(inputLine))
    s.WriteString("\n")

    // Inline error
    if m.ProjectPathError != "" {
        s.WriteString("\n")
        s.WriteString(ErrorStyle.Render("⚠ " + m.ProjectPathError))
        s.WriteString("\n")
    }

    s.WriteString("\n")
    s.WriteString(HelpStyle.Render("[Enter] confirm • [Esc] back • ~ is expanded to $HOME"))

    return s.String()
}
```

The `█` cursor blink can be simplified: always shown (no tick needed). If animation is desired, toggle based on `m.SpinnerFrame%2`.

---

### `renderProjectStack()` — view.go

**Pattern:** `renderSelection` exactly. Title + description + option list with cursor.

Description shows the auto-detected value:

```go
// In GetScreenDescription:
case ScreenProjectStack:
    if m.ProjectStack != "" && m.ProjectStack != "unknown" {
        return fmt.Sprintf("Auto-detected: %s — confirm or select a different stack", m.ProjectStack)
    }
    return "Select the primary technology stack for this project"
```

---

### `renderProjectMemory()` — view.go

**Pattern:** `renderSelection`. Single-select list.

---

### `renderProjectEngram()` — view.go

**Pattern:** `renderSelection` with 2 options (Yes/No). No separator lines needed.

---

### `renderProjectCI()` — view.go

**Pattern:** `renderSelection`.

---

### `renderProjectConfirm()` — view.go

**What it displays:** Summary of all choices, then action buttons.

```go
func (m Model) renderProjectConfirm() string {
    var s strings.Builder

    s.WriteString(TitleStyle.Render("📦 Initialize Project — Confirm"))
    s.WriteString("\n\n")

    s.WriteString(InfoStyle.Render(fmt.Sprintf("  Path:    %s", m.ProjectPathInput)))
    s.WriteString("\n")
    s.WriteString(InfoStyle.Render(fmt.Sprintf("  Stack:   %s", m.ProjectStack)))
    s.WriteString("\n")
    s.WriteString(InfoStyle.Render(fmt.Sprintf("  Memory:  %s", m.ProjectMemory)))
    if m.ProjectEngram {
        s.WriteString(InfoStyle.Render("  + Engram"))
    }
    s.WriteString("\n")
    s.WriteString(InfoStyle.Render(fmt.Sprintf("  CI:      %s", m.ProjectCI)))
    s.WriteString("\n\n")

    // Action buttons (options from GetCurrentOptions)
    options := m.GetCurrentOptions()
    for i, opt := range options {
        cursor := "  "
        style := UnselectedStyle
        if i == m.Cursor {
            cursor = "▸ "
            style = SelectedStyle
        }
        s.WriteString(style.Render(cursor + opt))
        s.WriteString("\n")
    }

    s.WriteString("\n")
    s.WriteString(HelpStyle.Render("↑/k up • ↓/j down • [Enter] select • [Esc] back"))

    return s.String()
}
```

---

### `renderProjectInstalling()` — view.go

**Pattern:** Reuses `renderInstalling` layout but reads from `m.ProjectLogLines` instead of `m.LogLines` and shows no step-progress bar.

```go
func (m Model) renderProjectInstalling() string {
    var s strings.Builder

    spinners := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
    spinner := spinners[m.SpinnerFrame%len(spinners)]

    s.WriteString(TitleStyle.Render(fmt.Sprintf("%s Initializing project...", spinner)))
    s.WriteString("\n\n")

    for _, line := range m.ProjectLogLines {
        s.WriteString(MutedStyle.Render("  " + line))
        s.WriteString("\n")
    }

    s.WriteString("\n")
    s.WriteString(HelpStyle.Render("[Ctrl+C] quit"))

    return s.String()
}
```

---

### `renderProjectResult()` — view.go

**Pattern:** `renderComplete` / `renderError` depending on `m.ErrorMsg`.

```go
func (m Model) renderProjectResult() string {
    var s strings.Builder

    if m.ErrorMsg != "" {
        s.WriteString(ErrorStyle.Render("✗ Initialization failed"))
        s.WriteString("\n\n")
        s.WriteString(MutedStyle.Render(m.ErrorMsg))
    } else {
        s.WriteString(SuccessStyle.Render("✓ Project initialized successfully"))
        s.WriteString("\n\n")
        s.WriteString(MutedStyle.Render(fmt.Sprintf("  %s is ready.", m.ProjectPathInput)))
    }

    s.WriteString("\n\n")
    s.WriteString(HelpStyle.Render("[Enter] return to main menu"))

    return s.String()
}
```

---

### `renderSkillMenu()` — view.go

**Pattern:** `renderMainMenu` (simple option list, no step progress).

---

### `renderSkillBrowse()` — view.go

**Pattern:** `renderAICategoryItems` scroll pattern — read-only, no checkboxes.

```go
func (m Model) renderSkillBrowse() string {
    var s strings.Builder

    s.WriteString(TitleStyle.Render("🎯 Skill Manager — Browse"))
    s.WriteString("\n")

    if m.SkillLoading {
        s.WriteString(MutedStyle.Render("  Fetching skill catalog..."))
        s.WriteString("\n")
        return s.String()
    }
    if m.SkillLoadError != "" {
        s.WriteString(ErrorStyle.Render("  ✗ " + m.SkillLoadError))
        s.WriteString("\n\n")
        s.WriteString(HelpStyle.Render("[Esc] back"))
        return s.String()
    }

    s.WriteString(MutedStyle.Render(fmt.Sprintf("  %d skills available", len(m.SkillList))))
    s.WriteString("\n\n")

    // Viewport scroll
    visibleItems := m.Height - 6
    if visibleItems < 5 { visibleItems = 5 }
    start := m.SkillScroll
    end := start + visibleItems
    if end > len(m.SkillList) { end = len(m.SkillList) }

    if start > 0 {
        s.WriteString(MutedStyle.Render(fmt.Sprintf("  ▲ %d more above", start)))
        s.WriteString("\n")
    }
    for i := start; i < end; i++ {
        s.WriteString(UnselectedStyle.Render("  " + m.SkillList[i]))
        s.WriteString("\n")
    }
    if end < len(m.SkillList) {
        s.WriteString(MutedStyle.Render(fmt.Sprintf("  ▼ %d more below", len(m.SkillList)-end)))
        s.WriteString("\n")
    }

    s.WriteString("\n")
    s.WriteString(HelpStyle.Render("↑/k up • ↓/j down • [Esc] back"))

    return s.String()
}
```

---

### `renderSkillInstall()` / `renderSkillRemove()` — view.go

**Pattern:** `renderAIToolSelection` (checkbox multi-select + scrolling from `renderAICategoryItems`).

```go
func (m Model) renderSkillInstall() string {
    var s strings.Builder
    s.WriteString(TitleStyle.Render("🎯 Skill Manager — Install"))
    s.WriteString("\n")

    if m.SkillLoading {
        s.WriteString(MutedStyle.Render("  Fetching skill catalog..."))
        s.WriteString("\n")
        return s.String()
    }
    if m.SkillLoadError != "" {
        s.WriteString(ErrorStyle.Render("  ✗ " + m.SkillLoadError))
        s.WriteString("\n\n")
        s.WriteString(HelpStyle.Render("[Esc] back"))
        return s.String()
    }

    options := m.GetCurrentOptions() // SkillList + separator + "✅ Confirm"

    visibleItems := m.Height - 6
    if visibleItems < 5 { visibleItems = 5 }
    start := m.SkillScroll
    end := start + visibleItems
    if end > len(options) { end = len(options) }

    if start > 0 {
        s.WriteString(MutedStyle.Render(fmt.Sprintf("  ▲ %d more above", start)))
        s.WriteString("\n")
    }

    for i := start; i < end; i++ {
        opt := options[i]
        if strings.HasPrefix(opt, "───") {
            s.WriteString(MutedStyle.Render(opt))
            s.WriteString("\n")
            continue
        }

        cursor := "  "
        style := UnselectedStyle
        if i == m.Cursor {
            cursor = "▸ "
            style = SelectedStyle
        }

        checkbox := "[ ] "
        if strings.HasPrefix(opt, "✅") {
            checkbox = ""
        } else if m.SkillSelected != nil && i < len(m.SkillSelected) && m.SkillSelected[i] {
            checkbox = "[✓] "
        }

        s.WriteString(style.Render(cursor + checkbox + opt))
        s.WriteString("\n")
    }

    if end < len(options) {
        s.WriteString(MutedStyle.Render(fmt.Sprintf("  ▼ %d more below", len(options)-end)))
        s.WriteString("\n")
    }

    s.WriteString("\n")
    s.WriteString(HelpStyle.Render("↑/k up • ↓/j down • [Enter] toggle/confirm • [Esc] back"))
    return s.String()
}
```

`renderSkillRemove` is identical but shows `InstalledSkills` instead of `SkillList` and says "Remove" in the title.

---

### `renderSkillResult()` — view.go

**Pattern:** `renderProjectResult`. Shows `m.SkillResultLog` lines and success/error state.

---

### View dispatch additions in `View()` — view.go

```go
case ScreenProjectPath:
    s.WriteString(m.renderProjectPath())
case ScreenProjectStack, ScreenProjectMemory, ScreenProjectEngram,
     ScreenProjectCI:
    s.WriteString(m.renderSelection()) // GetScreenTitle/Description handle per-screen text
case ScreenProjectConfirm:
    s.WriteString(m.renderProjectConfirm())
case ScreenProjectInstalling:
    s.WriteString(m.renderProjectInstalling())
case ScreenProjectResult:
    s.WriteString(m.renderProjectResult())
case ScreenSkillMenu:
    s.WriteString(m.renderSkillMenu())
case ScreenSkillBrowse:
    s.WriteString(m.renderSkillBrowse())
case ScreenSkillInstall:
    s.WriteString(m.renderSkillInstall())
case ScreenSkillRemove:
    s.WriteString(m.renderSkillRemove())
case ScreenSkillResult:
    s.WriteString(m.renderSkillResult())
```

---

## Installer Functions

### `stepInitProject(m *Model) error` — installer.go

**What it runs:** Reuses `/tmp/project-starter-framework-install` if present (from a prior AI Framework install); otherwise clones fresh. Runs `init-project.sh` with `--non-interactive` flags.

```go
func stepInitProject(m *Model) error {
    stepID := "initproject"

    clonePath := "/tmp/project-starter-framework-install"

    // Reuse existing clone if < 1 hour old, else re-clone
    info, err := os.Stat(clonePath)
    needsClone := err != nil || time.Since(info.ModTime()) > time.Hour

    if needsClone {
        system.Run("rm -rf "+clonePath, nil)
        SendLog(stepID, "Cloning project-starter-framework...")
        result := system.RunWithLogs(
            "git clone --depth 1 https://github.com/JNZader/project-starter-framework.git "+clonePath,
            nil, func(line string) { SendLog(stepID, line) },
        )
        if result.Error != nil {
            return wrapStepError(stepID, "Initialize Project",
                "Failed to clone project-starter-framework", result.Error)
        }
    }

    // Make script executable
    system.Run("chmod +x "+clonePath+"/init-project.sh", nil)

    // Build command
    memoryFlag := fmt.Sprintf("--memory=%s", m.Choices.ProjectMemory)
    ciFlag := fmt.Sprintf("--ci=%s", m.Choices.ProjectCI)

    cmd := fmt.Sprintf("bash %s/init-project.sh --non-interactive %s %s",
        clonePath, memoryFlag, ciFlag)

    if m.Choices.ProjectEngram {
        cmd += " --engram"
    }

    // Run from target project directory
    SendLog(stepID, fmt.Sprintf("Running: %s", cmd))
    SendLog(stepID, fmt.Sprintf("Target:  %s", m.Choices.ProjectPath))

    result := system.RunWithLogsInDir(cmd, m.Choices.ProjectPath, func(line string) {
        SendLog(stepID, line)
    })
    if result.Error != nil {
        return wrapStepError(stepID, "Initialize Project",
            "init-project.sh failed", result.Error)
    }

    SendLog(stepID, "✓ Project initialized")
    return nil
}
```

**Note:** `system.RunWithLogsInDir` may need to be added to the `system` package if it doesn't exist. It is equivalent to `system.RunWithLogs` but sets `cmd.Dir`. Check `system` package before adding.

---

### `stepSkillInstall(m *Model) error` — installer.go

```go
func stepSkillInstall(m *Model) error {
    stepID := "skillinstall"

    clonePath := "/tmp/project-starter-framework-install"
    // Reuse or re-clone (same logic as stepInitProject above)

    system.Run("chmod +x "+clonePath+"/add-skill.sh", nil)

    var failed []string
    for _, name := range m.Choices.SkillNames {
        SendLog(stepID, fmt.Sprintf("Installing skill: %s", name))
        cmd := fmt.Sprintf("bash %s/add-skill.sh gentleman %s", clonePath, name)
        result := system.RunWithLogsInDir(cmd, m.Choices.ProjectPath, func(line string) {
            SendLog(stepID, line)
        })
        if result.Error != nil {
            SendLog(stepID, fmt.Sprintf("⚠ Failed to install %s: %v", name, result.Error))
            failed = append(failed, name)
        } else {
            SendLog(stepID, fmt.Sprintf("✓ %s installed", name))
        }
    }

    if len(failed) > 0 {
        return fmt.Errorf("failed to install skills: %s", strings.Join(failed, ", "))
    }
    return nil
}
```

---

### `stepSkillRemove(m *Model) error` — installer.go

Identical to `stepSkillInstall` but runs `add-skill.sh remove <name>` instead:

```go
cmd := fmt.Sprintf("bash %s/add-skill.sh remove %s", clonePath, name)
```

---

### `executeStep` additions — installer.go

```go
case "initproject":
    return stepInitProject(m)
case "skillinstall":
    return stepSkillInstall(m)
case "skillremove":
    return stepSkillRemove(m)
```

---

### Bubbletea Cmd functions — update.go

These run in goroutines and return messages:

```go
// projectInstallStartMsg triggers the project init step
type projectInstallStartMsg struct{}

// projectInstallLogMsg carries a single log line from init-project.sh
type projectInstallLogMsg struct{ line string }

// projectInstallCompleteMsg signals completion (err == nil = success)
type projectInstallCompleteMsg struct{ err error }

// skillsLoadedMsg carries the result of async skill catalog fetch
type skillsLoadedMsg struct {
    available []string
    installed []string
    err       error
}

// skillActionCompleteMsg carries the result of install/remove
type skillActionCompleteMsg struct {
    logLines []string
    err      error
}

// runProjectInitCmd executes init-project.sh asynchronously
func runProjectInitCmd(m *Model) tea.Cmd {
    // Capture choices at call time (Model is a value, but we need the path etc.)
    choices := m.Choices
    sysInfo := m.SystemInfo
    return func() tea.Msg {
        model := &Model{SystemInfo: sysInfo, Choices: choices, LogLines: []string{}}
        err := stepInitProject(model)
        return projectInstallCompleteMsg{err: err}
    }
}

// loadSkillsCmd clones Gentleman-Skills and parses available + installed skills
func loadSkillsCmd() tea.Cmd {
    return func() tea.Msg {
        available, installed, err := fetchSkillCatalog()
        return skillsLoadedMsg{available: available, installed: installed, err: err}
    }
}

// runSkillInstallCmd installs selected skills
func runSkillInstallCmd(projectPath string, names []string) tea.Cmd {
    return func() tea.Msg {
        var logs []string
        // ... (calls stepSkillInstall internally or runs inline)
        return skillActionCompleteMsg{logLines: logs, err: nil}
    }
}
```

---

## GetCurrentOptions Additions

Add to `GetCurrentOptions()` switch in `model.go`:

```go
case ScreenProjectStack:
    return []string{
        "auto (" + m.ProjectStack + ")",
        "─────────────",
        "node", "angular", "nextjs", "react", "vue",
        "python", "go", "rust", "java", "unknown",
    }

case ScreenProjectMemory:
    return []string{
        "🧠 Obsidian Brain — Markdown vault with AI context",
        "📋 VibeKanban — Kanban board in markdown",
        "💾 Engram — Lightweight AI memory store",
        "📝 Simple — Plain project notes",
        "🚫 None — No memory module",
    }

case ScreenProjectEngram:
    return []string{
        "Yes, also add Engram for AI memory",
        "No, Obsidian Brain is enough",
    }

case ScreenProjectCI:
    return []string{
        "⚙️  GitHub Actions",
        "🦊 GitLab CI",
        "🪵 Woodpecker CI",
        "🚫 None",
    }

case ScreenProjectConfirm:
    return []string{
        "✅ Confirm & Run",
        "❌ Cancel",
    }

case ScreenSkillMenu:
    return []string{
        "📖 Browse available skills",
        "➕ Install skills",
        "➖ Remove skills",
        "─────────────",
        "← Back",
    }

case ScreenSkillInstall:
    opts := make([]string, 0, len(m.SkillList)+2)
    opts = append(opts, m.SkillList...)
    opts = append(opts, "─────────────")
    opts = append(opts, "✅ Confirm selection")
    return opts

case ScreenSkillRemove:
    opts := make([]string, 0, len(m.InstalledSkills)+2)
    opts = append(opts, m.InstalledSkills...)
    opts = append(opts, "─────────────")
    opts = append(opts, "✅ Confirm removal")
    return opts
```

---

## GetScreenTitle / GetScreenDescription Additions

Add to `GetScreenTitle()` in `model.go`:

```go
case ScreenProjectPath:
    return "📦 Initialize Project — Step 1: Project Path"
case ScreenProjectStack:
    return "📦 Initialize Project — Step 2: Tech Stack"
case ScreenProjectMemory:
    return "📦 Initialize Project — Step 3: Memory Module"
case ScreenProjectEngram:
    return "📦 Initialize Project — Step 3b: Add Engram?"
case ScreenProjectCI:
    return "📦 Initialize Project — Step 4: CI Provider"
case ScreenProjectConfirm:
    return "📦 Initialize Project — Confirm"
case ScreenProjectInstalling:
    return "📦 Initializing Project..."
case ScreenProjectResult:
    if m.ErrorMsg != "" {
        return "📦 Initialization Failed"
    }
    return "📦 Project Initialized!"
case ScreenSkillMenu:
    return "🎯 Skill Manager"
case ScreenSkillBrowse:
    return "🎯 Skill Manager — Browse"
case ScreenSkillInstall:
    return "🎯 Skill Manager — Install"
case ScreenSkillRemove:
    return "🎯 Skill Manager — Remove"
case ScreenSkillResult:
    if m.ErrorMsg != "" {
        return "🎯 Skill Operation Failed"
    }
    return "🎯 Skill Operation Complete"
```

Add to `GetScreenDescription()` in `model.go`:

```go
case ScreenProjectPath:
    return "Enter an absolute path or ~ path to an existing directory"
case ScreenProjectStack:
    if m.ProjectStack != "" && m.ProjectStack != "unknown" {
        return fmt.Sprintf("Auto-detected: %s — confirm or select a different stack", m.ProjectStack)
    }
    return "Select the primary technology stack for this project"
case ScreenProjectMemory:
    return "Choose how AI agents will store and recall project context"
case ScreenProjectEngram:
    return "Engram provides lightweight key-value memory alongside Obsidian Brain"
case ScreenProjectCI:
    return "Choose your CI/CD provider for automated workflows"
case ScreenProjectConfirm:
    return "Review your choices before running init-project.sh"
case ScreenSkillMenu:
    return "Manage AI skills in this project via Gentleman-Skills"
case ScreenSkillBrowse:
    return "Read-only view of all available skills"
case ScreenSkillInstall:
    return "Toggle skills to install. Press Enter on a skill to select."
case ScreenSkillRemove:
    return "Toggle skills to remove. Press Enter on a skill to select."
case ScreenSkillResult:
    return "Operation complete"
```

---

## CLI Flags (main.go)

Add to `cliFlags` struct:

```go
// Project Init flags
initProject   bool
projectPath   string
projectStack  string
projectMemory string
projectCI     string
projectEngram bool

// Skill Manager flags
skillAction string // "install" or "remove"
skillNames  string // comma-separated skill names
```

Add to `parseFlags()`:

```go
flag.BoolVar(&flags.initProject, "init-project", false, "Initialize a project with AI framework")
flag.StringVar(&flags.projectPath, "project-path", "", "Path to the project directory (required with --init-project)")
flag.StringVar(&flags.projectStack, "project-stack", "", "Tech stack override (auto-detected if empty)")
flag.StringVar(&flags.projectMemory, "project-memory", "simple", "Memory module: obsidian-brain, vibekanban, engram, simple, none")
flag.StringVar(&flags.projectCI, "project-ci", "none", "CI provider: github, gitlab, woodpecker, none")
flag.BoolVar(&flags.projectEngram, "project-engram", false, "Add Engram alongside Obsidian Brain")
flag.StringVar(&flags.skillAction, "skill", "", "Skill action: install or remove")
flag.StringVar(&flags.skillNames, "skill-names", "", "Comma-separated skill names (used with --skill)")
```

Validation in `runNonInteractive`:

```go
// --- Project Init ---
if flags.initProject {
    if flags.projectPath == "" {
        return fmt.Errorf("--project-path is required with --init-project")
    }
    validMemory := map[string]bool{"obsidian-brain": true, "vibekanban": true, "engram": true, "simple": true, "none": true}
    if !validMemory[flags.projectMemory] {
        return fmt.Errorf("invalid --project-memory: %s", flags.projectMemory)
    }
    validCI := map[string]bool{"github": true, "gitlab": true, "woodpecker": true, "none": true}
    if !validCI[flags.projectCI] {
        return fmt.Errorf("invalid --project-ci: %s", flags.projectCI)
    }
    choices.InitProject = true
    choices.ProjectPath = flags.projectPath
    choices.ProjectStack = flags.projectStack
    choices.ProjectMemory = flags.projectMemory
    choices.ProjectCI = flags.projectCI
    choices.ProjectEngram = flags.projectEngram
}

// --- Skill Manager ---
if flags.skillAction != "" {
    validActions := map[string]bool{"install": true, "remove": true}
    if !validActions[flags.skillAction] {
        return fmt.Errorf("invalid --skill: %s (valid: install, remove)", flags.skillAction)
    }
    if flags.skillNames == "" {
        return fmt.Errorf("--skill-names is required with --skill")
    }
    var names []string
    for _, n := range strings.Split(flags.skillNames, ",") {
        n = strings.TrimSpace(n)
        if n != "" {
            names = append(names, n)
        }
    }
    choices.SkillAction = flags.skillAction
    choices.SkillNames = names
}
```

`RunNonInteractive` dispatches these via separate functions `runProjectInit(choices)` and `runSkillAction(choices)` that create a minimal `Model` and call `stepInitProject` / `stepSkillInstall` / `stepSkillRemove` directly, following the `buildStepsForChoices` + `executeStep` pattern in `non_interactive.go`.

---

## Text Input Component

The project does not use `charmbracelet/bubbles` (not in `go.mod`). The path input screen uses the **same manual key-accumulation pattern** as `TrainerInput` in `handleTrainerExerciseKeys`. No new dependency is needed.

### Initialization

The field `ProjectPathInput string` on `Model` is initialized to `""` when entering `ScreenProjectPath` from the main menu dispatch:

```go
case strings.Contains(selected, "Initialize Project"):
    m.Screen = ScreenProjectPath
    m.ProjectPathInput = ""  // reset on each entry
    m.ProjectPathError = ""
    m.Cursor = 0
```

### Update delegation

In `handleKeyPress`, `ScreenProjectPath` is dispatched to `handleProjectPathKeys`. The space key guard in `handleKeyPress` must include `ScreenProjectPath` to prevent leader mode activation:

```go
// In the space-key block of handleKeyPress:
case ScreenProjectPath:
    // space is a valid path character (e.g., "/Users/foo/my project")
    // fall through — do NOT activate leader mode
    // (handleProjectPathKeys appends " " for len(key)==1 chars)
```

Actual key handling in `handleProjectPathKeys`:

```go
func (m Model) handleProjectPathKeys(key string) (tea.Model, tea.Cmd) {
    switch key {
    case "esc":
        m.Screen = ScreenMainMenu
        m.ProjectPathInput = ""
        m.ProjectPathError = ""
        m.Cursor = 0
        return m, nil

    case "backspace":
        // Handle multi-byte UTF-8 correctly using rune slices
        runes := []rune(m.ProjectPathInput)
        if len(runes) > 0 {
            m.ProjectPathInput = string(runes[:len(runes)-1])
        }
        m.ProjectPathError = "" // clear error on edit
        return m, nil

    case "enter":
        path := expandPath(strings.TrimSpace(m.ProjectPathInput))
        if path == "" {
            m.ProjectPathError = "Path cannot be empty"
            return m, nil
        }
        info, err := os.Stat(path)
        if err != nil {
            m.ProjectPathError = fmt.Sprintf("Path does not exist: %s", path)
            return m, nil
        }
        if !info.IsDir() {
            m.ProjectPathError = "Path must be a directory, not a file"
            return m, nil
        }
        absPath, _ := filepath.Abs(path)
        m.ProjectPathInput = absPath
        m.ProjectPathError = ""
        m.ProjectStack = detectStack(absPath)
        m.Screen = ScreenProjectStack
        m.Cursor = 0
        return m, nil

    default:
        // Append single printable characters (same pattern as TrainerInput)
        if len(key) == 1 {
            m.ProjectPathInput += key
            m.ProjectPathError = "" // clear error on edit
        }
        return m, nil
    }
}
```

### View rendering

The cursor `█` is rendered inline — no blinking animation required (the tick is already running but we use a static block cursor for simplicity):

```go
// In renderProjectPath():
inputLine := "> " + m.ProjectPathInput + "█"
s.WriteString(BoxStyle.Render(inputLine))
```

---

## Skill List Loading

### Async pattern

When the user selects Browse, Install, or Remove from `ScreenSkillMenu`, the handler immediately sets `m.SkillLoading = true` and returns a `loadSkillsCmd()` tea.Cmd that runs in the background:

```go
func loadSkillsCmd() tea.Cmd {
    return func() tea.Msg {
        available, installed, err := fetchSkillCatalog()
        return skillsLoadedMsg{
            available: available,
            installed: installed,
            err:       err,
        }
    }
}
```

### `fetchSkillCatalog()` implementation

```go
const gentlemanSkillsRepo = "https://github.com/Gentleman-Programming/Gentleman-Skills.git"
const skillsCachePath = "/tmp/gentleman-skills-cache"

func fetchSkillCatalog() (available []string, installed []string, err error) {
    // Use cached clone if < 1 hour old
    info, statErr := os.Stat(skillsCachePath)
    needsClone := statErr != nil || time.Since(info.ModTime()) > time.Hour

    if needsClone {
        system.Run("rm -rf "+skillsCachePath, nil)
        result := system.Run(
            "git clone --depth 1 "+gentlemanSkillsRepo+" "+skillsCachePath,
            nil,
        )
        if result.Error != nil {
            return nil, nil, fmt.Errorf("failed to fetch skill catalog: %w", result.Error)
        }
    }

    // Parse available skills: list directories in skillsCachePath/skills/
    skillsDir := filepath.Join(skillsCachePath, "skills")
    entries, err := os.ReadDir(skillsDir)
    if err != nil {
        return nil, nil, fmt.Errorf("could not read skills directory: %w", err)
    }
    for _, e := range entries {
        if e.IsDir() && !strings.HasPrefix(e.Name(), ".") {
            available = append(available, e.Name())
        }
    }
    sort.Strings(available)

    // Parse installed skills: check current working directory for .claude/skills/ or similar
    // The exact detection depends on init-project.sh conventions
    cwd, _ := os.Getwd()
    installedDir := filepath.Join(cwd, ".claude", "skills")
    if entries2, err2 := os.ReadDir(installedDir); err2 == nil {
        for _, e := range entries2 {
            if e.IsDir() && !strings.HasPrefix(e.Name(), ".") {
                installed = append(installed, e.Name())
            }
        }
        sort.Strings(installed)
    }

    return available, installed, nil
}
```

### Model population

When `skillsLoadedMsg` arrives in `Update()`:

```go
case skillsLoadedMsg:
    m.SkillLoading = false
    if msg.err != nil {
        m.SkillLoadError = msg.err.Error()
        m.SkillList = []string{}
        m.InstalledSkills = []string{}
    } else {
        m.SkillLoadError = ""
        m.SkillList = msg.available
        m.InstalledSkills = msg.installed
        m.SkillSelected = make([]bool, len(msg.available))
    }
    return m, nil
```

The view functions check `m.SkillLoading` before rendering the list — showing a spinner message while loading and the full list once `SkillLoading == false`. The spinner is driven by the existing `tickMsg` + `SpinnerFrame` mechanism.

---

## Test Strategy

### `project_screens_test.go` — new file in `installer/internal/tui/`

Key test cases:

| Test | What it verifies |
|------|-----------------|
| `TestProjectPathValidation` | Empty input shows error; non-existent path shows error; file (not dir) shows error; valid dir advances to `ScreenProjectStack` |
| `TestProjectPathTildeExpansion` | `~/projects` expands to `$HOME/projects` |
| `TestProjectPathBackspace` | Backspace removes last rune (UTF-8 safe) |
| `TestProjectStackNavigation` | Up/down moves cursor; enter saves stack and advances to `ScreenProjectMemory` |
| `TestProjectMemoryConditionalEngram` | Selecting "obsidian-brain" routes to `ScreenProjectEngram`; others skip to `ScreenProjectCI` |
| `TestProjectCIAdvancesToConfirm` | Selecting any CI option routes directly to `ScreenProjectConfirm` |
| `TestProjectEscapeBackNavigation` | Esc from each screen returns to the correct prior screen |
| `TestProjectConfirmCancel` | Cursor=1 on confirm returns to `ScreenMainMenu` |
| `TestGetCurrentOptionsProjectScreens` | Each new screen returns correct option count and labels |
| `TestGetScreenTitleProjectScreens` | Each screen returns non-empty title |
| `TestGetScreenDescriptionProjectStack` | Shows "Auto-detected:" prefix when `m.ProjectStack != ""` |

### `skill_screens_test.go` — new file in `installer/internal/tui/`

Key test cases:

| Test | What it verifies |
|------|-----------------|
| `TestSkillMenuNavigation` | Browse/Install/Remove options exist and dispatch correctly |
| `TestSkillMenuEscape` | Esc from `ScreenSkillMenu` returns to `ScreenMainMenu` |
| `TestSkillBrowseScrolling` | Up/down adjusts `SkillScroll`; esc resets scroll and returns to menu |
| `TestSkillInstallToggle` | Enter toggles `SkillSelected[cursor]`; confirm with nothing selected is a no-op |
| `TestSkillInstallConfirm` | After selecting skills and confirming, screen transitions to `ScreenSkillResult` |
| `TestSkillRemoveEmpty` | `ScreenSkillRemove` with empty `InstalledSkills` shows 0 items (no panic) |
| `TestSkillsLoadedMsg` | Delivering `skillsLoadedMsg` populates `SkillList`, clears `SkillLoading` |
| `TestSkillsLoadedMsgError` | Error in `skillsLoadedMsg` sets `SkillLoadError`, keeps list empty |
| `TestSkillResultEnterReturnsToMenu` | Enter on `ScreenSkillResult` navigates to `ScreenSkillMenu` |
| `TestGetCurrentOptionsSkillScreens` | `ScreenSkillInstall` options = `len(SkillList) + 2` (separator + confirm) |
| `TestRenderSkillInstallLoading` | While `SkillLoading == true`, view shows loading message (no panic on empty list) |

All tests follow the existing pattern in `update_test.go` — create a model with `NewModel()`, set fields directly, send a `tea.KeyMsg`, and assert the resulting screen and state. No `teatest` framework needed for unit-level handler tests.
