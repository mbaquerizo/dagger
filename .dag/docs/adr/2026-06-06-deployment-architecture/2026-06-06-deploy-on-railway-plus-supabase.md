---
status: proposed
date: 2026-06-06
tags:
  - infrastructure
  - deployment
  - phase-1
---

# ADR-001: Deploy on Railway + Supabase

## Context

Dagger needs a home for Phase 1 (Go API server) and the infrastructure to support all three phases: database, object storage, authentication, and eventually a web UI and chat interface. The choices must balance cost, simplicity, and forward-compatibility.

Key constraints:
- **Cost sensitive**: prefer near-zero cost in Phase 1
- **Resource conscious**: avoid Docker for local development where possible
- **Go-first**: primary API is a Go binary with `pgx`
- **Future needs**: auth (Phase 2), object storage (Phase 2+), vector search (Phase 3)
- **Learning curve**: operator is learning Go — infrastructure should be simple

See the [code exploration](2026-06-06-CE-deployment-platform-research.md) for detailed research findings.

## Options Considered

### Railway (app) + Supabase (DB/storage/auth) — chosen

App hosted on Railway ($5/mo Hobby). Data platform on Supabase (free tier: Postgres + Storage + Auth + pgvector). Each side is best-in-class: Railway for compute deployment, Supabase for managed data services.

| Component | Platform | Phase 1 cost |
|---|---|---|
| Go API server | Railway Hobby | $5/mo |
| TanStack Start web UI (P2) | Railway (same project) | included |
| PostgreSQL | Supabase free (500MB) | $0 |
| Storage (screenshots) | Supabase free (1GB) | $0 |
| Auth (OAuth, sessions) | Supabase Auth (50K MAU free) | $0 |
| pgvector (P3) | Supabase (extension) | included |

Tradeoffs: $5/mo baseline, two platforms to manage, Supabase free tier pauses after 7 days inactivity.

### Railway (app) + Neon (DB) + separate storage

App and DB on separate platforms. No built-in object storage or auth. Three platforms. Storage requires a third service (Cloudflare R2 or DO Spaces). Cheapest if Neon's free tier is sufficient, but sprawl grows.

### DigitalOcean all-in-one

App Platform + Managed PostgreSQL + Spaces on one cloud. Flat pricing, single bill. App Platform $5-12/mo + Managed PG $15/mo + Spaces $5/mo = $25-32/mo. Predictable but 5-6x more expensive than chosen option in Phase 1.

### Self-hosted (Hetzner VPS + Docker Compose)

Maximum control, lowest raw cost (~$3.50/mo). Full stack: Postgres, Go API, reverse proxy, SSL, storage, all on one box. Requires managing Postgres, backups, monitoring, security updates. Too much ops overhead for a solo founder learning Go.

### Render (free) + Supabase (free)

Truly $0, but Render's free tier spins down after 15 minutes of inactivity. Cold start latency (~30s) is unacceptable for an agent-facing API that needs snappy responses. Mitigatable with a keepalive ping but adds operational complexity.

## Decision

**Use Railway for compute (Go API in Phase 1, TanStack Start web UI in Phase 2) and Supabase for the data platform (Postgres, Storage, Auth, pgvector).**

Rationale:

1. **Two platforms, not one** — Railway for app hosting, Supabase for everything else. Avoids sprawl of 3+ platforms while keeping each side specialized.

2. **Railway's DX** — auto-detects Go and Node.js, git push deploy, no Dockerfile needed. Fits the "learn Go, not DevOps" goal.

3. **Supabase's feature density** — Postgres + Storage + Auth + pgvector in one platform. Auth saves building GitHub OAuth from scratch in Phase 2. Storage handles ticket screenshots. pgvector future-proofs for Phase 3.

4. **$5/mo in Phase 1** — Supabase is free (2 projects, 500MB DB, 1GB storage, 50K MAU). Railway Hobby is $5/mo. Only move Supabase to Pro ($25/mo) when outgrowing free limits.

5. **No Docker for local dev** — Go app connects to Supabase Postgres directly via connection string. The same connection pattern works in dev, staging, and production.

6. **Scales through all three phases** — Phase 2 web UI is a second service in the same Railway project. Phase 3 SSE chat extends the Go API. No platform migration needed.

Caveat: Supabase free projects pause after 7 days of inactivity. For Phase 1 development this is manageable (periodic use keeps it alive, unpausing is one click). Upgrade to Pro when the API is deployed and needs to stay warm.

## Consequences

**Positive:**
- Only two platforms to manage across all three phases
- Supabase Auth removes an entire feature build in Phase 2
- Storage is ready before we need it (just write to bucket)
- pgvector is available from day one, no migration needed
- Local dev has zero resource overhead (no Docker, no local Postgres)

**Negative:**
- $5/mo baseline cost
- Supabase free tier pause behavior requires attention during inactive periods
- Two credentials to rotate (Railway + Supabase)
- API key auth in Phase 1 is custom Go code (Supabase Auth doesn't replace API keys)

**Neutral:**
- Partition between compute (Railway) and data (Supabase) is clean but means two dashboards
- Migrating from Supabase to another Postgres provider later is possible but would lose Auth and Storage integration
