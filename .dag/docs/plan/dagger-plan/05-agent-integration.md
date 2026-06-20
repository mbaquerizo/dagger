# Agent Integration

## Interface: REST API

The primary interface for agents. Returns optimized markdown with full context assembled.

### `GET /api/v1/agent/issues/:displayId`

Returns a single markdown document containing the issue and all linked context.

```
GET /api/v1/agent/issues/DG-42
Authorization: Bearer dgr_abc123...
```

Response body:
```markdown
# DG-42: Add rate limiting to the API

**Status:** in-progress  |  **Type:** story
**Parent:** DG-40 (API hardening epic)

## Acceptance Criteria (Gherkin)
Scenario: User exceeds rate limit
  Given a user has made 100 requests in the current window
  When they make the 101st request
  Then they receive a 429 response
  And the response includes a Retry-After header

Scenario: Rate limit resets
  Given a user is rate-limited
  When the rate limit window expires
  Then they can make requests again

---

## Linked Context

### ADR-004: Rate limiting strategy
**Status:** approved

**Decision:** Use token bucket algorithm with per-user buckets.
Redis-backed with 1-second precision.

**Options considered:**
- Token bucket (chosen) — burst-friendly, simple
- Sliding window log — memory-intensive at scale
- Fixed window — stampeding herd on reset

**Consequences:**
- Redis becomes a dependency for API servers
- Need graceful degradation if Redis is down

### Code Exploration: Rate limiting implementation
**Scope:** `internal/api/middleware/ratelimit/`

**Relevant files:**
- `internal/api/middleware/ratelimit/bucket.go` — Token bucket logic
- `internal/api/middleware/ratelimit/redis.go` — Redis backend
- `internal/api/middleware/ratelimit/middleware.go` — HTTP middleware

**Patterns:**
- Middleware checks rate limit before handler runs
- Headers: `X-RateLimit-Remaining`, `X-RateLimit-Reset`
- `429` response uses standard `application/problem+json` format

**Constraints:**
- Existing Redis cluster has <1ms p99 latency
- Max bucket size: 1000 tokens (configurable per endpoint)
```

The assembly engine resolves:
1. The issue itself (`issues` table)
2. Linked docs via `doc_issues` (ADR, code exploration, pitch)
3. Parent issue via `parent_id` (epic-level context)
4. Child issues (sub-tasks, for the "big picture")

### `GET /api/v1/agent/docs/:id`

Returns a single document as optimized markdown.

### `GET /api/v1/agent/issues?status=open`

List issues with optional filters. Returns JSON with minimal metadata (no context assembly).

## Interface: MCP Server

For Claude Code, OpenCode, and other MCP-native agents.

### Implementation

A small Go server (~200 lines) that runs as a subprocess and communicates via stdio JSON-RPC 2.0.

```
dagger-mcp --api-key dgr_abc123
```

### Tools

| Tool | Description | Parameters |
|------|-------------|------------|
| `get_issue` | Fetch issue with full context | `issue_id: string` (e.g. `DG-42`) |
| `search_docs` | Search across docs and issues | `query: string`, `limit?: number` |
| `list_issues` | List issues by status | `status?: string`, `limit?: number` |
| `get_doc` | Fetch individual document | `doc_id: string` |

### User Configuration

In the agent's MCP config:

```json
{
  "mcpServers": {
    "dagger": {
      "command": "dagger-mcp",
      "args": ["--api-key", "dgr_abc123"]
    }
  }
}
```

Then in conversation: *"Work on DG-42"* → agent calls `get_issue("DG-42")` → receives full context → generates code.

## Interface: Web UI (Phase 2)

For humans who want to read the same context the agent sees. The ticket detail view renders the assembled markdown directly, with linked docs in a sidebar.

## Context Assembly Engine

The core value prop. Given an issue display ID:

1. **Fetch** the issue from `issues` table
2. **Traverse** `doc_issues` to find linked docs (type-filter: ADR, CE, pitch)
3. **Traverse** `parent_id` to find parent epic (include its linked docs)
4. **Traverse** `children` to find sub-tasks (include titles only for scope)
5. **Assemble** into a single markdown document:
   - Ticket header (title, status, type, parent)
   - Acceptance criteria (from ticket body)
   - "Linked Context" section
   - Each linked doc rendered with its key sections
   - Token-optimized (omit internal metadata, workspace IDs, timestamps)

The output is designed for token efficiency. No chrome, no navigation, no repeated headers.
