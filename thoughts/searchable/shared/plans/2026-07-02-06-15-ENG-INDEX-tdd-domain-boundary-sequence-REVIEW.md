# Plan Review Report: ENG-INDEX TDD Domain Boundary Sequence

Reviewed plan: `agentfield/thoughts/searchable/shared/plans/2026-07-02-06-15-ENG-INDEX-tdd-domain-boundary-sequence.md`

Review date: 2026-07-02

Decision: Needs Major Revision

## Review Summary

| Category | Status | Issues Found |
| --- | --- | ---: |
| Contracts | [CRITICAL] | 5 |
| Interfaces | [CRITICAL] | 5 |
| Promises | [CRITICAL] | 4 |
| Data Models | [CRITICAL] | 4 |
| APIs | [CRITICAL] | 3 |
| CodeCleanup gates | [WARN] | 5 |

The plan is directionally sound: it targets real boundary leaks around execution events, aggregate writes, VC issuance, callback discovery, and SDK helpers. It is not implementation-ready because it over-specifies new seams while under-specifying the durable contracts, migration shape, and compatibility behavior those seams need to preserve.

## Critical Findings

### 1. Durable execution event sequencing is not defined

The plan proposes `ExecutionEventPublisher.PublishExecutionEvent(...) (sequence int64, err error)` and `SubscribeExecutionEvents(...) <-chan ExecutionEvent` at plan lines 153-158, but it does not define sequence scope, allocation, replay, acknowledgements, or the mapping from current event stores into the new publisher.

Current evidence:

- `control-plane/internal/events/execution_events.go:78` publishes in-memory and drops on full subscriber channels.
- `control-plane/internal/events/event_bus.go:41` has the same drop-on-full behavior for generic event buses.
- `control-plane/internal/storage/storage.go:150` exposes `StoreWorkflowExecutionEvent`, but only for `types.WorkflowExecutionEvent`.
- `control-plane/internal/storage/local.go:7617` stores workflow execution events after callers supply `Sequence`; there is no allocator.
- `control-plane/migrations/013_workflow_execution_state.sql:19` defines `workflow_execution_events` as per-`execution_id`, with `UNIQUE(execution_id, sequence)`, not as a global execution outbox.
- Existing consumers rely on unsubscribe-capable buses through `GetExecutionEventBus()` and `GetWorkflowExecutionEventBus()` in `control-plane/internal/storage/storage.go:438`.

Why this blocks implementation:

The plan promises "durable sequence-bearing execution events", SSE correctness, and VC idempotency, but the proposed interface is still an in-memory subscription surface. Without a cursor/replay contract and a sequence allocator, implementers could satisfy the type shape while retaining dropped events, duplicate sequence `0`, or inconsistent per-execution ordering.

Required amendment:

Define the durable event contract before implementation. At minimum, specify event table ownership, sequence scope, allocation rule, replay API such as `ListExecutionEvents(afterSequence, limit)`, subscriber cancellation/unsubscribe semantics, publish-after-commit ordering, and whether existing `WorkflowExecutionEvent` records are reused or replaced.

### 2. Aggregate write boundary does not cover the real transactional problem

The plan proposes `ExecutionAggregateStore.CreateExecutionAggregate` and `UpdateExecutionAggregate` at plan lines 249-252, but it does not define ownership of `executions` versus `workflow_executions`, return values, transactional support, or error behavior across the existing writers.

Current evidence:

- `control-plane/internal/storage/execution_records.go:18` and `control-plane/internal/storage/local.go:2010` write the two tables independently.
- `control-plane/internal/handlers/execute.go:1317` updates the primary execution row, then `control-plane/internal/handlers/execute.go:1346` updates workflow execution state separately.
- `control-plane/internal/handlers/execute.go:1899` logs workflow final-state update failures rather than failing the enclosing operation.
- `control-plane/internal/handlers/execute_cancel.go:73` updates `executions`; `control-plane/internal/handlers/execute_cancel.go:101` treats missing workflow execution rows as non-fatal.
- `control-plane/internal/handlers/execute_awaiter_status.go:141` and `control-plane/internal/handlers/workflow_execution_events.go:35` are additional mutation paths.
- `control-plane/internal/storage/local.go:7612` leaves `BeginTransaction` unimplemented, while `control-plane/internal/storage/unitofwork.go:30` does not include primary execution records.

Why this blocks implementation:

The plan wants all state transitions to become ordered domain operations, but the interface only accepts a mutator and returns `error`. Existing handlers need updated records for events, webhooks, responses, and idempotency decisions. The plan also does not say whether "unknown current state can initialize to running" is intended, while `control-plane/internal/storage/execution_state_validation.go:26` currently allows `unknown -> pending` only.

Required amendment:

