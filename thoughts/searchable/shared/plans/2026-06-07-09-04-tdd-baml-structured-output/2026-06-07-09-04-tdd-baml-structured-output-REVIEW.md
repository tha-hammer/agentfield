# Plan Review Report: BAML Structured-Output Integration — TDD Implementation Plan

> **Plan tier: L1** (confidence: high) — 3 L1 keyword(s) + plan size 363 lines (bounded scope; L1 >= L2 signal)
>
> Findings below are scoped to this tier's lens. L1 plans receive correctness-and-idempotency findings; L2 plans add typed-interface and schema-compat findings; L3 plans add UnifiedError/retry/observability. Above-tier prescriptions were forbidden in each phase prompt.

## Review Summary

| Category | Status | Issues Found |
|----------|--------|--------------|
| Contracts | ❌ | 7 issues |
| Interfaces | ❌ | 8 issues |
| Promises | ❌ | 13 issues |
| Data Models | ❌ | 9 issues |
| APIs | ❌ | 6 issues |

## Contracts Review

### Well-Defined
- ✅ Codegen drift gate oracle (B-F1) — CI step: modify .baml, run generate, `git diff --exit-code`; binary pass/fail with no code-produced artifact as the oracle
- ✅ API-key skip contract for integration tests — B-B1, B-B2, B-D1: skip when key env var absent, real call when present, assert independently-known facts (e.g. 'Ada Lovelace') — not a spy oracle
- ✅ FailureType.SCHEMA preservation on unrecoverable parse (B-C1) — B-C1 Edge Cases preserves classification on unrecoverable input; blast-radius note keeps retry orchestration in _runner.py separate from _schema.py parse internals
- ✅ Fixture-based deterministic test oracle (B-C1, B-D2) — Recorded messy-CLI fixture is a pre-committed artifact; `parse(fixture) == embedded` compares against the fixture's own independently authored embedded object
- ✅ Provider env-var config boundary — clients.baml reads env.ANTHROPIC_API_KEY / env.OPENAI_API_KEY / env.OPENROUTER_API_KEY — same vars already in the SDK; no new credential surface
- ✅ Vertical-slice implementation order with named observable artifacts — Seven milestones each state a concrete user-observable output; no horizontal layer-only step

### Missing or Unclear
- ⚠️ test_baml_bridge.py oracle is potentially self-referential — B-B2 TDD Cycle, Red step: 'maps a nested Pydantic model → TypeBuilder and asserts the built type shape' with no independent expected value named
- ⚠️ BAML exception class → existing exception type mapping unspecified at the _runner.py boundary — B-C1 Edge Cases; B-B2 Edge Cases ('typed BAML validation error surfaced') never names the exception class BAML raises
- ⚠️ rate_limiter.py wrapping contract with the BAML call is unspecified — B-B2 Refactor step: 'ensure rate-limit retry still wraps the call' without identifying the 429 exception class
- ⚠️ B-A1 property test oracle risks self-reference — B-A1 Property specification: enumerating F from baml_src/ derives the oracle from the same artifact under generation
- ⚠️ Large-schema-on-disk behavior (schema.go:15,50) explicitly deferred without resolution — B-D2 Edge Cases / Risks & Open Items; neither preserve nor drop branch is chosen

### Critical
- ❌ BAML parser-only API endpoint is never named for Python or Go — B-C1/B-D2 Green steps; BAML Replacement Mechanics. Without a named API, the Red step is unwritable; the approach may be structurally incompatible with the fixture-based @@dynamic path or not exist in the installed baml-py version
- ❌ TypeBuilder → Pydantic instance cast mechanics unspecified — B-B2 Test Specification / BAML Replacement Mechanics. The boundary between BAML dynamic output and the Pydantic return type is undefined; callers may silently receive a raw dict and no test would catch it

