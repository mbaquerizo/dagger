# Business Model

## Phases

### Phase 1-2: Free

No monetization. Prove the workflow works. Delight early users.

- Dagger is free to use
- Users bring their own model API key (BYOK) if they use the chat
- Single user, no team features yet
- Goal: validate the agent API value prop, gather usage data

### Phase 3 Launch: Subscription

Introduce paid plans when the chat launches.

| Tier | Target | Price | Chat model | Limits |
|------|--------|-------|------------|--------|
| Free | Individuals, evaluation | $0 | BYOK only | 1 user, 3 projects |
| Team | Small engineering teams | $20/seat/mo | BYOK | Unlimited projects, teams |
| Business | Engineering orgs | $40/seat/mo | BYOK + bundled optional | SSO, audit logs, support |

### Post-Launch: Bundled Tier

Introduce a bundled model tier where Dagger carries the LLM cost.

- Chat model costs are ~$0.30/session at current pricing
- At 20 sessions/month/seat: ~$6/model cost per seat
- Easily absorbed into a $40/seat Business plan
- BYOK remains as an enterprise option

## Why BYOK First

1. **Zero marginal cost risk** — no surprise bills while proving the workflow
2. **Early adopters have model subscriptions** — the target audience already pays for Claude/GPT
3. **Simpler initial product** — no provider integration, no billing for model usage
4. **Usage data** — when bundled launches, pricing will be informed by real usage patterns

## Why Bundled Later

1. **Single bill is simpler for most users** — one subscription, everything included
2. **Dagger controls the full stack** — can optimize prompts per-model to reduce token costs
3. **Token optimization compounds** — DAG-shaped context reduces codegen iterations → user saves more on their Claude Code bill than Dagger costs in chat tokens
4. **Competitive pricing** — $20-40/seat is standard for dev tools (Linear: $8-14, Notion AI: $10, Cursor: $20)

## Go-to-Market

1. **Phase 1** — Personal network, DAG plugin users, developer Twitter
2. **Phase 2** — Open source the core? Publish on ProductHunt, Hacker News
3. **Phase 3** — Content marketing (dag-research findings as blog posts), engineering orgs

## Cost Structure Estimation (at 100 orgs, 10 seats each)

| Cost | Monthly |
|------|---------|
| PostgreSQL (Railway/Fly) | ~$50 |
| API server hosting | ~$50 |
| TanStack Start hosting | ~$50 |
| Domain, email, misc | ~$20 |
| **Total infra** | **~$170** |
| Per-org cost | ~$0.17 |

At $20/seat: $20,000/mo revenue vs ~$170 infra. Margins are exceptional because the heavy compute (LLM calls) is on the user's side during Phase 1-2, and still absorbed by subscription pricing in Phase 3.
