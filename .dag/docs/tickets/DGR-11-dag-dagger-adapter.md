---
id: DGR-11
issueType: story
status: open
tags:
  - phase-1
  - dag-plugin
  - adapter
parent: DGR-1
blockedBy:
  - DGR-5
---

# DAG "dagger" adapter

## Description

As a DAG plugin user, I want to publish docs and tickets to Dagger so they're available to AI agents via Dagger's API. This is a new adapter type in the DAG plugin repo — a separate project from Dagger itself.

## Acceptance criteria

1. New adapter type `"dagger"` added to the DAG plugin's adapter system
2. On publish, sends HTTP POST to `DAGGER_URL/api/v1/publish` with the publish payload
3. Reads `DAGGER_URL` and `DAGGER_API_KEY` from config or env
4. Sends `Authorization: Bearer <key>` header
5. Captures the response (displayId, URL) and reports it to the user
6. Publishes in dependency order: parent docs first, then children referencing returned IDs
7. Graceful error handling: reports failures, suggests retry

## Technical notes

- Implemented in the DAG plugin repo (TypeScript), not in this repo
- Maps the DAG publish contract to the Dagger API payload
- See `.dag/docs/plan/dagger-plan/04-adapter-design.md` for the contract

---

*This ticket was created by opencode and reviewed by a human before publishing.*
