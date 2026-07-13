# Domain Boundary Refactor Sequence TDD Implementation Plan

## Overview

This plan implements the `specs/INDEX.md` recommended sequence as five behavior-preserving TDD moves. It incorporates the findings from `2026-07-02-06-15-ENG-INDEX-tdd-domain-boundary-sequence-REVIEW.md`; the review blockers are treated as required contracts, not follow-up cleanup.

1. Establish one injected, durable execution event bus with typed events and sequence numbers.
2. Collapse execution writes onto one transactional write path guarded by a single transition validator.
3. Move VC issuance behind `ExecutionCompleted` event consumption.
4. Add active callback reachability probing during agent registration.
5. Decompose the Python SDK `Agent` god object behind the now-stable event boundary.

The order matters. Event durability and injection come first because Execution, Identity, Registry, SDK, UI SSE, cancel dispatch, tracing, and webhooks all depend on execution events. Each phase must preserve existing public API behavior before the next phase begins.

## Current State Analysis

### Key Discoveries

- `EventBus[T]` already exists but delivery is process-local and drop-on-full: `agentfield/control-plane/internal/events/event_bus.go:6`, `agentfield/control-plane/internal/events/event_bus.go:47`.
- `ExecutionEventBus` duplicates the generic bus and has a package global singleton: `agentfield/control-plane/internal/events/execution_events.go:42`, `agentfield/control-plane/internal/events/execution_events.go:114`.
- Node and reasoner events duplicate the same bus shape, with node dedup embedded in event-bus code: `agentfield/control-plane/internal/events/node_events.go:50`, `agentfield/control-plane/internal/events/node_events.go:249`, `agentfield/control-plane/internal/events/reasoner_events.go:31`.
- Storage owns a separate execution bus instance, creating split-brain delivery: `agentfield/control-plane/internal/storage/local.go:482`, `agentfield/control-plane/internal/storage/local.go:6269`.
- `execute.go` publishes completion events to both the controller bus and the global bus: `agentfield/control-plane/internal/handlers/execute.go:826`, `agentfield/control-plane/internal/handlers/execute.go:828`.
- Observability consumers currently subscribe to global event buses, not the storage-owned bus used by sync waiters: `agentfield/control-plane/internal/observability/execution_tracer.go:69`, `agentfield/control-plane/internal/services/observability_forwarder.go:371`.
- Execution notes publish an untyped `workflow_note_added` string through the storage-owned bus only: `agentfield/control-plane/internal/handlers/execution_notes.go:158`.
- `workflow_execution_events.sequence` and `previous_sequence` already exist, but the table is scoped to `workflow_executions` and has `UNIQUE(execution_id, sequence)`, so it is not the general execution outbox required by this plan: `agentfield/control-plane/migrations/013_workflow_execution_state.sql:19`, `agentfield/control-plane/pkg/types/types.go:702`.
- `StoreWorkflowExecutionEvent` persists caller-supplied sequence values and publishes after commit, which is useful as a compatibility pattern but not a sequence allocator: `agentfield/control-plane/internal/storage/local.go:7617`, `agentfield/control-plane/internal/storage/local.go:7640`, `agentfield/control-plane/internal/storage/local.go:7655`.
- `/api/v1/workflow/executions/events` mutates execution rows but does not append `WorkflowExecutionEvent` or publish a bus event today: `agentfield/control-plane/internal/handlers/workflow_execution_events.go:35`, `agentfield/control-plane/internal/handlers/workflow_execution_events.go:70`.
- Sync and async execute routes share `prepareExecution`, but route registration still exposes separate public endpoints: `agentfield/control-plane/internal/server/routes_core.go:91`, `agentfield/control-plane/internal/server/routes_core.go:110`.
- `prepareExecution` writes `executions` first, then `workflow_executions` separately: `agentfield/control-plane/internal/handlers/execute.go:1151`, `agentfield/control-plane/internal/handlers/execute.go:1182`.
- `CreateExecutionRecord` and `StoreWorkflowExecution` run in separate storage transactions: `agentfield/control-plane/internal/storage/execution_records.go:18`, `agentfield/control-plane/internal/storage/local.go:2010`.
- Completion and failure paths update `executions`, then separately update `workflow_executions`: `agentfield/control-plane/internal/handlers/execute.go:1304`, `agentfield/control-plane/internal/handlers/execute.go:1377`, `agentfield/control-plane/internal/handlers/execute.go:1899`.
- Status transition guards are duplicated in handlers: `agentfield/control-plane/internal/handlers/execute.go:516`, `agentfield/control-plane/internal/handlers/execute.go:537`.
- A workflow-only transition validator already exists but does not protect `UpdateExecutionRecord`: `agentfield/control-plane/internal/storage/execution_state_validation.go:21`, `agentfield/control-plane/internal/storage/local.go:2160`.
- VC generation is currently exposed as inline HTTP `POST /api/v1/execution/vc`: `agentfield/control-plane/internal/handlers/did_handlers.go:260`, `agentfield/control-plane/internal/handlers/did_handlers.go:327`.
- SDK reasoner tracking can fire-and-forget VC generation after success or error: `agentfield/sdk/python/agentfield/decorators.py:333`, `agentfield/sdk/python/agentfield/decorators.py:347`, `agentfield/sdk/python/agentfield/decorators.py:429`.
- Python `Agent` and Go SDK also have direct fire-and-forget VC paths: `agentfield/sdk/python/agentfield/agent.py:2260`, `agentfield/sdk/python/agentfield/agent.py:2894`, `agentfield/sdk/go/agent/agent_did.go:96`.
- TypeScript appears to expose manual DID credential generation rather than an automatic completion hook: `agentfield/sdk/typescript/src/did/DidInterface.ts:33`, `agentfield/sdk/typescript/src/did/DidClient.ts:190`.
- `types.ExecutionContext` requires DID and parent VC fields that current execution/workflow rows do not fully persist, so event-driven VC issuance must carry or reference that context at event time: `agentfield/control-plane/pkg/types/did_types.go:153`.
- Execution VC storage has a logical uniqueness index on `(execution_id, issuer_did, target_did)`, but current upserts are by `vc_id`, so pre-check-only idempotency is not sufficient under duplicate event delivery: `agentfield/control-plane/migrations/004_create_execution_vcs.sql:39`, `agentfield/control-plane/internal/storage/local.go:7140`, `agentfield/control-plane/internal/storage/local.go:7167`.
- Registry callback discovery normalizes candidates but returns no probe results today: `agentfield/control-plane/internal/handlers/nodes_register.go:187`, `agentfield/control-plane/internal/handlers/nodes_discovery_test.go:26`.
- `CallbackTestResult` already exists in the public type model with `Success`, not `Reachable`, so tests and response contracts must use `success`: `agentfield/control-plane/pkg/types/types.go:208`.
- Python sends `callback_candidates` in discovery payloads while Go decodes `candidates`, so callback probing must cover both keys for compatibility: `agentfield/sdk/python/agentfield/agent.py:1609`, `agentfield/control-plane/pkg/types/types.go:203`.
- Serverless registration already has an active `/discover` probe pattern that can inform long-running registration: `agentfield/control-plane/internal/handlers/nodes_register.go:796`, `agentfield/control-plane/internal/handlers/nodes_register.go:829`.
- The SDK has callback candidate detection tests already: `agentfield/sdk/python/tests/test_agent_networking.py`.
- Python `Agent` is the large composition root and still owns call, note, pause, context, VC, and invocation behavior: `agentfield/sdk/python/agentfield/agent.py:489`, `agentfield/sdk/python/agentfield/agent.py:3758`, `agentfield/sdk/python/agentfield/agent.py:4238`, `agentfield/sdk/python/agentfield/agent.py:4417`.
- Pause coordination is already partially extracted, but duplicate `_PauseManager` implementations remain: `agentfield/sdk/python/agentfield/agent.py:377`, `agentfield/sdk/python/agentfield/agent_pause.py:49`.
- Context is mirrored across the agent, client, decorators, and workflow helper: `agentfield/sdk/python/agentfield/agent.py:691`, `agentfield/sdk/python/agentfield/client.py:169`, `agentfield/sdk/python/agentfield/decorators.py:289`, `agentfield/sdk/python/agentfield/agent_workflow.py:58`.

