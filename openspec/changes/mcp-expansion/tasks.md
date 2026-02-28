# Tasks: mcp-expansion

**Change:** mcp-expansion
**Status:** Draft
**Spec:** [spec.md](./spec.md)
**Design:** [design.md](./design.md)

---

## TASK-01: Update mcp-servers.template.json

**File:** `GentlemanClaude/mcp-servers.template.json`
**Requirements:** REQ-MCP-01, REQ-MCP-02, REQ-MCP-03, REQ-MCP-04
**Estimated lines changed:** ~20

**Change:**
Add 3 new entries to the `mcpServers` object after the existing `figma` entry:

1. `"brave-search"` -- STDIO server via npx with `BRAVE_API_KEY` env placeholder
2. `"sentry"` -- HTTP server, `https://mcp.sentry.dev/mcp`
3. `"cloudflare"` -- HTTP server, `https://mcp.cloudflare.com/mcp`

**Exact edit:** Add a comma after the closing `}` of `figma`, then insert the 3 new JSON blocks as shown in design.md File 1.

**Validation:** `jq . GentlemanClaude/mcp-servers.template.json` exits 0 and shows 6 keys under `mcpServers`.

---

## TASK-02: Update opencode.json

**File:** `GentlemanOpenCode/opencode.json`
**Requirements:** REQ-MCP-05
**Estimated lines changed:** ~15

**Change:**
Add 3 new entries to the `mcp` object after the existing `jira` entry:

1. `"brave-search"` -- `type: "local"`, `command: ["npx", "-y", "@modelcontextprotocol/server-brave-search"]`, `enabled: true`
2. `"sentry"` -- `type: "remote"`, `url: "https://mcp.sentry.dev/mcp"`, `enabled: true`
3. `"cloudflare"` -- `type: "remote"`, `url: "https://mcp.cloudflare.com/mcp"`, `enabled: true`

**Exact edit:** Add a comma after the closing `}` of `jira`, then insert the 3 new JSON blocks as shown in design.md File 2.

**Validation:** `jq .mcp GentlemanOpenCode/opencode.json` exits 0 and shows 7 keys.

---

## TASK-03: Update TUI module registry

**File:** `installer/internal/tui/update.go`
**Requirements:** REQ-MCP-06
**Estimated lines changed:** ~3

**Change:**
Append 3 `ModuleItem` structs to the `mcp` category's `Items` slice, after `{ID: "mcp-notion", Label: "Notion"}`:

```go
{ID: "mcp-brave-search", Label: "Brave Search"},
{ID: "mcp-sentry", Label: "Sentry"},
{ID: "mcp-cloudflare", Label: "Cloudflare"},
```

**Validation:** `go build ./installer/...` compiles without errors.

---

## TASK-04: Update test expectations

**File:** `installer/internal/tui/ai_screens_test.go`
**Requirements:** REQ-MCP-06
**Estimated lines changed:** ~1

**Change:**
In `TestModuleCategoriesItemCount`, change the expected MCP count:

```go
// Before
"mcp":      6,
// After
"mcp":      9,
```

**Validation:** `go test ./installer/internal/tui/ -run TestModuleCategoriesItemCount` passes.

---

## TASK-05: Update ai-framework-modules.md

**File:** `docs/ai-framework-modules.md`
**Requirements:** REQ-MCP-07a
**Estimated lines changed:** ~10

**Changes (8 edits):**

1. Line 1: `203 modules` -> `206 modules` (header paragraph)
2. Line 15: TOC entry `(6 items)` -> `(9 items)`
3. Line 21: Overview `203 individual modules` -> `206 individual modules`
4. Line 33: Example `2/6 selected` -> `2/9 selected`
5. Line 51: Categories Summary table MCP row `6` -> `9`
6. Line 53: Total `203` -> `206`
7. Line 466: Section heading `(6 items)` -> `(9 items)`
8. After line 477: Add 3 table rows (Brave Search, Sentry, Cloudflare)

**Validation:** All counts are internally consistent; no broken markdown table formatting.

---

## TASK-06: Update ai-configuration.md

**File:** `docs/ai-configuration.md`
**Requirements:** REQ-MCP-07b
**Estimated lines changed:** ~10

**Changes (4 edits):**

1. Line 79: `203 modules` -> `206 modules`
2. Line 82: `203 modules` -> `206 modules`
3. Line 114: MCP template description add `Brave Search, Sentry, Cloudflare`
4. Lines 311-316: Expand MCP Integrations table from 2 rows to 7 rows (add Notion, Jira -- already in config but missing from docs -- plus Brave Search, Sentry, Cloudflare)

**Validation:** Table renders correctly in markdown preview.

---

## TASK-07: Update ai-tools-integration.md

**File:** `docs/ai-tools-integration.md`
**Requirements:** REQ-MCP-07c
**Estimated lines changed:** ~2

**Changes (2 edits):**

1. Line 106: `203 individual modules` -> `206 individual modules`
2. Line 118: `2/6 selected` -> `2/9 selected`

**Validation:** Counts match ai-framework-modules.md.

---

## Execution Order

All tasks are independent and can be executed in any order. No task blocks another.

Recommended order for review convenience:
1. TASK-01 + TASK-02 (config files -- the actual MCP changes)
2. TASK-03 + TASK-04 (Go code + test)
3. TASK-05 + TASK-06 + TASK-07 (docs)

**Total estimated lines changed:** ~61
