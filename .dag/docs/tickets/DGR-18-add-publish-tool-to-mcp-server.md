---
id: DGR-18
issueType: story
status: done
tags:
  - phase-1
  - mcp
parent: DGR-5
blockedBy:
  - DGR-9
  - DGR-12
---

# Add publish tool to MCP server

## Description

As an AI agent using the MCP interface, I want to publish documents and tickets
to Dagger so I can create new artifacts during a session without switching to
the REST API. This extends the existing MCP server with a `publish` tool that
wraps the existing `POST /api/v1/publish` service.

## Acceptance criteria

1. `publish` tool registered in `tools/list` with all 5 tools now returned
2. Accepts parameters: `type` (required), `title` (required), `body` (required),
   `project_id` (required), `issue_type` (required if type=issue),
   `status` (optional), `parent_id` (optional)
3. Calls `internal/publish.Publish` directly via the `DBService` (same pattern
   as existing MCP tools)
4. Returns the published entity's `id`, `displayId`, and `url` on success
5. Returns descriptive error messages for validation failures (missing required
   fields, invalid type, unknown project, etc.)
6. Scoped to the authenticated workspace and project via existing auth context

## Scenarios

```gherkin
Scenario: Publish an ADR
  Given a valid MCP session with project_id=1
  When I call publish with type="adr", title="My ADR", body="# Title", project_id=1
  Then the response includes id, displayId (e.g. "DGR-47"), and url

Scenario: Publish an issue
  Given a valid MCP session with project_id=1
  When I call publish with type="issue", issue_type="story", title="My Story", body="desc", project_id=1
  Then the response includes id, displayId, and url

Scenario: Missing required type parameter
  When I call publish without type
  Then the response error code is -32602 (invalid params)

Scenario: Validation error from publish service
  Given an invalid issue_type value
  When I call publish with type="issue", issue_type="invalid"
  Then the response error message describes the validation failure
```

## Technical notes

- `project_id` is required (not auto-detected) — the caller must specify it
- `parent_id` refers to the internal DB id, not the display_id
- The `DBService.poolIface` needs a `Begin` method added so the pool can be
  passed to `internal/publish.Publish`
- `baseURL` for constructing the response URL needs to be threaded through the
  `DBService` (via `NewDBService` or config) — similar to how the REST handler
  receives it

## Files

| File | Change |
|------|--------|
| `internal/mcp/mcp.go` | Add `Publish(...)` to `ToolService` interface; add `publish` tool definition in `ListTools()` |
| `internal/mcp/server.go` | Add `publish` case to `handleToolCall` switch |
| `internal/mcp/dbservice.go` | Add `Begin` to `poolIface`; implement `Publish` method calling `internal/publish.Publish` |
| `internal/mcp/mcp_test.go` | Update tool count assertion (4→5) |
| `internal/mcp/server_test.go` | Add test for publish tool dispatch |
| `internal/mcp/dbservice_test.go` | Add test for DBService.Publish |

## Related docs

- Publish endpoint: `DGR-9` (`.dag/docs/tickets/DGR-9-implement-post-api-v1-publish-endpoint.md`)
- MCP server: `DGR-12` (`.dag/docs/tickets/DGR-12-build-mcp-server.md`)
- Agent integration: `.dag/docs/plan/dagger-plan/05-agent-integration.md`

---

*This ticket was created by opencode and reviewed by a human before publishing.*
