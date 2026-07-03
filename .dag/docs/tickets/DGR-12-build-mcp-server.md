---
id: DGR-12
issueType: story
status: done
tags:
  - phase-1
  - mcp
  - agent-integration
parent: DGR-5
blockedBy:
  - DGR-9
  - DGR-10
  - DGR-11
  - DGR-16
---

# Build MCP server

## Description

As an AI agent (Claude Code / OpenCode), I want to fetch Dagger ticket context via MCP tools so I can generate code informed by full decision history without leaving the chat.

## Acceptance criteria

1. `POST /mcp` endpoint on the REST API that accepts JSON-RPC 2.0 requests
2. Tool dispatch logic behind a shared `ToolService` interface
3. Tools implemented:
   - `get_issue` — fetch ticket with full context by display_id
   - `get_doc` — fetch individual document by display_id
   - `list_issues` — list tickets with optional status filter
   - `update_issue_status` — updates issue status
4. Each tool calls the database directly via the existing `internal/issues` and `internal/docs` packages
5. Authentication via existing `Authorization: Bearer <key>` header (uses existing auth middleware)
6. Compatible with standard MCP HTTP transport config in `.mcp.json` or Claude Desktop

## User configuration

```json
{
  "mcpServers": {
    "dagger": {
      "type": "http",
      "url": "https://api.dagger.sh/mcp",
      "headers": {
        "Authorization": "Bearer dgr_abc123"
      }
    }
  }
}
```

## Technical notes

- Single `POST /mcp` route on the existing chi router, behind the existing auth middleware
- Tool logic lives in `internal/mcp/` — dispatch, tool definitions, JSON-RPC types
- `ToolService` interface extracted so dispatch is independent of data backend
- HTTP transport uses `DBService` (calls `internal/issues` and `internal/docs` directly — no extra HTTP hop)
- JSON-RPC 2.0 structured request/response format
- See `.dag/docs/plan/dagger-plan/05-agent-integration.md` for original tool definitions

## Files

| File | Purpose |
|------|---------|
| `internal/mcp/mcp.go` | JSON-RPC 2.0 types, tool schema types, `ListTools()`, `ToolService` interface |
| `internal/mcp/server.go` | `Server` dispatch, `tools/list` + `tools/call`, `Serve` stdio loop |
| `internal/mcp/dbservice.go` | `DBService` — `ToolService` impl backed by `internal/issues` / `internal/docs` |
| `cmd/api/main.go` | Add `POST /mcp` route |
| `.mcp.json` | Example user config |

## Related docs

- Agent integration: `.dag/docs/plan/dagger-plan/05-agent-integration.md`

---

*This ticket was created by opencode and reviewed by a human before publishing.*
