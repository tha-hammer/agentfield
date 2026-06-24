---
date: 2026-06-09T07:42:08-04:00
researcher: tha-hammer
git_commit: 7ea5107d3ef4973ae3e14b9ff47770b892cc6503
branch: fix/harness-retry-preserve-goal
repository: agentfield
topic: "How a run's state is saved (or not), and the interfaces/seams in the per-call-timeout → build-cancel cascade"
tags: [research, codebase, swe-af, harness, dag-executor, checkpoint, async-cancellation, control-plane, persistence]
status: complete
last_updated: 2026-06-09
last_updated_by: tha-hammer
last_updated_note: "Resolved the cancel-cascade trigger against the live control-plane DB: build AND execute both hit the 21600s active-time budget (watchdog), NOT the cancel route and NOT a connection drop. Corrected the cancellation-seam section and Open Questions."
---

# Research: Run-State Persistence and the Per-Call-Timeout → Build-Cancel Cascade

> **Follow-up sections appended below** — see "Follow-up: Why no replan fired during the 6 hours" (2026-06-09).

**Date**: 2026-06-09T07:42:08-04:00
**Researcher**: tha-hammer
**Git Commit**: 7ea5107d3ef4973ae3e14b9ff47770b892cc6503
**Branch**: fix/harness-retry-preserve-goal
**Repository**: agentfield (with cross-repo references to `/home/maceo/Dev/SWE-AF-anthropic`)

## Research Question

How is a run's state saved — or, where it is not saved, where *could* it be saved? Map all interfaces and seams involved, in the context of this observed scenario:

> "the issue-12 router-wiring turn got stuck for ~28m, hit the 60-min per-call cap, returned empty JSON, and that cascaded to a build cancel."

## Summary

A SWE-AF build is a tree of AgentField executions. Run state is persisted at **four independent layers**, each with its own store, its own write cadence, and its own (non-)resume behavior:

1. **Git** (`<repo>/.worktrees/` + branches + the integration/feature branch) — the only layer that survives a build cancel intact, because git writes are made by agent subprocesses and are never rolled back. The merged controllers from the observed run survived here.
2. **Disk artifacts under `<repo>/<artifacts_dir>/`** — `DAGState` checkpoint (`execution/checkpoint.json`), per-issue iteration checkpoints (`execution/iterations/<build_id>/<issue>.json`), per-call agent artifacts (`coding-loop/<id>/*.json`), and the planning docs (`prd.md`, `architecture.md`, `issues/*.md`). These are written by SWE-AF Python code and by the agents' `Write` tool.
3. **Control-plane database** (`executions`, `workflow_executions`, `workflow_runs`, `workflow_execution_events`, `execution_logs` tables) — status, payloads, and an append-only event log per execution node.
4. **In-process Python objects** (`DAGState`, `BuildResult`, the `_shared_memory` dict) — RAM-only, lost on cancel except where flushed to layer 2.

The central finding for the scenario: **the checkpoint machinery exists and is written, but the `build` reasoner never reloads it.** `execute()` accepts `resume: bool = False` (`swe_af/app.py:1549`), and `run_dag()` honors it (`dag_executor.py:1412-1425`), but `build()`'s call into `execute()` (`app.py:933-940`) does not pass `resume`, so a re-fired build always starts a fresh `DAGState`. The per-issue iteration checkpoint *does* auto-resume regardless of a flag (`coding_loop.py:589-600`), but only within a single live `run_dag()` process. Across a process-level cancel, the only durable carry-over is git.

The cascade itself crosses three seams: (a) the **harness per-call timeout** (`AGENTFIELD_HARNESS_TIMEOUT_SECONDS`, enforced by `asyncio.wait_for` in `harness/_cli.py:39-41`) which kills the opencode subprocess and yields a `FailureType.TIMEOUT` result; (b) the **child-reasoner error-isolation** in `run_coder` which normally swallows a failed harness result into a fallback dict (so a timeout alone does *not* cancel the build); and (c) the **async cancellation seam** (`agent.py:2247` / `agent.py:2875`) where an `asyncio.CancelledError` propagating into the build reasoner's endpoint emits exactly the string `"Execution cancelled by upstream client"`. That `CancelledError` is raised either by the active-time watchdog (`agent.py:2379`) or by the control-plane cancel route calling `task.cancel()` (`cancel.py:78`). **For the observed run, the live DB shows it was the watchdog: both `build` and `execute` exhausted the 21600s (6h) active-time budget** (`build` `duration_ms=21600428`, exact). The hardcoded `"Execution cancelled by upstream client"` string at `agent.py:2247` is a generic side-effect of `task.cancel()` and does not indicate the cancel route or a client disconnect — the authoritative `error_message` is `"Reasoner 'build' timed out after 21600.0s of active time"`. See the cancellation-seam section for the full table.

