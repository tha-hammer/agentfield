# Capability Playbook — How to Actually Use Coordination, and What It's Worth

This is the reference engineers and architects read **before** they design, to decide *which* coordination capability a problem wants and *whether the math closes*. The other three references tell you how the API works:

- `references/triggers.md` — events/schedules arrive from the outside world
- `references/memory-events.md` — reasoners react to internal state writes
- `references/shared-memory.md` — reasoners share state across time and agents

This file tells you **when each one pays for itself**, in the language a CTO and a CMO both have to sign off on.

---

## The north star

> **ROI% = (value − cost) ÷ cost × 100**

Every capability decision routes through this. The three AgentField coordination primitives move ROI on **both sides of the fraction**:

- They **raise value** — reactions happen in seconds not hours, context is never re-derived, events are never missed.
- They **cut cost** — the control plane already owns the plumbing you would otherwise build, staff, and operate: webhook receivers, HMAC verification, replay tables, a Redis/pub-sub tier, cache invalidation, and an audit trail.

The trap that *destroys* ROI is always the same: doing expensive work on the wrong cadence — firing on every write instead of the meaningful delta, polling what you could subscribe to, recomputing what a sibling already wrote. The playbook below is mostly about keeping the **cost** term small.

### The cost you delete by adopting the primitive

| Hand-rolled today | AgentField primitive | Cost removed |
|---|---|---|
| Webhook endpoint + signature verification + replay table + dispatch | **Triggers** | An entire ingress service per provider, and its on-call |
| Cron loops polling a row to detect change | **Memory events** | Wasted compute every tick + the latency of the poll interval |
| Redis/pub-sub tier + cache invalidation + recompute | **Shared memory** | A stateful infra tier, its ops, and duplicated LLM spend |

Each row is a line item that moves from your **cost** column into the control plane's. That is the structural ROI lever — before a single use case.

---

## Choosing the primitive (decision lens)

```
Where does the work come FROM?

├─ The outside world pushes an event (provider webhook, cron)   → TRIGGERS
├─ An internal state value changed and that should drive action → MEMORY EVENTS
└─ Several reasoners need the same context over time            → SHARED MEMORY
```

Most real systems use **all three at once**: a trigger *arrives*, writes to *shared memory*, and the write fires a *memory event* that fans out. The triage example at the bottom shows the composition. Design the key namespace and the event patterns together — good keys (`risk.scores.<acct>`) make good subscriptions (`risk.scores.*`) free.

---

## The five verticals

Each vertical below names the highest-ROI shape, the primitive it leans on, where the **value** comes from, where the **cost** trap hides, and the dominant ROI driver. Treat these as starting topologies to adapt — not templates to copy. (Insurance, Travel, Healthcare, Financial Services, Retail.)

### 🛡️ Insurance

| | |
|---|---|
| **Highest-ROI shape** | FNOL / claim webhook → triage → reactive re-underwriting on risk crossings |
| **Primitives** | Triggers (claim arrives) + Shared memory (claim context) + Memory events (`risk.scores.*` crossing) |
| **Where value comes from** | Straight-through processing of low-risk claims; humans see only the exceptions. Cycle time drops from days to seconds on the happy path. Every decision carries a verifiable credential chain — audit and dispute defense come free. |
| **Cost trap** | Re-underwriting on every score write instead of the threshold *crossing* (`prev < T <= data`). LLM spend balloons; ROI inverts. |
| **Dominant ROI driver** | **Value↑** — analyst hours reallocated from triage to exceptions; **Cost↓** — deleted claims-intake plumbing. |

Worked math (illustrative): if 60% of claims are low-risk and can flow straight through, and each avoided manual triage is worth ~$40 of analyst time, a 100k-claim/yr book saves ~$2.4M of triage labor. Against a build+run cost on the order of $200k, ROI ≈ (2.4M − 0.2M) ÷ 0.2M × 100 ≈ **1,100%**. The crossing-not-level discipline is what keeps the cost denominator honest.

### ✈️ Travel

| | |
|---|---|
| **Highest-ROI shape** | Airline/PNR-change webhook → rebooking agent reading the live itinerary from shared memory |
| **Primitives** | Triggers (schedule-change webhook) + Shared memory (`trip:<pnr>.itinerary` every agent reads) |
| **Where value comes from** | Disruptions (cancellations, delays) are handled before the traveler calls. The itinerary lives once in shared memory; the rebooker, the notifier, and the fare-rules checker all read it instead of re-fetching. |
| **Cost trap** | Webhook storms during a weather event re-triggering full re-plans. Keep the trigger reasoner thin; debounce by writing a "replanning" flag and diffing it. Over-scoping the itinerary to `global` when it's a per-trip (`workflow`/`session`) concern. |
| **Dominant ROI driver** | **Value↑** — disruptions resolved proactively lift NPS and cut call-center volume; **Cost↓** — no per-provider webhook ingress to operate. |

CMO read: proactive rebooking is a *marketable* differentiator ("we rebooked you before you noticed") — the value term includes retention and word-of-mouth, not just deflected support tickets.

### 🏥 Healthcare

| | |
|---|---|
| **Highest-ROI shape** | One patient/case context in shared memory, read by intake → triage → discharge agents over a long-running episode |
| **Primitives** | Shared memory (`case:<id>` context, scoped `session`/`workflow`) + Memory events (status change → next-step agent) |
| **Where value comes from** | No reasoner re-derives what an earlier one already established; the case context is the single source of truth across a multi-hour or multi-day episode. The verifiable-credential chain provides the provenance regulated workflows require. |
| **Cost trap** | Putting PHI-bearing context in `global` (wrong blast radius) or in memory at all when it should be ephemeral and `workflow`-scoped. Scope is a compliance control here, not just a perf knob. Secrets/PHI handling must follow your governing regime — memory scope is necessary, not sufficient. |
| **Dominant ROI driver** | **Cost↓** — eliminated recompute and eliminated coordination infra; **Value↑** — fewer hand-off errors across agents. |

