# Proposal: Project Init & Skill Manager TUI Flows

## Intent

The Javi.Dots TUI installer currently only handles global environment setup (one-time, fresh PC). The `project-starter-framework` has two powerful scripts — `init-project.sh` (per-project setup) and `add-skill.sh` (skill marketplace) — that are NOT accessible through the TUI. Users have to know about these scripts and run them manually.

We want to expose both as first-class TUI flows in the Main Menu, making the installer a **recurring tool** (not just a one-time setup), with the same granular drill-down UX used in the AI Framework flow.

## Scope

### In Scope

1. **"Initialize Project" main menu item** — TUI flow that collects:
   - Project path (user input with validation)
   - Stack detection (auto-detected, shown for confirmation)
   - Memory module selection (Obsidian Brain, VibeKanban, Engram, Simple, None)
   - CI provider selection (GitHub Actions, GitLab CI, Woodpecker, None)
   - Optional add-on: Engram for AI memory (if Obsidian Brain selected)
   - Executes `init-project.sh --non-interactive --memory=N --ci=N [--engram]`

2. **"Skill Manager" main menu item** — TUI flow with 3 sub-actions:
   - **Browse**: list available skills from Gentleman-Skills repo (cached clone)
   - **Install**: select skill(s) from the list → `add-skill.sh gentleman <name>`
   - **Remove**: list installed skills → select → `add-skill.sh remove <name>`

3. **New screen constants and handlers** for both flows
4. **Tests** for new screens, navigation, and edge cases

### Out of Scope

- Modifying `init-project.sh` or `add-skill.sh` themselves (we consume them as-is)
- `collect-skills.sh` integration (imports from local filesystem — different workflow)
- `sync-ai-config.sh` / `sync-skills.sh` / `doctor.sh` (maintenance tools, not user-facing setup)
- Per-project CI template customization beyond provider selection
- Codex CLI as a 5th AI tool option

## Approach

**Hybrid pattern** (same as AI Framework integration): TUI screens collect user choices with drill-down navigation, then delegate execution to the framework scripts with `--non-interactive` flags. The scripts handle all file operations, git config, and template copying.

### Initialize Project Flow

```
ScreenMainMenu
  → "📦 Initialize Project"
    → ScreenProjectPath         (text input: project directory path)
    → ScreenProjectStack        (auto-detected stack, confirm or override)
    → ScreenProjectMemory       (single-select: obsidian-brain, vibekanban, engram, simple, none)
    → [if obsidian-brain] ScreenProjectEngram  (yes/no: add Engram too?)
    → ScreenProjectCI           (single-select: github, gitlab, woodpecker, none)
    → ScreenProjectConfirm      (summary → execute)
    → ScreenProjectInstalling   (progress with log output)
    → ScreenProjectComplete     (done)
```

### Skill Manager Flow

```
ScreenMainMenu
  → "🎯 Skill Manager"
    → ScreenSkillMenu           (Browse / Install / Remove)
    → ScreenSkillBrowse         (scrollable list from Gentleman-Skills repo)
    → ScreenSkillInstall        (multi-select from available, confirm, execute)
    → ScreenSkillRemove         (multi-select from installed, confirm, execute)
    → ScreenSkillResult         (success/error output)
```

### Execution

- `init-project.sh` is cloned from `project-starter-framework` to `/tmp/` (same pattern as AI Framework), then run with flags derived from user choices
- `add-skill.sh` is run from the same cloned repo — the clone is already done if AI Framework was installed, otherwise we clone on demand
- Both flows return to Main Menu on completion or error

## Affected Areas

| Area | Impact | Description |
|------|--------|-------------|
| `installer/internal/tui/model.go` | Modified | ~10 new screen constants, new Model fields (ProjectPath, ProjectStack, ProjectMemory, ProjectCI, SkillList, InstalledSkills), new entries in GetCurrentOptions/GetScreenTitle/GetScreenDescription |
| `installer/internal/tui/update.go` | Modified | New handler functions for each screen, main menu dispatch for 2 new items, text input handling for path screen, skill list loading via script execution |
| `installer/internal/tui/view.go` | Modified | New render functions for each screen type (text input, single-select, multi-select with scroll, progress, result) |
| `installer/internal/tui/installer.go` | Modified | New step functions: `stepInitProject()` (clones framework, runs init-project.sh), `stepSkillInstall()` / `stepSkillRemove()` (runs add-skill.sh) |
| `installer/internal/tui/project_screens_test.go` | New | Tests for project init flow (path validation, stack options, memory options, CI options, conditional screens) |
| `installer/internal/tui/skill_screens_test.go` | New | Tests for skill manager flow (browse, install, remove, empty states) |
| `installer/cmd/gentleman-installer/main.go` | Modified | New CLI flags for non-interactive project init (`--init-project`, `--project-path`, `--stack`, `--memory`, `--ci`, `--engram`) |

## Risks

| Risk | Likelihood | Mitigation |
|------|------------|------------|
| Path input validation edge cases (spaces, ~, symlinks) | Medium | Expand `~` to `$HOME`, `filepath.Abs()`, check `os.Stat()` exists and is dir |
| `add-skill.sh` clone takes time on slow connections | Medium | Show spinner with "Fetching skill catalog...", reuse cached clone if <1h old |
| `init-project.sh` requires git repo at target path | Low | Auto-detect, offer to `git init` if not a repo (script already handles this) |
| Project path outside current filesystem access | Low | Validate with `os.Stat()` before proceeding |
| Screen count explosion (currently 41, adding ~11 more) | Low | Follow existing patterns, no architectural change needed — iota handles it |

## Rollback Plan

All changes are additive — new screens, new handlers, new menu items. No existing functionality is modified. Rollback = revert the commit(s). The Main Menu dispatch uses `strings.Contains()` matching, so removing items doesn't break anything.

## Dependencies

- `project-starter-framework` repo must be accessible at `https://github.com/JNZader/project-starter-framework.git` (already used by AI Framework flow)
- `init-project.sh` expects `bash` available (true for all supported platforms except pure Windows)
- `add-skill.sh` expects `git` for cloning Gentleman-Skills (already a prerequisite of the installer)

## Success Criteria

- [ ] "Initialize Project" appears in Main Menu and completes a full flow (path → stack → memory → CI → execute)
- [ ] "Skill Manager" appears in Main Menu with Browse/Install/Remove sub-actions
- [ ] Project init runs `init-project.sh` with correct `--non-interactive` flags
- [ ] Skill install/remove runs `add-skill.sh` with correct subcommands
- [ ] Path input validates correctly (existing dir, expandable ~, absolute paths)
- [ ] Conditional screen works (Engram prompt only after Obsidian Brain)
- [ ] All new screens have tests
- [ ] `go test ./...` passes with zero failures
- [ ] Non-interactive CLI flags work for both flows
