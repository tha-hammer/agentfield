# Python SDK

The AgentField Python SDK enables building AI agents that connect to the AgentField control plane. Built on **FastAPI/Uvicorn**, it provides agent lifecycle management, reasoner registration, cross-agent memory, AI model integration, and cryptographic identity.

**Package root:** `sdk/python/agentfield/`

## Architecture

```
┌─────────────────────────────────────────┐
│              Agent (agent.py)           │
│  ┌─────────┐ ┌──────────┐ ┌──────────┐ │
│  │ Lifecycle│ │ Reasoner │ │ Execution│ │
│  │ (init,   │ │ Registry │ │ Context  │ │
│  │  run,    │ │          │ │          │ │
│  │  shutdown│ │          │ │          │ │
│  └─────────┘ └──────────┘ └──────────┘ │
│  ┌─────────┐ ┌──────────┐ ┌──────────┐ │
│  │  Memory │ │    AI    │ │   DID    │ │
│  │ (4 scope│ │ (LLM     │ │ (Identity│ │
│  │  levels)│ │  calls)  │ │  + VCs)  │ │
│  └─────────┘ └──────────┘ └──────────┘ │
│  ┌─────────┐ ┌──────────┐ ┌──────────┐ │
│  │  Client │ │  Router  │ │ Triggers │ │
│  │ (HTTP to│ │ (Agent   │ │ (External│ │
│  │   CP)   │ │  routing)│ │  events) │ │
│  └─────────┘ └──────────┘ └──────────┘ │
├─────────────────────────────────────────┤
│         AgentServer (FastAPI)           │
│         agent_server.py                 │
├─────────────────────────────────────────┤
│         AgentFieldHandler               │
│         agent_field_handler.py          │
└─────────────────────────────────────────┘
```

## Core Modules

### Agent (`agent.py`)

The `Agent` class is the main entry point for building agents. It manages:

- **Lifecycle:** Initialization, startup, shutdown, lease renewal
- **Reasoner Registration:** Decorator-based registration of agent functions
- **Execution:** Incoming request handling, routing to reasoners
- **Configuration:** Agent identity, control plane URL, DID/VC settings

**Key imports** (see `agent.py:1-66`):

```python
from agentfield.agent_ai import AgentAI
from agentfield.agent_field_handler import AgentFieldHandler
from agentfield.agent_server import AgentServer
from agentfield.agent_workflow import AgentWorkflow
from agentfield.client import AgentFieldClient
from agentfield.execution_context import ExecutionContext
from agentfield.did_manager import DIDManager
from agentfield.vc_generator import VCGenerator
from agentfield.memory import MemoryClient
from agentfield.router import AgentRouter
```

**Code reference:** `sdk/python/agentfield/agent.py:1-66` — Agent class imports and structure

### Agent Server (`agent_server.py`)

FastAPI application setup and lifecycle management. Creates the HTTP server that exposes agent reasoners as REST endpoints. Handles middleware, error formatting, and request validation.

**Code reference:** `sdk/python/agentfield/agent_server.py`

### Agent AI (`agent_ai.py`)

AI model integration layer. Provides:
- Structured model calls with response validation
- Tool/function calling support
- Model configuration (provider, model ID, parameters)
- Cost tracking integration

**Code reference:** `sdk/python/agentfield/agent_ai.py`

### Memory (`memory.py`)

Cross-agent persistent memory with four scope levels:

| Scope | Lifetime | Use Case |
|-------|----------|----------|
| **Global** | Until explicitly deleted | Shared configuration, knowledge bases |
| **Session** | Duration of user session | Conversation context, session preferences |
| **Actor** | Across all sessions for an actor | Actor-specific learned data |
| **Workflow (Run)** | Single workflow execution | Intermediate results, execution state |

Lookup resolves from narrowest to widest scope: `workflow → session → actor → global`.

**MemoryClient** (`memory.py`): Primary memory interface for get/set/delete operations across scopes.
**MemoryEventClient** (`memory_events.py`): Event-driven memory change notifications — agents can react to memory changes.

**Code reference:** `sdk/python/agentfield/memory.py:1-60` — Memory scope hierarchy and lookup behavior

### Execution Context (`execution_context.py`)

Per-execution state management. Provides:
- Current execution tracking (workflow ID, agent ID, session ID)
- Context propagation across agent calls
- Thread-safe context access via `get_current_context()`

**Code reference:** `sdk/python/agentfield/execution_context.py`

### Client (`client.py`)

HTTP client for communicating with the AgentField control plane:
- Agent registration
- Cross-agent calls (`app.call("other_agent.reasoner", ...)`)
- Memory operations
- Workflow status queries
- Approval handling

**Code reference:** `sdk/python/agentfield/client.py`

### Router (`router.py`)

Agent discovery and routing. Resolves agent-to-agent call targets and manages the agent network topology.

**Code reference:** `sdk/python/agentfield/router.py`

### DID/VC System

Decentralized identity and verifiable credentials:

| Module | Purpose |
|--------|---------|
| `did_manager.py` | DID creation, resolution, key management |
| `vc_generator.py` | VC issuance, signing, verification |
| `did_auth.py` | DID-based authentication |
| `agent_vc.py` | Agent-specific VC integration |