## Detailed Findings

### Layer 1 — Git as durable run state (SWE-AF)

Git is the de-facto durable store for completed work. `DAGState` only tracks pointers; the bytes live in `.git`.

- Worktrees dir is `os.path.join(repo_path, ".worktrees")`, set in `_init_dag_state` (`swe_af/execution/dag_executor.py:755`), stored on `DAGState.worktrees_dir` (`swe_af/execution/schemas.py:315`).
- **Create**: `run_workspace_setup` reasoner (`swe_af/reasoners/execution_agents.py:682`) runs `git worktree add <dir>/issue-<BUILD_ID>-<NN>-<name> -b issue/<BUILD_ID>-<NN>-<name> <integration_branch>` via an LLM agent with `tools=["Bash","Write"]`; command template in `swe_af/prompts/workspace.py:29-32`. Called per level by `_setup_worktrees` (`dag_executor.py:54`, invoked at `dag_executor.py:1491`).
- **Merge**: `run_merger` reasoner (`execution_agents.py:749`) runs `git merge <branch> --no-ff -m "Merge <branch>: <title>"` one branch at a time; template in `swe_af/prompts/merger.py:18-19`. Only issues in `level_result.completed` with a non-empty `branch_name` are passed (`dag_executor.py:229-240`); `completed` includes both `COMPLETED` and `COMPLETED_WITH_DEBT` (`dag_executor.py:1192-1193`).
- **Result tracking**: merged names → `dag_state.merged_branches` (`dag_executor.py:276-278`); failed → `dag_state.unmerged_branches` (`dag_executor.py:281-283`); raw dict → `dag_state.merge_results` (`dag_executor.py:275`). `MergeResult` carries `merge_commit_sha` and `pre_merge_sha` (`schemas.py:354-366`).
- **Teardown**: `run_workspace_cleanup` (`execution_agents.py:897`) runs `git worktree remove --force` + `git branch -D` + `git worktree prune` (`workspace.py:63-73`), launched as a background task after merge (`dag_executor.py:1588-1596`) and awaited before the next level (`dag_executor.py:1757`); a final sweep runs at `dag_executor.py:1763-1784`.

**Cancel behavior (the seam that saved the observed run):** there is no `try/finally`, signal handler, or cancellation hook in `dag_executor.py` or `coding_loop.py` that runs cleanup on process termination. Cleanup only runs after *successful* level completion. So on a mid-issue cancel, the in-flight worktree and branch persist on disk — and so does every previously-merged branch on the integration/feature branch. This is exactly why the 8 merged controllers survived while the in-flight issue-12 worktree was left orphaned.

### Layer 2 — Disk artifacts under `<artifacts_dir>` (SWE-AF)

Directory layout for a build (`artifacts_dir` defaults to `.artifacts`, resolved absolute), per `_ensure_paths` (`swe_af/reasoners/pipeline.py:36`):

```
<repo_path>/<artifacts_dir>/
  plan/
    prd.md            ← written by PM agent via Write tool
    architecture.md   ← written by Architect agent via Write tool
    review.json       ← written by Python (pipeline.py:482-484)
    issues/issue-<seq>-<name>.md  ← written by issue-writer agents
  rationale.md        ← written by Python (app.py:1525-1526)
  logs/replanner_<n>_raw_<attempt>.txt  ← (execution_agents.py:382)
  execution/
    checkpoint.json                       ← DAGState snapshot
    iterations/<build_id>/<issue>.json    ← per-issue iteration checkpoint
  coding-loop/<iteration_id>/
    coder.json | review.json | qa.json | synthesis.json  ← per-call artifacts
  approval_state.json  ← HITL only (app.py:789, 804)
```

**`DAGState` checkpoint** — the build's central recoverable state object:
- `_checkpoint_path` → `<artifacts_dir>/execution/checkpoint.json` (`dag_executor.py:683-685`).
- `_save_checkpoint` → `json.dump(dag_state.model_dump(), ..., default=str)` (`dag_executor.py:688-697`).
- `_load_checkpoint` → `DAGState(**json.load(...))` (`dag_executor.py:700-706`).
- Save points inside `run_dag`: initial (`:1436`), pre-level with `in_flight_issues` populated (`:1505`), post-level barrier with `in_flight_issues` cleared (`:1517`), post-split (`:1676`), post-replan (`:1734`), level-failure abort (`:1557`), final (`:1799`).
- `DAGState` fields (`schemas.py:276`) include `completed_issues`, `failed_issues`, `in_flight_issues`, `current_level`, `merged_branches`/`unmerged_branches`/`pending_merge_branches`, `build_id`, `accumulated_debt`, and the full `all_issues` list (mutated in place with `worktree_path`/`branch_name`).

