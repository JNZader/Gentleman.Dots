# Proposal: zed-editor-support

## Intent

The Javi.Dots TUI installer currently supports Neovim as the sole editor option (Step 6). Zed is a high-performance GUI editor with native Vim mode, built-in AI agent support via ACP (Agent Control Protocol), and a growing user base among developers who want a modern IDE alongside a terminal-native workflow.

We want to add Zed as an **optional complement** to Neovim — not a replacement. Users who select Neovim can also opt into Zed, getting a consistent Gentleman-branded experience (Kanagawa theme, vim keybindings, Claude agent integration) across both editors. Users who skip Neovim can still install Zed standalone.

## Scope

### In Scope

1. **New TUI screen `ScreenZedSelect`** — Simple yes/no selection (same pattern as `ScreenNvimSelect`), inserted between `ScreenNvimSelect` and `ScreenAIToolsSelect` in the installation flow
2. **New config directory `GentlemanZed/zed/`** — Contains `settings.json` and `keymap.json` with:
   - `vim_mode: true` enabled by default
   - Kanagawa theme (consistent with the Gentleman aesthetic)
   - Claude agent integration via ACP
   - Sensible defaults matching the project's opinionated philosophy
3. **Installation step `stepInstallZed()`** — Installs Zed binary + copies config to `~/.config/zed/`
4. **Progress bar update** — Add "Zed" to the hardcoded `renderStepProgress()` steps array (currently: OS, Terminal, Font, Shell, WM, Nvim, AI Tools, Framework)
5. **Non-interactive CLI flag** — `--zed` boolean flag for headless installation
6. **Platform-specific installation**:
   - macOS: `brew install --cask zed`
   - Arch: `pacman -S zed` (or AUR)
   - Debian/Ubuntu: official curl install script
   - Fedora: `dnf install zed` (or Copr)
7. **Termux exclusion** — Zed requires a GUI with Vulkan support; the screen must be skipped entirely on Termux (same pattern used for AI tools)

### Out of Scope

