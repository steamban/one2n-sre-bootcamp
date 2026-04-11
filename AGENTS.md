# Agents

## Project Overview

Go/Postgres REST API for student CRUD. Rewrite to TypeScript is planned.

## Key Commands

```bash
cp example.env .env           # Required before running
make dev                     # Start server (go run cmd/server/main.go)
make build                   # Build binaries to bin/
make test                    # Run tests (no test files exist yet)
go run cmd/migrate/main.go up   # Run migrations
go run cmd/migrate/main.go down # Rollback migrations
```

## Architecture

- Entry point: `cmd/server/main.go`
- Handlers: `internal/api/handler/handler.go` (Gin-based)
- Repository: `internal/repository/student.go` (raw SQL, soft-delete)
- Model: `internal/model/student.go`
- Config: `internal/config/config.go` (reads from `.env`)
- Migrations: `migrations/*.sql` (000001*, 000002*)

## Routes

- `GET /healthcheck`
- `POST /api/v1/students`
- `GET /api/v1/students`
- `GET /api/v1/students/:id`
- `PATCH /api/v1/students/:id`
- `DELETE /api/v1/students/:id`

## Migration Notes

Migrations run via separate binary (`cmd/migrate/main.go`), NOT on server startup. Server only calls `db.InitDB()` without auto-migrating.

## Workflow Rules

- Do not start code without telling the plan and getting a go ahead
- Follow conventional commits (`feat:`, `fix:`, `refactor:`, etc.)
- Only make atomic commits (one logical change per commit)
- Do not create any other git branches — keep all changes within the same branch
- Do not push branches — you will do that after reviewing code
- Stick with best practices and 12-factor app methodology
- Only make minimal required changes — do not overcomplicate
- Ask if you have doubts

## Rewrite to TypeScript

- Use **Prisma** — model file is single source of truth, migrations auto-generated via `prisma migrate`
- Do NOT use hand-written SQL migration files; the TypeScript model drives schema
- Target: Express
- Preserve routes and repository pattern from Go version