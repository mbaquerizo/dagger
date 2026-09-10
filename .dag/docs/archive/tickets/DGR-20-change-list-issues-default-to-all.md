---
id: DGR-20
issueType: story
status: done
tags:
  - mcp
blockedBy: []
parent: DGR-5
---

# Change list_issues default to "all" when status omitted

## Description

As an AI agent exploring a project, I want `list_issues` to return all issues when no status is specified, so that I can get a complete picture of the project without needing to know what status values exist.

Currently, `list_issues` defaults to filtering by `"open"` when no status parameter is provided. This is surprising — an agent trying to inventory all work misses closed/done/completed issues.

## Acceptance criteria

1. `list_issues` with no `status` argument returns all issues regardless of status
2. `list_issues` with `status: "done"` returns only done issues
3. Existing explicit-filter behavior is unchanged
4. Tool description in `tools/list` clearly says "when status is omitted, returns all issues"
5. HTTP handler in `internal/issues/handler.go` also updated to not default to "open"

## Scenarios

```gherkin
Scenario 1: No filter returns all issues
  Given issues with various statuses exist (open, in-progress, done)
  When I call list_issues without a status parameter
  Then I receive all issues regardless of status

Scenario 2: Explicit filter still works
  Given issues with various statuses exist
  When I call list_issues with status "done"
  Then I receive only issues with status "done"

Scenario 3: HTTP endpoint also defaults to all
  Given the HTTP API at GET /api/v1/issues
  When I call it without a ?status= query parameter
  Then it returns all issues instead of defaulting to "open"
```

## Technical notes

- Remove `if status == "" { status = "open" }` from `internal/mcp/dbservice.go` (the MCP fallback)
- Remove same fallback from `internal/issues/handler.go` (the HTTP fallback)
- Modify `ListIssues` SQL in `internal/issues/issues.go` to use `WHERE ($1 = '' OR i.status = $1)` — this keeps a single query shape
- Passing `""` as a query parameter to Postgres with the `($1 = '' OR i.status = $1)` pattern means the `$1 = ''` clause matches everything

## Related docs

- ADR: (none)
- Code exploration: (none)
- Parent epic: DGR-5

---

*This ticket was created by opencode and reviewed by a human before publishing.*
