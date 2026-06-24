---
date: 2026-06-04T19:49:16-04:00
researcher: tha-hammer
git_commit: fb8559d5d88d42e8b15043ff746611ee482152ba
branch: main
repository: agentfield
topic: "How the code makes AgentField 'the AI Backend' — API-addressable, scalable, autonomous, coordinated, governed agents"
tags: [research, codebase, control-plane, sdk-python, sdk-go, execution, triggers, did-vc, observability, storage, memory]
status: complete
last_updated: 2026-06-04
last_updated_by: tha-hammer
---

# Research: How the code makes AgentField "the AI Backend"

**Date**: 2026-06-04T19:49:16-04:00
**Researcher**: tha-hammer
**Git Commit**: fb8559d5d88d42e8b15043ff746611ee482152ba
**Branch**: main
**Repository**: agentfield (github.com/Agent-Field/agentfield)

> GitHub permalink base for every `file:line` below (HEAD is pushed to `origin/main`):
> `https://github.com/Agent-Field/agentfield/blob/fb8559d5d88d42e8b15043ff746611ee482152ba/<file>#L<line>`

## Research Question

This app claims: *"AgentField is the AI Backend: open-source infrastructure for turning agents into API-addressable, scalable services that can run autonomously, coordinate across systems, and be governed in production."* Research how the code makes this possible.

## Summary

AgentField backs the claim with a control plane (Go/Gin) plus two SDKs (Python on FastAPI, Go on `net/http`) that together turn ordinary agent functions into governed production services. The mechanism behind each phrase of the claim is concrete in the code:

- **API-addressable** — An agent function decorated `@reasoner` (Python) or registered via `RegisterReasoner` (Go) becomes a real HTTP route, `POST /reasoners/{name}`, on the agent's own server (`agent.py:1808`; `agent.go:797`). The control plane never holds agent logic; it is a router that maps a dotted target `node_id.reasoner_name` to that registered endpoint and POSTs to it (`execute.go:1010-1302`, `buildAgentURL` `:1540-1560`).
- **Scalable services** — A `sync.Once`-initialized async worker pool (default `NumCPU` workers, 1024-deep queue) absorbs load, returns `202` immediately, and sheds with HTTP `503 concurrency_limit` when full (`execute.go:2215-2230`). A per-agent atomic concurrency limiter (`429` on overflow), a unified storage interface that switches SQLite+BoltDB ↔ PostgreSQL by config (`storage.go:789-835`), a streaming file payload store, and a stateless per-container SDK rate limiter with circuit breaker round out the scale story.
- **Run autonomously** — A self-registering Source plugin system (GitHub/Slack/Stripe HMAC-verified webhooks + a cron `LoopSource`) converts external events into deduplicated `InboundEvent` rows that a `TriggerDispatcher` POSTs straight to an agent reasoner — no human, no caller (`triggers.go:59-173`, `trigger_dispatcher.go:56-186`, `source_manager.go`). Completion webhooks, DB-busy retries, and a full pause/resume/cancel state machine (including bottom-up workflow-tree cancel bridged to `asyncio.Task.cancel()`) keep long-running work alive.
- **Coordinate across systems** — `agent.call("other.fn", ...)` always routes through the control plane's async execute endpoint (never agent→agent directly), propagating run/session/parent headers that stitch every hop into one workflow DAG (`agent.py:3710`, `client.py:850-889`, `execute.go:1242-1252`). Operator-facing connectors expose fleet management under token auth.
- **Governed in production** — Every server/agent/reasoner/skill gets a deterministic `did:key` derived by HKDF from one encrypted master seed (`did_service.go:562-585`); each execution can mint an Ed25519-signed W3C Verifiable Credential, chained per-workflow and exportable as a self-verifying audit bundle (`GET /api/v1/did/workflow/{id}/vc-chain`) that `af verify` checks fully offline (`vc_chain.go`, `verify_provenance.go`). Layered on top: tag-based access policies, a tag-approval lifecycle gate (signed `AgentTagCredential`s), human-in-the-loop execution approval over HMAC webhooks, Prometheus metrics, OTel tracing, a batching observability forwarder with dead-letter queue, and scoped memory (run/session/actor/global).

