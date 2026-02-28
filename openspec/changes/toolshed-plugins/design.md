# Design: toolshed-plugins

## 1. Directory Layout

### Repository Structure (`GentlemanClaude/plugins/`)

```
GentlemanClaude/
  skills/                          # existing: 23 pure .md skills
  plugins/                         # NEW: executable plugins
    merge-checks/
      PLUGIN.md                    # entry point
      checks/
        docs.md                    # Check 1 agent instructions
        comments.md                # Check 2 agent instructions
        shared.md                  # Check 12 agent instructions
      scripts/
        ensure-deps.sh
        gather-context.sh
        precompute.sh
        build-manifest.sh
        detect-mode.sh
        detect-features.sh
        detect-env-vars.sh
        detect-suppressions.sh
        list-hardcoded-strings.sh
        check-debug-artifacts.sh
        check-env-coverage.sh
        check-i18n-consistency.sh
        check-migration-exists.sh
        check-route-registered.sh
        check-seed-imported.sh
        check-shared-types.sh
        check-story-exists.sh
        check-test-exists.sh
      assets/
        activity-merge-checks-workflow.mmd
        activity-merge-checks-workflow.svg
        activity-script-execution.mmd
        activity-script-execution.svg
    trim-md/
      PLUGIN.md                    # entry point
      scripts/
        ensure-deps.sh
        trim-md.sh
      config/                      # renamed from "reference/" for clarity
        full.markdownlint-cli2.jsonc
        safe.markdownlint-cli2.jsonc
    mermaid/
      PLUGIN.md                    # entry point (hub router)
      LICENSE
      skills/                      # 5 sub-skills
        mermaid-architect/
          SKILL.md
        mermaid-config/
          SKILL.md
        mermaid-diagram/
          SKILL.md
        mermaid-render/
          SKILL.md
        mermaid-validate/
          SKILL.md
      specialists/                 # 7 diagram type experts
        mermaid-activity.md
        mermaid-architecture.md
        mermaid-class.md
        mermaid-deployment.md
        mermaid-er.md
        mermaid-sequence.md
        mermaid-state.md
      agents/
        diagram-architect.md       # autonomous agent definition
      references/
        guides/
          common-mistakes.md
          quick-decision-matrix.md
          resilient-workflow.md
          styling-guide.md
          troubleshooting.md
          code-to-diagram/
            README.md
          unicode-symbols/
            guide.md
      examples/
        fastapi/README.md
        java-webapp/README.md
        node-webapp/README.md
        python-etl/README.md
        react/README.md
        spring-boot/README.md
      assets/                      # design doc templates
        api-design-template.md
        architecture-design-template.md
        database-design-template.md
        feature-design-template.md
        local-config-template.md
        system-design-template.md
      scripts/
        ensure-deps.sh
        extract_mermaid.js
        render.js
        resilient_diagram.js
        package.json
```

### Installed Location (`~/.claude/plugins/`)

```
~/.claude/
  skills/                          # existing: symlinks to ~/.gentleman/skills/
  plugins/                         # NEW: plugin directories (full copies, NOT symlinks)
    merge-checks/                  # entire dir copied from repo
      PLUGIN.md
      checks/
      scripts/
      assets/
    trim-md/
      PLUGIN.md
      scripts/
      config/
    mermaid/
      PLUGIN.md
      skills/
      specialists/
      agents/
      references/
      examples/
      assets/
      scripts/
        node_modules/              # created by ensure-deps.sh on first run
          beautiful-mermaid/
        package.json
```

**Key difference from skills:** Skills are installed as symlinks (`~/.claude/skills/react-19 -> ~/.gentleman/skills/curated/react-19`). Plugins are full directory copies because they contain executable scripts that need `chmod +x` and may generate local `node_modules/`.

---

## 2. PLUGIN.md Frontmatter Format

### Gentleman Skill Frontmatter (existing)

```yaml
---
name: react-19
description: >
  React 19 patterns with React Compiler.
  Trigger: When writing React components - no useMemo/useCallback needed.
license: Apache-2.0
metadata:
  author: gentleman-programming
  version: "1.0"
---
```

