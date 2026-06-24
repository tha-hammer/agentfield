---
date: 2026-06-10
title: AgentPlane Parity Work-List (Go AgentField → Elixir port)
status: source-of-truth
repo_under_port: AgentPlane (Elixir/Phoenix umbrella, /workspaces/project/agentplane)
spec_source: Go AgentField control-plane (/home/maceo/Dev/agentfield/control-plane, READ-ONLY spec)
builder: SWE-AF on :8090 (swe-planner node) — SWE-AF builds, not hand-coding
---

# AgentPlane Parity Work-List

## END-STATE (locked — this is "done", not a number)
**The port works when (1) a real SDK agent registers with AgentPlane, the control plane routes a reasoner call from one agent to another (sync AND async), the workflow is tracked end-to-end, and pause/resume/cancel propagate; AND (2) the existing web UI `dist/` copied over from the Go control plane loads and works FLAWLESSLY against AgentPlane — every page renders real data, real-time SSE updates flow, no broken calls.** Verified by a live agent-to-agent E2E + a browse-the-whole-UI pass, not by the conformance fixture runner.
The conformance runner is a **regression check only**, permanently demoted from "definition of done."

## HARD CONSTRAINTS (NON-GOALS — locked)
- **The web UI = a COPY OF THE SOURCE, built inside AgentPlane. NOT a prebuilt dist, NOT reimplemented.** Copy the `control-plane/web/client/` React/TS **source tree** into the AgentPlane repo (vendored, lives with the port). AgentPlane's build compiles it from source (pnpm install + `vite build`, with `VITE_BASE_PATH=/ui` and output into Phoenix's static dir); Phoenix serves it at `/ui/` (Plug.Static + SPA index.html fallback). The React/TS source is **copied UNCHANGED — not rewritten/reimplemented** ("don't build from scratch" = don't rewrite the UI); only build config (env vars) is set for the new home. A prebuilt/frozen dist is REJECTED as unmaintainable. So: a real build toolchain (Node + pnpm) is part of the AgentPlane build. Frontend code work = ZERO (it's the same source); the real work is the BACKEND endpoints the UI calls + matching response shapes.
- The Go control plane is the READ-ONLY spec; do not modify it.

## STAY-ON-TRACK RULES
1. This file is the single source of truth. Every codepath below is a row with an explicit status; "done" = the row's acceptance met **by a real request**, not a fixture.
2. **SWE-AF builds it.** Each build's goal references rows in this file; the Go handler (file:line) is the spec. No hand-coding the port.
3. The gate that flips Phase-A rows to ✅ is the **live E2E** (below), not the conformance count.
4. No pivoting, re-scoping, or "done" declarations until the targeted phase's rows are all ✅. If a measurement looks good but rows are open, the port is NOT done.
5. Scope is **tiered**: build the agent-runtime critical path first; UI/admin/connector/agentic are explicitly deferred — they don't make the port "work," they make it manageable.

## THE LIVE E2E GATE (Phase-A acceptance)
Boot AgentPlane on :4000 (PG up, migrations applied). Run two real SDK agents (`af init`, `AGENTFIELD_SERVER=http://localhost:4000`): caller + callee. Then:
- callee registers + heartbeats; `GET /discovery/capabilities` shows it.
- caller does `agent.call("callee.reasoner", input=…)` → CP dispatches to callee, returns the result. Workflow row created with parent/child.
- async variant: 202 + execution_id, callee 202-acks, posts status back, caller's wait resolves.
- pause/resume/cancel a running execution → propagates to the agent.
All four pass = Phase A ✅ = **the port runs a real project.**