### Existing Test Patterns

- Event publisher tests: `agentfield/control-plane/internal/events/publishers_test.go`.
- Generic bus tests, including ordering and slow subscriber behavior: `agentfield/control-plane/internal/events/event_bus_test.go`.
- Execution note handler/event tests: `agentfield/control-plane/internal/handlers/execution_notes_test.go`.
- Persisted workflow/log sequence tests: `agentfield/control-plane/internal/storage/coverage_identity_events_reasoner_test.go`, `agentfield/control-plane/internal/handlers/execution_logs_test.go`.
- Status normalization tests: `agentfield/control-plane/pkg/types/status_test.go`.
- Status update and terminal regression tests: `agentfield/control-plane/internal/handlers/execute_status_update_test.go`.
- Sync/async execute tests: `agentfield/control-plane/internal/handlers/execute_async_test.go`, `agentfield/control-plane/internal/handlers/execute_handler_test.go`.
- Workflow state machine invariant tests: `agentfield/control-plane/internal/storage/execution_state_invariant_test.go`.
- Callback discovery tests: `agentfield/control-plane/internal/handlers/nodes_discovery_test.go`.
- SDK networking/callback tests: `agentfield/sdk/python/tests/test_agent_networking.py`.
- SDK decorator and workflow tracking tests: `agentfield/sdk/python/tests/test_decorators.py`.
- SDK `app.call()` tests: `agentfield/sdk/python/tests/test_agent_call.py`.
- SDK pause/async callback tests: `agentfield/sdk/python/tests/test_async_execution.py`, `agentfield/sdk/python/tests/test_agent_coverage_additions.py`.

## Desired End State

The platform has one execution event backbone that producers receive by injection, stores events durably in a global execution outbox with monotonically increasing sequence numbers, and delivers typed event DTOs to existing subscribers. Execution state changes go through one validator and one transactional write path. Identity consumes `ExecutionCompleted` to issue VCs with retry and atomic idempotency, while SDK reasoner code no longer calls the VC endpoint in fire-and-forget tasks. Registry only marks a node live after selecting a callback URL that is reachable under a bounded SSRF-safe probing policy. SDK decomposition happens after the event boundary is stable, preserving the developer-facing `Agent` API.

### Observable Behaviors

- Given an execution note is posted, when the handler appends the note, then subscribers on the injected execution bus receive a typed `ExecutionNoteAdded` event and the durable event log records it with a sequence number.
- Given a slow subscriber is present, when a critical execution event is published, then publishing remains non-blocking and the event is still recoverable from the durable outbox.
- Given an execution status update attempts terminal-to-non-terminal regression, when any write path applies it, then the same validator rejects it and persisted state is unchanged.
- Given a sync or async execution starts, when the initial write succeeds, then `executions` and `workflow_executions` reflect the same execution state atomically or not at all.
- Given an execution completes, when an `ExecutionCompleted` event is committed, then the VC issuer consumes the event and stores exactly one execution VC for that execution.
- Given SDK reasoner tracking completes, when VC issuance is event-driven, then decorators do not call `_generate_vc_async` or create a VC task.
- Given multiple callback candidates, when registration runs auto-discovery, then the first reachable candidate is selected and probe results are included in `callback_discovery.tests`.
- Given `app.call()`, `app.note()`, and `app.pause()` are used after SDK extraction, when tests exercise existing APIs, then behavior and signatures remain backward compatible.

## What We're NOT Doing

- No external message broker. The durable bus starts with an in-process bus plus storage-backed outbox/append-log.
- No physical table merge of `executions` and `workflow_executions` in this plan. End the split-brain by transactional write coordination and explicit projection ownership, not by a risky migration.
- No DID key custody redesign in this sequence. The VC move removes hot-path coupling first; local SDK signing remains a later Identity phase.
- No public route removals. Existing endpoints must keep working.
- No silent public API redefinition. Existing route paths, JSON response shapes, auth behavior, and SSE `from` semantics must be preserved unless a behavior explicitly adds a compatibility test for a deliberate change.
- No durable consumer ack/redrive table in this sequence. Consumers must be idempotent and recover by replaying the outbox cursor; retry/backoff is owned by the consumer that needs it.
- No SDK web framework rewrite. `Agent` may remain a FastAPI subclass while services are extracted.

## Review Closure Contracts

These contracts close every critical and warning from the review. Each behavior below must add red tests for its relevant contract before implementation starts.

### Event Outbox Contract

- **Outbox owner:** Add a new narrow `execution_event_outbox` storage model for general execution events. Do not use `workflow_execution_events` as the primary outbox because it is scoped by `execution_id` and remains the workflow-event compatibility table.
- **Sequence scope:** `sequence` is a single global, monotonically increasing `int64` across all execution events. It is the durable replay cursor.
- **Allocation:** sequence allocation happens inside the same storage transaction as the outbox insert. Use a DB-enforced allocator such as an autoincrement primary key or a locked one-row sequence table; caller-supplied sequence values are prohibited for the general outbox.
- **Stored event shape:** store `sequence`, `previous_sequence`, `execution_id`, optional `workflow_id`, `event_type`, `status`, `occurred_at`, `payload_json`, `metadata_json`, and `idempotency_key`. `payload_json` must carry the typed event DTO needed by consumers, not only a free-form log string.
- **Replay API:** expose `ListExecutionEvents(ctx, afterSequence int64, limit int) ([]StoredExecutionEvent, error)` and `GetExecutionEvent(ctx, sequence int64)`.
- **Live API:** expose subscribe with cancellation/unsubscribe, for example `SubscribeExecutionEvents(ctx context.Context, subscriberID string, afterSequence int64, buffer int) (ExecutionEventSubscription, error)`. The subscription must first replay persisted events with `sequence > afterSequence`, then switch to live events.
- **Publish order:** append and commit first, then publish the stored event to the in-memory bus. A live publish failure must not roll back committed state.
- **Slow subscribers:** live channels remain non-blocking and may drop, but drops are recoverable because the subscriber has a durable `sequence` cursor.
- **Compatibility:** existing `GetExecutionEventBus`, global helper functions, and `workflow_execution_events` ingestion remain wrappers/projections during migration, but new producers must use the outbox-backed publisher.
- **SSE compatibility:** existing `/api/v1/executions/:id/events` and UI SSE response shapes remain stable. If durable outbox sequence is exposed to clients, it must be added as a new field; it must not silently change existing `id` or `from` semantics.

### Execution Aggregate Contract