### Toolshed Skill Frontmatter (upstream)

```yaml
---
name: merge-checks
description: Audit code changes across 13 quality dimensions before or after merge
argument-hint: [scope-argument]
allowed-tools: Bash, Read, Write, Task, AskUserQuestion
model: opus
---
```

### Gentleman Plugin Frontmatter (NEW)

```yaml
---
name: merge-checks
description: Audit code changes across 13 quality dimensions before or after merge
type: plugin
license: MIT
metadata:
  author: diego-marino
  upstream: https://github.com/dmarino/claude-toolshed
  version: "1.0"
dependencies:
  required:
    - git
    - bash
  optional:
    - node  # for i18n consistency checks
permissions:
  - "Bash(~/.claude/plugins/merge-checks/scripts/*:*)"
---
```

### Field Mapping

| Toolshed field | Gentleman plugin field | Action |
|---|---|---|
| `name` | `name` | Keep as-is |
| `description` | `description` | Keep as-is |
| `argument-hint` | -- | DROP (not supported by Gentleman) |
| `allowed-tools` | `permissions` | Map to `settings.json` format |
| `model` | -- | DROP (not supported by Gentleman) |
| `user-invocable` | -- | DROP (plugins are always invocable) |
| -- | `type: plugin` | ADD (distinguishes from skills) |
| -- | `license` | ADD (MIT for all toolshed ports) |
| -- | `metadata.upstream` | ADD (attribution to source repo) |
| -- | `dependencies.required` | ADD (runtime requirements) |
| -- | `dependencies.optional` | ADD (nice-to-have runtime deps) |

### Parser Changes

`parseSkillFrontmatter()` in `update.go` already does line-by-line YAML parsing for `name:` and `description:`. It needs two additions:

1. Accept both `SKILL.md` and `PLUGIN.md` filenames
2. Parse `type:` field and populate `SkillInfo.Type`

The parser does NOT need to parse `dependencies` or `permissions` -- those are consumed by the install logic, not the TUI display.

---

## 3. Path Resolution Rewrite Strategy

### Pattern: Dynamic `find` (toolshed)

Every toolshed SKILL.md and sub-skill uses this pattern to locate its own directory:

```bash
# BEFORE (toolshed)
SKILL_DIR="$(find "$HOME/.claude/plugins/cache" -type d -name "merge-checks" -path "*/skills/merge-checks" 2>/dev/null | head -1)"
[[ -z "$SKILL_DIR" ]] && SKILL_DIR="$(find "$HOME" -maxdepth 8 -type d -name "merge-checks" -path "*/skills/merge-checks" 2>/dev/null | head -1)"
```

### Pattern: Static path (Gentleman)

```bash
# AFTER (gentleman)
PLUGIN_DIR="$HOME/.claude/plugins/merge-checks"
```

### Concrete Rewrites

#### merge-checks PLUGIN.md

```markdown
## Git context

<!-- BEFORE -->
!`bash -c 'script="$(find $HOME/.claude/plugins/cache -name gather-context.sh -path "*/merge-checks/*" 2>/dev/null | head -1)"; ...; bash "$script" '"$ARGUMENTS"`

<!-- AFTER -->
!`bash "$HOME/.claude/plugins/merge-checks/scripts/gather-context.sh" "$ARGUMENTS"`
```

```markdown
## Phase 0 — Step 5

<!-- BEFORE -->
script="$(find $HOME/.claude/plugins/cache -name precompute.sh -path '*/merge-checks/*' 2>/dev/null | head -1)"
bash "$script" [ARGS]

<!-- AFTER -->
bash "$HOME/.claude/plugins/merge-checks/scripts/precompute.sh" [ARGS]
```

#### trim-md PLUGIN.md