## THE UI E2E GATE (Phase-UI acceptance)
Copy the `control-plane/web/client/` SOURCE into AgentPlane; AgentPlane builds it from source (pnpm + vite, `VITE_BASE_PATH=/ui`, output → Phoenix static); Phoenix serves it at `/ui/` (Plug.Static + SPA index.html fallback). Open `http://localhost:4000/ui/` in a browser against AgentPlane (with real agents running per Phase A). Every page must render with live data and no failed requests:
- Dashboard (summary/enhanced, llm/health, queue/status, trigger metrics)
- Agents (nodes summary/details/status, start/stop, env, config, logs stream)
- Executions + Workflows (summary/enhanced/stats/timeline/recent, DAG, run list v2, details, logs/notes streams, cancel/pause/resume)
- Reasoners (all/details/metrics/history/templates, reasoner events stream)
- Identity/DID + Authorization + Triggers + Settings (observability-webhook, node-log-proxy)
- The **6 SSE streams** push live updates: nodes/events, executions/events, executions/:id/logs/stream, executions/:id/notes/stream, reasoners/events, nodes/:id/mcp/health/stream.
**The UI uses relative paths and same-origin** — no client rebuild needed; it just needs the endpoints + matching response shapes.

---

## ELIXIR CURRENT STATE (2026-06-10, branch feature/1f313100)
- **Deps:** phoenix, ecto_sql, postgrex, oban, rustler, jason, phoenix_pubsub, plug_cowboy. **MISSING: outbound HTTP client (no finch/req/httpoison)** — cannot call agents.
- **Migrations (tables):** health_checks, memory_entries, nodes. **MISSING: executions, workflow_executions, workflow_execution_events, agent_nodes(full), triggers, vectors.**
- **Real controllers:** nodes (register/heartbeat/status/lifecycle), memory (set/get/delete/list), triggers (CRUD), executions (CRUD-only: show/cancel/pause/resume/status_update — **no dispatch, no real storage**), did (status/resolve/verify-audit/execution-vc), workflow (events/dag).
- **Contexts:** Nodes, Memory, Triggers, Executions (record stubs), DID, Crypto (Rustler NIF), Workflow (record_event), Storage(kv), ExecutionManager (GenServer skeleton), HealthPoller.
- **Conformance:** 22/24 fixtures (peripheral surface). Does NOT exercise the dispatch engine.

---

## BUILD PLAN (phased; each codepath = a discrete SWE-AF issue, Go file:line = spec)

