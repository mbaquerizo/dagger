# DAG → Dagger Adapter Design

## The Data Flow

```
DAG Plugin (CLI)
    │
    │  dag:plan-feature → ADR + CE
    │  dag:create-tickets → ticket hierarchy
    │
    ▼
.dag/config.json → { type: "dagger", ... }
    │
    ▼
POST /api/v1/publish  ◄── Dagger API
    │
    ▼
PostgreSQL (docs, issues, doc_issues)
    │
    ├── GET /api/v1/tickets/:id  → Agent
    └── Web UI (Phase 2)        → Human
```

## Adapter Configuration

In the DAG plugin's `.dag/config.json`:

```json
{
  "adapters": {
    "documentation": {
      "type": "dagger",
      "config": {
        "url": "https://api.dagger.dev",
        "workspaceId": "ws_..."
      }
    },
    "issue-tracking": {
      "type": "dagger",
      "config": {
        "url": "https://api.dagger.dev",
        "workspaceId": "ws_..."
      }
    }
  }
}
```

Authentication: `Authorization: Bearer dgr_...` header passed on every request.

## Publish Contract

Same contract defined in the DAG plugin's `contract.md`:

| Field | Required | Description |
|-------|----------|-------------|
| `type` | yes | `adr`, `pitch`, `ce`, `ticket` |
| `title` | yes | Display heading |
| `body` | yes | Markdown content |
| `parent` | no | Grouping identifier (maps to `group_id` in schema) |
| `metadata` | no | Structured fields: `status`, `tags`, `docType`, `issueType`, etc. |

## API Endpoint

**`POST /api/v1/publish`**

Request body:
```json
{
  "type": "adr",
  "title": "Use PostgreSQL for document storage",
  "body": "# ADR-004: Use PostgreSQL...\n\n## Context...",
  "metadata": {
    "status": "approved",
    "tags": ["database", "architecture"],
    "relationships": [
      {
        "targetId": 127,
        "type": "motivates"
      }
    ]
  }
}
```

Relationships reference Dagger internal IDs (returned from previous publish calls). The DAG plugin publishes in dependency order: parent first, then children referencing the returned IDs.

Response:
```json
{
  "id": 128,
  "displayId": "DGR-5",
  "url": "https://api.dagger.dev/tickets/DGR-5"
}
```

The response includes the Dagger-assigned internal ID, the human-readable display ID, and a URL. The DAG plugin captures these for referencing in subsequent publishes.

## Implementation Order

1. **Build Dagger API first** — implement `POST /api/v1/publish` endpoint
2. **Write the "dagger" adapter** — in the DAG plugin repo, implement the adapter that calls this endpoint
3. **Iterate** — as the API evolves, update the adapter

The adapter itself is thin — it's an HTTP client that maps the DAG publish contract to the Dagger API payload. Most complexity is on the Dagger side (validation, deduplication, relationship resolution).

## Relationship Resolution

The DAG plugin publishes docs and tickets in dependency order. Dagger tracks relationships via:

- **`parent_id`** — set in the body, references a Dagger internal ID from a previous publish response.
- **`metadata.relationships`** — array of `{ targetId, type }` entries. `targetId` is a Dagger internal ID from a previous response. Maps to `doc_issues` join table.
- **`group_id`** — if the DAG plugin publishes an ADR and CE from the same planning session, they share a `group_id` UUID for sibling grouping.

The caller is responsible for ordering publishes correctly (parents before children). Dagger validates that referenced IDs exist and returns 400 if not.

## Error Handling

| Error | HTTP | Adapter behavior |
|-------|------|------------------|
| Invalid payload | 422 | Report error, halt publish |
| Referenced ID not found | 400 | Report error, halt |
| Auth failure | 401 | Report error, halt |
| Server error | 500 | Retry with backoff |
