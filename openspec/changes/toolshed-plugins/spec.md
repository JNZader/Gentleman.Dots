# Spec: toolshed-plugins

## Requirements

### REQ-PLUG-01: Plugin Directory Structure

**Description:** Create a `GentlemanClaude/plugins/` directory as a first-class category alongside `skills/`. Each plugin lives in its own subdirectory with a `PLUGIN.md` entry point.

**Acceptance Criteria:**

1. `GentlemanClaude/plugins/` directory exists at the repo root level
2. Each plugin has its own subdirectory: `merge-checks/`, `trim-md/`, `mermaid/`
3. Each plugin directory contains a `PLUGIN.md` entry point
4. Plugin subdirectories can contain: `scripts/`, `checks/`, `agents/`, `specialists/`, `references/`, `examples/`, `assets/`, `config/`, `templates/`
5. Plugins are NOT mixed into `skills/` -- the separation is explicit and enforced by directory structure
6. The installer copies plugins to `~/.claude/plugins/<name>/` (NOT `~/.claude/skills/`)

---

### REQ-PLUG-02: PLUGIN.md Format

**Description:** Define the `PLUGIN.md` frontmatter format. It must be recognizable by the TUI parser but clearly distinct from `SKILL.md`.

**Acceptance Criteria:**

1. Frontmatter uses YAML between `---` delimiters (same as SKILL.md)
2. Required fields: `name`, `description`, `type: plugin`
3. Optional fields: `dependencies` (list of runtime deps like `git`, `node`, `npm`), `permissions` (list of `settings.json` permission entries needed)
4. Fields from toolshed that are DROPPED: `allowed-tools` (we use `settings.json`), `argument-hint` (not supported), `model` (not supported), `user-invocable` (plugins are always invocable)
5. Body follows the same markdown conventions as SKILL.md (headings, code blocks, tables)
6. `type: plugin` field distinguishes plugins from skills for the TUI parser

**Example frontmatter:**

```yaml
---
name: merge-checks
description: Audit code changes across 13 quality dimensions before or after merge
type: plugin
dependencies:
  - git
  - bash >= 4.0
permissions:
  - "Bash(~/.claude/plugins/merge-checks/scripts/*:*)"
---
```

---

### REQ-PLUG-03: merge-checks Adaptation

**Description:** Port the merge-checks plugin (22 files, ~90KB) from toolshed format to Gentleman plugin format. Requires only `git` + `bash`.

**Acceptance Criteria:**

1. All 18 shell scripts are copied to `GentlemanClaude/plugins/merge-checks/scripts/`
2. All 3 check definition files are copied to `GentlemanClaude/plugins/merge-checks/checks/`
3. `SKILL.md` is renamed to `PLUGIN.md` with adapted frontmatter (drop `allowed-tools`, `argument-hint`, `model`; add `type: plugin`, `dependencies`, `permissions`)
4. ALL `find $HOME/.claude/plugins/cache` patterns in PLUGIN.md are rewritten to static `~/.claude/plugins/merge-checks/` paths
5. Script cross-references (`$SCRIPTS`, `$(dirname "$0")`) are verified to work from the installed location `~/.claude/plugins/merge-checks/scripts/`
6. `gather-context.sh` uses static path: `bash ~/.claude/plugins/merge-checks/scripts/gather-context.sh`
7. `precompute.sh` uses static path: `bash ~/.claude/plugins/merge-checks/scripts/precompute.sh`
8. Agent instruction file paths (`checks/docs.md`, `checks/comments.md`, `checks/shared.md`) use relative paths from `~/.claude/plugins/merge-checks/`
9. `ensure-deps.sh` is preserved as-is (already checks for `git`)
10. Plugin works with only `git` and `bash` installed (no Node.js required)
11. Assets directory (`assets/`) with SVG/MMD workflow diagrams is included

---

### REQ-PLUG-04: trim-md Adaptation