### Recommendations
- B-C1 and B-D2 — BAML parser-only API: Add a spike step before B-C1's Red to confirm the exact Python (`b.parse.FunctionName()` / `TypeBuilder.parse_text()`) and Go function that parses arbitrary text without an LLM call; if no generic entry exists, pre-generate a single `ParseRaw(text) -> DynamicOutput` function
- B-B2 — TypeBuilder dynamic output → Pydantic instance cast: In baml_bridge.py specify a `(builder, deserialize)` pair where `deserialize(raw_dict) -> T` calls `schema.model_validate` and lets ValidationError propagate; agent_ai.py:792 already catches it
- B-B2 — test_baml_bridge.py oracle: Replace 'asserts the built type shape' with a concrete independent oracle — `deserialize(hardcoded_dict) == M(field=value, ...)` authored before implementation
- B-B2 — rate_limiter.py boundary: Confirm baml-py's HTTP 429 exception class; translate to the existing rate-limit type at the agent_ai.py call site
- B-C1 and B-B2 — BAML exception → existing error type: At each replacement call site, catch baml-py's parse-failure exception and translate to ValueError / FailureType.SCHEMA so _runner.py's existing logic receives types it handles
- B-D2 — large-schema-on-disk decision: Decide before implementation; document the chosen branch in B-D2's spec, not as an open item

## Interfaces Review

### Well-Defined
- ✅ baml-cli generate output check (B-A1) — `from baml_client import b` + `hasattr(b, F)` for every F declared in baml_src/
- ✅ ExtractPerson static BAML function signature — B-B1: `ExtractPerson(bio: string) -> Person`; result.name == 'Ada Lovelace' is an independent factual oracle
- ✅ try_parse_from_text existing signature (to be replaced) — _schema.py:209 `(text: str, schema: Any) -> Optional[Any]`
- ✅ WithSchema existing Go public signature (to be removed/replaced) — request.go:258; JSON() at response.go:99
- ✅ CI drift gate observable behavior — B-F1: `baml-cli generate` then `git diff --exit-code`; Makefile `baml-generate` target
- ✅ agent_ai.py parse ladder to be deleted — agent_ai.py:788-824 fully shown; deletion confirmed at B-B2 refactor step

### Missing or Unclear
- ⚠️ Go BAML-generated client call expression never shown — B-D1, sdk/go/ai/baml_structured_integration_test.go (new); plan says 'route callers through it' without naming the new exported identifier or its signature
- ⚠️ agent_ai.py BAML call expression undefined — B-B2; plan does not state the Python call expression or how the TypeBuilder is passed
- ⚠️ Exception-type compatibility between BAML parse errors and _runner.py catch logic not addressed — B-C1 refactor; _runner.py:332/:496; if BAML raises baml_py.BamlError instead of ValueError, _handle_schema_with_retry silently stops catching
- ⚠️ Large-schema-on-disk decision for B-D2 explicitly deferred — schema.go:15,50; public contract of the Go harness is undefined for callers
- ⚠️ baml_bridge.py module-level visibility not stated — B-B2; no __all__, module path (agentfield.baml_bridge vs. agentfield.harness.baml_bridge) not given

### Critical
- ❌ `baml_bridge.pydantic_to_typebuilder` parameter and return types are not stated anywhere — B-B2, TDD Red step 1. The 'load-bearing piece' has no signature; test_baml_bridge.py is unwritable and agent_ai.py has no callable target
- ❌ Replacement for `try_parse_from_text` in _schema.py has no stated name or signature — B-C1; _runner.py:332 calls by name, return type changes from `Optional[Any]` to an undefined 'typed parse error'; _runner.py breaks silently or fails mypy
- ❌ Go public functions `WithSchema` and `(*Response).JSON` removed with no replacement interface — B-D1; request.go:258, response.go:99. The 'or delegate' ambiguity leaves the exported interface undefined; any Go call site fails to compile with no migration path

### Recommendations
- baml_bridge.py (B-B2): Add `def pydantic_to_typebuilder(model: type[BaseModel]) -> baml_py.TypeBuilder`, raises TypeError for unsupported field types; confirm internal-only module path
- _schema.py replacement entry point (B-C1): Keep the name and update return type to `Any` raising a named BAML exception, or rename and update the one _runner.py call site; state the exception type
- Go WithSchema / JSON() fate (B-D1): Pick (a) keep both as exported wrappers delegating to BAML (lower-blast-radius L1 choice) or (b) delete and name the single new exported function
- Go BAML client call expression (B-D1): Add a one-line example showing the concrete Go expression so the integration test is unambiguously implementable
- Exception-type compatibility (B-C1, B-B2): Add a one-line refactor note to update _runner.py catch clause at :332 to also catch the BAML validation exception

