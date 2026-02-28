# Design: mcp-expansion

**Change:** mcp-expansion
**Status:** Draft
**Spec:** [spec.md](./spec.md)

---

## File 1: `GentlemanClaude/mcp-servers.template.json`

### Current content (FULL)

```json
{
  "mcpServers": {
    "context7": {
      "type": "http",
      "url": "https://mcp.context7.com/mcp"
    },
    "mcp-atlassian": {
      "type": "stdio",
      "command": "uvx",
      "args": ["--python=3.12", "mcp-atlassian"],
      "env": {
        "JIRA_URL": "https://YOUR-COMPANY.atlassian.net",
        "JIRA_USERNAME": "YOUR_EMAIL@company.com",
        "JIRA_API_TOKEN": "YOUR_API_TOKEN_HERE"
      }
    },
    "figma": {
      "type": "stdio",
      "command": "npx",
      "args": ["-y", "@anthropic/mcp-figma"],
      "env": {
        "FIGMA_ACCESS_TOKEN": "YOUR_FIGMA_TOKEN_HERE"
      }
    }
  }
}
```

### New content (FULL)

```json
{
  "mcpServers": {
    "context7": {
      "type": "http",
      "url": "https://mcp.context7.com/mcp"
    },
    "mcp-atlassian": {
      "type": "stdio",
      "command": "uvx",
      "args": ["--python=3.12", "mcp-atlassian"],
      "env": {
        "JIRA_URL": "https://YOUR-COMPANY.atlassian.net",
        "JIRA_USERNAME": "YOUR_EMAIL@company.com",
        "JIRA_API_TOKEN": "YOUR_API_TOKEN_HERE"
      }
    },
    "figma": {
      "type": "stdio",
      "command": "npx",
      "args": ["-y", "@anthropic/mcp-figma"],
      "env": {
        "FIGMA_ACCESS_TOKEN": "YOUR_FIGMA_TOKEN_HERE"
      }
    },
    "brave-search": {
      "type": "stdio",
      "command": "npx",
      "args": ["-y", "@modelcontextprotocol/server-brave-search"],
      "env": {
        "BRAVE_API_KEY": "YOUR_BRAVE_API_KEY_HERE"
      }
    },
    "sentry": {
      "type": "http",
      "url": "https://mcp.sentry.dev/mcp"
    },
    "cloudflare": {
      "type": "http",
      "url": "https://mcp.cloudflare.com/mcp"
    }
  }
}
```

### Diff

```diff
       "env": {
         "FIGMA_ACCESS_TOKEN": "YOUR_FIGMA_TOKEN_HERE"
       }
+    },
+    "brave-search": {
+      "type": "stdio",
+      "command": "npx",
+      "args": ["-y", "@modelcontextprotocol/server-brave-search"],
+      "env": {
+        "BRAVE_API_KEY": "YOUR_BRAVE_API_KEY_HERE"
+      }
+    },
+    "sentry": {
+      "type": "http",
+      "url": "https://mcp.sentry.dev/mcp"
+    },
+    "cloudflare": {
+      "type": "http",
+      "url": "https://mcp.cloudflare.com/mcp"
     }
   }
 }
```

**Notes:**
- `brave-search` follows the exact `figma` pattern (npx STDIO + env placeholder)
- `sentry` and `cloudflare` follow the exact `context7` pattern (HTTP, 2 fields)
- Existing entries untouched

---

## File 2: `GentlemanOpenCode/opencode.json`

### Current MCP section

```json
  "mcp": {
    "context7": {
      "type": "remote",
      "url": "https://mcp.context7.com/mcp",
      "enabled": true
    },
    "engram": {
      "command": ["engram", "mcp"],
      "enabled": true,
      "type": "local"
    },
    "notion": {
      "type": "remote",
      "url": "https://mcp.notion.com/mcp",
      "enabled": true,
      "timeout": 5000
    },
    "jira": {
      "type": "remote",
      "url": "https://mcp.atlassian.com/v1/mcp",
      "enabled": true
    }
  },
```

### New MCP section

```json
  "mcp": {
    "context7": {
      "type": "remote",
      "url": "https://mcp.context7.com/mcp",
      "enabled": true
    },
    "engram": {
      "command": ["engram", "mcp"],
      "enabled": true,
      "type": "local"
    },
    "notion": {
      "type": "remote",
      "url": "https://mcp.notion.com/mcp",
      "enabled": true,
      "timeout": 5000
    },
    "jira": {
      "type": "remote",
      "url": "https://mcp.atlassian.com/v1/mcp",
      "enabled": true
    },
    "brave-search": {
      "type": "local",
      "command": ["npx", "-y", "@modelcontextprotocol/server-brave-search"],
      "enabled": true
    },
    "sentry": {
      "type": "remote",
      "url": "https://mcp.sentry.dev/mcp",
      "enabled": true
    },
    "cloudflare": {
      "type": "remote",
      "url": "https://mcp.cloudflare.com/mcp",
      "enabled": true
    }
  },
```

