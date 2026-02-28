# Spec: mcp-expansion

**Change:** mcp-expansion
**Status:** Draft
**Proposal:** [proposal.md](./proposal.md)

---

## Requirements

### REQ-MCP-01: Brave Search server config (STDIO)

Add a Brave Search MCP server entry to the Claude Code template (`mcp-servers.template.json`).

- **Transport:** STDIO via `npx`
- **Package:** `@modelcontextprotocol/server-brave-search`
- **Args:** `["-y", "@modelcontextprotocol/server-brave-search"]`
- **Env:** `BRAVE_API_KEY` with placeholder value `YOUR_BRAVE_API_KEY_HERE`
- **Pattern match:** Identical to the existing `figma` entry (npx + env var placeholder)

**Acceptance criteria:**
- JSON is valid and parseable
- Entry key is `"brave-search"`
- `type` is `"stdio"`, `command` is `"npx"`
- `env.BRAVE_API_KEY` exists with placeholder value

### REQ-MCP-02: Sentry server config (HTTP, OAuth)

Add a Sentry MCP server entry to the Claude Code template.

- **Transport:** HTTP remote
- **URL:** `https://mcp.sentry.dev/mcp`
- **Auth:** OAuth browser flow (no API key needed in config)
- **Pattern match:** Identical to the existing `context7` entry (2 fields: type + url)

**Acceptance criteria:**
- Entry key is `"sentry"`
- `type` is `"http"`, `url` is `"https://mcp.sentry.dev/mcp"`
- No `env`, `command`, or `args` fields

### REQ-MCP-03: Cloudflare server config (HTTP, OAuth)

Add a Cloudflare MCP server entry to the Claude Code template.

- **Transport:** HTTP remote
- **URL:** `https://mcp.cloudflare.com/mcp`
- **Auth:** OAuth browser flow (no API key needed in config)
- **Pattern match:** Identical to `context7` and `sentry`

**Acceptance criteria:**
- Entry key is `"cloudflare"`
- `type` is `"http"`, `url` is `"https://mcp.cloudflare.com/mcp"`
- No `env`, `command`, or `args` fields

### REQ-MCP-04: Claude Code template update

The file `GentlemanClaude/mcp-servers.template.json` must contain exactly 6 entries after the change: context7, mcp-atlassian, figma (existing) + brave-search, sentry, cloudflare (new).

**Acceptance criteria:**
- File has 6 top-level keys under `mcpServers`
- Existing 3 entries are unchanged (byte-identical)
- New entries follow the ordering: existing first, then brave-search, sentry, cloudflare
- JSON is valid (parseable with `jq .`)

### REQ-MCP-05: OpenCode config update (different format)

Add Brave Search, Sentry, and Cloudflare to `GentlemanOpenCode/opencode.json` under the `mcp` key.

**Format differences from Claude Code:**
- HTTP servers use `"type": "remote"` (not `"http"`)
- STDIO servers use `"type": "local"` (not `"stdio"`)
- STDIO commands use `"command": ["npx", "-y", "@pkg"]` array (not separate `command` + `args`)
- No `env` field (environment variables are assumed external)
- Each entry has `"enabled": true`

**Acceptance criteria:**
- `mcp.brave-search` has `type: "local"`, `command: ["npx", "-y", "@modelcontextprotocol/server-brave-search"]`, `enabled: true`
- `mcp.sentry` has `type: "remote"`, `url: "https://mcp.sentry.dev/mcp"`, `enabled: true`
- `mcp.cloudflare` has `type: "remote"`, `url: "https://mcp.cloudflare.com/mcp"`, `enabled: true`
- Existing 4 entries (context7, engram, notion, jira) are unchanged
- JSON is valid

### REQ-MCP-06: TUI registry update

Add 3 `ModuleItem` entries to the `mcp` category in `installer/internal/tui/update.go`.

**Acceptance criteria:**
- Items appended after `{ID: "mcp-notion", Label: "Notion"}`
- New items: `{ID: "mcp-brave-search", Label: "Brave Search"}`, `{ID: "mcp-sentry", Label: "Sentry"}`, `{ID: "mcp-cloudflare", Label: "Cloudflare"}`
- `IsAtomic: true` remains on the category (not changed)
- Test expectation in `ai_screens_test.go` updated from `6` to `9`