### PHASE A — THE DISPATCH ENGINE (makes the port work; all currently STUB) ⬅ build first
| # | Codepath | Go spec (file:line) | Elixir status | Acceptance |
|---|---|---|---|---|
| A0 | Add outbound HTTP client (Req or Finch) + executions/workflow_executions/workflow_execution_events Ecto schemas+migrations | models.go:5-30,115-180 | MISSING | tables migrate; client callable |
| A1 | **Sync execute** `POST /execute/:target`: parse `node.reasoner`, look up node base_url, forward POST to `{base_url}/reasoners/{name}` with X-Execution-ID/X-Run-ID/X-Workflow-ID/X-Parent-Execution-ID/X-Session-ID/X-Actor-ID headers, create execution+workflow rows (status=running), return result | execute.go:169,1011-1302 | STUB 501 | caller gets callee's result; rows created |
| A2 | **Async execute** `POST /execute/async/:target`: 202 + execution_id immediately, background task does the dispatch, status reaches terminal | execute.go:175,358-419,2117-2183 | STUB 501 | 202 then pollable terminal status |
| A3 | **Status callback** `POST /executions/:id/status`: agent posts terminal/progress; update rows; resolve waiting sync caller via a registry/PubSub (replace Go's ExecutionEventBus) | execute.go:193,474-684,858-931 | partial (record only) | sync caller blocked on 202-ack unblocks on callback |
| A4 | **Execution status/get** `GET /executions/:id` (+ batch-status) | execute.go:181,421-440 | partial | returns live status/result |
| A5 | **Pause/resume/cancel propagation**: cancel → POST `{base_url}/_internal/executions/{id}/cancel`; pause blocks dispatch via wait-for-resume | execute_cancel.go:31, execute_pause.go:38, cancel_dispatcher.go:107 | record-only | running execution actually cancels/pauses at the agent |
| A6 | **Workflow cancel-tree** `POST /workflows/:id/cancel-tree` + **workflow/executions/events** intake (local sub-exec reporting) | execute_cancel_tree.go:62, workflow_execution_events.go:35 | partial | whole run cancels; local sub-execs recorded |
| A7 | **Execution events SSE** `GET /executions/:id/events` | reasoner_catalog.go:130 | MISSING | live event stream to a tailing client |

### PHASE B — AGENT LIFECYCLE COMPLETENESS (registration/discovery the SDK relies on)
| Codepath | Go (file:line) | Elixir | 
|---|---|---|
| `GET /discovery/capabilities` (SDK + agentic use it) | discovery.go:188 | STUB |
| `GET /reasoners` catalog | reasoner_catalog.go:38 | MISSING |
| nodes: `GET /nodes`, `GET /nodes/:id`, status/refresh, start/stop, lease PATCH, actions ack/claim, shutdown, register-serverless | nodes_*.go | STUB (4 of 18 done) |
| executions: logs intake, notes, approval-request/status, awaiter-status, batch-status | execution_logs.go, execution_notes.go, execute_approval.go, execute_awaiter_status.go | MISSING |

### PHASE C — PERIPHERAL AGENT-FACING (round out parity)
| Domain | Done | Missing |
|---|---|---|
| memory | set/get/delete/list ✓ | vector set/get/search/delete, events ws/sse/history, append, query |
| did | status/resolve/verify-audit/execution-vc ✓ | register, verify, vc-chain, workflow vc, export, document, issuer-key, well-known did:web |
| triggers | CRUD ✓ | public ingest `/sources/:id`, events list/stream/replay, pause/resume, test, sources list |
| webhooks | — | approval-response ingest |

### PHASE UI — WEB UI PARITY (REQUIRED — the dist/ must copy over and work flawlessly)
**Spec = exactly what the client calls (125 endpoints) + the client's TS response types (must match byte-shape).** Source: web/client TS types (`src/types/*.ts`, `services/*.ts` inline interfaces) + the Go `internal/handlers/ui/*` handlers. Relative same-origin paths → no client rebuild. Sequenced AFTER Phase A (UI reads real execution/workflow data).

| Group | Endpoints client calls | Spec source | Notes |
|---|---|---|---|
| UI-BUILD+SERVE | copy `control-plane/web/client/` SOURCE into AgentPlane repo; AgentPlane build = pnpm install + `vite build` (`VITE_BASE_PATH=/ui`, outDir → Phoenix static); serve at `/ui/` + SPA fallback | vite.config.ts, embedded/ui.go:18-54 | adds Node+pnpm to AgentPlane build toolchain; React source vendored UNCHANGED, built in-repo (NOT a prebuilt dist) |
| UI-DASH (5) | `/api/ui/v1/dashboard/summary\|enhanced`, `/llm/health`, `/queue/status`, `/triggers/metrics` | ui/dashboard.go, ui/execution_logs.go | health strip |
| UI-NODES (12) | `/api/ui/v1/nodes/summary\|:id/details\|:id/status\|status/bulk\|status/refresh\|:id/start\|:id/stop`, `nodes/events`(SSE), `:id/logs`(NDJSON), `:id/did`, `:id/mcp/*` | ui/nodes.go, ui/node_logs.go, ui/did.go | |
| UI-AGENTS (14) | `/api/ui/v1/agents/running\|packages\|packages/:id/details\|:id/status\|:id/start\|:id/stop\|:id/reconcile\|:id/env(GET/PUT/PATCH/DELETE)\|:id/config(GET/POST)\|:id/config/schema` | ui/lifecycle.go, ui/config.go, ui/env.go, ui/packages.go | start/stop/env/config |
| UI-EXEC (18) | `/api/ui/v1/executions/summary\|enhanced\|stats\|timeline\|recent\|filter-options\|view-stats\|:id/details\|:id/logs\|:id/notes\|note\|:id/cancel\|:id/pause\|:id/resume\|:id/webhook/retry`, `executions/events`(SSE), `:id/logs/stream`(SSE), `:id/notes/stream`(SSE) | ui/executions.go, ui/execution_timeline.go, ui/recent_activity.go | reads Phase-A data |
| UI-WORKFLOW (8 + v2) | `/api/ui/v1/workflows/:id/details\|dag\|dag?mode=lightweight\|cancel-tree\|cleanup\|vc-chain\|vc-status\|verify-vc`; `/api/ui/v2/workflow-runs`, `/workflow-runs/:id` | workflow_dag.go, ui/workflow_runs.go | DAG + run list |
| UI-REASONERS (6) | `/api/ui/v1/reasoners/all\|:id/details\|:id/metrics\|:id/executions\|:id/templates(GET/POST)`, `reasoners/events`(SSE) | ui/reasoners.go | |
| UI-DID/IDENTITY (17) | `/api/ui/v1/did/*` (register/resolve/document/status/agents/verify/verify-audit/workflow vc/export/resolution-bundle), `/identity/*`, `/authorization/agents`, exec/workflow `vc`/`vc-status`/`verify-vc` | ui/did.go, ui/identity.go, ui/authorization.go | |
| UI-SETTINGS (2) | `/api/ui/v1/settings/node-log-proxy (GET/PUT)` | ui/node_log_settings.go | |
| UI-via-/api/v1 (32) | client also calls agent routes: `execute(/async)/:id`, `executions/:id`, `triggers*`, `admin/policies*`, `admin/tags`, `admin/agents/*-tags`, `settings/observability-webhook*`, `webhooks/approval-response`, `nodes/register-serverless`, `did/agentfield-server`, `sources` | (covered by Phases A/B/C + admin) | overlaps other phases |
| UI-AUTH | every call sends `X-API-Key` (and `X-Admin-Token` for admin); SSE sends `?api_key=` | — | auth plug |

**Gate:** the UI E2E above — browse every page flawlessly, SSE live. **Critical shape-match risks** (silently break UI): execution `status` string values, workflow-run `status_counts` map, `WorkflowRunListResponse.has_more`, `WorkflowVCChainResponse.component_vcs`, `ExecutionTimelineResponse.cache_timestamp`.

### PHASE Z — TRULY DEFERRED (client never calls these)
agentic API (11), connector API (~29), memory vector (UI doesn't use memory at all), did:web well-known, some admin/config. Build only on explicit request.

---

## FULL ROUTE SURFACE (counts, for completeness — see catalog agents for every file:line)
nodes 18 · executions 15 · execute 4 · workflow 3 · webhooks 1 · memory 14 · did 16 · triggers 18 · discovery 2 · agentic 11 · admin 16 · connector ~29 · observability 7 · health 2 · did:web 2 · UI v1 ~80 · UI v2 2 · ui-static 2 = **~240 routes**. Core-orchestration (agent-runtime) = nodes+execute+executions+workflow+webhooks+discovery+memory ≈ **57 routes**; that is the working-port target. The rest is peripheral/management.

## BUILD SEQUENCE (each = one SWE-AF build on :8090, Go handlers/TS types as spec)
1. **Phase A** — dispatch engine (A0-A7). Gate: live agent-to-agent E2E. *Makes executions/workflows real — everything else reads this data.*
2. **Phase B** — agent lifecycle + discovery (capabilities, reasoner catalog, full nodes). Gate: SDK registers + discovers.
3. **Phase UI** — serve dist/ + the 125 client endpoints + 6 SSE + matching TS shapes + auth. Gate: browse the whole UI flawlessly. *Depends on A+B for real data.*
4. **Phase C** — remaining agent-facing peripherals (triggers ingest, did completeness, webhooks/approval, admin policies/tags for the UI).
5. **Phase Z** — deferred (agentic/connector/memory-vector), only on request.

NEXT: fire **Phase A** first (the engine); UI parity (Phase UI) is required but sequenced after the engine produces real data. Go handlers + the client TS types are the spec; acceptance is the E2E + UI gates, never the conformance count.
