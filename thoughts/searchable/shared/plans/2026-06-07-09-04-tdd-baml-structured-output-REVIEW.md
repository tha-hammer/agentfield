---
date: 2026-06-08
reviewer: claude (review_plan)
plan: thoughts/searchable/shared/plans/2026-06-07-09-04-tdd-baml-structured-output.md
type: pre-implementation-review
status: needs-minor-revision
---

# Plan Review Report: BAML Structured-Output Integration (TDD)

Pre-implementation review of the BAML replacement plan. Every load-bearing
`file:line` claim and the BAML API assumptions were verified against the working
tree and the installed `baml-py` on 2026-06-08 via four parallel verification
agents.

## Review Summary

| Category | Status | Issues Found |
|----------|--------|--------------|
| Contracts | ⚠️ | 2 (exception translation, parser-only return type) |
| Interfaces | ⚠️ | 1 (Go line drift, cosmetic) |
| Promises | ✅ | 0 (timeout/retry/concurrency all addressed) |
| Data Models | ✅ | 0 (Pydantic→TypeBuilder table is complete + has independent oracle) |
| APIs | ⚠️ | 1 critical (env mismatch) + plan's own tracked risks |

**Verdict: Needs Minor Revision.** The plan is exceptionally rigorous —
independent oracles, blast-radius notes, resolved decisions, verified provenance.
27 of ~30 spot-checked `file:line` claims are EXACT. But one factual claim is
**wrong**, two exception/return-type contracts are **underspecified or
contradictory**, and the BAML environment was verified against the **wrong venv**.
None are architectural; all are fixable with text edits before B-A1.

---

## Critical Issues (Must Address Before Implementation)

### ❌ C1 — BAML verified in the wrong venv; SDK's real test env has no baml-py and excludes py3.14

The plan repeatedly asserts "baml-py is already in the venv (`.venv/bin/baml-cli`)"
(B-A1 green, B-F1, line 118, 292, 546). Verification found:

- `baml-py==0.222.0` **is** installed — but in `/home/maceo/Dev/agentfield/.venv`
  (**python 3.14**, repo root).
- The SDK's own venv `/home/maceo/Dev/agentfield/sdk/python/.venv` (**python 3.12**)
  has **no baml_py, no dist-info, no `baml-cli`**. This is where `pytest` /
  `pip install -e .[dev]` actually run.
- `sdk/python/pyproject.toml:26` declares `requires-python = ">=3.10, <3.14"` —
  the py3.14 root venv where BAML was verified is **outside the SDK's supported
  range**.

**Impact:** Every "baml-py already present, just run `baml-cli generate`" green
step is false for the environment the SDK actually builds/tests in. B-A1 cannot go
green as written. The B-F1 Makefile `$(VENV)/bin/baml-cli` will resolve to a venv
without the binary. Whether `baml-py==0.222.0` even installs and its `.abi3.so`
loads on **python 3.12** is unverified (it almost certainly does — abi3 wheel —
but it is an unproven precondition for the whole plan).

**Recommendation:**
1. Add `baml-py==0.222.0` to `pyproject.toml` runtime deps (the plan already
   commits to this in Resolved Decisions §1 — make it B-A1's *first* green step,
   before `baml-cli generate`, not an afterthought).
2. Add a B-A1 precondition test: `baml-py` imports and `baml-cli` runs **in
   `sdk/python/.venv` (py3.12)** — re-verify the API symbols there, not in the
   root py3.14 venv.
3. Pin down which venv `$(VENV)` is in the Makefile target (B-F1). State it
   explicitly as `sdk/python/.venv`.
4. Decide the py3.14 question: either the root venv is irrelevant to the SDK
   (likely) and should not be cited as evidence, or `requires-python` needs
   revisiting. Today the plan's evidence venv contradicts the package's declared
   support range.

---

## Contract Review

### ⚠️ I1 — Exception-translation contract is wrong and self-contradictory (claim at `agent_ai.py:808`)

The `deserialize` docstring (plan lines 465-469) says it returns
`model.model_validate(...)` *"letting pydantic.ValidationError propagate (caught at
agent_ai.py:808)"*. Verified against real code:

- `agent_ai.py:808` does **not catch** `ValidationError` — it **raises**
  `ValueError("Could not parse structured response: ...")`. `ValidationError` is
  caught and swallowed earlier (lines 792-796 and 806).
