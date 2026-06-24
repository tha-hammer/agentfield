---
date: 2026-06-12
topic: "AgentPlane UI-API surface work-list — make the existing web UI drop in unchanged"
status: ready-to-build
repository: agentplane (Elixir port) / agentfield (Go spec + UI source)
---

# AgentPlane UI-API Work-List — "the UI just drops in"

## END STATE (non-negotiable)

The **existing, already-built** React UI at `control-plane/web/client` is **COPIED** into
AgentPlane and served by Phoenix — **not ported, not reworked, not rewritten**. To make that
true, AgentPlane's backend must expose the **same `/api/ui/v1` (+ `/api/ui/v2`) HTTP surface**
the UI already calls, with the **same request/response shapes** as Go AgentField. When that
surface exists, copying the UI is:

1. `vite build` the **unchanged** `control-plane/web/client` source (base `/ui`, `VITE_API_BASE_URL=/api/ui/v1`)
2. drop the `dist/` into AgentPlane and serve it via Phoenix `Plug.Static` at `/ui`
3. open `http://localhost:4000/ui/` → it works, because every `/api/ui/v1/*` call resolves

**This work-list is the backend surface (steps that make #1–#3 possible). The UI code itself is COPIED, never edited.**

The UI's base URL is hard-coded to `/api/ui/v1` (`src/services/configurationApi.ts:4`,
`src/components/WorkflowDAG/hooks/useNodeDetails.ts:58`, `src/components/StepDetail.tsx:240`),
overridable by `VITE_API_BASE_URL`. A few flows also call `/api/v1/*` directly (execute, triggers).

## CURRENT GAP

AgentPlane's `/api/ui/v1` scope exposes **exactly 1** route today:
`GET /api/ui/v1/workflows/:workflowId/dag`. The UI calls **~74** `/api/ui/v1` + **2** `/api/ui/v2`
routes. So **~75 endpoints are missing** — copying the UI over today would 404 on nearly every call.

Spec source (Go): `control-plane/internal/server/routes_ui.go` (the `uiAPI := Router.Group("/api/ui/v1")`
block) + handlers under `control-plane/internal/handlers/ui/`. Build AgentPlane controllers to match
each handler's response shape exactly (same JSON keys, same status codes), as Phase A did for `/execute`.

Legend: ✅ exists in AgentPlane · ♻️ exists under `/api/v1` (reuse/alias into `/api/ui/v1`) · ❌ missing

---

## WORK ITEM U1 — Serve the UI (the copy-over mechanism)

- [ ] Vendor the **unchanged** `control-plane/web/client` source into AgentPlane (e.g. `apps/agentplane_web/ui_src/`) — a COPY, no edits.
- [ ] Build step: `pnpm install && pnpm build` with `VITE_BASE_PATH=/ui`, `VITE_API_BASE_URL=/api/ui/v1`, `VITE_BUILD_OUT_DIR=../priv/static/ui` (or copy dist into `apps/agentplane_web/priv/static/ui`).
- [ ] Phoenix `endpoint.ex`: add `Plug.Static` at `"/ui"` from `priv/static/ui`, with SPA fallback (serve `index.html` for unknown `/ui/*` paths so client-side routing works).
- [ ] Verify: `GET /ui/` returns the app shell; static assets 200.
- [ ] **Acceptance:** the copied UI loads at `/ui/` with no source edits; failing API calls are 404s on *missing endpoints* below, not on the static serving.

## WORK ITEM U2 — Dashboard + executions read surface (the landing pages)

Go spec: `routes_ui.go:105-151`, handlers `internal/handlers/ui/{executions,dashboard,recent_activity,timeline}.go`.

