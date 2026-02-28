# Proposal: toolshed-plugins

## Intent

GentlemanClaude's skill system is pure markdown -- 23 `.md` files that provide instructions, conventions, and workflows. This is by design: skills are stateless knowledge that Claude reads and follows.

However, diego marino's [claude-toolshed](https://github.com/dmarino/claude-toolshed) (MIT license) offers three **executable plugins** that extend Claude's capabilities with scripts, hooks, and external tool integrations:

- **merge-checks**: 13 quality checks via shell scripts, 3 sub-agents for automated PR review
- **trim-md**: Markdown linting via `markdownlint-cli2` with a PostToolUse hook for automatic cleanup
- **mermaid**: Diagram generation via `beautiful-mermaid`, 6 sub-skills, 7 specialist agent types

These don't fit the "pure markdown skill" model. They have **bash scripts, Node.js dependencies, git hooks, and runtime behavior**. Cramming them into `skills/` would violate the existing architecture. We need a new category.

The goal is to port all three plugins into a new `GentlemanClaude/plugins/` directory, adapting them to Gentleman conventions (path resolution, settings.json permissions, TUI discovery) while preserving their full functionality.

## Scope

### In Scope

1. **New `GentlemanClaude/plugins/` directory** as a first-class category alongside `skills/`
2. **Port merge-checks** (22+ files, ~75KB): shell scripts for 13 quality checks, 3 sub-agent definitions, SKILL.md entry point. Requires only git + bash.
3. **Port trim-md** (5 files, ~8KB): markdownlint wrapper with ensure-deps pattern. Requires Node.js for `markdownlint-cli2`.
4. **Port mermaid** (50+ files, ~350KB): diagram generation with `beautiful-mermaid`. Requires Node.js. 6 sub-skills, 7 specialist types.
5. **Path resolution rewrite**: toolshed uses dynamic `find $HOME/.claude/plugins/cache` patterns. Rewrite to static `~/.claude/plugins/<name>/` paths.
6. **Frontmatter adaptation**: toolshed uses `allowed-tools`, `argument-hint`, `model` fields not present in Gentleman. Map or drop as appropriate.
7. **settings.json updates**: add necessary `permissions.allow` entries for plugin scripts (e.g., `Bash(~/.claude/plugins/merge-checks/scripts/*:*)`)
8. **TUI Skill Manager extension**: add plugin discovery/install/remove alongside existing skill management (ScreenSkillMenu already exists per project-init-and-skill-manager change)
9. **Module registry update**: `docs/ai-framework-modules.md` count 203 -> 206 (3 new plugins)
10. **Dependency management**: `ensure-deps.sh` pattern for Node.js-dependent plugins (trim-md, mermaid)

### Out of Scope

- Modifying the upstream claude-toolshed repo
- Writing a generic plugin runtime/loader (plugins are self-contained, not dynamically loaded)
- Adding a hook system to GentlemanClaude (see Key Decisions)
- Porting any toolshed plugins beyond these three (scope is explicitly merge-checks, trim-md, mermaid)
- Auto-updating plugins from upstream (manual port, Gentleman-owned copies)
- Windows support (plugins are shell/Node.js, same as the rest of the installer)

## Approach

### Plugin Directory Structure

```
GentlemanClaude/
  skills/           # existing: pure .md files (23 skills)
  plugins/          # NEW: executable plugins with scripts + deps
    merge-checks/
      PLUGIN.md     # entry point (analogous to SKILL.md)
      scripts/      # 18 bash scripts for quality checks
      agents/       # 3 sub-agent definitions
      config/       # default config
    trim-md/
      PLUGIN.md
      scripts/      # wrapper + ensure-deps.sh
      config/       # markdownlint config
    mermaid/
      PLUGIN.md
      skills/       # 6 sub-skills for diagram types
      agents/       # 7 specialist definitions
      scripts/      # ensure-deps.sh + render pipeline
      templates/    # diagram templates
```

### Naming Convention

Plugins use `PLUGIN.md` as their entry point (not `SKILL.md`) to make the distinction explicit. CLAUDE.md will reference plugins separately from skills in its auto-load table.

### Path Resolution Strategy

Toolshed uses dynamic `find` to locate files. We rewrite ALL paths to be static and relative to `~/.claude/plugins/<name>/`. This means:
- No runtime `find` calls
- Paths are predictable and testable
- Installer copies plugins to the fixed location

### Dependency Management

Both trim-md and mermaid need Node.js packages. Approach:
1. Each plugin ships an `ensure-deps.sh` that checks for and installs its npm dependency
2. `ensure-deps.sh` runs on first plugin invocation (lazy install)
3. Dependencies install to `~/.claude/plugins/<name>/node_modules/` (isolated, not global)
4. If Node.js is not available, plugin degrades gracefully with a clear error message

### TUI Integration

