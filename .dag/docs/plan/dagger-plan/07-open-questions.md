# Open Questions & Risks

## Pre-Phase 1

### 1. DAG adapter design

The DAG plugin needs a new adapter to publish to Dagger. Should it publish directly via HTTP (cleaner for SaaS) or should Dagger consume the local filesystem output?

**Decision:** HTTP adapter. Filesystem consumption is fragile for remote deployment. Needs DAG plugin changes.

### 2. DAG plugin end state

If the DAG plugin's prompting eventually folds into Dagger's chat backend (Phase 3), how does the plugin transition from standalone CLI tool to internal module? License implications for the plugin's open source code being used in a commercial product?

### 3. Go ramp-up

Phase 1 timeline depends on Go learning curve. How much buffer? First Go project size and scope should be calibrated — start with a small, well-scoped API rather than a full framework.

### 4. TanStack Start maturity

TanStack Start is v1 RC. What's the fallback if framework bugs block development? React Router v7 (framework mode) is the most natural backup — same React ecosystem, similar mental model.

## Pre-Phase 2

### 5. Frontend framework confirmation

Confirmed: TanStack Start. But worth re-evaluating at Phase 2 start depending on community maturity at that point.

### 6. Graph visualization scope

Simple hyperlinks on the ticket detail view (Phase 2) vs interactive graph visualization (post-Phase 3). Current plan: hyperlinks. The graph visualization is a differentiator but not required for the MVP.

### 7. Vector DB timing

pgvector for semantic search: Phase 3 (for chat context) or Phase 2 (for web UI search)? Basic full-text search via PostgreSQL `tsvector` covers Phase 2 needs. pgvector is Phase 3.

### 8. Open source decision

Personal tool vs open core vs closed source. Affects:
- Licensing (MIT/Apache vs proprietary)
- Contribution model
- Community building strategy
- Competitive positioning (open source tools have different adoption patterns)

Decision not urgent. Can be deferred to Phase 2.

## Pre-Phase 3

### 9. BYOK provider support

Phase 3 needs: which providers at launch? Anthropic (Claude) and OpenAI (GPT) are table stakes. Google (Gemini) and any others depend on dag-research findings.

### 10. Chat cost transparency

BYOK means users pay for chat tokens. If a user runs 100 sessions, they might get a surprise bill from their AI provider. Dagger should:
- Show estimated token usage before each session
- Show cumulative usage in the dashboard
- Send alerts at configurable thresholds

### 11. Import from existing tools

Linear, Jira, GitHub Issues import is table-stakes for org adoption. Phase 3+ or earlier?
- Without import: teams start fresh (high friction)
- With import: teams can try Dagger without committing (lower friction)
- Consider building import as a Phase 2.5 feature if early users request it

## Cross-Phase

### 12. API versioning strategy

Agent API versioning from day 1 (`/api/v1/`). What triggers a version bump? Breaking changes to the markdown output format or the relationship model.

### 13. Agent notification

How does an agent know a ticket is ready? Polling (simple) vs webhooks (reliable). Polling is fine for Phase 1. Webhooks can be Phase 2+.

### 14. Versioning for docs/tickets

Track history from Phase 1 or defer? Deferred to post-Phase 3. Versioning adds significant schema and API complexity.

### 15. Deployment as Phase 1 scope

Not confirmed. Options:
- Fly.io — minimal config, managed Postgres, simple deploys
- Docker Compose on VPS — more control, lower cost, more ops burden

Decision: start with Docker Compose for maximum control and zero vendor dependency. Graduate to Pulumi when scaling requires it.