- **Ownership:** `executions` remains the authoritative public execution record for existing APIs, payloads, and status. `workflow_executions` is the workflow projection that owns workflow state version, child counters, and `last_event_sequence`.
- **Aggregate shape:** define an `ExecutionAggregate` with `Execution *types.Execution`, `WorkflowExecution *types.WorkflowExecution`, and `CommittedEvents []events.StoredExecutionEvent`.
- **Write interface:** aggregate commands must return the updated aggregate, not only `error`, so handlers can build existing responses, webhooks, and events without reloading stale state.
- **Transaction:** `CreateExecutionAggregate` and `UpdateExecutionAggregate` must write `executions`, `workflow_executions`, and the required outbox event in one transaction. If workflow projection or event append fails, the primary execution row is rolled back.
- **Local storage:** do not rely on the currently unimplemented `LocalStorage.BeginTransaction` unless this behavior implements it. A private local aggregate transaction method is acceptable if it is fully tested.
- **Validation:** `types.ValidateExecutionTransition(from, to)` is the only transition decision point. It must be pure, table- or switch-backed, and must not perform storage or event side effects.
- **Legacy unknown state:** preserve current behavior: `unknown -> pending` is allowed. `unknown -> running` is not allowed directly unless the aggregate creates a missing legacy projection as `pending` in the same transaction before applying `pending -> running`, with a test documenting that compatibility path.
- **Writer inventory:** migrate or wrap all execution mutation paths: sync prepare, async prepare, completion, failure, status callback, cancel, cancel tree, pause, resume, approval, awaiter status, workflow event ingestion, webhook-driven updates, and any reaper/retry code found by `rg "UpdateExecutionRecord|StoreWorkflowExecution|UpdateWorkflowExecution|StoreWorkflowExecutionEvent"`.
- **Error mapping:** public HTTP status codes and JSON shapes must be preserved for execute, async execute, status callback, cancel, pause/resume, approval, workflow event ingestion, and UI routes unless a test explicitly documents a deliberate compatibility change.

### VC Issuance Contract

- **Source event:** `ExecutionCompleted` must carry or reference enough context to build `types.ExecutionContext`: `execution_id`, `workflow_id`, `session_id`, `caller_did`, `target_did`, `agent_node_did`, optional `parent_vc_id`, timestamps, input/output payload references, status, error message, duration, and reasoner identity.
- **No row-only reconstruction:** the VC issuer must not assume current execution/workflow rows contain all DID and parent VC data. If the data is missing from the event, issuance must skip with a structured retry/error event rather than fabricate context.
- **Idempotency key:** the logical key is `(execution_id, issuer_did, target_did)`. The issuer must atomically claim or upsert this key before signing, or use deterministic VC IDs plus conflict-as-success. A check-before-generate alone is forbidden.
- **Conflict behavior:** duplicate `ExecutionCompleted` delivery or legacy manual endpoint calls must return/load the existing VC for the logical key instead of creating another VC or retrying indefinitely on a uniqueness error.
- **Manual endpoint:** `POST /api/v1/execution/vc` remains available and uses the same idempotency policy. Its request and response shape must stay compatible.
- **SDK automatic hooks:** Python decorators, Python `Agent`, and Go `maybeGenerateVC` automatic fire-and-forget paths are disabled by default or routed through the server event path. TypeScript manual DID APIs remain available but must not become an automatic completion hook without tests.
- **Consumer lifecycle:** `ExecutionCompletedVCIssuer` has explicit `Start(ctx)`, `Stop(ctx)`, retry/backoff, and subscription cursor behavior. Execution completion latency must not depend on VC generation success.

### Callback Probing Contract

- **Pure normalization:** URL normalization and candidate collection remain pure functions. Network probing happens in a separate bounded prober step.
- **Request DTO compatibility:** decode both `callback_discovery.candidates` and legacy `callback_discovery.callback_candidates`. Response DTOs use the existing `CallbackTestResult.Success` / `json:"success"` field, not `Reachable`.
- **Security:** only `http` and `https` schemes are eligible. Reject userinfo, opaque URLs, malformed hosts, query/fragment-only candidates, and metadata-address candidates. Private, loopback, and link-local addresses are allowed only in test/dev mode or when explicitly allowlisted by registration config.
- **Allowlist:** reuse or extend registration config for callback probing allowlists. Production probing must be deny-by-default for hosts outside the configured policy when an allowlist is present.
- **Redirects:** do not follow redirects by default. Record redirect responses as failed tests unless explicitly allowed by config.
- **Bounds:** define named constants/config for `CallbackProbeHealthPath`, `CallbackProbeRootPath`, `CallbackProbeTimeout`, `CallbackProbeMaxCandidates`, and total registration probe budget.
- **Probe order:** probe `/health` first, then root only for candidates that pass security checks and fail health with a non-terminal error. Select the first candidate with `Success == true`.
- **Failure policy:** invalid or blocked candidates are recorded as failed `CallbackTestResult` entries. When no candidate succeeds, preserve the existing registration policy for pending/accepted status and make that behavior explicit in tests.

### SDK Decomposition Contract

- **Public API preservation:** `Agent.call`, `Agent.note`, `Agent.pause`, decorators, and workflow helpers keep their signatures and existing observable behavior.
- **Note adapter:** add `AgentFieldClient.post_note` and route note transport through it. `Agent.note()` becomes a flat adapter: gather context, build payload, delegate, and schedule/await according to the existing sync/async behavior. It must not retain direct `aiohttp` or `requests` fallback logic.
- **Behavior coverage:** tests must cover note context selection, no-context behavior, payload shape, sync path, async path, timeout/error mapping, and task completion/fire-and-forget semantics.
- **Context:** preserve `_current_execution_context` as a compatibility facade while routing new writes through one context manager/port.
- **CodeCleanup gates:** avoid side effects in conditionals, avoid mutation in control expressions, flatten nested public methods, name all new literals/constants, and keep validators/probers pure until the side-effect boundary.

### Review Issue Closure Matrix

| Review finding | Plan closure |
| --- | --- |
| Durable event sequencing, replay, and unsubscribe undefined | Event Outbox Contract plus Behavior 1 red tests for global sequence, replay from cursor, unsubscribe, and workflow-ingestion append |
| `workflow_execution_events` mistaken for general outbox | Event Outbox Contract explicitly adds `execution_event_outbox` and keeps `workflow_execution_events` as compatibility/projection storage |
| Aggregate interface lacks transaction, ownership, return values, and writer inventory | Execution Aggregate Contract plus Behavior 2 aggregate DTO, transactional command, full writer inventory, rollback tests, and response-shape tests |
| `unknown -> running` conflicts with current validator | Execution Aggregate Contract and Behavior 2 tests preserve `unknown -> pending` only, with a documented legacy two-step compatibility path |
| VC issuer lacks full context and exactly-once semantics | VC Issuance Contract plus Behavior 3 event payload context, logical idempotency key, conflict-as-success, missing-context skip, and legacy endpoint sharing |
| SDK/direct VC paths can still duplicate event-driven issuance | VC Issuance Contract plus Behavior 3 Python decorator, Python Agent, Go SDK, and TypeScript compatibility tests |
| Callback probing lacks SSRF, bounds, and compatibility | Callback Probing Contract plus Behavior 4 security filter, allowlist, redirect policy, timeout/candidate constants, and blocked-host tests |
| Callback DTO mismatch `Reachable` vs `Success` | Callback Probing Contract and Behavior 4 tests require `CallbackTestResult.Success` and prohibit `Reachable` |
| SDK note decomposition tests too shallow | SDK Decomposition Contract plus Behavior 5 note tests for context, no-context, sync/async, payload, timeout/error, and task completion |
| CodeCleanup gates missing from plan | Review Closure Contracts and each behavior's refactor/success criteria require pure validators/probers, serial side effects, named constants, and flat SDK adapters |

## Testing Strategy

- **Go unit tests:** `go test ./internal/events ./internal/handlers ./internal/storage ./internal/services ./internal/config ./pkg/types`.
- **Go focused integration tests:** `go test -tags=integration ./internal/storage`.
- **Python SDK tests:** `cd agentfield/sdk/python && pytest tests/test_agent_networking.py tests/test_decorators.py tests/test_agent_call.py tests/test_agent_note.py tests/test_agent_workflow.py`.
- **Go SDK tests:** `cd agentfield/sdk/go && go test ./agent -run 'Test(DID|VC|MaybeGenerateVC|RegisterNode)'`.
- **TypeScript SDK tests:** `cd agentfield/sdk/typescript && CI=1 npm run test:core`.
- **Contract style tests:** Prefer tests at the existing package boundaries before modifying implementation. Add storage-level tests for durability, global sequence allocation, replay, aggregate transactions, and VC idempotency; handler-level tests for public HTTP behavior and response shapes; SDK tests for API compatibility and direct VC scheduling removal.
- **Review closure tests:** Every review blocker must have a red test before implementation begins: event replay/unsubscribe, aggregate rollback/return shape, VC conflict-as-success, callback SSRF and legacy key decoding, and SDK note behavior.
- **Mocking/setup:** Use existing `newTestExecutionStorage`, `httptest.Server`, Gin test routers, SDK `create_test_agent`, and monkeypatch patterns. Avoid live network except local `httptest` servers.