The Skill Manager (from the project-init-and-skill-manager change) gets a new section:
- ScreenSkillMenu gains a "Plugins" option alongside Browse/Install/Remove
- Plugins are listed from `GentlemanClaude/plugins/` (same pattern as skills from Gentleman-Skills repo)
- Install copies the plugin directory to `~/.claude/plugins/<name>/`
- Remove deletes the plugin directory
- Plugin state (installed/not) checked via directory existence

### settings.json Permissions

Plugins need script execution permission. Add to `permissions.allow`:
```json
"Bash(~/.claude/plugins/merge-checks/scripts/*:*)",
"Bash(~/.claude/plugins/trim-md/scripts/*:*)",
"Bash(~/.claude/plugins/mermaid/scripts/*:*)"
```

## Key Decisions

| Decision | Choice | Rationale |
|----------|--------|-----------|
| Separate `plugins/` directory | Yes | Skills are pure markdown; plugins have scripts, deps, configs. Mixing them violates the existing architecture. |
| `PLUGIN.md` entry point | Yes | Makes plugin vs skill distinction explicit at the filesystem level. |
| Drop trim-md PostToolUse hook | Yes, for now | GentlemanClaude has NO hook system. Adding one is a separate, larger change. trim-md works fine as a manual-invoke tool ("run markdownlint on this file"). The hook is a convenience, not a requirement. |
| Static paths over dynamic find | Yes | Predictable, testable, no runtime overhead. Toolshed's `find` pattern was for flexibility we don't need since we control the install location. |
| Local node_modules per plugin | Yes | Avoids global npm pollution. Each plugin is self-contained. Follows the same isolation principle as the existing framework. |
| Adapt frontmatter, don't copy | Yes | Drop `allowed-tools` (we use settings.json permissions), drop `argument-hint` and `model` (not supported by Gentleman). Keep only fields that map to our conventions. |
| TUI integration depends on project-init-and-skill-manager | Yes | The Skill Manager screens are being built in that change. We extend them, not duplicate them. |

## Risks & Mitigations

| Risk | Likelihood | Impact | Mitigation |
|------|------------|--------|------------|
| Node.js not installed on user's machine | Medium | trim-md and mermaid won't work | `ensure-deps.sh` checks for `node`/`npm` first, prints clear error with install instructions. merge-checks works without Node.js. |
| Mermaid plugin is large (350KB+, 50+ files) | Low | Bloats the GentlemanClaude directory | Plugin is only copied when explicitly installed via TUI. Not included by default. |
| Upstream toolshed changes after our port | Medium | Our copy drifts from upstream | We own the Gentleman copies. Upstream is MIT, 1-day-old repo, single author. We can cherry-pick useful changes manually. |
| Path rewrite introduces bugs in script cross-references | Medium | Scripts fail at runtime | Comprehensive path audit during port. Each script tested individually after rewrite. |
| TUI Skill Manager change not yet merged | High | Plugin TUI integration blocked | Plugin port (files + scripts) can proceed independently. TUI integration is a separate task that depends on Skill Manager being available. |
| Permission entries in settings.json need manual sync | Low | New scripts not executable | Plugin install step appends permission entries. Document the required entries in PLUGIN.md as fallback. |

## Dependencies

| Dependency | Type | Status |
|------------|------|--------|
| `project-init-and-skill-manager` change | Internal | In progress -- TUI plugin screens depend on this |
| `claude-toolshed` repo (MIT) | External | Available -- source for the port |
| `bash` >= 4.0 | Runtime | Required by merge-checks scripts |
| `git` | Runtime | Required by merge-checks (diff analysis, branch detection) |
| `node` >= 18 + `npm` | Runtime | Required by trim-md (`markdownlint-cli2`) and mermaid (`beautiful-mermaid`). Optional -- merge-checks works without it. |
| GentlemanClaude `settings.json` | Internal | Must be updated with plugin script permissions |
| GentlemanClaude `CLAUDE.md` | Internal | Must be updated with plugin auto-load table |
| `docs/ai-framework-modules.md` | Internal | Module count update (203 -> 206) |

## Estimated Effort

| Phase | Effort | Notes |
|-------|--------|-------|
| Plugin directory structure + conventions | Small (1-2h) | Define PLUGIN.md format, directory layout |
| merge-checks port | Medium (3-4h) | 22 files, path rewrites, script testing |
| trim-md port | Small (1-2h) | 5 files, straightforward, drop hook |
| mermaid port | Large (4-6h) | 50+ files, complex sub-skill structure, template system |
| Path resolution rewrite (all 3) | Medium (2-3h) | Audit every `find`/dynamic path, rewrite to static |
| settings.json + CLAUDE.md updates | Small (1h) | Permission entries, auto-load table |
| TUI plugin screens | Medium (2-3h) | Extends Skill Manager, depends on that change landing |
| Module registry update | Trivial (15min) | Count update + 3 new entries |
| Testing + validation | Medium (2-3h) | Each plugin end-to-end, ensure-deps flows |
| **Total** | **~16-24 hours** | Can parallelize: merge-checks and trim-md have no interdependency |
