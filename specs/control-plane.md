# Control Plane

The central orchestration server for Silmari. Written in **Go** using the **Gin** web framework. It provides REST/gRPC APIs, workflow execution, agent discovery, cryptographic identity, and observability.

## Entry Points

Two binaries exist in `control-plane/cmd/`:

| Binary | Path | Purpose |
|--------|------|---------|
| `af` | `cmd/af/` | Unified CLI with `dev`, `init`, `run`, and server commands |
| `agentfield-server` | `cmd/agentfield-server/` | Standalone server binary (production) |

The unified CLI is the primary development interface. In production, `agentfield-server` is used directly or via the `af` binary with server subcommands.

**Code reference:** `control-plane/cmd/af/` — CLI entry, `control-plane/cmd/agentfield-server/` — server entry

## Request Lifecycle

```
HTTP Request
    │
    ▼
┌──────────────┐
│  Middleware   │  internal/server/middleware/
│  (Auth, CORS,│
│   Logging)   │
└──────┬───────┘
       │
       ▼
┌──────────────┐
│   Handlers   │  internal/handlers/
│  (REST/gRPC) │  - admin/, agentic/, connector/, ui/
└──────┬───────┘
       │
       ▼
┌──────────────┐
│   Services   │  internal/services/
│  (Business   │  - workflow execution, agent registry,
│   Logic)     │    DID/VC generation, workflow status
└──────┬───────┘
       │
       ▼
┌──────────────┐
│   Storage    │  internal/storage/
│  (Persistence│  - SQLite, PostgreSQL, BoltDB
│   Layer)     │
└──────┬───────┘
       │
       ▼
┌──────────────┐
│  Event Bus   │  internal/events/
│  (SSE/Notify)│  - Execution events, node events,
└──────────────┘     trigger events, reasoner events
```

## Server Layer (`internal/server/`)

The HTTP server is built on Gin with middleware for authentication, CORS, logging, and request tracing.

**Key files:**
- `internal/server/middleware/` — Authentication, CORS, request logging middleware
- `internal/server/apicatalog/` — API catalog and discovery endpoints
- `internal/server/knowledgebase/` — Knowledge base endpoints

Routes are registered in the server setup, mapping HTTP methods and paths to handler functions. The server supports both REST and gRPC protocols.

## Handler Layer (`internal/handlers/`)

Handlers are organized by domain:

| Domain | Path | Responsibilities |
|--------|------|-----------------|
| **Agentic** | `handlers/agentic/` | Agent execution, reasoning, tool calls |
| **Admin** | `handlers/admin/` | Administrative operations, system management |
| **Connector** | `handlers/connector/` | External service connectors |
| **UI** | `handlers/ui/` | UI-specific API endpoints |

**Key handler files (top-level):**

| File | Purpose |
|------|---------|
| `handlers/nodes_rest.go` | Agent node registration and lifecycle REST endpoints |
| `handlers/nodes_status.go` | Agent node status reporting |
| `handlers/execute_approval.go` | Human-in-the-loop approval handling |
| `handlers/reasoner_catalog.go` | Reasoner/skill discovery catalog |
| `handlers/verify_audit.go` | Audit trail verification endpoints |
| `handlers/workflow_cleanup.go` | Workflow cleanup and garbage collection |
| `handlers/execution_cleanup.go` | Execution resource cleanup |
| `handlers/agent_concurrency.go` | Concurrency control for agent execution |
| `handlers/vector_memory.go` | Vector memory operations |
| `handlers/memory_events.go` | Memory change event handlers |
| `handlers/errors.go` | Structured error responses |

**Code reference:** `control-plane/internal/handlers/nodes_rest.go` — node registration, `control-plane/internal/handlers/reasoner_catalog.go` — reasoner discovery

## Service Layer (`internal/services/`)

Business logic lives in services. Each service encapsulates a domain concern:

| Service | Path | Purpose |
|---------|------|---------|
| **Workflow Execution** | `services/` | DAG-based workflow orchestration |
| **Agent Registry** | `services/` | Agent registration and lifecycle |
| **DID/VC Generation** | `services/` | Cryptographic identity and credentials |
| **Workflow Status** | `services/workflowstatus/` | Workflow state tracking and queries |

Services are stateless — they delegate persistence to the storage layer and emit events through the event bus.

## Storage Layer (`internal/storage/`)

Abstracts persistence behind Go interfaces. Supports multiple backends:

| Backend | Mode | Storage Engine |
|---------|------|---------------|
| **Local** | `AGENTFIELD_MODE=local` | SQLite (relational) + BoltDB (key-value) |
| **PostgreSQL** | `AGENTFIELD_STORAGE_MODE=postgresql` | PostgreSQL 15+ |
| **Cloud** | `AGENTFIELD_STORAGE_MODE=cloud` | Managed PostgreSQL |