## Behavior 1: Typed Durable Injected Execution Event Bus

### Test Specification

**Given**: An execution event producer receives an injected outbox-backed publisher.
**When**: It publishes an execution note, status update, or completion event.
**Then**: the event is appended to the global execution outbox with an allocated sequence, committed, and only then delivered to live subscribers. Subscribers can recover missed events by replaying from their last durable sequence, and no producer writes to both global and storage-owned buses.

**Edge Cases**:
- Slow subscriber channel is full.
- Subscriber attaches after event publication and replays from `afterSequence`.
- Subscriber cancels or unsubscribes without leaking the bus registration.
- Existing publisher helpers still produce the same event payload shape.
- Note event uses a typed `ExecutionNoteAdded` constant, not raw string.
- Workflow execution event ingestion appends a durable event and publishes to the same execution boundary.
- `workflow_execution_events` remains a workflow compatibility table and is not used as the general outbox.
- Existing SSE `from` and event `id` behavior remains compatible while the durable outbox sequence is available for internal catch-up.

### TDD Cycle

#### Red: Write Failing Tests

**Files**:
- `agentfield/control-plane/internal/events/event_bus_test.go`
- `agentfield/control-plane/internal/events/execution_events_test.go`
- `agentfield/control-plane/internal/handlers/execution_notes_test.go`
- `agentfield/control-plane/internal/handlers/workflow_execution_events_test.go`
- `agentfield/control-plane/internal/handlers/reasoner_catalog_test.go`
- `agentfield/control-plane/internal/storage/execution_event_outbox_test.go`

```go
func TestExecutionEventOutboxAssignsGlobalMonotonicSequence(t *testing.T) {
    store := newTestExecutionStorage(nil)
    publisher := newTestExecutionEventPublisher(store)
    first := events.ExecutionEvent{Type: events.ExecutionStarted, ExecutionID: "exec-1"}
    second := events.ExecutionEvent{Type: events.ExecutionCompleted, ExecutionID: "exec-2"}

    stored1, err := publisher.PublishExecutionEvent(context.Background(), first)
    require.NoError(t, err)
    stored2, err := publisher.PublishExecutionEvent(context.Background(), second)
    require.NoError(t, err)

    require.Greater(t, stored2.Sequence, stored1.Sequence)
}

func TestAddExecutionNoteHandlerPublishesTypedNoteEvent(t *testing.T) {
    // Extend existing TestAddExecutionNoteHandler_AppendsNoteAndPublishesEvent.
    // Assert evt.Type == events.ExecutionNoteAdded instead of "workflow_note_added".
}

func TestExecutionEventSubscriptionReplaysFromCursorAndUnsubscribes(t *testing.T) {
    store := newTestExecutionStorage(nil)
    publisher := newTestExecutionEventPublisher(store)
    stored, err := publisher.PublishExecutionEvent(context.Background(), events.ExecutionEvent{
        Type: events.ExecutionStarted, ExecutionID: "exec-1",
    })
    require.NoError(t, err)

    sub, err := publisher.SubscribeExecutionEvents(context.Background(), "test-sub", stored.Sequence-1, 8)
    require.NoError(t, err)
    require.Equal(t, stored.Sequence, (<-sub.Events()).Sequence)

    require.NoError(t, sub.Unsubscribe())
    // Assert later publishes do not block and the subscriber is removed.
}

func TestWorkflowExecutionEventIngestionPreservesResponseAndAppendsOutbox(t *testing.T) {
    // POST /api/v1/workflow/executions/events keeps {success, created|updated}
    // while appending one execution outbox event with a generated sequence.
}
```

#### Green: Minimal Implementation

**Files**:
- `agentfield/control-plane/internal/events/event_bus.go`
- `agentfield/control-plane/internal/events/execution_events.go`
- `agentfield/control-plane/internal/events/execution_outbox.go`
- `agentfield/control-plane/internal/storage/storage.go`
- `agentfield/control-plane/internal/storage/local.go`
- `agentfield/control-plane/migrations/<next>_execution_event_outbox.sql`
- `agentfield/control-plane/internal/handlers/execution_notes.go`
- `agentfield/control-plane/internal/handlers/workflow_execution_events.go`

Minimal steps:
- Add `ExecutionNoteAdded ExecutionEventType = "execution_note_added"`.
- Add stored-event and subscription interfaces in `events` or `storage`:

```go
type StoredExecutionEvent struct {
    Sequence         int64
    PreviousSequence int64
    Event            ExecutionEvent
    Payload          json.RawMessage
    Metadata         json.RawMessage
    OccurredAt       time.Time
}

type ExecutionEventPublisher interface {
    PublishExecutionEvent(ctx context.Context, event ExecutionEvent) (StoredExecutionEvent, error)
    ListExecutionEvents(ctx context.Context, afterSequence int64, limit int) ([]StoredExecutionEvent, error)
    SubscribeExecutionEvents(ctx context.Context, subscriberID string, afterSequence int64, buffer int) (ExecutionEventSubscription, error)
}

type ExecutionEventSubscription interface {
    Events() <-chan StoredExecutionEvent
    Unsubscribe() error
    Err() error
}
```

- Add the new `execution_event_outbox` table or storage model required by the Event Outbox Contract. Do not make caller code set `sequence`.
- Back `PublishExecutionEvent` with an outbox append transaction, then publish the committed `StoredExecutionEvent` to the in-memory bus.
- Change `AddExecutionNoteHandler` to publish `ExecutionNoteAdded`.
- Change `/workflow/executions/events` ingestion to append through the same durable publisher instead of only mutating rows.
- Make storage expose one bus instance only. During the transition, keep global helper functions as compatibility wrappers that delegate to the injected default bus configured at process bootstrap.
- Keep existing SSE response shapes and `from` semantics. If an SSE handler uses outbox catch-up internally, map durable events back to the existing stream contract.

#### Refactor: Improve Code

**Files**:
- `agentfield/control-plane/internal/events/execution_events.go`
- `agentfield/control-plane/internal/services/cancel_dispatcher.go`
- `agentfield/control-plane/internal/services/observability_forwarder.go`
- `agentfield/control-plane/internal/observability/execution_tracer.go`
- `agentfield/control-plane/internal/handlers/execute.go`

Refactor steps:
- Replace direct `events.GlobalExecutionEventBus` subscriptions in cancel/observability with injected bus dependencies.
- Collapse `ExecutionEventBus` onto `EventBus[ExecutionEvent]`.
- Move sequence/replay behavior behind an outbox adapter so in-memory event subscribers remain non-blocking.
- Keep network, storage, and bus side effects out of control expressions. Allocate/store/publish in serial statements with explicit error handling.
- Add named constants for default subscription buffer, replay limit, and any outbox retention limit.
- Add a migration for `execution_event_outbox`.

### Success Criteria

**Automated:**
- [ ] Red tests fail for missing `ExecutionNoteAdded`, append/replay, or injection.
- [ ] Red tests fail if `workflow_execution_events` is used as the general outbox or caller-supplied sequence values are accepted.
- [ ] `cd agentfield/control-plane && go test ./internal/events ./internal/handlers -run 'Test(ExecutionEventOutbox|ExecutionEventSubscription|WorkflowExecutionEventIngestion|AddExecutionNoteHandler|ExecutionPublishers|ExecutionSSE)'`.
- [ ] `cd agentfield/control-plane && go test ./internal/services -run 'Test(CancelDispatcher|ObservabilityForwarder)'`.
- [ ] `cd agentfield/control-plane && go test ./internal/storage -run 'TestExecutionEventOutbox'`.

