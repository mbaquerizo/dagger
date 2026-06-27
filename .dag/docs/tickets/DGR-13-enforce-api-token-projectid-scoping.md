---
id: DGR-13
issueType: story
status: open
tags:
  - auth
  - security
  - phase-2
blockedBy: []
---

# Enforce API token projectID scoping

## Description

As an API consumer using a project-scoped API token, I want my access to be restricted to issues and docs within that project only, so that the token scoping is actually enforced.

Currently the auth middleware injects a `projectID` into the request context when the API key has a non-null `project_id`, but no handler or query uses it. A scoped token can read any issue/doc in the workspace and publish to any project.

## Acceptance criteria

1. `GET /api/v1/agent/issues/{displayId}` returns 404 when the issue belongs to a different project than the token's projectID
2. `GET /api/v1/issues` only returns issues in the token's project when the token has a projectID
3. `GET /api/v1/agent/docs/{displayId}` returns 404 when the doc belongs to a different project than the token's projectID
4. `POST /api/v1/publish` returns 403 when the token's projectID doesn't match the request's `projectId`
5. All sub-queries (parent issues, children, linked docs, related issues, publish parent/relationship checks) also filter by projectID for defense-in-depth
6. Auth checks (workspaceID, projectID) happen before request body validation in the publish handler

## Scenarios

```gherkin
Scenario: Project-scoped token reads issue in its project
  Given an API token scoped to project "P1"
  When the user requests GET /api/v1/agent/issues/DGR-42
  And issue DGR-42 belongs to project "P1"
  Then the response is 200 with the issue rendered as markdown

Scenario: Project-scoped token reads issue outside its project
  Given an API token scoped to project "P1"
  When the user requests GET /api/v1/agent/issues/DGR-99
  And issue DGR-99 belongs to project "P2"
  Then the response is 404 Not Found

Scenario: Unscoped token reads issue in any project
  Given an API token with no projectID restriction
  When the user requests GET /api/v1/agent/issues/DGR-42
  And issue DGR-42 belongs to project "P2"
  Then the response is 200 with the issue rendered as markdown

Scenario: Project-scoped token lists issues
  Given an API token scoped to project "P1"
  When the user requests GET /api/v1/issues
  Then only issues in project "P1" are returned

Scenario: Project-scoped token publishes to wrong project
  Given an API token scoped to project "P1"
  When the user POSTs to /api/v1/publish with projectId: 2
  Then the response is 403 Forbidden

Scenario: Project-scoped token reads doc outside its project
  Given an API token scoped to project "P1"
  When the user requests GET /api/v1/agent/docs/DGR-3
  And doc DGR-3 belongs to project "P2"
  Then the response is 404 Not Found
```

## Technical notes

Thread an optional `*int` projectID through all query functions. When non-nil, append `AND project_id = $N` to SQL WHERE clauses (incrementing the parameter index). This matches the existing `workspaceID` pattern but optional.

### Files to modify

| File | Change |
|------|--------|
| `internal/issues/issues.go` | Add `*int` param to `GetIssueContext`, `ListIssues`, `queryLinkedDocs`, `queryParent`, `queryChildren`, `queryRelatedIssues`; conditionally append SQL filter |
| `internal/issues/handler.go` | Extract projectID from auth context, pass to business logic |
| `internal/docs/docs.go` | Add `*int` param to `GetDoc`; conditionally append SQL filter |
| `internal/docs/handler.go` | Extract projectID from auth context, pass to business logic |
| `internal/publish/handler.go` | Add auth projectID vs request projectID validation (before body validation) |
| `internal/publish/publish.go` | Add `authProjectID *int` param; filter parent/relationship EXISTS queries by projectID |
| `internal/issues/handler_test.go` | Update expectations with new args, add scoped test cases |
| `internal/issues/issues_test.go` | Update calls with `nil`, add scoped test cases |
| `internal/docs/handler_test.go` | Update expectations with new args, add scoped test cases |
| `internal/docs/docs_test.go` | Update calls with `nil`, add scoped test cases |
| `internal/publish/handler_test.go` | Add projectID mismatch test |
| `internal/publish/publish_test.go` | Add projectID filter test |

### Handler auth flow (publish)

1. Decode JSON → 400
2. Check workspaceID → 401
3. Check authProjectID vs req.ProjectID mismatch → 403
4. Validate body → 422
5. Publish with authProjectID filter

## Related docs

- Auth middleware: `internal/auth/middleware.go`
- Auth context: `internal/auth/context.go`
- Issues handler: `internal/issues/handler.go`
- Docs handler: `internal/docs/handler.go`
- Publish handler: `internal/publish/handler.go`

---

*This ticket was created by opencode and reviewed by a human before publishing.*
