---
id: DGR-6
issueType: story
status: done
tags:
  - phase-1
  - api
  - issues
  - context-assembly
parent: DGR-1
blockedBy:
  - DGR-5
---

# Implement GET /api/v1/agent/issues/:displayId endpoint

## Description

As an AI agent, I want to fetch a ticket by its display ID and receive full context as optimized markdown so I can generate correct code. This is the core value prop of Dagger.

## Acceptance criteria

1. `GET /api/v1/agent/issues/:displayId` accepts display IDs like `DGR-42`
2. Returns optimized markdown with ticket header (title, status, type, parent)
3. Includes "Linked Context" section with all related docs (ADRs, CEs, pitches)
4. Includes parent epic info and child subtask titles
5. Omits internal metadata (workspace IDs, timestamps) for token efficiency
6. Returns `404` for unknown display IDs
7. `GET /api/v1/agent/issues?status=open` returns JSON list of issues matching the filter (no context assembly)
8. Protected by API key auth

## Scenarios

```gherkin
Scenario: Fetch existing ticket
  Given a ticket with displayId "DGR-42" exists, linked to an ADR and a CE
  When I GET /api/v1/agent/issues/DGR-42
  Then the response is 200
  And the body is markdown with the ticket header and linked ADR/CE

Scenario: Ticket not found
  Given no ticket has displayId "DGR-999"
  When I GET /api/v1/agent/issues/DGR-999
  Then the response is 404
```

## Technical notes

- Context assembly engine traverses: doc_issues → linked docs, parent_id → epic, children → subtasks
- Renders each linked doc with key sections, token-optimized
- See `.dag/docs/plan/dagger-plan/05-agent-integration.md` for output format

## Related docs

- Agent integration: `.dag/docs/plan/dagger-plan/05-agent-integration.md`
- Data model: `.dag/docs/plan/dagger-plan/03-data-model.md`

---

*This ticket was created by opencode and reviewed by a human before publishing.*
