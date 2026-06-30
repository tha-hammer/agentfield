# Go SDK

The Silmari Go SDK enables building agents in Go that connect to the Silmari control plane. Uses the standard library `net/http` and provides agent lifecycle management, skill registration, LLM harness integration, DID/VC identity, and cross-agent communication.

**Module path:** `github.com/Agent-Field/agentfield/sdk/go`

## Architecture

```
┌─────────────────────────────────────────┐
│              Agent (agent/)             │
│  ┌──────────┐ ┌──────────┐ ┌─────────┐ │
│  │Lifecycle │ │  Skills  │ │ Router  │ │
│  │(Init,Run,│ │(Register,│ │(Request │ │
│  │ Shutdown)│ │ Execute) │ │ routing)│ │
│  └──────────┘ └──────────┘ └─────────┘ │
│  ┌──────────┐ ┌──────────┐ ┌─────────┐ │
│  │  Memory  │ │   CLI    │ │  DID    │ │
│  │(CP-backed│ │(CLI mode │ │(Identity│ │
│  │  store)  │ │ support) │ │ + VCs)  │ │
│  └──────────┘ └──────────┘ └─────────┘ │
│  ┌──────────┐ ┌──────────────────────┐  │
│  │Verification│ Process Logs         │  │
│  │(Agent     │ (Execution logging)   │  │
│  │ verify)   │                       │  │
│  └──────────┘ └──────────────────────┘  │
├─────────────────────────────────────────┤
│         Client (client/)                │
│         HTTP client for CP API          │
├─────────────────────────────────────────┤
│         Types (types/)                  │
│         Shared type definitions         │
├─────────────────────────────────────────┤
│         Harness (harness/)              │
│         LLM provider integrations       │
│         Claude, Gemini, Codex, OpenAI   │
├─────────────────────────────────────────┤
│         DID (did/)                      │
│         Identity, VCs, crypto           │
└─────────────────────────────────────────┘
```

## Agent Module (`agent/`)

### Lifecycle (`agent/agent_lifecycle.go`)

The `Agent` struct manages the full agent lifecycle:

```go
// Initialize registers the agent with the control plane.
func (a *Agent) Initialize(ctx context.Context) error

// Run intelligently routes between CLI and server modes.
func (a *Agent) Run(ctx context.Context) error
```

**Initialization flow (see `agent/agent_lifecycle.go:18-76`):**
1. Register node with control plane
2. Auto-register DIDs if enabled (`EnableDID` or `VCEnabled`)
3. Mark agent as ready (`markReady`)
4. Start lease renewal loop (`startLeaseLoop`)

**Code reference:** `sdk/go/agent/agent_lifecycle.go:18-76` — Initialize, `sdk/go/agent/agent_lifecycle.go:79` — Run

### Agent Core (`agent/agent.go`)

Main `Agent` struct definition and configuration:

```go
type Agent struct {
    cfg         Config
    client      *client.Client
    reasoners   map[string]ReasonerFunc
    // ...
}
```

**Code reference:** `sdk/go/agent/agent.go`

### Skill Registration

Skills are Go functions registered with the agent:

```go
agent.RegisterSkill("greet", func(ctx context.Context, input map[string]any) (any, error) {
    return map[string]any{"message": "hello"}, nil
})
```

Skills become callable via the control plane at `<node_id>.greet`.

**Code reference:** `sdk/go/agent/agent.go` — RegisterSkill

### Router (`agent/router.go`)

Routes incoming execution requests to registered skills. Handles input validation, context propagation, and response formatting.

**Code reference:** `sdk/go/agent/router.go`

### Memory Backend (`agent/control_plane_memory_backend.go`)

Control plane-backed memory implementation. Supports the same four scopes as the Python SDK (global, session, actor, workflow). Memory operations are proxied to the control plane's memory API.

**Code reference:** `sdk/go/agent/control_plane_memory_backend.go`

### DID Integration (`agent/agent_did.go`)

Opt-in decentralized identity for Go agents. When `cfg.EnableDID` or `cfg.VCEnabled` is true, the agent automatically initializes DIDs during `Initialize()`.

**Code reference:** `sdk/go/agent/agent_did.go`

### CLI Mode (`agent/cli.go`)

Supports running agents in CLI mode for local development and testing. The `Run()` method intelligently routes between CLI mode (when `os.Args` has subcommands) and server mode (long-running HTTP listener).