### Diff

```diff
     "jira": {
       "type": "remote",
       "url": "https://mcp.atlassian.com/v1/mcp",
       "enabled": true
+    },
+    "brave-search": {
+      "type": "local",
+      "command": ["npx", "-y", "@modelcontextprotocol/server-brave-search"],
+      "enabled": true
+    },
+    "sentry": {
+      "type": "remote",
+      "url": "https://mcp.sentry.dev/mcp",
+      "enabled": true
+    },
+    "cloudflare": {
+      "type": "remote",
+      "url": "https://mcp.cloudflare.com/mcp",
+      "enabled": true
     }
   },
```

**Notes:**
- `brave-search` uses `"type": "local"` + `"command"` array (matches `engram` pattern)
- No `env` field -- OpenCode assumes env vars are set externally
- `sentry` and `cloudflare` use `"type": "remote"` (matches `context7` and `jira` pattern)
- All new entries have `"enabled": true`

---

## File 3: `installer/internal/tui/update.go`

### Current MCP category block (lines ~1812-1822)

```go
	{
		ID: "mcp", Label: "MCP Servers", Icon: "\xf0\x9f\x94\x8c", IsAtomic: true,
		Items: []ModuleItem{
			{ID: "mcp-context7", Label: "Context7"},
			{ID: "mcp-engram", Label: "Engram"},
			{ID: "mcp-jira", Label: "Jira"},
			{ID: "mcp-atlassian", Label: "Atlassian"},
			{ID: "mcp-figma", Label: "Figma"},
			{ID: "mcp-notion", Label: "Notion"},
		},
	},
```

### New MCP category block

```go
	{
		ID: "mcp", Label: "MCP Servers", Icon: "\xf0\x9f\x94\x8c", IsAtomic: true,
		Items: []ModuleItem{
			{ID: "mcp-context7", Label: "Context7"},
			{ID: "mcp-engram", Label: "Engram"},
			{ID: "mcp-jira", Label: "Jira"},
			{ID: "mcp-atlassian", Label: "Atlassian"},
			{ID: "mcp-figma", Label: "Figma"},
			{ID: "mcp-notion", Label: "Notion"},
			{ID: "mcp-brave-search", Label: "Brave Search"},
			{ID: "mcp-sentry", Label: "Sentry"},
			{ID: "mcp-cloudflare", Label: "Cloudflare"},
		},
	},
```

### Diff

```diff
 			{ID: "mcp-notion", Label: "Notion"},
+			{ID: "mcp-brave-search", Label: "Brave Search"},
+			{ID: "mcp-sentry", Label: "Sentry"},
+			{ID: "mcp-cloudflare", Label: "Cloudflare"},
 		},
```

---

## File 4: `installer/internal/tui/ai_screens_test.go`

### Current (line ~142)

```go
		"mcp":      6,
```

### New

```go
		"mcp":      9,
```

---

## File 5: `docs/ai-framework-modules.md`

### Change 1: Header and overview (lines 1-3)

**Current:**
```markdown
Complete reference of all 203 modules across 6 categories available in the [project-starter-framework](...).
```

**New:**
```markdown
Complete reference of all 206 modules across 6 categories available in the [project-starter-framework](...).
```

### Change 2: Table of Contents MCP entry (line 15)

**Current:**
```markdown
- [MCP Servers (6 items)](#-mcp-servers-6-items)
```

**New:**
```markdown
- [MCP Servers (9 items)](#-mcp-servers-9-items)
```

### Change 3: Overview paragraph (line 21)

**Current:**
```markdown
The installer presents 203 individual modules organized into 6 categories.
```

**New:**
```markdown
The installer presents 206 individual modules organized into 6 categories.
```

### Change 4: How Features Work example (line 33)

**Current:**
```
🔌 MCP (2/6 selected)         →    --features=mcp
```

**New:**
```
🔌 MCP (2/9 selected)         →    --features=mcp
```

### Change 5: Categories Summary table (line 51)

**Current:**
```markdown
| MCP Servers | 🔌 | 6 | `mcp` | Yes |
```

**New:**
```markdown
| MCP Servers | 🔌 | 9 | `mcp` | Yes |
```