- [ ] ❌ `GET /api/ui/v1/dashboard/summary` — DashboardSummary
- [ ] ❌ `GET /api/ui/v1/dashboard/enhanced` — EnhancedDashboardSummary
- [ ] ❌ `GET /api/ui/v1/executions/summary` — ExecutionsSummary
- [ ] ❌ `GET /api/ui/v1/executions/stats` — ExecutionStats
- [ ] ❌ `GET /api/ui/v1/executions/enhanced` — EnhancedExecutions (list + filters)
- [ ] ❌ `GET /api/ui/v1/executions/recent` — RecentActivity
- [ ] ❌ `GET /api/ui/v1/executions/timeline` — ExecutionTimeline
- [ ] ❌ `GET /api/ui/v1/executions/filter-options` — filter dropdown options (UI calls it; confirm Go handler)
- [ ] ❌ `GET /api/ui/v1/executions/view-stats` — view stats (UI calls it; confirm Go handler)
- [ ] ❌ `GET /api/ui/v1/executions/:execution_id/details` — GetExecutionDetailsGlobal
- [ ] ♻️ `POST /api/ui/v1/executions/:execution_id/cancel|pause|resume` — reuse Phase A `ExecuteController`/`Executions`
- [ ] ❌ `POST /api/ui/v1/executions/note` + `GET /api/ui/v1/executions/:execution_id/notes` — execution notes (needs a notes table)
- [ ] ❌ `POST /api/ui/v1/executions/:execution_id/webhook/retry` — RetryExecutionWebhook

## WORK ITEM U3 — Nodes/agents management surface

Go spec: `routes_ui.go:53-101`, handlers `internal/handlers/ui/` (nodes, lifecycle, config, env, packages).

- [ ] ❌ `GET /api/ui/v1/nodes/summary`
- [ ] ❌ `GET /api/ui/v1/nodes/:nodeId/status` · `POST /nodes/:nodeId/status/refresh` · `POST /nodes/status/bulk` · `POST /nodes/status/refresh`
- [ ] ❌ `GET /api/ui/v1/nodes/:nodeId/details`
- [ ] ❌ `GET /api/ui/v1/nodes/:nodeId/logs` — ProxyNodeLogs
- [ ] ❌ `GET /api/ui/v1/nodes/:nodeId/did` · `GET /nodes/:nodeId/vc-status`
- [ ] ❌ `GET /api/ui/v1/agents/running`
- [ ] ❌ `GET /api/ui/v1/agents/:agentId/details` · `/status` · `POST /start` · `/stop` · `/reconcile`
- [ ] ❌ `GET /api/ui/v1/agents/:agentId/config/schema` · `GET /config` · `POST /config`
- [ ] ❌ `GET /api/ui/v1/agents/:agentId/env` · `PUT /env` · `PATCH /env` · `DELETE /env/:key`
- [ ] ❌ `GET /api/ui/v1/agents/packages` · `GET /agents/packages/:packageId/details`
- [ ] ❌ `GET /api/ui/v1/agents/:agentId/executions` · `GET /:agentId/executions/:executionId`
- [ ] ❌ `GET /api/ui/v1/settings/node-log-proxy` · `PUT /settings/node-log-proxy`

## WORK ITEM U4 — Reasoners surface

Go spec: `routes_ui.go:170-179`, handler `internal/handlers/ui/reasoners*.go`.

- [ ] ❌ `GET /api/ui/v1/reasoners/all`
- [ ] ❌ `GET /api/ui/v1/reasoners/:reasonerId/details` · `/metrics` · `/executions` · `GET/POST /templates`

## WORK ITEM U5 — Workflows surface (DAG is the one piece already done)

Go spec: `routes_ui.go:154-166`.

- [ ] ✅ `GET /api/ui/v1/workflows/:workflowId/dag` — already in AgentPlane
- [ ] ♻️ `POST /api/ui/v1/workflows/:workflowId/cancel-tree` — reuse Phase A cancel-tree
- [ ] ❌ `DELETE /api/ui/v1/workflows/:workflowId/cleanup` — CleanupWorkflow
- [ ] ❌ `POST /api/ui/v1/workflows/vc-status` · `GET /:workflowId/vc-chain` · `POST /:workflowId/verify-vc` (DID/VC — see U7, may DEFER as stub)
- [ ] ♻️/❌ `GET /api/ui/v2/workflow-runs` · `GET /api/ui/v2/workflow-runs/:run_id` — workflow-runs list/detail

## WORK ITEM U6 — Real-time SSE streams (6 of them; the UI opens EventSource)