**Code reference:** `sdk/go/agent/cli.go`

### Verification (`agent/verification.go`)

Agent verification logic — validates agent state, configuration, and connectivity to the control plane.

**Code reference:** `sdk/go/agent/verification.go`

### Process Logs (`agent/process_logs.go`)

Structured execution logging with ring buffer for in-memory log retention. Logs are annotated with execution context (workflow ID, node ID, timestamp).

**Code reference:** `sdk/go/agent/process_logs.go`

### Cancel (`agent/cancel.go`)

Execution cancellation support. Allows in-flight executions to be cancelled via the control plane.

**Code reference:** `sdk/go/agent/cancel.go`

## Client Module (`client/`)

HTTP client for the Silmari control plane API. Handles:
- Agent registration and heartbeat
- Cross-agent execution calls
- Memory operations
- Workflow status queries
- Authentication

**Code reference:** `sdk/go/client/`

## Types Module (`types/`)

Shared type definitions used across the SDK:
- `AgentStatus` — agent lifecycle states
- `AIConfig` — AI model configuration
- `DiscoveryResult` — agent discovery response
- `HarnessConfig` — LLM harness configuration
- `MemoryConfig` — memory backend configuration

**Code reference:** `sdk/go/types/`

## Harness Module (`harness/`)

LLM provider integration harness. Provides a unified interface for calling different AI models:

| Provider | File | Purpose |
|----------|------|---------|
| **Claude Code** | `harness/claudecode.go` | Anthropic Claude via Claude Code CLI |
| **Codex** | `harness/codex.go` | OpenAI Codex integration |
| **Gemini** | `harness/gemini.go` | Google Gemini integration |
| **OpenCode** | `harness/opencode.go` | OpenCode integration |
| **Factory** | `harness/factory.go` | Provider factory and selection |
| **Runner** | `harness/runner.go` | Test runner with harness integration |
| **Schema** | `harness/schema.go` | Structured output schema definitions |
| **Result** | `harness/result.go` | Standardized result types |
| **CLI** | `harness/cli.go` | CLI harness interface |
| **Provider** | `harness/provider.go` | Provider interface definition |

The harness is designed for both production use and testing — it enables deterministic testing of agent behavior across different LLM providers.

**Code reference:** `sdk/go/harness/runner.go` — test runner, `sdk/go/harness/factory.go` — provider factory

## DID Module (`did/`)

Decentralized identity implementation:

| File | Purpose |
|------|---------|
| `did/did_manager.go` | DID creation, resolution, key rotation |
| `did/did_client.go` | HTTP client for DID operations |
| `did/vc_generator.go` | Verifiable credential issuance and signing |
| `did/types.go` | DID/VC type definitions |

**Code reference:** `sdk/go/did/did_manager.go`, `sdk/go/did/vc_generator.go`

## Agent Lifecycle

```
┌──────────┐    ┌──────────┐    ┌──────────┐    ┌──────────┐
│  New()   │───▶│Initialize│───▶│  Run()   │───▶│ Shutdown │
│(Create   │    │(Register │    │(HTTP or  │    │(Cleanup, │
│ struct)  │    │ with CP) │    │ CLI mode)│    │ drain)   │
└──────────┘    └──────────┘    └──────────┘    └──────────┘
                     │
                     ▼
              ┌──────────┐
              │  Lease   │ (background goroutine)
              │  Loop    │ renews agent lease with CP
              └──────────┘
```

1. **New(config):** Create Agent struct with node ID, control plane URL, skills
2. **Initialize(ctx):** Register with control plane, set up DIDs, start lease loop
3. **Run(ctx):** Start HTTP server (server mode) or execute CLI command (CLI mode)
4. **Shutdown:** Graceful drain, deregistration, cleanup

**Code reference:** `sdk/go/agent/agent_lifecycle.go:18-76` — Initialize, `sdk/go/agent/agent_lifecycle.go:79` — Run

## Testing

Standard Go testing with table-driven tests:

```bash
cd sdk/go
go test ./...                              # All tests
go test ./agent/                           # Agent package tests
go test ./harness/                         # Harness tests
go test ./did/                             # DID tests
```

Tests include invariant tests (`*_invariant_test.go`) for property-based testing, integration tests for provider harnesses, and branch coverage tests.

**Code reference:** `sdk/go/agent/agent_test.go`, `sdk/go/harness/runner_test.go`
