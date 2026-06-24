---
date: 2026-06-16T11:14:06-04:00
researcher: tha-hammer
git_commit: 1471c373
branch: fix/harness-retry-preserve-goal
repository: agentfield
topic: "AgentPlane Elixir port — SWE-AF/Codex e2e + UI-API data wiring"
tags: [implementation, agentplane, swe-af, codex, ui-api, drop-in, cosmic-hr]
status: complete
last_updated: 2026-06-16
last_updated_by: tha-hammer
type: implementation_strategy
---

# Handoff: AgentPlane port — Codex e2e build + UI-API data wiring (IN PROGRESS)

## Task(s)

Overarching goal: prove the **Elixir AgentPlane port** (a reimplementation of the Go AgentField control plane) can run real **SWE-AF** builds, then drive an actual build (the cosmic-HR "Meetings surface" TDD plan) through it using the **Codex** runtime, and make the **copied React UI** fully functional against the port.

1. **AgentPlane runtime-contract fixes — DONE & validated.** status_update handler, async await timeout, UI bind, set_terminal cancel guard. Suite green.
2. **Codex e2e build of the cosmic-HR meetings plan — RUNNING.** Repointed an SWE-AF planner at AgentPlane with Codex runtime; after a chain of fixes the build now does real work (`run-0055a5c7`): git_init→PM→architect→tech_lead review loop. Architect is executing on the plan.
3. **UI-API data wiring — IN PROGRESS (the active problem).** The "drop-in" work wired correct JSON **shapes** but **static empty stubs** — ~50 endpoints return hardcoded empty data, so the UI loads but is non-functional. **Fixed so far:** the run-list (`executions/enhanced`, `executions/recent`, `workflow-runs index`). **STILL BROKEN (user's live complaints):** nodes/agents show "0/0 online / no agent nodes found", top nav "disconnected", execution detail + logs empty, dashboard counts all 0.

## Critical References

- **The plan being built** (READ-ONLY, committed on `cosmic-HR04`, remote `github.com/Cosmic-HQ/cosmic-HR`): `/home/maceo/Dev/cosmic-HR04/thoughts/searchable/shared/plans/2026-06-16-09-18-tdd-meetings-surface-specs-19-20.md` (16-behavior TDD plan, Meetings on existing `meetings` table mig 017).
- **Go spec (READ-ONLY ground truth)** for UI-API shapes: `/home/maceo/Dev/agentfield/control-plane/internal/handlers/ui/` — `executions.go:202` (EnhancedExecution), `recent_activity.go:40` (ActivityExecution), `workflow_runs.go:27` (WorkflowRunSummary), `dashboard*.go`, `nodes_*` handlers. Also visible in-container at `/workspaces/agentfield-src/`.
- **UI client (READ-ONLY, what the UI actually calls)**: `/home/maceo/Dev/agentfield/control-plane/web/client/src/` — `components/HealthStrip.tsx` (top nav), `hooks/queries/useAgents.ts` → `getNodesSummary` → `GET /api/ui/v1/nodes/summary`; `hooks/useSSEQuerySync.tsx:66` → SSE `GET /api/ui/v1/executions/events` (`execConnected` set by EventSource `onopen`); `types/agentfield.ts:16` (`AgentNodeSummary`).

## Recent changes

All AgentPlane edits are on branch **`feature/1aafa39b-build-ui-api-surface-for-react-ui-copy`** in the git repo at host `/home/maceo/Dev/silmariAgentPlane/agentplane` (= container `/workspaces/project/agentplane`, a **bind mount** — edits persist regardless of container state). Files are root-owned in-container; patch via `docker exec -i swe-af-anthropic-swe-agent-1 python3 - <<'PY'` (host Edit gives EACCES).

- `apps/agentplane_web/lib/agentplane_web/controllers/api_v1/executions_controller.ex` (or core `ExecutionsController`): `status_update/2` was an always-500 stub (called arity-1 `update_status`); now calls real `update_status(id, params)`, returns 404 on not_found (was 500).
- `apps/agentplane_core/lib/agentplane_core/executions.ex`: added `set_terminal` cancel guard (won't overwrite `cancelled` with `succeeded`); added **`list_recent/1`, `count_all/0`, `list_run_summaries/1`** + `run_summary/2` + `_dt_iso` helpers (UI list queries).
- `apps/agentplane_core/lib/agentplane_core/dispatch/types.ex:63` `default_callback_timeout_ms`: now reads `AGENTFIELD_ASYNC_DEFAULT_EXECUTION_TIMEOUT` (secs→ms, fallback 90s) — was hardcoded 90_000.
- `apps/agentplane_web/lib/agentplane_web/controllers/ui/executions_controller.ex`: `enhanced` + `recent` wired to `Executions.list_recent` with Go shapes (`EnhancedExecution`/`ActivityExecution`) + helpers (`paginate`, `enhanced_row`, `activity_row`, `relative_time`, `duration_display`). **`summary`, `stats`, `timeline`, `filter_options`, `view_stats` STILL STUBS.**
- `apps/agentplane_web/lib/agentplane_web/controllers/ui/workflow_runs_controller.ex`: `index` wired to `Executions.list_run_summaries` (WorkflowRunSummary shape). `show` already queried real data.
- Tests updated (relaxed empty-stub asserts → type checks): `executions_controller_test.exs`, `u2_dashboard_executions_test.exs`. Also earlier: `nodes.ex` mode→`auto`, cancel-message tests, `stub_routes_test.exs` capture_io, `dispatch_test.exs` cancel race, `types_test.exs`.
- **SWE-AF Codex harness fix** — `swe_af/runtime/codex_harness_patch.py` (BOTH host `/home/maceo/Dev/SWE-AF-anthropic/...` AND container `/app/swe_af/...` — container `/app` is the **image copy, NOT bind-mounted**, so the container copy is what runs): added `_ensure_typed/1` and applied it in the `anyOf/oneOf/allOf` loop of `_codex_strict_json_schema`. Gives typeless composed-schema branches a `type` so OpenAI `--output-schema` accepts them.

## Learnings

- **The drop-in stubbed shapes, not data.** Audit command: `grep -rn "GoShapes\." apps/agentplane_web/lib/agentplane_web/controllers/ui/*.ex`. Every action piping a static `GoShapes.*` builder is an empty stub. ENTIRE nodes/agents/dashboard/DID/identity surface + several executions actions are stubs. `GoShapes` builders live in `apps/agentplane_web/lib/agentplane_web/ui/go_shapes.ex`.
- **Node IS registered & heartbeating** but every read endpoint returns empty because the handlers are stubs. DB: `swe-planner | http://swe-agent:8004 | last_seen ~now`. The Node schema (`apps/agentplane_core/lib/agentplane_core/nodes/node.ex`) has `state, health_score, lifecycle_status, health_status, last_seen, reasoners, skills`. The `Nodes` context (`nodes.ex`) has NO `list/0` yet — `register_go_map` hardcodes status fields. **UI `onlineCount` (HealthStrip.tsx:50) counts a node online iff `health_status ∈ {ready,active}` OR `lifecycle_status ∈ {running,ready}`** — so the nodes/summary handler must compute online from `last_seen` freshness and set those fields.
- **SSE actually opens.** `curl -sN http://localhost:4000/api/ui/v1/executions/events` returns `event: snapshot\ndata: {...}` (rc=0). So `execConnected` should be true; "disconnected" may be browser-only (the EventSource appends `?api_key=` — check auth) OR a stream that drops. SSE controller: `apps/agentplane_web/lib/agentplane_web/controllers/ui/sse_controller.ex`. Verify in the browser before assuming the SSE is broken.
- **Codex root-causes (all fixed, in order discovered):** (1) BuildConfig is `extra=forbid` — `branch` belongs in `repos:[{repo_url,role,branch}]`, not top-level; (2) runtime=codex but model resolved to deepseek — env `SWE_DEFAULT_MODEL` overrides runtime preset; explicit `config.models.default` wins (schemas.py:636 `resolve_runtime_models` precedence); (3) `gpt-5.3-codex` is NOT valid on a ChatGPT-auth Codex account — use **`gpt-5.5`** (config default), `gpt-5.4`, `gpt-5.4-mini` (all validated via `codex exec -m <model>`); (4) **the real blocker** — codex `--output-schema` rejects reasoner schemas because Pydantic emits `{"anyOf":[{},{"type":"null"}]}` for `Optional[Any]` (`AskUserFormField.default_value`); OpenAI structured output requires every branch to have a `type`. Fixed by `_ensure_typed`.
- **`docker exec -d` + `docker exec -i`** required for detach / heredoc-stdin respectively. Host background processes (socat via `&`/run_in_background) get reaped (exit 143/144) — use a **docker sidecar** for anything that must persist.

## Artifacts

- This handoff.
- AgentPlane fixes on branch `feature/1aafa39b-build-ui-api-surface-for-react-ui-copy` (see Recent changes for files).
- Codex build payload (ready to re-fire): container `/tmp/codex_build_payload.json` — goal pins architect to the plan; `config: {runtime:"codex", repos:[{repo_url, role:"primary", branch:"cosmic-HR04"}], enable_github_pr:false, check_ci:false, models:{default:"gpt-5.5", coder:"gpt-5.4-mini", qa:"gpt-5.4", code_reviewer:"gpt-5.4", qa_synthesizer:"gpt-5.4"}}`.
- Memory: `/home/maceo/.claude/projects/-home-maceo-Dev-agentfield/memory/agentplane-swe-af-integration.md` (repoint recipe + fixes).

## Action Items & Next Steps

**Resume the UI-API wiring — this is what the user is angry about. Make the whole UI functional, not one endpoint at a time.** Wire each stub to real data in the exact Go shape. Restart AgentPlane after edits, verify via `localhost:4000`.

1. **Nodes/agents online (most visible — "0/0", "no nodes"):**
   - Add `Nodes.list/0` (`Repo.all(Node)`) + an `online?` helper (`last_seen` within ~60s).
   - Wire `nodes_controller.ex` `summary` → `{nodes: [AgentNodeSummary...], count}` with `health_status`/`lifecycle_status` reflecting online; `status`, `details`.
   - Wire `agents_controller.ex` `running` → `{running_agents, total_count}`.
   - Wire `dashboard_controller.ex` `summary`/`enhanced` counts (nodes online/total, executions today/yesterday, agents running) — `dashboard*.go` shapes.
   - Map fields per `AgentNodeSummary` (types/agentfield.ts:16): `id`(=node_id), `base_url`, `version`, `team_id`, `health_status`, `lifecycle_status`, `last_heartbeat`(=last_seen iso), `reasoner_count`(=length reasoners), `skill_count`.
2. **"Disconnected":** verify the executions/events SSE in the browser (curl shows it opens). If broken in-browser, check the `?api_key=` auth append and that the stream sends periodic keepalives in `stream_loop` (sse_controller.ex). Also `nodes/events` SSE (nodeConnected).
3. **Execution + logs empty:** wire `executions_controller.ex` `details` fully; verify `workflow_runs_controller.ex` `show` returns the execution tree the `RunTrace.tsx` component expects (`node.execution_id`, `node.reasoner_id`, `node.parent_execution_id`, `node.started_at`, children). Logs: AgentPlane likely does NOT capture reasoner logs (SDK POSTs to `/executions/:id/logs` which 404s) — decide whether to add that route + storage or return empty gracefully.
4. **Re-run the full suite** after wiring (`MIX_ENV=test mix test` in container) — expect to relax more empty-stub test asserts (drop-in gate tests). Baseline was core 182/0, web 496/0/1-skipped.
5. The codex build (`run-0055a5c7`) was progressing through architect/tech_lead. Check its final state; confirm the architect output actually follows the plan. Re-fire with `curl -s -X POST http://localhost:4000/api/v1/execute/async/swe-planner.build -d @/tmp/codex_build_payload.json` if needed.

## Other Notes

**Infra topology (shared Docker host — see DANGER below):**
- Container `swe-af-anthropic-swe-agent-1` (up; `docker start <name>` is the only safe single-container op). Its PID 1 = the **Go-pointed** swe-planner (`python -m swe_af`, NODE_ID=swe-planner, `AGENTFIELD_SERVER=control-plane:8080`, port 8003).
- **AgentPlane** runs in-container: `docker exec -d ... bash -lc 'cd /workspaces/project/agentplane && MIX_ENV=dev PORT=4000 mix phx.server > /tmp/ap_e2e3.log 2>&1'`, bound `0.0.0.0:4000`. Needs Postgres: `pg_ctlcluster 17 main start`, DB `agentplane_dev` (user postgres). Verify: `curl localhost:4000/ui/` (via sidecar) or `curl <container-ip>:4000/ui/`.
- **Codex planner** (the one driving the build): a SECOND `python -m swe_af` started with `docker exec -d -e GH_TOKEN="$(gh auth token)" ... 'cd /app && AGENTFIELD_SERVER=http://localhost:4000 SWE_DEFAULT_RUNTIME=codex NODE_ID=swe-planner python -m swe_af > /tmp/planner_codex.log 2>&1'`. Registers `swe-planner@swe-agent:8004` on AgentPlane (port auto-increments if 8004 busy → 8005). GH_TOKEN (host `gh auth token`, tha-hammer) is in its process env (not persisted) to clone the private cosmic-HR repo.
- **UI host access:** docker sidecar **`ap-ui-fwd`** (`docker run -d --name ap-ui-fwd --restart unless-stopped --network swe-af-anthropic_default -p 4000:4000 alpine/socat TCP-LISTEN:4000,fork,reuseaddr TCP:swe-af-anthropic-swe-agent-1:4000`) → `localhost:4000` reaches the UI. Host can also reach container IP directly (`172.19.0.4:4000`). Remove with `docker rm -f ap-ui-fwd`.
- **Codex CLI:** `codex-cli 0.137.0`, ChatGPT auth (`~/.codex/auth.json`), config model `gpt-5.5` xhigh. Verify usable: `codex exec --skip-git-repo-check "Reply OK"`.

**⚠️ DANGER (memory `docker-compose-shared-host-orphan-removal.md`):** NEVER `docker compose up`/`--force-recreate` on this host — it removed 4 other-projects' containers (~36h/$330 lost). Single-container `docker start`/`docker run --name` only. `/app` in the container is the SWE-AF **image copy** (not the host repo) — the running codex planner uses `/app/swe_af`, so the codex harness fix had to be applied there too.

**Test status (AgentPlane, MIX_ENV=test):** core 182/0, web 496/0/1-skipped after the run-list wiring (before the nodes/dashboard/execution wiring still to come).

**Beads:** not used this session (no `bd` issues created/touched).
