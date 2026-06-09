---
id: DGR-4
issueType: task
status: open
tags:
  - phase-1
  - auth
  - api-keys
parent: DGR-1
blockedBy:
  - DGR-3
---

# Add API key authentication middleware

## Description

Add API key authentication to the Go API server. Every request (except health check) must include an `Authorization: Bearer dgr_...` header. Keys are hashed with bcrypt and stored in the `api_keys` table in Supabase. This protects all endpoints from day one.

## Acceptance criteria

1. Auth middleware extracts `Bearer` token from `Authorization` header
2. Middleware looks up key hash in `api_keys` table and verifies with bcrypt
3. Requests with missing or invalid keys return `401 Unauthorized`
4. Valid requests proceed with workspace context attached to the request
5. Health check endpoint (`GET /healthz`) is excluded from auth
6. A seed API key is generated for development (set via env var or printed on first startup)
7. Key secrets follow the prefix format `dgr_` for identification

## Scenarios

```gherkin
Scenario: Request with valid API key
  Given a valid API key exists in the api_keys table
  When a request includes Authorization: Bearer dgr_xxx...
  Then the request proceeds to the handler
  And the handler receives the workspace context

Scenario: Request with no API key
  Given no Authorization header is present
  When a request is made
  Then the response is 401 Unauthorized

Scenario: Request with invalid API key
  Given an API key that does not exist in the database
  When a request includes Authorization: Bearer dgr_invalid...
  Then the response is 401 Unauthorized

Scenario: Health check bypasses auth
  Given no Authorization header
  When a GET /healthz request is made
  Then the response is 200 OK
```

## Technical notes

- Use `golang.org/x/crypto/bcrypt` for key hashing
- Keys are generated as random 32-byte values, base64-encoded, prefixed `dgr_`
- Store only `key_hash` (bcrypt), `prefix` (first 8 chars), never the raw key
- Middleware attaches `workspace_id` to `context.Context` for downstream handlers
- Seed key for dev can be set via `DAGGER_DEV_API_KEY` env var or auto-generated on startup

## Related docs

- Data model: `.dag/docs/plan/dagger-plan/03-data-model.md` (api_keys table)
- ADR-001: `.dag/docs/adr/2026-06-06-deployment-architecture/2026-06-06-deploy-on-railway-plus-supabase.md`

---

*This ticket was created by opencode and reviewed by a human before publishing.*
