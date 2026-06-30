# Deployment

Silmari deployment models — from local development to production Kubernetes clusters.

## Deployment Models

| Mode | Storage | Dependencies | Use Case |
|------|---------|-------------|----------|
| **Local (dev)** | SQLite + BoltDB | None | Development, testing, single-machine |
| **Docker Compose** | PostgreSQL | Docker | Local production-like testing |
| **PostgreSQL** | PostgreSQL | PostgreSQL 15+ | Production, single-node |
| **Cloud** | Managed PostgreSQL | Cloud provider | Distributed production |
| **Kubernetes** | PostgreSQL | K8s cluster | Orchestrated production |

## Local Mode

Simplest deployment — no external dependencies:

```bash
cd control-plane
go run ./cmd/af dev
# or
go run ./cmd/agentfield-server
```

Control plane runs at `http://localhost:8080`.
Web UI at `http://localhost:8080/ui/`.

### Architecture

```
┌────────────────────────────────────┐
│         Control Plane (Go)         │
│         :8080                      │
│  ┌──────────┐  ┌────────────────┐  │
│  │  SQLite  │  │    BoltDB      │  │
│  │(relations│  │  (key-value)   │  │
│  │  data)   │  │                │  │
│  └──────────┘  └────────────────┘  │
│  ┌────────────────────────────────┐ │
│  │    Embedded Web UI             │ │
│  │    (if AGENTFIELD_UI_ENABLED)  │ │
│  └────────────────────────────────┘ │
└────────────────────────────────────┘
```

- **SQLite:** Stores relational data (agents, workflows, executions, memory)
- **BoltDB:** Stores key-value data (configuration, internal state)
- **Embedded UI:** If `AGENTFIELD_UI_ENABLED=true`, the production-built UI is served from the same process

No database setup required — SQLite and BoltDB files are created on first run.

**Code reference:** `control-plane/cmd/af/` — dev command, `control-plane/internal/storage/` — local storage backends

## PostgreSQL Mode

Production-ready single-node deployment:

```bash
# 1. Start PostgreSQL
# 2. Run migrations
cd control-plane
export AGENTFIELD_DATABASE_URL="postgres://silmari:silmari@localhost:5432/silmari?sslmode=disable"
goose -dir ./migrations postgres "$AGENTFIELD_DATABASE_URL" up

# 3. Start server
AGENTFIELD_STORAGE_MODE=postgresql \
AGENTFIELD_DATABASE_URL="postgres://silmari:silmari@localhost:5432/silmari?sslmode=disable" \
go run ./cmd/agentfield-server
```

### Architecture

```
┌────────────────────────────────────┐
│         Control Plane (Go)         │
│         :8080                      │
│  ┌──────────────────────────────┐  │
│  │     Connection Pool          │  │
│  │  (infrastructure/storage/)   │  │
│  └──────────────┬───────────────┘  │
└─────────────────┼──────────────────┘
                  │
         ┌────────▼────────┐
         │  PostgreSQL 15+  │
         │  (all data)      │
         └─────────────────┘
```

- **PostgreSQL** stores ALL data (no BoltDB in this mode)
- **Connection pooling** via `internal/infrastructure/storage/`
- **Migrations** must be run before first start

Key environment variables:
- `AGENTFIELD_STORAGE_MODE=postgresql`
- `AGENTFIELD_DATABASE_URL` — PostgreSQL connection string
- `AGENTFIELD_STORAGE_POSTGRES_MAX_CONNECTIONS` — pool size

**Code reference:** `control-plane/internal/infrastructure/storage/` — connection pooling, `control-plane/migrations/` — SQL migrations

## Docker Compose

Includes PostgreSQL and the control plane:

```bash
cd deployments/docker
docker compose up
```

### Architecture

```
┌──────────────────────┐
│  agentfield-server   │
│  container :8080     │
└──────────┬───────────┘
           │
┌──────────▼───────────┐
│  postgres container  │
│  :5432               │
└──────────────────────┘
```

**Code reference:** `deployments/docker/` — Docker Compose configuration

## Cloud Mode

For distributed production deployments. Uses managed PostgreSQL with the same application binary:

```
                    ┌──────────────┐
                    │ Load Balancer│
                    └──────┬───────┘
                           │
              ┌────────────┼────────────┐
              │            │            │
    ┌─────────▼──┐ ┌──────▼─────┐ ┌─────▼────────┐
    │ CP Instance│ │CP Instance │ │ CP Instance  │
    │ :8080      │ │:8080       │ │ :8080        │
    └─────┬──────┘ └──────┬─────┘ └──────┬───────┘
          │               │              │
          └───────────────┼──────────────┘
                          │
                 ┌────────▼────────┐
                 │ Managed PG     │
                 │ (Cloud SQL /   │
                 │  RDS / Supabase)│
                 └────────────────┘
```

Each control plane instance is stateless — all state is in PostgreSQL. Horizontal scaling is limited only by PostgreSQL connection capacity.

**Code reference:** `control-plane/internal/infrastructure/storage/` — connection pooling for cloud

## Kubernetes

For orchestrated production deployments. The control plane runs as a Kubernetes Deployment with PostgreSQL as a StatefulSet or external managed service.

Key considerations:
- Control plane pods are stateless → safe to scale horizontally
- PostgreSQL needs persistent volume or external managed service
- Web UI is embedded in control plane binary → no separate frontend deployment
- Agents run as separate pods or external services, registering with control plane via HTTP

## Configuration Reference

See `control-plane/.env.example` for the complete list. Key deployment variables:

| Variable | Values | Default | Description |
|----------|--------|---------|-------------|
| `AGENTFIELD_PORT` | integer | `8080` | HTTP server port |
| `AGENTFIELD_MODE` | `local`, `cloud` | `local` | Operating mode |
| `AGENTFIELD_STORAGE_MODE` | `local`, `postgresql`, `cloud` | `local` | Storage backend |
| `AGENTFIELD_DATABASE_URL` | connection string | — | PostgreSQL connection |
| `AGENTFIELD_UI_ENABLED` | `true`, `false` | `true` | Enable embedded UI |
| `AGENTFIELD_UI_MODE` | `embedded`, `development` | `embedded` | UI serving mode |
| `AGENTFIELD_CONFIG_FILE` | file path | — | Override the default `config/agentfield.yaml` path while keeping the legacy env var name stable |
| `GIN_MODE` | `debug`, `release` | `debug` | Gin framework mode |
| `LOG_LEVEL` | `debug`, `info`, `warn`, `error` | `info` | Log verbosity |
| `AGENTFIELD_STORAGE_POSTGRES_MAX_CONNECTIONS` | integer | — | DB pool size |

## Build & Release

### Building

```bash
make build                 # All components
make control-plane         # Control plane only
./control-plane/build-single-binary.sh  # Unified binary with embedded UI
```

### Release Process

Releases via GoReleaser (`.goreleaser.yml`) and GitHub Actions (`.github/workflows/release.yml`):

```bash
git tag v0.1.0
git push origin v0.1.0
```

GitHub Actions builds binaries for Linux, macOS, and Windows. The unified binary includes the embedded web UI.

## Serverless Agents

Agents can run in serverless environments (AWS Lambda, Cloudflare Workers, etc.) using the SDK's serverless adapter:

**Python:** `sdk/python/agentfield/agent_serverless.py`

The control plane communicates with serverless agents via HTTP callbacks rather than persistent connections. The serverless adapter handles the function lifecycle (cold start, invocation, teardown) and translates between the serverless event model and Silmari's execution model.

**Code reference:** `sdk/python/agentfield/agent_serverless.py`