**Per-issue iteration checkpoint** (`coding_loop.py`):
- `_iteration_state_path` → `<artifacts_dir>/execution/iterations/<build_id>/<issue>.json` (`coding_loop.py:47-56`).
- `_save_iteration_state` writes `{iteration, feedback, files_changed, iteration_history}` after every iteration (`coding_loop.py:730-735`).
- `_load_iteration_state` is called unconditionally at issue start (`coding_loop.py:589-600`): if the file exists, `start_iteration = existing_state["iteration"] + 1` and the coder loop resumes mid-issue. **This resume is independent of the `resume` flag** — but it only helps within a single live `run_dag()` process, since it keys on the in-memory `dag_state.build_id` and the iteration files are under the same artifacts dir.

**Per-call agent artifacts**: `_save_artifact` (`coding_loop.py:76-85`) writes `coder.json`/`review.json`/`qa.json`/`synthesis.json` per `iteration_id`; append-only, never read back.

`IssueResult` (`schemas.py:222`) — outcome, `branch_name`, `files_changed`, `attempts`, `iteration_history`, `debt_items` — is stored *inside* `DAGState.completed_issues`/`failed_issues` and serialized only via the `checkpoint.json`, not as separate per-issue result files.

### Layer 3 — Control-plane execution persistence (Go)

Every reasoner call (build, plan, execute, run_coder, run_merger, …) is an execution row in the control plane. Two parallel tables share `execution_id`:

- `executions` / `types.Execution` (`control-plane/internal/storage/models.go:5-30`, `control-plane/pkg/types/execution.go:10-48`): `execution_id` (PK), `run_id`, `parent_execution_id`, `status`, `status_reason`, `error_message`, `started_at`/`completed_at`/`duration_ms`, `input_payload`/`result_payload` (raw JSON), `session_id`, `notes`.
- `workflow_executions` (`models.go:115-162`): adds `workflow_id`, `root_workflow_id`/`parent_workflow_id`, `workflow_depth`, `state_version` (optimistic concurrency), `last_event_sequence`, `active_children`/`pending_children`, `pending_terminal_status`, `lease_owner`/`lease_expires_at`, `retry_count`, and HITL approval columns.
- `workflow_runs` (`models.go:224-240`): the parent grouping per `run_id` — `status`, `total_steps`/`completed_steps`/`failed_steps`, `state_version`.
- `workflow_execution_events` (`models.go:164-180`): append-only per-execution event log (`sequence`, `event_type`, `status`, `payload`).
- `execution_logs` (`models.go:182-207`): append-only structured log lines.

**Status model** (`control-plane/pkg/types/status.go:8-19`): `unknown, pending, queued, waiting, running, paused, succeeded, failed, cancelled, timeout`. Terminal = `succeeded|failed|cancelled|timeout` (`status.go:76-83`). Transition table at `control-plane/internal/storage/execution_state_validation.go:22-57` (notably `running → {waiting, paused, succeeded, failed, cancelled, timeout}`).

**Lifecycle** (`control-plane/internal/handlers/execute.go`): `prepareExecution` creates the row with `status=running` (`execute.go:1119,1151`); `completeExecution` sets `succeeded` (`execute.go:1304`); `failExecution` sets `failed` + `error_message` + `status_reason` (`execute.go:1377`); async path blocks in `waitForExecutionCompletion` on the event bus (`execute.go:856`).

**Pause/resume + checkpoint**: pause is persisted simply as `status=paused` in both tables (no separate checkpoint table); `PauseExecutionHandler`/`ResumeExecutionHandler` → `handlePauseResume` (`execute_pause.go:38-50`). The control-plane goroutine detects a paused row in `callAgent` (`execute.go:1223-1232`) and blocks in `waitForResume` on the event bus until `ExecutionResumed` or `ExecutionCancelledEvent` (`execute.go:935-991`).

**Stale/orphan reaping**: `MarkStaleExecutions` → `timeout` (`execution_records.go:986`); `MarkAgentExecutionsOrphaned` → `failed` for an agent that re-registers with a new `instance_id` (`execution_records.go:1234-1270`); `RetryStaleWorkflowExecutions` resets to `pending` under `retry_count < max` (`execution_records.go:1272`).

### Layer 4 — In-process state and AgentField memory (not used for build state)

- `DAGState`, `BuildResult` (`schemas.py:851`), and `PlanResult` (`reasoners/schemas.py:96`) are passed by reference between reasoner calls in-process; only `DAGState` is flushed to disk (layer 2).
- The cross-issue learning store is an in-process dict `_shared_memory = {}` created at `dag_executor.py:1452` and accessed via the `_memory_fn` closure (`dag_executor.py:1454-1457`). It is RAM-only.
- **The AgentField SDK memory API (global/agent/session/run scopes) is not used anywhere in the build pipeline.** On the control-plane side those scopes exist (`control-plane/internal/handlers/memory.go:402-434` `resolveScope`; BoltDB buckets `workflow/session/actor/reasoner/global` at `local.go:943`, keyed `{scope_id}:{key}`; Postgres `kv_store` PK `(scope, scope_id, key)`), but SWE-AF never writes build state into them.

