# Architecture Overview

## System Design

Silmari is a **Kubernetes-style control plane for AI agents** — a three-tier monorepo providing production infrastructure for deploying, orchestrating, and observing multi-agent systems.

### Three-Tier Architecture

```
                    ┌──────────────────────────┐
                    │        Web UI            │
                    │   React / TypeScript     │
                    │   (Embedded SPA)         │
                    └──────────┬───────────────┘
                               │ HTTP/SSE
                    ┌──────────▼───────────────┐
                    │     Control Plane        │
                    │        Go / Gin          │
                    │                          │
                    │  ┌────────────────────┐  │
                    │  │   REST/gRPC API    │  │
                    │  │ handlers/          │  │
                    │  └────────┬───────────┘  │
                    │  ┌────────▼───────────┐  │
                    │  │   Business Logic   │  │
                    │  │   services/        │  │
                    │  └────────┬───────────┘  │
                    │  ┌────────▼───────────┐  │
                    │  │   Storage Layer    │  │
                    │  │   storage/         │  │
                    │  │  SQLite | PG | Bolt│  │
                    │  └────────────────────┘  │
                    │  ┌────────────────────┐  │
                    │  │   Event Bus (SSE)  │  │
                    │  │   events/          │  │
                    │  └────────────────────┘  │
                    └──────────┬───────────────┘
                               │ HTTP
              ┌────────────────┼────────────────┐
              │                                 │
    ┌─────────▼────────┐            ┌───────────▼────────┐
    │   Python Agent   │            │     Go Agent       │
    │   SDK (FastAPI)  │            │   SDK (net/http)   │
    │                  │            │                    │
    │ agent.py:Agent   │            │ agent/:Agent       │
    │ memory.py        │            │ harness/:LLM       │
    │ agent_ai.py      │            │ did/:Identity      │
    │ execution_ctx.py │            │ client/:HTTP       │
    └──────────────────┘            └────────────────────┘
```

## Component Inventory

### Control Plane (`control-plane/`)

The central orchestration server. Written in Go using the Gin web framework.

| Component | Path | Purpose |
|-----------|------|---------|
| **Entry Points** | `cmd/af/`, `cmd/agentfield-server/` | CLI and server binaries |
| **Server** | `internal/server/` | HTTP server setup, routing, middleware |
| **Handlers** | `internal/handlers/` | REST/gRPC request handlers (admin, agentic, connector, UI) |
| **Services** | `internal/services/` | Business logic — workflow execution, agent registry, DID/VC |
| **Storage** | `internal/storage/` | Data persistence — SQLite, PostgreSQL, BoltDB backends |
| **Events** | `internal/events/` | Event bus for workflow notifications and SSE streaming |
| **Core** | `internal/core/` | Domain models (`core/domain/models.go`) and interfaces |
| **Config** | `internal/config/` | Viper-based configuration management |
| **Encryption** | `internal/encryption/` | Cryptographic primitives for DID/VC |
| **MCP** | Not in internal/ — separate MCP integration | Model Context Protocol support |
| **CLI** | `internal/cli/` | CLI command definitions and routing |
| **Logger** | `internal/logger/` | Structured logging via zerolog |
| **Templates** | `internal/templates/` | Code generation templates for `af init` |
| **Embedded** | `internal/embedded/` | Embedded assets (web UI dist) |
| **Infrastructure** | `internal/infrastructure/` | DB connection pooling, process management, communication |
| **Sources** | `internal/sources/` | Trigger sources — GitHub, Stripe, Slack, cron, webhooks |
| **SkillKit** | `internal/skillkit/` | Skill data and definitions |
| **Packages** | `internal/packages/` | Shared internal packages |

### Python SDK (`sdk/python/agentfield/`)

Agent builder for Python. Built on FastAPI/Uvicorn.

| Module | Path | Purpose |
|--------|------|---------|
| **Agent** | `agent.py` | Main `Agent` class — lifecycle, reasoner registration, execution |
| **Agent Server** | `agent_server.py` | FastAPI server setup and lifecycle |
| **Agent AI** | `agent_ai.py` | AI model integration, structured outputs |
| **Memory** | `memory.py` | Cross-agent persistent memory — 4 scopes (global, session, actor, workflow) |
| **Memory Events** | `memory_events.py` | Event-driven memory change notifications |
| **Execution Context** | `execution_context.py` | Per-execution state, current context tracking |
| **Client** | `client.py` | HTTP client for calling the control plane |
| **Router** | `router.py` | Agent routing and discovery |
| **DID Manager** | `did_manager.py` | Decentralized identity management |
| **VC Generator** | `vc_generator.py` | Verifiable credential generation |
| **Agent Workflow** | `agent_workflow.py` | Multi-step workflow orchestration |
| **Multimodal** | `multimodal.py`, `multimodal_response.py` | Vision, image, audio processing |
| **Tool Calling** | `tool_calling.py` | Structured tool/function calling |
| **Triggers** | `triggers.py` | External event triggers |
| **Rate Limiter** | `rate_limiter.py` | Request rate limiting |
| **Cost Tracker** | `cost_tracker.py` | LLM cost tracking |
| **Logger** | `logger.py` | Agent logging |
| **Types** | `types.py` | Shared type definitions |
| **Harness** | `harness/` | LLM harness for testing and evaluation |

