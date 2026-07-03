---
id: DGR-12
issueType: story
status: open
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

1. `dagger-mcp` binary that communicates via stdio JSON-RPC 2.0
2. Tools implemented:
   - `get_issue` — fetch ticket with full context by displayId
   - `get_doc` — fetch individual document by id/displayId
   - `list_issues` — list tickets with optional status filter
   - `update_issue_status` — updates issue status
   - `search_docs` — search across docs and tickets (basic text match)
3. Each tool calls the Dagger REST API internally (not direct DB access)
4. Configured via `--api-key` flag and optionally `--api-url`
5. Compatible with standard MCP config format for `.mcp.json` or Claude Desktop

## Scenarios

```gherkin
Scenario: Agent fetches ticket
  Given dagger-mcp is running with a valid API key
  When the agent calls get_ticket({ ticket_id: "DGR-42" })
  Then the tool returns the ticket with full context as markdown

Scenario: Agent lists open tickets
  When the agent calls list_tickets({ status: "open" })
  Then the tool returns a list of open tickets with metadata
```

## Technical notes

- Small Go binary (~200-300 lines max), uses stdlib `net/http` to call the REST API
- Keep thin — all business logic lives in the REST API, not the MCP server
- JSON-RPC 2.0 structured request/response format
- See `.dag/docs/plan/dagger-plan/05-agent-integration.md` for tool definitions

## Related docs

- Agent integration: `.dag/docs/plan/dagger-plan/05-agent-integration.md`

---

*This ticket was created by opencode and reviewed by a human before publishing.*