**Key files:**
- `internal/storage/execution_records.go` — Execution record persistence
- `internal/storage/migrations.go` — Schema migration logic
- `internal/storage/observability_webhook.go` — Observability webhook storage
- `internal/infrastructure/storage/` — Connection pooling and DB setup

The unified interface means services never know which backend is active. Backend selection happens at startup via configuration.

**Code reference:** `control-plane/internal/storage/execution_records.go` — execution persistence, `control-plane/internal/storage/migrations.go` — migration logic

## Event Bus (`internal/events/`)

Publish-subscribe event system powering real-time updates via Server-Sent Events (SSE):

| File | Purpose |
|------|---------|
| `events/event_bus.go` | Core event bus — publish/subscribe infrastructure |
| `events/execution_events.go` | Workflow execution lifecycle events |
| `events/node_events.go` | Agent node registration and heartbeat events |
| `events/trigger_events.go` | External trigger events (webhook, cron, etc.) |
| `events/reasoner_events.go` | Reasoner/skill execution events |
| `events/publishers_test.go` | Publisher tests |

The event bus is the backbone of observability — workflow state changes, node status updates, and trigger activations all flow through it. The Web UI consumes events via SSE for real-time updates.

**Code reference:** `control-plane/internal/events/event_bus.go` — core event bus, `control-plane/internal/events/execution_events.go` — execution events

## MCP Integration

Model Context Protocol integration enables Silmari to interoperate with MCP-compatible tools and services. The integration layer provides MCP server/client capabilities for tool discovery and execution.

**Code reference:** `control-plane/internal/mcp/` — MCP protocol integration

## Configuration System (`internal/config/`)

Uses [Viper](https://github.com/spf13/viper) for configuration with this precedence:

1. Environment variables (highest)
2. Config file (`config/agentfield.yaml` by default, or another path via the legacy-compatible `AGENTFIELD_CONFIG_FILE`)
3. Code defaults (lowest)

**Reference:** `control-plane/.env.example` documents all supported variables.

Key configuration areas:
- Server port and mode (`AGENTFIELD_PORT`, `AGENTFIELD_MODE`)
- Storage backend selection (`AGENTFIELD_STORAGE_MODE`)
- Database connection (`AGENTFIELD_DATABASE_URL`)
- UI embedding (`AGENTFIELD_UI_ENABLED`, `AGENTFIELD_UI_MODE`)
- Logging (`LOG_LEVEL`, `GIN_MODE`)

## CLI System (`internal/cli/`)

The `af` CLI is built with a custom command framework:

- `internal/cli/commands/` — Command implementations (dev, init, run, server)
- `internal/cli/framework/` — CLI framework utilities

Commands are registered and dispatched through the framework. The `af dev` command starts a local development server with SQLite+BoltDB. The `af init` command scaffolds new agent projects using templates.

**Code reference:** `control-plane/internal/cli/commands/` — CLI commands, `control-plane/internal/templates/` — code generation templates

## Trigger Sources (`internal/sources/`)

External event sources that can activate agent workflows:

| Source | Path | Purpose |
|--------|------|---------|
| **GitHub** | `sources/github/` | GitHub webhook events |
| **Stripe** | `sources/stripe/` | Stripe webhook events |
| **Slack** | `sources/slack/` | Slack event integration |
| **Cron** | `sources/cron/` | Scheduled time-based triggers |
| **Generic Bearer** | `sources/genericbearer/` | Custom webhook with bearer auth |
| **Generic HMAC** | `sources/generichmac/` | Custom webhook with HMAC auth |
| **All** | `sources/all/` | Combined source registration |

**Code reference:** `control-plane/internal/sources/github/` — GitHub trigger source, `control-plane/internal/sources/cron/` — cron trigger source

## Encryption (`internal/encryption/`)

Cryptographic primitives for the DID/VC system:

- `internal/encryption/encryption.go` — Core cryptographic operations

Used by the DID/VC services to generate keys, sign credentials, and verify proofs.

**Code reference:** `control-plane/internal/encryption/encryption.go`

## Database Migrations

Goose-managed SQL migrations in `control-plane/migrations/`. Each migration has up and down SQL files. Migrations define the schema for PostgreSQL mode, including tables for agents, workflows, executions, memory, and audit trails.

**Code reference:** `control-plane/migrations/` — all migration files

## Observability

Structured logging via [zerolog](https://github.com/rs/zerolog):

- `internal/logger/` — Logger initialization and configuration
- `internal/observability/` — Observability utilities

Log levels: `debug`, `info`, `warn`, `error`. Configured via `LOG_LEVEL` environment variable.

## Embedded Assets (`internal/embedded/`)

The production build embeds the Web UI's compiled static assets (HTML, JS, CSS) into the Go binary using Go's `embed` package. In development mode, the UI is served by a separate Vite dev server with hot module replacement.