- Zed plugin/extension management (users install extensions via Zed's built-in marketplace)
- Hooks support via ACP (documented limitation — ACP does not support hooks as of 2026)
- Zed as a replacement for Neovim in any flow
- Zed keymaps reference screen in Learn & Practice (can be added later)
- Zed Vim Trainer integration
- Windows/WSL support for Zed (Zed's Linux support covers native Linux only)

## Approach

**Mirror the Neovim pattern exactly.** The Neovim integration touches 7 files with a clear, repeatable pattern. Zed follows the same structure:

### Screen Flow

```
ScreenWMSelect
  → ScreenNvimSelect          (existing — Step 6)
    → ScreenZedSelect          (NEW — Step 7: "Yes, install Zed with config" / "No, skip Zed")
      → ScreenAIToolsSelect    (existing — renumbered to Step 8)
        → ScreenAIFrameworkConfirm (existing — renumbered to Step 9)
```

### Config Copy

```
GentlemanZed/zed/settings.json  →  ~/.config/zed/settings.json
GentlemanZed/zed/keymap.json    →  ~/.config/zed/keymap.json
```

### Installation Logic

Same pattern as `stepInstallNvim()`:
1. Install Zed binary via platform-specific package manager
2. Create `~/.config/zed/` if it doesn't exist
3. Copy config files from cloned repo's `GentlemanZed/zed/` directory
4. Log each operation via `SendLog()`

### Progress Bar

Update `renderStepProgress()` in `view.go`:
```go
// Before:
steps := []string{"OS", "Terminal", "Font", "Shell", "WM", "Nvim", "AI Tools", "Framework"}

// After:
steps := []string{"OS", "Terminal", "Font", "Shell", "WM", "Nvim", "Zed", "AI Tools", "Framework"}
```

All subsequent screen-to-index mappings shift by 1 for AI Tools (6→7) and Framework (7→8).

## Affected Areas

| Area | Impact | Description |
|------|--------|-------------|
| `installer/internal/tui/model.go` | Modified | Add `ScreenZedSelect` constant (between `ScreenNvimSelect` and `ScreenInstalling`), add `InstallZed bool` to `UserChoices`, add entries in `GetCurrentOptions`, `GetScreenTitle`, `GetScreenDescription` |
| `installer/internal/tui/update.go` | Modified | Handle `ScreenZedSelect` enter key (set `InstallZed`), update `ScreenNvimSelect` to navigate to `ScreenZedSelect` instead of `ScreenAIToolsSelect`, skip Zed screen on Termux, add to Esc/back handler, add to cursor clamp list |
| `installer/internal/tui/view.go` | Modified | Add `ScreenZedSelect` to the single-select render case, update `renderStepProgress()` steps array and index mapping |
| `installer/internal/tui/installer.go` | Modified | Add `stepInstallZed()` function, add Zed step to `SetupInstallSteps()` (between nvim and aitools steps) |
| `installer/cmd/gentleman-installer/main.go` | Modified | Add `--zed` flag to `cliFlags` struct, `parseFlags()`, `runNonInteractive()` choices, `printHelp()`, and summary output |
| `installer/internal/tui/non_interactive.go` | Modified | Wire `InstallZed` in `RunNonInteractive()` to call `stepInstallZed()` |
| `GentlemanZed/zed/settings.json` | New | Zed settings: vim_mode, Kanagawa theme, font, Claude ACP agent config |
| `GentlemanZed/zed/keymap.json` | New | Custom Vim-mode keybindings for Zed |

## Key Decisions

| Decision | Rationale |
|----------|-----------|
| **Complement, not replacement** | Neovim remains the primary editor. Zed is additive for users who want a GUI option alongside their terminal workflow. |
| **Simple yes/no, not multi-select** | Matches the Neovim pattern. One editor = one boolean. No complexity. |
| **vim_mode enabled by default** | The target audience (Gentleman.Dots users) is already invested in Vim keybindings. Consistency across editors. |
| **Kanagawa theme** | Matches the Gentleman aesthetic used across all configs (terminals, Neovim, shells). |
| **Screen placed after Nvim, before AI Tools** | Logical grouping: editors together, then AI tools. Keeps the "editor block" contiguous in the flow. |
| **Skip on Termux entirely** | Zed requires Vulkan GPU rendering. Termux runs on Android without desktop GPU access. Same skip pattern as AI tools. |
| **ACP for Claude integration, no hooks** | ACP is Zed's protocol for AI agents. Hooks are NOT supported via ACP (documented limitation). Config will include ACP setup but explicitly note this gap. |

## Risks & Mitigations

| Risk | Likelihood | Impact | Mitigation |
|------|------------|--------|------------|
| Zed not available on all Linux distros | Medium | Medium | Provide fallback curl install script; log clear error with manual install URL if all methods fail |
| Zed config format changes between versions | Low | Medium | Pin to stable config schema; settings.json and keymap.json are simple JSON, rarely breaking |
| Progress bar gets crowded with 9 steps | Low | Low | Steps already render as a single horizontal line; one more item fits within 80-col terminals. Monitor and consider abbreviating if more steps are added later |
| Step renumbering breaks existing tests | Medium | Low | Tests reference `Screen*` constants (not step numbers). The iota values change but all references are symbolic. Search-and-verify all `ScreenAIToolsSelect` assertions in tests |
| Zed binary is large (~100MB) | Low | Low | Same as any brew cask install. User already expects downloads during installation. Show progress via `SendLog` |
| Config overwrite conflicts | Low | Low | Existing backup system handles this — if `~/.config/zed/` exists, backup step captures it before overwrite |

## Dependencies

- **Zed binary availability**: Zed must be installable via `brew install --cask zed` (macOS), package manager (Linux), or official install script
- **Vulkan-capable GPU**: Required for Zed rendering (excludes Termux, some headless servers)
- **Existing patterns**: Relies on `system.RunBrewWithLogs`, `system.RunSudoWithLogs`, `system.RunWithLogs`, `system.CopyDir` — all already used by the Neovim step
- **Cloned repo**: Zed config files live in the repo under `GentlemanZed/zed/`, copied during the clone step (same as `GentlemanNvim/nvim/`)

## Estimated Effort

| Component | Estimate |
|-----------|----------|
| `GentlemanZed/zed/` config files | ~1 hour (settings.json + keymap.json with Kanagawa + vim_mode + ACP) |
| `model.go` changes | ~30 min (screen constant, UserChoices field, 3 switch cases) |
| `update.go` changes | ~30 min (enter handler, navigation rewiring, esc handler, cursor clamp) |
| `view.go` changes | ~15 min (render case, progress bar update) |
| `installer.go` changes | ~1 hour (stepInstallZed with 4 platform variants + config copy) |
| `main.go` + `non_interactive.go` | ~30 min (flag, validation, wiring) |
| Test updates | ~1 hour (fix broken assertions from screen reordering, add Zed-specific tests) |
| **Total** | **~5 hours** (~300-400 lines of Go, following existing patterns) |

## Success Criteria

- [ ] `ScreenZedSelect` appears as Step 7 between Neovim and AI Tools in the TUI flow
- [ ] Selecting "Yes" installs Zed binary and copies config to `~/.config/zed/`
- [ ] Selecting "No" skips Zed entirely (no residual config)
- [ ] Screen is automatically skipped on Termux
- [ ] Progress bar shows 9 steps: OS → Terminal → Font → Shell → WM → Nvim → Zed → AI Tools → Framework
- [ ] `--zed` flag works in non-interactive mode
- [ ] `GentlemanZed/zed/settings.json` enables vim_mode, Kanagawa theme, and Claude ACP agent
- [ ] Backup system captures existing `~/.config/zed/` before overwrite
- [ ] `go test ./...` passes with zero failures
- [ ] Installation works on macOS (brew cask), Arch (pacman), Debian/Ubuntu (curl script), and Fedora (dnf)