## Promises Review

### Well-Defined
- ✅ Missing API key → integration test skipped, not failed (B-B1, B-D1) — edge cases sections
- ✅ Unparseable harness text → typed parse error, not silent None; FailureType.SCHEMA preserved — B-C1 edge cases
- ✅ Real-LLM tests use known-fact assertions as independent oracle; spy/mock never the passing oracle — B-B1, B-B2, B-D1; Testing Strategy
- ✅ Drift gate: stale committed client → loud CI failure; baml-cli absent → build error with install hint — B-F1 edge cases
- ✅ rate_limiter.py continues to wrap BAML calls after migration — B-B2 refactor step
- ✅ Old parse ladders deleted (not coexisted) once BAML path is green — What We're NOT Doing; B-B2/B-C1 refactor
- ✅ TypeBuilder mapper has a standalone deterministic unit test independent of any LLM call — B-B2 blast-radius note
- ✅ Vertical slice implementation order with visible user-observable artifact per milestone — Implementation Order section

### Missing or Unclear
- ⚠️ Retry behavior after BAML parse failure in B-C1 unspecified — B-C1; _runner.py:46,59,332; whether a parse failure re-invokes the provider CLI or no-op loops on identical fixed text is unstated
- ⚠️ Sync vs. async BAML Python client choice not specified for B-B2 — agent_ai.py:816; sync client blocks the event loop
- ⚠️ Python exception type emitted on BAML parse/validation failure not declared — B-B2; agent_ai.py:808 catches ValueError and breaks silently if BAML differs
- ⚠️ Idempotency of baml-cli generate not stated — B-F1/B-A1; if generated client embeds timestamps, git diff is always non-empty
- ⚠️ [VERIFY] BAML Go runtime module path and go 1.21 compatibility never named — B-D1 files touched; sdk/go/go.mod
- ⚠️ Go BAML client goroutine-safety not addressed — B-D1; sdk/go/ai may serve concurrent callers
- ⚠️ Large-schema token gate (schema.go:50) disposition deferred with no decision recorded — B-D2 edge cases
- ⚠️ B-B2 unit test oracle for test_baml_bridge.py structurally self-referential — B-B2 Red step; expected shape derived from the same Pydantic model the mapper processes
- ⚠️ Recursive Pydantic models and Union types with >2 variants through TypeBuilder unspecified — B-B2 edge cases; silent regression for callers using such models
- ⚠️ pyproject.toml placement of baml-py undecided, contingent on the unresolved commit-vs-generate decision — B-A1; pyproject.toml:31-32

### Critical
- ❌ Caller-specified timeout silently dropped after B-B2 migration — agent_ai.py:544 sets litellm_params['timeout']; after replacing litellm with a BAML client that parameter no longer reaches the HTTP call. Callers passing explicit timeout= silently get the wrong value; invisible in tests
- ❌ response.JSON() / Into() disposition in B-D1 is an unresolved OR — response.go:99; 'remove ... or delegate to BAML' leaves Go callers of resp.JSON(&dest) with undefined post-migration behavior
- ❌ .gitignore / commit-vs-generate decision deferred in B-A1 is a load-bearing unresolved promise — every downstream behavior imports from baml_client; CI is broken in a clean checkout until resolved, and dependency classification is undefined

### Recommendations
- B-B2 timeout forwarding: Pass effective_timeout as a per-call BAML client option (options= parameter); document the forwarding in baml_bridge.py
- B-D1 response.JSON disposition: Pick (a) remove and require callers to use the BAML-typed return, or (b) keep delegating to the BAML Go parser; enforce with a deletion comment in response.go:99
- B-A1 .gitignore decision: Make it explicit — committing the generated client avoids build-before-test ordering; if generate-on-build, add an autouse fixture / Makefile prerequisite
- B-C1 retry behavior: State whether parse failure re-invokes the provider CLI or halts immediately with FailureType.SCHEMA
- B-B2 async BAML client: Use baml-py's async client since agent_ai.py is already async; name the constraint explicitly
- B-B2 unit test oracle: Hard-code an expected TypeBuilder output-shape dict derived independently, not from the mapper's own output
- B-D2 large-schema gate: Make the preserve-or-drop decision explicit, recorded as a constant/comment not prose

