# Memory Events — Reactive Programming Over Shared State

Reasoners are the unit of intelligence. Triggers are the unit of arrival. **Memory events are the unit of reaction** — the way one reasoner wakes up *because another reasoner wrote a value*, with no polling, no queue wiring, and no direct call between them.

Pick memory events when the system is **state-driven**: "when the risk score for any account crosses a threshold, re-underwrite"; "when warehouse stock for a SKU drops, re-plan fulfillment". The producer just writes to shared memory; the consumer subscribes to a key pattern and fires on the write. The two reasoners never reference each other — they're coupled only through the key namespace.

This is the capability no static chain framework has: a LangGraph/CrewAI DAG must declare the edge from producer to consumer at build time. AgentField lets the edge **emerge from a memory write at runtime**.

This guide is the **local, offline** reference. It does not require a network fetch.

---

## The one decorator that matters

```python
@app.on_change("warehouse.stock.*")
async def on_stock_change(event):
    # fires whenever ANY key matching warehouse.stock.* is written
    ...
```

`@app.on_change(pattern)` marks a function as a memory-change listener. When a write lands on a key the pattern matches, the control plane delivers the change to the handler. `pattern` is a single glob string or a list of them.

```python
# One pattern
@app.on_change("ticket:*:sentiment")
async def on_ticket_sentiment(event):
    ...

# Several patterns, one handler
@app.on_change(["session.user_id", "session.permissions"])
async def on_session_change(event):
    ...
```

### Scope-bound subscription

`@app.on_change(...)` listens across the agent's default resolution. To bind a subscription to **one scope**, subscribe through the scope object:

```python
# Only fire on writes to the GLOBAL scope
@app.memory.global_scope.on_change("risk.scores.*")
async def on_risk_spike(event):
    ...
```

The scoped variants — `global_scope`, `session(id)`, `actor(id)`, `workflow(id)` — each expose `.on_change(...)`. Scope binding matters when the same logical key (`risk.scores.acct-9`) can be written in more than one scope and you only care about one of them. See `references/shared-memory.md` for the four scopes.

---

## The event object

The handler receives one `event` argument. Verified fields:

| Field | Meaning |
|---|---|
| `event.key` | The memory key that changed (e.g. `warehouse.stock.SKU-42`) |
| `event.path` | The full dotted path — use `event.path.endswith("theme")` to branch on the leaf |
| `event.data` | The **new** value just written |
| `event.previous_data` | The value before this write (`None` on first write) |
| `event.scope` | Isolation scope of the write — `global` / `session` / `actor` / `workflow` |
| `event.scope_id` | The scope identifier (session id, actor id, workflow/run id) |
| `event.action` | The change kind (e.g. `set`) |

```python
@app.on_change("user.preferences.*")
async def handle_preference_change(event):
    log_info(f"{event.key}: {event.previous_data} -> {event.data}")
    if event.path.endswith("theme"):
        await app.call(f"{app.node_id}.repaint", theme=event.data)
```

`previous_data` is what makes events *useful* and not just notifications — it lets a handler act on the **delta** (crossed a threshold, flipped a flag, grew past a budget) rather than the absolute value.

---

## Pattern syntax

Glob over the dotted/colon key namespace. Design your **keys** so the patterns you'll want to subscribe to fall out naturally.

| Pattern | Matches |
|---|---|
| `warehouse.stock.*` | `warehouse.stock.SKU-1`, `warehouse.stock.SKU-2`, … |
| `ticket:*:sentiment` | `ticket:T-123:sentiment`, `ticket:T-9:sentiment` |
| `risk.scores.*` | every account's risk score |
| `session.permissions` | exactly that key, no wildcard |

**Namespace the keys for the subscriptions you want.** `warehouse.stock.<sku>` is a good key precisely because `warehouse.stock.*` is a meaningful subscription. A flat key like `stock_SKU_42` cannot be watched as a group.

---

## Architectural rules

**The producer and consumer must stay decoupled.** The whole point is that the writer doesn't know who's listening. If a reasoner writes `risk.scores.X` *and then* directly `app.call`s the re-underwriter, you've built a chain and the event is dead weight. Write the value, return. Let the subscription do the rest.

**Handlers are entry points — keep them thin.** Like triggered reasoners, an `on_change` handler is a router: read the event, decide if it's interesting, hand off to specialists via `app.call(...)`. Heavy synthesis inside the handler is the same anti-pattern as heavy synthesis inside a webhook reasoner.

