# Contributing to Dagger

## Prerequisites

- Go 1.26+

## Environment

Copy `.env.example` to `.env`. Only two keys:

- `DB_URL` — Supabase Postgres connection string (pooled, port 6543).
  Used by the API, `migrate`, and `seedkey`.
- `BASE_URL` — public base URL of the API (e.g. `http://localhost:8080`
  locally). Used to build display URLs in publish responses.

No local database is set up — there is no Docker Compose or local
Postgres. All environments (local dev included) connect directly to the
Supabase database via `DB_URL`.

## Setup

```
git clone <repo>
cd dagger
cp .env.example .env  # then fill in real values
go mod tidy
```

## Migrations

Run against whatever `DB_URL` points at (currently the shared Supabase DB —
there is no separate local DB, so `make migrate` touches the shared
database):

```
make migrate
# runs: go run ./cmd/migrate (applies db/migrations in order)
```

## Seed an API key

Prints the raw key once — store it, it is hashed in the DB:

```
go run ./cmd/seedkey --workspace-id 1 --name "dev key"
```

## Run

```
make run
# or: go run ./cmd/api
```

Health check:

```
curl http://localhost:8080/healthz
# → OK
```

Other useful commands: `make build`, `make test`, `make lint`.

## Deployment (Railway + Supabase)

- Supabase project holds Postgres. Copy the pooled connection string into
  `DB_URL`.
- Railway project hosts the Go API service (auto-detects Go, no Dockerfile).
  Set `DB_URL` and `BASE_URL` (public Railway URL) in Railway env vars.
- MCP is served from the same service at `POST /api/v1/agent/mcp`.