**Manual:**
- [ ] Grep confirms no producer publishes the same execution event to both controller and global buses.
- [ ] Grep confirms new producers call `PublishExecutionEvent` and do not call `StoreWorkflowExecutionEvent` as a general outbox substitute.
- [ ] Existing SSE behavior still receives execution updates.

## Behavior 2: Single Execution Transition Validator and Transactional Write Path

### Test Specification

**Given**: Any handler or service updates execution status.
**When**: The update transitions from one canonical status to another.
**Then**: `types.ValidateExecutionTransition(from, to)` decides validity for every writer, terminal states remain final, waiting transitions remain constrained, and `executions` plus `workflow_executions` update in one transaction or not at all.

**Edge Cases**:
- Terminal to same terminal remains idempotent.
- Terminal to non-terminal is rejected from status callback, completion callback, cancellation, and awaiter status paths.
- `waiting` can transition only to `running`, `cancelled`, or `failed`.
- Unknown current state can initialize to `pending` only. A legacy missing projection that must become `running` must be created as `pending` and then transitioned through the normal validator inside one aggregate transaction.
- Workflow write failure rolls back the `executions` row.
- Durable execution event append failure rolls back both state rows.
- Existing HTTP error codes and JSON shapes are preserved for invalid transitions unless explicitly covered by a compatibility-change test.

### TDD Cycle

#### Red: Write Failing Tests

**Files**:
- `agentfield/control-plane/pkg/types/status_test.go`
- `agentfield/control-plane/internal/handlers/execute_status_update_test.go`
- `agentfield/control-plane/internal/storage/execution_records_transaction_test.go`
- `agentfield/control-plane/internal/storage/execution_aggregate_test.go`
- `agentfield/control-plane/internal/handlers/workflow_execution_events_test.go`

```go
func TestValidateExecutionTransitionCentralRules(t *testing.T) {
    require.NoError(t, types.ValidateExecutionTransition("running", "succeeded"))
    require.NoError(t, types.ValidateExecutionTransition("failed", "failed"))
    require.Error(t, types.ValidateExecutionTransition("failed", "running"))
    require.Error(t, types.ValidateExecutionTransition("waiting", "succeeded"))
    require.NoError(t, types.ValidateExecutionTransition("waiting", "running"))
    require.NoError(t, types.ValidateExecutionTransition("unknown", "pending"))
    require.Error(t, types.ValidateExecutionTransition("unknown", "running"))
}

func TestCreateExecutionAggregateRollsBackWhenWorkflowWriteFails(t *testing.T) {
    // Use a storage stub or local storage fault injection.
    // Given workflow_executions write fails, assert GetExecutionRecord returns not found.
}

func TestUpdateExecutionAggregateReturnsUpdatedRowsAndStoredEvent(t *testing.T) {
    // Given a running aggregate, when it transitions to succeeded,
    // assert the returned aggregate includes updated executions,
    // workflow_executions, and a committed ExecutionCompleted outbox event.
}

func TestCancelInvalidTerminalStateKeepsExistingErrorShape(t *testing.T) {
    // Preserve 409 {"error":"invalid_state", ...} while moving through the validator.
}
```

#### Green: Minimal Implementation

**Files**:
- `agentfield/control-plane/pkg/types/status.go`
- `agentfield/control-plane/internal/storage/execution_state_validation.go`
- `agentfield/control-plane/internal/storage/storage.go`
- `agentfield/control-plane/internal/storage/execution_records.go`
- `agentfield/control-plane/internal/handlers/execute.go`
- `agentfield/control-plane/internal/handlers/execute_cancel.go`
- `agentfield/control-plane/internal/handlers/execute_cancel_tree.go`
- `agentfield/control-plane/internal/handlers/execute_awaiter_status.go`

Minimal steps:
- Move or mirror `validateExecutionStateTransition` into exported `types.ValidateExecutionTransition`.
- Replace handler-local transition guards with calls to the single validator.
- Add storage commands and aggregate DTOs such as:

```go
type ExecutionAggregate struct {
    Execution         *types.Execution
    WorkflowExecution *types.WorkflowExecution
    CommittedEvents   []events.StoredExecutionEvent
}

type ExecutionAggregateMutation func(exec *types.Execution, workflow *types.WorkflowExecution) ([]events.ExecutionEvent, error)

type ExecutionAggregateStore interface {
    CreateExecutionAggregate(ctx context.Context, exec *types.Execution, workflow *types.WorkflowExecution, initialEvents []events.ExecutionEvent) (*ExecutionAggregate, error)
    UpdateExecutionAggregate(ctx context.Context, executionID string, mutator ExecutionAggregateMutation) (*ExecutionAggregate, error)
}
```

- Inside each aggregate command, update `executions`, update/create `workflow_executions`, append required execution outbox events, and commit or roll back as one unit.
- Use the aggregate command from `prepareExecution`, sync/async completion, failure, status callback, cancel, cancel tree, pause/resume, approval, awaiter status updates, workflow event ingestion, webhook-driven updates, and any reaper/retry writer found by the writer inventory grep.
- Keep existing read/query APIs unchanged.
- Preserve public response payloads by building them from the returned `ExecutionAggregate`.

#### Refactor: Improve Code

**Files**:
- `agentfield/control-plane/internal/handlers/execute.go`
- `agentfield/control-plane/internal/storage/local.go`
- `agentfield/control-plane/internal/storage/unitofwork.go`

Refactor steps:
- Reuse existing unit-of-work infrastructure where practical.
- If the existing unit of work cannot include primary `executions`, either extend it under tests or keep a private aggregate transaction helper local to storage.
- Extract execution application service boundaries only after behavior is covered.
- Make `workflow_executions` a projection candidate by ensuring all final state changes also emit durable execution events from Behavior 1.
- Keep `types.ValidateExecutionTransition` pure and flat. Do not put storage reads, event appends, or mutation inside validator conditionals.
- Replace magic retry counts or timeout literals touched in this phase with named constants.

### Success Criteria

**Automated:**
- [ ] Red validator tests fail before central validator exists.
- [ ] Red aggregate tests fail if updated records or committed events are not returned.
- [ ] Red compatibility tests fail if cancel/status/workflow event ingestion response shapes change unexpectedly.
- [ ] `cd agentfield/control-plane && go test ./pkg/types -run 'Test(ValidateExecutionTransition|NormalizeExecutionStatus|IsTerminalExecutionStatus)'`.
- [ ] `cd agentfield/control-plane && go test ./internal/handlers -run 'Test(UpdateExecutionStatusHandler|CancelExecutionHandler|UpdateAwaiterStatusHandler|ExecuteHandler|ExecuteAsyncHandler)'`.
- [ ] `cd agentfield/control-plane && go test ./internal/storage -run 'Test(ExecutionAggregate|WorkflowUnitOfWork|Invariant_ExecutionState)'`.
- [ ] `cd agentfield/control-plane && go test -tags=integration ./internal/storage -run 'TestWorkflowUnitOfWork'`.

**Manual:**
- [ ] Grep shows no duplicate ad hoc transition guard remains in handlers.
- [ ] Grep inventory for `UpdateExecutionRecord|StoreWorkflowExecution|UpdateWorkflowExecution|StoreWorkflowExecutionEvent` is either migrated to aggregate commands or explicitly documented as a read/projection compatibility path.
- [ ] Public execute, async execute, status, cancel, pause/resume, and workflow UI endpoints still respond with the same JSON shapes.

## Behavior 3: VC Issuance Consumes ExecutionCompleted

