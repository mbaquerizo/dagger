# Three-Phase Roadmap

## Phase 1 (Months 1-2): Agent API

**Goal:** DAG plugin can publish to Dagger, and agents can retrieve full ticket context via API.

### Deliverables

- **Go API server** — HTTP service with PostgreSQL backend
- **Database schema** — based on the DB adapter design doc (docs, issues, relationships, tags)
- **`POST /api/v1/publish`** — ingest docs and tickets from the DAG plugin
- **`GET /api/v1/agent/issues/:id`** — return issue + all linked context as a single optimized markdown document
- **`GET /api/v1/agent/docs/:id`** — return individual doc as markdown
- **API key authentication** — hashed keys stored in DB, scoped to projects
- **MCP server** — thin Go server wrapping the REST API for Claude Code / OpenCode integration
- **DAG "dagger" adapter** — new adapter type in the DAG plugin repo that POSTs to Dagger's API
- **Landing page** — static site explaining what Dagger is, with API reference
- **Deployment** — Fly.io or VPS with Docker Compose

### Dependencies

- PostgreSQL schema finalized
- DAG plugin changes (new adapter type)
- Go project scaffolded

### Success Criteria

- DAG plugin publishes a ticket + linked ADR to Dagger
- `curl` retrieves that ticket with full context as markdown
- Claude Code via MCP can fetch ticket context and generate code from it

---

## Phase 2 (Months 3-4): Web UI

**Goal:** Humans can browse, search, and manage the doc/ticket graph.

### Deliverables

- **TanStack Start web app** — connects to the same Go API
- **User accounts** — GitHub OAuth, email + password, session management
- **Teams and workspaces** — multi-user, multi-project
- **Doc/ticket browser** — hierarchical tree view of the decision graph
- **Ticket detail view** — clean reading view with linked context sidebar
- **Full-text search** — across all docs and tickets
- **Manual CRUD** — create and edit tickets/docs through the UI
- **Visual design system** — UI language, component library setup
- **Agent API polish** — rate limiting, pagination, response caching

### Dependencies

- TanStack Start project scaffolded
- Go API may need updates for user management endpoints

### Success Criteria

- User signs up, creates a workspace
- Browses tickets and docs from the web UI
- Creates a new ticket manually
- Search returns relevant results

---

## Phase 3 (Months 5-6): DAG-powered Chat

**Goal:** Full feature lifecycle from conversation to published context, all in the browser.

### Deliverables

- **Chat interface** — SSE streaming chat UI
- **BYOK LLM integration** — abstraction over Anthropic, OpenAI, Google APIs
- **DAG workflow prompts ported to web** — plan-feature, create-tickets, code-exploration as chat-driven workflows
- **Meeting transcript analysis** — identify stakeholder/product/engineering voices
- **Conversation threads and history** — persistent, searchable
- **pgvector** — semantic search across all docs for chat context
- **Agent notification** — webhooks or polling when tickets are created/updated
- **Dag-research integration** — prompt templates optimized per-model based on research findings

### Dependencies

- BYOK integration built and tested
- DAG plugin prompts ported to Go/TS
- pgvector extension enabled in PostgreSQL

### Success Criteria

- User pastes a meeting transcript → Dagger produces ADR + tickets
- User has a planning conversation → Dagger guides through structured decision-making
- Agent fetches the resulting ticket with full context → generates correct code
- Token usage is transparent and predictable

---

## Post-Phase 3

- Versioning for docs and tickets (history, diffs, rollbacks)
- Import from existing tools (Linear, Jira, GitHub Issues)
- Bundled model pricing (absorb LLM costs into subscription)
- Enterprise features (SSO, RBAC, audit logs)
- On-prem deployment option