### The harness per-call seam (`AGENTFIELD_HARNESS_TIMEOUT_SECONDS`)

This is the "60-min cap" in the scenario (the env was overridden to 3600s; SDK default is 1800s).

- Read in the opencode provider: `int(os.environ.get("AGENTFIELD_HARNESS_TIMEOUT_SECONDS", "1800"))` (`sdk/python/agentfield/harness/providers/opencode.py:228-230`), passed as `timeout=` to `run_cli()` (`opencode.py:236-237`).
- Enforced via `asyncio.wait_for(proc.communicate(), timeout=timeout)` (`sdk/python/agentfield/harness/_cli.py:39-41`). On fire: `proc.kill()` then re-raise as `TimeoutError(f"CLI command timed out after {timeout}s: …")` (`_cli.py:43-45`).
- Mapped to a result: `RawResult(is_error=True, failure_type=FailureType.TIMEOUT, error_message=…)` (`opencode.py:249-255`). `FailureType` enum at `sdk/python/agentfield/harness/_result.py:8-25`.
- **Scope note**: this wall-clock cap is implemented only for the CLI-subprocess providers (opencode, and codex/gemini which actually pass `timeout=None`). `ClaudeCodeProvider` has no `asyncio.wait_for` wrapper (`harness/providers/claude.py:89`).

**Structured output and the "empty JSON":**
- Output file is `<cwd>/.agentfield_output.json` (`harness/_schema.py:22,29-31`). It is written by the agent subprocess (instructed via `build_prompt_suffix`, `_schema.py:69-99`), read by the harness via `parse_and_validate` (`_schema.py:185-206`).
- When the subprocess is killed, the file is typically absent or empty. `read_and_parse`/`read_repair_and_parse` guard `if not content.strip(): return None` *before* `json.loads` (`_schema.py:148-149, 160-161`), so they silently return `None` rather than raising.
- **Citation correction (verified):** the exact string `"Failed to parse JSON response: Expecting value: line 1 column 1 (char 0)"` does **not** exist in the harness source. `Expecting value: line 1 column 1 (char 0)` is CPython's standard message for `json.loads("")`; the harness explicitly guards against that path. If that message appeared in the build logs it originated from a caller doing `json.loads` directly (e.g. a Pydantic validator or the control-plane response parsing), not from the harness parse pipeline. This is labeled here rather than asserted as a harness behavior.
- **No durable per-call state:** `.agentfield_output.json` and `.agentfield_schema.json` are deleted in the `run()` `finally` via `cleanup_temp_files` (`harness/_runner.py:279-280`, `_schema.py:264-271`). The per-call opencode `XDG_DATA_HOME` tempdir is `rmtree`'d in `finally` (`opencode.py:217,259-260`). Retry transcripts (`all_raws`) live only in memory (`_runner.py:346,455`).

**Why a timeout alone does not cancel the build** (`run_coder` isolation): `run_coder` (`swe_af/reasoners/execution_agents.py:1005-1041`) calls `router.harness(...)`, runs `check_fatal_harness_error(result)`, and if `result.parsed is None` (the timeout case) falls through to `return CoderResult(files_changed=[], complete=False, …).model_dump()` — a non-raising fallback dict. So a single timed-out coder call yields `complete=False`, the control plane marks that child `succeeded` (it returned normally), and the build continues into the next coding-loop iteration. The stuck-loop / exhausted-loop logic in `coding_loop.py:808-887` then governs whether the issue ends `COMPLETED_WITH_DEBT` or `FAILED_UNRECOVERABLE`.

### The async cancellation seam (where the cascade string is emitted)

The string `"Execution cancelled by upstream client"` is emitted at exactly two sites, both inside `except asyncio.CancelledError` blocks that call `notify_call_error` then re-raise:

- **Site 1** — `_execute_reasoner_endpoint` (the `@app.reasoner()` request handler), `sdk/python/agentfield/agent.py:2247-2259` (verified by direct read). This is the handler the `build` reasoner runs in.
- **Site 2** — the `@app.skill()` endpoint closure, `agent.py:2875-2887` (verified by direct read).

`notify_call_error` (`agent_workflow.py:177-211`) sets `status="failed"`, `event_type="reasoner.failed"`, and fires `POST {server}/api/v1/workflow/executions/events`, then the `raise cancel_err` lets cooperative cancellation propagate.

**Who raises the `CancelledError` into the build reasoner — two mechanisms:**

