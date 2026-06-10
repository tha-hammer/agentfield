---
title: Enhance BAML TDD plan — fix all 43 review findings
task: Enhance BAML TDD plan, fix all review findings
slug: baml-plan-enhance-fix-review-findings
effort: E3
phase: complete
progress: 39/39
mode: ALGORITHM
started: 2026-06-07
updated: 2026-06-07
---

## Context

The PlanReview pipeline produced a 28 KB report on
`thoughts/searchable/shared/plans/2026-06-07-09-04-tdd-baml-structured-output.md`
with verdict **Needs Major Revision**: 11 critical + 32 warning findings across
Contracts, Interfaces, Promises, Data Models, APIs, plus 3 self-referential
test-oracle violations in the Scaffolding & Realness hostile check.

Goal: edit the plan markdown in place so every finding is resolved. The dominant
theme is **undefined/unnamed interfaces** (BAML `@@dynamic` function schema,
`pydantic_to_typebuilder` signature, parser-only API, Go replacement API) and
**deferred decisions** (`.gitignore` commit-vs-generate, large-schema preserve/drop,
`response.JSON()` disposition). Fixes MUST be grounded in (a) the real codebase
signatures and (b) the real BAML Python/Go API — inventing fake interfaces would
reproduce the exact bug the review caught.

### Plan / approach
1. Extract ground-truth signatures from the 8 referenced code files (parallel agents).
2. Verify BAML's real Python + Go API surface from docs (parallel web research) —
   baml-py is NOT installed locally, so no on-disk introspection is possible.
3. Rewrite plan sections to specify every missing interface/schema/decision and
   replace the 3 self-referential oracles with independent ground-truth oracles.
4. Verify by re-reading edited sections against each finding.

### Risks
- (per premortem) Hallucinating BAML API names — mitigated by mandatory doc
  verification before writing any BAML symbol into the plan.
- Inventing codebase signatures — mitigated by reading actual files first.
- Over-prescribing above L1 tier — keep fixes correctness/idempotency-scoped.

## Criteria

### Critical findings (11)
- [x] ISC-1: ExtractDynamic @@dynamic function signature + TypeBuilder registration specified
- [x] ISC-2: Pydantic→TypeBuilder mapping table added (str/Optional/list/Enum/nested/union)
- [x] ISC-3: baml_bridge.pydantic_to_typebuilder full signature + raises stated
- [x] ISC-4: try_parse_from_text replacement name+signature+return type stated
- [x] ISC-5: Go WithSchema replacement interface decided and named
- [x] ISC-6: caller timeout forwarding through BAML call specified
- [x] ISC-7: response.JSON()/Into() disposition resolved to one branch
- [x] ISC-8: .gitignore commit-vs-generate decision resolved
- [x] ISC-9: BAML parser-only API endpoint named for Python and Go
- [x] ISC-10: TypeBuilder→Pydantic instance cast (deserialize) mechanics specified
- [x] ISC-11: Go SDK consumer-facing API after B-D1 fully specified

### Self-referential oracle fixes (3)
- [x] ISC-12: test_baml_bridge.py oracle uses independent hardcoded expected value
- [x] ISC-13: B-B2 Red unit step rewritten to deserialize(hardcoded_dict)==M(...)
- [x] ISC-14: B-A1 property oracle uses independent hardcoded function-name list