### REQ-MCP-07: Documentation updates

Update 3 documentation files with new counts and table rows.

**REQ-MCP-07a: `docs/ai-framework-modules.md`**
- Total module count: 203 -> 206 (in heading, overview paragraph, and table references)
- MCP section heading: "6 items" -> "9 items"
- MCP row in Categories Summary table: `6` -> `9`
- Total row: `203` -> `206`
- Add 3 table rows to MCP section: Brave Search, Sentry, Cloudflare
- Category drill-down example: `2/6 selected` -> `2/9 selected`

**REQ-MCP-07b: `docs/ai-configuration.md`**
- MCP Integrations table in OpenCode section: add Sentry, Cloudflare, Brave Search rows
- Update `mcp-servers.template.json` description to include new servers
- Update total module count references (203 -> 206)

**REQ-MCP-07c: `docs/ai-tools-integration.md`**
- Category drill-down example: `2/6 selected` -> `2/9 selected`
- Module count references: 203 -> 206

---

## Scenarios

### SCN-01: User selects Complete preset

**Given** a user running the TUI installer
**When** they reach Step 8b (Preset Selection) and choose "Complete"
**Then** the `mcp` feature flag is included in `--features=`
**And** all 9 MCP servers are installed (context7, mcp-atlassian, figma, engram, notion, jira, brave-search, sentry, cloudflare)
**And** the user sees "MCP (9 items)" if they inspect the category

### SCN-02: User inspects MCP category items in Custom mode

**Given** a user in Step 8c (Custom Category Drill-Down)
**When** they navigate to the MCP category
**Then** they see 9 items listed: Context7, Engram, Jira, Atlassian, Figma, Notion, Brave Search, Sentry, Cloudflare
**And** all items show checkboxes (since MCP is atomic, toggling any enables all)

### SCN-03: Non-interactive with --ai-modules=mcp

**Given** a CI/CD pipeline running:
```bash
gentleman.dots --non-interactive --shell=zsh --ai-tools=claude --ai-framework --ai-modules=mcp
```
**When** the installer runs
**Then** `setup-global.sh --features=mcp` is executed
**And** both `mcp-servers.template.json` (6 servers) and `opencode.json` MCP section (7 servers) are deployed
**And** no errors occur from the new JSON entries

### SCN-04: Brave Search requires API key

**Given** a user who installed with MCP enabled
**When** they open Claude Code and Brave Search MCP is configured
**Then** the config contains `BRAVE_API_KEY: "YOUR_BRAVE_API_KEY_HERE"`
**And** the user must replace the placeholder with a real key for Brave Search to function
**And** other MCP servers (Sentry, Cloudflare) work immediately via OAuth

### SCN-05: OpenCode format differs from Claude Code

**Given** both config files are deployed
**When** comparing Sentry config between the two files
**Then** Claude Code has `"type": "http"` while OpenCode has `"type": "remote"`
**And** Brave Search in Claude Code has `"command": "npx", "args": [...]` while OpenCode has `"command": ["npx", "-y", "..."]`
**And** OpenCode entries have `"enabled": true` while Claude Code entries do not

---

## Traceability

| Requirement | Files | Tasks |
|-------------|-------|-------|
| REQ-MCP-01 | `GentlemanClaude/mcp-servers.template.json` | TASK-01 |
| REQ-MCP-02 | `GentlemanClaude/mcp-servers.template.json` | TASK-01 |
| REQ-MCP-03 | `GentlemanClaude/mcp-servers.template.json` | TASK-01 |
| REQ-MCP-04 | `GentlemanClaude/mcp-servers.template.json` | TASK-01 |
| REQ-MCP-05 | `GentlemanOpenCode/opencode.json` | TASK-02 |
| REQ-MCP-06 | `installer/internal/tui/update.go`, `ai_screens_test.go` | TASK-03, TASK-04 |
| REQ-MCP-07a | `docs/ai-framework-modules.md` | TASK-05 |
| REQ-MCP-07b | `docs/ai-configuration.md` | TASK-06 |
| REQ-MCP-07c | `docs/ai-tools-integration.md` | TASK-07 |