These power live updates; the UI keeps `EventSource` connections open. Phase A already has the
PubSub plumbing + `GET /api/v1/executions/:id/events`; extend to the UI's stream paths.

- [ ] ❌ `GET /api/ui/v1/nodes/events` — node status stream
- [ ] ❌ `GET /api/ui/v1/executions/events` — execution stream
- [ ] ❌ `GET /api/ui/v1/executions/:execution_id/logs/stream` — log stream
- [ ] ❌ `GET /api/ui/v1/executions/:execution_id/notes/stream` (a.k.a. `/workflows/:id/notes/events`) — notes stream
- [ ] ❌ `GET /api/ui/v1/reasoners/events` — reasoner stream
- [ ] ❌ `GET /api/ui/v1/nodes/:nodeId/mcp/health/stream` — MCP health stream
- [ ] **Acceptance:** each returns `text/event-stream`, sends a snapshot on connect, pushes change events, heartbeats; UI live panels update without reload.

## WORK ITEM U7 — DID / VC surface (DEFER-able as 501 if DID is out of scope)

Go spec: `routes_ui.go:100,141-143,160-162,191-218`. AgentPlane has DID stubs today.
The UI tolerates these gracefully if they return the Go "not configured" shapes. Decide: implement vs stub.

- [ ] ❌ `GET /api/ui/v1/did/status` · `GET /did/export/vcs` · `POST /did/verify` · `POST /did/verify-audit`
- [ ] ❌ `GET /api/ui/v1/did/:did/resolution-bundle` (+ `/download`)
- [ ] ❌ `GET /api/ui/v1/vc/:vcId/download` · `POST /vc/verify`
- [ ] ❌ `GET /api/ui/v1/executions/:execution_id/vc` · `/vc-status` · `POST /verify-vc`
- [ ] ❌ `GET /api/ui/v1/authorization/agents` — GetAgentsWithTags
- [ ] ❌ identity endpoints the UI calls under `/api/ui/v1/identity/*` (agents, dids/search, credentials/search, dids/stats) — confirm against Go `identity.go`

## WORK ITEM U8 — `/api/v1/*` endpoints the UI also calls directly (not under /api/ui)

The UI hits these straight (execute, triggers, settings). Phase A covered execute; the rest:

- [ ] ✅ `/api/v1/execute/:target`, `/execute/async/:target`, `GET /executions/:id` — Phase A done
- [ ] ♻️/❌ `/api/v1/triggers`, `/triggers/metrics`, `/triggers/:id`, `/triggers/:id/events`, `/triggers/:id/events/:id/replay` — AgentPlane has triggers CRUD; **metrics/events/replay missing**
- [ ] ❌ `/api/v1/sources` (list) — AgentPlane only has `POST /sources/:trigger_id` stub
- [ ] ❌ `/api/v1/settings/observability-webhook` (+ `/status`, `/dlq`, `/redrive`)
- [ ] ❌ `/api/v1/nodes/register-serverless`
- [ ] ❌ `/api/v1/admin/policies` · `/api/v1/webhooks/approval-response`
- [ ] DEFER `/api/v1/agentic/kb/articles/*` (agentic KB — out of core scope)

---

## BUILD SEQUENCE (for SWE-AF)

1. **U1** (serve the copied UI) — do first so progress is visible in the browser.
2. **U2 + U3** (dashboard/executions/nodes read surface) — unblocks the main pages.
3. **U4 + U5** (reasoners + workflows).
4. **U6** (SSE) — live updates.
5. **U7 + U8** (DID/VC + /api/v1 extras) — implement or stub per scope.

## GATE (how we know the UI "just dropped in")

- Copied UI builds with **zero source edits** and loads at `/ui/`.
- Every page the UI renders makes its `/api/ui/v1` calls and gets **2xx with Go-shaped bodies** (no 404/501 on a path the UI needs).
- The 6 SSE panels stream live.
- Cross-check each response shape against the Go handler in `internal/handlers/ui/` (same keys/status as the conformance approach).

Related: [agentplane-parity-worklist.md](agentplane-parity-worklist.md) (Phase A dispatch engine — DONE, live E2E passed 2026-06-12).