### Test Specification

**Given**: An execution reaches `succeeded` or another terminal status that should produce a VC.
**When**: The durable `ExecutionCompleted` event is committed.
**Then**: a VC issuer consumer uses the DID context carried by the stored event or an explicitly referenced persisted context, atomically claims the logical VC key, calls `VCService.GenerateExecutionVC`, stores the VC once, and emits `VCIssued`.

**Edge Cases**:
- Duplicate `ExecutionCompleted` delivery does not create duplicate VCs.
- Concurrent legacy manual endpoint and event-driven issuance resolve to the same logical VC.
- Missing DID or parent VC context does not fabricate a credential; issuance records a structured skip/error and can retry when configured.
- VC generation failure retries without blocking execution completion.
- VC disabled config skips issuance.
- Legacy `POST /api/v1/execution/vc` remains available during migration.
- SDK decorator no longer calls `_generate_vc_async` after tracked reasoner completion.
- Python `Agent` and Go SDK automatic fire-and-forget VC hooks are disabled or routed through event-driven compatibility shims.
- TypeScript manual DID credential generation stays manual and is covered by compatibility tests.

### TDD Cycle

#### Red: Write Failing Tests

**Files**:
- `agentfield/control-plane/internal/services/vc_issuer_test.go`
- `agentfield/control-plane/internal/handlers/did_handlers_test.go`
- `agentfield/sdk/python/tests/test_decorators.py`
- `agentfield/sdk/python/tests/test_agent_vc.py`
- `agentfield/sdk/go/agent/*_test.go`
- `agentfield/sdk/typescript/src/did/*.test.ts`

```go
func TestVCIssuerConsumesExecutionCompletedOnce(t *testing.T) {
    store := newTestExecutionStorage(agent)
    publisher := newTestExecutionEventPublisher(store)
    issuer := services.NewExecutionCompletedVCIssuer(store, vcService, publisher)

    publishCompletedTwiceWithDIDContext(publisher, "exec-1")

    require.Eventually(t, func() bool {
        vcs, _ := store.ListExecutionVCs(context.Background(), types.VCFilters{ExecutionID: ptr("exec-1")})
        return len(vcs) == 1
    }, time.Second, 10*time.Millisecond)
}

func TestVCIssuerDuplicateDeliveryUsesLogicalKeyConflictAsSuccess(t *testing.T) {
    // Publish the same ExecutionCompleted event concurrently.
    // Assert one stored VC for (execution_id, issuer_did, target_did)
    // and no repeated uniqueness-error retry loop.
}

func TestCreateExecutionVCManualEndpointSharesIssuerIdempotency(t *testing.T) {
    // First call POST /api/v1/execution/vc.
    // Then deliver ExecutionCompleted for the same logical key.
    // Assert the existing VC is returned/loaded and no duplicate row is created.
}

func TestVCIssuerSkipsWhenCompletedEventLacksDIDContext(t *testing.T) {
    // Assert no fabricated CallerDID/TargetDID/AgentNodeDID is used.
}
```

```python
@pytest.mark.asyncio
async def test_execute_with_tracking_does_not_fire_vc_task(monkeypatch):
    tasks = []
    monkeypatch.setattr(asyncio, "create_task", lambda coro: tasks.append(coro))
    # Execute tracked reasoner with vc enabled.
    # Assert no task targeting _generate_vc_async is scheduled.
```

#### Green: Minimal Implementation

**Files**:
- `agentfield/control-plane/internal/services/vc_issuer.go`
- `agentfield/control-plane/internal/services/vc_issuance.go`
- `agentfield/control-plane/internal/storage/execution_vc_idempotency.go`
- `agentfield/control-plane/internal/server/routes_core.go` or bootstrap file
- `agentfield/sdk/python/agentfield/decorators.py`
- `agentfield/sdk/python/agentfield/agent.py`
- `agentfield/sdk/go/agent/agent_did.go`
- `agentfield/sdk/typescript/src/did/DidClient.ts`

Minimal steps:
- Add `ExecutionCompletedVCIssuer` service with explicit `Start(ctx)` and `Stop(ctx)` lifecycle that subscribes to the injected outbox-backed publisher from Behavior 1.
- Extend the `ExecutionCompleted` typed event payload so it carries or references the complete `types.ExecutionContext` fields needed for VC generation. Do not reconstruct DID fields from execution/workflow rows unless those fields are explicitly persisted by this behavior.
- If the event lacks `CallerDID`, `TargetDID`, or `AgentNodeDID`, skip issuance with a structured event/log and do not generate a fake context.
- Add an atomic idempotency command around the logical key `(execution_id, issuer_did, target_did)`. Implement either a pre-signing claim row or deterministic VC ID plus conflict-as-success. A read-before-generate check alone is not allowed.
- Call `VCService.GenerateExecutionVC` and persist through existing `ShouldPersistExecutionVC` behavior.
- Update `CreateExecutionVC` handler to share the same idempotency command while keeping the public request/response contract intact.
- Disable SDK fire-and-forget VC generation behind a feature flag defaulted to event-driven mode, then remove direct decorator scheduling after tests prove compatibility.
- Add Python `Agent`, Python decorator, Go SDK, and TypeScript DID regression coverage so automatic direct issuance is not reintroduced while manual VC APIs remain available.

#### Refactor: Improve Code

**Files**:
- `agentfield/control-plane/internal/services/vc_issuer.go`
- `agentfield/sdk/python/agentfield/agent_vc.py`
- `agentfield/sdk/python/agentfield/decorators.py`
- `agentfield/sdk/go/agent/agent_did.go`

Refactor steps:
- Move retry/backoff into the consumer.
- Emit typed `VCIssued` event after successful generation.
- Keep SDK metadata registration so the control plane still knows per-reasoner VC policy.
- Remove duplicated VC helper code from `agent.py` only after `agent_vc.py` tests cover it.
- Compute VC eligibility booleans in separate statements before side-effect calls. Remove or route all automatic scheduling paths serially, without scheduling side effects inside compound conditionals.

### Success Criteria

**Automated:**
- [ ] Red tests fail because no consumer exists and decorators still schedule VC work.
- [ ] Red tests fail if duplicate event delivery, manual endpoint calls, or concurrent issuance creates more than one VC for `(execution_id, issuer_did, target_did)`.
- [ ] Red tests fail if the issuer fabricates DID context from incomplete stored rows.
- [ ] `cd agentfield/control-plane && go test ./internal/services -run 'TestVCIssuer|TestVCService'`.
- [ ] `cd agentfield/control-plane && go test ./internal/handlers -run 'TestDIDHandlers|TestCreateExecutionVC'`.
- [ ] `cd agentfield/sdk/python && pytest tests/test_decorators.py tests/test_agent_networking.py tests/test_client_execution_vc_payload.py tests/test_agent_vc.py`.
- [ ] `cd agentfield/sdk/go && go test ./agent -run 'Test(DID|VC|MaybeGenerateVC)'`.
- [ ] `cd agentfield/sdk/typescript && CI=1 npm run test:core -- --runInBand`.

**Manual:**
- [ ] Execution completion latency is unchanged by VC issuance failures.
- [ ] Existing clients using `/api/v1/execution/vc` still work during migration.
- [ ] Grep confirms no Python decorator, Python `Agent`, or Go automatic completion path schedules direct VC generation by default.

## Behavior 4: Callback Reachability Probe During Registration

### Test Specification

**Given**: A registering agent provides callback candidates.
**When**: registration is in auto-discovery mode.
**Then**: the control plane normalizes candidates, applies the callback probing security policy, probes bounded candidates on named paths, records `CallbackTestResult` entries using `Success`, selects the first reachable candidate, and rejects or marks pending candidates when none are reachable according to existing policy.