**Description:** Port the trim-md plugin (5 core files, ~8KB) from toolshed format. Requires Node.js for `markdownlint-cli2`. The PostToolUse hook is DROPPED.

**Acceptance Criteria:**

1. `SKILL.md` is renamed to `PLUGIN.md` with adapted frontmatter
2. `scripts/trim-md.sh` is copied to `GentlemanClaude/plugins/trim-md/scripts/`
3. `scripts/ensure-deps.sh` is copied to `GentlemanClaude/plugins/trim-md/scripts/`
4. `reference/` directory with markdownlint configs is copied to `GentlemanClaude/plugins/trim-md/config/`
5. ALL `find` path resolution patterns are rewritten to static `~/.claude/plugins/trim-md/` paths
6. The `SKILL_DIR` resolution block in PLUGIN.md is replaced with: `PLUGIN_DIR="$HOME/.claude/plugins/trim-md"`
7. `trim-md.sh` internal reference to `$SCRIPT_DIR/../reference` is rewritten to `$SCRIPT_DIR/../config`
8. `hooks/hooks.json` and `hooks/posttool-trim-md.sh` are NOT ported (GentlemanClaude has no hook system)
9. `tests/` directory from upstream is NOT ported (tests are for CI, not for user install)
10. If `node`/`npx` is not available, `ensure-deps.sh` exits with clear error message and install instructions
11. No global npm pollution -- `markdownlint-cli2` is auto-downloaded by `npx --yes` on first use

---

### REQ-PLUG-05: mermaid Adaptation

**Description:** Port the mermaid plugin (50+ files, ~350KB) from toolshed format. Requires Node.js for `beautiful-mermaid`. Has 6 sub-skills (mermaid-architect, mermaid-config, mermaid-diagram, mermaid-render, mermaid-validate, plus the main SKILL.md hub) and 7 specialist agent definitions.

**Acceptance Criteria:**

1. Main `SKILL.md` becomes `PLUGIN.md` in `GentlemanClaude/plugins/mermaid/`
2. 5 sub-skills are ported as sub-directories under `GentlemanClaude/plugins/mermaid/skills/`: `mermaid-architect/`, `mermaid-config/`, `mermaid-diagram/`, `mermaid-render/`, `mermaid-validate/`
3. 7 specialist definitions are copied to `GentlemanClaude/plugins/mermaid/specialists/`
4. 1 agent definition (`diagram-architect.md`) is copied to `GentlemanClaude/plugins/mermaid/agents/`
5. Reference guides (6 files) are copied to `GentlemanClaude/plugins/mermaid/references/guides/`
6. Code-to-diagram examples (6 framework dirs + 1 master README) are copied to `GentlemanClaude/plugins/mermaid/examples/`
7. Design document templates (6 files) are copied to `GentlemanClaude/plugins/mermaid/assets/`
8. JavaScript scripts (`extract_mermaid.js`, `render.js`, `resilient_diagram.js`) are copied to `GentlemanClaude/plugins/mermaid/scripts/`
9. `package.json` is copied to `GentlemanClaude/plugins/mermaid/scripts/`
10. `ensure-deps.sh` is copied to `GentlemanClaude/plugins/mermaid/scripts/` and adapted: `npm install --prefix` uses `~/.claude/plugins/mermaid/scripts/` for `node_modules/`
11. ALL `find $HOME/.claude/plugins/cache` patterns across ALL sub-skills and the agent are rewritten to static `~/.claude/plugins/mermaid/` paths
12. Internal references between sub-skills (e.g., "load `references/guides/troubleshooting.md`") use paths relative to `~/.claude/plugins/mermaid/`
13. `tests/` directory from upstream is NOT ported
14. `LICENSE` file is preserved in the plugin directory
15. Unicode symbols guide is copied to `GentlemanClaude/plugins/mermaid/references/guides/unicode-symbols/`

---

### REQ-PLUG-06: Dependency Handling

