# AgentField Control Plane

The AgentField control plane orchestrates agent workflows, manages verifiable credentials, serves the admin UI, and exposes REST/gRPC APIs consumed by the SDKs.

## Which path is for you?

These are three separate setups for three goals — not competing ways to install the same thing:

| Goal | Use | Where |
|---|---|---|
| Build & run agents / try AgentField | the `af` CLI: `curl -fsSL https://agentfield.ai/install.sh \| bash`, then `af init` / `af dev` | root [`README.md`](../README.md) |
| Run the control plane (ops) | Docker: `docker run -p 8080:8080 agentfield/control-plane:latest`, or `docker compose -f deployments/docker/docker-compose.yml up` | [`deployments/docker/`](../deployments/docker) |
| Develop the control plane (this guide) | build from source — Quick Start below | this file + [`docs/DEVELOPMENT.md`](../docs/DEVELOPMENT.md) |

> `af` and the Docker image are different layers, not alternatives: `af` is the CLI you use to scaffold and run agents; the control plane is the Go service those agents talk to. Run the control plane via `af dev` (local), Docker, or from source below.

## Requirements

- Go 1.23+
- Node.js 20+ (for the web UI under `web/client`)
- PostgreSQL 15+

## Quick Start

```bash
# From the repository root
cd control-plane
go mod download

# Build the embedded web UI
( cd web/client && npm install && npm run build )

# Stand up PostgreSQL matching the default DSN below
# (skip this if you use SQLite mode — see the note after this block)
docker run -d --name agentfield-pg -p 5432:5432 \
  -e POSTGRES_USER=agentfield -e POSTGRES_PASSWORD=agentfield -e POSTGRES_DB=agentfield \
  pgvector/pgvector:pg16

# Run database migrations
export AGENTFIELD_DATABASE_URL="postgres://agentfield:agentfield@localhost:5432/agentfield?sslmode=disable"
goose -dir ./migrations postgres "$AGENTFIELD_DATABASE_URL" up

# Start the control plane
go run ./cmd/agentfield-server
```

Visit `http://localhost:8080/ui/` to access the embedded admin UI.

> **No PostgreSQL?** Run `./dev.sh` instead for SQLite mode (default — no database to provision, no migrations needed); see [Local Development](#local-development-with-hot-reload) below. The DSN above is the project default and matches the `pgvector` container in `deployments/docker/docker-compose.yml` — it is not a placeholder.

## Local Development with Hot-Reload

For development with hot-reload, use the `dev.sh` script. This automatically rebuilds and restarts the server when Go files change.

```bash
cd control-plane
./dev.sh            # SQLite mode (default, no dependencies)
./dev.sh postgres   # PostgreSQL mode (set AGENTFIELD_DATABASE_URL first)
```

The server runs at `http://localhost:8080` and will automatically reload when you modify `.go`, `.yaml`, or `.yml` files.

**Notes:**
- Uses [Air](https://github.com/air-verse/air) for hot-reload (auto-installed if missing)
- Web UI is not included in dev mode; run `npm run dev` separately in `web/client/` if needed

## Configuration

Environment variables override `config/agentfield.yaml`. Common options:

- `AGENTFIELD_DATABASE_URL` – PostgreSQL DSN
- `AGENTFIELD_HTTP_ADDR` – HTTP listen address (`0.0.0.0:8080` by default)
- `AGENTFIELD_LOG_LEVEL` – log verbosity (`info`, `debug`, etc.)

Sample config files live in `config/`.

## Web UI Development

```bash
cd control-plane/web/client
npm install
npm run dev
# Build production assets embedded in Go binaries (run from web/client)
npm run build
```

Run the Go server alongside the UI so API calls resolve locally. During production builds the UI is embedded via Go's `embed` package.

## Database Migrations

Migrations use [Goose](https://github.com/pressly/goose):

```bash
AGENTFIELD_DATABASE_URL=postgres://agentfield:agentfield@localhost:5432/agentfield?sslmode=disable \
goose -dir ./migrations postgres "$AGENTFIELD_DATABASE_URL" status
```

## Testing

```bash
go test ./...
```

## Linting

```bash
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
golangci-lint run
```

## Releases

The `build-single-binary.sh` script creates platform-specific binaries and README artifacts. CI-driven releases are defined in `.github/workflows/release.yml`.
