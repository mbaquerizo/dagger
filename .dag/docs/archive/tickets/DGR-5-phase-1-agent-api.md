---
id: DGR-5
issueType: epic
status: open
tags:
  - phase-1
  - api
children:
  - DGR-7
  - DGR-8
  - DGR-6
  - DGR-9
  - DGR-10
  - DGR-11
  - DGR-12
  - DGR-13
  - DGR-14
  - DGR-15
  - DGR-16
  - DGR-17
  - DGR-18
  - DGR-19
  - DGR-20
  - DGR-21
  - DGR-22
  - DGR-23
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

- DGR-7: Scaffold Go project skeleton
- DGR-8: Set up PostgreSQL schema and migration infrastructure
- DGR-6: Add API key authentication middleware
- DGR-9: Implement POST /api/v1/publish endpoint
- DGR-10: Implement GET /api/v1/tickets/:displayId endpoint
- DGR-11: Implement GET /api/v1/documents/:id endpoint
- DGR-12: Build MCP server
- DGR-13: Create deployment configuration
- DGR-14: Build landing page
- DGR-15: DAG "dagger" adapter
- DGR-16: Add PATCH endpoint to update issue status
- DGR-17: Enforce API token project_id scoping
- DGR-18: Add publish tool to MCP server
- DGR-19: Show internal DB ID in rendered markdown outputs
- DGR-20: Change list_issues default to "all" when status omitted
- DGR-21: Support issue_relations in publish tool
- DGR-22: Add add_issue_relation tool + API endpoint
- DGR-23: Polish MCP tool descriptions

---

*This ticket was created by opencode and reviewed by a human before publishing.*