In short, the claim is realized by a strict separation: agent frameworks author behavior; AgentField wraps that behavior in HTTP addressing, an async/queued gateway, an event-driven trigger fabric, a header-propagated workflow DAG, and a cryptographic identity + policy + observability layer. The unifying thread across all of it is a single `ExecutionContext` carried as `X-*` headers. One notable gap versus documentation: the `internal/mcp/` integration referenced in `CLAUDE.md` is absent from the current tree.

## Official Positioning (vendor docs)

The AgentField docs (https://agentfield.ai/docs/learn/vs-frameworks) frame the product as operating at the **infrastructure** layer, not the authoring layer: *"Agent frameworks help you author agent behavior... AgentField is the AI Backend,"* exposing agents *"as production services with APIs, async execution, distributed calls, deployment, identity, and audit."* Named primitives — Agent, Reasoner, Skill, Memory, Harness, Execution — map 1:1 onto the code below. Each documented capability has a concrete code locus: Agent APIs → reasoner/skill HTTP routes; distributed agent services → `node.function` targets; async execution → queue/retry/webhook/resume; service discovery → `/discover` + node registry; fleet observability → DAG/metrics/traces; memory scopes → run/session/actor/global; identity & policy → DID + access policy; proof & audit → signed VC chains.

## Detailed Findings

### 1. API-addressability — functions become HTTP endpoints; the control plane routes to them

The Python `Agent` subclasses `FastAPI`, so the agent instance *is* the ASGI app (`agent.py:446`). `@reasoner` stamps metadata onto the function (`decorators.py:49-170`); `Agent.reasoner()` then computes `endpoint_path = f"/reasoners/{reasoner_id}"` (`agent.py:1775`) and registers `@self.post(endpoint_path)` (`agent.py:1808`), recording a `ReasonerEntry` in `self._reasoner_registry` (`agent.py:1971-1980`). The endpoint closure validates input and branches sync vs. async on the presence of an `X-Execution-ID` header (`agent.py:1858-1909`). Infra routes (`/health`, `/reasoners`, `/skills`, `/info`, `/status`, `/shutdown`, `/webhooks/approval`) are wired by `AgentServer.setup_agentfield_routes()`, and boot ends in `uvicorn.run(self.agent)` (`agent_server.py:133-958`).

The Go SDK is symmetric: `RegisterReasoner` stores `a.reasoners[name]` (`agent_register.go:6-45`), and `agent.handler()` builds an `http.ServeMux` (`agent.go:791`) wiring `POST /reasoners/{name}` → `handleReasoner` (`agent.go:797,1190`) plus `/health`, `/discover`, `/execute`, `/_internal/executions/`. `handleReasoner` reads the same `X-*` execution headers and dispatches async via goroutine (HTTP 202) when both `ExecutionID` and `AgentFieldURL` are set (`agent.go:1190-1397`).

The **control plane is purely a router**. `routes_core.go:27-150` registers `POST /api/v1/execute/:target` and `POST /api/v1/execute/async/:target` (plus legacy `/reasoners/:reasoner_id`, `/skills/:skill_id`, and `GET /reasoners`). `prepareExecution` calls `parseTarget` to split `node_id.reasoner_name` (`execute.go:1512-1524`), `store.GetAgent` to fetch the registered node's `BaseURL`, and `determineTargetType` to classify reasoner vs. skill (`execute.go:1526-1533`). `callAgent` → `buildAgentURL` then constructs `{BaseURL}/reasoners/{name}` (or `/skills/`, or `/execute` for serverless, or a stored `InvocationURL`) and POSTs with the full `X-*` header set + `Authorization` bearer + DID headers (`execute.go:1210-1560`).

**Connection:** the same dotted-target grammar used by external API callers is what `agent.call()` uses internally — addressing is uniform whether the caller is `curl` or another agent.

### 2. Registration & discovery — how an agent joins the fleet

At startup the SDK `ConnectionManager.start()` calls `register_with_agentfield_server()`, which POSTs `node_id` + `reasoners` + `skills` + `base_url` + `vc_metadata` + `instance_id` to `POST /api/v1/nodes/register` (`agent_field_handler.py:41-185`, `connection_manager.py:35-83`). On `pending_approval`, the SDK polls `GET /api/v1/nodes/{id}` every 5s; heartbeats POST to `/api/v1/nodes/{id}/heartbeat` from a daemon thread.

`RegisterNodeHandler` validates the `AgentNode`, normalizes `proposed_tags` → `tags`, resolves callback URLs, preserves `ApprovedTags` on re-registration (`nodes_register.go:527`), detects redeploys via `InstanceID` mismatch and orphans in-flight executions (`:641-666`), persists via `RegisterAgent` (`:630`), publishes `NodeRegistered` + invalidates the discovery cache (`:635-636`), runs DID registration (`:682`), and upserts code-managed triggers (`:741`). Serverless agents are discovered by GETting `{invocation_url}/discover` behind an SSRF guard (`nodes_register.go:776-1072`), setting `DeploymentType="serverless"`. A 2-second `HeartbeatCache` throttles heartbeat DB writes (`nodes_heartbeat.go:25-72`).

**Connection:** registration is the single point that fans out into three subsystems at once — DID identity (§6), tag approval (§7), and trigger wiring (§4).

### 3. Workflow execution & DAG — header propagation builds the call graph

`prepareExecution` is the hub for both sync and async paths (`execute.go:1010-1197`). It derives `runID` from `X-Run-ID` (or generates one), always mints a fresh `executionID`, and reads `X-Parent-Execution-ID`. The `Execution.ParentExecutionID` field is the DAG edge. Records are written to **two tables** — `executions` and `workflow_executions` — and `deriveWorkflowHierarchy` inherits `RootWorkflowID` and increments `WorkflowDepth` from the parent row (`execute.go:1799-1897`). When forwarding, `callAgent` sets `X-Run-ID` / `X-Workflow-ID` / `X-Parent-Execution-ID` (`execute.go:1242-1252`); the receiving agent echoes them on any onward `call()`, so a tree of HTTP calls self-correlates without a central coordinator.

`buildExecutionDAG` reconstructs the tree from all executions sharing a `RunID` — building `execMap` + `childrenMap`, finding the root (`ParentExecutionID == nil`), recursing with cycle detection, and computing depth (`workflow_dag.go:465-733`). SDK-internal sub-steps that never hit the gateway can still appear in the DAG via `POST /api/v1/workflow/executions/events` (`workflow_execution_events.go:35-231`). Lifecycle is broadcast over an in-process `ExecutionEventBus` singleton (`GlobalExecutionEventBus`, `events/execution_events.go:114`) consumed by sync-waiters, the OTel tracer, the observability forwarder, and SSE fan-out to the UI.

**Connection:** the event bus is the spine — pause/resume waits, sync-call completion, tracing, and live UI all subscribe to the same stream.

### 4. Autonomous operation — triggers, retries, pause/cancel

External autonomy comes from the **Source plugin registry** (`sources/source.go:79-110`): `HTTPSource` plugins (GitHub `X-Hub-Signature-256` HMAC, Slack `v0` signing with replay window, Stripe `t=/v1=`) and a `LoopSource` cron plugin that sleeps until `schedule.Next` and emits `tick` events (`cron/cron.go:80-239`). Plugins self-register via `init()` and are aggregated by `sources/all/all.go`. `upsertCodeManagedTriggers` wires reasoner trigger bindings at registration and starts loop sources immediately (`triggers_register.go:57-130`).

Inbound webhooks hit the no-auth `POST /sources/:trigger_id`, which resolves the secret from `os.Getenv(SecretEnvVar)`, calls `sources.HandleHTTP`, filters by event type (prefix match), dedups by idempotency key, persists an `InboundEvent`, and fires `go dispatcher.DispatchEvent` after returning 200 to the provider (`triggers.go:59-173`). `TriggerDispatcher.DispatchEvent` loads the agent, checks health/lifecycle, verifies the `TargetReasoner`, optionally mints a trigger-event VC (passed as `X-Parent-VC-ID`), and POSTs `{event,_meta}` **directly** to `{BaseURL}/reasoners/{TargetReasoner}` with a 30s timeout (`trigger_dispatcher.go:56-186`). `SourceManager` owns the `KindLoop` goroutines, re-fetching each trigger per tick so runtime pauses take effect (`source_manager.go:26-212`).

Durability: the async worker pool's `process()` checks pause/cancel before each `callAgent` and waits on the bus when paused (`execute.go:2117-2183`). Completion webhooks retry with exponential backoff and a restart-surviving poller (`webhook_dispatcher.go:113-469`). DB-busy errors (`database is locked`/`SQLITE_BUSY`/`deadlock`) retry 5× with jittered backoff (`retry.go:16-29`). The pause/resume/cancel state machine spans control plane and SDK: `CancelDispatcher` subscribes to the bus and POSTs `/_internal/executions/{id}/cancel`, which the SDK turns into `asyncio.Task.cancel()` (`cancel_dispatcher.go:107-232`, `cancel.py:65-115`); `CancelWorkflowTreeHandler` cancels a whole run leaf-first by depth (`execute_cancel_tree.go:62-254`). `PauseClock` subtracts human-approval wait time so a child awaiting approval doesn't exhaust the parent's timeout (`agent_pause.py:10-95`).

**Connection:** the trigger path mirrors `callAgent` but is driven by the world, not a caller; the trigger VC becomes the chain root for the resulting execution VC (§6).

### 5. Cross-system coordination — calls always go through the plane

`Agent.call("node.fn", ...)` requires a dotted target, refuses to run if disconnected, builds headers from the current `ExecutionContext` (setting `X-Parent-Execution-ID` to itself), submits to `POST /api/v1/execute/async/{target}`, then polls for the result, propagating multi-hop `waiting`/`running` transitions upward via `on_child_waiting`/`on_child_running` (`agent.py:3710-4097`, `client.py:774-889`). The Go SDK has **no** outbound `Call()` on `Agent`; cross-agent calls there go through `client.Client` directly.

Operator coordination is via **connectors** — `/api/v1/connector/*` under `ConnectorTokenAuth`, only when `Features.Connector.Enabled` and a token are set (`routes_connector.go:17-44`). Capabilities are individually gated: `reasoner_management` (CRUD, versions, traffic weights 0–10000, restart), `status_read`, `policy_management`, `tag_management`; `GET /manifest` is always available (`connector/handlers.go:52-157`).

**Gap (documented vs. present):** `CLAUDE.md` references `control-plane/internal/mcp/`, but that directory does not exist and there are zero `mcp` route/handler matches in `internal/server` or `internal/handlers`. Only stale shell scripts and a coverage artifact reference MCP endpoints. MCP is **absent** from the current tree.

### 6. Governance via cryptographic identity (DID/VC)

One 32-byte master seed is generated once and stored AES-256-GCM-encrypted (PBKDF2, 600k rounds, binary magic `AFENC2`) (`encryption.go:19-189`, `did_service.go:41-86`). Every identity is derived deterministically: `HKDF(masterSeed, salt="agentfield-did-key-derivation-v1", info=derivationPath)` → Ed25519 key (`did_service.go:562-574`), encoded as `did:key:z` + base64url(multicodec `0xed 0x01` || pubkey) (`did_service.go:577-585`; note base64url, not the spec's base58btc). Derivation paths are BIP32-shaped: agent `m/44'/{serverHash}'/{idx}'`, reasoner `.../0'/{idx}'`, skill `.../1'/{idx}'`. Registration mints one DID per agent, reasoner, and skill (`did_service.go:152-388`), persisting the registry with the master seed re-encrypted (`did_registry.go:329-418`).

Per execution, `GenerateExecutionVC` (gated on `config.Enabled && RequireVCForExecution`) resolves the caller's signing identity, optionally SHA-256-hashes input/output, builds a W3C `AgentFieldExecutionCredential`, and signs the proof-stripped JSON with Ed25519, attaching an `Ed25519Signature2020` proof with `proofPurpose: assertionMethod` (`vc_issuance.go:16-276`). VCs carry `ParentVCID`, so trigger VCs chain to execution VCs (`vc_issuance_trigger.go:37-158`). `GetWorkflowVCChain` assembles all execution VCs + an on-demand signed `AgentFieldWorkflowCredential` + a `DIDResolutionBundle` snapshot of every public-key JWK, making the exported audit **self-verifying offline** (`vc_chain.go:14-123`). `af verify <file>` / `POST /api/v1/did/verify-audit` parse enhanced/legacy/bare formats, resolve each DID from the embedded bundle (online `did:web` resolution disabled), re-check Ed25519 signatures, and validate ~11 metadata fields with a 0–100 score (`cli/vc.go:36-197`, `verify_provenance.go:17-196`, `verify_audit.go:15-47`).

**Connection:** the same execution-context headers that build the DAG (§3) carry the caller/target DIDs, so identity and graph topology are issued from one propagation mechanism.

### 7. Governance via access control, approval & observability

Tag-based `AccessPolicyService` evaluates cross-agent calls against a priority-sorted policy cache — caller/target tag intersection, deny-first, allow-lists, parameter constraints (`<=,>=,==,!=,<,>`, fail-closed when params are nil), wildcards — first match wins (`access_policy_service.go:26-424`). `TagApprovalService` gates capability tags at registration into `auto`/`manual`/`forbidden`, parks agents in `PendingApproval`, and on approval issues a signed `AgentTagCredential` verified later by `TagVCVerifier` (revoked/expired/signature/issuer-subject binding) (`tag_approval_service.go:400-567`, `tag_vc_verifier.go:34-77`).

Human-in-the-loop: `RequestApprovalHandler` moves a running execution to `waiting` with `approval_status=pending` (72h default expiry) and publishes `ExecutionWaiting` (`execute_approval.go:52-214`). `ApprovalWebhookHandler` resolves it over an HMAC-verified webhook — `normalizeDecision` maps approve/continue/confirm → approved, reject/deny/abort/cancel → rejected; approved→running, rejected/expired→cancelled, request_changes→running with IDs cleared — publishing `ExecutionApprovalResolved` (`webhook_approval.go:104-537`). The SDK `_PauseManager` resolves the corresponding approval future.

Observability: Prometheus gateway metrics (`agentfield_gateway_queue_depth`, `agentfield_waiters_inflight`, `agentfield_gateway_backpressure_total`) (`execution_metrics.go`); OTel spans driven off the event buses (`execution_tracer.go:17-143`); 10 MiB structured-log ingestion with retention pruning (`execution_logs.go:34-155`); a 10s agent `HealthMonitor` with heartbeat-stale gating + recovery debounce (`health_monitor.go:47-461`); a 3-state LLM circuit breaker feeding execution preconditions (`llm_health_monitor.go:19-343`); a batching `observabilityForwarder` with HMAC signing, dead-letter queue, and `Redrive` (`observability_forwarder.go:59-741`); and anonymized OSS telemetry with a hashed install ID and a 16-key property allowlist (`telemetry.go:47-447`).

**Connection:** access policy, concurrency, and LLM health all funnel through one gate — `CheckExecutionPreconditions` enforces `429` (concurrency), `503` (LLM circuit open), and policy denial before any agent is called (`execution_guards.go:40-127`).

### 8. Scalability & storage backends

The async worker pool (`getAsyncWorkerPool`, `sync.Once`) runs `AGENTFIELD_EXEC_ASYNC_WORKERS` goroutines (default `runtime.NumCPU()`) over a `AGENTFIELD_EXEC_ASYNC_QUEUE_CAPACITY` channel (default 1024); a full queue returns HTTP `503 concurrency_limit` (`execute.go:2215-2230`). The per-agent `AgentConcurrencyLimiter` uses a `sync.Map` of atomic counters; overflow → `429 ErrorCategoryConcurrencyLimit` (`agent_concurrency.go:12-78`, `execution_guards.go:40-127`).

Storage is backend-agnostic: `StorageProvider` is one large interface, and `StorageFactory.CreateStorage` switches on `Mode` (`local`/`postgres`, overridable by `AGENTFIELD_STORAGE_MODE`) (`storage.go:789-835`). A single `*LocalStorage` struct backs both — SQLite (WAL, `synchronous=NORMAL`, `cache_size=10000`, 256 MB mmap, `MaxOpenConns` default 1, write-serializing mutex) + BoltDB for KV, or a PostgreSQL pool (`MaxOpenConns` 25, `MaxIdleConns` 5, `ConnMaxLifetime` 30m) (`local.go:469-767`). `FilePayloadStore` streams large payloads to disk via `io.MultiWriter` + SHA-256, atomic rename, `payload://` URIs, and ctx-cancellable copy (`payload_store.go:35-172`). The SDK `StatelessRateLimiter` adds a per-container circuit breaker (open at 5 failures, 30s timeout), up to 6 retries, 429/503 detection, `Retry-After` honoring, and `HOSTNAME+PID`-seeded jitter so distinct containers diverge (`rate_limiter.py:18-280`).

### 9. State & memory model (cross-cutting)

Memory is keyed `(scope, scope_id, key)`. `resolveScope` reads an explicit body scope, else header priority `X-Workflow-ID` → `workflow`, `X-Session-ID` → `session`, `X-Actor-ID` → `actor`, default `global`; `getScopeID` maps each scope to its header value (`memory.go:402-434`). Unscoped reads walk the hierarchy `workflow → session → actor → global`. The Python `MemoryInterface` exposes `.session()/.actor()/.workflow()` scoped clients plus WS/SSE change streaming with `filepath.Match` key patterns (`memory.py:509-884`, `memory_events.go:56-204`). `ExecutionContext` — `run_id, execution_id, parent_execution_id, session_id, actor_id, workflow_id (=run_id), root_workflow_id, caller_did, target_did, agent_node_did, parent_vc_id` — is the single `contextvars`-backed dataclass marshaled to/from the `X-*` headers that drive routing, DAG, memory scoping, and VC issuance alike (`execution_context.py:27-302`).

## Code References

- `sdk/python/agentfield/agent.py:446` — Agent subclasses FastAPI; the agent instance is the ASGI app
- `sdk/python/agentfield/decorators.py:49-170` — `@reasoner` stamps metadata, delegates to `_execute_with_tracking`
- `sdk/python/agentfield/agent.py:1775-1980` — `reasoner()` registers `POST /reasoners/{id}` (`:1808`) and `ReasonerEntry` registry (`:1971`)
- `sdk/python/agentfield/agent_server.py:133-958` — standard infra routes and uvicorn boot
- `sdk/go/agent/agent_register.go:6-45` — Go `RegisterReasoner` makes handler available at `/reasoners/{name}`
- `sdk/go/agent/agent.go:791-1397` — `http.NewServeMux` (`:791`), `/reasoners/` mux (`:797`), `handleReasoner` (`:1190`), async goroutine + status callback
- `control-plane/internal/server/routes_core.go:27-150` — execute/reasoner/skill route registration + DID middleware
- `control-plane/internal/handlers/execute.go:1010-1302` — `prepareExecution` (`:1010`), `callAgent` (`:1210`), `parseTarget` (`:1512`)
- `control-plane/internal/handlers/execute.go:1540-1560` — `buildAgentURL` (serverless/skill/reasoner/InvocationURL)
- `control-plane/internal/handlers/execute.go:2215-2230` — async worker pool (`sync.Once`, `NumCPU` workers, 1024 queue, 503 shedding)
- `control-plane/internal/handlers/execute.go:1799-1897` — dual-table persistence + `deriveWorkflowHierarchy`
- `control-plane/internal/handlers/workflow_dag.go:465-733` — `buildExecutionDAG` reconstructs tree from `RunID`
- `control-plane/internal/events/execution_events.go:13-114` — `ExecutionEventType` constants + `GlobalExecutionEventBus` singleton (`:114`)
- `control-plane/internal/handlers/nodes_register.go:350-1072` — `RegisterNodeHandler` + serverless `/discover` (SSRF guard)
- `sdk/python/agentfield/agent_field_handler.py:41-250` — `register_with_agentfield_server` + heartbeat
- `sdk/python/agentfield/connection_manager.py:35-83` — connect/reconnect/health-check loops
- `control-plane/internal/handlers/agent_concurrency.go:12-78` — per-agent atomic concurrency limiter (`Acquire` `:44`, `atomic.AddInt64` `:52`)
- `control-plane/internal/storage/storage.go:789-835` — `StorageFactory.CreateStorage` switching local (`:805`) / postgres (`:819`) via `AGENTFIELD_STORAGE_MODE`
- `control-plane/internal/storage/local.go:469-767` — one `LocalStorage` struct backs SQLite(WAL)+BoltDB and Postgres pools
- `control-plane/internal/services/payload_store.go:35-172` — streaming file payload store (SHA-256, atomic rename)
- `sdk/python/agentfield/rate_limiter.py:18-280` — stateless per-container rate limiter + circuit breaker
- `control-plane/internal/sources/source.go:79-110` — `Source` (`:79`) / `HTTPSource` (`:99`) / `LoopSource` (`:107`) plugin contract
- `control-plane/internal/sources/cron/cron.go:80-239` — cron `LoopSource` emits tick on `schedule.Next`
- `control-plane/internal/sources/github/github.go:41-78` — GitHub HMAC-SHA256 webhook verification
- `control-plane/internal/handlers/triggers.go:59-173` — `IngestSourceHandler` (`:59`), dedup (`:134`), dispatch (`:172`)
- `control-plane/internal/services/trigger_dispatcher.go:56-186` — `DispatchEvent` (`:56`) POSTs direct to `{BaseURL}/reasoners/{TargetReasoner}` (`:137`)
- `control-plane/internal/services/source_manager.go:26-212` — loop-source goroutine lifecycle
- `control-plane/internal/services/webhook_dispatcher.go:113-469` — completion webhooks, HMAC, backoff, restart-safe poller
- `control-plane/internal/handlers/retry.go:16-29` — `isRetryableDBError` + `backoffDelay` for DB contention
- `control-plane/internal/handlers/execute_cancel_tree.go:62-254` — leaf-first workflow tree cancel
- `control-plane/internal/services/cancel_dispatcher.go:107-232` — bus→SDK `/_internal` cancel bridge
- `sdk/python/agentfield/cancel.py:65-115` — `install_cancel_route` → `asyncio.Task.cancel()`
- `sdk/python/agentfield/agent_pause.py:10-95` — `_PauseManager` approval futures + `PauseClock`
- `sdk/python/agentfield/agent.py:3710-4097` — `Agent.call()` routes through control plane, multi-hop wait propagation
- `sdk/python/agentfield/client.py:774-889` — `_submit_execution_async` + header injection (`X-Run-ID`/`X-Caller-Agent-ID`)
- `control-plane/internal/handlers/connector/handlers.go:52-157` — `RegisterRoutes` (`:52`), `GetManifest` (`:136`), capability-gated connector API
- `control-plane/internal/encryption/encryption.go:19-189` — AES-256-GCM + PBKDF2 600k master-seed encryption (`AFENC2` `:19`, `600000` `:22`)
- `control-plane/internal/services/did_service.go:41-585` — master seed, HKDF derivation (salt `:563`, `hkdf.New` `:566`), did:key `0xed01` encoding (`:580`)
- `control-plane/internal/services/vc_issuance.go:16-276` — `GenerateExecutionVC` (`:16`) + `Ed25519Signature2020` proof (`:80`) + `signVC` (`:240`)
- `control-plane/internal/services/vc_chain.go:14-123` — `GetWorkflowVCChain` (`:14`) + offline `collectDIDResolutionBundle` (`:63`)
- `control-plane/internal/cli/vc.go:36-197` — `NewVerifyAliasCommand` (`:36`) → `verifyVC` (`:175`); `verify_provenance.go:17-196` offline VC verification
- `control-plane/internal/services/access_policy_service.go:26-424` — tag-based cross-agent access evaluation
- `control-plane/internal/services/tag_approval_service.go:400-567` — tag approval lifecycle + signed `AgentTagCredential`
- `control-plane/internal/handlers/execute_approval.go:52-214` — `RequestApprovalHandler` human approval gate (execution→waiting)
- `control-plane/internal/handlers/webhook_approval.go:104-537` — `normalizeDecision` (`:104`), `ApprovalWebhookHandler` (`:126`), HMAC `verifySignature` (`:432`)
- `control-plane/internal/handlers/execution_guards.go:40-127` — `CheckExecutionPreconditions` (`:40`), `checkLLMEndpointHealth` (`:62`), `ErrorCategoryConcurrencyLimit` (`:127`)
- `control-plane/internal/services/health_monitor.go:47-461` — agent health loop, `CheckInterval` 10s default (`:67`), heartbeat-stale gating
- `control-plane/internal/services/llm_health_monitor.go:19-343` — 3-state LLM circuit breaker
- `control-plane/internal/services/observability_forwarder.go:59-741` — batching event forwarder, `sendBatch` (`:636`), DLQ (`:693`), `Redrive` (`:249`)
- `control-plane/internal/observability/execution_tracer.go:17-143` — bus→OTel span bridge
- `control-plane/internal/handlers/memory.go:402-434` — `resolveScope` (`:402`) + `getScopeID` (`:421`) scope resolution & hierarchical lookup
- `sdk/python/agentfield/execution_context.py:27-302` — `ExecutionContext` dataclass ↔ `X-*` header marshaling (unifying thread)
- `control-plane/internal/mcp` — **VERIFIED ABSENT**: directory does not exist; zero `mcp` route/handler matches (documentation references it, current tree has none)

## Architecture Documentation

Current patterns observed in the codebase:

| Pattern | Where | What it does |
|---|---|---|
| Control plane as router, agents own logic | `execute.go:1010-1560` | The plane resolves `node.reasoner` → registered `BaseURL` and POSTs; it never executes agent code itself. |
| `ExecutionContext` as `X-*` header envelope | `execution_context.py:27-302`, `execute.go:1242-1252` | One dataclass drives routing, DAG edges, memory scoping, and DID/VC issuance; children echo headers to self-correlate. |
| In-process event bus as spine | `events/execution_events.go:114` | A single `GlobalExecutionEventBus` fans lifecycle events to sync-waiters, tracer, forwarder, and SSE. |
| Source plugin registry | `sources/source.go:79-110`, `sources/all/all.go` | `HTTPSource`/`LoopSource` implementations self-register via `init()`; the gateway dispatches by name. |
| Deterministic key derivation from one seed | `did_service.go:562-585` | HKDF over an encrypted master seed yields every agent/reasoner/skill `did:key` without per-identity key storage. |
| Self-verifying audit bundle | `vc_chain.go:14-123`, `verify_provenance.go:17-196` | The VC chain embeds a `DIDResolutionBundle` so `af verify` checks signatures fully offline. |
| Unified preconditions gate | `execution_guards.go:40-127` | Concurrency (429), LLM circuit (503), and access policy are all enforced at one chokepoint before dispatch. |
| Single struct, two storage backends | `local.go:469-767`, `storage.go:789-835` | One `LocalStorage` implements both SQLite+BoltDB and PostgreSQL, selected by config. |
| Async queue with backpressure | `execute.go:2215-2230` | Bounded worker pool returns 202 fast and sheds with 503 when saturated. |

## Historical Context (from thoughts/)

No prior `thoughts/` research documents on this topic were located during this pass.

## Related Research

None found in `thoughts/searchable/shared/research/` at the time of writing. This is the first research document on AgentField's "AI Backend" claim.

## Open Questions

None blocking the answer. The one documentation/code divergence found — `internal/mcp/` referenced in `CLAUDE.md` but absent from the tree — is recorded as a verified fact in the findings, not an unresolved question.
