---
id: DGR-9
issueType: story
status: done
tags:
  - phase-1
  - api
  - ingest
parent: DGR-5
blockedBy:
  - DGR-8
---

# Implement POST /api/v1/publish endpoint

## Description

As a DAG plugin user, I want to publish documents and tickets to Dagger via the API so they're available for agent context retrieval. This is the primary ingest endpoint that the DAG "dagger" adapter calls.

## Acceptance criteria

1. `POST /api/v1/publish` accepts publish payload (type, title, body, parent, metadata)
2. Validates required fields, returns `422 Unprocessable Entity` for invalid payloads
3. Inserts new docs/issues into the database with auto-assigned display ID (`DGR-N`)
4. Supports parent references by internal ID (not external_id)
5. Supports `metadata.relationships` with `targetId` references to Dagger internal IDs
6. Supports `metadata.docType` (for docs: adr, pitch, ce) and `metadata.issueType` (for tickets: epic, story, task, bug, spike)
7. Returns response with `id`, `displayId`, and `url`
8. Referenced parent/relationship IDs that don't exist return `400`
9. Protected by API key auth

## Scenarios

```gherkin
Scenario: Publish a new ADR
  Given a valid API key
  When I POST an ADR payload to /api/v1/publish
  Then the response is 201 Created
  And the response includes id, displayId (e.g. "DGR-3"), and url

Scenario: Publish a child referencing a parent by ID
  Given a parent doc exists with id 127
  When I POST a child with parent: 127
  Then the child's parent_id is set to 127

Scenario: Publish with doc_issue relationship
  Given a ticket exists with id 128
  When I POST an ADR with metadata.relationships [{ targetId: 128, type: "motivates" }]
  Then a doc_issues row links the ADR to the ticket

Scenario: Referenced ID not found
  Given no doc exists with id 999
  When I POST a payload with parent: 999
  Then the response is 400 Bad Request

Scenario: Missing required fields
  Given a payload missing "type"
  When I POST to /api/v1/publish
  Then the response is 422 Unprocessable Entity
```

## Technical notes

- Body field `parent` maps to `parent_id` on docs or issues
- Body field `projectId` selects the project; default to first project if not specified
- Display ID allocated atomically in the INSERT transaction
- See `.dag/docs/plan/dagger-plan/04-adapter-design.md` for full contract

## Related docs

- Adapter design: `.dag/docs/plan/dagger-plan/04-adapter-design.md`
- Agent integration: `.dag/docs/plan/dagger-plan/05-agent-integration.md`
- Data model: `.dag/docs/plan/dagger-plan/03-data-model.md`

---

*This ticket was created by opencode and reviewed by a human before publishing.*