Define the aggregate ownership contract, transaction mechanism, return shape, and full writer inventory. The plan should explicitly say which paths must use the aggregate store, which paths remain compatibility shims, and how state validation errors map to existing HTTP responses and logs.

### 3. Event-driven VC issuance cannot be made exactly-once from the current data

The plan says the VC issuer should build `types.ExecutionContext` from the stored execution at plan line 343 and use idempotency checks at plan line 345. That is not sufficient with the current model.

Current evidence:

- `control-plane/pkg/types/did_types.go:153` defines `ExecutionContext` with DID and parent-VC fields: `CallerDID`, `TargetDID`, `AgentNodeDID`, and `ParentVCID`.
- Current execution/workflow rows do not persist a complete `ExecutionContext` equivalent for reconstruction.
- `control-plane/internal/services/vc_issuance.go:106` creates a fresh VC ID for each generation.
- `control-plane/internal/storage/local.go:7140` and `control-plane/internal/storage/local.go:7167` upsert execution VC records by `vc_id`, not by the logical uniqueness key.
- `control-plane/migrations/004_create_execution_vcs.sql:39` has a unique index on `(execution_id, issuer_did, target_did)`, but the plan does not specify conflict-as-success or concurrent issuance handling.
- Direct SDK issuance still exists in Python and Go, for example `sdk/python/agentfield/decorators.py:333`, `sdk/python/agentfield/agent.py:1564`, and `sdk/go/agent/agent_did.go:94`.
- The public manual endpoint `POST /api/v1/execution/vc` exists in `control-plane/internal/handlers/did_handlers.go:260`.

Why this blocks implementation:

The plan can remove decorator-level VC calls and still leave direct SDK generation paths producing duplicates or missing credentials. It also cannot reliably issue server-side VCs if the event does not carry or reference all DID context needed by `GenerateExecutionVC`.

Required amendment:

Define the VC source-of-truth event payload, persistence key, and conflict handling. Include all SDK languages and the manual VC endpoint in the compatibility matrix. Specify whether SDK direct issuance is deprecated, disabled, routed to the server, or preserved as a separate feature.

### 4. Callback probing needs a security and compatibility contract

The plan adds probing for callback candidates at plan lines 429-434, but it does not define SSRF controls, redirect handling, candidate limits, timeout constants, or response compatibility. It also references a non-existent `Reachable` field in tests at plan lines 415-416.

Current evidence:

- `control-plane/pkg/types/types.go:208` defines `CallbackTestResult.Success`; there is no `Reachable`.
- `control-plane/internal/handlers/nodes_register.go:91` reads `discovery.Candidates`, not the legacy Python key `callback_candidates`.
- `sdk/python/agentfield/agent.py:1603` still sends `callback_candidates`.
- Existing callback URL validation at `control-plane/internal/handlers/nodes_register.go:37` checks basic scheme/host shape only.
- The stricter serverless path already has safer URL parsing and host allowlist handling in `control-plane/internal/handlers/nodes_register.go:235` and `control-plane/internal/config/config.go:65`.
- Existing test `control-plane/internal/handlers/nodes_discovery_test.go:25` expects callback candidate resolution to return nil probe results because probing is currently disabled.

Why this blocks implementation:

Callback probing will make outbound network requests during registration. Without a bounded and allowlisted probing contract, implementers can create SSRF risk or slow registration behavior. Without DTO compatibility, the existing Python SDK can keep sending candidates that the server ignores.

Required amendment:

Define the callback prober contract in terms of allowed schemes, host allowlists, redirect policy, timeout, max candidates, probe order, fallback paths, result DTO field names, and legacy payload decoding. Add tests for `callback_candidates`, `candidates`, `Success`, probe failure, invalid URLs, and blocked hosts.

### 5. SDK decomposition tests are too shallow for the promised behavior

The plan correctly identifies `Agent.call`, `Agent.note`, and `Agent.pause` as oversized public methods, but the proposed tests do not yet cover the public behavior that can regress.

Current evidence:

- `sdk/python/agentfield/agent.py:4238` implements `Agent.note` with nested sync/async HTTP paths, direct `aiohttp` and `requests` usage, and fire-and-forget scheduling.
- `sdk/python/agentfield/client.py` has no `post_note` method today.
- `sdk/python/agentfield/agent.py:4417` implements `pause`; timeout/default behavior also appears in `wait_for_resume` at `sdk/python/agentfield/agent.py:4559`.
- The plan test at lines 499-507 only asserts that `Agent.note` delegates to `client.post_note`.

Why this matters:

Moving transport into `AgentFieldClient.post_note` is good, but a delegation-only test can pass while dropping sync/async behavior, context resolution, timeout handling, error handling, or task scheduling semantics.

Required amendment:

