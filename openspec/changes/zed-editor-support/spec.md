# Spec: zed-editor-support

## Requirements

### REQ-ZED-01: TUI Screen

A new screen `ScreenZedSelect` must be presented to the user during the installation wizard, positioned between `ScreenNvimSelect` (Step 6) and `ScreenAIToolsSelect` (Step 7, renumbered to Step 8).

- The screen title must read: `"Step 7: Zed Editor"`
- The screen description must read: `"High-performance editor with Vim mode and AI agent support"`
- The screen must offer exactly two selectable options:
  1. `"Yes, install Zed with config"` (index 0)
  2. `"No, skip Zed"` (index 1)
- Selecting index 0 sets `UserChoices.InstallZed = true`
- Selecting index 1 sets `UserChoices.InstallZed = false`
- After selection, the user proceeds to `ScreenAIToolsSelect`
- The screen must follow the existing `renderSelection()` pattern (same as `ScreenNvimSelect`, `ScreenFontSelect`)
- Navigation: up/k, down/j, Enter to select, Esc to go back to `ScreenNvimSelect`

### REQ-ZED-02: Installation

Zed must be installed via platform-specific package managers:

| Platform | Command | Notes |
|----------|---------|-------|
| macOS | `brew install --cask zed` | Via Homebrew cask |
| Arch Linux | `pacman -S --noconfirm zed` | Official repos or AUR |
| Debian/Ubuntu | `curl -f https://zed.dev/install.sh \| sh` | Official install script |
| Fedora | `dnf install -y zed` | Via Copr or official repo |

- If installation fails, log a warning with `SendLog()` including the manual install URL (`https://zed.dev/download`) but do NOT abort the entire installation
- The step ID must be `"zed"`
- The step must NOT be marked as `Interactive` (no sudo required beyond what package managers handle)

### REQ-ZED-03: Configuration

Two config files must be created under `GentlemanZed/zed/` in the repository and copied to `~/.config/zed/` during installation:

1. **`settings.json`** containing:
   - `vim_mode: true`
   - Theme: Kanagawa (matching the Gentleman aesthetic)
   - Buffer font family: `"Iosevka Term"` (same Nerd Font used across all configs)
   - Buffer font size: `14`
   - Relative line numbers enabled
   - Claude agent configuration via ACP (Agent Control Protocol)
   - Tab size: 2, soft tabs enabled
   - Format on save enabled
   - Telemetry disabled

2. **`keymap.json`** containing:
   - Vim-mode keybindings consistent with the Gentleman Neovim workflow
   - Window split navigation (`ctrl-h/j/k/l`)
   - Leader key bindings for common operations (file finder, buffer close, etc.)

- Config copy must use `system.CopyDir()` (same pattern as Neovim: `GentlemanNvim/nvim` -> `~/.config/nvim`)
- Source: `Gentleman.Dots/GentlemanZed/zed/`
- Destination: `~/.config/zed/`
- The destination directory must be created with `system.EnsureDir()` if it does not exist

### REQ-ZED-04: Non-Interactive Mode

- A `--zed` boolean CLI flag must be added to `cliFlags` struct
- When `--non-interactive` is used with `--zed`, the Zed installation step must execute
- The flag must default to `false`
- The summary output must include a `Zed:` line showing `true` or `false`
- The help text must document the `--zed` flag under "Non-Interactive Options"

### REQ-ZED-05: Termux Exclusion

- When the system is detected as Termux (`m.SystemInfo.IsTermux == true`), `ScreenZedSelect` must be skipped entirely
- The forward navigation from `ScreenNvimSelect` on Termux must bypass `ScreenZedSelect` and go directly to `proceedToBackupOrInstall()` (existing behavior, unchanged)
- The `stepInstallZed()` function must early-return with a skip log if called on Termux
- Rationale: Zed requires a GUI with Vulkan GPU rendering, which Termux does not provide

## Scenarios

### SC-01: User selects Zed on macOS

1. User reaches `ScreenNvimSelect`, makes a choice (yes or no for Nvim)
2. TUI advances to `ScreenZedSelect` (Step 7)
3. Progress bar shows: `... WM -> Nvim -> [Zed] -> AI Tools -> Framework`
4. User sees title "Step 7: Zed Editor" and two options
5. User selects "Yes, install Zed with config"
6. `UserChoices.InstallZed` is set to `true`
7. TUI advances to `ScreenAIToolsSelect` (Step 8)
8. During installation, `stepInstallZed()` runs:
   - Executes `brew install --cask zed`
   - Creates `~/.config/zed/` if missing
   - Copies `GentlemanZed/zed/settings.json` and `keymap.json` to `~/.config/zed/`
   - Logs "Zed configured with Gentleman setup"

### SC-02: User selects Zed on Arch Linux

1. User reaches `ScreenZedSelect` after `ScreenNvimSelect`
2. User selects "Yes, install Zed with config"
3. During installation, `stepInstallZed()` runs:
   - Executes `pacman -S --noconfirm zed` via `system.RunSudoWithLogs()`
   - Creates `~/.config/zed/` if missing
   - Copies config files
   - Logs success

### SC-03: User skips Zed

1. User reaches `ScreenZedSelect`
2. User selects "No, skip Zed" (index 1)
3. `UserChoices.InstallZed` is set to `false`
4. TUI advances to `ScreenAIToolsSelect`
5. During installation, no `stepInstallZed()` step is registered in `SetupInstallSteps()`
6. No Zed-related files are touched

### SC-04: Non-interactive with --zed flag

1. User runs: `gentleman.dots --non-interactive --shell=fish --zed`
2. `parseFlags()` sets `flags.zed = true`
3. `runNonInteractive()` creates `UserChoices{InstallZed: true}`
4. Summary prints `Zed: true`
5. `buildStepsForChoices()` includes `{ID: "zed", Name: "Install Zed editor"}`
6. `stepInstallZed()` executes platform-specific installation + config copy

### SC-05: Termux user

1. User selects Termux as OS in `ScreenOSSelect`
2. User proceeds through the wizard to `ScreenNvimSelect`
3. After making Nvim choice, TUI skips `ScreenZedSelect` entirely
4. TUI goes directly to `proceedToBackupOrInstall()` (same as current behavior)
5. `ScreenZedSelect` never appears
6. `UserChoices.InstallZed` remains `false`
7. No Zed step is registered or executed

### SC-06: User navigates back from Zed screen

1. User is on `ScreenZedSelect`
2. User presses Esc or Backspace
3. TUI returns to `ScreenNvimSelect`
4. `UserChoices.InstallZed` is reset to `false`
5. Cursor resets to 0

### SC-07: User navigates back from AI Tools to Zed

1. User is on `ScreenAIToolsSelect`
2. User presses Esc
3. TUI returns to `ScreenZedSelect` (NOT `ScreenNvimSelect`)
4. `UserChoices.AITools` is reset to nil
5. `AIToolSelected` is reset to nil

### SC-08: Backup system detects existing Zed config

1. User has existing `~/.config/zed/` directory
2. During backup detection, `system.ConfigPaths()` includes `"zed": home + "/.config/zed"`
3. Backup system captures the existing directory before overwrite
4. If user later restores backup, Zed config is restored from backup
