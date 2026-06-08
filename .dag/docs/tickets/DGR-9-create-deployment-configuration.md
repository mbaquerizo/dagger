---
id: DGR-9
issueType: task
status: open
tags:
  - phase-1
  - deployment
  - railway
  - supabase
parent: DGR-1
blockedBy:
  - DGR-2
  - DGR-3
---

# Create deployment configuration

## Description

Configure the project for deployment on Railway with Supabase as the database. This includes the Railway project setup, environment variable configuration, Supabase project setup, and local dev instructions.

## Acceptance criteria

1. Supabase project configured, connection string documented
2. Railway project created with Go API service
3. Environment variables documented in `.env.example`:
   - `DATABASE_URL` — Supabase Postgres connection string (pooled)
   - `DAGGER_DEV_API_KEY` — seed API key for development
4. README section with deployment instructions:
   - Create Supabase project → copy connection string
   - Create Railway project → connect GitHub repo → set env vars → deploy
5. Railway-compatible start command configured (Railway auto-detects Go)
6. Local dev instructions use Supabase connection string directly (no Docker)

## Technical notes

- Railway auto-detects Go apps; no Dockerfile needed
- Use Supabase connection pooling (PgBouncer) via port 6543

## Related docs

- ADR-001: `.dag/docs/adr/2026-06-06-deployment-architecture/2026-06-06-deploy-on-railway-plus-supabase.md`

---

*This ticket was created by opencode and reviewed by a human before publishing.*