Enable opt-in per agent: `app.vc_generator.set_enabled(True)`.

**Code reference:** `sdk/python/agentfield/did_manager.py`, `sdk/python/agentfield/vc_generator.py`

### Agent Workflow (`agent_workflow.py`)

Multi-step workflow orchestration within an agent. Supports DAG-based execution with dependencies, parallel branches, and error handling.

**Code reference:** `sdk/python/agentfield/agent_workflow.py`

### Multimodal Support

Vision, image, and audio processing for multimodal agents:

| Module | Purpose |
|--------|---------|
| `multimodal.py` | Multimodal input processing |
| `multimodal_response.py` | Multimodal output formatting |
| `vision.py` | Vision/image analysis |
| `media_providers.py` | Media provider integrations |
| `media_router.py` | Media routing logic |

**Code reference:** `sdk/python/agentfield/multimodal.py`, `sdk/python/agentfield/multimodal_response.py`

### Tool Calling (`tool_calling.py`)

Structured tool/function calling support. Enables agents to invoke external tools with typed inputs and validated outputs.

**Code reference:** `sdk/python/agentfield/tool_calling.py`

### Triggers (`triggers.py`)

External event trigger support. Agents can react to webhooks, scheduled events, and system events.

**Code reference:** `sdk/python/agentfield/triggers.py`

### Additional Modules

| Module | Purpose |
|--------|---------|
| `agent_cli.py` | CLI mode support for agents |
| `agent_registry.py` | Agent registration with control plane |
| `agent_schema.py` | Pydantic schemas for agent data |
| `agent_pause.py` | Agent pause/resume for human-in-the-loop |
| `agent_serverless.py` | Serverless deployment support |
| `agent_discovery.py` | Agent discovery protocol |
| `agent_utils.py` | Utility functions |
| `async_config.py` | Async execution configuration |
| `async_execution_manager.py` | Background task execution |
| `cancel.py` | Execution cancellation |
| `connection_manager.py` | Connection pooling and management |
| `cost_tracker.py` | LLM cost tracking |
| `decorators.py` | Reasoner registration decorators |
| `exceptions.py` | SDK-specific exceptions |
| `execution_state.py` | Execution state machine |
| `http_connection_manager.py` | HTTP connection management |
| `litellm_adapters.py` | LiteLLM provider adapters |
| `logger.py` | Agent logging (debug, info, warn, error) |
| `node_logs.py` | Node-level log management |
| `pydantic_utils.py` | Pydantic validation utilities |
| `rate_limiter.py` | Request rate limiting |
| `result_cache.py` | Execution result caching |
| `status.py` | Agent status reporting |
| `testing.py` | Test utilities and fixtures |
| `types.py` | Shared type definitions (AgentStatus, AIConfig, etc.) |
| `utils.py` | General utilities |
| `verification.py` | Agent verification logic |

## Agent Lifecycle

```
┌──────────┐    ┌──────────┐    ┌──────────┐    ┌──────────┐
│  Create  │───▶│ Register │───▶│   Run    │───▶│ Shutdown │
│  Agent() │    │ with CP  │    │ (FastAPI │    │ (cleanup)│
│          │    │          │    │  server) │    │          │
└──────────┘    └──────────┘    └──────────┘    └──────────┘
```

1. **Create:** Instantiate `Agent(config)` with node ID, control plane URL, and optional DID/VC settings
2. **Register:** Agent registers its reasoners with the control plane, which adds them to the service catalog
3. **Run:** FastAPI server starts, exposing reasoners as HTTP endpoints. Control plane routes incoming calls.
4. **Shutdown:** Graceful shutdown, deregistration from control plane, cleanup

## Reasoner Registration

Reasoners are the core abstraction — they're Python functions decorated to become callable agent capabilities:

```python
@agent.reasoner()
async def score_claim(claim: dict) -> dict:
    """Score an insurance claim for risk."""
    # AI logic here
    return {"score": 0.85, "risk_level": "low"}
```

This registers `score_claim` as an endpoint accessible at `POST /api/v1/execute/<node_id>.score_claim`.

**Code reference:** `sdk/python/agentfield/decorators.py` — reasoner decorator implementation

## Memory API

```python
# Store in global scope (shared across all agents/sessions)
await agent.memory.set("global", "api_key", "sk-...")

# Store in session scope (conversation context)
await agent.memory.set("session", "user_prefs", {"theme": "dark"})

# Retrieve (narrowest-scope-first resolution)
val = await agent.memory.get("user_prefs")

# Delete
await agent.memory.delete("session", "user_prefs")
```

**Code reference:** `sdk/python/agentfield/memory.py` — MemoryClient class

## Testing

Tests are organized with pytest markers:

```bash
pytest                                    # Unit + functional (default)
pytest -m unit                            # Unit tests only
pytest -m integration                     # Integration tests only
pytest -m mcp                             # MCP tests (excluded by default)
pytest --cov=agentfield --cov-report=term-missing  # With coverage
```

Test markers defined in `pyproject.toml`: `unit`, `functional`, `integration`, `mcp`.

**Code reference:** `sdk/python/pyproject.toml` — test configuration
