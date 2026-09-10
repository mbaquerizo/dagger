---
id: DGR-15
issueType: story
status: in-progress
tags:
  - phase-1
  - dag-plugin
  - adapter
parent: DGR-5
blockedBy:
  - DGR-9
---

# DAG "dagger" adapter

## Description

As a DAG plugin user, I want to publish docs and tickets to Dagger so they're available to AI agents via Dagger's API. This is a new adapter type in the DAG plugin repo — a separate project from Dagger itself.

The DAG plugin repo is a Claude Code plugin (markdown skills + `config.schema.json` + `scripts/dag-init.sh`). The adapter is new adapter docs plus a schema entry, executed by the agent following published instructions.

## Acceptance criteria

### Registration
1. `config.schema.json` gains a `dagger` variant under both `adapters.documentation` and `adapters.issue-tracking` (alongside `local`/`db`).
2. New `documentation-dagger.md` + `issue-tracking-dagger.md` in `plugins/dag/skills/dag-publish/`, following the existing adapter-doc structure (config table, contract mapping, step-by-step instructions, examples). `SKILL.md` resolution already handles external types by name — no change needed there.

### Config
3. `url` (Dagger base URL) + Dagger `project_id` (publish requires it) live in `.dag/config.json`, both required with descriptive failures per `adapter-conventions.md` §5.
4. API key from `DAGGER_API_KEY` env only — never committed to config.

### HTTP path
5. `POST {url}/api/v1/publish` with `Authorization: Bearer <key>`.
6. Capture `id` (internal, for subsequent relationship refs) + `displayId` + `url` from the response; report `displayId`/`url` to the user.
7. Publish parents first, children referencing returned internal IDs.
8. Errors mapped per `adapter-conventions.md` §6 and the plan docd (§Error Handling: 401 auth / 400 bad ref / 422 validation / 500 retry with backoff).

### MCP path
9. Same flow via the `publish` tool; reads via `get_issue` / `get_doc` / `list_issues`.
10. Issue↔issue edges via follow-up `add_issue_relation` calls — one perdone-way relation (each call writes both directions). Single-call `issue_relations` over MCP is not supported (schema + arg parsing lack it); the HTTP path supports it in one call.

### Contract mapping
11. Adapter docs include a mapping table: `type`/`title`/`body` direct; `parent` (contract grouping string) distinguished from Dagger internal-ID references; `blockedBy`/`children` → `add_issue_relation` (MCP) or `issue_relations` (HTTP); `status`/`tags`/`issueType` → `metadata`.

## Out of scope (follow-ups)
- Updating existing docs in place (no Dagger API for it).
- MCP path requires MCP configured before first publish (bootstrap note).
- `dag-init` dagger option: `screen_adapter` in `scripts/dag-init.sh` is hardcoded local-only ("Cloud adapters coming later") with only a `local` branch in `generate_config_json`. Offering dagger there (prompts for `url` + `project_id`, env note for the key) is a separate change in the DAG repo.

## Technical notes
- Work spans the DAG plugin repo (adapter docs, schema, later `dag-init`), not this repo.
- See `.dag/docs/plan/dagger-plan/04-adapter-design.md` for the contract. Note the plan doc predates the MCP server: where it says HTTP-only, this ticket now covers both transports per above.
