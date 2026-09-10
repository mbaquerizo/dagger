---
status: accepted
date: 2026-06-06
docType: code-exploration
tags:
  - infrastructure
  - deployment
  - research
---

# Code Exploration: Deployment platform research

## Goal

Evaluate deployment options for Dagger's Phase 1 Go API server with an eye toward Phase 2-3 requirements (web UI, auth, object storage, vector search). Key constraints: cost sensitivity, avoid Docker for local dev, simple operations while learning Go.

## Platforms evaluated

### App hosting

| Platform | Min cost | Go support | Managed Postgres | Free tier |
|----------|----------|------------|------------------|-----------|
| Railway | $5/mo Hobby | Auto-detect (Nixpacks) | Containerized (no PITR) | No (30d trial) |
| Fly.io | ~$7-10/mo | Dockerfile only | Community (unmanaged) | No (2h trial) |
| Render | $7/mo Starter | Buildpacks | Full managed (PITR) | Yes (spins down 15min) |
| Hetzner VPS | ~$3.50/mo | Manual | Manual | N/A |
| DigitalOcean App Platform | $5/mo Starter | Buildpacks + Docker | Managed ($15/mo) | Static sites only |

### Data platforms

| Platform | Free DB | Object storage | Auth | pgvector | Free tier caveat |
|----------|---------|---------------|------|----------|-----------------|
| Supabase | 500MB Postgres | 1GB S3 + CDN | 50K MAU (OAuth, email) | Included | Pauses after 7d inactivity |
| Neon | 0.5GB Postgres | None | JS-only | Included | 100 compute-hr/mo limit |
| Railway | None (separate) | None | None | Manual | — |

## Key findings

1. **Railway has the best DX for Go** — auto-detects Go apps, git push deploy, no Dockerfile required. Also handles Node.js for Phase 2 web UI.
2. **Supabase has the most feature density** — Postgres + Storage + Auth + pgvector in one platform. The only option that covers all three non-compute needs.
3. **Cost sweet spot** — Railway Hobby ($5/mo) + Supabase free covers all Phase 1 needs. Supabase free scales to Pro ($25/mo) when limits are hit.
4. **Local dev** — Both Railway and Supabase provide connection strings. Go app with pgx connects directly, no Docker needed.

## Recommendation

Railway for compute (Go API + future web UI). Supabase for data platform (Postgres + Storage + Auth + pgvector). Two platforms, scales through all three phases without migration.
