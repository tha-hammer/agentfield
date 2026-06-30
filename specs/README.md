# Silmari — Specifications & Architecture

Comprehensive architecture documentation and technical specifications for the Silmari platform. These docs map the codebase structure, design decisions, data flows, and operational characteristics of the three-tier monorepo.

## Document Index

| Document | Description |
|----------|-------------|
| [Architecture Overview](architecture-overview.md) | Three-tier system architecture, component diagram, design philosophy |
| [Control Plane](control-plane.md) | Go orchestration server — handlers, services, storage, events, MCP, routing |
| [Python SDK](sdk-python.md) | Python agent SDK — Agent lifecycle, reasoners, memory, AI integration |
| [Go SDK](sdk-go.md) | Go agent SDK — Agent lifecycle, skills, client, types, AI harness |
| [Web UI](web-ui.md) | React/TypeScript admin interface — components, routing, state, embedded build |
| [Data Flow](data-flow.md) | Agent-to-agent communication, workflow DAG execution, memory synchronization |
| [Security](security.md) | DID/VC identity, cryptographic audit trails, access control, IAM |
| [Deployment](deployment.md) | Local mode, PostgreSQL/cloud mode, Docker, Kubernetes |

## Quick Reference

- **Control plane entry:** `control-plane/cmd/af/` (unified CLI) and `control-plane/cmd/agentfield-server/` (standalone server)
- **Python SDK root:** `sdk/python/agentfield/agent.py` — `Agent` class
- **Go SDK root:** `sdk/go/agent/agent_lifecycle.go` — `Agent.Initialize()` and `Agent.Run()`
- **Web UI root:** `control-plane/web/client/src/` — React + Vite + Tailwind
- **Config:** `control-plane/.env.example` — all environment variables
- **Migrations:** `control-plane/migrations/` — Goose-managed SQL migrations

## Architecture at a Glance

```
┌─────────────────────────────────────────────────────────┐
│                      Web UI (React)                      │
│              control-plane/web/client/src/                │
├─────────────────────────────────────────────────────────┤
│                  Control Plane (Go/Gin)                   │
│              control-plane/cmd/agentfield/                │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐               │
│  │ Handlers │  │ Services │  │ Storage  │               │
│  │ (REST)   │  │ (Logic)  │  │ (SQLite/ │               │
│  │          │  │          │  │  PG/Bolt)│               │
│  └──────────┘  └──────────┘  └──────────┘               │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐               │
│  │  Events  │  │   MCP    │  │   DID    │               │
│  │  (SSE)   │  │ (Plugin) │  │  (Crypto)│               │
│  └──────────┘  └──────────┘  └──────────┘               │
├─────────────────────────────────────────────────────────┤
│              Agent SDKs (Python / Go)                     │
│  ┌──────────────────┐  ┌──────────────────┐              │
│  │   Python SDK     │  │     Go SDK       │              │
│  │ agent.py:Agent   │  │ agent_lifecycle  │              │
│  │ memory.py:Memory │  │ agent.go:Agent   │              │
│  │ agent_ai.py:AI   │  │ harness/:LLM     │              │
│  └──────────────────┘  └──────────────────┘              │
└─────────────────────────────────────────────────────────┘
```

## Key Design Principles

1. **Agents as backend services** — Each agent is an independently deployable process that registers with the control plane, owns a capability, and can scale independently.
2. **Control plane as orchestrator** — The control plane routes calls, tracks workflow DAGs, injects observability, and manages identity — agents never call each other directly.
3. **Cryptographic identity by default** — DID/VC-based identity and audit trails are built into the platform, not bolted on.
4. **Storage abstraction** — Unified storage interface supports SQLite+BoltDB (local dev) and PostgreSQL (production) with the same API.
5. **SDK parity** — Python and Go SDKs provide equivalent capabilities with idiomatic APIs for each language.

## Conventions

- Go code references use module path `github.com/Agent-Field/agentfield/control-plane`
- Python code references use import path `agentfield.*`
- File paths are relative to repository root unless absolute
- Code references use format: `path/to/file.go:line_number`