```markdown
## Step 1: Resolve skill directory

<!-- BEFORE -->
SKILL_DIR="$(find "$HOME/.claude/plugins/cache" -type d -name "trim-md" -path "*/skills/trim-md" 2>/dev/null | head -1)"
[[ -z "$SKILL_DIR" ]] && SKILL_DIR="$(find "$HOME" -maxdepth 8 -type d -name "trim-md" -path "*/skills/trim-md" 2>/dev/null | head -1)"
echo "SKILL_DIR=$SKILL_DIR"

<!-- AFTER -->
PLUGIN_DIR="$HOME/.claude/plugins/trim-md"
echo "PLUGIN_DIR=$PLUGIN_DIR"
```

```markdown
## Step 2-3

<!-- BEFORE -->
bash "$SKILL_DIR/scripts/ensure-deps.sh"
bash "$SKILL_DIR/scripts/trim-md.sh" [--dry-run] <paths>

<!-- AFTER -->
bash "$PLUGIN_DIR/scripts/ensure-deps.sh"
bash "$PLUGIN_DIR/scripts/trim-md.sh" [--dry-run] <paths>
```

#### trim-md scripts/trim-md.sh internal reference

```bash
# BEFORE
REF_DIR="$SCRIPT_DIR/../reference"

# AFTER
REF_DIR="$SCRIPT_DIR/../config"
```

#### mermaid -- ALL sub-skills + agent

Every file that contains a `find` path resolution block gets the same rewrite:

```bash
# BEFORE (appears in: PLUGIN.md, mermaid-architect/SKILL.md, mermaid-config/SKILL.md,
#         mermaid-diagram/SKILL.md, mermaid-render/SKILL.md, mermaid-validate/SKILL.md,
#         agents/diagram-architect.md)
find "$HOME/.claude/plugins/cache" -type d -name "mermaid" -path "*/skills/mermaid" 2>/dev/null | head -1

# AFTER
PLUGIN_DIR="$HOME/.claude/plugins/mermaid"
```

All downstream references change from `$PLUGIN_DIR/scripts/...` to `$HOME/.claude/plugins/mermaid/scripts/...` or use the `$PLUGIN_DIR` variable set at the top of each file.

#### Script self-references (no change needed)

Scripts that use `$(cd "$(dirname "$0")" && pwd)` or `$(dirname "${BASH_SOURCE[0]}")` to find sibling scripts do NOT need path changes. These patterns work from any installed location because they're relative to the script file itself. Files affected:

- `merge-checks/scripts/precompute.sh`: `SCRIPTS="$(cd "$(dirname "$0")" && pwd)"` -- works
- `mermaid/scripts/ensure-deps.sh`: `SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"` -- works
- `trim-md/scripts/trim-md.sh`: `SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"` -- works

---

## 4. Dependency Handling Design

### Per-Plugin Dependency Matrix

| Plugin | Required | Optional | Node.js needed? |
|---|---|---|---|
| merge-checks | git, bash >= 4.0 | node (i18n checks) | No (core) |
| trim-md | npx (Node.js 18+) | -- | Yes |
| mermaid | node, npm | -- | Yes |

### ensure-deps.sh Behavior

```
┌─────────────────────┐
│ Plugin invoked       │
└──────────┬──────────┘
           │
           ▼
┌─────────────────────┐
│ ensure-deps.sh runs │
└──────────┬──────────┘
           │
     ┌─────┴─────┐
     │ Deps OK?  │
     └─────┬─────┘
       yes │     no
           │     │
           ▼     ▼
     ┌─────┐  ┌──────────────────┐
     │ OK  │  │ Print error msg  │
     │     │  │ + install hints  │
     └─────┘  │ + exit 1         │
              └──────────────────┘
```

### Node.js Module Isolation

```
~/.claude/plugins/mermaid/scripts/
  package.json           # { "dependencies": { "beautiful-mermaid": "^1.1.3" } }
  node_modules/          # created by ensure-deps.sh: npm install --prefix .
    beautiful-mermaid/
  extract_mermaid.js     # requires beautiful-mermaid from local node_modules
  render.js
  resilient_diagram.js
```

