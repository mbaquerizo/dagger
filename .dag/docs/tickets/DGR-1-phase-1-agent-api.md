---
id: DGR-1
issueType: epic
status: open
tags:
  - phase-1
  - api
---

# Phase 1 — Agent API

## Description

Build the Dagger API server: an HTTP service with PostgreSQL backend that allows the DAG plugin to publish documents and tickets, and AI agents to retrieve full ticket context as optimized markdown. This is the foundation of the entire Dagger platform — API-first, agent-first.

## Scope

- **In scope:**
  - Go API server with chi router and pgx PostgreSQL driver
  - Database schema for docs, issues, relationships, API keys, projects, and multi-tenancy
  - `POST /api/v1/publish` — ingest docs and tickets from DAG plugin
  - `GET /api/v1/tickets/:displayId` — return ticket with full linked context
  - `GET /api/v1/documents/:id` — return individual document
  - API key authentication middleware
  - MCP server wrapping the REST API for Claude Code / OpenCode
  - Railway + Supabase deployment configuration
  - Static landing page with product explanation and API reference
  - DAG "dagger" adapter in the DAG plugin repo
- **Out of scope:**
  - Web UI (Phase 2)
  - User accounts and session management (Phase 2)
  - Full-text search (Phase 2)
  - Chat interface (Phase 3)
  - pgvector / semantic search (Phase 3)

## Related docs

- ADR-001: `.dag/docs/adr/2026-06-06-deployment-architecture/2026-06-06-deploy-on-railway-plus-supabase.md`
- Product plan: `.dag/docs/plan/dagger-plan/index.md`
- Data model: `.dag/docs/plan/dagger-plan/03-data-model.md`
- Adapter design: `.dag/docs/plan/dagger-plan/04-adapter-design.md`
- Agent integration: `.dag/docs/plan/dagger-plan/05-agent-integration.md`
- Go learning plan: `.dag/docs/plan/learning-go/README.md`

## Stories

- DGR-2: Scaffold Go project skeleton
- DGR-3: Set up PostgreSQL schema and migration infrastructure
- DGR-4: Add API key authentication middleware
- DGR-5: Implement POST /api/v1/publish endpoint
- DGR-6: Implement GET /api/v1/tickets/:displayId endpoint
- DGR-7: Implement GET /api/v1/documents/:id endpoint
- DGR-8: Build MCP server
- DGR-9: Create deployment configuration
- DGR-10: Build landing page
- DGR-11: DAG "dagger" adapter

---

*This ticket was created by opencode and reviewed by a human before publishing.*