Add behavior tests for note context selection, synchronous and asynchronous call paths, timeout/error mapping, payload shape, and task completion. Keep the refactor flat: public SDK methods should gather context and delegate to typed client methods, not retain nested transport fallbacks.

## Category Review

### Contracts

Status: [CRITICAL]

What is well-defined:

- The plan has a clear phased structure and small behavior slices.
- It correctly names the main boundary leaks: global execution buses, direct paired table writes, decorator-issued VCs, callback discovery normalization, and SDK public method size.

Missing or unclear:

- Durable execution event contract: sequence scope, allocation, replay, unsubscribe, and retention.
- Aggregate state transition contract: table ownership, all writer paths, transaction guarantees, and updated-record returns.
- VC idempotency contract: logical uniqueness key, concurrent conflict behavior, and manual endpoint coexistence.
- Callback probing security contract: SSRF controls, host allowlists, redirect policy, and bounded candidate probing.
- Workflow event ingestion compatibility: whether existing loose request/response behavior is preserved.

### Interfaces

Status: [CRITICAL]

What is well-defined:

- The proposed seams point at the right ownership boundaries.
- The plan intends to replace globals and direct HTTP transport with injectable domain interfaces.

Missing or unclear:

- `ExecutionEventPublisher` lacks replay/list and unsubscribe/cancel support.
- The event publisher interface does not map to `StoreWorkflowExecutionEvent` or current SSE consumers.
- `ExecutionAggregateStore` does not return the updated aggregate state needed by handlers.
- Bootstrap and dependency injection updates are incomplete for event forwarders, tracers, handlers, and SDK compatibility paths.
- `CallbackProber` needs a typed result contract aligned with `CallbackTestResult.Success`.

### Promises

Status: [CRITICAL]

Promise gaps:

- "Durable sequence-bearing events" is promised, but current storage and interfaces cannot guarantee no drops or replay.
- "All paired execution/workflow writes go through an aggregate" is promised, but the plan does not inventory all mutation paths.
- "VC issuance becomes event-driven and exactly-once" is promised, but direct SDK/manual paths and idempotency are unresolved.
- "Callback discovery probes candidates" is promised, but outbound-request safety and legacy SDK payloads are unresolved.

### Data Models

Status: [CRITICAL]

Data model gaps:

- `workflow_execution_events` is a workflow execution event table, not a general execution outbox.
- The plan does not say whether to introduce a new outbox table, extend the existing table, or migrate events into one stream.
- The execution aggregate has no explicit record shape that owns both `executions` and `workflow_executions`.
- Server-side VC issuance needs persisted or event-carried DID context fields that current execution rows do not provide.

### APIs

Status: [CRITICAL]

API gaps:

- SSE replay semantics and cursor parameter behavior are not defined.
- `/api/v1/workflow/executions/events` compatibility and response shape should be preserved or deliberately migrated.
- Callback registration must decode both `candidates` and legacy `callback_candidates`.
- Status/cancel/pause/resume endpoints need explicit error-shape compatibility when moved through the aggregate store.
- The public VC endpoint must be included in the VC migration contract.

### CodeCleanup Plan-Hygiene Gates

Status: [WARN]

Loaded gates:

- No side effects in conditionals.
- No mutation in control expressions.
- Never nesting.
- Named constants over literals.
- Control-expression discipline.
- Maintainability recovery.

Plan-hygiene concerns:

- Callback URL normalization should remain pure; probing I/O should happen in a separate bounded step.
- Probe constants such as `/health`, root fallback, timeout, max candidates, and redirect policy need names/config rather than inline literals.
- Aggregate state validation should remain a pure, table- or switch-backed decision with explicit precedence.
- VC scheduling removal should compute eligibility booleans before side-effect calls and remove all direct scheduling paths serially.
- `Agent.note()` should become a flat public adapter around `AgentFieldClient.post_note`; direct nested HTTP fallbacks should not remain in the public method.

## Required Plan Amendments Before Implementation

1. Add a formal event-outbox contract section covering table shape, sequence allocation, replay, retention, publish-after-commit, subscriber cancellation, and migration from current buses.
2. Add an execution aggregate contract section covering record shape, transaction support, updated-record return values, error mapping, and the full inventory of writers to migrate.
3. Add a VC issuance contract section covering event payload context, idempotency key, conflict-as-success behavior, manual endpoint behavior, and all SDK direct issuance paths.
4. Add a callback probing security section covering allowlists, private/loopback handling, redirect policy, timeout, max candidates, fallback paths, DTO field names, and legacy payload compatibility.
5. Expand SDK tests beyond delegation to cover context, sync/async execution, payload shape, timeouts, and error behavior.
6. Update the test plan so every new contract has a red test before implementation starts.

## Approval

Not approved for implementation as written.

The next revision should preserve the plan's phase structure, but each phase needs explicit contracts and compatibility obligations before implementation begins.