`trim-md` uses `npx --yes markdownlint-cli2` which auto-downloads on first use and caches in npm's global cache. No local `node_modules/` needed for trim-md.

---

## 5. TUI Skill Manager Changes

### Current Flow (skills only)

```
ScreenSkillMenu
  ├── Browse Skills     → fetchSkillCatalog() → list from ~/.gentleman/skills/ + ~/.claude/skills/
  ├── Install Skills    → symlink to ~/.claude/skills/<name>
  ├── Remove Skills     → remove symlink from ~/.claude/skills/<name>
  └── Update Catalog    → git pull on ~/.gentleman/skills/
```

### Extended Flow (skills + plugins)

```
ScreenSkillMenu
  ├── Browse Skills     → fetchSkillCatalog() → list from ~/.gentleman/skills/ + ~/.claude/skills/ + GentlemanClaude/plugins/
  ├── Install Skills    → skill: symlink to ~/.claude/skills/<name>
  │                       plugin: copy dir to ~/.claude/plugins/<name>/ + update settings.json
  ├── Remove Skills     → skill: remove symlink
  │                       plugin: rm -rf ~/.claude/plugins/<name>/ + remove settings.json entries
  └── Update Catalog    → git pull on ~/.gentleman/skills/ (plugins are repo-bundled, not in external catalog)
```

### Model Changes (`model.go`)

```go
// BEFORE
type SkillInfo struct {
    Name        string
    Description string
    Category    string // "curated" or "community"
    DirName     string
    FullPath    string
    Installed   bool
}

// AFTER
type SkillInfo struct {
    Name        string
    Description string
    Category    string // "curated", "community", "plugin", "local"
    Type        string // "skill" or "plugin"
    DirName     string
    FullPath    string
    Installed   bool
    Permissions []string // only for plugins: settings.json entries
}
```

### Catalog Fetch Changes (`update.go`)

```go
func fetchSkillCatalog() ([]SkillInfo, error) {
    // ... existing skill scanning ...

    // NEW: Scan GentlemanClaude/plugins/ for PLUGIN.md files
    repoDir := getRepoDir() // from installer context
    pluginsDir := filepath.Join(repoDir, "GentlemanClaude", "plugins")
    if entries, err := os.ReadDir(pluginsDir); err == nil {
        for _, entry := range entries {
            if !entry.IsDir() { continue }
            pluginFile := filepath.Join(pluginsDir, entry.Name(), "PLUGIN.md")
            if _, err := os.Stat(pluginFile); err != nil { continue }

            name, desc := parsePluginFrontmatter(pluginFile)
            if name == "" { name = entry.Name() }

            installed := isPluginInstalled(home, name) // checks ~/.claude/plugins/<name>/PLUGIN.md
            skills = append(skills, SkillInfo{
                Name:     name,
                Description: desc,
                Category: "plugin",
                Type:     "plugin",
                DirName:  entry.Name(),
                FullPath: filepath.Join(pluginsDir, entry.Name()),
                Installed: installed,
            })
        }
    }

    return skills, nil
}
```

### Install/Remove Logic Branching

```go
func installSkillSymlinks(skills []SkillInfo) ([]string, error) {
    for _, s := range skills {
        if s.Type == "plugin" {
            // Copy entire directory (not symlink)
            dst := filepath.Join(home, ".claude", "plugins", s.DirName)
            system.CopyDir(s.FullPath, dst)
            // Make scripts executable
            chmodScripts(dst)
            // Append permissions to settings.json
            appendPluginPermissions(s)
            logLines = append(logLines, fmt.Sprintf("📦 Installed plugin: %s → ~/.claude/plugins/%s/", s.Name, s.DirName))
        } else {
            // Existing symlink behavior for skills
            // ...
        }
    }
}
```

### Browse Screen Category Ordering

```go
func getSkillCategoryOrder(skills []SkillInfo) []string {
    // BEFORE: ["curated", "community", "local"]
    // AFTER:  ["curated", "community", "plugin", "local"]
}
```

