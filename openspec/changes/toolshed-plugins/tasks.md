# Tasks: toolshed-plugins

## Phase 1: Infrastructure

### TASK-1.1: Create plugin directory structure

**Description:** Create `GentlemanClaude/plugins/` with subdirectories for all 3 plugins. Establish the PLUGIN.md format convention.

**Files affected:**
- `GentlemanClaude/plugins/` (new directory)
- `GentlemanClaude/plugins/merge-checks/` (new directory)
- `GentlemanClaude/plugins/trim-md/` (new directory)
- `GentlemanClaude/plugins/mermaid/` (new directory)

**Effort:** Small (30 min)
**Dependencies:** None

---

### TASK-1.2: Define PLUGIN.md frontmatter spec

**Description:** Create a reference/template `PLUGIN.md` with the standardized frontmatter format. Document the field mapping from toolshed's SKILL.md to Gentleman's PLUGIN.md. This is informational -- the actual PLUGIN.md files are created per-plugin in later tasks.

**Files affected:**
- (Informational -- conventions documented in design.md, applied in TASK-2.1, TASK-3.1, TASK-4.1)

**Effort:** Trivial (15 min)
**Dependencies:** TASK-1.1

---

## Phase 2: merge-checks Port

### TASK-2.1: Create merge-checks PLUGIN.md

**Description:** Convert toolshed's `skills/merge-checks/SKILL.md` to `GentlemanClaude/plugins/merge-checks/PLUGIN.md`. Rewrite frontmatter (drop `allowed-tools`, `argument-hint`, `model`; add `type: plugin`, `dependencies`, `permissions`, `license`, `metadata`). Rewrite ALL `find` path resolution patterns to static `~/.claude/plugins/merge-checks/` paths.

**Source:** `/tmp/claude-toolshed/plugins/merge-checks/skills/merge-checks/SKILL.md` (11,160 bytes)
**Target:** `GentlemanClaude/plugins/merge-checks/PLUGIN.md`