- That `ValueError` (with the specific substring) is what the retry loop at
  `:816-824` matches on (`"Could not parse structured response" in str(e)`) to
  decide whether to retry.
- **B-B2 deletes lines 788-824 entirely** (the parse ladder *and* that retry
  loop). So the line being cited as the catch site no longer exists post-change.

**Why it matters (real behavior change):** Today callers of `agent.ai(schema=)`
observe `ValueError` on parse failure. If `deserialize` lets a raw
`pydantic.ValidationError` (or BAML's `BamlValidationError`) propagate, the public
exception type **silently changes**, and any caller catching `ValueError` breaks.
The B-B2 edge-case text *contradicts* the docstring: it says
"`BamlValidationError` surfaced then **re-raised as the existing `ValueError`**".

**Recommendation:** Pick one and state it as the contract:
- If preserving the `ValueError` contract: `deserialize` (or the agent_ai branch)
  must catch `pydantic.ValidationError` / `baml_py.errors.BamlValidationError` and
  re-raise `ValueError("Could not parse structured response: ...")`. Remove the
  "propagate, caught at 808" sentence (that line is deleted).
- If intentionally changing to `BamlValidationError`: call it out as a breaking
  exception-contract change in "What We're NOT Doing" / Risks, and grep for
  existing `except ValueError` callers first.
- Also specify: with the `:816-824` retry loop deleted, **what now retries a
  schema-parse failure?** BAML's `clients.baml` `max_retries` covers HTTP, not
  necessarily SAP validation failures. Either add a BAML `retry_policy` or state
  that parse-failure retries are intentionally dropped.

### ⚠️ I2 — B-C1 parser-only returns a BAML `DynamicOutput`, not a `schema`-typed instance (inconsistent with B-B2)

B-C1's replacement body (plan lines 374-379) is:
```python
return b.parse.ExtractDynamic(text, baml_options={"tb": tb})
```
This returns BAML's dynamic `DynamicOutput` object. But B-B2 establishes that
callers should *always* get a `model`-validated instance via `deserialize(raw,
schema)` (plan lines 471-473: "callers always get `model`-validated instances …
never a raw dict"). B-C1 skips `deserialize`, so `try_parse_from_text` would
return a `DynamicOutput`, not an instance of `schema`. Its callers
(`parse_and_validate`, `_runner.py:355/471`) and everything downstream expect the
target schema type.

**Recommendation:** Make B-C1 symmetric with B-B2:
```python
raw = b.parse.ExtractDynamic(text, baml_options={"tb": tb})
return deserialize(raw, schema)
```
or explicitly document that the harness path consumes `DynamicOutput` and verify
downstream (`parse_and_validate`'s isinstance/validation) accepts it. As written
the two milestones disagree on the return type of the parsed object.

---

## Interface Review

### ⚠️ I3 — Go `file:line` citations have drifted (cosmetic, but fix to avoid misedits)

Substance of every Go claim is correct; two line numbers are stale:
- `goTypeToJSONType` is at `request.go:456`, not 455 (455 is its doc comment).
- `FailureSchema` mapping is at `runner.go:410`, not 401 (401 is an
  `accumulateMetrics` call).

Python citations all matched exactly except the harness call sites:
`try_parse_from_text` is invoked at `_runner.py:355` and `:471`, **not** `:348`
(348 is `parse_and_validate`). `rate_limiter.py` lives at
`agentfield/rate_limiter.py`, not under `harness/`.

**Recommendation:** Implementers should locate by symbol name, not raw line
number (the plan's own blast-radius notes already encourage this). Optionally
refresh the three drifted citations. Non-blocking.

---

## Promise Review — ✅ Well-Defined

- ✅ **Timeout** — correctly identifies BAML has no per-call timeout (issue #1630,
  confirmed open at 0.222.0); preserves the real `asyncio.wait_for(...,
  effective_timeout*2)` guard at `agent_ai.py:673` (verified EXACT) + client-level
  `request_timeout_ms`. Sound.
- ✅ **Rate-limit retry** — routes through `execute_with_retry` at
  `rate_limiter.py:209` (verified EXACT); the detector `_is_rate_limit_error`
  (`rate_limiter.py:55-98`, verified) matches `RateLimitError`/429/503. Plan
  correctly adds `BamlClientHttpError`/429 to it.
- ✅ **Concurrency** — Go goroutine-safety unknown ⇒ B-D1 adds an N=8 concurrency
  smoke test rather than assuming safety. Correct posture.
- ✅ **`FailureType.SCHEMA` preservation** — verified: the bare `except Exception`
  in `parse_and_validate` (`_schema.py:196,203`) absorbs `BamlValidationError`,
  and `_runner.py:498/504` maps `None`→`FailureType.SCHEMA`. The plan's
  refutation of the prior review finding here is **correct**.

---

## Data Model Review — ✅ Well-Defined

- ✅ Pydantic→TypeBuilder mapping table (plan lines 223-234) is complete:
  primitives, optional, list, dict, nested BaseModel, enum, literal, union.
- ✅ Independent hand-authored oracle (`test_baml_bridge.py`) is genuinely
  independent — uses `b.parse` (no network), asserts against hardcoded `HARDCODED`
  / `M(**HARDCODED)`, not the mapper's own output. This resolves the prior
  review's self-reference finding correctly.
- ✅ Edge limits documented (recursive models registered-then-referenced;
  `Any`/callables raise `TypeError`, not silent coercion).
- ⚠️ Minor: the bridge's *top-level* attachment to the `@@dynamic DynamicOutput`
  class (vs `add_class` for nested models) isn't shown explicitly. The table
  covers nested classes; the dynamic-root attachment mechanism
  (`tb.DynamicOutput.add_property(...)`) is implied but not spelled out. Worth one
  concrete line in B-B2.

---

## API Review

- ✅ BAML Python API symbols (`errors.BamlValidationError`, `BamlClientHttpError`,
  `AbortController`, `ClientRegistry`, `b.parse.Fn`, `TypeBuilder`) — **verified
  present** in the installed `baml_py` source/stubs. No fictional API.
- ✅ `baml-cli` ships in the wheel (`.venv/bin/baml-cli`) — confirmed.
- ⚠️ Go runtime `github.com/boundaryml/baml v0.222.0` existence as a Go module tag
  was **not** verified (no `go get` run). Plan flags it pre-stable. `go.mod` is
  `go 1.21` (verified) but BAML's module declares `go 1.24.0` — if the build
  toolchain is <1.24, `go build` fails. Plan tracks this as a B-D1 green gate;
  acceptable, but it is a genuine unproven precondition (parallel to C1 on the Go
  side). Confirm the Go toolchain version in CI ≥ what BAML requires.

---

## Suggested Plan Amendments

```diff
# B-A1 (Generate Python BAML client)
+ Green step 0: add baml-py==0.222.0 to sdk/python/pyproject.toml deps, then
+   `pip install -e .[dev]` INTO sdk/python/.venv (py3.12); re-verify baml_py
+   imports + baml-cli runs THERE (not the root py3.14 venv).
~ Replace "baml-py is already in the venv" with the SDK venv it actually targets.

# B-B2 (agent.ai reroute)
~ Fix deserialize contract: catch pydantic.ValidationError / BamlValidationError
+   and re-raise ValueError("Could not parse structured response: ...") to keep
+   the public exception contract — OR declare the exception-type change in Risks.
- Remove "(caught at agent_ai.py:808)" — that line raises ValueError and is deleted.
+ State what retries a parse failure now that :816-824 is removed (BAML retry_policy
+   or "retries intentionally dropped").
+ Show the @@dynamic root attachment line (tb.DynamicOutput.add_property...).

# B-C1 (harness parser-only)
~ try_parse_from_text: return deserialize(b.parse.ExtractDynamic(...), schema),
+   not the raw DynamicOutput — make it symmetric with B-B2's return-type contract.

# B-F1 (build/CI gate)
~ Pin $(VENV) = sdk/python/.venv explicitly.

# B-D1 (Go)
+ Confirm CI Go toolchain >= BAML's required go directive (module declares 1.24.0)
+   as an explicit green gate, not just go.mod's 1.21.
```

## Approval Status

- [ ] Ready for Implementation
- [x] **Needs Minor Revision** — address C1 (env), I1 (exception contract), I2
      (parser-only return type) before starting B-A1. I3 (line drift) and the Go
      toolchain confirmation are non-blocking but should be folded in.
- [ ] Needs Major Revision

**Bottom line:** Architecture is sound and the replacement strategy is correct.
The blockers are precondition/contract precision, not design. Fix the four text
amendments above and this is ready.