## Data Models Review

### Well-Defined
- ✅ Existing Python inline-AI structured-output path — agent_ai.py:546-557 and :788-824 fully cited with types and control-flow
- ✅ Existing harness parse entry-point and retry constant — _schema.py:209, _runner.py:46,59,332
- ✅ Existing Go structured-output path — request.go:258,390, response.go:99
- ✅ Existing Go harness large-schema gate constants — schema.go:15,16,50
- ✅ Provider env-var names — ANTHROPIC_API_KEY / OPENAI_API_KEY / OPENROUTER_API_KEY for clients.baml
- ✅ SUPPORTED_PROVIDERS set — _factory.py:9: {claude-code, codex, gemini, opencode}
- ✅ Python test marker convention — pyproject.toml:77-81 markers unit/functional/integration; pytest-asyncio present
- ✅ Delete-not-coexist strategy — What We're NOT Doing

### Missing or Unclear
- ⚠️ [VERIFY] FailureType enum full variant list not shown — B-C1 Edge Cases; _runner.py:46/:496 cited but enum definition not in cited block
- ⚠️ Person model — only one field implied (name); full field list, types, required/optional not specified — B-B1/B-B2
- ⚠️ messy_cli_output.txt fixture embedded JSON schema not defined — B-C1/B-D2 assert 'the correct typed object' but no field list given
- ⚠️ Go baml_src/ BAML function and model definitions not specified — B-D1 only gestures at 'generators/clients/functions'
- ⚠️ clients.baml provider/model configuration schema not specified — only env-var names listed
- ⚠️ baml-py initial version pin not specified — B-A1 pyproject.toml; B-F1 enforces parity only after first commit
- ⚠️ Large-schema preserve-or-drop decision deferred with no schema for resulting behavior — B-D2 'decide & document' not bound to a deliverable

### Critical
- ❌ ExtractDynamic @@dynamic BAML function has no schema — B-B2 Test Specification / Files touched. The load-bearing function for all agent.ai(schema=) calls has no input params, no output type signature, no TypeBuilder registration pattern; baml_bridge.py and the agent_ai.py rewrite have no target interface
- ❌ Pydantic→BAML TypeBuilder type mapping rules unspecified — B-B2 Red / Edge Cases. No mapping table; test_baml_bridge.py oracle becomes whatever the implementation produces (self-referential verification), passing even for a mapping that drops optional fields or coerces enums to strings

### Recommendations
- ExtractDynamic function schema: Add a concrete signature block — e.g. `function ExtractDynamic(prompt: string) -> dynamic` with a note on TypeBuilder registration before the call
- Pydantic→BAML TypeBuilder mapping: Prepend a mapping table (str → tb.string(), Optional[X] → tb.optional(...), list[X] → tb.list(...), Enum → tb.enum(...)) as the independent oracle for test_baml_bridge.py
- Person model: Enumerate all fields with types and required/optional status in B-B1
- messy_cli_output.txt fixture schema: Add a one-line schema comment specifying embedded JSON field names and types
- Large-schema schema.go:50 decision: Make preserve/drop explicit before implementation, recorded as a constant/bool flag not prose
- baml-py version pin: Choose and record the initial version in pyproject.toml

## APIs Review

### Well-Defined
- ✅ baml-cli generate CLI invocation — B-A1: package exists, import succeeds, hasattr(b, F) for every declared F
- ✅ b.ExtractPerson(bio) static call returns typed Pydantic object — B-B1 fixed-input oracle 'result.name == "Ada Lovelace"'; skip-not-fail specified
- ✅ FailureType.SCHEMA preserved for unrecoverable parse input — B-C1 edge cases and refactor step
- ✅ make baml-generate target + git diff --exit-code drift gate — B-F1: Makefile target, CI step, stale-client and missing-PATH failure modes
- ✅ clients.baml reads existing env vars — BAML Replacement Mechanics
- ✅ Integration tests skip (not fail) when API key env var absent — B-B1/B-D1 edge cases

