---
id: DGR-7
issueType: story
status: open
tags:
  - phase-1
  - api
  - documents
parent: DGR-1
blockedBy:
  - DGR-2
  - DGR-3
  - DGR-4
---

# Implement GET /api/v1/documents/:id endpoint

## Description

As an AI agent, I want to fetch an individual document by its ID or display ID so I can read its full content without the context assembly that tickets provide.

## Acceptance criteria

1. `GET /api/v1/documents/:id` accepts internal IDs (integer) or display IDs (`DGR-5`)
2. Returns the document as JSON with: id, displayId, type, title, body, status, createdAt, updatedAt
3. Returns `404` for unknown IDs
4. Protected by API key auth

## Scenarios

```gherkin
Scenario: Fetch by internal ID
  Given a doc exists with id 128
  When I GET /api/v1/documents/128
  Then the response is 200
  And the body includes id: 128, displayId, title, and body

Scenario: Fetch by display ID
  Given a doc exists with displayId "DGR-3"
  When I GET /api/v1/documents/DGR-3
  Then the response is 200

Scenario: Document not found
  When I GET /api/v1/documents/999
  Then the response is 404
```

## Related docs

- Agent integration: `.dag/docs/plan/dagger-plan/05-agent-integration.md`

---

*This ticket was created by opencode and reviewed by a human before publishing.*
