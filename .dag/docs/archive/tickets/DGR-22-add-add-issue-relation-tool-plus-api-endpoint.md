---
id: DGR-22
issueType: story
status: done
tags:
  - mcp
  - api
blockedBy:
  - DGR-19
parent: DGR-5
---

# Add add_issue_relation tool + API endpoint

## Description

As an AI agent, I want to relate two existing issues without re-publishing one of them, so that I can update the relationship graph incrementally as I learn about connections between issues.

This adds a new MCP tool `add_issue_relation` and a corresponding HTTP API endpoint `POST /api/v1/issues/relations`, following the same three-layer pattern used by other tools (core function in `internal/issues`, HTTP handler, MCP dbservice).

## Acceptance criteria

1. `add_issue_relation(source_id, target_id, relation_type)` creates bidirectional rows in `issue_relations`
2. Corresponding HTTP endpoint `POST /api/v1/issues/relations` accepts `{source_id, target_id, relation_type}` JSON body
3. Unknown relation_type returns an error
4. Non-existent source or target issue ID returns an error
5. `source_id == target_id` returns an error
6. Tool is listed in `tools/list` with clear parameter descriptions
7. Relation types: blocks, blocked_by, duplicates, duplicated_from, relates_to, causes, caused_by

## Scenarios

```gherkin
Scenario 1: Create bidirectional relation
  Given issues with internal IDs 5 and 6 exist
  When I call add_issue_relation with source_id=5, target_id=6, relation_type="blocks"
  Then issue_relations has row (source_id=5, target_id=6, relation="blocks")
  And issue_relations has row (source_id=6, target_id=5, relation="blocked_by")

Scenario 2: Unknown relation type
  When I call add_issue_relation with relation_type="unknown"
  Then the call fails with an error

Scenario 3: Source equals target
  When I call add_issue_relation with source_id=5, target_id=5
  Then the call fails with an error

Scenario 4: Non-existent issue
  When I call add_issue_relation with source_id=9999
  Then the call fails with a not-found error

Scenario 5: HTTP endpoint works
  Given valid parameters
  When I POST to /api/v1/issues/relations with JSON body
  Then I receive a 200 OK response
```

## Technical notes

### Core function (`internal/issues/issues.go`)
```go
func AddIssueRelation(ctx context.Context, pool poolIface, sourceID, targetID int, relationType string, workspaceID int, authProjectID *int) error
```
- Resolve `relationType` → `relation_id` via `SELECT id FROM relations WHERE name = $1`
- Resolve inverse type using `relationInverse` map → `inverse_relation_id`
- Validate `sourceID != targetID`
- Validate both issues exist in workspace/auth scope
- Begin transaction, insert forward row, insert inverse row, commit

### HTTP handler (`internal/issues/handler.go`)
- `POST /api/v1/issues/relations`
- Standard auth/error pattern (same as `NewUpdateIssueStatusHandler`)
- JSON body: `{"source_id": 5, "target_id": 6, "relation_type": "blocks"}`

### MCP tool (`internal/mcp/mcp.go` + `server.go` + `dbservice.go`)
- Tool definition: `add_issue_relation(source_id, target_id, relation_type)`
- Interface method on `ToolService`
- Dispatch case in `server.go`
- Implementation in `dbservice.go` calls `issues.AddIssueRelation`

### Route (`cmd/api/main.go`)
```go
r.Post("/api/v1/issues/relations", issues.NewAddIssueRelationHandler(pool))
```

## Related docs

- ADR: (none)
- Code exploration: (none)
- Parent epic: DGR-5

---

*This ticket was created by opencode and reviewed by a human before publishing.*