### Warning findings (resolved or explicitly bound)
- [x] ISC-15: BAML parse-failure exception class named + mapped to ValueError/FailureType.SCHEMA
- [x] ISC-16: rate_limiter.py wrapping contract with BAML async call specified
- [x] ISC-17: large-schema-on-disk (schema.go:50) preserve-or-drop decided as constant
- [x] ISC-18: _runner.py:332 catch-clause update note added for BAML exception
- [x] ISC-19: sync-vs-async BAML Python client choice specified
- [x] ISC-20: agent_ai.py concrete BAML call expression shown
- [x] ISC-21: Go BAML client concrete call expression shown
- [x] ISC-22: baml_bridge.py module path + visibility stated
- [x] ISC-23: B-C1 retry behavior after BAML parse failure specified
- [x] ISC-24: baml-cli generate idempotency stated (no nondeterministic output)
- [x] ISC-25: BAML Go runtime module path + go 1.21 compatibility named
- [x] ISC-26: Go BAML client goroutine-safety addressed
- [x] ISC-27: recursive Pydantic models + Union>2 variants behavior specified
- [x] ISC-28: pyproject.toml placement of baml-py decided
- [x] ISC-29: FailureType enum full variant list shown
- [x] ISC-30: Person model full field list specified
- [x] ISC-31: messy_cli_output.txt fixture embedded JSON schema specified
- [x] ISC-32: Go baml_src/ function + model definitions specified
- [x] ISC-33: clients.baml provider/model config schema specified
- [x] ISC-34: baml-py version pin specified

### Grounding + meta
- [x] ISC-35: every codebase signature written verified against the actual file
- [x] ISC-36: every BAML API symbol written verified against BAML docs
- [x] ISC-37: edited plan re-reads clean against all 11 critical findings
- [x] ISC-A1: NO fix introduces an unverified/invented API name (anti-criterion)
- [x] ISC-A2: settled decisions preserved — replace-not-coexist, Web-UI out of scope (anti-criterion)

## Decisions

- PRD placed under `thoughts/.../work/` (untracked) rather than `MEMORY/WORK/`
  to avoid polluting the agentfield OSS repo tree.
- Silmari: recall empty, no in-progress cards. Memory integration available but no priors.
- baml-py not installed → BAML API claims must come from docs, verified per ISC-36.

## Verification

All 39 criteria resolved by editing the plan (363 → 725 lines). Evidence:

- **Grounding (ISC-35/36):** 5 parallel agents extracted real signatures (3
  codebase Explore + 2 BAML-API web). Confirmed against installed `baml-py==0.222.0`
  in `.venv`: `baml-cli` present, `errors` classes incl. `BamlValidationError`/
  `BamlTimeoutError`/`BamlClientHttpError`, `ClientRegistry` present, `TypeBuilder`
  NOT in `baml_py` (generated). Web research corrected: per-call timeout
  unsupported (#1630), BAML Go pre-stable/Linux-mac-only.
- **No invented symbols (ISC-A1):** `grep -niE "TypeBuilder.parse_text|parse_raw|
  options=\{timeout|b.request"` → `NONE-FORBIDDEN`.
- **No deferrals (ISC-7/8/17 etc.):** consolidated into "Resolved Decisions" §1–5;
  only the header references the review.
- **Self-ref oracles fixed (ISC-12/13/14):** B-A1 uses hardcoded `EXPECTED_FUNCTIONS`;
  B-B2 unit uses `deserialize(b.parse.ExtractDynamic(json.dumps(HARDCODED),...), M)
  == M(**HARDCODED)` (hand-authored RHS, `b.parse` = no LLM).
- **Decisions preserved (ISC-A2):** replace-not-coexist + Web-UI-out-of-scope intact.
- **Capability audit:** Mass Parallelism / Background Agents invoked via 5 `Agent`
  tool calls (3 Explore + 2 web-search-researcher) — all returned. No phantom
  selections.

### Memory cards

#### Q-learning
Save block, `mode: root`, `source: algorithm-baml-plan-enhance-fix-review-findings-learn`,
`kind: signal`, `trunk: 5`:
> When fixing "undefined interface" plan-review findings, the fix is NOT to invent
> a plausible signature — it is to verify the real API (installed pkg + docs) and
> the real codebase file:line, then write the verified symbol. Inventing reproduces
> the exact bug the review caught. baml-py example: `TypeBuilder` is in the
> generated `baml_client`, not `baml_py`; per-call timeout does not exist (#1630).

No prior cards matched (recall empty) → novel, root. Resume hook: not needed
(39/39 complete, no open criteria).
