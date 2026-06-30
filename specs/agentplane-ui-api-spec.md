# Silmari UI API Specification

**Purpose:** the exact HTTP surface the Silmari UI API must expose so the **existing** `control-plane/web/client`
UI is **copied in unchanged** and works. Derived from the UI's actual calls + the Go spec
(`control-plane/internal/server/routes_ui.go`, handlers under `control-plane/internal/handlers/ui/`).

**UI base URL:** `/api/ui/v1` (hard-coded; overridable via `VITE_API_BASE_URL`). A few flows call `/api/v1/*` directly.

**Status legend:** ✅ implemented in Silmari · ♻️ exists under `/api/v1` (alias/reuse) · ❌ missing

**Totals:** `/api/ui/v1` = 74 · `/api/ui/v2` = 2 · **core UI-API = 76**. Plus UI-observed extras (§3) and `/api/v1` extras the UI calls (§4).

---

## 1. `/api/ui/v1` (74 endpoints) — Go spec: `routes_ui.go`, handlers `internal/handlers/ui/`

### 1.1 agents/ — lifecycle, config, env, packages (17)

| # | Method | Path | Go handler | Status |
|---|--------|------|-----------|--------|
| 1 | GET | `/agents/packages` | packages.ListPackagesHandler | ❌ |
| 2 | GET | `/agents/packages/:packageId/details` | packages.GetPackageDetailsHandler | ❌ |
| 3 | GET | `/agents/running` | lifecycle.ListRunningAgentsHandler | ❌ |
| 4 | GET | `/agents/:agentId/details` | (inline agent details) | ❌ |
| 5 | GET | `/agents/:agentId/status` | lifecycle.GetAgentStatusHandler | ❌ |
| 6 | POST | `/agents/:agentId/start` | lifecycle.StartAgentHandler | ❌ |
| 7 | POST | `/agents/:agentId/stop` | lifecycle.StopAgentHandler | ❌ |
| 8 | POST | `/agents/:agentId/reconcile` | lifecycle.ReconcileAgentHandler | ❌ |
| 9 | GET | `/agents/:agentId/config/schema` | config.GetConfigSchemaHandler | ❌ |
| 10 | GET | `/agents/:agentId/config` | config.GetConfigHandler | ❌ |
| 11 | POST | `/agents/:agentId/config` | config.SetConfigHandler | ❌ |
| 12 | GET | `/agents/:agentId/env` | env.GetEnvHandler | ❌ |
| 13 | PUT | `/agents/:agentId/env` | env.PutEnvHandler | ❌ |
| 14 | PATCH | `/agents/:agentId/env` | env.PatchEnvHandler | ❌ |
| 15 | DELETE | `/agents/:agentId/env/:key` | env.DeleteEnvVarHandler | ❌ |
| 16 | GET | `/agents/:agentId/executions` | agentExecution.ListExecutionsHandler | ❌ |
| 17 | GET | `/agents/:agentId/executions/:executionId` | agentExecution.GetExecutionDetailsHandler | ❌ |

### 1.2 nodes/ — status, details, logs, did (10)

| # | Method | Path | Go handler | Status |
|---|--------|------|-----------|--------|
| 18 | GET | `/nodes/summary` | uiNodes.GetNodesSummaryHandler | ❌ |
| 19 | GET | `/nodes/events` *(SSE)* | uiNodes.StreamNodeEventsHandler | ❌ |
| 20 | GET | `/nodes/:nodeId/status` | uiNodes.GetNodeStatusHandler | ❌ |
| 21 | POST | `/nodes/:nodeId/status/refresh` | uiNodes.RefreshNodeStatusHandler | ❌ |
| 22 | POST | `/nodes/status/bulk` | uiNodes.BulkNodeStatusHandler | ❌ |
| 23 | POST | `/nodes/status/refresh` | uiNodes.RefreshAllNodeStatusHandler | ❌ |
| 24 | GET | `/nodes/:nodeId/details` | uiNodes.GetNodeDetailsHandler | ❌ |
| 25 | GET | `/nodes/:nodeId/logs` | nodeLogs.ProxyNodeLogsHandler | ❌ |
| 26 | GET | `/nodes/:nodeId/did` | did.GetNodeDIDHandler | ❌ |
| 27 | GET | `/nodes/:nodeId/vc-status` | did.GetNodeVCStatusHandler | ❌ |

### 1.3 settings/ (2)

| # | Method | Path | Go handler | Status |
|---|--------|------|-----------|--------|
| 28 | GET | `/settings/node-log-proxy` | nodeLogSettings.GetNodeLogProxySettingsHandler | ❌ |
| 29 | PUT | `/settings/node-log-proxy` | nodeLogSettings.PutNodeLogProxySettingsHandler | ❌ |

### 1.4 executions/ — lists, details, notes, lifecycle, logs, vc (18)