1. **Active-time watchdog** (`_execute_async_with_callback`, `agent.py:2314-2458`): a `PauseClock` (`agent_pause.py:10-46`) tracks paused seconds; a watchdog computes `active_elapsed = (now - start) - pause_clock.total_paused()` every `min(5, timeout/4)`s and, when `active_elapsed > reasoner_timeout`, sets `pause_clock.timed_out=True` and calls `reasoner_task.cancel()` (`agent.py:2351-2382`). `reasoner_timeout` = `async_config.default_execution_timeout` (`agent.py:2335-2337`), default `7200.0`, env `AGENTFIELD_ASYNC_DEFAULT_EXECUTION_TIMEOUT` (overridden to `21600` in `SWE-AF-anthropic/docker-compose.yml:31`), max `21600.0` (`async_config.py:33-34`). When this fires, the outer handler sees `pause_clock.timed_out==True` and posts `status="failed"` with `"Reasoner '<name>' timed out after <t>s of active time"` (`agent.py:2395-2426`) — a *different* string from the upstream-client one.

2. **Control-plane cancel route** (`cancel.py:65-112`): `POST /_internal/executions/{id}/cancel` → `cancel_execution` → `task.cancel()` on the task registered via `register_execution_task` (`agent.py:1859-1874` async path / `1888-1907` sync path). When the control plane cancels the build execution — e.g. because the upstream HTTP client/connection that invoked the build went away (SSE/HTTP `Request.Context().Done()` on the control-plane side, `reasoner_catalog.go:130-196`; or `waitForResume`/disconnect paths marking the row), or an explicit cancel — `task.cancel()` injects `CancelledError` into the build reasoner's `await func(...)` (`agent.py:2208`), caught at Site 1, emitting `"Execution cancelled by upstream client"`.

**Resolved against the live control-plane DB (2026-06-09, run `run_20260608_202801_560683uw`, control plane `swe-af-anthropic-control-plane-1` :8090, local SQLite).** The authoritative execution records show mechanism **1 (active-time watchdog)** is what fired — for *both* the build reasoner and its execute child, each independently exhausting the 21600s budget:

| Execution | Reasoner | `started_at` → `completed_at` | `duration_ms` | `status` | `error` (authoritative) |
|---|---|---|---|---|---|
| `exec_20260608_202801_lo8k200j` | `build` | 20:28:01Z → 02:28:02Z | **21600428** (6h 0m 0.428s) | failed | `Reasoner 'build' timed out after 21600.0s of active time` |
| `exec_20260608_204927_eyc4ugq1` | `execute` | 20:49:27Z → 02:49:30Z | 21603319 (≈6h 0m 3s) | failed | `Reasoner 'execute' timed out after 21600.0s of active time` |

`build`'s `duration_ms` is `21600428` — the 6-hour active-time budget to the millisecond; the ~0.4–3s overage matches the watchdog's `min(5, timeout/4)=5s` poll interval firing just past the threshold. So the `build` task was cancelled by **its own active-time watchdog** (`agent.py:2351-2382` → `reasoner_task.cancel()`), and the outer handler recorded the authoritative `error_message = "Reasoner 'build' timed out after 21600.0s of active time"` via the timeout branch (`agent.py:2395-2426`, `pause_clock.timed_out == True`).

**The `"Execution cancelled by upstream client"` string that appears in the worker logs is a red herring.** It is the *inner, generic* `notify_call_error` emitted at `agent.py:2247` during the same `task.cancel()` — it fires for *any* `CancelledError` and carries no information about the trigger. For these executions the outer timeout-branch status post overwrote the recorded `error_message` with the budget-timeout message. Neither the control-plane cancel route (`cancel.py:78`) nor the `CancelDispatcher` was involved (the dispatcher acts only on `ExecutionCancelledEvent`; these executions flipped to `failed`, not `cancelled`). There was no "upstream connection going away" — both executions ran their full 6-hour active-time budget to completion of the watchdog.

**Child→parent error propagation** (`agent.py:3912-4069`, `async_execution_manager.py:525-727`): a child reasoner that the control plane marks `FAILED` surfaces in the parent as `ExecutionFailedError` (`async_execution_manager.py:649-651`), re-raised immediately out of `app.call()` (`agent.py:4027-4069`), bypassing `_unwrap` (`swe_af/execution/envelope.py:25-66`) and propagating as an unhandled exception through `build()` → caught by the generic `except Exception` at `agent.py:2296` → `notify_call_error(str(exc))` (a failure, not a cancel). This is the path for a *raised* child failure; the `run_coder` fallback-dict path (above) avoids it entirely.

## Code References

