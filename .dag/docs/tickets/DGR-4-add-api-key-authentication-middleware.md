---
id: DGR-4
issueType: task
status: done
tags:
  - phase-1
  - auth
  - api-keys
parent: DGR-1
---

# Add API key authentication middleware

## Description

Add API key authentication to the Go API server. Every request (except health check) must include an `Authorization: Bearer dgr_...` header. Keys are hashed with SHA-256 and stored in the `api_keys` table. This protects all endpoints from day one.

## Acceptance criteria

1. Auth middleware extracts `Bearer` token from `Authorization` header
2. Middleware looks up key hash in `api_keys` table and verifies via SHA-256 hash lookup
3. Requests with missing or invalid keys return `401 Unauthorized`
4. Valid requests proceed with workspace context attached to the request
5. Health check endpoint (`GET /healthz`) is excluded from auth
6. A seed API key is generated via a separate CLI tool: `make seedkey ARGS="--workspace-id=<id>"`
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

- Use `crypto/sha256` (stdlib) for key hashing — no `golang.org/x/crypto` dependency needed
- Keys are generated as random 32-byte values, base64-encoded, prefixed `dgr_`
- Store only `key_hash` (SHA-256 hex), `prefix` (first 8 chars), never the raw key
- First 8 chars stored as `prefix` for indexed lookup; full hash comparison done in Go
- No server-side pepper, no salt — SHA-256 is deterministic; a unique raw key is the only secret
- Server validates keys only, never creates them — no startup side effects
- Key generation is via a separate CLI tool at `cmd/seedkey`
- Middleware attaches `workspace_id`, `project_id`, and `key_id` to `context.Context`
- Middleware accepts a `keyQuerier` interface, testable without a real database via pgxmock
- Health check endpoint (`GET /healthz`) bypasses auth

## Related docs

- Data model: `.dag/docs/plan/dagger-plan/03-data-model.md` (api_keys table)
- ADR-001: `.dag/docs/adr/2026-06-06-deployment-architecture/2026-06-06-deploy-on-railway-plus-supabase.md`

---

*This ticket was created by opencode and reviewed by a human before publishing.*
