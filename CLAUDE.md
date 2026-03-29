# CLAUDE.md

## Commands

- `make run` — run the server (localhost:8090)
- `make dev` — hot-reload dev mode via nodemon (requires Node.js)
- `make build` — compile to `bin/app`
- `make docker` — build Docker image (requires `IMAGE_TAG` env var)
- `go run ./cmd/server/main.go migrate create "name"` — generate a new migration

## Architecture

The entry point (`cmd/server/main.go`) creates a PocketBase instance and passes it to `app.Register`, which wires everything together. New code belongs in one of three layers:

| Layer | Package | Add here when… |
|---|---|---|
| Bootstrap | `internal/app` | registering a new plugin or global middleware |
| Hooks | `internal/hooks` | reacting to PocketBase record/auth lifecycle events |
| Routes | `internal/routes` | adding a custom API endpoint or server-rendered view |

`routes.Register` calls `RegisterAPIs` and `RegisterViews` — keep that split when adding new routes.

## Migrations

- Live in `migrations/` as Go files
- Each file must call `migrate.Register` in its `init()` function (the package is blank-imported from `main.go` to trigger init)
- Auto-applied on startup via `migratecmd`

## Static Files & SPA

- Static files are served from `pb_public/`
- `indexFallback` is enabled by default, so unknown paths fall back to `index.html` — suitable for SPAs

## Notes

- The module name is `pocketbase-starter` (not a domain-style path)
- `internal/app/bootstrap.go` mirrors `examples/base/main.go` from the official PocketBase repo — check there for upstream patterns
- `pb_data/` is the PocketBase data directory (SQLite + file storage); never commit it
