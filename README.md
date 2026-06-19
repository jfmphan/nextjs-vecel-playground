# Home Inventory

A lightweight app to keep track of things around the house — items, where they
live, how many you have, and what's running low or expiring. Reachable from
anywhere, runs entirely on free tiers.

## Stack

| Layer    | Choice                                                              |
| -------- | ------------------------------------------------------------------ |
| Frontend | Next.js (App Router) + React + [Mantine](https://mantine.dev)      |
| Backend  | Go serverless function (chi router)                                |
| Database | SQLite locally / [Turso](https://turso.tech) (libSQL) in prod      |
| Auth     | Single shared password → signed, stateless session cookie          |
| Hosting  | Vercel (frontend + Go function), free tier                         |

The frontend and the Go API both deploy to Vercel; the API talks to Turso, a
SQLite-compatible cloud database. Because the SQL dialect is identical, local
development uses a plain SQLite file (pure-Go driver, zero setup) and only the
connection string changes in production.

## Architecture

```
Browser ──/──────────────▶ Next.js (Mantine UI)        app/, components/, lib/
        └─/api/v1/* ──────▶ Go function (chi router)    api/index.go → internal/
                                  └─▶ Turso / SQLite

internal/
  domain/    entities + value objects + repository interfaces (no infra deps)
  service/   application use cases (validation, business rules)
  storage/   libSQL/SQLite repository implementations + migrations
  httpapi/   HTTP handlers, DTO mapping, auth middleware, router
  auth/      shared-password check + signed session cookie
  config/    environment configuration
  app/       composition root (wires everything together)
```

Dependencies point inward: `httpapi → service → domain ← storage`. The domain
layer knows nothing about HTTP or SQL, so services are unit-testable with fake
repositories and the storage backend can change without touching business logic.

## Prerequisites

- [mise](https://mise.jdx.dev) (manages the project-local Go and Node versions)

```bash
mise install   # provisions Go and Node as pinned in mise.toml
```

All commands below are prefixed with `mise exec --` (or run them directly if you
have mise activated in your shell).

## Local development

The frontend and backend run as two processes. The Go API listens on `:8080`;
`next dev` proxies `/api/*` to it (configured via `API_PROXY_TARGET` in
`mise.toml`) so the browser stays same-origin and the session cookie works.

```bash
# Terminal 1 — API (auto-creates and migrates ./data/inventory.db)
APP_PASSWORD=changeme mise exec -- go run ./cmd/server

# Terminal 2 — web app
mise exec -- npm run dev
```

Open http://localhost:3000 and sign in with the password from `APP_PASSWORD`
(defaults to `changeme`).

### Useful commands

```bash
mise exec -- go vet ./...        # vet backend
mise exec -- go test ./...       # backend tests
mise exec -- npm run typecheck   # frontend type check
mise exec -- npm run build       # production build of the frontend
mise exec -- go run ./cmd/migrate  # apply DB migrations (used for remote/Turso)
```

## Environment variables

See `.env.example`. The Go API reads:

| Variable              | Purpose                                              | Default (local)            |
| --------------------- | ---------------------------------------------------- | -------------------------- |
| `DATABASE_URL`        | `file:...` (SQLite) or `libsql://...` (Turso)        | a file under `./data`      |
| `DATABASE_AUTH_TOKEN` | Turso auth token (remote only)                       | —                          |
| `APP_PASSWORD`        | the shared password gating the app                   | `changeme`                 |
| `SESSION_SECRET`      | HMAC key signing the session cookie                  | a dev placeholder          |
| `AUTO_MIGRATE`        | force migrations on startup for a remote DB          | on for file DBs only       |

## Deployment (Vercel + Turso)

The fastest path is the **Vercel plugin for coding agents** (already installed):
restart your agent and run `/vercel-plugin:bootstrap`, which links the project,
provisions env vars, and deploys. Manual steps for reference:

1. **Create the Turso database** and apply migrations:
   ```bash
   turso db create home-inventory
   turso db show home-inventory --url        # → DATABASE_URL
   turso db tokens create home-inventory     # → DATABASE_AUTH_TOKEN
   DATABASE_URL=<url> DATABASE_AUTH_TOKEN=<token> mise exec -- go run ./cmd/migrate
   ```
2. **Create a Vercel project** from this repo and set the env vars:
   `DATABASE_URL`, `DATABASE_AUTH_TOKEN`, `APP_PASSWORD`, `SESSION_SECRET`.
3. **Deploy.** `vercel.json` rewrites `/api/v1/*` to the Go function
   (`api/index.go`); Next.js serves everything else. The app is reachable over
   HTTPS and gated by the shared password.

## Features

See [FEATURES.md](./FEATURES.md).