**Description:** Each plugin that requires runtime dependencies ships an `ensure-deps.sh` script that validates prerequisites before first use.

**Acceptance Criteria:**

1. Every plugin has a `scripts/ensure-deps.sh` script
2. `merge-checks/ensure-deps.sh`: checks for `git`, warns (non-fatal) if `node` is missing (i18n checks skipped)
3. `trim-md/ensure-deps.sh`: checks for `npx` (ships with Node.js 18+), exits with error + install instructions if missing
4. `mermaid/ensure-deps.sh`: checks for `node_modules/beautiful-mermaid`, runs `npm install --prefix` to install locally if missing
5. Dependencies install to `~/.claude/plugins/<name>/scripts/node_modules/` (isolated per plugin)
6. No global npm installs -- all dependencies are local to the plugin's scripts directory
7. `ensure-deps.sh` provides platform-specific install hints (macOS: `brew install node`, Linux: `sudo apt install nodejs npm`)
8. If Node.js is completely absent, error messages clearly state which plugins are affected and which still work (merge-checks works fine)

---

### REQ-PLUG-07: TUI Integration

**Description:** The existing Skill Manager TUI (ScreenSkillMenu) is extended to discover, install, and remove plugins alongside skills.

**Acceptance Criteria:**

1. `fetchSkillCatalog()` is extended to also scan `GentlemanClaude/plugins/` for `PLUGIN.md` files
2. `SkillInfo` struct gains a `Type` field: `"skill"` or `"plugin"` to distinguish entries
3. `parseSkillFrontmatter()` is adapted to also parse `PLUGIN.md` files (same YAML frontmatter format, different filename)
4. Plugins appear in the Browse screen with a distinct category label (e.g., "Plugins" category header)
5. Plugin install copies the entire plugin directory to `~/.claude/plugins/<name>/` (NOT a symlink -- plugins have scripts that need to be executable)
6. Plugin remove deletes the `~/.claude/plugins/<name>/` directory
7. `installSkillSymlinks()` is extended: for `type: plugin` entries, it copies the directory instead of creating a symlink
8. Plugin installed state is checked via directory existence: `~/.claude/plugins/<name>/PLUGIN.md` exists
9. Skill Manager menu options remain the same ("Browse", "Install", "Remove", "Update Catalog") -- plugins are integrated into existing flows, not a separate section
10. `settings.json` permission entries for installed plugins are automatically appended during install
11. **Dependency:** This requirement depends on the `project-init-and-skill-manager` change being merged first

---

### REQ-PLUG-08: Documentation Updates

**Description:** Update all relevant documentation to reflect the new plugin system.

**Acceptance Criteria:**

1. `GentlemanClaude/CLAUDE.md` auto-load table gains a new section for plugins with entries for merge-checks, trim-md, and mermaid
2. `docs/ai-framework-modules.md` module count updates from 203 to 206
3. `docs/ai-framework-modules.md` gains 3 new entries in a "Plugins" category
4. `GentlemanClaude/settings.json` gains plugin script permission entries in `permissions.allow`
5. Each `PLUGIN.md` includes a "Common Mistakes" section documenting known issues and fixes
6. Each `PLUGIN.md` includes dependency requirements in its frontmatter

---

## Scenarios

### SCN-01: User installs all 3 plugins

**Preconditions:** User has Javi.Dots installed, Claude CLI configured, Node.js 18+ available.

**Steps:**

1. User opens Skill Manager TUI
2. User navigates to "Install Skills"
3. User sees "Plugins" category with 3 entries: merge-checks, trim-md, mermaid
4. User selects all 3 and confirms
5. Installer copies each plugin directory to `~/.claude/plugins/<name>/`
6. Installer appends permission entries to `~/.claude/settings.json`
7. User runs `/merge-checks` -- gather-context.sh executes, scope selection works
8. User runs `/trim-md .` -- ensure-deps.sh runs, markdownlint executes
9. User runs `/mermaid-diagram "sequence of login flow"` -- ensure-deps.sh installs beautiful-mermaid, diagram generates