State stores:
- `swe_af/execution/dag_executor.py:683-706` — `_checkpoint_path` / `_save_checkpoint` / `_load_checkpoint`
- `swe_af/execution/dag_executor.py:1408-1436` — resume-from-checkpoint gate + initial save
- `swe_af/execution/coding_loop.py:47-85` — iteration-state + artifact persistence helpers
- `swe_af/execution/coding_loop.py:589-600` — unconditional per-issue iteration resume
- `swe_af/execution/coding_loop.py:730-735` — per-iteration checkpoint write
- `swe_af/execution/schemas.py:276` — `DAGState`; `:222` `IssueResult`; `:149-158` `IssueOutcome`; `:851` `BuildResult`; `:1002`/`:702` `max_coding_iterations=5`
- `swe_af/reasoners/schemas.py:96` — `PlanResult`

Build orchestration / resume seam:
- `swe_af/app.py:458` — `@app.reasoner() build`; `:933-940` `app.call(execute, …)` **without `resume`**
- `swe_af/app.py:1542-1552` — `execute(... resume: bool = False ...)`; `:1550` passes `resume` to `run_dag`

Worktree / merge lifecycle:
- `swe_af/reasoners/execution_agents.py:682` `run_workspace_setup`; `:749` `run_merger`; `:897` `run_workspace_cleanup`
- `swe_af/execution/dag_executor.py:54,206,485,755,1491,1564,1588-1596,1757,1763-1784`
- `swe_af/prompts/workspace.py:29-32,63-73`; `swe_af/prompts/merger.py:18-19`

Coding-loop outcome paths:
- `swe_af/execution/coding_loop.py:808-829` (stuck → `COMPLETED_WITH_DEBT`), `:849-874` (exhausted → `COMPLETED_WITH_DEBT`), `:839-847`/`:887` (`FAILED_UNRECOVERABLE`)

Harness per-call seam:
- `sdk/python/agentfield/harness/providers/opencode.py:228-230,236-237,249-255` — timeout read/enforce/result
- `sdk/python/agentfield/harness/_cli.py:39-45` — `asyncio.wait_for` + kill + re-raise
- `sdk/python/agentfield/harness/_result.py:8-25` — `FailureType`
- `sdk/python/agentfield/harness/_schema.py:22,29-31,143-165,185-206,264-271` — output file, parse, cleanup
- `sdk/python/agentfield/harness/_runner.py:279-280` — `cleanup_temp_files` in `finally`
- `swe_af/reasoners/execution_agents.py:1005-1041` — `run_coder` fallback-dict isolation

Async cancellation seam:
- `sdk/python/agentfield/agent.py:2247-2259` — Site 1 `"Execution cancelled by upstream client"` (verified)
- `sdk/python/agentfield/agent.py:2875-2887` — Site 2 (verified)
- `sdk/python/agentfield/agent.py:2314-2458` — `_execute_async_with_callback` watchdog + timeout/cancel distinction
- `sdk/python/agentfield/agent.py:2335-2337,2351-2382,2395-2426` — active-time budget computation and branch
- `sdk/python/agentfield/agent.py:1859-1874,1888-1907` — task registration for cancel
- `sdk/python/agentfield/cancel.py:65-112` — cancel route → `task.cancel()`
- `sdk/python/agentfield/agent_workflow.py:177-211` — `notify_call_error` → POST events
- `sdk/python/agentfield/agent_pause.py:10-46` — `PauseClock`
- `sdk/python/agentfield/async_config.py:33-34,140-143,217-219` — execution-timeout config
- `swe_af/docker-compose.yml:29-31` — `AGENTFIELD_HARNESS_TIMEOUT_SECONDS=3600`, `AGENTFIELD_ASYNC_DEFAULT_EXECUTION_TIMEOUT=21600`

Control-plane persistence:
- `control-plane/internal/storage/models.go:5-30,115-162,164-180,182-207,224-240`
- `control-plane/pkg/types/execution.go:10-48`; `control-plane/pkg/types/status.go:8-19,76-83`
- `control-plane/internal/storage/execution_state_validation.go:22-57`
- `control-plane/internal/handlers/execute.go:856,1119,1151,1223-1233,1304,1377`; `execute_cancel.go:31-160`; `execute_pause.go:38-98`
- `control-plane/internal/events/execution_events.go:42-114,133-231`
- `control-plane/internal/handlers/reasoner_catalog.go:130-196` — SSE; client disconnect = `Request.Context().Done()`
- `control-plane/internal/handlers/memory.go:402-434`; `control-plane/internal/storage/local.go:943,3965-4076`

## Architecture Documentation

**Persistence cadence by layer** (what is on disk/DB at any instant):

