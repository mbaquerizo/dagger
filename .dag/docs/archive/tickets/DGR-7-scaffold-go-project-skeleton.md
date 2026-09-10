---
id: DGR-7
issueType: task
status: done
blockedBy: []
tags:
  - phase-1
  - go
  - project-setup
parent: DGR-5
---

# Scaffold Go project skeleton

## Description

Initialize the Go project with the standard `cmd/` + `internal/` layout, a working HTTP server with a health check endpoint, and all the tooling needed for development. This establishes the project structure that every subsequent ticket builds on.

## Acceptance criteria

1. `go mod init` creates the module, `go build ./cmd/api` produces a binary
2. HTTP server starts on `:8080` with a `GET /healthz` endpoint returning `200 OK`
3. Project layout follows standard Go conventions: `cmd/api/main.go`, `internal/` packages flat
4. Dependencies declared: `chi` router, `pgx` PostgreSQL driver
5. Makefile with targets: `build`, `run`, `test`, `lint`
6. `.gitignore` excludes Go binaries, IDE files, environment files
7. README updated with prerequisites, setup, and run instructions

## Related docs

- Product plan: `.dag/docs/plan/dagger-plan/index.md`
- Go project layout reference: `github.com/golang-standards/project-layout`
- Go learning plan: `.dag/docs/plan/learning-go/README.md`

---

*This ticket was created by opencode and reviewed by a human before publishing.*