**Expected Result:** All 3 plugins functional. Each plugin's scripts execute from `~/.claude/plugins/<name>/scripts/`.

---

### SCN-02: User installs only merge-checks

**Preconditions:** User has Javi.Dots installed, Claude CLI configured, NO Node.js.

**Steps:**

1. User opens Skill Manager TUI
2. User navigates to "Install Skills"
3. User selects only merge-checks and confirms
4. Installer copies merge-checks to `~/.claude/plugins/merge-checks/`
5. Installer appends merge-checks permission entry to `~/.claude/settings.json`
6. User runs `/merge-checks`
7. `gather-context.sh` executes successfully (requires only `git` + `bash`)
8. `precompute.sh` runs all checks -- i18n consistency check is skipped with a warning (no Node.js), all other 12 checks pass

**Expected Result:** merge-checks fully functional except i18n consistency check. No errors about missing Node.js beyond the one-line warning.

---

### SCN-03: User without Node.js tries trim-md

**Preconditions:** User has Javi.Dots installed, Claude CLI configured, NO Node.js or npm.

**Steps:**

1. User has trim-md plugin installed
2. User runs `/trim-md .`
3. `ensure-deps.sh` detects `npx` is missing
4. Script outputs: "trim-md: missing dependencies: - npx (install Node.js 18+: https://nodejs.org)"
5. Script includes platform-specific hint: "Quick install (macOS): brew install node" or "Quick install (Linux): sudo apt install nodejs npm"
6. Script exits with code 1
7. Claude presents the error message to the user and stops

**Expected Result:** Clear, actionable error. No cryptic stack traces. User knows exactly what to install and how.

---

### SCN-04: Non-interactive mode (CLI install)

**Preconditions:** User runs the Javi.Dots installer in non-interactive mode with `--ai-plugins merge-checks,mermaid` flag (or equivalent).

**Steps:**

1. Installer detects the plugin names from CLI args
2. Installer locates `GentlemanClaude/plugins/merge-checks/` and `GentlemanClaude/plugins/mermaid/` in the repo
3. Installer copies both plugin directories to `~/.claude/plugins/`
4. Installer appends permission entries to `~/.claude/settings.json`
5. No TUI screens are shown

**Expected Result:** Plugins installed silently. Same end state as TUI install. Exit code 0.

---

### SCN-05: Plugin removal

**Preconditions:** User has mermaid plugin installed at `~/.claude/plugins/mermaid/`.

**Steps:**

1. User opens Skill Manager TUI
2. User navigates to "Remove Skills"
3. User sees mermaid listed under "Plugins" with "installed" indicator
4. User selects mermaid and confirms removal
5. Installer deletes `~/.claude/plugins/mermaid/` directory (including `node_modules/` if present)
6. Permission entries for mermaid scripts are removed from `~/.claude/settings.json`

**Expected Result:** Plugin fully removed. No orphaned files. Permission entries cleaned up.

---

### SCN-06: Path resolution verification

**Preconditions:** merge-checks plugin installed.

**Steps:**

1. User runs `/merge-checks`
2. PLUGIN.md inline bash executes: `bash ~/.claude/plugins/merge-checks/scripts/gather-context.sh`
3. gather-context.sh resolves `SCRIPTS` dir via `$(cd "$(dirname "$0")" && pwd)` (works because the script is at a fixed location)
4. precompute.sh calls sibling scripts: `bash "$SCRIPTS/check-debug-artifacts.sh"` etc.
5. PLUGIN.md references check files: `checks/docs.md` resolves relative to `~/.claude/plugins/merge-checks/`

**Expected Result:** No `find` commands executed at runtime. All paths resolve statically. Scripts find each other via `$SCRIPTS` (dirname-based) or absolute `~/.claude/plugins/` paths.