| Layer | Store | Written when | Survives process cancel? | Reloaded on re-fire? |
|---|---|---|---|---|
| Git worktrees/branches | `<repo>/.git`, `.worktrees/` | agent `Bash` git commands | Yes (no cancel cleanup) | Manual (branch persists; build re-clones) |
| `DAGState` checkpoint | `<artifacts_dir>/execution/checkpoint.json` | level boundaries + replan/split/final | Yes (last boundary only) | **No** — `build()` never passes `resume=True` |
| Iteration checkpoint | `…/execution/iterations/<build_id>/<issue>.json` | after every coder iteration | Yes | Only within a live `run_dag()` (same `build_id`) |
| Per-call artifacts | `…/coding-loop/<id>/*.json` | each agent call | Yes | Never read back |
| Planning docs | `…/plan/*.md`, `review.json` | during planning | Yes | Re-generated by a fresh build |
| Control-plane rows | `executions`/`workflow_executions`/… | each reasoner lifecycle event | Yes (status→failed/cancelled) | Read for status/audit, not for build resume |
| `DAGState`/`_shared_memory` in RAM | process memory | continuously | No | n/a |

**The seams (interfaces where state crosses a boundary):**
1. Harness↔subprocess: `.agentfield_output.json` file (write by agent, read by harness, deleted in `finally`); `asyncio.wait_for` wall-clock cap.
2. Coder↔coding-loop: `CoderResult` dict return (fallback on harness error) + iteration checkpoint file.
3. Coding-loop↔DAG: `IssueResult` return value, folded into `DAGState`, flushed at level boundaries.
4. DAG↔build: `DAGState.model_dump()` returned through `execute()`; `resume` parameter is the dormant seam (exists on `execute`/`run_dag`, not threaded from `build`).
5. Reasoner↔control-plane: `notify_call_complete`/`notify_call_error` → `POST /api/v1/workflow/executions/events`; status + payload rows.
6. Control-plane↔upstream client: SSE/HTTP request context; `task.cancel()` via the cancel route; `Request.Context().Done()` on disconnect.

## Historical Context (from thoughts/)

No prior `thoughts/` documents specifically cover SWE-AF run-state persistence or the cancellation cascade; the `thoughts/searchable/shared/research/` directory contains unrelated topics (`baml-structured-output-integration-surface.md`, `agentfield-ai-backend-claim.md`). This is the first research artifact on this topic.

A related operational record lives in the session memory file `agentplane-conformance-baseline.md` (outside `thoughts/`), which documents the observed build run `run_20260608_202801_560683uw`: build merged 8 issues to `feature/9d2ec5b4-implement-agentplane-endpoints`, then failed on `"Execution cancelled by upstream client"` after the issue-12 router-wiring coder turn hit the per-call cap. The merged controllers persisted on the feature branch (consistent with Layer 1 findings here).

## Related Research

- `thoughts/searchable/shared/research/2026-06-07-08-32-baml-structured-output-integration-surface.md` — touches the same harness structured-output surface (`.agentfield_output.json`, schema validation) from the BAML-integration angle.

## Resolved Questions (verified against live DB, 2026-06-09)

- **Which mechanism cancelled the build — RESOLVED.** The control-plane execution records (`GET /api/v1/executions/{id}` on `:8090`) show both `build` (`exec_20260608_202801_lo8k200j`, `duration_ms=21600428`, `error="Reasoner 'build' timed out after 21600.0s of active time"`) and `execute` (`exec_20260608_204927_eyc4ugq1`, `error="Reasoner 'execute' timed out after 21600.0s of active time"`) hit the **21600s active-time budget**. The cause was the SDK active-time watchdog (`agent.py:2351-2382`), **not** the control-plane cancel route (`cancel.py:78`) and **not** a client/connection disconnect. The `"Execution cancelled by upstream client"` log string is the generic inner side-effect of `task.cancel()` (`agent.py:2247`), overwritten in the record by the authoritative timeout `error_message`.

## Open Questions

- Whether the `"Expecting value: line 1 column 1"` message observed in logs came from a control-plane response parse or an SWE-AF caller (not the harness, per the verified guard at `_schema.py:148-149`) — would need the exact log line's source field to pin down. (Note: the `/api/v1/executions/{id}/events` endpoint is an SSE stream, so an empty/streamed read of it reproduces exactly this `json.loads("")` message — one concrete non-harness origin for the string.)

## Follow-up: Why no replan fired during the 6 hours (2026-06-09)

**Question:** during the 6-hour run, why was the work never sent to `run_replanner` — was it a big plan with successful builds, did issues fail too slowly to hit retries, or something else?

**Answer: it was a big plan with slow-but-successful level-by-level builds; no issue ever failed unrecoverably, so the replan gate was never eligible — and the 6-hour active-time budget expired mid-level-4.** Not a retry trap.

### Mechanism (source)

