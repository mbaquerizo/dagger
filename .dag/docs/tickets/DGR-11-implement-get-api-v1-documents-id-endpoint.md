---
id: DGR-11
issueType: story
status: done
tags:
  - phase-1
  - api
  - docs
parent: DGR-5
blockedBy:
  - DGR-9
---

# Implement GET /api/v1/agent/docs/:displayID endpoint

## Description

As an AI agent, I want to fetch an individual document by its display ID so I can read its full content without the context assembly that tickets provide.

## Acceptance criteria

1. `GET /api/v1/agent/docs/:id` accepts internal IDs (integer) or display IDs (`DGR-5`)
2. Returns the document as optimized markdown with header (title, type, status) and full body
3. Returns `404` for unknown IDs
4. Protected by API key auth

## Scenarios

```gherkin
Scenario: Fetch by internal ID
  Given a doc exists with id 128
  When I GET /api/v1/agent/docs/128
  Then the response is 200
  And the body is markdown with the doc header and full content

Scenario: Fetch by display ID
  Given a doc exists with displayId "DGR-3"
  When I GET /api/v1/agent/docs/DGR-3
  Then the response is 200

Scenario: Document not found
  When I GET /api/v1/agent/docs/999
  Then the response is 404
```

## Related docs

- Agent integration: `.dag/docs/plan/dagger-plan/05-agent-integration.md`

---

*This ticket was created by opencode and reviewed by a human before publishing.*