| # | Method | Path | Go handler | Status |
|---|--------|------|-----------|--------|
| 30 | GET | `/executions/summary` | uiExecutions.GetExecutionsSummaryHandler | ❌ |
| 31 | GET | `/executions/stats` | uiExecutions.GetExecutionStatsHandler | ❌ |
| 32 | GET | `/executions/enhanced` | uiExecutions.GetEnhancedExecutionsHandler | ❌ |
| 33 | GET | `/executions/events` *(SSE)* | uiExecutions.StreamExecutionEventsHandler | ❌ |
| 34 | GET | `/executions/timeline` | timeline.GetExecutionTimelineHandler | ❌ |
| 35 | GET | `/executions/recent` | recentActivity.GetRecentActivityHandler | ❌ |
| 36 | GET | `/executions/:execution_id/details` | uiExecutions.GetExecutionDetailsGlobalHandler | ❌ |
| 37 | POST | `/executions/:execution_id/webhook/retry` | uiExecutions.RetryExecutionWebhookHandler | ❌ |
| 38 | POST | `/executions/:execution_id/cancel` | handlers.CancelExecutionHandler | ♻️ |
| 39 | POST | `/executions/:execution_id/pause` | handlers.PauseExecutionHandler | ♻️ |
| 40 | POST | `/executions/:execution_id/resume` | handlers.ResumeExecutionHandler | ♻️ |
| 41 | POST | `/executions/note` | handlers.AddExecutionNoteHandler | ❌ |
| 42 | GET | `/executions/:execution_id/notes` | handlers.GetExecutionNotesHandler | ❌ |
| 43 | GET | `/executions/:execution_id/logs` | execLogs.GetExecutionLogsHandler | ❌ |
| 44 | GET | `/executions/:execution_id/logs/stream` *(SSE)* | execLogs.StreamExecutionLogsHandler | ❌ |
| 45 | GET | `/executions/:execution_id/vc` | did.GetExecutionVCHandler | ❌ |
| 46 | GET | `/executions/:execution_id/vc-status` | did.GetExecutionVCStatusHandler | ❌ |
| 47 | POST | `/executions/:execution_id/verify-vc` | did.VerifyExecutionVCComprehensiveHandler | ❌ |

### 1.5 llm / queue (2)

| # | Method | Path | Go handler | Status |
|---|--------|------|-----------|--------|
| 48 | GET | `/llm/health` | llm.GetLLMHealthHandler | ❌ |
| 49 | GET | `/queue/status` | llm.GetExecutionQueueStatusHandler | ❌ |

### 1.6 workflows/ (7)

| # | Method | Path | Go handler | Status |
|---|--------|------|-----------|--------|
| 50 | GET | `/workflows/:workflowId/dag` | handlers.GetWorkflowDAGHandler | ✅ |
| 51 | POST | `/workflows/:workflowId/cancel-tree` | handlers.CancelWorkflowTreeHandler | ♻️ |
| 52 | DELETE | `/workflows/:workflowId/cleanup` | handlers.CleanupWorkflowHandler | ❌ |
| 53 | POST | `/workflows/vc-status` | did.GetWorkflowVCStatusBatchHandler | ❌ |
| 54 | GET | `/workflows/:workflowId/vc-chain` | did.GetWorkflowVCChainHandler | ❌ |
| 55 | POST | `/workflows/:workflowId/verify-vc` | did.VerifyWorkflowVCComprehensiveHandler | ❌ |
| 56 | GET | `/workflows/:workflowId/notes/events` *(SSE)* | workflowNotes.StreamWorkflowNodeNotesHandler | ❌ |

### 1.7 reasoners/ (7)

| # | Method | Path | Go handler | Status |
|---|--------|------|-----------|--------|
| 57 | GET | `/reasoners/all` | reasoners.GetAllReasonersHandler | ❌ |
| 58 | GET | `/reasoners/events` *(SSE)* | reasoners.StreamReasonerEventsHandler | ❌ |
| 59 | GET | `/reasoners/:reasonerId/details` | reasoners.GetReasonerDetailsHandler | ❌ |
| 60 | GET | `/reasoners/:reasonerId/metrics` | reasoners.GetPerformanceMetricsHandler | ❌ |
| 61 | GET | `/reasoners/:reasonerId/executions` | reasoners.GetExecutionHistoryHandler | ❌ |
| 62 | GET | `/reasoners/:reasonerId/templates` | reasoners.GetExecutionTemplatesHandler | ❌ |
| 63 | POST | `/reasoners/:reasonerId/templates` | reasoners.SaveExecutionTemplateHandler | ❌ |

### 1.8 dashboard/ (2)

| # | Method | Path | Go handler | Status |
|---|--------|------|-----------|--------|
| 64 | GET | `/dashboard/summary` | dashboard.GetDashboardSummaryHandler | ❌ |
| 65 | GET | `/dashboard/enhanced` | dashboard.GetEnhancedDashboardSummaryHandler | ❌ |

