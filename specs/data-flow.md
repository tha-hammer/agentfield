# Data Flow

How data moves through the AgentField system — agent-to-agent communication, workflow execution, memory synchronization, and event streaming.

## Agent-to-Agent Communication

Agents never call each other directly. All communication flows through the control plane:

```
┌──────────┐  1. Call("agent_b.score")  ┌──────────────┐
│ Agent A  │───────────────────────────▶│ Control Plane│
│ (caller) │                            │              │
│          │◀───────────────────────────│              │
└──────────┘  6. Return result          │  2. Lookup   │
                                        │     agent_b  │
┌──────────┐                            │  3. Route    │
│ Agent B  │◀───────────────────────────│     call     │
│ (target) │  4. Execute score()        │              │
│          │───────────────────────────▶│              │
└──────────┘  5. Return output          └──────────────┘
```

### Call Flow (Python)

```python
# Agent A calls Agent B's reasoner
result = await app.call("agent_b.score_claim", claim={"id": "123"})
```

1. `app.call()` → HTTP POST to control plane at `POST /api/v1/execute/agent_b.score_claim`
2. Control plane looks up `agent_b` in agent registry
3. Control plane routes request to Agent B's HTTP server
4. Agent B executes `score_claim` reasoner with input
5. Agent B returns output to control plane
6. Control plane returns output to Agent A

### Call Flow (Go)

```go
result, err := agent.Call(ctx, "agent_b.score_claim", map[string]any{"claim": claim})
```

Same flow — Go client sends to control plane, control plane routes to target agent.

**Code reference:** `control-plane/internal/handlers/agentic/` — execution routing, `sdk/python/agentfield/client.py` — Python client, `sdk/go/client/` — Go client

## Workflow Execution DAG

AgentField models multi-agent executions as Directed Acyclic Graphs (DAGs):

```
        ┌─────────┐
        │ Trigger │  (webhook, cron, manual)
        └────┬────┘
             │
        ┌────▼────┐
        │Agent A  │  score_claim()
        └────┬────┘
             │
     ┌───────┼───────┐
     │               │
┌────▼────┐    ┌─────▼────┐
│Agent B  │    │ Agent C  │  (parallel branch)
│fraud_   │    │risk_     │
│check()  │    │assess()  │
└────┬────┘    └─────┬────┘
     │               │
     └───────┬───────┘
             │
        ┌────▼────┐
        │Agent D  │  approve_claim()
        │(depends │
        │ on B,C) │
        └────┬────┘
             │
        ┌────▼────┐
        │  Done   │
        └─────────┘
```

### Execution Model

1. **Trigger:** External event (webhook, cron, manual) initiates workflow
2. **Root Execution:** Control plane dispatches root agent call
3. **Fan-out:** Agent calls spawn child executions. Parallel calls execute concurrently.
4. **Dependency Resolution:** Downstream agents execute when all upstream dependencies complete
5. **Completion:** Workflow completes when all leaf nodes finish (or fail/cancel)

### State Machine

```
PENDING ──▶ RUNNING ──▶ COMPLETED
                │
                ├──▶ FAILED
                │
                └──▶ CANCELLED
```

Intermediate states: `PENDING_APPROVAL` (human-in-the-loop gate).

**Code reference:** `control-plane/internal/services/` — workflow execution service, `control-plane/internal/events/execution_events.go` — execution lifecycle events

## Memory Synchronization

AgentField maintains four memory scopes synchronized through the control plane:

```
┌──────────────────────────────────────────┐
│           Control Plane                  │
│  ┌────────────────────────────────────┐  │
│  │         Memory Store               │  │
│  │  ┌────────┐ ┌────────┐ ┌────────┐  │  │
│  │  │ Global │ │Session │ │ Actor  │  │  │
│  │  │ Scope  │ │ Scope  │ │ Scope  │  │  │
│  │  └────────┘ └────────┘ └────────┘  │  │
│  │  ┌────────┐                        │  │
│  │  │Workflow│ (per-execution)        │  │
│  │  │ Scope  │                        │  │
│  │  └────────┘                        │  │
│  └────────────────────────────────────┘  │
│  ┌────────────────────────────────────┐  │
│  │       Memory Event Bus             │  │
│  │   (change notifications via SSE)   │  │
│  └────────────────────────────────────┘  │
└──────────┬───────────────┬───────────────┘
           │               │
    ┌──────▼──────┐ ┌──────▼──────┐
    │  Agent A    │ │  Agent B    │
    │ memory.get()│ │ memory.set()│
    │ memory.set()│ │ memory.get()│
    └─────────────┘ └─────────────┘
```

### Read Path

1. Agent calls `memory.get(key)`
2. SDK sends request to control plane memory API
3. Control plane resolves from narrowest to widest scope: `workflow → session → actor → global`
4. Returns first match or null

### Write Path

1. Agent calls `memory.set(scope, key, value)`
2. SDK sends request to control plane memory API
3. Control plane persists to storage backend
4. Control plane emits memory change event on event bus
5. Other agents subscribed to memory events receive notification via SSE

### Memory Events

Agents can react to memory changes — enabling event-driven agent behavior:

```python
# Agent triggers when memory changes
@agent.on_memory_change("global", "config")
async def on_config_change(key, value, old_value):
    await agent.update_configuration(value)
```

**Code reference:** `control-plane/internal/handlers/memory_events.go` — memory event handler, `sdk/python/agentfield/memory.py` — Python memory client, `sdk/python/agentfield/memory_events.py` — memory event client

## Event Streaming (SSE)

The control plane pushes real-time updates to the Web UI and agents via Server-Sent Events:

```
┌──────────┐  SSE stream    ┌──────────────┐
│ Web UI   │◀───────────────│ Control Plane│
│          │  workflow       │              │
│          │  status, node   │  Event Bus   │
│          │  health, exec   │  (publish/   │
│          │  state changes  │  subscribe)  │
└──────────┘                 └──────────────┘
```

### Event Types

| Event | Source File | Payload |
|-------|-------------|---------|
| Execution started | `events/execution_events.go` | workflow_id, node_id, input |
| Execution completed | `events/execution_events.go` | workflow_id, node_id, output, duration |
| Execution failed | `events/execution_events.go` | workflow_id, node_id, error |
| Node registered | `events/node_events.go` | node_id, name, port |
| Node heartbeat | `events/node_events.go` | node_id, status, timestamp |
| Trigger fired | `events/trigger_events.go` | trigger_id, source, payload |
| Reasoner executed | `events/reasoner_events.go` | node_id, reasoner, input, output |

**Code reference:** `control-plane/internal/events/event_bus.go` — event bus core, `control-plane/internal/events/execution_events.go` — execution events

## Observability Data Flow

```
┌──────────┐     ┌──────────────┐     ┌──────────────┐
│  Agents  │────▶│ Control Plane│────▶│ Storage      │
│ (logs,   │     │ (event bus,  │     │ (execution   │
│  traces) │     │  metrics)    │     │  records,    │
└──────────┘     └──────┬───────┘     │  audit logs) │
                        │             └──────────────┘
                        │ SSE
                 ┌──────▼───────┐
                 │   Web UI     │
                 │ (real-time   │
                 │  dashboard)  │
                 └──────────────┘
```

1. Agents emit logs and execution metadata to control plane
2. Control plane publishes events on event bus
3. Events are persisted to storage for historical queries
4. Web UI consumes events via SSE for real-time updates
5. Audit trails are generated from stored execution records + cryptographic proofs

**Code reference:** `control-plane/internal/storage/execution_records.go` — execution persistence, `control-plane/internal/storage/observability_webhook.go` — observability webhooks