The `skillCategoryHeader()` function returns:
- `"curated"` -> `"━━━ Curated Skills ━━━"`
- `"community"` -> `"━━━ Community Skills ━━━"`
- `"plugin"` -> `"━━━ Plugins ━━━"`   (NEW)
- `"local"` -> `"━━━ Local Skills ━━━"`

---

## 6. Symlink vs Copy Strategy

| Item | Install method | Location | Reason |
|---|---|---|---|
| Skills (curated/community) | Symlink | `~/.claude/skills/<name>` -> `~/.gentleman/skills/<path>` | Pure markdown, no execution, central updates via `git pull` |
| Plugins | Full copy | `~/.claude/plugins/<name>/` | Contains executable scripts, needs `chmod +x`, may generate `node_modules/`, must be self-contained |

### Why NOT symlinks for plugins

1. **Scripts need execute permission**: `chmod +x` on a symlink only affects the link, not the target (on some systems). Copying ensures we own the file.
2. **node_modules generated at install location**: `npm install --prefix` writes to the script's directory. With symlinks, this would write into the repo clone.
3. **Isolation**: If a user modifies a plugin (e.g., custom markdownlint config), changes don't affect the repo.
4. **Settings.json permissions use absolute paths**: `Bash(~/.claude/plugins/merge-checks/scripts/*:*)` -- the path must exist and be stable.

### Install Step in `installer.go`

The existing `setupCentralizedSkills()` function handles skills. Plugins are handled separately, triggered by TUI selection:

```go
// In the Claude AI tool install step (installer.go line ~1098)
if hasAITool(m.Choices.AITools, "claude") {
    // ... existing CLAUDE.md, settings.json, skills setup ...

    // NEW: ensure plugins directory exists
    system.EnsureDir(filepath.Join(claudeDir, "plugins"))
}
```

Plugin installation is deferred to the Skill Manager -- plugins are not installed by default during the initial installer run. The user explicitly chooses which plugins to install.

---

## 7. CLAUDE.md Auto-Load Table Extension

### Current Table (skills only)

```markdown
| Context                                | Read this file                         |
| -------------------------------------- | -------------------------------------- |
| React components, hooks, JSX           | `~/.claude/skills/react-19/SKILL.md`   |
| ...                                    | ...                                    |
```

### Extended Table (skills + plugins)

```markdown
### Plugin Detection

| Context                                       | Read this file                                |
| --------------------------------------------- | --------------------------------------------- |
| PR review, merge audit, code quality checks   | `~/.claude/plugins/merge-checks/PLUGIN.md`    |
| Markdown cleanup, lint, token optimization     | `~/.claude/plugins/trim-md/PLUGIN.md`         |
| Mermaid diagrams, architecture docs, SVG gen   | `~/.claude/plugins/mermaid/PLUGIN.md`         |
```

### How it works

Same mechanism as skills: Claude reads `CLAUDE.md` on session start, sees the context-to-file mapping, and loads the relevant `PLUGIN.md` when it detects matching context in the user's request.

---

## 8. settings.json Permission Entries

### New Entries

```json
{
  "permissions": {
    "allow": [
      "Bash(~/.claude/plugins/merge-checks/scripts/*:*)",
      "Bash(~/.claude/plugins/trim-md/scripts/*:*)",
      "Bash(~/.claude/plugins/mermaid/scripts/*:*)"
    ]
  }
}
```

These entries are appended to the existing `permissions.allow` array. They use glob patterns to allow execution of any script within the plugin's scripts directory.

### Append Strategy

During plugin install, the TUI reads `~/.claude/settings.json`, parses the JSON, checks if the permission entry already exists, and appends if missing. During remove, it removes matching entries.

The existing `settings.json` already has broad `Bash` permissions (e.g., `"Bash(npm:*)"`, `"Bash(npx:*)"`), so plugin scripts would technically work without explicit entries. However, explicit plugin permissions are added for:
1. Documentation: makes it clear which plugins have been installed
2. Future-proofing: if permissions become stricter
3. Removal: clean uninstall can remove exactly the entries it added
