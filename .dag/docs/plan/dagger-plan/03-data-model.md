# Data Model

Based on the DB adapter design doc (`docs-augmented-generation/.dag/docs/design/db-adapter-schema.md`) with modifications for Dagger SaaS.

## Core Tables

### `projects`

Namespaces for display IDs. Docs and issues share a single auto-incrementing sequence per project (e.g. `DGR-1`, `DGR-2`).

| Column | Type | Notes |
|--------|------|-------|
| `id` | `INTEGER PRIMARY KEY` | Auto-increment |
| `workspace_id` | `INTEGER NOT NULL REFERENCES workspaces(id)` | Multi-tenant isolation |
| `name` | `TEXT NOT NULL` | Display name |
| `slug` | `TEXT NOT NULL UNIQUE` | URL-friendly prefix for display IDs (e.g. `DGR`) |
| `next_display_number` | `INTEGER NOT NULL DEFAULT 1` | Auto-incrementing counter shared by docs and issues |
| `created_at` | `TEXT NOT NULL` | ISO 8601 |

### `docs`

| Column | Type | Notes |
|--------|------|-------|
| `id` | `INTEGER PRIMARY KEY` | Auto-increment |
| `display_id` | `TEXT NOT NULL` | Human-readable ID like `DGR-1` (shared sequence with issues) |
| `type` | `TEXT NOT NULL` | Free string: `adr`, `ce` (code exploration), `pitch`, `design-doc` |
| `title` | `TEXT NOT NULL` | Display heading |
| `body` | `TEXT` | Markdown content. Null only on partial writes. |
| `status` | `TEXT NOT NULL DEFAULT 'proposed'` | `proposed`, `approved`, `superseded`, `rejected` |
| `parent_id` | `INTEGER REFERENCES docs(id)` | Self-referencing FK for trees |
| `group_id` | `TEXT` | Optional UUID for sibling grouping (ADR + CE from same planning session) |
| `project_id` | `INTEGER NOT NULL REFERENCES projects(id)` | Scopes display ID namespace |
| `workspace_id` | `INTEGER NOT NULL REFERENCES workspaces(id)` | Multi-tenant isolation (denormalized for query perf) |
| `created_at` | `TEXT NOT NULL` | ISO 8601 |
| `updated_at` | `TEXT NOT NULL` | ISO 8601 |

### `issues`

| Column | Type | Notes |
|--------|------|-------|
| `id` | `INTEGER PRIMARY KEY` | Auto-increment |
| `display_id` | `TEXT NOT NULL` | Human-readable ID like `DGR-2` (shared sequence with docs) |
| `type` | `TEXT NOT NULL` | `epic`, `story`, `task`, `bug`, `spike` |
| `title` | `TEXT NOT NULL` | Display heading |
| `body` | `TEXT` | Markdown content |
| `status` | `TEXT NOT NULL DEFAULT 'open'` | Workflow states |
| `parent_id` | `INTEGER REFERENCES issues(id)` | Epic contains stories |
| `project_id` | `INTEGER NOT NULL REFERENCES projects(id)` | Scopes display ID namespace |
| `workspace_id` | `INTEGER NOT NULL REFERENCES workspaces(id)` | Multi-tenant isolation (denormalized for query perf) |
| `created_at` | `TEXT NOT NULL` | ISO 8601 |
| `updated_at` | `TEXT NOT NULL` | ISO 8601 |

### `doc_issues`

Doc-to-issue cross-references. The core relationship that powers context assembly.

| Column | Type | Notes |
|--------|------|-------|
| `id` | `INTEGER PRIMARY KEY` | |
| `doc_id` | `INTEGER NOT NULL REFERENCES docs(id)` | |
| `issue_id` | `INTEGER NOT NULL REFERENCES issues(id)` | |
| `relationship_type` | `TEXT NOT NULL` | `motivates`, `implements`, `references` |

### `issue_relations`

Typed directional edges between issues.

| Column | Type | Notes |
|--------|------|-------|
| `id` | `INTEGER PRIMARY KEY` | |
| `source_issue_id` | `INTEGER NOT NULL REFERENCES issues(id)` | |
| `relation_id` | `INTEGER NOT NULL REFERENCES relations(id)` | |
| `target_issue_id` | `INTEGER NOT NULL REFERENCES issues(id)` | |

