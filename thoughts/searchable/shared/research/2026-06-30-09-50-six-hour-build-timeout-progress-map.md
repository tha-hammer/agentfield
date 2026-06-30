---
date: 2026-06-30T09:50:27-04:00
researcher: Codex
git_commit: cbd2ddf7a14ac21524994c600ef513d163d919fe
branch: main
repository: agentfield
topic: "Map how the 6 hour build timeout is fed and where progress signals enter"
tags: [research, codebase, timeout, async-execution, control-plane, sdk-python, harness]
status: complete
last_updated: 2026-06-30
last_updated_by: Codex
beads: unavailable - bd list --status=open returned "no beads database found"
---

# Research: Six-hour Build Timeout and Progress Signals

**Date**: 2026-06-30 09:50:27 -04:00  
**Researcher**: Codex  
**Git Commit**: cbd2ddf7a14ac21524994c600ef513d163d919fe  
**Branch**: main  
**Repository**: agentfield

## Research Question

The code uses a 6 hour timeout as a guard against builds that are not progressing. The timer kills builds that are long and are progressing. Map the timer interfaces, contracts, and seams to understand how the 6 hour timer is fed. We need to consider a better method to let builds run if it is progressing.

## Summary

There is no dedicated persisted `builds` model or `/builds` API in this repo. "Build" state is represented through the same execution, workflow execution, workflow run, status, event, log, and note contracts used by all reasoners.

The six-hour value appears in the Python SDK async configuration as `max_execution_timeout = 21600.0`. The reasoner-level watchdog that cancels an in-flight async reasoner uses `default_execution_timeout`, which defaults to 7200 seconds but can be set from `AGENTFIELD_ASYNC_DEFAULT_EXECUTION_TIMEOUT`. The `Agent.call()` child wait path chooses `max_execution_timeout` first, so its default wait budget is six hours.

The watchdog is fed by elapsed wall-clock time minus a local `PauseClock.total_paused()`. That pause clock is fed only by explicit waiting states: `Agent.pause()`, `Agent.wait_for_resume()`, and `Agent.call()` observing an awaited child in `WAITING`. It is not fed by ordinary logs, status progress payloads, notes, workflow events, or harness stdout/stderr activity.

The Go control plane has a separate stale-execution cleanup guard based on persisted `updated_at` activity. That guard is explicitly intended to time out non-terminal rows with no recent activity. Its feeders are persisted record updates such as status callbacks, approval/waiting transitions, awaiter-status transitions, pause/resume, cancel, notes, and terminal completion/failure. Structured execution logs are stored and streamed, but the current storage path does not update `executions.updated_at` or `workflow_executions.updated_at`.

The harness has one narrower progress-like timer today: the OpenCode provider passes both a total subprocess timeout and an idle no-output timeout into `run_cli()`. That idle timeout is based on stdout/stderr bytes, and it applies to one OpenCode subprocess, not to the six-hour SDK reasoner/watch budget.

## Detailed Findings

### Timer Inventory

| Surface | Current value/config | What it bounds | Fed by |
| --- | --- | --- | --- |
| Python SDK `max_execution_timeout` | `21600.0` seconds, env `AGENTFIELD_ASYNC_MAX_EXECUTION_TIMEOUT` | Client polling/session total and default `Agent.call()` child wait | Wall-clock in polling/waiting code |
| Python SDK `default_execution_timeout` | `7200.0` seconds, env `AGENTFIELD_ASYNC_DEFAULT_EXECUTION_TIMEOUT` | Async reasoner watchdog active-time budget | Wall-clock minus `PauseClock.total_paused()` |
| Control-plane `agent_call_timeout` | YAML `agentfield.execution_queue.agent_call_timeout`, sample config `1800s` | Go controller HTTP call to agent and sync wait for async callback | Go `http.Client` timeout and `time.NewTimer(timeout)` |
| Control-plane stale cleanup | YAML `agentfield.execution_cleanup.stale_execution_timeout`, sample config `10m`, struct default `30m` | Non-terminal persisted rows with no recent activity | `updated_at`, falling back to `created_at` or `started_at` |
| Go request middleware | hard-coded `3600*time.Second` | Whole HTTP request context | Context deadline |
| Async worker background context | hard-coded `24*time.Hour` | Control-plane async worker goroutine lifetime | Context deadline |
| OpenCode harness subprocess | `AGENTFIELD_HARNESS_TIMEOUT_SECONDS`, default `1800`; `AGENTFIELD_HARNESS_IDLE_TIMEOUT_SECONDS`, default `600` | One OpenCode CLI process | Total wall-clock and stdout/stderr idle activity |

Key code references:

