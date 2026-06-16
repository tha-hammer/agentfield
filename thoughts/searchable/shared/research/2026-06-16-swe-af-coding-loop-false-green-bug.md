---
date: 2026-06-16
author: tha-hammer
repository: SWE-AF-anthropic (bug); observed via agentfield/AgentPlane port
topic: "SWE-AF coding loop reports false-green when agents crash in the sandbox"
tags: [bug, swe-af, coding-loop, false-green, sandbox, bwrap]
severity: high
status: open
---

# Bug: SWE-AF coding loop rubber-stamps agent crashes as approved → false green

## Summary

When the code-reviewer (and QA-synthesizer) agent **throws any exception** — including
an infrastructure crash where the agent never actually ran — the coding loop records the
iteration as `approved=True, blocking=False`. A run in which **every** coder/QA/reviewer
agent crashed in the sandbox is therefore reported as **completed / approved**: a false
green that masks total failure as success.

Observed on run `run-0055a5c7-a52a-4f3a-a977-05e71e92418b` (target
`/workspaces/cosmic-HR-ed8f73b9`, Meetings TDD plan B1–B16), driven through the Elixir
AgentPlane port. The run never exercised the port's code-writing path — it died one layer
below, in the SWE-AF coding harness's sandbox (`swe-af-anthropic-swe-agent` image).

## Root cause (code locus)

`swe_af/reasoners/execution_agents.py:1185-1198` — the code-reviewer reasoner's exception
handler converts **any** failure into an approval:

```python
except FatalHarnessError:
    raise  # Non-retryable — propagate immediately
except Exception as e:
    router.note(f"Code reviewer agent failed: {issue_name}: {e}", tags=["code_reviewer", "error"])

return CodeReviewResult(
    approved=True,            # <-- BUG: reviewer crash treated as approval
    summary=f"Code reviewer agent failed for {issue_name} — not blocking",
    blocking=False,
    iteration_id=iteration_id,
).model_dump()
```

A reviewer that **never ran** (sandbox crash) is not the same as a reviewer that
**chose not to block**. The bwrap sandbox error (`bwrap: No permissions to create a new
namespace`) is caught by the generic `except Exception` rather than being classified as a
harness/infra failure, so it falls through to `approved=True`.

Compounding paths:
- `swe_af/reasoners/execution_agents.py:1262-1289` — QA-synthesizer fallback decides
  `approve` from `tests_passed and review_approved and not review_blocking`; `review_approved`
  is now `True` from the crash above.
- `swe_af/execution/coding_loop.py:496-498` and `:351` — iteration accept gate is
  `qa_passed and review_approved and not review_blocking` / `approved and not blocking`.
- `swe_af/execution/coding_loop.py:718` — iteration record stores
  `"review_approved": review_result.get("approved", False)`.

## Evidence (run-0055a5c7, queried from AgentPlane DB `agentplane_dev.executions`)

- **16** executions' `result_payload` contain `bwrap … No permissions to create a new namespace`.
- **11** executions carry the literal `"not blocking"`.
- Terminal `build` execution: DB `status = succeeded` while its own payload says
  `success: false`, summary: *"Partial: 4/16 issues completed … every local shell command
  failed before execution with a bwrap namespace error … none of the acceptance criteria
  can be marked passed with evidence."*
- 9 coding iterations each: `files_changed: []`, `qa_passed: null`, `review_approved: true`,
  summary `"Code reviewer agent failed … not blocking"`.
- On disk (`/workspaces/cosmic-HR-ed8f73b9`): `git status` clean, no migration 038, no
  `meetings/` UI, no route/helper changes. AC3–AC18 unmet.

## Why it's reportable independent of the sandbox

Same host, identical Docker privileges across containers (`Privileged=false`, no `CapAdd`,
no `SecurityOpt`). The `swe-af-baml:smoke` image's sandbox **works** (no bwrap error; 15 real
commits of B1–B12 on `feature/e5f93eb7-implement-meetings-tdd`). So the sandbox break is
image-specific — but the **false-green logic** (crash → approved) is a defect in the coding
loop that would mask *any* class of agent crash, not just bwrap.

## Proposed fix (surgical)

1. **Distinguish reviewer infra-failure from reviewer approval.** In the
   `execution_agents.py:1193` fallback, return `approved=False, blocking=True` (or add a
   distinct `errored=True` status) on exception. A reviewer that never ran must not count
   as an approval.
2. **Classify harness/sandbox errors as fatal.** Errors like `bwrap: No permissions to
   create a new namespace` and "every local shell command failed before execution" should
   raise `FatalHarnessError` (already special-cased at `:1185`) rather than fall through the
   generic `except Exception`.
3. **Guard the iteration-accept gate.** Never auto-approve an iteration whose coder produced
   `files_changed == []` together with tool/shell failures. Empty diff + agent error ⇒ fail
   the iteration.
4. **Propagate to run terminal status.** N consecutive infra-errored, zero-file-change
   iterations ⇒ mark the run `failed`, not `succeeded`/`completed`.

## Separate follow-up (not this bug)

Fix the sandbox in the `swe-af-anthropic-swe-agent` image so bwrap can create a namespace
(likely `--security-opt seccomp=unconfined` and/or userns on a **single-container**
`docker run` — never `docker compose` on this shared host), or adopt the no-sandbox/bypass
mode the BAML image already uses.
