---
id: DGR-8
issueType: task
status: done
blockedBy:
  - DGR-7
tags:
  - phase-1
  - database
  - postgresql
  - supabase
parent: DGR-5
---

# Set up PostgreSQL schema and migration infrastructure

## Description

Create the database schema for all Phase 1 tables and the migration runner that applies them on application startup. Uses raw SQL with `pgx` — no ORM. Database is hosted on Supabase (free tier).

## Acceptance criteria

1. SQL migration files under `db/migrations/` include all Phase 1 tables:
   - `workspaces`, `users`, `workspace_members`
   - `projects` (with `next_display_number` counter, `slug` like `DGR`)
   - `docs` (with `display_id`, `project_id`, `workspace_id`, no `external_id`)
   - `issues` (with `display_id`, `project_id`, `workspace_id`)
   - `doc_issues`, `issue_relations`, `relations`
   - `api_keys` (with optional `project_id` for scoped keys)
2. Display ID sequence: atomically increment `projects.next_display_number` on INSERT, format as `{slug}-{number}`
3. Migration runner applies pending migrations on startup in order
4. `pgx` connection pool configured from `DATABASE_URL` env var (from Supabase)
5. Migrations are idempotent (safe to re-run)
6. Local dev uses Supabase free tier Postgres — no Docker, no local Postgres install

## Scenarios

```gherkin
Scenario: Fresh database
  Given no tables exist in the Supabase project
  When the server starts
  Then all migration files are applied in order
  And all tables exist with correct columns

Scenario: Display ID increments across types
  Given a project with slug "DGR" exists
  When a doc is inserted
  Then its display_id is "DGR-1"
  When an issue is inserted next
  Then its display_id is "DGR-2"

Scenario: Already migrated database
  Given tables already exist
  When the server starts
  Then no migrations are re-applied
```

## Technical notes

- Use `pgx/v5` for PostgreSQL driver and connection pooling
- Migration files are plain `.sql` files, applied in filename order
- Track applied migrations in a `schema_migrations` table
- Display ID allocated atomically: `UPDATE projects SET next_display_number = next_display_number + 1 WHERE id = $1 RETURNING next_display_number - 1`
- See `.dag/docs/plan/dagger-plan/03-data-model.md` for full schema

## Related docs

- Data model: `.dag/docs/plan/dagger-plan/03-data-model.md`
- ADR-001: `.dag/docs/adr/2026-06-06-deployment-architecture/2026-06-06-deploy-on-railway-plus-supabase.md`

---

*This ticket was created by opencode and reviewed by a human before publishing.*