- [`sdk/python/agentfield/async_config.py:33`](https://github.com/Agent-Field/agentfield/blob/cbd2ddf7a14ac21524994c600ef513d163d919fe/sdk/python/agentfield/async_config.py#L33) defines `max_execution_timeout: float = 21600.0`.
- [`sdk/python/agentfield/async_config.py:34`](https://github.com/Agent-Field/agentfield/blob/cbd2ddf7a14ac21524994c600ef513d163d919fe/sdk/python/agentfield/async_config.py#L34) defines `default_execution_timeout: float = 7200.0`.
- [`sdk/python/agentfield/async_config.py:137`](https://github.com/Agent-Field/agentfield/blob/cbd2ddf7a14ac21524994c600ef513d163d919fe/sdk/python/agentfield/async_config.py#L137) reads the max timeout env var, and [`sdk/python/agentfield/async_config.py:140`](https://github.com/Agent-Field/agentfield/blob/cbd2ddf7a14ac21524994c600ef513d163d919fe/sdk/python/agentfield/async_config.py#L140) reads the default timeout env var.
- [`control-plane/internal/config/config.go:149`](https://github.com/Agent-Field/agentfield/blob/cbd2ddf7a14ac21524994c600ef513d163d919fe/control-plane/internal/config/config.go#L149) defines execution cleanup config; [`control-plane/config/agentfield.yaml:11`](https://github.com/Agent-Field/agentfield/blob/cbd2ddf7a14ac21524994c600ef513d163d919fe/control-plane/config/agentfield.yaml#L11) sets `stale_execution_timeout: 10m`.
- [`control-plane/internal/config/config.go:161`](https://github.com/Agent-Field/agentfield/blob/cbd2ddf7a14ac21524994c600ef513d163d919fe/control-plane/internal/config/config.go#L161) defines execution queue config; [`control-plane/config/agentfield.yaml:13`](https://github.com/Agent-Field/agentfield/blob/cbd2ddf7a14ac21524994c600ef513d163d919fe/control-plane/config/agentfield.yaml#L13) sets `agent_call_timeout: 1800s`.

### The Six-hour SDK Budget

When an agent receives a request with `X-Execution-ID`, the FastAPI reasoner endpoint creates a background task and returns `202`. The task is registered for cancellation and then run through `_execute_async_with_callback()`.

- [`sdk/python/agentfield/agent.py:1900`](https://github.com/Agent-Field/agentfield/blob/cbd2ddf7a14ac21524994c600ef513d163d919fe/sdk/python/agentfield/agent.py#L1900) detects `X-Execution-ID`.
- [`sdk/python/agentfield/agent.py:1902`](https://github.com/Agent-Field/agentfield/blob/cbd2ddf7a14ac21524994c600ef513d163d919fe/sdk/python/agentfield/agent.py#L1902) creates `_execute_async_with_callback(...)`.
- [`sdk/python/agentfield/agent.py:1916`](https://github.com/Agent-Field/agentfield/blob/cbd2ddf7a14ac21524994c600ef513d163d919fe/sdk/python/agentfield/agent.py#L1916) registers the task with the cancel registry.

Inside `_execute_async_with_callback()`, the watchdog takes `reasoner_timeout` from `async_config.default_execution_timeout`, creates a `PauseClock`, records `start_time`, and periodically computes:

```text
active_elapsed = (time.time() - start_time) - pause_clock.total_paused()
```

If `active_elapsed > reasoner_timeout`, it marks the pause clock as timed out, cancels the reasoner task, and posts a failed terminal status whose error details include `reasoner_timeout`.

- [`sdk/python/agentfield/agent.py:2376`](https://github.com/Agent-Field/agentfield/blob/cbd2ddf7a14ac21524994c600ef513d163d919fe/sdk/python/agentfield/agent.py#L2376) documents the active-time budget.
- [`sdk/python/agentfield/agent.py:2383`](https://github.com/Agent-Field/agentfield/blob/cbd2ddf7a14ac21524994c600ef513d163d919fe/sdk/python/agentfield/agent.py#L2383) reads `default_execution_timeout`.
- [`sdk/python/agentfield/agent.py:2389`](https://github.com/Agent-Field/agentfield/blob/cbd2ddf7a14ac21524994c600ef513d163d919fe/sdk/python/agentfield/agent.py#L2389) creates and stores the pause clock.
- [`sdk/python/agentfield/agent.py:2399`](https://github.com/Agent-Field/agentfield/blob/cbd2ddf7a14ac21524994c600ef513d163d919fe/sdk/python/agentfield/agent.py#L2399) defines the watchdog loop.
- [`sdk/python/agentfield/agent.py:2409`](https://github.com/Agent-Field/agentfield/blob/cbd2ddf7a14ac21524994c600ef513d163d919fe/sdk/python/agentfield/agent.py#L2409) computes active elapsed time.
- [`sdk/python/agentfield/agent.py:2410`](https://github.com/Agent-Field/agentfield/blob/cbd2ddf7a14ac21524994c600ef513d163d919fe/sdk/python/agentfield/agent.py#L2410) fires the timeout and cancels the task.
- [`sdk/python/agentfield/agent.py:2444`](https://github.com/Agent-Field/agentfield/blob/cbd2ddf7a14ac21524994c600ef513d163d919fe/sdk/python/agentfield/agent.py#L2444) builds the timeout failure callback.

The parent/child `Agent.call()` path has its own wait budget. It chooses `max_execution_timeout` first, falling back to `default_execution_timeout` and then `600.0`, so the default child wait is six hours even though the reasoner watchdog default is two hours.

- [`sdk/python/agentfield/agent.py:3938`](https://github.com/Agent-Field/agentfield/blob/cbd2ddf7a14ac21524994c600ef513d163d919fe/sdk/python/agentfield/agent.py#L3938) computes `execution_timeout`.
- [`sdk/python/agentfield/agent.py:3961`](https://github.com/Agent-Field/agentfield/blob/cbd2ddf7a14ac21524994c600ef513d163d919fe/sdk/python/agentfield/agent.py#L3961) passes that timeout to `execute_async`.
- [`sdk/python/agentfield/agent.py:4061`](https://github.com/Agent-Field/agentfield/blob/cbd2ddf7a14ac21524994c600ef513d163d919fe/sdk/python/agentfield/agent.py#L4061) waits with `wait_for_execution_result`.
- [`sdk/python/agentfield/async_execution_manager.py:421`](https://github.com/Agent-Field/agentfield/blob/cbd2ddf7a14ac21524994c600ef513d163d919fe/sdk/python/agentfield/async_execution_manager.py#L421) stores a timeout on local `ExecutionState`.
- [`sdk/python/agentfield/async_execution_manager.py:581`](https://github.com/Agent-Field/agentfield/blob/cbd2ddf7a14ac21524994c600ef513d163d919fe/sdk/python/agentfield/async_execution_manager.py#L581) chooses the wait timeout.
- [`sdk/python/agentfield/async_execution_manager.py:623`](https://github.com/Agent-Field/agentfield/blob/cbd2ddf7a14ac21524994c600ef513d163d919fe/sdk/python/agentfield/async_execution_manager.py#L623) loops until active elapsed reaches the timeout.

### What Feeds the SDK Active-time Clock

`PauseClock` is the only subtraction mechanism for the SDK watchdog and the `Agent.call()` wait loop.

- [`sdk/python/agentfield/agent_pause.py:10`](https://github.com/Agent-Field/agentfield/blob/cbd2ddf7a14ac21524994c600ef513d163d919fe/sdk/python/agentfield/agent_pause.py#L10) defines `PauseClock`.
- [`sdk/python/agentfield/agent_pause.py:33`](https://github.com/Agent-Field/agentfield/blob/cbd2ddf7a14ac21524994c600ef513d163d919fe/sdk/python/agentfield/agent_pause.py#L33) starts a pause.
- [`sdk/python/agentfield/agent_pause.py:37`](https://github.com/Agent-Field/agentfield/blob/cbd2ddf7a14ac21524994c600ef513d163d919fe/sdk/python/agentfield/agent_pause.py#L37) ends a pause.
- [`sdk/python/agentfield/agent_pause.py:42`](https://github.com/Agent-Field/agentfield/blob/cbd2ddf7a14ac21524994c600ef513d163d919fe/sdk/python/agentfield/agent_pause.py#L42) reports total paused seconds.

Direct human-approval pauses feed the clock:

- [`sdk/python/agentfield/agent.py:4472`](https://github.com/Agent-Field/agentfield/blob/cbd2ddf7a14ac21524994c600ef513d163d919fe/sdk/python/agentfield/agent.py#L4472) tells the control plane to transition to `waiting`.
- [`sdk/python/agentfield/agent.py:4505`](https://github.com/Agent-Field/agentfield/blob/cbd2ddf7a14ac21524994c600ef513d163d919fe/sdk/python/agentfield/agent.py#L4505) gets the current execution's pause clock.
- [`sdk/python/agentfield/agent.py:4507`](https://github.com/Agent-Field/agentfield/blob/cbd2ddf7a14ac21524994c600ef513d163d919fe/sdk/python/agentfield/agent.py#L4507) starts it before awaiting approval.
- [`sdk/python/agentfield/agent.py:4521`](https://github.com/Agent-Field/agentfield/blob/cbd2ddf7a14ac21524994c600ef513d163d919fe/sdk/python/agentfield/agent.py#L4521) ends it in `finally`.
- [`sdk/python/agentfield/agent.py:4560`](https://github.com/Agent-Field/agentfield/blob/cbd2ddf7a14ac21524994c600ef513d163d919fe/sdk/python/agentfield/agent.py#L4560) does the same for `wait_for_resume()`.

Child waits feed the parent's clock when the awaited child enters `WAITING`:

- [`sdk/python/agentfield/agent.py:3968`](https://github.com/Agent-Field/agentfield/blob/cbd2ddf7a14ac21524994c600ef513d163d919fe/sdk/python/agentfield/agent.py#L3968) documents cross-reasoner pause propagation.
- [`sdk/python/agentfield/agent.py:3976`](https://github.com/Agent-Field/agentfield/blob/cbd2ddf7a14ac21524994c600ef513d163d919fe/sdk/python/agentfield/agent.py#L3976) resolves the parent execution id.
- [`sdk/python/agentfield/agent.py:3981`](https://github.com/Agent-Field/agentfield/blob/cbd2ddf7a14ac21524994c600ef513d163d919fe/sdk/python/agentfield/agent.py#L3981) finds the parent pause clock.
- [`sdk/python/agentfield/agent.py:4011`](https://github.com/Agent-Field/agentfield/blob/cbd2ddf7a14ac21524994c600ef513d163d919fe/sdk/python/agentfield/agent.py#L4011) passes the pause clock into the wait.
- [`sdk/python/agentfield/async_execution_manager.py:684`](https://github.com/Agent-Field/agentfield/blob/cbd2ddf7a14ac21524994c600ef513d163d919fe/sdk/python/agentfield/async_execution_manager.py#L684) starts the parent clock when the child is waiting.
- [`sdk/python/agentfield/async_execution_manager.py:702`](https://github.com/Agent-Field/agentfield/blob/cbd2ddf7a14ac21524994c600ef513d163d919fe/sdk/python/agentfield/async_execution_manager.py#L702) ends the clock when the child is no longer waiting.
- [`sdk/python/agentfield/execution_state.py:234`](https://github.com/Agent-Field/agentfield/blob/cbd2ddf7a14ac21524994c600ef513d163d919fe/sdk/python/agentfield/execution_state.py#L234) makes the polling task's local overdue check pause-aware.

Multi-hop propagation is a control-plane contract. A parent watching an execution only sees that child become `WAITING` if the intermediate awaiter pushes its own execution into `WAITING` through the awaiter-status endpoint.

- [`sdk/python/agentfield/agent.py:4014`](https://github.com/Agent-Field/agentfield/blob/cbd2ddf7a14ac21524994c600ef513d163d919fe/sdk/python/agentfield/agent.py#L4014) documents the multi-hop propagation need.
- [`sdk/python/agentfield/agent.py:4030`](https://github.com/Agent-Field/agentfield/blob/cbd2ddf7a14ac21524994c600ef513d163d919fe/sdk/python/agentfield/agent.py#L4030) pushes self to `waiting`.
- [`sdk/python/agentfield/agent.py:4044`](https://github.com/Agent-Field/agentfield/blob/cbd2ddf7a14ac21524994c600ef513d163d919fe/sdk/python/agentfield/agent.py#L4044) pushes self back to `running`.
- [`sdk/python/agentfield/client.py:1787`](https://github.com/Agent-Field/agentfield/blob/cbd2ddf7a14ac21524994c600ef513d163d919fe/sdk/python/agentfield/client.py#L1787) implements `notify_awaiter_status()`.
- [`control-plane/internal/server/routes_core.go:134`](https://github.com/Agent-Field/agentfield/blob/cbd2ddf7a14ac21524994c600ef513d163d919fe/control-plane/internal/server/routes_core.go#L134) wires `/agents/:node_id/executions/:execution_id/awaiter-status`.
- [`control-plane/internal/handlers/execute_awaiter_status.go:21`](https://github.com/Agent-Field/agentfield/blob/cbd2ddf7a14ac21524994c600ef513d163d919fe/control-plane/internal/handlers/execute_awaiter_status.go#L21) describes the awaiter-status request.
- [`control-plane/internal/handlers/execute_awaiter_status.go:139`](https://github.com/Agent-Field/agentfield/blob/cbd2ddf7a14ac21524994c600ef513d163d919fe/control-plane/internal/handlers/execute_awaiter_status.go#L139) updates the lightweight execution record.
- [`control-plane/internal/handlers/execute_awaiter_status.go:155`](https://github.com/Agent-Field/agentfield/blob/cbd2ddf7a14ac21524994c600ef513d163d919fe/control-plane/internal/handlers/execute_awaiter_status.go#L155) updates the workflow execution row.

### What Does Not Feed the SDK Active-time Clock

The following signals are observable and may update persisted state, but they do not currently subtract time from the six-hour SDK watch/wait budget:

- Status callbacks and `progress` payloads.
- Workflow execution events.
- Structured execution logs.
- `app.note()` notes.
- Harness stdout/stderr activity.
- Node heartbeat or UI SSE heartbeat activity.

Status callbacks update persisted execution rows and publish event data, including a `progress` field, but `progress` is not part of the Python SDK watchdog's elapsed-time calculation.

- [`control-plane/internal/handlers/execute.go:474`](https://github.com/Agent-Field/agentfield/blob/cbd2ddf7a14ac21524994c600ef513d163d919fe/control-plane/internal/handlers/execute.go#L474) handles `/executions/:id/status`.
- [`control-plane/internal/handlers/execute.go:511`](https://github.com/Agent-Field/agentfield/blob/cbd2ddf7a14ac21524994c600ef513d163d919fe/control-plane/internal/handlers/execute.go#L511) updates the execution row.
- [`control-plane/internal/handlers/execute.go:614`](https://github.com/Agent-Field/agentfield/blob/cbd2ddf7a14ac21524994c600ef513d163d919fe/control-plane/internal/handlers/execute.go#L614) syncs workflow execution status.
- [`control-plane/internal/handlers/execute.go:623`](https://github.com/Agent-Field/agentfield/blob/cbd2ddf7a14ac21524994c600ef513d163d919fe/control-plane/internal/handlers/execute.go#L623) includes `progress` in event data.
- [`sdk/python/agentfield/agent.py:2508`](https://github.com/Agent-Field/agentfield/blob/cbd2ddf7a14ac21524994c600ef513d163d919fe/sdk/python/agentfield/agent.py#L2508) posts async status callbacks from the SDK.

Workflow execution events mirror agent-emitted events into the `executions` table. Existing rows are updated through `UpdateExecutionRecord`; a missing row creates both an execution row and a workflow execution row. The current handler does not use these events as SDK timer input.

- [`control-plane/internal/handlers/workflow_execution_events.go:32`](https://github.com/Agent-Field/agentfield/blob/cbd2ddf7a14ac21524994c600ef513d163d919fe/control-plane/internal/handlers/workflow_execution_events.go#L32) describes the mirror contract.
- [`control-plane/internal/handlers/workflow_execution_events.go:52`](https://github.com/Agent-Field/agentfield/blob/cbd2ddf7a14ac21524994c600ef513d163d919fe/control-plane/internal/handlers/workflow_execution_events.go#L52) creates missing records.
- [`control-plane/internal/handlers/workflow_execution_events.go:70`](https://github.com/Agent-Field/agentfield/blob/cbd2ddf7a14ac21524994c600ef513d163d919fe/control-plane/internal/handlers/workflow_execution_events.go#L70) updates existing execution records.

Structured execution logs are best-effort execution telemetry. The SDK prints a JSON line and dispatches it to the control plane, where it is stored in `execution_logs` and published to subscribers. The storage path inserts into `execution_logs`; it does not update the execution/workflow `updated_at` fields in the snippets below.

- [`sdk/python/agentfield/logger.py:176`](https://github.com/Agent-Field/agentfield/blob/cbd2ddf7a14ac21524994c600ef513d163d919fe/sdk/python/agentfield/logger.py#L176) emits a structured record.
- [`sdk/python/agentfield/client.py:1936`](https://github.com/Agent-Field/agentfield/blob/cbd2ddf7a14ac21524994c600ef513d163d919fe/sdk/python/agentfield/client.py#L1936) posts execution logs.
- [`control-plane/internal/handlers/execution_logs.go:32`](https://github.com/Agent-Field/agentfield/blob/cbd2ddf7a14ac21524994c600ef513d163d919fe/control-plane/internal/handlers/execution_logs.go#L32) handles log ingestion.
- [`control-plane/internal/storage/local.go:7769`](https://github.com/Agent-Field/agentfield/blob/cbd2ddf7a14ac21524994c600ef513d163d919fe/control-plane/internal/storage/local.go#L7769) stores a log entry.
- [`control-plane/internal/storage/local.go:7853`](https://github.com/Agent-Field/agentfield/blob/cbd2ddf7a14ac21524994c600ef513d163d919fe/control-plane/internal/storage/local.go#L7853) inserts the `execution_logs` row.

Notes update the execution row and broadcast an event, but they are separate from the SDK watchdog.

- [`sdk/python/agentfield/agent.py:4238`](https://github.com/Agent-Field/agentfield/blob/cbd2ddf7a14ac21524994c600ef513d163d919fe/sdk/python/agentfield/agent.py#L4238) implements `app.note()`.
- [`control-plane/internal/handlers/execution_notes.go:117`](https://github.com/Agent-Field/agentfield/blob/cbd2ddf7a14ac21524994c600ef513d163d919fe/control-plane/internal/handlers/execution_notes.go#L117) updates the execution row.
- [`control-plane/internal/handlers/execution_notes.go:156`](https://github.com/Agent-Field/agentfield/blob/cbd2ddf7a14ac21524994c600ef513d163d919fe/control-plane/internal/handlers/execution_notes.go#L156) publishes a note event.

### Control-plane Activity Timer

The stale cleanup service is separate from the six-hour SDK timer. It is the codebase's persisted "no activity" guard.

The contract is explicit: stale execution cleanup uses `updated_at` as the last activity timestamp, not `started_at`, and callers must bump it on every meaningful activity.

- [`control-plane/internal/storage/execution_records.go:979`](https://github.com/Agent-Field/agentfield/blob/cbd2ddf7a14ac21524994c600ef513d163d919fe/control-plane/internal/storage/execution_records.go#L979) documents `MarkStaleExecutions`.
- [`control-plane/internal/storage/execution_records.go:983`](https://github.com/Agent-Field/agentfield/blob/cbd2ddf7a14ac21524994c600ef513d163d919fe/control-plane/internal/storage/execution_records.go#L983) states the `updated_at` invariant.
- [`control-plane/internal/storage/execution_records.go:997`](https://github.com/Agent-Field/agentfield/blob/cbd2ddf7a14ac21524994c600ef513d163d919fe/control-plane/internal/storage/execution_records.go#L997) selects stale execution rows by `COALESCE(updated_at, created_at, started_at)`.
- [`control-plane/internal/storage/execution_records.go:1036`](https://github.com/Agent-Field/agentfield/blob/cbd2ddf7a14ac21524994c600ef513d163d919fe/control-plane/internal/storage/execution_records.go#L1036) marks stale execution rows as `timeout`.
- [`control-plane/internal/storage/execution_records.go:1088`](https://github.com/Agent-Field/agentfield/blob/cbd2ddf7a14ac21524994c600ef513d163d919fe/control-plane/internal/storage/execution_records.go#L1088) does the same for workflow executions.
- [`control-plane/internal/handlers/execution_cleanup.go:119`](https://github.com/Agent-Field/agentfield/blob/cbd2ddf7a14ac21524994c600ef513d163d919fe/control-plane/internal/handlers/execution_cleanup.go#L119) runs cleanup passes.
- [`control-plane/internal/handlers/execution_cleanup.go:134`](https://github.com/Agent-Field/agentfield/blob/cbd2ddf7a14ac21524994c600ef513d163d919fe/control-plane/internal/handlers/execution_cleanup.go#L134) uses `StaleExecutionTimeout`.

`CreateExecutionRecord` initializes `created_at` and `updated_at`; `UpdateExecutionRecord` unconditionally rewrites `updated_at` after a non-nil update callback returns.

- [`control-plane/internal/storage/execution_records.go:18`](https://github.com/Agent-Field/agentfield/blob/cbd2ddf7a14ac21524994c600ef513d163d919fe/control-plane/internal/storage/execution_records.go#L18) creates execution rows.
- [`control-plane/internal/storage/execution_records.go:30`](https://github.com/Agent-Field/agentfield/blob/cbd2ddf7a14ac21524994c600ef513d163d919fe/control-plane/internal/storage/execution_records.go#L30) sets `CreatedAt` and `UpdatedAt`.
- [`control-plane/internal/storage/execution_records.go:112`](https://github.com/Agent-Field/agentfield/blob/cbd2ddf7a14ac21524994c600ef513d163d919fe/control-plane/internal/storage/execution_records.go#L112) updates execution rows.
- [`control-plane/internal/storage/execution_records.go:154`](https://github.com/Agent-Field/agentfield/blob/cbd2ddf7a14ac21524994c600ef513d163d919fe/control-plane/internal/storage/execution_records.go#L154) rewrites `updated_at`.

Current persisted activity feeders include:

- Status callbacks: [`control-plane/internal/handlers/execute.go:474`](https://github.com/Agent-Field/agentfield/blob/cbd2ddf7a14ac21524994c600ef513d163d919fe/control-plane/internal/handlers/execute.go#L474)
- Approval waiting transitions: [`control-plane/internal/handlers/execute_approval.go:125`](https://github.com/Agent-Field/agentfield/blob/cbd2ddf7a14ac21524994c600ef513d163d919fe/control-plane/internal/handlers/execute_approval.go#L125)
- Awaiter waiting/running transitions: [`control-plane/internal/handlers/execute_awaiter_status.go:139`](https://github.com/Agent-Field/agentfield/blob/cbd2ddf7a14ac21524994c600ef513d163d919fe/control-plane/internal/handlers/execute_awaiter_status.go#L139)
- Pause/resume transitions: [`control-plane/internal/handlers/execute_pause.go:90`](https://github.com/Agent-Field/agentfield/blob/cbd2ddf7a14ac21524994c600ef513d163d919fe/control-plane/internal/handlers/execute_pause.go#L90)
- Cancellation: [`control-plane/internal/handlers/execute_cancel.go:72`](https://github.com/Agent-Field/agentfield/blob/cbd2ddf7a14ac21524994c600ef513d163d919fe/control-plane/internal/handlers/execute_cancel.go#L72)
- Completion/failure persistence: [`control-plane/internal/handlers/execute.go:1304`](https://github.com/Agent-Field/agentfield/blob/cbd2ddf7a14ac21524994c600ef513d163d919fe/control-plane/internal/handlers/execute.go#L1304) and [`control-plane/internal/handlers/execute.go:1377`](https://github.com/Agent-Field/agentfield/blob/cbd2ddf7a14ac21524994c600ef513d163d919fe/control-plane/internal/handlers/execute.go#L1377)
- Notes: [`control-plane/internal/handlers/execution_notes.go:117`](https://github.com/Agent-Field/agentfield/blob/cbd2ddf7a14ac21524994c600ef513d163d919fe/control-plane/internal/handlers/execution_notes.go#L117)

### Harness Subprocess Timers

The OpenCode harness has the closest existing "progress" timer. It treats stdout/stderr output as activity for one subprocess.

- [`sdk/python/agentfield/harness/providers/opencode.py:221`](https://github.com/Agent-Field/agentfield/blob/cbd2ddf7a14ac21524994c600ef513d163d919fe/sdk/python/agentfield/harness/providers/opencode.py#L221) sets a total wall-clock cap for one OpenCode subprocess.
- [`sdk/python/agentfield/harness/providers/opencode.py:228`](https://github.com/Agent-Field/agentfield/blob/cbd2ddf7a14ac21524994c600ef513d163d919fe/sdk/python/agentfield/harness/providers/opencode.py#L228) reads `AGENTFIELD_HARNESS_TIMEOUT_SECONDS`.
- [`sdk/python/agentfield/harness/providers/opencode.py:232`](https://github.com/Agent-Field/agentfield/blob/cbd2ddf7a14ac21524994c600ef513d163d919fe/sdk/python/agentfield/harness/providers/opencode.py#L232) documents the idle no-output cap.
- [`sdk/python/agentfield/harness/providers/opencode.py:242`](https://github.com/Agent-Field/agentfield/blob/cbd2ddf7a14ac21524994c600ef513d163d919fe/sdk/python/agentfield/harness/providers/opencode.py#L242) reads `AGENTFIELD_HARNESS_IDLE_TIMEOUT_SECONDS`.
- [`sdk/python/agentfield/harness/providers/opencode.py:250`](https://github.com/Agent-Field/agentfield/blob/cbd2ddf7a14ac21524994c600ef513d163d919fe/sdk/python/agentfield/harness/providers/opencode.py#L250) passes both values to `run_cli()`.
- [`sdk/python/agentfield/harness/_cli.py:19`](https://github.com/Agent-Field/agentfield/blob/cbd2ddf7a14ac21524994c600ef513d163d919fe/sdk/python/agentfield/harness/_cli.py#L19) reads pipes and stamps output activity.
- [`sdk/python/agentfield/harness/_cli.py:45`](https://github.com/Agent-Field/agentfield/blob/cbd2ddf7a14ac21524994c600ef513d163d919fe/sdk/python/agentfield/harness/_cli.py#L45) documents total cap and idle cap semantics.
- [`sdk/python/agentfield/harness/_cli.py:107`](https://github.com/Agent-Field/agentfield/blob/cbd2ddf7a14ac21524994c600ef513d163d919fe/sdk/python/agentfield/harness/_cli.py#L107) kills when no output arrives for `idle_timeout`.
- [`sdk/python/agentfield/harness/_cli.py:110`](https://github.com/Agent-Field/agentfield/blob/cbd2ddf7a14ac21524994c600ef513d163d919fe/sdk/python/agentfield/harness/_cli.py#L110) kills when total timeout is exceeded.

Other harness runtimes have timeout surfaces, but not the same idle-output behavior in the snippets reviewed:

- [`sdk/python/agentfield/harness/providers/codex.py:50`](https://github.com/Agent-Field/agentfield/blob/cbd2ddf7a14ac21524994c600ef513d163d919fe/sdk/python/agentfield/harness/providers/codex.py#L50) calls `run_cli()` without timeout arguments.
- [`sdk/python/agentfield/harness/providers/gemini.py:46`](https://github.com/Agent-Field/agentfield/blob/cbd2ddf7a14ac21524994c600ef513d163d919fe/sdk/python/agentfield/harness/providers/gemini.py#L46) calls `run_cli()` without timeout arguments.
- [`sdk/go/harness/provider.go:65`](https://github.com/Agent-Field/agentfield/blob/cbd2ddf7a14ac21524994c600ef513d163d919fe/sdk/go/harness/provider.go#L65) exposes a Go harness timeout option.
- [`sdk/go/harness/provider.go:120`](https://github.com/Agent-Field/agentfield/blob/cbd2ddf7a14ac21524994c600ef513d163d919fe/sdk/go/harness/provider.go#L120) defaults the Go harness timeout to 600 seconds.
- [`sdk/typescript/src/harness/cli.ts:31`](https://github.com/Agent-Field/agentfield/blob/cbd2ddf7a14ac21524994c600ef513d163d919fe/sdk/typescript/src/harness/cli.ts#L31) sets an optional TypeScript CLI timeout.

### Cancellation Path

When a control-plane cancel is requested, the control plane updates the execution row and publishes a cancel event. The `CancelDispatcher` best-effort POSTs to the worker-side internal cancel route, which cancels the registered asyncio task.

- [`control-plane/internal/handlers/execute_cancel.go:31`](https://github.com/Agent-Field/agentfield/blob/cbd2ddf7a14ac21524994c600ef513d163d919fe/control-plane/internal/handlers/execute_cancel.go#L31) handles execution cancellation.
- [`control-plane/internal/handlers/execute_cancel.go:73`](https://github.com/Agent-Field/agentfield/blob/cbd2ddf7a14ac21524994c600ef513d163d919fe/control-plane/internal/handlers/execute_cancel.go#L73) updates execution status to cancelled.
- [`control-plane/internal/services/cancel_dispatcher.go:76`](https://github.com/Agent-Field/agentfield/blob/cbd2ddf7a14ac21524994c600ef513d163d919fe/control-plane/internal/services/cancel_dispatcher.go#L76) constructs the cancel dispatcher.
- [`control-plane/internal/services/cancel_dispatcher.go:134`](https://github.com/Agent-Field/agentfield/blob/cbd2ddf7a14ac21524994c600ef513d163d919fe/control-plane/internal/services/cancel_dispatcher.go#L134) only handles cancel events.
- [`control-plane/internal/services/cancel_dispatcher.go:202`](https://github.com/Agent-Field/agentfield/blob/cbd2ddf7a14ac21524994c600ef513d163d919fe/control-plane/internal/services/cancel_dispatcher.go#L202) builds the worker cancel URL.
- [`sdk/python/agentfield/cancel.py:38`](https://github.com/Agent-Field/agentfield/blob/cbd2ddf7a14ac21524994c600ef513d163d919fe/sdk/python/agentfield/cancel.py#L38) registers in-flight execution tasks.
- [`sdk/python/agentfield/cancel.py:65`](https://github.com/Agent-Field/agentfield/blob/cbd2ddf7a14ac21524994c600ef513d163d919fe/sdk/python/agentfield/cancel.py#L65) cancels the matching task.
- [`sdk/python/agentfield/cancel.py:92`](https://github.com/Agent-Field/agentfield/blob/cbd2ddf7a14ac21524994c600ef513d163d919fe/sdk/python/agentfield/cancel.py#L92) installs `/_internal/executions/{execution_id}/cancel`.

## Architecture Documentation

### Current Contracts

1. SDK active-time contract

   The reasoner watchdog and `Agent.call()` wait logic use local active time, defined as wall-clock elapsed minus pause-clock time. The only current pause-clock inputs are explicit `WAITING` states created by approvals or child-wait propagation.

2. Control-plane no-activity contract

   The control plane uses persisted `updated_at` as last activity for stale cleanup. This is a storage-level contract and is independent of the SDK's local watch/wait budget.

3. Visibility contract

   Logs, notes, workflow events, status updates, and UI streams expose execution activity to humans and API consumers. They do not currently feed the SDK active-time clock.

4. Harness subprocess contract

   The OpenCode provider has a subprocess-level idle detector where stdout/stderr bytes reset last-output activity. This exists below the reasoner execution layer and does not reset the SDK's six-hour wait/watch budget.

### Seams Relevant to Progress-aware Behavior

These are the existing integration points where progress-like information already crosses a boundary:

| Seam | Existing contract | Current relation to 6-hour SDK budget |
| --- | --- | --- |
| `PauseClock` | local active-time subtraction | Direct input |
| `Agent.pause()` / `wait_for_resume()` | approval wait state | Direct input through `PauseClock` |
| `Agent.call()` child WAITING observation | parent pause-clock pause | Direct input through `PauseClock` |
| `notify_awaiter_status()` | multi-hop WAITING/RUNNING propagation | Indirect input by making ancestors see WAITING |
| `/executions/:id/status` | status/result/error/progress callback | Persists activity and publishes events; not an SDK timer input |
| `/workflow/executions/events` | event mirror into executions | Persists activity for execution rows; not an SDK timer input |
| `/executions/:id/logs` | structured log ingestion/streaming | Stores logs and streams to UI; not an SDK timer input |
| `/executions/note` | app note with execution context | Updates execution row; not an SDK timer input |
| OpenCode `run_cli(... idle_timeout=...)` | stdout/stderr idle detector | Applies to one subprocess; not an execution-layer timer input |
| stale cleanup `updated_at` | persisted no-activity reaper | Separate control-plane timer |

## Historical Context

- `thoughts/shared/research/2026-06-09-07-42-run-state-persistence-and-cancel-cascade.md` documented a production run where `build` and `execute` exhausted a 21600-second active-time budget. It identified the SDK active-time watchdog, the `AGENTFIELD_ASYNC_DEFAULT_EXECUTION_TIMEOUT`/harness timeout environment surfaces, and the distinction between generic cancellation logs and the authoritative reasoner timeout callback.
- `thoughts/shared/research/2026-06-16-swe-af-coding-loop-false-green-bug.md` is related historical context for long-running SWE-AF coding loops.
- `thoughts/shared/handoffs/general/2026-06-16_11-14-06_agentplane-port-codex-ui-wiring.md` noted that the AgentPlane port had its own callback timeout wiring from `AGENTFIELD_ASYNC_DEFAULT_EXECUTION_TIMEOUT`.
- `thoughts/shared/handoffs/general/2026-06-24_08-44-25_agentplane-ui-provenance-swe-af-debrand.md` recorded later AgentPlane/SWE-AF run-state and UI streaming context.

## Related Research

- `thoughts/shared/research/2026-06-09-07-42-run-state-persistence-and-cancel-cascade.md`
- `thoughts/shared/research/2026-06-16-swe-af-coding-loop-false-green-bug.md`

## Open Questions

- Which deployment path sets `AGENTFIELD_ASYNC_DEFAULT_EXECUTION_TIMEOUT` to `21600` for the reasoner watchdog is not visible in this repo.
- Whether current build reasoners emit `/executions/:id/status` progress callbacks, only structured logs, only harness stdout/stderr, or a mixture depends on the build implementation outside the generic execution contracts mapped here.
- There is no dedicated build heartbeat or progress contract in this repo; all current activity signals are expressed through generic execution/status/log/note/event interfaces.
