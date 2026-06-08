# Dagger — Product Plan

**Status:** Draft
**Date:** 2026-05-25

Dagger is the DAG-powered product management tool built for AI. A graph of connected documents — pitches → ADRs → code explorations → tickets — served as optimized context to code generation agents.

Built on the [DAG documentation conventions](https://github.com/mbaquerizo/docs-augmented-generation) and the evidence base from [dag-research](https://github.com/mbaquerizo/dag-research).

## Documents

| # | Document | Description |
|---|----------|-------------|
| 01 | [Why Dagger](01-why-dagger.md) | Problem statement, differentiation, target user |
| 02 | [Phases Overview](02-phases-overview.md) | Three-phase roadmap with deliverables |
| 03 | [Data Model](03-data-model.md) | PostgreSQL schema, multi-tenancy, auth |
| 04 | [Adapter Design](04-adapter-design.md) | DAG → Dagger HTTP adapter |
| 05 | [Agent Integration](05-agent-integration.md) | REST API + MCP server + context assembly |
| 06 | [Business Model](06-business-model.md) | BYOK first → bundled later, pricing |
| 07 | [Open Questions](07-open-questions.md) | Risks and unresolved decisions |

## Phases

| Phase | Months | Focus | Delivers |
|-------|--------|-------|----------|
| 1 | 1-2 | Agent API | DAG publishes to Dagger, agents fetch full ticket context |
| 2 | 3-4 | Web UI | Doc/ticket browser, teams, search, manual CRUD |
| 3 | 5-6 | Chat | DAG-powered chat with BYOK LLM integration |

## Stack

| Layer | Phase 1 | Phase 2+ |
|-------|---------|----------|
| API server | Go | Go |
| Database | PostgreSQL | PostgreSQL (+ pgvector Phase 3) |
| Frontend | Static landing page | TanStack Start |
| Agent integration | REST API + MCP | Same |
| DAG integration | HTTP ingest API | Same |
| Deployment | Fly.io or Docker Compose | Pulumi (Go or TS) |
| Auth | API keys | API keys + user accounts |

## Key Design Decisions

1. **BYOK first, bundled later** — Phase 3 chat launches with BYOK to prove workflow without marginal cost risk. Bundled as default when usage patterns are understood.
2. **API-first** — Phase 1 delivers agent value before human UI. Agents are the primary customer.
3. **DAG adapter follows Dagger API** — Build the ingest endpoint first, then write the "dagger" adapter type in the DAG plugin repo.
4. **Go long-running server** — Simpler than serverless for this workload. SSE streaming and goroutine-per-request model are natural fits.