### 1.9 did/ (6)

| # | Method | Path | Go handler | Status |
|---|--------|------|-----------|--------|
| 66 | GET | `/did/status` | did.GetDIDSystemStatusHandler | ❌ |
| 67 | GET | `/did/export/vcs` | did.ExportVCsHandler | ❌ |
| 68 | POST | `/did/verify` | did.VerifyVCHandler | ❌ |
| 69 | POST | `/did/verify-audit` | did.VerifyAuditBundleHandler | ❌ |
| 70 | GET | `/did/:did/resolution-bundle` | did.GetDIDResolutionBundleHandler | ❌ |
| 71 | GET | `/did/:did/resolution-bundle/download` | did.DownloadDIDResolutionBundleHandler | ❌ |

### 1.10 vc/ (2)

| # | Method | Path | Go handler | Status |
|---|--------|------|-----------|--------|
| 72 | GET | `/vc/:vcId/download` | did.DownloadVCHandler | ❌ |
| 73 | POST | `/vc/verify` | did.VerifyVCHandler | ❌ |

### 1.11 authorization/ (1)

| # | Method | Path | Go handler | Status |
|---|--------|------|-----------|--------|
| 74 | GET | `/authorization/agents` | authorization.GetAgentsWithTagsHandler | ❌ |

---

## 2. `/api/ui/v2` (2 endpoints) — Go spec: `routes_ui.go` (`uiAPIV2` group)

| # | Method | Path | Go handler | Status |
|---|--------|------|-----------|--------|
| 75 | GET | `/api/ui/v2/workflow-runs` | workflowRuns.ListWorkflowRunsHandler | ❌ |
| 76 | GET | `/api/ui/v2/workflow-runs/:run_id` | workflowRuns.GetWorkflowRunDetailHandler | ❌ |

---

## 3. UI-observed extras under `/api/ui/v1/*` (confirm Go handler before building)

These appear in the UI source but are not in the `routes_ui.go` core block (may be in a sibling
Go route file / identity service). Confirm exact handler + shape, then implement.

| Method | Path | Notes | Status |
|--------|------|-------|--------|
| GET | `/api/ui/v1/identity/agents` | identity service | ❌ |
| GET | `/api/ui/v1/identity/agents/:id/details` | | ❌ |
| GET | `/api/ui/v1/identity/agents/:id/dids` | | ❌ |
| GET | `/api/ui/v1/identity/dids/search` | | ❌ |
| GET | `/api/ui/v1/identity/dids/stats` | | ❌ |
| GET | `/api/ui/v1/identity/credentials/search` | | ❌ |
| GET | `/api/ui/v1/nodes/:nodeId/mcp/health/stream` *(SSE)* | MCP health stream | ❌ |
| GET | `/api/ui/v1/executions/:execution_id/notes/stream` *(SSE)* | execution notes stream | ❌ |
| GET | `/api/ui/v1/executions/filter-options` | filter dropdowns | ❌ |
| GET | `/api/ui/v1/executions/view-stats` | | ❌ |

## 4. `/api/v1/*` endpoints the UI calls directly (not under /api/ui)

| Method | Path | Status |
|--------|------|--------|
| POST | `/api/v1/execute/:target` · `/execute/async/:target` | ✅ (Phase A) |
| GET | `/api/v1/executions/:id` | ✅ (Phase A) |
| GET/POST | `/api/v1/triggers` · `/triggers/:id` (CRUD) | ♻️ partial |
| GET | `/api/v1/triggers/metrics` | ❌ |
| GET/POST | `/api/v1/triggers/:id/events` · `/events/:id/replay` | ❌ |
| GET | `/api/v1/sources` | ❌ |
| GET/PUT/POST | `/api/v1/settings/observability-webhook` (+ `/status`, `/dlq`, `/redrive`) | ❌ |
| POST | `/api/v1/nodes/register-serverless` | ❌ |
| GET | `/api/v1/admin/policies` | ❌ |
| POST | `/api/v1/webhooks/approval-response` | ❌ |
| GET | `/api/v1/did/agentfield-server` · `POST /api/v1/did/verify-audit` | ♻️ partial |
| GET | `/api/v1/agentic/kb/articles/*` | DEFER (out of core scope) |

---

## Conformance approach

Build each Silmari controller to match the Go handler's response **exactly** (same JSON keys,
status codes, error shapes), validated the same way Phase A validated `/execute` — diff against the
Go handler output. The UI is the consumer of record: any shape mismatch breaks a page.

**Drop-in gate:** the unedited `web/client` builds, loads at `/ui/`, and every page's `/api/ui/v1`
call returns a Go-shaped 2xx (no 404/501 on a path the UI needs); the 6 SSE panels stream live.

See also the historical planning note `thoughts/searchable/shared/plans/agentplane-ui-api-worklist.md` (build sequence U1–U8).