### Go SDK (`sdk/go/`)

Agent builder for Go. Uses standard library `net/http`.

| Module | Path | Purpose |
|--------|------|---------|
| **Agent** | `agent/agent_lifecycle.go` | Agent lifecycle — Initialize, Run, register, lease |
| **Agent Core** | `agent/agent.go` | Main Agent struct and configuration |
| **Router** | `agent/router.go` | Request routing to skills |
| **Client** | `client/` | HTTP client for control plane API |
| **Types** | `types/` | Shared type definitions |
| **DID** | `did/` | Decentralized identity — client, manager, VC generator |
| **Harness** | `harness/` | LLM provider harness — Claude, Gemini, Codex, OpenAI |
| **Memory Backend** | `agent/control_plane_memory_backend.go` | Control plane memory integration |
| **CLI** | `agent/cli.go` | CLI mode support |
| **Verification** | `agent/verification.go` | Agent verification logic |
| **Process Logs** | `agent/process_logs.go` | Execution logging |

### Web UI (`control-plane/web/client/`)

React + TypeScript admin interface. Embedded in Go binary via `embed`.

| Area | Path | Purpose |
|------|------|---------|
| **Pages** | `src/pages/` | Route-level page components |
| **Components** | `src/components/` | Reusable UI components by domain |
| **Dashboard** | `src/components/dashboard/` | Main dashboard widgets |
| **Workflow DAG** | `src/components/WorkflowDAG/` | Workflow graph visualization |
| **Execution** | `src/components/execution/` | Execution monitoring UI |
| **Nodes** | `src/components/nodes/` | Agent node management |
| **Triggers** | `src/components/triggers/` | Trigger configuration |
| **DID/VC** | `src/components/did/`, `src/components/vc/` | Identity and credentials UI |
| **Contexts** | `src/contexts/` | React context providers |
| **Hooks** | `src/hooks/` | Custom React hooks and queries |
| **Services** | `src/services/` | API client services |
| **Types** | `src/types/` | TypeScript type definitions |
| **Config** | `src/config/` | Frontend configuration |

## Storage Architecture

Silmari supports three storage backends behind a unified interface:

```
┌──────────────────────────────────────┐
│         Storage Interface            │
│    internal/storage/ (Go interfaces) │
├────────────┬────────────┬────────────┤
│  Local     │ PostgreSQL │   Cloud    │
│  SQLite +  │  Full RDBMS│  Managed   │
│  BoltDB    │            │  Postgres  │
├────────────┴────────────┴────────────┤
│     internal/infrastructure/storage/ │
│     Connection pooling, migrations   │
└──────────────────────────────────────┘
```

- **Local mode** (`AGENTFIELD_MODE=local`): SQLite for relational data + BoltDB for key-value storage. No external dependencies. Used for development.
- **PostgreSQL mode** (`AGENTFIELD_STORAGE_MODE=postgresql`): Full PostgreSQL backend. Requires running Goose migrations from `control-plane/migrations/`. Production-ready.
- **Cloud mode** (`AGENTFIELD_STORAGE_MODE=cloud`): Managed PostgreSQL. Used in distributed deployments.

The storage interface is defined by Go interfaces in `internal/storage/`. Services call these interfaces; the storage layer switches backends based on configuration. See `control-plane/internal/infrastructure/storage/` for connection pooling logic.

## Configuration System

Configuration uses [Viper](https://github.com/spf13/viper) with this precedence:

1. **Environment variables** (highest priority) — e.g., `AGENTFIELD_PORT`, `AGENTFIELD_MODE`
2. **Config file** — `config/agentfield.yaml` by default, or the path set via the legacy-compatible `AGENTFIELD_CONFIG_FILE`
3. **Defaults** — defined in `internal/config/`

Key configuration entry point: `control-plane/internal/config/` (Viper initialization).

Reference: `control-plane/.env.example` for all supported variables.

## Database Migrations

Migrations are managed by [Goose](https://github.com/pressly/goose):

- **Directory:** `control-plane/migrations/`
- **Format:** SQL files with up/down migrations
- **Execution:** `goose -dir ./migrations postgres "$AGENTFIELD_DATABASE_URL" up`

Always run migrations before starting the server in PostgreSQL mode. Local mode handles schema automatically via SQLite.
