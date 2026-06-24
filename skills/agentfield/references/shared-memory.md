# Shared Memory — Distributed State Without the Redis Tax

Reasoners are stateless functions. **Shared memory is how they share context** — a distributed key-value (and vector) store with four isolation scopes, managed by the control plane. No Redis to stand up, no pub/sub to wire, no cache-invalidation code. One reasoner's write becomes another reasoner's context.

Pick shared memory when work spans more than one reasoner call and the later calls need what the earlier calls learned: a sentiment read taken at intake that the resolution reasoner needs an hour later; a document corpus embedded once and queried by many specialists; organization-wide config every agent reads.

This guide is the **local, offline** reference. It does not require a network fetch. For *reacting* to writes (not just reading them), see `references/memory-events.md`.

---

## The four scopes

Scope is the blast radius of a value — who can see it and how long it lives. Verified from the SDK; pick the **narrowest** scope that still reaches every reader.

| Scope | Visible to | Lives until | Use for |
|---|---|---|---|
| `global` | every agent, every session | explicitly deleted | config, shared knowledge bases, cross-agent state |
| `session` | one user session / conversation | session ends | conversation context, in-session preferences |
| `actor` | one actor, across all its sessions | explicitly deleted | actor-specific learned data, persona config |
| `workflow` | one workflow execution | the run completes | intermediate results, per-run scratch state |

`actor` is sometimes written `agent`; `workflow` is sometimes written `run` — same scopes. Conceptually they nest widest→narrowest: **global ⊃ session ⊃ actor ⊃ workflow**.

### Hierarchical read resolution

`memory.get(key)` **without an explicit scope** resolves most-specific → least-specific and returns the first hit:

```
workflow  →  session  →  actor  →  global
```

So a narrower scope **overrides** a broader one for reads. Write a default in `global`, override it per-session, and `get` transparently returns the most specific value in context. This is the design lever: put org defaults in `global`, let `session`/`workflow` shadow them, and readers never branch on scope.

---

## Key-value API

```python
# Explicit scope via the scope objects
await app.memory.global_scope.set("config", {"temperature": 0.2})
await app.memory.session(session_id).set("context", {"topic": "billing"})
await app.memory.actor(actor_id).set("preferences", {"tone": "concise"})
await app.memory.workflow(workflow_id).set("step1_output", {"ok": True})

# Hierarchical lookup from current context (workflow→session→actor→global)
value = await app.memory.get("preferences", default={})

# Low-level client with explicit scope + scope_id
await memory_client.set("context", {"topic": "billing"}, scope="session", scope_id=session_id)
```

Core operations on every scope: **`set(key, data)`**, **`get(key, default=None)`**, **`delete(key)`**. Values are JSON-serializable (dicts, lists, scalars).

```python
await app.memory.global_scope.set("ticket:T-123.sentiment", {"mood": "angry", "urgency": "high"})
sentiment = await app.memory.global_scope.get("ticket:T-123.sentiment")
```

---

## Vector API

Each scope also carries vector memory for semantic retrieval — embed once, search by similarity. Same scope semantics apply (a `workflow`-scoped vector dies with the run; a `global` one persists).

```python
# Store an embedding with metadata for filtering
await app.memory.global_scope.set_vector(
    "doc:chunk-1",
    embedding,                                   # Sequence[float]
    metadata={"source": "contracts.pdf"},
)

# Top-k similarity search, optionally filtered on metadata
hits = await app.memory.global_scope.similarity_search(
    query_embedding,
    top_k=5,
    filters={"source": "contracts.pdf"},
)

await app.memory.global_scope.delete_vector("doc:chunk-1")
```

`similarity_search` returns a list of hit dicts. `top_k` defaults to 10. `filters` narrows the candidate set by the metadata you stored with `set_vector`. This is the canonical RAG substrate: chunk + embed into a scope once, fan out specialists that each `similarity_search` the slice they need.

---

## Architectural rules

**Choose the narrowest scope that reaches every reader.** `workflow` for per-run scratch that should vanish when the run ends. `global` only for truly cross-agent, long-lived state. Over-scoping to `global` turns transient scratch into a leak you must remember to `delete`; under-scoping hides a value the next reasoner needed.

**Keys are a namespace you design, not random strings.** `ticket:<id>.sentiment`, `warehouse.stock.<sku>`, `risk.scores.<account>`. Good keys make `references/memory-events.md` subscriptions (`warehouse.stock.*`) fall out for free. Design keys and event patterns together.

**Memory is shared context, not a message bus.** Use `set`/`get` to share *state*; use `app.call(...)` to request *work*; use `on_change` to *react* to a write. Don't smuggle RPC through memory (write a "request" key and poll for a "response" key) — that's a queue you didn't need.