### `relations`

Relation type definitions.

| Column | Type | Notes |
|--------|------|-------|
| `id` | `INTEGER PRIMARY KEY` | |
| `name` | `TEXT NOT NULL UNIQUE` | `blocks`, `is_blocked_by`, `duplicates`, `relates_to`, `causes` |

## Multi-Tenancy

### `workspaces`

| Column | Type | Notes |
|--------|------|-------|
| `id` | `INTEGER PRIMARY KEY` | |
| `name` | `TEXT NOT NULL` | |
| `slug` | `TEXT NOT NULL UNIQUE` | URL-friendly |
| `created_at` | `TEXT NOT NULL` | |

### `users`

| Column | Type | Notes |
|--------|------|-------|
| `id` | `INTEGER PRIMARY KEY` | |
| `email` | `TEXT NOT NULL UNIQUE` | |
| `name` | `TEXT` | |
| `created_at` | `TEXT NOT NULL` | |

### `workspace_members`

| Column | Type | Notes |
|--------|------|-------|
| `workspace_id` | `INTEGER REFERENCES workspaces(id)` | |
| `user_id` | `INTEGER REFERENCES users(id)` | |
| `role` | `TEXT NOT NULL DEFAULT 'member'` | `admin`, `member` |

## Auth

### `api_keys`

| Column | Type | Notes |
|--------|------|-------|
| `id` | `INTEGER PRIMARY KEY` | |
| `key_hash` | `TEXT NOT NULL` | Hashed key (never store raw) |
| `prefix` | `TEXT NOT NULL` | First 8 chars for identification |
| `name` | `TEXT NOT NULL` | User-given label |
| `workspace_id` | `INTEGER NOT NULL REFERENCES workspaces(id)` | |
| `project_id` | `INTEGER REFERENCES projects(id)` | Optional project-scoped key (null = full workspace access) |
| `created_by` | `INTEGER REFERENCES users(id)` | |
| `last_used_at` | `TEXT` | |
| `created_at` | `TEXT NOT NULL` | |
| `expires_at` | `TEXT` | Optional expiration |

## Key Design Decisions

1. **Shared display ID sequence** — Docs and issues share a single auto-incrementing counter per project (`next_display_number`). A published ADR gets `DGR-1`, the next ticket gets `DGR-2`, the next code exploration gets `DGR-3`. This gives every content item a unique, human-readable ID regardless of type.
2. **`display_id` on both docs and issues** — Both tables have a user-facing ID (`DGR-1`). Generated atomically on insert via `UPDATE projects SET next_display_number = next_display_number + 1 ... RETURNING`.
3. **`project_id` on docs and issues** — Scopes the display ID namespace per project. Also enables project-level API key scoping.
4. **`workspace_id` denormalized** — Present on both docs and issues despite being reachable via `project_id → projects.workspace_id`. Avoids join on every tenant-scoped query.
5. **No `external_id`** — Removed. Dagger is the system of record. Relationships are resolved by Dagger internal IDs, not external IDs. Idempotent retries can be handled via an `Idempotency-Key` header later if needed.
6. **`type` as free text** — No enums. New doc or issue types can be introduced without schema migrations.
7. **`relationship_type` as free text** — Same reasoning. Predefined values (`motivates`, `implements`, `references`) but not constrained.
8. **Parent-child via self-referencing FK** — Simple tree model. No join table needed for hierarchy. Many-to-many relationships use `doc_issues` and `issue_relations`.

## Migration from DB Adapter Schema

The original DB adapter design used SQLite. Dagger uses PostgreSQL. The schema is compatible with minor adjustments:

- `datetime('now')` defaults → `NOW()` or `CURRENT_TIMESTAMP`
- `INTEGER PRIMARY KEY` auto-increment works in both (SERIAL in PG, autoincrement in SQLite)
- Add `workspace_id` to all content tables for multi-tenancy
- Add `projects` table (new — not in original design)
- Add `api_keys` table (not in original design)