**Edge Cases**:
- Explicit/manual mode preserves the explicit URL without probing.
- Invalid URL candidates are recorded as failed tests.
- Blocked private, link-local, metadata, or non-allowlisted hosts are recorded as failed tests without a network request.
- Redirect responses are not followed by default and are recorded according to config.
- Probe timeout is bounded.
- Probe candidate count and total registration probe time are bounded.
- Local `httptest.Server` candidate succeeds.
- First candidate unreachable, second reachable.
- Both `candidates` and legacy Python `callback_candidates` JSON keys are accepted.
- Go and TypeScript SDK localhost defaults are documented as compatibility debt and guarded by SDK-specific tests until fixed.

### TDD Cycle

#### Red: Write Failing Tests

**Files**:
- `agentfield/control-plane/internal/handlers/nodes_discovery_test.go`
- `agentfield/control-plane/internal/handlers/nodes_register_test.go`
- `agentfield/control-plane/internal/config/config_test.go`
- `agentfield/sdk/python/tests/test_agent_networking.py`

```go
func TestProbeCallbackCandidatesSelectsFirstSuccessfulCandidate(t *testing.T) {
    unreachable := "http://127.0.0.1:1"
    reachable := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        require.Equal(t, "/health", r.URL.Path)
        w.WriteHeader(http.StatusOK)
    }))
    defer reachable.Close()

    candidates := normalizeCallbackCandidates([]string{unreachable, reachable.URL}, "")
    resolved, results := probeCallbackCandidates(context.Background(), candidates, testProbeConfig())

    require.Equal(t, reachable.URL, resolved)
    require.Len(t, results, 2)
    require.False(t, results[0].Success)
    require.True(t, results[1].Success)
}

func TestCallbackDiscoveryAcceptsLegacyCallbackCandidatesKey(t *testing.T) {
    body := `{"callback_discovery":{"mode":"auto","callback_candidates":["http://127.0.0.1:8080"]}}`
    // Bind registration request and assert the candidate is present.
}

func TestCallbackProbeBlocksMetadataAddressWithoutDialing(t *testing.T) {
    candidates := normalizeCallbackCandidates([]string{"http://169.254.169.254/latest/meta-data"}, "")
    resolved, results := probeCallbackCandidates(context.Background(), candidates, productionProbeConfig())

    require.Empty(t, resolved)
    require.Len(t, results, 1)
    require.False(t, results[0].Success)
    require.Contains(t, results[0].Error, "blocked")
}
```

#### Green: Minimal Implementation

**Files**:
- `agentfield/control-plane/internal/handlers/nodes_register.go`
- `agentfield/control-plane/internal/config/config.go`
- `agentfield/control-plane/pkg/types/types.go`
- `agentfield/sdk/python/agentfield/agent.py`
- `agentfield/sdk/python/agentfield/agent_field_handler.py`

Minimal steps:
- Add named constants/config: `CallbackProbeHealthPath`, `CallbackProbeRootPath`, `CallbackProbeTimeout`, `CallbackProbeMaxCandidates`, `CallbackProbeTotalBudget`, and redirect policy.
- Add a `CallbackProber` interface with default HTTP implementation that does not follow redirects by default.
- Keep `normalizeCallbackCandidates` pure; it only parses, cleans, deduplicates, and returns candidates plus normalization errors.
- Add a security filter before dialing: only `http`/`https`, no userinfo, no opaque URLs, no malformed hosts, no metadata addresses, and private/loopback/link-local allowed only in test/dev mode or when explicitly allowlisted.
- Probe normalized allowed candidates with bounded timeout against `/health` first, then fallback to root only if needed.
- Populate existing `CallbackDiscovery.Tests`.
- Select first reachable URL in auto mode.
- Ensure SDK `base_url` uses the same detection ladder already tested by `_build_callback_candidates`.
- Decode both `callback_discovery.candidates` and `callback_discovery.callback_candidates` during the transition.
- Keep response field names aligned with `CallbackTestResult.Success`. Do not add `Reachable` unless a versioned compatibility plan is created.

#### Refactor: Improve Code

**Files**:
- `agentfield/control-plane/internal/handlers/nodes_register.go`
- `agentfield/control-plane/internal/services/agent_registration.go`
- `agentfield/sdk/python/agentfield/agent.py`

Refactor steps:
- Extract `AgentRegistrationService` after probe tests are green.
- Keep route handler as request/response glue.
- Centralize SDK callback URL resolution to avoid `agent.py` and `agent_field_handler.py` duplicating candidate logic.
- Keep probing side effects out of conditional expressions. Normalize, filter, probe, and select in separate named steps.

### Success Criteria

**Automated:**
- [ ] Existing `TestResolveCallbackCandidatesSuccess` changes from expecting nil probe results to expecting explicit reachability data.
- [ ] Red tests fail if a test references `Reachable` instead of `Success`.
- [ ] Red tests fail if `callback_candidates` is ignored.
- [ ] Red tests fail if metadata/private/blocked hosts are dialed under production probe config.
- [ ] `cd agentfield/control-plane && go test ./internal/handlers -run 'Test(ResolveCallbackCandidates|RegisterAgent|Callback)'`.
- [ ] `cd agentfield/control-plane && go test ./internal/config -run 'TestCallbackProbe'`.
- [ ] `cd agentfield/sdk/python && pytest tests/test_agent_networking.py`.

**Manual:**
- [ ] Registration response includes `resolved_base_url` and `callback_discovery.tests`.
- [ ] Off-host deployment no longer stores `localhost` when a reachable non-local candidate is available.
- [ ] Probe path, timeout, max-candidate, total-budget, and redirect behavior are named constants/config, not inline literals.

## Behavior 5: SDK Agent Decomposition Behind Stable Events

### Test Specification

**Given**: Existing SDK users call `app.call()`, `app.note()`, `app.pause()`, and decorated reasoners.
**When**: internal services are extracted from `Agent`.
**Then**: public methods keep the same signatures and behavior, outbound note/status traffic uses the client adapter, and context propagation is observable through a single `ExecutionContext` path.

**Edge Cases**:
- `app.call()` with parent context preserves run ID, parent execution ID, session, actor, and DID headers.
- `app.note()` with no execution context remains a no-op or existing logged behavior.
- `app.note()` with context sends the same UI note payload shape through `AgentFieldClient.post_note`.
- `app.note()` preserves existing sync and async/fire-and-forget scheduling behavior, including task completion visibility in tests.
- `app.note()` preserves timeout and error handling semantics while moving transport out of `Agent`.
- `app.pause()` requires callback URL and resumes through existing webhook flow.
- Exceptions restore context.
- Legacy tests that set `agent._current_execution_context` continue to pass during migration.

### TDD Cycle

#### Red: Write Failing Tests

**Files**:
- `agentfield/sdk/python/tests/test_agent_call.py`
- `agentfield/sdk/python/tests/test_agent_note.py`
- `agentfield/sdk/python/tests/test_agent_workflow.py`
- `agentfield/sdk/python/tests/test_execution_context_core.py`

