# Go Ramp-Up Plan

**Background:** TypeScript/Node, some Java and Rust exposure but never built anything with them.
**Pace:** ~5-10 hours/week.
**Goal:** Comfortable enough to build Dagger's Phase 1 API server.

---

## Phase 0: Structured Learning

Go through these resources before touching the project. No rush — go at your own pace.

- [ ] **Go Tour** ([tour.golang.org](https://go.dev/tour/)) — syntax, types, control flow, concurrency basics (~2-3h)
- [ ] **Learn Go with Tests** ([quii.gitbook.io](https://quii.gitbook.io/learn-go-with-tests)) — focus chapters (~10-15h):
  - [ ] Hello, World
  - [ ] Integers
  - [ ] Arrays and Slices
  - [ ] Structs, Methods, Interfaces
  - [ ] Pointers and Errors
  - [ ] Maps
  - [ ] Dependency Injection
  - [ ] Mocking
  - [ ] HTTP Server
  - [ ] HTTP Handlers (Routing with `net/http`)
- [ ] **Go by Example** ([gobyexample.com](https://gobyexample.com/)) — reference as needed, not a front-to-back read

---

## Phase 1: Dagger Scaffolding

Build these in order. Each is a milestone that teaches a new concept. Project layout TBD when you get there.

- [ ] **Project skeleton** — `go mod init`, `cmd/api/main.go`, `internal/` packages, HTTP server on `:8080`
- [ ] **Config & DB connection** — env-var config, `pgx` connection pool, SQL migration runner
- [ ] **Models** — Go structs for docs, issues, relationships (translating from DB adapter schema)
- [ ] **POST /api/v1/publish** — JSON handler, validate input, insert to DB
- [ ] **GET /api/v1/tickets/:id** — path param, assemble context, return markdown
- [ ] **GET /api/v1/documents/:id** — individual doc retrieval
- [ ] **API key auth** — `Authorization` header, bcrypt hash, middleware pattern
- [ ] **Dockerfile + docker-compose.yml** — Go binary + PostgreSQL, one `docker compose up`
- [ ] **MCP server** — thin Go wrapper around the REST API (stdlib HTTP client)
- [ ] **Integration tests** — spin up test DB, hit endpoints, verify end-to-end

---

## Key TypeScript → Go Mappings

| TypeScript | Go |
|---|---|
| `interface Foo { ... }` | `type Foo struct { ... }` |
| `const foo = (x: T): R => ...` | `func foo(x T) R { ... }` |
| `async/await` | Not needed initially. Handlers call DB synchronously. Goroutines come in Phase 3. |
| `throw / catch` | `return err` everywhere. `if err != nil` is your new `try/catch`. |
| `null` / `undefined` | Zero values (`""`, `0`, `nil`). `*string` for "absent" fields. |
| `?.foo ?? default` | Explicit nil check + `if` block |
| Prisma / ORM | Raw SQL with `pgx`. Go philosophy is explicit SQL. |
| `npm install foo` | `go get github.com/foo/bar` |
| `tsc` / `vite build` | `go build ./cmd/api` — one command, no config file |
| Jest / Vitest | `go test ./...` — built-in, no assertion library |
| `import x from 'y'` | `import "module/path"` |
| `export` | Capitalized names are exported. Lowercase = private. |

## Common Rookie Traps

- **Don't use `interface{}`** — use `any` (Go 1.18+) or better, concrete types
- **Don't use goroutines yet** — start synchronous. Goroutines + channels come in Phase 3 (chat/SSE).
- **Don't reach for a framework** — `chi` router + `pgx` + stdlib `net/http` is the sweet spot
- **Don't fight `if err != nil`** — it's verbose but honest. Every function documents exactly where it can fail.
- **Don't overuse `json.RawMessage`** — define proper structs and let the stdlib marshal/unmarshal
- **Don't nest packages** — flat `internal/` packages are fine. Deep nesting is a Go anti-pattern.
- **Don't use `init()` functions** — explicit initialization is clearer

## Resource Links

| Resource | URL |
|---|---|
| Go Tour | https://go.dev/tour/ |
| Learn Go with Tests | https://quii.gitbook.io/learn-go-with-tests |
| Go by Example | https://gobyexample.com/ |
| Effective Go | https://go.dev/doc/effective_go |
| chi router | https://github.com/go-chi/chi |
| pgx (Postgres driver) | https://github.com/jackc/pgx |
| Standard Go project layout | https://github.com/golang-standards/project-layout |
| Let's Go (book, paid) | https://lets-go.alexedwards.info/ |
| For the Love of Go (book, paid) | https://fortheloveofgo.com/ |

---

*Decide project layout when you get there. The standard `cmd/` + `internal/` pattern will feel natural after the first HTTP handler.*
