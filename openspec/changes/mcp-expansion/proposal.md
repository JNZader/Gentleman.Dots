# Proposal: mcp-expansion

## Intent

Expand the MCP server catalog from 6 to 9 servers by adding Brave Search, Sentry, and Cloudflare. All three are production-ready, publicly available MCP endpoints that follow patterns already established in the codebase. The goal is to give users broader out-of-the-box connectivity to common developer services without changing the installation model.

## Scope

### In scope

- Add **Brave Search** (STDIO, npx, requires `BRAVE_API_KEY`), **Sentry** (HTTP remote, OAuth), and **Cloudflare** (HTTP remote, OAuth) to:
  - `GentlemanClaude/mcp-servers.template.json` (Claude Code format: `type: http` / `type: stdio`)
  - `GentlemanOpenCode/opencode.json` (OpenCode format: `type: remote` / `type: local`)
  - `installer/internal/tui/update.go` (TUI module registry `ModuleItem` entries)
- Update `installer/internal/tui/ai_screens_test.go` expected MCP count from 6 to 9.
- Update documentation counts and tables in:
  - `docs/ai-framework-modules.md` (total 203 -> 206, MCP section 6 -> 9, add 3 table rows)
  - `docs/ai-configuration.md` (MCP description and table)
  - `docs/ai-tools-integration.md` (TUI example counts)

### Out of scope

- Changes to `project-starter-framework` (`setup-global.sh` script) -- separate repo, separate PR.
- Changing `IsAtomic: true` behavior for the MCP category. All-or-nothing installation stays as-is.
- Any new MCP servers beyond the three listed (GitHub, Playwright, Sequential Thinking, PostgreSQL, Filesystem were explicitly evaluated and rejected in the explore phase).
- Per-project MCP configuration or selective MCP installation -- future work.

## Approach

1. **JSON config entries** -- Sentry and Cloudflare are HTTP remote servers (identical pattern to Context7: 2 fields each). Brave Search is STDIO via npx (identical pattern to Figma: command + args + env).
2. **TUI registry** -- Append 3 `ModuleItem` structs to the `mcp` category slice in `update.go`. No new categories, no structural changes.
3. **Tests** -- Single numeric assertion change (`6` -> `9`) in the expected counts map.
4. **Docs** -- Mechanical updates: increment counts, add rows to markdown tables. No prose restructuring.

The two config files use different MCP schemas:

| Field | Claude Code (`mcp-servers.template.json`) | OpenCode (`opencode.json`) |
|-------|-------------------------------------------|---------------------------|
| HTTP server type | `"type": "http"` | `"type": "remote"` |
| STDIO server type | `"type": "stdio"` | `"type": "local"` |
| URL field | `"url": "..."` | `"url": "..."` |
| Command field | `"command": "npx"` | `"command": ["npx", ...]` |
| Args | `"args": ["-y", "@pkg"]` | (inlined in command array) |
| Env vars | `"env": { "KEY": "VALUE" }` | Not used (env assumed external) |

Both files must be updated with the correct format for each server.

## Key Decisions

| Decision | Rationale |
|----------|-----------|
| Only 3 new servers | Evaluated 8 candidates. Rejected GitHub MCP (redundant with `gh` CLI), Playwright (per-project concern), Sequential Thinking (no proven value), PostgreSQL (deprecated), Filesystem (redundant with Claude Code built-in). |
| Keep `IsAtomic: true` | Selective MCP installation adds complexity to `setup-global.sh` for minimal user benefit at this scale. Revisit if catalog exceeds ~15 servers. |
| OAuth servers need no env vars | Sentry and Cloudflare use OAuth browser flow -- no API keys in config templates. This matches the Context7 pattern exactly. |
| Brave Search uses placeholder env var | Same pattern as Figma (`FIGMA_ACCESS_TOKEN`) -- user replaces `YOUR_BRAVE_API_KEY_HERE` after install. |

## Risks & Mitigations

| Risk | Impact | Mitigation |
|------|--------|------------|
| MCP endpoint URLs change | Servers stop connecting | URLs are from official provider docs (sentry.dev, cloudflare.com). Pin to current stable endpoints. Breakage is detectable and fixable in a single-line change. |
| Brave Search API key friction | Users skip Brave Search because it requires a key | Document how to get a free API key in the MCP servers table description. Same friction already exists for Figma and Atlassian. |
| OpenCode `command` array format incorrect | Brave Search fails in OpenCode | Verified against existing `engram` entry which uses `"command": ["engram", "mcp"]` format. Follow same pattern: `"command": ["npx", "-y", "@modelcontextprotocol/server-brave-search"]`. |
| Test hardcoding of MCP category index | Tests reference `SelectedModuleCategory = 5` (MCP is last). Adding categories later would break. | Not a concern for THIS change -- we are adding items within MCP, not new categories. Index stays 5. |

## Dependencies

| Dependency | Type | Status |
|------------|------|--------|
| `@modelcontextprotocol/server-brave-search` npm package | External runtime | Available on npm, actively maintained |
| `https://mcp.sentry.dev/mcp` endpoint | External service | Live, documented by Sentry |
| `https://mcp.cloudflare.com/mcp` endpoint | External service | Live, documented by Cloudflare |
| `project-starter-framework` `setup-global.sh` | Separate repo | Must be updated to copy the new JSON entries. OUT OF SCOPE for this PR but required for end-to-end functionality. |

## Estimated Effort

**Size: XS** (< 1 hour)

- 7 files changed, all mechanical edits
- No new logic, no new patterns, no new dependencies in the Go codebase
- Every addition follows an existing pattern verbatim
- Risk of regression: minimal (one test assertion change, rest is config/docs)

| File | Change type | Lines |
|------|------------|-------|
| `GentlemanClaude/mcp-servers.template.json` | Add 3 JSON entries | ~20 |
| `GentlemanOpenCode/opencode.json` | Add 3 JSON entries | ~15 |
| `installer/internal/tui/update.go` | Add 3 `ModuleItem` structs | ~3 |
| `installer/internal/tui/ai_screens_test.go` | Change `6` to `9` | ~1 |
| `docs/ai-framework-modules.md` | Update counts + add 3 table rows | ~10 |
| `docs/ai-configuration.md` | Update MCP description/table | ~5 |
| `docs/ai-tools-integration.md` | Update example counts | ~2 |