### Missing or Unclear
- ⚠️ baml_bridge.pydantic_to_typebuilder return type and full signature not specified — B-B2; red step unwritable without the TypeBuilder return shape
- ⚠️ ExtractDynamic @@dynamic function input parameter name(s) and output class declaration absent — B-B2 functions.baml
- ⚠️ Exception class raised by BAML on parse failure unspecified — B-B2/B-C1 edge cases; if it differs from ValueError the FailureType.SCHEMA mapping silently breaks
- ⚠️ Large-schema-on-disk preserve-or-drop decision explicitly deferred — B-D2 edge cases; Risks section
- ⚠️ How rate_limiter.py wraps BAML async calls only described as 'ensure rate-limit retry still wraps the call' — B-B2 refactor; integration seam undefined

### Critical
- ❌ Go SDK consumer-facing API after B-D1 unspecified — B-D1 removes WithSchema (request.go:258) and JSON()/Into() (response.go:99) but never defines the replacement calling convention; the word 'or' signals the decision is open. Any existing Go caller has no specified post-migration call site

### Recommendations
- Go SDK API migration (B-D1): Add one paragraph picking exactly (a) thin wrappers delegating to BAML (zero callers broken) or (b) deletion with callers switching to `b.ExtractFoo()` directly; resolve the ambiguity before the red step
- baml_bridge.py API contract (B-B2): Add `def pydantic_to_typebuilder(model: type[BaseModel]) -> baml_py.TypeBuilder` and one sentence on unsupported constructs
- BAML parse-error mapping (B-C1, B-B2): Name the BAML exception class and state it is caught and re-raised as the existing ValueError / mapped to FailureType.SCHEMA
- Large-schema-on-disk decision (B-D2): Resolve preserve-or-drop before the B-D2 red step with a single sentence

## Scaffolding & Realness (Hostile Check)

This plan's central risk is a cluster of **self-referential verification** oracles whose tests would pass even if the real mapping/parse functionality were wrong:

- ❌ **Pydantic→BAML TypeBuilder mapping has no independent ground-truth oracle** (B-B2 Red step / Edge Cases; test_baml_bridge.py). The plan labels this the "riskiest unit" yet gives no mapping table. The test "asserts the built type shape" of a TypeBuilder the code under test produced — its oracle is the same fiction the implementation emits. It passes for any mapping, including one that silently drops optional fields or coerces enums to strings. This is the textbook self-referential violation.
- ❌ **test_baml_bridge.py oracle inspects the code's own output** (B-B2 TDD Cycle, Red step: "maps a nested Pydantic model → TypeBuilder and asserts the built type shape"). No independently authored expected value is named; the test confirms the code produced something shaped like itself.
- ❌ **B-A1 property-test oracle derives expected names from the artifact under generation** (B-A1 Property specification: "for every function name F declared in baml_src/, hasattr(b, F)"). If F is read from baml_src/ at test time, the oracle is the same source the generator consumes — it cannot fail to find what it just enumerated.

Note: the fixture-based oracles (B-C1, B-D2 `parse(fixture) == embedded`) and the known-fact LLM oracles (B-B1 `result.name == 'Ada Lovelace'`) are genuinely independent and pass this check; the milestones do name human-observable artifacts. The violations above are isolated to the three TypeBuilder/property oracles, but each is sufficient to ship a green-but-wrong mapper.

## Critical Issues (Must Address Before Implementation)

Contracts:
- BAML parser-only API endpoint is never named for Python or Go (B-C1/B-D2 Green; BAML Replacement Mechanics)
- TypeBuilder → Pydantic instance cast mechanics unspecified (B-B2 Test Specification)

Interfaces:
- `baml_bridge.pydantic_to_typebuilder` parameter and return types not stated (B-B2)
- Replacement for `try_parse_from_text` in _schema.py has no name or signature (B-C1; _runner.py:332)
- Go `WithSchema` and `(*Response).JSON` removed with no replacement interface (B-D1)

Promises:
- Caller-specified timeout silently dropped after B-B2 migration (agent_ai.py:544)
- response.JSON() / Into() disposition in B-D1 is an unresolved OR (response.go:99)
- .gitignore / commit-vs-generate decision deferred in B-A1 — breaks CI in a clean checkout (B-A1)

Data Models:
- ExtractDynamic @@dynamic BAML function has no schema (B-B2)
- Pydantic→BAML TypeBuilder type mapping rules unspecified (B-B2)

APIs:
- Go SDK consumer-facing API after B-D1 unspecified (B-D1)