**Path rewrites required:**
- `gather-context.sh` find pattern -> `bash "$HOME/.claude/plugins/merge-checks/scripts/gather-context.sh"`
- `precompute.sh` find pattern -> `bash "$HOME/.claude/plugins/merge-checks/scripts/precompute.sh"`
- Agent instruction paths (checks/*.md) -> relative from `~/.claude/plugins/merge-checks/`

**Effort:** Medium (1h)
**Dependencies:** TASK-1.1, TASK-1.2

---

### TASK-2.2: Copy merge-checks scripts

**Description:** Copy all 18 shell scripts from toolshed to `GentlemanClaude/plugins/merge-checks/scripts/`. Verify that internal cross-references (`$SCRIPTS`, `$(dirname "$0")`) still work from the new location. No path changes should be needed in the scripts themselves since they use relative dirname patterns.

**Source files (18):**
- `ensure-deps.sh` (603 bytes)
- `gather-context.sh` (8,279 bytes)
- `precompute.sh` (18,114 bytes)
- `build-manifest.sh` (1,494 bytes)
- `detect-mode.sh` (4,891 bytes)
- `detect-features.sh` (11,308 bytes)
- `detect-env-vars.sh` (1,336 bytes)
- `detect-suppressions.sh` (1,586 bytes)
- `list-hardcoded-strings.sh` (2,934 bytes)
- `check-debug-artifacts.sh` (3,817 bytes)
- `check-env-coverage.sh` (1,944 bytes)
- `check-i18n-consistency.sh` (2,924 bytes)
- `check-migration-exists.sh` (2,671 bytes)
- `check-route-registered.sh` (2,828 bytes)
- `check-seed-imported.sh` (2,994 bytes)
- `check-shared-types.sh` (2,288 bytes)
- `check-story-exists.sh` (1,848 bytes)
- `check-test-exists.sh` (2,921 bytes)

**Target:** `GentlemanClaude/plugins/merge-checks/scripts/`

**Effort:** Small (30 min -- mostly copy + verify)
**Dependencies:** TASK-1.1

---

### TASK-2.3: Copy merge-checks check definitions

**Description:** Copy the 3 agent instruction files for Phase 2 reasoning agents.

**Source files (3):**
- `checks/docs.md` (1,240 bytes) -- Check 1: Documentation gaps
- `checks/comments.md` (1,243 bytes) -- Check 2: Comment quality
- `checks/shared.md` (1,403 bytes) -- Check 12: Shared contracts

**Target:** `GentlemanClaude/plugins/merge-checks/checks/`

**Effort:** Trivial (10 min)
**Dependencies:** TASK-1.1

---

### TASK-2.4: Copy merge-checks assets

**Description:** Copy workflow diagram assets (SVG + MMD source files).

**Source files (4):**
- `assets/activity-merge-checks-workflow.mmd`
- `assets/activity-merge-checks-workflow.svg`
- `assets/activity-script-execution.mmd`
- `assets/activity-script-execution.svg`

**Target:** `GentlemanClaude/plugins/merge-checks/assets/`

**Effort:** Trivial (5 min)
**Dependencies:** TASK-1.1

---

### TASK-2.5: Copy merge-checks README

**Description:** Copy the upstream README for reference/attribution.

**Source:** `README.md`
**Target:** `GentlemanClaude/plugins/merge-checks/README.md`

**Effort:** Trivial (5 min)
**Dependencies:** TASK-1.1

---

## Phase 3: trim-md Port

### TASK-3.1: Create trim-md PLUGIN.md

**Description:** Convert toolshed's `skills/trim-md/SKILL.md` to `GentlemanClaude/plugins/trim-md/PLUGIN.md`. Rewrite frontmatter. Replace the `find`-based `SKILL_DIR` resolution with static `PLUGIN_DIR="$HOME/.claude/plugins/trim-md"`. Update all downstream `$SKILL_DIR` references to `$PLUGIN_DIR`.

**Source:** `/tmp/claude-toolshed/plugins/trim-md/skills/trim-md/SKILL.md`
**Target:** `GentlemanClaude/plugins/trim-md/PLUGIN.md`

**Path rewrites required:**
- `SKILL_DIR` find block -> `PLUGIN_DIR="$HOME/.claude/plugins/trim-md"`
- `$SKILL_DIR/scripts/ensure-deps.sh` -> `$PLUGIN_DIR/scripts/ensure-deps.sh`
- `$SKILL_DIR/scripts/trim-md.sh` -> `$PLUGIN_DIR/scripts/trim-md.sh`

**Effort:** Small (30 min)
**Dependencies:** TASK-1.1, TASK-1.2

---

### TASK-3.2: Copy and adapt trim-md scripts

**Description:** Copy `ensure-deps.sh` and `trim-md.sh` to `GentlemanClaude/plugins/trim-md/scripts/`. Modify `trim-md.sh` line 9: change `REF_DIR="$SCRIPT_DIR/../reference"` to `REF_DIR="$SCRIPT_DIR/../config"` (matching our renamed config directory).

**Source files (2):**
- `scripts/ensure-deps.sh` (27 lines)
- `scripts/trim-md.sh` (275 lines)

**Target:** `GentlemanClaude/plugins/trim-md/scripts/`

**Modification:** `trim-md.sh` line 9: `reference` -> `config`

**Effort:** Small (15 min)
**Dependencies:** TASK-1.1

---

### TASK-3.3: Copy trim-md config files

**Description:** Copy markdownlint configuration files. Rename directory from `reference/` to `config/` for clarity.

**Source files (2):**
- `reference/full.markdownlint-cli2.jsonc`
- `reference/safe.markdownlint-cli2.jsonc`

**Target:** `GentlemanClaude/plugins/trim-md/config/`

**Effort:** Trivial (5 min)
**Dependencies:** TASK-1.1

---

### TASK-3.4: Verify trim-md hook exclusion

**Description:** Confirm that `hooks/hooks.json` and `hooks/posttool-trim-md.sh` are NOT included in the port. Document in PLUGIN.md that the PostToolUse hook is intentionally dropped and trim-md is manual-invoke only.

**Source files to SKIP:**
- `hooks/hooks.json`
- `hooks/posttool-trim-md.sh`

**Effort:** Trivial (5 min -- verification only)
**Dependencies:** TASK-3.1

---

## Phase 4: mermaid Port

### TASK-4.1: Create mermaid PLUGIN.md

**Description:** Convert toolshed's `skills/mermaid/SKILL.md` to `GentlemanClaude/plugins/mermaid/PLUGIN.md`. Rewrite frontmatter. The hub PLUGIN.md routes to sub-skills -- update all internal file references to use paths relative to `~/.claude/plugins/mermaid/`.

**Source:** `/tmp/claude-toolshed/plugins/mermaid/skills/mermaid/SKILL.md` (6,948 bytes)
**Target:** `GentlemanClaude/plugins/mermaid/PLUGIN.md`

**Path rewrites required:**
- All `references/guides/*.md` paths -> relative to `~/.claude/plugins/mermaid/`
- All `examples/*/README.md` paths -> relative to `~/.claude/plugins/mermaid/`
- All `assets/*.md` paths -> relative to `~/.claude/plugins/mermaid/`
- All `scripts/*.js` paths -> relative to `~/.claude/plugins/mermaid/`

**Effort:** Medium (45 min)
**Dependencies:** TASK-1.1, TASK-1.2

---

### TASK-4.2: Port mermaid sub-skills (5 files)

**Description:** Copy and adapt the 5 sub-skill SKILL.md files. Each has a `find` path resolution block that must be rewritten to static `PLUGIN_DIR="$HOME/.claude/plugins/mermaid"`.

**Source files (5):**
- `skills/mermaid-architect/SKILL.md` -- codebase analysis + multi-diagram generation
- `skills/mermaid-config/SKILL.md` -- config wizard
- `skills/mermaid-diagram/SKILL.md` -- single diagram from description
- `skills/mermaid-render/SKILL.md` -- render .mmd/.md to SVG
- `skills/mermaid-validate/SKILL.md` -- validate Mermaid syntax

**Target:** `GentlemanClaude/plugins/mermaid/skills/<name>/SKILL.md`

**Path rewrites per file:**
- Replace `find "$HOME/.claude/plugins/cache"...` block with `PLUGIN_DIR="$HOME/.claude/plugins/mermaid"`
- Replace `find "$HOME" -maxdepth 8...` fallback block (remove entirely)
- Update `$PLUGIN_DIR/scripts/ensure-deps.sh` to `$HOME/.claude/plugins/mermaid/scripts/ensure-deps.sh`

**Effort:** Medium (1h)
**Dependencies:** TASK-1.1

---

### TASK-4.3: Copy mermaid specialists (7 files)

**Description:** Copy the 7 diagram-type specialist definitions. These are pure markdown -- no path changes needed.

**Source files (7):**
- `specialists/mermaid-activity.md` (2,832 bytes)
- `specialists/mermaid-architecture.md` (3,064 bytes)
- `specialists/mermaid-class.md` (3,154 bytes)
- `specialists/mermaid-deployment.md` (3,013 bytes)
- `specialists/mermaid-er.md` (3,138 bytes)
- `specialists/mermaid-sequence.md` (2,785 bytes)
- `specialists/mermaid-state.md` (3,161 bytes)

**Target:** `GentlemanClaude/plugins/mermaid/specialists/`

**Effort:** Small (15 min)
**Dependencies:** TASK-1.1

---

### TASK-4.4: Port mermaid agent definition (1 file)

**Description:** Copy and adapt `agents/diagram-architect.md`. This file has a `find` path resolution block that must be rewritten.

**Source:** `agents/diagram-architect.md` (8,960 bytes)
**Target:** `GentlemanClaude/plugins/mermaid/agents/diagram-architect.md`

**Path rewrite:**
- `find "$HOME/.claude/plugins/cache" -type d -name "mermaid"...` -> `PLUGIN_DIR="$HOME/.claude/plugins/mermaid"`
- All downstream `$PLUGIN_DIR/...` references adjusted accordingly

**Effort:** Small (20 min)
**Dependencies:** TASK-1.1

---

### TASK-4.5: Copy mermaid reference guides (7 files)

**Description:** Copy all reference guide markdown files. Pure markdown -- no path changes needed.

**Source files (7):**
- `references/guides/common-mistakes.md` (10,067 bytes)
- `references/guides/quick-decision-matrix.md` (1,440 bytes)
- `references/guides/resilient-workflow.md` (16,461 bytes)
- `references/guides/styling-guide.md` (21,440 bytes)
- `references/guides/troubleshooting.md` (15,384 bytes)
- `references/guides/code-to-diagram/README.md`
- `references/guides/unicode-symbols/guide.md`

**Target:** `GentlemanClaude/plugins/mermaid/references/guides/`

**Effort:** Small (15 min)
**Dependencies:** TASK-1.1

---

### TASK-4.6: Copy mermaid examples (6 directories)

**Description:** Copy framework-specific code-to-diagram examples. Pure markdown -- no path changes needed.

**Source files (6):**
- `examples/fastapi/README.md`
- `examples/java-webapp/README.md`
- `examples/node-webapp/README.md`
- `examples/python-etl/README.md`
- `examples/react/README.md`
- `examples/spring-boot/README.md`

**Target:** `GentlemanClaude/plugins/mermaid/examples/`

**Effort:** Small (10 min)
**Dependencies:** TASK-1.1

---

### TASK-4.7: Copy mermaid design templates (6 files)

**Description:** Copy design document templates.

**Source files (6):**
- `assets/api-design-template.md`
- `assets/architecture-design-template.md`
- `assets/database-design-template.md`
- `assets/feature-design-template.md`
- `assets/local-config-template.md`
- `assets/system-design-template.md`

**Target:** `GentlemanClaude/plugins/mermaid/assets/`

**Effort:** Trivial (10 min)
**Dependencies:** TASK-1.1

---

### TASK-4.8: Copy and adapt mermaid scripts (4 files + package.json)

**Description:** Copy JavaScript scripts and `ensure-deps.sh`. Adapt `ensure-deps.sh` so that `npm install --prefix` targets the installed location `~/.claude/plugins/mermaid/scripts/` for node_modules.

**Source files (5):**
- `scripts/ensure-deps.sh` (15 lines)
- `scripts/extract_mermaid.js` (2,318 bytes)
- `scripts/render.js` (2,269 bytes)
- `scripts/resilient_diagram.js` (5,559 bytes)
- `scripts/package.json` (6 lines)

**Target:** `GentlemanClaude/plugins/mermaid/scripts/`

**Adaptation:** `ensure-deps.sh` uses `SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"` + `npm install --prefix "$SCRIPT_DIR"` -- this is already correct and works from any installed location. No changes needed.

**Effort:** Small (15 min)
**Dependencies:** TASK-1.1

---

### TASK-4.9: Copy mermaid LICENSE

**Description:** Preserve the MIT license file.

**Source:** `LICENSE`
**Target:** `GentlemanClaude/plugins/mermaid/LICENSE`

**Effort:** Trivial (2 min)
**Dependencies:** TASK-1.1

---

### TASK-4.10: Copy mermaid README

**Description:** Copy the upstream README for reference.

**Source:** `README.md`
**Target:** `GentlemanClaude/plugins/mermaid/README.md`

**Effort:** Trivial (2 min)
**Dependencies:** TASK-1.1

---

## Phase 5: TUI Integration

### TASK-5.1: Extend SkillInfo struct with Type field

**Description:** Add `Type string` and `Permissions []string` fields to the `SkillInfo` struct in `model.go`. Update `NewModel()` defaults.

**Files affected:**
- `installer/internal/tui/model.go` (SkillInfo struct at line 765)

**Effort:** Trivial (10 min)
**Dependencies:** None (can be done in parallel with Phase 2-4)

---

### TASK-5.2: Create parsePluginFrontmatter function

**Description:** Create `parsePluginFrontmatter()` that reads `PLUGIN.md` YAML frontmatter. Similar to `parseSkillFrontmatter()` but also extracts `type:` and `permissions:` fields.

**Files affected:**
- `installer/internal/tui/update.go` (after `parseSkillFrontmatter` at line ~497)

**Effort:** Small (30 min)
**Dependencies:** TASK-5.1

---

### TASK-5.3: Extend fetchSkillCatalog to scan plugins

**Description:** Add plugin scanning to `fetchSkillCatalog()`. After scanning `~/.gentleman/skills/` and `~/.claude/skills/`, also scan the repo's `GentlemanClaude/plugins/` directory for `PLUGIN.md` files. Set `Category: "plugin"` and `Type: "plugin"` on discovered entries.

**Files affected:**
- `installer/internal/tui/update.go` (`fetchSkillCatalog` function at line ~314)

**Effort:** Medium (45 min)
**Dependencies:** TASK-5.1, TASK-5.2

---

### TASK-5.4: Extend installSkillSymlinks for plugins

**Description:** Branch on `Type == "plugin"`: instead of creating a symlink, copy the entire plugin directory to `~/.claude/plugins/<name>/`. Make all scripts executable. Append permission entries to `~/.claude/settings.json`.

**Files affected:**
- `installer/internal/tui/update.go` (`installSkillSymlinks` function at line ~568)

**Effort:** Medium (1h)
**Dependencies:** TASK-5.1, TASK-5.2

---

### TASK-5.5: Extend removeSkillSymlinks for plugins

**Description:** Branch on `Type == "plugin"`: instead of removing a symlink, `os.RemoveAll()` the `~/.claude/plugins/<name>/` directory. Remove corresponding permission entries from `~/.claude/settings.json`.

**Files affected:**
- `installer/internal/tui/update.go` (`removeSkillSymlinks` function at line ~616)

**Effort:** Small (30 min)
**Dependencies:** TASK-5.1

---

### TASK-5.6: Add isPluginInstalled helper

**Description:** Create `isPluginInstalled(home, name string) bool` that checks for `~/.claude/plugins/<name>/PLUGIN.md` existence.

**Files affected:**
- `installer/internal/tui/update.go` (near `isSkillInstalled` at line ~550)

**Effort:** Trivial (10 min)
**Dependencies:** None

---

### TASK-5.7: Update category ordering and headers

**Description:** Add `"plugin"` to `getSkillCategoryOrder()` between `"community"` and `"local"`. Add header string for `skillCategoryHeader("plugin")` returning `"━━━ Plugins ━━━"`.

**Files affected:**
- `installer/internal/tui/update.go` (category functions)
- `installer/internal/tui/view.go` (if header rendering is there)

**Effort:** Trivial (15 min)
**Dependencies:** TASK-5.1

---

### TASK-5.8: Ensure ~/.claude/plugins/ directory in installer

**Description:** In the Claude AI tool install step (`installer.go` line ~1106), add `system.EnsureDir(filepath.Join(claudeDir, "plugins"))` to create the plugins directory alongside the existing skills directory.

**Files affected:**
- `installer/internal/tui/installer.go` (line ~1108, after `system.EnsureDir(filepath.Join(claudeDir, "skills"))`)

**Effort:** Trivial (5 min)
**Dependencies:** None

---

### TASK-5.9: Update TUI tests

**Description:** Add test cases for plugin discovery, install, and remove in `skill_screens_test.go`. Test that plugins show under "plugin" category, install via copy (not symlink), and remove via directory deletion.

**Files affected:**
- `installer/internal/tui/skill_screens_test.go`

**Effort:** Medium (1h)
**Dependencies:** TASK-5.1 through TASK-5.7

---

## Phase 6: Documentation

### TASK-6.1: Update CLAUDE.md auto-load table

**Description:** Add a "Plugin Detection" section to the auto-load table in `GentlemanClaude/CLAUDE.md` with entries for merge-checks, trim-md, and mermaid.

**Files affected:**
- `GentlemanClaude/CLAUDE.md` (auto-load table, after "Framework/Library Detection")

**New entries:**

| Context | Read this file |
|---|---|
| PR review, merge audit, code quality checks | `~/.claude/plugins/merge-checks/PLUGIN.md` |
| Markdown cleanup, lint, token optimization | `~/.claude/plugins/trim-md/PLUGIN.md` |
| Mermaid diagrams, architecture docs, SVG gen | `~/.claude/plugins/mermaid/PLUGIN.md` |

**Effort:** Small (15 min)
**Dependencies:** TASK-2.1, TASK-3.1, TASK-4.1

---

### TASK-6.2: Update settings.json with plugin permissions

**Description:** Add plugin script permission entries to `GentlemanClaude/settings.json` in the `permissions.allow` array.

**Files affected:**
- `GentlemanClaude/settings.json`

**New entries:**
```json
"Bash(~/.claude/plugins/merge-checks/scripts/*:*)",
"Bash(~/.claude/plugins/trim-md/scripts/*:*)",
"Bash(~/.claude/plugins/mermaid/scripts/*:*)"
```

**Effort:** Trivial (5 min)
**Dependencies:** None

---

### TASK-6.3: Update module registry

**Description:** Update `docs/ai-framework-modules.md` module count from 203 to 206 and add 3 new plugin entries.

**Files affected:**
- `docs/ai-framework-modules.md`

**New entries:**

| Module key | Display name | Description |
|---|---|---|
| `plugin:merge-checks` | Plugin: Merge Checks | 13-dimension code quality audit for pre/post merge review |
| `plugin:trim-md` | Plugin: Trim MD | Markdown linting and optimization for LLM consumption |
| `plugin:mermaid` | Plugin: Mermaid | Diagram generation with 7 diagram types and code-to-diagram |

**Effort:** Small (15 min)
**Dependencies:** None

---

## Phase 7: Verification

### TASK-7.1: Verify merge-checks path resolution

**Description:** Manually trace every path reference in `merge-checks/PLUGIN.md` and verify it resolves correctly when the plugin is at `~/.claude/plugins/merge-checks/`. Test `gather-context.sh` and `precompute.sh` execution from the installed location.

**Verification steps:**
1. Copy plugin to `~/.claude/plugins/merge-checks/`
2. Run `bash ~/.claude/plugins/merge-checks/scripts/gather-context.sh` in a git repo
3. Run `bash ~/.claude/plugins/merge-checks/scripts/precompute.sh` with a base branch
4. Verify all 18 scripts are found by precompute.sh via `$SCRIPTS` variable

**Effort:** Medium (30 min)
**Dependencies:** TASK-2.1 through TASK-2.5

---

### TASK-7.2: Verify trim-md path resolution + ensure-deps

**Description:** Test trim-md from the installed location. Verify the `config/` directory rename is correctly referenced by `trim-md.sh`.

**Verification steps:**
1. Copy plugin to `~/.claude/plugins/trim-md/`
2. Run `bash ~/.claude/plugins/trim-md/scripts/ensure-deps.sh` (should pass if Node.js installed)
3. Run `bash ~/.claude/plugins/trim-md/scripts/trim-md.sh --dry-run .` on a directory with .md files
4. Verify it finds `config/full.markdownlint-cli2.jsonc` or `config/safe.markdownlint-cli2.jsonc`

**Effort:** Small (20 min)
**Dependencies:** TASK-3.1 through TASK-3.3

---

### TASK-7.3: Verify mermaid path resolution + ensure-deps

**Description:** Test mermaid scripts from the installed location. Verify `ensure-deps.sh` installs `beautiful-mermaid` locally and JS scripts can find it.

**Verification steps:**
1. Copy plugin to `~/.claude/plugins/mermaid/`
2. Run `bash ~/.claude/plugins/mermaid/scripts/ensure-deps.sh` (should install beautiful-mermaid)
3. Verify `~/.claude/plugins/mermaid/scripts/node_modules/beautiful-mermaid/` exists
4. Run `node ~/.claude/plugins/mermaid/scripts/extract_mermaid.js --help` (or equivalent)
5. Verify each sub-skill SKILL.md has no remaining `find` patterns

**Effort:** Medium (30 min)
**Dependencies:** TASK-4.1 through TASK-4.10

---

### TASK-7.4: Audit for remaining `find` patterns

**Description:** Run `rg 'find.*plugins/cache' GentlemanClaude/plugins/` to confirm zero remaining dynamic path resolution patterns across all ported files.

**Verification steps:**
1. `rg 'find.*\.claude' GentlemanClaude/plugins/` -- should return 0 matches
2. `rg 'plugins/cache' GentlemanClaude/plugins/` -- should return 0 matches
3. `rg 'maxdepth' GentlemanClaude/plugins/` -- should return 0 matches

**Effort:** Trivial (10 min)
**Dependencies:** TASK-2.1, TASK-3.1, TASK-4.1, TASK-4.2, TASK-4.4

---

### TASK-7.5: Verify no-Node.js degradation

**Description:** Test each plugin's behavior when Node.js is not available. merge-checks should work (with i18n warning). trim-md and mermaid should fail gracefully with actionable error messages.

**Verification steps:**
1. Temporarily rename `node` binary (or use a clean env)
2. Run merge-checks ensure-deps.sh -- should warn about missing node, exit 0
3. Run trim-md ensure-deps.sh -- should error with install instructions, exit 1
4. Run mermaid ensure-deps.sh -- should error (npm not found), exit 1
5. Restore node binary

**Effort:** Small (20 min)
**Dependencies:** TASK-7.1, TASK-7.2, TASK-7.3

---

## Summary

| Phase | Tasks | Estimated Total |
|---|---|---|
| 1. Infrastructure | 2 tasks | 45 min |
| 2. merge-checks port | 5 tasks | 2h 10min |
| 3. trim-md port | 4 tasks | 55 min |
| 4. mermaid port | 10 tasks | 3h |
| 5. TUI integration | 9 tasks | 4h 25min |
| 6. Documentation | 3 tasks | 35 min |
| 7. Verification | 5 tasks | 1h 50min |
| **Total** | **38 tasks** | **~13.5h** |

### Dependency Graph (Phase Level)

```
Phase 1 (Infrastructure)
    │
    ├── Phase 2 (merge-checks)  ──┐
    ├── Phase 3 (trim-md)       ──┼── Phase 6 (Documentation)
    ├── Phase 4 (mermaid)       ──┘         │
    │                                       │
    └── Phase 5 (TUI integration) ──────────┤
                                            │
                                      Phase 7 (Verification)
```

Phases 2, 3, 4, and 5 can run in parallel after Phase 1 completes. Phase 6 and 7 depend on all prior phases.