**Reconstruct on read, don't trust shape.** A value round-trips as JSON. A Pydantic model written in goes out as a plain dict. Reconstruct (`Model(**value)`) on read if you need methods, exactly like the `app.call` cross-boundary gotcha.

**Idempotent writes for replayed work.** Triggered and reactive reasoners can re-run. Scope run-local writes by the run id (`app.memory.workflow(run_id)`) so a replay overwrites cleanly instead of accreting.

---

## When NOT to use shared memory

- The value is needed once, synchronously, by the immediate caller — pass it as the `app.call(...)` return value. Memory adds a round trip and a lifetime you must manage.
- It's secret material (API keys, webhook secrets) — those belong in env (`secret_env=...`), never in memory.
- You're reaching for memory to implement request/response between two reasoners — use `app.call`. Memory is state, not transport.

---

## Concrete example — shared context across a multi-agent ticket flow

Intake writes sentiment once; a downstream resolution agent reads it without re-deriving it; a notifier reacts to the same write. One write, three consumers, zero direct coupling.

```python
# Intake agent — writes shared context
@app.reasoner(tags=["entry"])
async def intake(ticket_id: str, body: str) -> dict:
    sentiment = await app.call(f"{app.node_id}.read_sentiment", body=body)
    await app.memory.global_scope.set(f"ticket:{ticket_id}.sentiment", sentiment)
    return {"ticket_id": ticket_id, "sentiment": sentiment}


# Resolution agent — reads the context instead of recomputing it
@app.reasoner()
async def resolve(ticket_id: str) -> dict:
    sentiment = await app.memory.global_scope.get(f"ticket:{ticket_id}.sentiment", default={})
    urgency = sentiment.get("urgency", "normal")
    plan = await app.call(f"{app.node_id}.draft_resolution", ticket_id=ticket_id, urgency=urgency)
    return plan


# Notifier — reacts to the very same write (see memory-events.md)
@app.on_change("ticket:*.sentiment")
async def on_sentiment(event):
    if (event.data or {}).get("urgency") == "high":
        await app.call(f"{app.node_id}.page_oncall", key=event.key, change=event.data)
```

`resolve` never calls `intake` — it reads the state `intake` left behind. `on_sentiment` never imports either — it reacts to the write. That's composite intelligence held together by a key namespace.

---

## Shared-memory anti-patterns

| ❌ | ✅ |
|---|---|
| Everything in `global` scope | Narrowest scope that reaches every reader; `workflow` for per-run scratch |
| Random flat keys (`sentiment_T123`) | Namespaced keys (`ticket:<id>.sentiment`) that also make `*` subscriptions work |
| Write a "request" key, poll for a "response" key | `app.call(...)` — memory is state, not a message queue |
| Storing secrets in memory | `secret_env=...` and real env vars |
| Passing a Pydantic model through memory and expecting methods back | Reconstruct `Model(**value)` on read |
| `global` scratch never cleaned up | `workflow` scope auto-clears on run completion |
| Re-deriving a value a sibling already wrote | `get` the shared key; compute once, read many |
| Branching reader code on scope | Hierarchical `get` (workflow→session→actor→global) resolves it for you |

---

## Smoke test for a shared-state build

After `docker compose up`:

1. Write then read across reasoners:
   ```bash
   # intake writes ticket:T-1.sentiment
   curl -sS -X POST http://localhost:8080/api/v1/execute/<slug>.intake \
     -H 'Content-Type: application/json' \
     -d '{"input":{"ticket_id":"T-1","body":"this is the third time..."}}'

   # resolve reads it back — should reflect intake's write, not a recompute
   curl -sS -X POST http://localhost:8080/api/v1/execute/<slug>.resolve \
     -H 'Content-Type: application/json' \
     -d '{"input":{"ticket_id":"T-1"}}' | jq '.result'
   ```
2. Confirm scope lifetime: a `workflow`-scoped key written in one run must **not** be readable from a different run id. A `global` key must survive across both.

If `resolve` returns urgency derived from `intake`'s write — across two separate executions — shared context is healthy.

---

## ROI & vertical fit

Shared memory pays off wherever **the same context is read by many reasoners over time** — Healthcare (one patient/case context read by intake, triage, and discharge agents), Travel (a trip's evolving itinerary state every agent reads), Financial Services (org-wide risk config in `global`, per-deal overrides in `workflow`). The value is the eliminated Redis/pub-sub infrastructure plus the avoided recompute; the cost trap is over-scoping scratch to `global`. See **`references/capability-playbook.md`** for the full ROI = (value − cost) ÷ cost × 100 treatment across all five verticals.
