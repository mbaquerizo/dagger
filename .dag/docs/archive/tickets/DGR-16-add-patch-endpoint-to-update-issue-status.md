---
id: DGR-12
issueType: task
status: done
tags:
  - api
blockedBy: []
---

# Add PATCH endpoint to update issue status

## Description

Add an endpoint to update the status of an issue. Currently the API supports
listing issues (`GET /api/v1/issues`) and fetching a single issue with full
context (`GET /api/v1/agent/issues/{displayId}`), but there is no way to
transition an issue through its workflow states.

## Acceptance criteria

1. `PATCH /api/v1/issues/{displayId}/status` accepts `{ "status": "<new-status>" }`
2. Validates status is one of: `open`, `in-progress`, `in-review`, `done`, `closed`
3. Returns 404 if the issue does not exist
4. Returns 422 if the status value is invalid
5. Returns 200 on success
6. Protected by the existing API-key auth middleware

## Files

- `internal/issues/handler.go` — new handler function
- `internal/issues/issues.go` — `UpdateIssueStatus` query function
- `internal/issues/models.go` — request/response types
- `cmd/api/main.go` — route registration

---

*This ticket was created by opencode and reviewed by a human before publishing.*