Architect note: here the *narrowest-scope* rule is a safety rule. `workflow` scope that auto-clears on run completion is the default; promote to `session`/`actor` deliberately, never to `global` for case data.

### 💳 Financial Services

| | |
|---|---|
| **Highest-ROI shape** | Payment/risk webhook → reconciliation, with `risk.scores.*` memory events driving freeze/alert on crossings |
| **Primitives** | Triggers (Stripe/payment events) + Memory events (risk crossing) + Shared memory (org risk config in `global`, per-deal overrides in `workflow`) |
| **Where value comes from** | Sub-second reaction to risk crossings (freeze, alert, escalate) where minutes meant loss. Org-wide policy lives in `global`; a deal shadows it in `workflow`; readers never branch on scope because hierarchical `get` resolves it. Every action is cryptographically auditable for the regulator. |
| **Cost trap** | Acting on the risk *level* every write instead of the *crossing*; smuggling request/response through memory instead of `app.call`. Both inflate cost with no value. |
| **Dominant ROI driver** | **Value↑** — loss prevented in the seconds saved; **Cost↓** — deleted reconciliation plumbing + audit trail you'd otherwise build. |

CFO read: the value term is dominated by *avoided loss per second of reaction time* — model it as (incidents × loss-rate × seconds-saved), and the reactive primitive's ROI is usually decided there, not on labor.

### 🛒 Retail

| | |
|---|---|
| **Highest-ROI shape** | `warehouse.stock.*` memory events → re-plan fulfillment; order webhooks → orchestration |
| **Primitives** | Memory events (stock deltas) + Triggers (order/inventory webhooks) + Shared memory (catalog/inventory snapshot) |
| **Where value comes from** | Stockouts and re-routes are handled the instant inventory crosses a threshold, not on the next nightly batch. One inventory snapshot in shared memory serves pricing, fulfillment, and availability reasoners. |
| **Cost trap** | Subscribing `warehouse.stock.*` and firing a full re-plan on *every* decrement instead of on the low-water crossing. At catalog scale this is the single biggest cost leak in the playbook. |
| **Dominant ROI driver** | **Value↑** — fewer oversells and faster re-routes lift conversion and margin; **Cost↓** — no polling tier over the inventory table. |

CMO read: "in stock when we said it was" is a conversion and trust lever — the value term includes abandoned-cart recovery, not only operational savings.

---

## Cross-vertical summary

| Vertical | Lead primitive | Value term is mostly… | Cost trap to watch |
|---|---|---|---|
| Insurance | Memory events | Analyst hours → straight-through | Re-underwrite on level, not crossing |
| Travel | Triggers + shared memory | Retention + deflected support | Webhook storms; itinerary over-scoped |
| Healthcare | Shared memory | Eliminated recompute + fewer hand-off errors | PHI in `global`; wrong scope |
| Financial Services | Memory events | Loss avoided per second saved | Acting on level; memory-as-queue |
| Retail | Memory events | Conversion + margin | Re-plan on every decrement |

The pattern across all five: **the primitive cuts a structural cost line (plumbing/infra), and the discipline of acting on the *delta/crossing* keeps the variable cost (LLM spend) from eating the gain.** Get the cadence right and ROI is decided by the value term — which is where the business case actually lives.

---

## A worked composition — the three primitives together

Insurance triage, end to end, using all three:

```python
# 1. TRIGGER — a claim arrives from the outside world (thin entry)
@app.reasoner()
@on_event(source="generic_hmac", types=["claim.filed"], secret_env="CLAIMS_WEBHOOK_SECRET")
async def on_claim(event: dict, trigger: TriggerContext | None = None):
    claim_id = event["claim_id"]
    # 2. SHARED MEMORY — write the claim context once; everyone reads it
    await app.memory.workflow(claim_id).set(f"claim:{claim_id}.context", event)
    score = await app.call(f"{app.node_id}.score_risk", claim_id=claim_id, claim=event)
    # write the score into a watched namespace
    await app.memory.global_scope.set(f"risk.scores.{claim_id}", score["value"])
    return {"claim_id": claim_id, "score": score["value"]}


# 3. MEMORY EVENT — re-underwrite ONLY on the threshold crossing
@app.memory.global_scope.on_change("risk.scores.*")
async def on_risk_spike(event):
    THRESHOLD = 0.8
    if (event.previous_data or 0.0) < THRESHOLD <= event.data:
        claim_id = event.key.rsplit(".", 1)[-1]
        ctx = await app.memory.workflow(claim_id).get(f"claim:{claim_id}.context", default={})
        await app.call(f"{app.node_id}.reunderwrite", claim_id=claim_id, context=ctx)
```

The trigger stays thin. The context is written once and read by the watcher. The expensive `reunderwrite` fires only on the crossing. That is the ROI math made into code: plumbing deleted (the trigger), recompute deleted (shared context), variable cost contained (delta-gated reaction).

---

## Bottom line for the architect

1. **Name the source of work** (outside event / state change / shared context) and the primitive falls out.
2. **Design keys and event patterns together** so subscriptions are free.
3. **Gate every expensive reaction on the delta or crossing**, never the level — this is where ROI is won or lost.
4. **State the value and cost terms explicitly** before building. If `(value − cost) ÷ cost × 100` doesn't clear your hurdle rate with the cost-deletion lines above already counted, the use case isn't ready — change the shape, not the spreadsheet.