```python
@pytest.mark.asyncio
async def test_app_call_delegates_to_call_gateway(monkeypatch):
    agent, _ = create_test_agent(monkeypatch)
    calls = []

    class StubGateway:
        async def call(self, target, *args, **kwargs):
            calls.append((target, args, kwargs))
            return {"ok": True}

    agent.call_gateway = StubGateway()
    assert await agent.call("node.reasoner", {"x": 1}) == {"ok": True}
    assert calls[0][0] == "node.reasoner"

def test_note_uses_client_adapter(monkeypatch):
    agent, _ = create_test_agent(monkeypatch)
    sent = []
    agent.client.post_note = lambda *args, **kwargs: sent.append((args, kwargs))
    agent._current_execution_context = ExecutionContext.create_new(agent.node_id, "reasoner")

    agent.note("hello", tags=["debug"])

    assert sent

def test_note_without_execution_context_preserves_noop_behavior(monkeypatch):
    agent, _ = create_test_agent(monkeypatch)
    sent = []
    agent.client.post_note = lambda *args, **kwargs: sent.append((args, kwargs))

    agent.note("hello")

    assert sent == []

@pytest.mark.asyncio
async def test_note_async_path_completes_scheduled_client_call(monkeypatch):
    agent, _ = create_test_agent(monkeypatch)
    calls = []

    async def post_note(*args, **kwargs):
        calls.append((args, kwargs))

    agent.client.post_note = post_note
    agent._current_execution_context = ExecutionContext.create_new(agent.node_id, "reasoner")

    agent.note("hello", tags=["debug"])
    await drain_scheduled_tasks()

    assert calls

def test_note_preserves_timeout_and_error_mapping(monkeypatch):
    agent, _ = create_test_agent(monkeypatch)
    agent._current_execution_context = ExecutionContext.create_new(agent.node_id, "reasoner")
    agent.client.post_note = Mock(side_effect=TimeoutError("boom"))

    # Assert the same logged/no-raise behavior as the current public method.
    agent.note("hello")
```

#### Green: Minimal Implementation

**Files**:
- `agentfield/sdk/python/agentfield/agent.py`
- `agentfield/sdk/python/agentfield/call_gateway.py`
- `agentfield/sdk/python/agentfield/client.py`
- `agentfield/sdk/python/agentfield/agent_pause.py`
- `agentfield/sdk/python/agentfield/execution_context.py`

Minimal steps:
- Extract `CallGateway` and have `Agent.call()` delegate without changing signature.
- Add `AgentFieldClient.post_note()` with the same payload, timeout, and error mapping currently implemented inside `Agent.note()`.
- Have `Agent.note()` become a flat adapter: resolve current execution context, build note payload, call or schedule `client.post_note`, and return according to existing behavior.
- Remove direct `aiohttp` and `requests` note transport from `Agent.note()` once adapter tests are green.
- Keep `Agent.pause()` delegating to the existing pause manager/service.
- Preserve `_current_execution_context` as a compatibility facade while `ExecutionContextManager` becomes the source of truth.
- Keep automatic VC scheduling disabled/routed as established in Behavior 3 while extracting call/note/pause services.

#### Refactor: Improve Code

**Files**:
- `agentfield/sdk/python/agentfield/agent.py`
- `agentfield/sdk/python/agentfield/decorators.py`
- `agentfield/sdk/python/agentfield/agent_workflow.py`
- `agentfield/sdk/python/agentfield/router.py`

Refactor steps:
- Delete dead duplicate pause manager only after pause tests pass.
- Remove direct context writes from decorators and workflow helper by routing through an execution context port.
- Collapse reasoner registration paths after the call/note/pause extractions are covered.
- Flatten public SDK methods touched in this phase. Avoid nested helper definitions and transport fallback branches inside `Agent.note()`, `Agent.call()`, and `Agent.pause()`.
- Replace raw timeout/default literals touched in this phase with named constants such as note timeout and pause default hours.

### Success Criteria

**Automated:**
- [ ] Red tests fail because `CallGateway` or `client.post_note` does not exist.
- [ ] Red tests fail if `Agent.note()` only delegates but drops no-context behavior, async scheduling, timeout/error behavior, or payload shape.
- [ ] `cd agentfield/sdk/python && pytest tests/test_agent_call.py tests/test_agent_workflow.py tests/test_decorators.py tests/test_execution_context_core.py tests/test_agent_networking.py`.
- [ ] `cd agentfield/sdk/python && pytest tests/test_agent_note.py tests/test_agent_coverage_additions.py tests/test_async_execution.py`.
- [ ] `cd agentfield/sdk/python && pytest -q -m 'not harness_live'`.

**Manual:**
- [ ] Existing sample agents do not need code changes.
- [ ] `rg "_current_execution_context|_current_workflow_context" agentfield/sdk/python/agentfield` shows only compatibility accessors until the final cleanup cycle.
- [ ] `rg "aiohttp|requests" agentfield/sdk/python/agentfield/agent.py` shows `Agent.note()` no longer owns direct note transport fallback logic.

## Integration and E2E Testing

- **Control plane focused:** `cd agentfield/control-plane && go test ./internal/events ./internal/handlers ./internal/storage ./internal/services ./internal/config ./pkg/types`.
- **Storage integration:** `cd agentfield/control-plane && go test -tags=integration ./internal/storage -run 'TestWorkflowUnitOfWork|TestExecutionEventOutbox|TestExecutionAggregate'`.
- **SDK focused:** `cd agentfield/sdk/python && pytest tests/test_agent_networking.py tests/test_decorators.py tests/test_agent_call.py tests/test_agent_note.py tests/test_agent_workflow.py tests/test_client_execution_vc_payload.py tests/test_agent_vc.py`.
- **Go SDK compatibility:** `cd agentfield/sdk/go && go test ./agent -run 'Test(Initialize_RegistersNodeAndMarksReady|RegisterNode|DID|VC)'`.
- **TypeScript SDK compatibility:** `cd agentfield/sdk/typescript && CI=1 npm run test:core`.
- **Full repo bounded check:** `make -C agentfield test` after all five phases, with live/harness tests excluded by existing config.

## Implementation Order and Checkpoints

1. **Event contract checkpoint:** `execution_event_outbox` owns global sequence allocation, replay/list, unsubscribe, injected bus subscribers, publish-after-commit, SSE compatibility, and no dual publish in `execute.go`.
2. **Execution aggregate checkpoint:** central pure validator, transactional create/update aggregate path, committed event append, updated aggregate return values, full writer inventory migrated or documented, and existing public endpoints green.
3. **VC event checkpoint:** `ExecutionCompleted` carries complete VC context, `ExecutionCompletedVCIssuer` idempotently writes VCs by `(execution_id, issuer_did, target_did)`, manual endpoint shares idempotency, and SDK automatic issuance no longer owns completion credentials.
4. **Registry probe checkpoint:** auto-discovery decodes both candidate keys, probes under SSRF-safe bounded policy, returns real `CallbackTestResult.Success` values, and chooses the first successful URL.
5. **SDK decomposition checkpoint:** public SDK tests pass while `Agent` delegates call/note/pause/invocation concerns to services; `Agent.note()` is a flat adapter with no direct HTTP fallback transport.
6. **Review closure checkpoint:** rerun this review checklist and confirm every critical/warning item is represented by a contract, red test, success criterion, or explicit non-goal.

## References

- Source index: `specs/INDEX.md`.
- Eventing spec: `specs/eventing-observability.domain.md`.
- Execution spec: `specs/execution-orchestration.domain.md`.
- Identity spec: `specs/identity-credentials.domain.md`.
- Registry spec: `specs/agent-registry.domain.md`.
- SDK spec: `specs/agent-runtime-sdk.domain.md`.
- Event bus: `agentfield/control-plane/internal/events/event_bus.go`.
- Execution events: `agentfield/control-plane/internal/events/execution_events.go`.
- Execution handler: `agentfield/control-plane/internal/handlers/execute.go`.
- Execution status model: `agentfield/control-plane/pkg/types/status.go`.
- Registry callback discovery: `agentfield/control-plane/internal/handlers/nodes_register.go`.
- VC handler/service: `agentfield/control-plane/internal/handlers/did_handlers.go`, `agentfield/control-plane/internal/services/vc_issuance.go`.
- SDK agent/decorators: `agentfield/sdk/python/agentfield/agent.py`, `agentfield/sdk/python/agentfield/decorators.py`.
- SDK cross-language callback/VC paths: `agentfield/sdk/go/agent/agent_did.go`, `agentfield/sdk/typescript/src/did/DidClient.ts`, `agentfield/sdk/typescript/src/agent/Agent.ts`.
