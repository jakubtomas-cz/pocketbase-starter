# PocketBase Starter

A production-ready Go starter for [PocketBase](https://pocketbase.io) projects with a clean, layered architecture.

## Overview

PocketBase Starter gives you a well-organized foundation for building applications on top of PocketBase — with proper Go package structure, server-side templating, migration support, Docker, and a hot-reload dev workflow out of the box.

## Project Structure

```
.
├── cmd/
│   └── server/
│       └── main.go          # Entry point
├── internal/
│   ├── app/
│   │   └── bootstrap.go     # Wires plugins, hooks, and routes onto PocketBase
│   ├── hooks/
│   │   └── records.go       # Record lifecycle hooks
│   └── routes/
│       └── api.go           # API routes and server-side view routes
├── migrations/
│   └── migrations.go        # Go-based database migrations
├── views/
│   └── index.html           # Server-side HTML templates
├── pb_public/               # Static assets (SPA-friendly)
├── Dockerfile
└── Makefile
```

### Layer Responsibilities

| Layer | Package | Purpose |
|---|---|---|
| Entry point | `cmd/server` | Creates PocketBase instance, registers app, starts server |
| Bootstrap | `internal/app` | Registers CLI flags, plugins, hooks, and routes |
| Hooks | `internal/hooks` | Record and auth lifecycle logic |
| Routes | `internal/routes` | Custom API endpoints and server-rendered views |
| Migrations | `migrations` | Schema and data migrations (auto-applied on startup) |

## Getting Started

**Requirements:** Go 1.25+

```bash
# Clone and install dependencies
git clone https://github.com/jakubtomas-cz/pocketbase-starter.git
cd pocketbase-starter
go mod download

# Run
make run
```

The server starts at `http://localhost:8090`. The admin UI is available at `http://localhost:8090/_/`.

## Development

Hot-reload via [nodemon](https://nodemon.io/) (requires Node.js):

```bash
make dev
```

Watches all `.go` files and restarts on change, ignoring `pb_data/`, `pb_public/`, and `bin/`.

## Building

```bash
make build      # Compiles to ./app
```

## Docker

```bash
# Build image
make docker IMAGE_TAG=latest

# Run container
docker run -p 8090:8090 -v $(pwd)/pb_data:/app/pb_data pocketbase-starter:latest
```

The image uses a two-stage build (Go builder → Alpine production). `pb_data` is exposed as a volume for persistence.

## Migrations

Migrations live in the `migrations/` package and are auto-applied on startup.

Generate a new migration:

```bash
go run ./cmd/server/main.go migrate create "your_migration_name"
```

Each migration file registers itself via `init()` — see [PocketBase migration docs](https://pocketbase.io/docs/go-migrations/) for the full API.

## Adding Routes

**API route** in `internal/routes/api.go`:

```go
func RegisterAPIs(pb core.App) {
    pb.OnServe().BindFunc(func(se *core.ServeEvent) error {
        se.Router.GET("/api/hello/{name}", func(e *core.RequestEvent) error {
            name := e.Request.PathValue("name")
            return e.JSON(http.StatusOK, map[string]string{"message": "Hello, " + name})
        })
        return se.Next()
    })
}
```

**Server-rendered view** using PocketBase's template engine:

```go
func RegisterViews(pb core.App) {
    pb.OnServe().BindFunc(func(se *core.ServeEvent) error {
        se.Router.GET("/{$}", func(e *core.RequestEvent) error {
            registry := template.NewRegistry()
            html, err := registry.LoadFiles("views/index.html").Render(map[string]any{})
            if err != nil {
                return err
            }
            return e.HTML(http.StatusOK, html)
        })
        return se.Next()
    })
}
```

## Adding Hooks

Hooks go in `internal/hooks/records.go` (or split into additional files within the `hooks` package):

```go
func Register(pb core.App) {
    pb.OnRecordCreateRequest().BindFunc(func(e *core.RecordRequestEvent) error {
        log.Printf("[hook] record create: collection=%s", e.Collection.Name)
        return e.Next()
    })
}
```

See the [PocketBase hooks reference](https://pocketbase.io/docs/go-event-hooks/) for all available events.

## Extending the Starter

As your project grows, here are conventional places to add new code:

| Folder / Package | Purpose |
|---|---|
| `internal/pages/` | UI-related route handlers (server-rendered pages, form submissions) — keeps view logic separate from API routes |
| `internal/middlewares/` | Custom middleware (auth guards, rate limiting, request logging) — register these in `bootstrap.go` |
| `internal/models/` | Domain structs and helpers that wrap PocketBase records with typed fields and business logic |
| `internal/jobs/` | Cron jobs and scheduled background tasks — wire them up in `bootstrap.go` using PocketBase's scheduler |

None of these packages exist in the starter by default — create them when you need them and register their entry functions in `internal/app/bootstrap.go`.

## Static Files (SPA)

Drop files into `pb_public/` to serve them statically. Index fallback is enabled by default, so SPAs with client-side routing work without extra configuration.

### Full-Stack in One Repo

You can scaffold a frontend app (React, Vue, Svelte, etc.) directly inside this repository and configure its build output to `pb_public/`. This keeps your entire stack — backend, frontend, and database migrations — in one place.

For example, with Vite:

```bash
npm create vite@latest ui -- --template react
```

Then set the build output in `ui/vite.config.ts`:

```ts
export default {
  build: {
    outDir: '../pb_public',
    emptyOutDir: true,
  },
}
```

Run `npm run build` inside `ui/` and PocketBase will serve the compiled app from `pb_public/` automatically. During development, run the Vite dev server alongside `make dev` and proxy API requests to `http://localhost:8090`.

## CLI Flags

| Flag | Default | Description |
|---|---|---|
| `--migrationsDir` | `""` | Path to additional migrations directory |
| `--automigrate` | `true` | Auto-apply migrations on startup |
| `--hooksDir` | `""` | Directory for JS app hooks (`pb_hooks`) |
| `--hooksWatch` | `true` | Auto-restart on JS hook file changes |
| `--hooksPool` | `15` | Prewarm pool size for JS hook runtimes |
| `--publicDir` | `./pb_public` | Directory for static file serving |
| `--indexFallback` | `true` | Fallback to `index.html` for SPA routing |

## Plugins Included

- **jsvm** — JavaScript hooks and migrations via `pb_hooks/`
- **migratecmd** — `migrate` CLI command with Go template generation
- **ghupdate** — `update` CLI command for GitHub-hosted binary releases

## License

MIT © [Jakub Tomáš (jakubtomas-cz)](https://github.com/jakubtomas-cz)