### Change 6: Total count (line 53)

**Current:**
```markdown
**Total: 203 modules**
```

**New:**
```markdown
**Total: 206 modules**
```

### Change 7: MCP section heading (line 466)

**Current:**
```markdown
## 🔌 MCP Servers (6 items)
```

**New:**
```markdown
## 🔌 MCP Servers (9 items)
```

### Change 8: MCP table -- add 3 rows after Notion (after line 477)

**Current (last row):**
```markdown
| `mcp-notion` | Notion | Notion workspace integration |
```

**New (append 3 rows):**
```markdown
| `mcp-notion` | Notion | Notion workspace integration |
| `mcp-brave-search` | Brave Search | Web search via Brave Search API (requires BRAVE_API_KEY) |
| `mcp-sentry` | Sentry | Error monitoring and performance tracking via Sentry |
| `mcp-cloudflare` | Cloudflare | Cloudflare Workers, Pages, and infrastructure management |
```

---

## File 6: `docs/ai-configuration.md`

### Change 1: Module count reference (line 79)

**Current:**
```markdown
The installer also optionally configures the **AI Framework** (Step 8) with 203 modules across 6 categories: hooks, commands, agents, skills, SDD, and MCP servers.
```

**New:**
```markdown
The installer also optionally configures the **AI Framework** (Step 8) with 206 modules across 6 categories: hooks, commands, agents, skills, SDD, and MCP servers.
```

### Change 2: Module count in "See also" reference (line 82)

**Current:**
```markdown
> See [AI Framework Module Registry](ai-framework-modules.md) for the complete list of 203 modules.
```

**New:**
```markdown
> See [AI Framework Module Registry](ai-framework-modules.md) for the complete list of 206 modules.
```

### Change 3: Claude Code template description (line 114)

**Current:**
```markdown
| `mcp-servers.template.json` | MCP server templates (Context7, Jira, Figma) |
```

**New:**
```markdown
| `mcp-servers.template.json` | MCP server templates (Context7, Jira, Figma, Brave Search, Sentry, Cloudflare) |
```

### Change 4: OpenCode MCP Integrations table (lines 311-316)

**Current:**
```markdown
| Server | Description |
|--------|-------------|
| **Context7** | Remote MCP for fetching up-to-date documentation |
| **Engram** | Local MCP backend for persistent SDD artifacts |

This is enabled by default and enhances the agent's ability to verify information with current docs.
```

**New:**
```markdown
| Server | Description |
|--------|-------------|
| **Context7** | Remote MCP for fetching up-to-date documentation |
| **Engram** | Local MCP backend for persistent SDD artifacts |
| **Notion** | Notion workspace integration |
| **Jira** | Jira integration via Atlassian MCP |
| **Brave Search** | Web search via Brave Search API (requires `BRAVE_API_KEY` env var) |
| **Sentry** | Error monitoring and performance tracking via Sentry (OAuth) |
| **Cloudflare** | Cloudflare Workers, Pages, and infrastructure management (OAuth) |

This is enabled by default and enhances the agent's ability to verify information with current docs and external services.
```

**Note:** The existing table was missing Notion and Jira which are already in `opencode.json`. This change brings the docs in sync with reality AND adds the 3 new servers.

---

## File 7: `docs/ai-tools-integration.md`

### Change 1: Category drill-down example (line 118)

**Current:**
```
  🔌 MCP (2/6 selected)
```

**New:**
```
  🔌 MCP (2/9 selected)
```

### Change 2: Module count in Step 8c description (line 106)

**Current:**
```markdown
The custom selection uses a **two-level drill-down** instead of a flat checkbox list, making it possible to navigate 203 individual modules across 6 categories.
```

**New:**
```markdown
The custom selection uses a **two-level drill-down** instead of a flat checkbox list, making it possible to navigate 206 individual modules across 6 categories.
```

---

## Design Decisions

| Decision | Rationale |
|----------|-----------|
| Append new entries after existing ones | Preserves existing order, minimizes diff, no behavioral impact |
| `brave-search` as key (hyphenated) | Matches existing patterns: `mcp-atlassian`, `context7`. Consistent with npm package naming |
| No `timeout` on new remote entries | Only `notion` has a timeout; Sentry and Cloudflare don't need one (OAuth flow handles latency) |
| Fix OpenCode MCP docs table | The existing table only lists 2 of 4 current servers. Fixing this while adding 3 new ones is the right thing to do |
| No changes to `IsAtomic: true` | Per proposal, atomic behavior stays. All MCP servers are installed together |
