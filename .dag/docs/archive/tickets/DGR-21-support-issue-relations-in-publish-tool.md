---
id: DGR-21
issueType: story
status: done
tags:
  - mcp
  - publish
blockedBy:
  - DGR-19
parent: DGR-5
---

# Support issue_relations in publish tool

## Description

As a user creating an issue via the `publish` tool, I want to specify issue-to-issue relationships (e.g. "blocks", "duplicates") alongside the cross-type doc↔issue relationships, so that I can set up the full relationship graph in a single publish call.

Currently, `metadata.relationships` only creates rows in `doc_issues` (doc↔issue cross-type links). The `issue_relations` table exists in the schema but is never written to by any code.

## Acceptance criteria

1. Publishing an issue with `metadata.issue_relations: [{target_id: 5, relation_type: "blocks"}]` creates rows in `issue_relations`
2. Both forward and inverse rows are inserted (e.g. "blocks" + "blocked_by")
3. Unknown `relation_type` returns a validation error
4. Existing `metadata.relationships` (doc_issues) still works unchanged
5. In a single publish call, both doc_issues and issue_relations are created in the same transaction

## Scenarios

```gherkin
Scenario 1: One issue_relations entry creates forward and inverse rows
  Given I am publishing a new issue with internal ID 6
  And I specify metadata.issue_relations = [{target_id: 5, relation_type: "blocked_by"}]
  When the publish succeeds
  Then issue_relations has row (source_id=6, target_id=5, relation="blocked_by")
  And issue_relations has row (source_id=5, target_id=6, relation="blocks")

Scenario 2: relates_to is its own inverse
  Given I am publishing a new issue
  And I specify metadata.issue_relations = [{target_id: 5, relation_type: "relates_to"}]
  When the publish succeeds
  Then two rows are inserted, both with relation "relates_to"

Scenario 3: Invalid relation type returns error
  Given I am publishing a new issue
  When I specify metadata.issue_relations = [{target_id: 5, relation_type: "invalid"}]
  Then the publish fails with a validation error

Scenario 4: Coexists with existing relationships
  Given I am publishing a new issue
  When I specify both metadata.relationships and metadata.issue_relations
  Then both doc_issues and issue_relations rows are created in the same transaction
```

## Technical notes

- Add `IssueRelation` type and `IssueRelations []IssueRelation` field to `publish.Metadata` in `internal/publish/models.go`
- Validate `RelationType` against known values in `internal/publish/validate.go`
- Add inverse mapping in `internal/publish/publish.go`:
  ```go
  var relationInverse = map[string]string{
      "blocks":          "blocked_by",
      "blocked_by":      "blocks",
      "duplicates":      "duplicated_from",
      "duplicated_from": "duplicates",
      "relates_to":      "relates_to",
      "causes":          "caused_by",
      "caused_by":       "causes",
  }
  ```
- Insert loop resolves `relationType` → `relation_id`, inserts forward row, resolves inverse → `inverse_relation_id`, inserts inverse row
- Cross-type `doc_issues` field name is `relationships`, issue↔issue field name is `issue_relations` (distinct names for clarity)

## Related docs

- ADR: (none)
- Code exploration: (none)
- Parent epic: DGR-5

---

*This ticket was created by opencode and reviewed by a human before publishing.*