Scaffolding & Realness:
- Pydantic→BAML TypeBuilder mapping has no independent ground-truth oracle (B-B2)
- test_baml_bridge.py oracle inspects the code's own output (B-B2)
- B-A1 property-test oracle derives expected names from the artifact under generation (B-A1)

## Approval Status

- [ ] Ready for Implementation — no critical issues AND Scaffolding & Realness is clean
- [ ] Needs Minor Revision — address warnings before proceeding
- [x] Needs Major Revision — critical issues OR any Scaffolding & Realness violation must be resolved first


## Beads Tracking

Epic: **[Plan Review Epic] 2026-06-07-09-04-tdd-baml-structured-output — 11 critical, 32 warnings** — _(not minted; pass `--mint-beads` to create)_
Issues: **11** critical, **32** warning.

| Phase | Severity | Issue | Bead ID |
|---|---|---|---|
| contracts | critical | [Plan Review §contracts] BAML parser-only API endpoint is never named for Python or … <!-- fp:e18a1bbe --> | _(not minted)_ |
| contracts | critical | [Plan Review §contracts] TypeBuilder → Pydantic instance cast mechanics are unspecif… <!-- fp:7197f213 --> | _(not minted)_ |
| contracts | warning | [Plan Review §contracts] test_baml_bridge.py oracle is potentially self-referential.… <!-- fp:060b850e --> | _(not minted)_ |
| contracts | warning | [Plan Review §contracts] BAML exception class → existing exception type mapping is u… <!-- fp:c9ba88f0 --> | _(not minted)_ |
| contracts | warning | [Plan Review §contracts] rate_limiter.py wrapping contract with the BAML call is uns… <!-- fp:5387802c --> | _(not minted)_ |
| contracts | warning | [Plan Review §contracts] B-A1 property test oracle risks self-reference. 'For every … <!-- fp:eeab6920 --> | _(not minted)_ |
| contracts | warning | [Plan Review §contracts] Large-schema-on-disk behavior (schema.go:15,50) is explicit… <!-- fp:9c01794f --> | _(not minted)_ |
| interfaces | critical | [Plan Review §interfaces] `baml_bridge.pydantic_to_typebuilder` parameter and return … <!-- fp:09fda200 --> | _(not minted)_ |
| interfaces | critical | [Plan Review §interfaces] Replacement for `try_parse_from_text` in _schema.py has no … <!-- fp:b8fc1b85 --> | _(not minted)_ |
| interfaces | critical | [Plan Review §interfaces] Go public functions `WithSchema` and `(*Response).JSON` are… <!-- fp:eaf1cd66 --> | _(not minted)_ |
| interfaces | warning | [Plan Review §interfaces] Go BAML-generated client call expression never shown — afte… <!-- fp:92ea4ce3 --> | _(not minted)_ |
| interfaces | warning | [Plan Review §interfaces] agent_ai.py BAML call expression undefined — plan says 'cal… <!-- fp:f730e28a --> | _(not minted)_ |
| interfaces | warning | [Plan Review §interfaces] Exception-type compatibility between BAML parse errors and … <!-- fp:071bb25e --> | _(not minted)_ |
| interfaces | warning | [Plan Review §interfaces] Large-schema-on-disk decision for B-D2 is explicitly deferr… <!-- fp:46878a14 --> | _(not minted)_ |
| interfaces | warning | [Plan Review §interfaces] baml_bridge.py module-level visibility not stated — only py… <!-- fp:e35413f7 --> | _(not minted)_ |
| promises | critical | [Plan Review §promises] Caller-specified timeout is silently dropped after B-B2 mig… <!-- fp:a2618999 --> | _(not minted)_ |
| promises | critical | [Plan Review §promises] response.JSON() / Into() disposition in B-D1 is an unresolv… <!-- fp:49338844 --> | _(not minted)_ |
| promises | critical | [Plan Review §promises] .gitignore / commit-vs-generate decision deferred in B-A1 i… <!-- fp:fc4d5e9d --> | _(not minted)_ |
| promises | warning | [Plan Review §promises] Retry behavior after BAML parse failure in B-C1 is unspecif… <!-- fp:d59530d5 --> | _(not minted)_ |
| promises | warning | [Plan Review §promises] Sync vs. async BAML Python client choice not specified for … <!-- fp:743e8882 --> | _(not minted)_ |
| promises | warning | [Plan Review §promises] Python exception type emitted on BAML parse/validation fail… <!-- fp:d65e8e50 --> | _(not minted)_ |
| promises | warning | [Plan Review §promises] Idempotency of baml-cli generate not stated. The drift gate… <!-- fp:04c6da6b --> | _(not minted)_ |
| promises | warning | [Plan Review §promises] [VERIFY] BAML Go runtime module path and go 1.21 compatibil… <!-- fp:699ae59a --> | _(not minted)_ |
| promises | warning | [Plan Review §promises] Go BAML client goroutine-safety not addressed. sdk/go/ai ma… <!-- fp:f7ff09d0 --> | _(not minted)_ |
| promises | warning | [Plan Review §promises] Large-schema token gate (schema.go:50, largeSchemaTokenThre… <!-- fp:0a8b3715 --> | _(not minted)_ |
| promises | warning | [Plan Review §promises] B-B2 unit test oracle for test_baml_bridge.py is structural… <!-- fp:3b474c27 --> | _(not minted)_ |
| promises | warning | [Plan Review §promises] Recursive Pydantic models and Union types with >2 variants … <!-- fp:b917482a --> | _(not minted)_ |
| promises | warning | [Plan Review §promises] pyproject.toml placement of baml-py (runtime dep vs. build-… <!-- fp:54ea9d1c --> | _(not minted)_ |
| data_models | critical | [Plan Review §data_models] ExtractDynamic @@dynamic BAML function has no schema — the … <!-- fp:d3fe4a12 --> | _(not minted)_ |
| data_models | critical | [Plan Review §data_models] Pydantic→BAML TypeBuilder type mapping rules are unspecifie… <!-- fp:306d62f6 --> | _(not minted)_ |
| data_models | warning | [Plan Review §data_models] [VERIFY] FailureType enum — full variant list not shown in … <!-- fp:fe936ae2 --> | _(not minted)_ |
| data_models | warning | [Plan Review §data_models] Person model — only one field implied (name == 'Ada Lovelac… <!-- fp:754b57bc --> | _(not minted)_ |
| data_models | warning | [Plan Review §data_models] messy_cli_output.txt fixture embedded JSON schema not defin… <!-- fp:f02e52d9 --> | _(not minted)_ |
| data_models | warning | [Plan Review §data_models] Go baml_src/ BAML function and model definitions not specif… <!-- fp:48509a97 --> | _(not minted)_ |
| data_models | warning | [Plan Review §data_models] clients.baml provider/model configuration schema not specif… <!-- fp:ea3c634b --> | _(not minted)_ |
| data_models | warning | [Plan Review §data_models] baml-py initial version pin not specified — plan adds baml-… <!-- fp:48524846 --> | _(not minted)_ |
| data_models | warning | [Plan Review §data_models] Large-schema preserve-or-drop decision deferred with no sch… <!-- fp:b83ab4f1 --> | _(not minted)_ |
| apis | critical | [Plan Review §apis] Go SDK consumer-facing API after B-D1 is unspecified. The p… <!-- fp:923c829a --> | _(not minted)_ |
| apis | warning | [Plan Review §apis] baml_bridge.pydantic_to_typebuilder return type and full si… <!-- fp:7a0f09c0 --> | _(not minted)_ |
| apis | warning | [Plan Review §apis] ExtractDynamic @@dynamic BAML function input parameter name… <!-- fp:6e31d61d --> | _(not minted)_ |
| apis | warning | [Plan Review §apis] Exception class raised by BAML on parse failure is unspecif… <!-- fp:cc330445 --> | _(not minted)_ |
| apis | warning | [Plan Review §apis] Large-schema-on-disk behavior (schema.go:15 schemaFilename … <!-- fp:5bd83fca --> | _(not minted)_ |
| apis | warning | [Plan Review §apis] How rate_limiter.py wraps BAML async calls is described onl… <!-- fp:d6324379 --> | _(not minted)_ |

Run `bun SAI/skills/PlanReview/Tools/PlanReview.ts --plan-path <path> --mint-beads` to mint these as a beads epic + child issues. The CLI will substitute real Bead IDs into this section.