The replan gate is evaluated **exclusively at level boundaries** (`dag_executor.py:1678-1753`), after `_execute_level`'s `asyncio.gather` over the whole level returns — never mid-level or mid-issue. It fires only when the level's `unrecoverable` set is non-empty:

```python
unrecoverable = [f for f in level_result.failed
                 if f.outcome in (IssueOutcome.FAILED_UNRECOVERABLE, IssueOutcome.FAILED_ESCALATED)]
```
(`dag_executor.py:1679-1682`), AND `config.enable_replanning` AND `replan_count < max_replans` (default `2`, `schemas.py:998`). `COMPLETED_WITH_DEBT` goes into `level_result.completed` (`dag_executor.py:1192`) and is invisible to the gate. `max_coding_iterations=5` (`schemas.py:1002`) bounds the inner coder↔reviewer loop; each agent call also has a 45-min `agent_timeout_seconds=2700` guard (`schemas.py:1005`, `coding_loop.py:616`).

### Empirical run data (worker `swe-af-anthropic-swe-agent-1` logs, full 6h retained 20:28:01Z→02:54:35Z)

Reasoner lifecycle tally for `run_20260608_202801_560683uw`:

| Reasoner | started | completed | failed |
|---|---|---|---|
| run_product_manager / architect / tech_lead / sprint_planner | 1 each | 1 each | 0 |
| run_issue_writer | 13 | 13 | 0 |
| run_coder | 15 | 15 | **0** |
| run_code_reviewer | 15 | 15 | **0** |
| run_qa / run_qa_synthesizer | 9 / 9 | 9 / 9 | 0 |
| run_merger | 6 | 6 | 0 |
| run_integration_tester | 5 | 5 | 0 |
| run_workspace_setup / cleanup | 4 / 3 | 4 / 3 | 0 |
| **run_replanner** | **0** | **0** | **0** |
| execute / build | 1 / 1 | 0 / 0 | **1 / 1** (21600s timeout) |

**`run_replanner` was invoked zero times.** Every coder and reviewer call completed; no `run_coder failed`, no `FAILED_UNRECOVERABLE`. 15 coder calls for 13 issues = only 2 issues took a 2nd coding iteration — nowhere near the 5-iteration cap. So no issue ever reached an outcome the replan gate reads.

Milestone timeline (one `workspace_setup` per level = 4 levels planned):

| Time (UTC) | Event |
|---|---|
| 20:28:01 | build started |
| 20:45:50 | sprint_planner completed (13 issues, 4 levels) |
| 20:49:27 | level-1 workspace_setup |
| 21:40:51 / 21:43:53 | level-1 mergers completed |
| 21:55:27 | level-1 integration_tester completed |
| 21:56:05 | level-2 workspace_setup |
| 22:51:15 | level-2 merger; 23:18:30 / 23:34:42 integration_testers |
| 23:36:19 | level-3 workspace_setup |
| 00:28:49 | level-3 merger; 01:04:46 / 01:07:09 integration_testers |
| 01:07:49 | level-4 workspace_setup |
| **02:28:02** | **build failed** (20:28:01 + 21600s active-time budget) |
| 02:49:29 | level-4 merger completed (orphaned `execute` still running) |
| **02:49:30** | **execute failed** (20:49:27 + 21600s) |

Levels 1, 2, 3 each completed cleanly (setup → code → review → QA → merge → integration-test), ~60-100 min apiece. The build was **mid-level-4** when `build`'s watchdog fired at 02:28:02. (`execute`, a separate execution with its own watchdog, kept running ~21 min longer and even completed the level-4 merge at 02:49:29 — one second before its own 02:49:30 timeout.)

### The structural nuance ("something else")

Even setting aside that nothing failed: the slow issue-12 (router-wiring) sat in **level 4**, and the replan gate is only evaluated when a level's `asyncio.gather` returns. Level 4 never reached its boundary before the 6-hour wall, so a level-4 issue **could not** have reached the replan gate within the time regardless of its outcome. Combined with the fact that issue-12's coder turns were *completing* (slowly, ~28 min each — under the 45-min `agent_timeout_seconds`), nothing routed it toward `FAILED_UNRECOVERABLE` either. Replan was both **unwarranted** (no unrecoverable failure) and **structurally unreachable in time** (level 4 incomplete). The limiting resource was wall-clock, not the retry/replan budget.

### Verdict on the three hypotheses

- "Big plan with successful builds" — **yes**, this is what happened: 4 levels / 13 issues, all succeeding, ~90 min/level, out of time mid-level-4.
- "Runs failed but too slow to hit the 3× retry" — **no**: nothing failed; coding iterations (15 for 13 issues) were far under the 5-iteration cap; no retry/replan path was ever entered.
- "Something else" — the structural fact that replan is gated at level boundaries and only on unrecoverable outcomes, so a slow-but-not-failing issue in an unfinished level is invisible to it.