**Guard against write storms and loops.** A handler that writes a key its own pattern matches will re-fire itself. Always either (a) write to a *different* namespace than you subscribe to, or (b) compare `event.data` against `event.previous_data` and return early when the change is a no-op.

**Act on the delta, not the level.** Threshold logic belongs in the handler: `if event.previous_data < T <= event.data:`. Firing expensive downstream work on every write — when only crossings matter — is the most common cost leak.

**Subscriptions are at-least-once.** The same write may be delivered more than once under retry. Make the handler idempotent: scope its own writes by `event.scope_id`, or read-before-write.

**Don't poll what you can subscribe to.** A cron reasoner that reads a key every minute to see if it changed is a memory event waiting to be written. Cron for *time*; `on_change` for *state*.

---

## When NOT to use memory events

- The consumer needs the result **synchronously in the same request** — just `app.call(...)` it. Events are fire-and-forget reactions, not return values.
- Only one reasoner ever reads the value, right after one reasoner writes it, in the same flow — a direct call is clearer than an event.
- The "change" is an external event (Stripe, GitHub) rather than an internal state write — that's a **trigger** (`references/triggers.md`), not a memory event.

---

## Concrete example — reactive risk re-underwriting

A scoring reasoner writes risk scores; a separate watcher re-underwrites only on a threshold crossing. They share no code path.

```python
# Producer — knows nothing about who watches
@app.reasoner()
async def score_account(account_id: str, signals: dict) -> dict:
    score = await app.call(f"{app.node_id}.compute_risk", signals=signals)
    await app.memory.global_scope.set(f"risk.scores.{account_id}", score["value"])
    return {"account_id": account_id, "score": score["value"]}


# Consumer — fires on the write, acts only on the crossing
@app.memory.global_scope.on_change("risk.scores.*")
async def on_risk_spike(event):
    THRESHOLD = 0.8
    prev = event.previous_data or 0.0
    if prev < THRESHOLD <= event.data:          # crossed UP into danger
        account_id = event.key.rsplit(".", 1)[-1]
        await app.call(
            f"{app.node_id}.reunderwrite",
            account_id=account_id,
            score=event.data,
        )
```

`score_account` and `on_risk_spike` are independently deployable, independently testable, and the edge between them is the key namespace — not a hardcoded call.

---

## Memory-event anti-patterns

| ❌ | ✅ |
|---|---|
| Polling a key on a cron to detect change | `@app.on_change("key.*")` fires on the write |
| Producer `app.call`s the consumer right after writing | Write and return; the subscription is the edge |
| Heavy reasoning inside the `on_change` handler | Handler routes; specialists via `app.call` + `asyncio.gather` |
| Handler writes a key its own pattern matches | Write to a different namespace, or diff `data` vs `previous_data` and bail |
| Firing downstream work on every write | Act on the delta — `prev < T <= data` |
| Flat keys (`stock_SKU_42`) you later want to watch as a group | Namespaced keys (`warehouse.stock.<sku>`) so `*` is meaningful |
| Assuming exactly-once delivery | Idempotent handler — scope writes by `event.scope_id` |
| Using `on_change` for external webhooks | External events are triggers — `references/triggers.md` |

---

## Smoke test for a reactive build

After `docker compose up`:

1. Write a value the pattern matches, twice, straddling the threshold:
   ```bash
   curl -sS -X POST http://localhost:8080/api/v1/execute/<slug>.score_account \
     -H 'Content-Type: application/json' \
     -d '{"input":{"account_id":"acct-9","signals":{"...":"..."}}}'
   ```
2. Confirm the watcher ran **without anyone calling it**: the re-underwrite execution should appear on `/runs` with no inbound `app.call` from `score_account` — its parent is the memory write, not the producer reasoner.
3. Write a *below-threshold* value and confirm the watcher fired but took the early-return path (no downstream run). That proves the delta logic, not just the subscription.

If the second write triggers a re-underwrite run that you never called directly, the reactive edge is healthy.

---

## ROI & vertical fit

The reactive pattern earns its keep wherever **a state change must drive expensive action, but only on the meaningful change** — Insurance (risk crossings → re-underwrite), Financial Services (`risk.scores.*` → freeze/alert), Retail (`warehouse.stock.*` → re-plan fulfillment). The value is the eliminated polling loop and the eliminated producer→consumer coupling; the cost trap is firing on every write instead of the delta. See **`references/capability-playbook.md`** for the full ROI = (value − cost) ÷ cost × 100 treatment across all five verticals.
