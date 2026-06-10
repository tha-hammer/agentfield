---
date: 2026-06-07T09:04:53-04:00
author: maceo
git_commit: fb8559d5d88d42e8b15043ff746611ee482152ba
branch: main
repository: agentfield
topic: "Replace hand-rolled structured-output parsing with BAML across Python + Go SDKs"
tags: [tdd, plan, baml, structured-output, python-sdk, go-sdk, harness]
status: ready
research: thoughts/searchable/shared/research/2026-06-07-08-32-baml-structured-output-integration-surface.md
tier: L1
---

# BAML Structured-Output Integration — TDD Implementation Plan

## Overview

Replace the codebase's hand-rolled structured-output generation and parsing with
**BAML** ("prompt as function") across the **Python SDK** (inline AI + harness)
and the **Go SDK** (`ai` + harness), plus a build/CI codegen gate. BAML defines
LLM calls as typed functions in `.baml` files, generates a typed client via
`baml-cli generate`, and parses model output with a tolerant schema parser —
removing the existing `json.loads + regex-fallback + retry` ladders.

This is a **replace**, not a coexistence: the old parse paths are deleted once the
BAML path is green. Each milestone ships a **real, user-observable** capability
(an importable client, a real `.ai()` call returning a typed object, a messy
fixture parsed correctly) — no scaffolding, no shape-only tests.

**Scope decided with the user:** all SDK surfaces (Python + Go), replace strategy,
provider config via BAML's own env-var reads. **Web UI is out of scope** — its
zod usage validates *runtime-dynamic* form input built from backend JSON Schema,
which BAML (static, build-time codegen) structurally cannot replace, and the UI
makes no direct LLM calls.

## Current State Analysis

BAML is **wholly absent** today (no `baml` refs in `sdk/`, `control-plane/`,
`Makefile`; no `baml_src/`/`baml_client/`; not in `pyproject.toml`/`go.mod`).
Structured output is hand-rolled at four surfaces; this plan touches three of
them (Python inline, Python harness, Go SDK+harness).

### Key Discoveries (verified, file:line)

- **Python inline AI** — `sdk/python/agentfield/agent_ai.py:546-557` attaches a
  strict `response_format: json_schema` from `schema.model_json_schema()`;
  `:788-810` does `json.loads(...)` → `schema(**data)` with a `r"\{.*\}"` regex
  fallback; `:729` `max_parse_retries = 2`; `:816-824` the active retry loop.
- **Python harness (parser-only target)** —
  `sdk/python/agentfield/harness/_schema.py:209` `try_parse_from_text` (3-strategy
  text scrape); `_runner.py:46` `DEFAULT_SCHEMA_RETRIES=2`, `:59` `_ai_schema_repair`,
  `:332` `_handle_schema_with_retry`, `:496` final schema-failure message.
  Providers shell out to real CLIs (`harness/providers/opencode.py:152`
  `cmd=[self._bin,"run"]`), wired in `harness/providers/_factory.py:9`.
- **Go SDK** — `sdk/go/ai/request.go:258` `WithSchema`, `:390`
  `structToJSONSchema` (reflection), `:290` `Strict:true`; `sdk/go/ai/response.go:99`
  `Into()`/`JSON()` via `json.Unmarshal`; `sdk/go/harness/schema.go` writes
  `.agentfield_schema.json` (`:15`) with a large-schema gate (`:50`).
- **Deps / tooling** — `sdk/python/pyproject.toml:31-32` `pydantic>=2.0`,
  `litellm`; pytest markers `unit`/`functional`/`integration` (`:77-81`),
  `pytest-asyncio`. Go module `github.com/Agent-Field/agentfield/sdk/go` (go 1.21),
  `go test` with `*_integration_test.go` convention. `bd` present for beads.
- **Web UI (out of scope, for the record)** —
  `control-plane/web/client/src/utils/jsonSchemaToZod.ts` +
  `schemaUtils.ts:318` build zod from backend JSON Schema at runtime; no AI deps
  in `package.json`.

## Desired End State

An agent author calls structured LLM functions through BAML and gets validated,
typed objects with no hand-rolled parsing in the SDK. Verify by:

### Observable Behaviors
- Running `baml-cli generate` produces an importable `baml_client`.
- A real `agent.ai(prompt, schema=Model)` returns the correct typed object via
  BAML; the old `json.loads + regex + retry` code is gone.
- Messy recorded CLI output parses into the correct typed object via BAML's
  parser-only path; `try_parse_from_text` and the Go `schema.go` post-hoc path
  are replaced.
- `make build` regenerates the client and a CI drift check stays green.
- `go test` returns a correct struct from a real call and parses a recorded
  fixture; the reflection path is replaced.

## What We're NOT Doing

- **Web UI zod** — runtime-dynamic form validation; no BAML target. Left as-is.
- **Control plane (Go server)** — makes no LLM calls; untouched.
- **No coexistence layer** — the old parse ladders are deleted, not kept as
  fallback (per "replace" decision).
- **No new providers/models** — reuse existing API keys via BAML env-var reads.

## BAML Replacement Mechanics (design notes)

- **Static functions** (`baml_src/*.baml`) cover fixed, known schemas.
- **Dynamic schemas** — `agent.ai(schema=RuntimePydanticModel)` passes a *runtime*
  Pydantic class. Route it through a **generic `@@dynamic` BAML function** plus
  **TypeBuilder**, which maps the Pydantic model's fields into BAML's output type
  at call time. This is the load-bearing piece that lets a single BAML function
  serve arbitrary caller schemas.
- **Parser-only** — for harness output (third-party CLI made the call, not BAML),
  feed raw text to BAML's tolerant parser to produce a typed object. BAML can
  parse text it did not generate.
- **Provider config** — `clients.baml` references `env.ANTHROPIC_API_KEY` /
  `env.OPENAI_API_KEY` / `env.OPENROUTER_API_KEY`, reading the same environment
  the SDK already uses.

## BAML API Reference (verified against installed `baml-py==0.222.0`)

> Every BAML symbol used below was verified against the installed runtime
> (`baml_py` version `0.222.0`) and/or the BoundaryML docs + generated-client
> source. Symbols that do **not** exist are called out explicitly so no milestone
> assumes a fictional API.
>
> **Venv caveat (Resolved Decisions §7):** `baml-py==0.222.0` was originally
> confirmed in the **repo-root** `/home/maceo/Dev/agentfield/.venv` (**python
> 3.14**) — which is *not* where the SDK builds or tests. The SDK's own venv
> `sdk/python/.venv` (**python 3.12**, where `pip install -e .[dev]` and `pytest`
> run) did **not** have `baml-py` until this plan adds it. The root py3.14 venv is
> also **outside** the SDK's declared `requires-python = ">=3.10, <3.14"`
> (`sdk/python/pyproject.toml:26`) and is therefore **not** valid evidence for the
> SDK target. **B-A1 re-verifies every symbol in `sdk/python/.venv` (py3.12)**
> after installing `baml-py` there (abi3 wheel — expected to load on 3.12, proven
> in B-A1 rather than assumed).

### Tooling
- **`baml-cli`** ships inside the `baml-py` wheel — once `baml-py` is installed
  into the SDK venv it lands at `sdk/python/.venv/bin/baml-cli` (and
  `sdk/python/.venv/bin/baml`), NOT a separate global install. CI/`make` must
  invoke it via that venv (`python -m baml_py` is **not** an entrypoint; use the
  `baml-cli` console script the wheel installs). The Makefile target pins
  `$(VENV) := sdk/python/.venv` (B-F1).
- **Version pin:** `baml-py==0.222.0` (latest stable, 2026-04-27). The
  `generator.version` field in every `.baml` MUST equal this string.

### Python (generated `baml_client`)
- **Sync vs async:** both clients are always generated. agentfield's `agent.ai`
  is `async def` (`agent_ai.py:257`), so use the **async** client:
  `from baml_client.async_client import b` → `await b.Fn(...)`. Harness
  parser-only paths are synchronous and use `from baml_client.sync_client import b`.
- **TypeBuilder** (dynamic schemas): imported from the **generated** package, not
  `baml_py`: `from baml_client.type_builder import TypeBuilder`. Builder methods
  (verified): `tb.string() | tb.int() | tb.float() | tb.bool()` →`FieldType`;
  `tb.list(ft) | tb.union([ft,...]) | tb.map(k,v)` →`FieldType`;
  `tb.literal_string(s)/literal_int(i)/literal_bool(b)`; `tb.add_class(name)`
  →`ClassBuilder`; `tb.add_enum(name)` →`EnumBuilder`. On `FieldType`:
  `.optional()`, `.list()`. On `ClassBuilder`: `.add_property(name, ft)`
  →`ClassPropertyBuilder`, `.type()` →`FieldType`. On `EnumBuilder`:
  `.add_value(v)`.
- **Passing the TypeBuilder at call time:** the last param of every generated
  function is `baml_options: BamlCallOptions = {}` (a `TypedDict`). Keys:
  `"tb"`, `"client_registry"`, `"client"`, `"collector"`, `"env"`, `"tags"`,
  `"on_tick"`. Call: `await b.ExtractDynamic(prompt, baml_options={"tb": tb})`.
- **Parser-only (no LLM call):** the generated client exposes `b.parse` →
  `LlmResponseParser`, with one method per function:
  `b.parse.Fn(llm_response: str, baml_options={}) -> ReturnType`. It runs BAML's
  Schema-Aligned Parser on arbitrary text with **no** model call. Verified present
  in the generated client (canary source) and available on both sync and async
  clients. For a dynamic schema, pass the TypeBuilder too:
  `b.parse.ExtractDynamic(text, baml_options={"tb": tb})`.
- **Exceptions** (`from baml_py import errors`, verified present in 0.222.0):
  `BamlValidationError` (output fails schema validation / `@@assert` — has
  `.prompt`, `.raw_output`, `.message`), `BamlClientError`, `BamlClientHttpError`,
  `BamlClientFinishReasonError`, `BamlInvalidArgumentError`, `BamlTimeoutError`,
  `BamlAbortError`, and base `BamlError`. **Parse/validation failures raise
  `errors.BamlValidationError`** — this is the type harness/inline code must catch
  and translate.
- **Per-call timeout / ClientRegistry:** `ClientRegistry` exists in `baml_py`
  and is passable as `baml_options={"client_registry": cr}`. **There is no
  per-call `timeout` option** (BAML issue #1630, still open at 0.222.0). Timeouts
  are configured in the `client<llm>` block in `clients.baml`; cancellation uses
  `AbortController` (`baml_py.AbortController`). See "Timeout handling" decision below.

### Go (generated `baml_client`)
- **Status:** Go is a real, documented, but **pre-stable** target (`v0.x`, BETA
  since v0.84.0). Runtime is a downloaded Rust shared lib via C-FFI:
  **Linux & macOS, amd64/arm64 only — no Windows** (the runtime exports
  `ErrNotSupportedPlatform`). This is acceptable: agentfield CI/dev is Linux.
- **Runtime dep:** `github.com/boundaryml/baml v0.222.0` (`go get`). BAML's own
  module declares `go 1.24.0`, which **forces the consuming toolchain to be
  ≥1.24** — so `sdk/go/go.mod`'s current `go 1.21` **must be bumped to `go 1.24`**
  (the generics floor of 1.18 is not the binding constraint). Resolved as a B-D1
  green gate (see B-D1), not an open item.
- **Generator block needs `client_package_name`** (the Go module path) in addition
  to `output_type "go"`.
- **Call expression:** generated functions are **package-level** (no client
  struct/constructor), first arg `context.Context`, then inputs, then variadic
  options: `result, err := b.ExtractPerson(ctx, bio)`; with options
  `b.ExtractPerson(ctx, bio, b.WithClient("anthropic/claude"), b.WithTypeBuilder(tb))`.
- **Parser-only (no LLM call):** exposed as `b.Parse.Fn(...)` (added v0.204.0),
  same arg shape, routes to the parse-only path.
- **TypeBuilder:** `tb, err := b.NewTypeBuilder()`; `tb.AddClass/AddEnum/...`;
  passed via `b.WithTypeBuilder(tb)`.
- **Goroutine safety:** not formally documented. A CFFI deadlock/panic under high
  concurrency was fixed in v0.220.0 (#3185). Treat as concurrent-capable at
  0.222.0 but add a concurrency smoke test (B-D1) rather than assume safety.

## Resolved Decisions (previously deferred — now bound)

These close every "decide & document" / unresolved-OR finding from the plan review:

1. **Codegen: commit the generated client (`baml_client/`), do NOT gitignore it.**
   Committing avoids a build-before-test ordering hazard and lets a clean checkout
   run tests without a generate step. The drift gate (B-F1) keeps it honest. So:
   `.gitignore` does **not** list `baml_client/`; `sdk/python/pyproject.toml` adds
   `baml-py==0.222.0` as a **runtime** dependency (the generated client imports
   `baml_py` at import time). Adding the dep and installing it **into
   `sdk/python/.venv`** is the **first** green step of B-A1 (before `baml-cli
   generate`) — not an afterthought; see Resolved Decisions §7.
2. **Go `WithSchema` + `(*Response).JSON`/`Into`: keep as thin deprecated wrappers
   that delegate to the BAML path** (lower-blast-radius L1 choice — zero Go callers
   break). `request.go:258 WithSchema` keeps its `func(schema interface{}) Option`
   signature but builds the request to route through the generated BAML function;
   `response.go:100 JSON(dest)` / `:109 Into(dest)` keep `func(*Response, interface{}) error`
   but unmarshal via `b.Parse.<Fn>` instead of `json.Unmarshal`. The new
   first-class API is the generated `b.ExtractFoo(ctx, ...)`; the wrappers carry a
   `// Deprecated: use the generated baml_client function directly` comment.
3. **Harness large-schema-on-disk gate (`schema.go:16 largeSchemaTokenThreshold = 4000`):
   DROP it.** BAML embeds the output schema in its own generated prompt; the
   manual `.agentfield_schema.json` write + prompt-suffix is obsolete once parsing
   is BAML's job. Remove `schemaFilename`/`largeSchemaTokenThreshold` and the
   >4000-token branch. Recorded as code deletion, not a runtime flag.
4. **Timeout handling:** because BAML has no per-call timeout option, set
   `request_timeout`/`max_retries` in the `client<llm>` blocks in `clients.baml`,
   and preserve the existing `asyncio.wait_for(...)` wall-clock guard in
   `agent_ai.py:673` around the `await b.Fn(...)` call so the caller-supplied
   `timeout=` is still honored. The `effective_timeout` value (`agent_ai.py:539-544`)
   continues to drive `asyncio.wait_for`; it no longer reaches the HTTP layer
   (acceptable — the wall-clock guard is the contract callers actually observe).
5. **Sync/async split:** inline `agent.ai` (async) → `baml_client.async_client`;
   harness parser-only (sync) → `baml_client.sync_client`. Stated per behavior.
6. **Exception contract + parse-retry preservation (PRESERVE, do not change the
   public type).** Today callers of `agent.ai(schema=)` observe a
   `ValueError("Could not parse structured response: ...")` on parse failure
   (`agent_ai.py:808` **raises** it — it does *not* catch `ValidationError`; the
   `ValidationError`/`json.JSONDecodeError` are swallowed earlier at `:792-796`,
   `:806`). The `:814-827` retry loop keys on that exact substring
   (`"Could not parse structured response" in str(e)`) to re-invoke the model.
   **Decision:** keep the public `ValueError` contract. `baml_bridge.deserialize`
   **catches** `pydantic.ValidationError` **and** `baml_py.errors.BamlValidationError`
   and **re-raises** `ValueError("Could not parse structured response: ...")` (same
   message prefix). Because the message is preserved, the **existing `:814-827`
   parse-retry loop + `:729 max_parse_retries` are KEPT as-is** — they still fire
   on a BAML parse failure and re-run the model for fresh output (current
   behavior, zero contract change). Only the `:788-810` `json.loads`+regex parse
   *ladder* is deleted, not the retry orchestration. (Resolves review I1: removes
   the false "propagate, caught at agent_ai.py:808" claim and the B-B2/docstring
   contradiction; answers "what retries a parse failure now" → the preserved
   `ValueError`-keyed loop.)
7. **Build/test environment: the SDK target is `sdk/python/.venv` (py3.12), not
   the repo-root `.venv` (py3.14).** B-A1's first green step adds
   `baml-py==0.222.0` to `sdk/python/pyproject.toml` and runs `pip install -e
   .[dev]` **into `sdk/python/.venv`**, then re-verifies (`import baml_py`,
   `baml-cli --version`, and the API symbols in the API Reference) **there** — the
   abi3 wheel is expected to load on py3.12 but this is **proven in B-A1, not
   assumed**. The root py3.14 venv is irrelevant to the SDK and is not cited as
   evidence (it falls outside `requires-python = ">=3.10, <3.14"`). `$(VENV)` in
   the Makefile (B-F1) is pinned to `sdk/python/.venv`. (Resolves review C1.)

## Pydantic → TypeBuilder mapping (independent oracle for `test_baml_bridge.py`)

`baml_bridge.pydantic_to_typebuilder` maps a runtime Pydantic model to a BAML
dynamic class on a `TypeBuilder`. This table is the **independent, hand-authored
ground truth** the unit test asserts against (it is NOT derived from the mapper's
own output — see B-B2):

| Pydantic field type | TypeBuilder expression |
|---|---|
| `str` | `tb.string()` |
| `int` | `tb.int()` |
| `float` | `tb.float()` |
| `bool` | `tb.bool()` |
| `Optional[X]` / `X \| None` | `<map(X)>.optional()` |
| `list[X]` | `tb.list(<map(X)>)` |
| `dict[K, V]` | `tb.map(<map(K)>, <map(V)>)` |
| nested `BaseModel` `M` | `tb.add_class("M")` + `.add_property` per field, then `.type()` |
| `Enum E` | `tb.add_enum("E")` + `.add_value` per member, then `.type()` |
| `Literal["a","b"]` | `tb.union([tb.literal_string("a"), tb.literal_string("b")])` |
| `Union[A, B, ...]` (≥2) | `tb.union([<map(A)>, <map(B)>, ...])` |

**Top-level (`@@dynamic` root) attachment — explicit.** The table above covers
how each *field type* maps. The model's own fields attach to the dynamic root
class `DynamicOutput` (the `@@dynamic` class in `functions.baml`) via
`tb.DynamicOutput.add_property(field_name, <map(field_type)>)` for each top-level
field — distinct from `tb.add_class("M")` which is used only for **nested**
BaseModels. So `pydantic_to_typebuilder(M)` walks `M.model_fields`, mapping each
field type per the table, and calls `tb.DynamicOutput.add_property(name, ft)` for
each; nested models become their own `add_class` registrations referenced by the
parent property. (Resolves review Data-Model minor: the dynamic-root attachment
is now spelled out, not implied.)

**Explicit coverage limits (documented, not silently dropped):** recursive
models (a model referencing itself) are supported by registering the class name
first, then adding the self-referential property by name. Unsupported field
constructs (e.g. arbitrary callables, `Any`) raise `TypeError` from
`pydantic_to_typebuilder` rather than silently coercing to `string`.

## Testing Strategy

- **Framework**: Python `pytest` (markers `unit`/`functional`/`integration`,
  `pytest-asyncio`); Go `go test` (`*_integration_test.go` for live calls).
- **Deterministic tests, always run in CI**: codegen (B-A1), parser-only against
  **recorded fixtures** (B-C1, B-D2), drift gate (B-F1). Oracles: fixture content,
  empty `git diff`, `.baml` source names.
- **Real-LLM tests, API-key-gated**: B-B1, B-B2, B-D1 make a **real** model call,
  skipped when the key env var is absent but a genuine call when present, and
  assert **known facts** in a fixed input (the independent oracle). A spy/mock is
  never the passing oracle (Law 1).
- **Property tests** where an input domain exists (e.g. parser round-trips) via
  Hypothesis (Python) / table inputs (Go).

---

## Behavior B-A1: Generate Python BAML client

### Test Specification
**Given**: `baml_src/` containing `generators.baml` (`output_type "python/pydantic"`)
and at least one `function`.
**When**: `baml-cli generate` runs.
**Then**: a real `baml_client/` package exists; `from baml_client import b`
imports without error and the declared function is an attribute of `b`.

**Precondition (Resolved Decisions §7 — verified, not assumed):** before any
codegen, `baml-py==0.222.0` is added to `sdk/python/pyproject.toml` and installed
via `pip install -e .[dev]` **into `sdk/python/.venv` (py3.12)**. A precondition
assertion confirms, *in that venv*: `import baml_py` succeeds (abi3 wheel loads on
3.12), `sdk/python/.venv/bin/baml-cli --version` runs, and the API-Reference
symbols (`baml_py.errors.BamlValidationError`, `ClientRegistry`, `AbortController`,
`TypeBuilder` on the generated client) resolve **there** — not in the repo-root
py3.14 venv, which is outside `requires-python` and is not evidence.

**Edge Cases**: missing `generators.baml` → clear non-zero error; `.baml` syntax
error → non-zero exit with message; `baml-py` absent from `sdk/python/.venv` →
precondition fails loudly before `baml-cli generate`.

**Property (independent oracle — not self-referential):** the test asserts a
**hardcoded** expected-function list authored by hand, e.g.
`EXPECTED_FUNCTIONS = ["ExtractPerson", "ExtractDynamic"]`, and checks
`all(hasattr(b, f) for f in EXPECTED_FUNCTIONS)`. It does **not** enumerate
function names from `baml_src/` at test time (that would derive the oracle from
the artifact under generation and could never fail). Adding a new `.baml`
function therefore requires a test edit — intentional, so the test stays a real
gate. (Resolves review finding: B-A1 property oracle self-reference.)

**generators.baml syntax** (verified): `generator target { output_type "python/pydantic"  output_dir "../"  default_client_mode "async"  version "0.222.0" }`.

**Files touched**: `sdk/python/baml_src/generators.baml`,
`sdk/python/baml_src/clients.baml`, `sdk/python/baml_src/functions.baml`,
`sdk/python/pyproject.toml` (add `baml-py==0.222.0` as a **runtime** dependency —
the committed `baml_client` imports `baml_py` at import time; see Resolved
Decisions §1), `sdk/python/baml_client/**` (generated, **committed** — not
gitignored; Resolved Decisions §1), `sdk/python/tests/test_baml_codegen.py`.

### TDD Cycle
🔴 **Red** — `tests/test_baml_codegen.py`: assert `from baml_client.async_client import b`
and `hasattr(b, f)` for the hardcoded `EXPECTED_FUNCTIONS`. Fails (no client yet).
🟢 **Green** — *step 0:* add `baml-py==0.222.0` to `sdk/python/pyproject.toml`
deps, `pip install -e .[dev]` into `sdk/python/.venv` (py3.12), assert the
precondition above (import + `baml-cli --version` + symbols) **in that venv**.
*step 1:* author minimal `baml_src/`; run `sdk/python/.venv/bin/baml-cli generate`.
Commit `baml_client/`. Test passes.
🔵 **Refactor** — split `clients.baml`/`functions.baml`/`generators.baml`; confirm
`generator.version == "0.222.0"` (must equal installed `baml-py`).

> **Codegen determinism (drift-gate prerequisite, see B-F1):** `baml-cli generate`
> output carries no observed timestamp/random ordering in 0.222.0, so
> `git diff --exit-code` is a valid drift oracle. This is **verified empirically**
> in B-A1 green by running `baml-cli generate` twice and asserting an empty diff;
> if a future BAML version introduces nondeterministic output, the drift gate
> would false-positive and this assumption must be revisited.

---

## Behavior B-B1: Real structured call via a static BAML function

### Test Specification
**Given**: a static BAML extract function and `Person` class, and an LLM API key
in env. Concrete BAML definitions:
```baml
class Person {
  name string
  birth_year int
  field string  // area of work, e.g. "mathematics"
}
function ExtractPerson(bio: string) -> Person {
  client AnthropicClient
  prompt #"Extract the person from: {{ bio }} {{ ctx.output_format }}"#
}
```
**When**: called on a fixed factual paragraph about Ada Lovelace.
**Then**: returns a validated `Person` with the **correct known field values**.

**Edge Cases**: model wraps JSON in prose → BAML still parses; missing API key →
integration test **skipped**, not failed.

**clients.baml config schema** (verified syntax): each client is
`client<llm> AnthropicClient { provider anthropic  options { model "claude-..."  api_key env.ANTHROPIC_API_KEY  request_timeout_ms 120000 } }`; analogous
`OpenAIClient` (`provider openai`, `env.OPENAI_API_KEY`) and `OpenRouterClient`
(`provider openai-generic`, `env.OPENROUTER_API_KEY`). `request_timeout_ms` is
where the per-client timeout lives (Resolved Decisions §4).

**Person field oracle (independent):** `result.name == "Ada Lovelace"` AND
`result.birth_year == 1815` — two independently-known facts, not a spy.

**Files touched**: `sdk/python/baml_src/functions.baml`,
`sdk/python/baml_src/clients.baml`,
`sdk/python/tests/test_baml_inline_integration.py`.

### TDD Cycle
🔴 **Red** — integration test (marker `integration`) does
`from baml_client.async_client import b` then
`result = await b.ExtractPerson(<fixed Ada Lovelace bio>)`, asserts
`result.name == "Ada Lovelace"` and `result.birth_year == 1815`. Fails (function absent).
🟢 **Green** — define `ExtractPerson` + `Person` + `AnthropicClient` in `.baml`,
regenerate. Real async call returns the typed object; assertions pass.
🔵 **Refactor** — add `@description` hints; confirm provider via `clients.baml`.

---

## Behavior B-C1: Harness parser-only replaces `try_parse_from_text` ladder

### Test Specification
**Given**: a recorded fixture of messy third-party CLI stdout (JSON inside a
```` ```json ```` fence plus surrounding prose).
**When**: the Python harness parses it for a target schema **via BAML**.
**Then**: returns the correct typed object **and** the old `try_parse_from_text` /
repair ladder (`_schema.py:209`, `_runner.py` repair tiers) is removed.

**Edge Cases**: recoverable-malformed JSON → parsed; totally unparseable text →
**typed** parse error (not silent `None`); preserve `FailureType.SCHEMA`
classification on unrecoverable input.

**Replacement entry point (named):** keep `_schema.py:209`
`try_parse_from_text(text: str, schema: Any) -> Optional[Any]` — **same name and
signature** (lowest blast radius; `parse_and_validate` and the two `_runner.py`
call sites `:355` and `:471` call it unchanged — **not** `:348`, which is
`parse_and_validate` itself). Swap its body: instead of the 3 string-scrape
strategies (json fence / largest-brace / `cosmetic_repair`), call BAML
parser-only on the sync client and return a **schema-typed instance** via the
shared `deserialize` (symmetric with B-B2 — callers expect an instance of
`schema`, never a raw `DynamicOutput`):
```python
import baml_py
from baml_client.sync_client import b
from agentfield.baml_bridge import pydantic_to_typebuilder, deserialize
def try_parse_from_text(text: str, schema: Any) -> Optional[Any]:
    tb = pydantic_to_typebuilder(schema)
    try:
        raw = b.parse.ExtractDynamic(text, baml_options={"tb": tb})
        return deserialize(raw, schema)   # -> instance of `schema`, not DynamicOutput
    except (baml_py.errors.BamlValidationError, ValueError):
        return None   # parse_and_validate maps None → FailureType.SCHEMA
```
> Note: `deserialize` re-raises validation failures as `ValueError` (Resolved
> Decisions §6), so the `except` catches **both** `BamlValidationError` (raised by
> `b.parse`) and that `ValueError`, collapsing either to `None` →
> `FailureType.SCHEMA`. (Resolves review I2: B-C1's return type now matches B-B2.)

> **`baml_bridge.py` is introduced here (B-C1)** — it is the first consumer of
> both `pydantic_to_typebuilder` and `deserialize`. B-B2 reuses the same module
> and adds the dedicated `test_baml_bridge.py` mapper unit test. B-C1's fixture is
> a flat `Person` (primitives only), so B-C1 only exercises the primitive rows of
> the mapping table; B-B2 extends coverage to nested/optional/union/enum.
>
> **Ordering note:** the `@@dynamic DynamicOutput` class + `ExtractDynamic`
> function (shown in B-B2) are authored in **B-A1** — B-A1's `EXPECTED_FUNCTIONS`
> already lists `ExtractDynamic`, so the generated client exposes
> `b.parse.ExtractDynamic` by the time B-C1 (#3, before B-B2 #4) runs. B-B2 wires
> the bridge + reroutes `agent.ai`; it does not first-author the BAML function.
**Exception → FailureType mapping (named):** BAML raises
`baml_py.errors.BamlValidationError` on unrecoverable parse. `parse_and_validate`
(`_schema.py:196,203`) already catches bare `Exception` and returns `None`, and
`_runner._handle_schema_with_retry` (`_runner.py:332`) maps a `None` parse to
`FailureType.SCHEMA` (`_runner.py:499`). So **no new catch clause is needed at
`_runner.py:332`** — the existing bare-`Exception`/`None` path absorbs
`BamlValidationError` and preserves `FailureType.SCHEMA`. (Resolves review
finding: "_runner.py:332 silently stops catching" — verified false against the
real code: the catch is `except Exception`, which includes `BamlValidationError`.)
The full enum is unchanged: `FailureType` (in `harness/_result.py`) = `NONE,
CRASH, TIMEOUT, API_ERROR, NO_OUTPUT, SCHEMA`; only `SCHEMA` is in play here.

**Retry behavior after parse failure (named):** a BAML parse failure on **fixed**
captured text is deterministic — re-parsing identical text is a no-op. So the
existing `_runner` retry orchestration (`DEFAULT_SCHEMA_RETRIES=2`, `_runner.py:46`)
**re-invokes the provider CLI for fresh output** and re-parses that, exactly as
today; it does NOT loop on identical text. The `_ai_schema_repair` LLM-repair tier
(`_runner.py:59`) is **deleted** (BAML's SAP subsumes cosmetic repair).

**Fixture schema (`messy_cli_output.txt`):** the embedded JSON is, verbatim,
` ```json\n{"name": "Ada Lovelace", "birth_year": 1815, "field": "mathematics"}\n``` `
wrapped in surrounding prose ("Here is the result:\n…\nLet me know if you need more.").
The independent oracle is the hand-authored
`EXPECTED = {"name": "Ada Lovelace", "birth_year": 1815, "field": "mathematics"}`.

**Property**: for any fixture whose embedded JSON validates against the schema,
the parsed object equals the embedded object (`parse(fixture) == EXPECTED`).

**Files touched**: `sdk/python/agentfield/baml_bridge.py` (**new** — introduced
here with `pydantic_to_typebuilder` + `deserialize`; reused by B-B2),
`sdk/python/agentfield/harness/_schema.py` (swap `try_parse_from_text` body to
BAML parser-only; delete `cosmetic_repair` if now unused),
`sdk/python/agentfield/harness/_runner.py` (delete `_ai_schema_repair` tier at
`:59`; keep `_handle_schema_with_retry`/`FailureType` orchestration; the call
sites at `:355`/`:471` are unchanged), `sdk/python/tests/test_harness_schema.py`,
`sdk/python/tests/fixtures/messy_cli_output.txt` (new).

> Blast-radius note: this is two files of real logic + one fixture. Keep the
> `FailureType`/retry orchestration in `_runner.py`; only the *parsing* swaps to
> BAML. If the diff starts touching every provider, stop and reconsider.

### TDD Cycle
🔴 **Red** — deterministic test (no marker; runs in CI) feeds the recorded fixture
to `try_parse_from_text`, asserts `== EXPECTED`. Fails / still uses old ladder.
🟢 **Green** — replace `try_parse_from_text` internals with
`b.parse.ExtractDynamic(text, baml_options={"tb": ...})`. Passes.
🔵 **Refactor** — delete dead regex/`cosmetic_repair`/`_ai_schema_repair` code;
add a test asserting unparseable input (e.g. `"not json at all"`) → `None` →
`FailureType.SCHEMA` preserved.

---

## Behavior B-B2: Reroute `agent.ai(schema=)` through BAML (core replacement)

### Test Specification
**Given**: an agent calling `agent.ai(prompt, schema=RuntimePydanticModel)`.
**When**: it executes against a real model.
**Then**: returns a validated instance produced **via BAML** (generic `@@dynamic`
function + TypeBuilder maps the runtime Pydantic model), and the old
`json.loads`+regex parse **ladder** (`agent_ai.py:788-810`) is **deleted**. The
`:814-827` parse-retry loop and `:729 max_parse_retries` are **KEPT** — they key
on the `ValueError("Could not parse structured response: ...")` substring, which
`deserialize` still raises, so retry-on-parse-failure behavior is unchanged
(Resolved Decisions §6).

**Edge Cases**: nested / optional / `list` Pydantic fields map through TypeBuilder
(per the mapping table above); model omits a required field → `BamlValidationError`
(or `pydantic.ValidationError`) caught **inside `deserialize`** and re-raised as
`ValueError("Could not parse structured response: ...")` — the **same public type
and message** callers see today (not a silent `None`, not a leaked
`BamlValidationError`); enums/literals map to BAML dynamic enums / literal unions.

**The `@@dynamic` function + bridge interfaces (named, the load-bearing pieces):**
```baml
// baml_src/functions.baml
class DynamicOutput {
  @@dynamic          // no static fields; all added at call time via TypeBuilder
}
function ExtractDynamic(prompt: string) -> DynamicOutput {
  client AnthropicClient
  prompt #"{{ prompt }} {{ ctx.output_format }}"#
}
```
```python
# sdk/python/agentfield/baml_bridge.py   (internal module: agentfield.baml_bridge)
from baml_client.type_builder import TypeBuilder
from pydantic import BaseModel
def pydantic_to_typebuilder(model: type[BaseModel]) -> TypeBuilder:
    """Map a runtime Pydantic model onto the @@dynamic root class.
    For each top-level field of `model`, attach it to the dynamic root via
    `tb.DynamicOutput.add_property(name, <map(field_type)>)` (see mapping table);
    nested BaseModels become their own `tb.add_class(...)` registrations.
    Raises TypeError for unsupported field constructs (Any/callables)."""
    ...
def deserialize(raw: object, model: type[BaseModel]) -> BaseModel:
    """Cast BAML's dynamic output into an instance of the caller's model.
    `raw` is BAML's DynamicOutput. Returns `model.model_validate(raw.model_dump())`.
    PRESERVES the public exception contract (Resolved Decisions §6): catches
    pydantic.ValidationError AND baml_py.errors.BamlValidationError and re-raises
    ValueError("Could not parse structured response: ...") so existing callers
    (and the :814-827 retry loop) keep working unchanged."""
    import baml_py
    from pydantic import ValidationError
    try:
        return model.model_validate(raw.model_dump())
    except (ValidationError, baml_py.errors.BamlValidationError) as e:
        raise ValueError(f"Could not parse structured response: {e}") from e
```
This closes the review's "TypeBuilder → Pydantic instance cast unspecified" /
"callers may silently receive a raw dict" gap: callers always get
`model`-validated instances via `deserialize`, never a raw dict.

**agent_ai.py call expression (concrete):** replace the `:546-557` `response_format`
attach and the `:788-810` `json.loads`+regex parse **ladder** (the body of the
`if schema:` branch inside `_execute_and_parse`) with the BAML call below. The
surrounding `:814-827` retry loop and `:729 max_parse_retries` stay in place and
keep wrapping `_execute_and_parse` — because `deserialize` re-raises the same
`ValueError("Could not parse structured response: ...")`, the loop still retries
parse failures by re-running the model:
```python
from baml_client.async_client import b
from agentfield.baml_bridge import pydantic_to_typebuilder, deserialize
tb = pydantic_to_typebuilder(schema)
raw = await b.ExtractDynamic(prompt, baml_options={"tb": tb})
return deserialize(raw, schema)   # raises ValueError on validation failure → retry loop catches it
```

**Timeout forwarding (Resolved Decisions §4):** keep the existing
`asyncio.wait_for(..., timeout=effective_timeout * 2)` wall-clock guard
(`agent_ai.py:673`) wrapping the `await b.ExtractDynamic(...)` call, and set
`request_timeout_ms` in the `clients.baml` client block. The caller's `timeout=`
(via `effective_timeout`, `agent_ai.py:539-544`) therefore still bounds the call —
it is **not** silently dropped.

**Rate-limit wrapping (concrete):** the BAML call goes through the same
`StatelessRateLimiter.execute_with_retry(func)` seam used today
(`agent_ai.py:733-741`): wrap the `await b.ExtractDynamic(...)` (plus `deserialize`)
in the `_execute_with_fallbacks` closure that `execute_with_retry` already
receives. BAML 429s surface as `baml_py.errors.BamlClientHttpError`; add its
status-429 case to the rate-limiter's existing rate-limit-error detector
(`rate_limiter.py:55-98`, which matches on `"RateLimitError"`/429) so retry/backoff
still triggers.

**Property**: for a Pydantic model `M` and a hand-authored JSON instance `J`,
`deserialize(b.parse.ExtractDynamic(json.dumps(J), baml_options={"tb":
pydantic_to_typebuilder(M)}), M) == M(**J)` — using `b.parse` (no LLM call) so the
property is deterministic and CI-runnable.

**Files touched**: `sdk/python/agentfield/agent_ai.py` (replace `:546-557` attach
+ `:788-810` parse ladder with the BAML dynamic call; **keep** the `:814-827`
retry loop and `:729 max_parse_retries`),
`sdk/python/agentfield/baml_bridge.py` (reused — **introduced in B-C1**; B-B2 adds
its dedicated mapper unit test `test_baml_bridge.py`; internal module imported by
full path, no `__all__` needed),
`sdk/python/agentfield/rate_limiter.py` (add `BamlClientHttpError`/429 to the
detector), `sdk/python/baml_src/functions.baml` (`DynamicOutput` + `ExtractDynamic`
— **authored in B-A1**; only refined here if needed),
`sdk/python/tests/test_agent_ai_baml_integration.py`,
`sdk/python/tests/test_baml_bridge.py` (unit, deterministic, for the mapper).

> Blast-radius note: the Pydantic→TypeBuilder mapping is the riskiest unit — it
> gets its **own deterministic unit test** (`test_baml_bridge.py`) independent of
> any LLM call, so the mapping is verified without network.

### TDD Cycle
🔴 **Red** —
(1) unit `test_baml_bridge.py` (deterministic, no LLM): for a nested model `M`,
assert `deserialize(b.parse.ExtractDynamic(json.dumps(HARDCODED), baml_options=
{"tb": pydantic_to_typebuilder(M)}), M) == M(**HARDCODED)`, where `HARDCODED` and
the expected `M(**HARDCODED)` are **hand-authored independent oracles** — the test
does NOT inspect the TypeBuilder's own output shape. (Resolves review findings:
"oracle inspects the code's own output" / "self-referential verification".)
(2) integration: `agent.ai(<Ada bio>, schema=Person)` returns
`Person(name="Ada Lovelace", birth_year=1815, ...)`. Both fail initially.
🟢 **Green** — implement `pydantic_to_typebuilder` + `deserialize`; rewrite the
`agent_ai.py` structured branch to the BAML dynamic call above. Tests pass.
🔵 **Refactor** — delete only the `:788-810` `json.loads`+regex parse ladder;
**keep** the `:814-827` `ValueError`-keyed retry loop and `:729 max_parse_retries`
(they still function — `deserialize` preserves the `ValueError` contract). Confirm
rate-limit retry still wraps the call and the `asyncio.wait_for` guard remains.

---

## Behavior B-F1: Build integration & no-drift gate

### Test Specification
**Given**: a changed `baml_src/`.
**When**: `make build` (or a `generate` target) runs.
**Then**: `baml_client` regenerates and a CI drift check (`baml-cli generate`
then `git diff --exit-code`) passes.

**Edge Cases**: stale committed client → drift check fails loudly; `baml-cli` not
found → build errors with an install hint. Note: `baml-cli` is installed **by the
`baml-py` wheel** into the SDK venv, so the Makefile pins
`$(VENV) := sdk/python/.venv` and invokes `$(VENV)/bin/baml-cli generate`
explicitly (Resolved Decisions §7) — **not** the repo-root `.venv` and not a
system binary. The drift gate depends on `baml-cli generate` being deterministic
(no timestamps) — asserted empirically in B-A1; if it regresses the gate
false-positives.

**Files touched**: `Makefile` (pin `$(VENV) := sdk/python/.venv`; add
`baml-generate` target invoking `$(VENV)/bin/baml-cli generate`; hook into
`sdk-python`/`build`),
`.github/workflows/*.yml` (drift check step), `sdk/python/pyproject.toml`
(optional `on_generate` formatter).

### TDD Cycle
🔴 **Red** — CI/script: modify a `.baml`, run generate, `git diff --exit-code`
without committing → fails (drift). Encodes the gate.
🟢 **Green** — add `make baml-generate` (venv `baml-cli`); regenerate + commit;
gate green. Also assert generate-twice → empty diff (determinism check).
🔵 **Refactor** — make `make build` depend on `baml-generate`; document in
`CLAUDE.md`/`docs`.

---

## Behavior B-D1: Go BAML client + real structured call replaces reflection path

### Test Specification
**Given**: the BAML go generator configured (`output_type "go"`) and an API key
in env.
**When**: a Go BAML structured function runs against a real model.
**Then**: returns the correct struct values **and** the old `request.go`
`structToJSONSchema`→`Into` path (`request.go:390`, `response.go:99`) is replaced.

**Edge Cases**: nested/optional struct fields; missing key → integration test
skipped; `omitempty` fields handled by BAML, not the old required-array logic.

**BAML Go is pre-stable (caveat, see API Reference):** Go support is `v0.x`,
runtime is a C-FFI Rust shared lib, **Linux/macOS amd64-arm64 only (no Windows)**.
agentfield CI/dev is Linux so this is acceptable; document the platform limit in
`sdk/go` README. Runtime dep `github.com/boundaryml/baml v0.222.0`.

**Go toolchain gate (firm green step, resolves review API item).** BAML's module
declares `go 1.24.0`. A Go module that requires `go 1.24.0` **cannot be built by a
toolchain older than 1.24** — so `sdk/go/go.mod`'s current `go 1.21` is **not**
sufficient and `go build` will fail until bumped (the "≥1.18 for generics" figure
is the generated *code's* floor, not the binding constraint — BAML's `go 1.24.0`
directive dominates). Therefore B-D1's green step **bumps `sdk/go/go.mod` to `go
1.24` (or higher)** and **confirms the CI Go toolchain is ≥ 1.24** before the
build can pass. This is verified by an actual `go build ./...` against the
imported runtime, not assumed — it is a hard green gate, not an open item. (See
parallel precondition C1/§7 on the Python side.)

**Go baml_src definitions (concrete):**
```baml
// sdk/go/baml_src/generators.baml
generator gotarget {
  output_type "go"
  output_dir "../"
  version "0.222.0"
  client_package_name "github.com/Agent-Field/agentfield/sdk/go/baml_client"
}
// functions.baml — same Person/ExtractPerson as Python, Go-generated
class Person { name string  birth_year int  field string }
function ExtractPerson(bio: string) -> Person { client AnthropicClient  prompt #"...{{ ctx.output_format }}"# }
```

**Consumer-facing API after migration (Resolved Decisions §2 — no open OR):**
the new first-class API is the generated package-level
`result, err := b.ExtractPerson(ctx, bio)`. `request.go:258 WithSchema(schema
interface{}) Option` and `response.go:100 JSON(dest)` / `:109 Into(dest)` are
**kept as thin deprecated wrappers** that delegate to the BAML path (zero existing
Go callers break), each marked `// Deprecated: call the generated baml_client
function directly`. The reflection internals
`structToJSONSchema` (`request.go:390`) / `goTypeToJSONType` (`request.go:456` —
`:455` is its doc comment) and the `Strict:true` ResponseFormat assembly
(`request.go:286-291`) are deleted.

**Goroutine-safety smoke test:** since BAML Go's concurrency safety is not
formally documented (CFFI deadlock fixed v0.220.0), B-D1 adds a test firing N=8
concurrent `b.ExtractPerson` goroutines (key-gated) asserting no panic/deadlock.

**Files touched**: `sdk/go/baml_src/{generators,clients,functions}.baml`,
`sdk/go/baml_client/**` (generated, committed), `sdk/go/go.mod`/`go.sum`
(`github.com/boundaryml/baml v0.222.0`; **bump `go` directive 1.21 → 1.24** to
satisfy BAML's `go 1.24.0` requirement),
`sdk/go/ai/request.go` (delete reflection; `WithSchema` → deprecated delegating
wrapper), `sdk/go/ai/response.go` (`JSON`/`Into` → deprecated wrappers delegating
to `b.Parse`), `sdk/go/ai/baml_structured_integration_test.go` (new, incl.
concurrency smoke test).

### TDD Cycle
🔴 **Red** — `*_integration_test.go` calls `b.ExtractPerson(ctx, <fixed bio>)` on a
fixed input, asserts `result.Name == "Ada Lovelace"`. Fails (no client).
🟢 **Green** — author `sdk/go/baml_src/`, `go get github.com/boundaryml/baml@v0.222.0`,
**bump `sdk/go/go.mod` to `go 1.24`** and confirm the CI Go toolchain is ≥ 1.24,
run `baml-cli generate`, commit `baml_client/`; confirm `go build ./...` passes
against the runtime (the 1.21 toolchain would fail BAML's `go 1.24.0` directive).
Repoint `WithSchema`/`JSON`/`Into` wrappers to the BAML path. Passes.
🔵 **Refactor** — delete `structToJSONSchema`/`goTypeToJSONType` reflection +
`Strict:true` ResponseFormat assembly now owned by BAML.

---

## Behavior B-D2: Go harness parser-only replaces `schema.go` post-hoc path

### Test Specification
**Given**: a recorded messy-CLI fixture.
**When**: the Go harness parses it via BAML.
**Then**: returns the correct struct **and** the old `sdk/go/harness/schema.go`
post-hoc parse path is replaced.

**Edge Cases**: unparseable fixture → typed error (BAML Go parse error), mapped to
the existing Go `FailureSchema` (`runner.go:410`, not `:401` — `:401` is an
`accumulateMetrics` call). **Large-schema-on-disk gate:
DROPPED (Resolved Decisions §3)** — delete `schemaFilename` (`schema.go:15`),
`largeSchemaTokenThreshold = 4000` (`schema.go:16`), and the >4000-token
write-to-disk + prompt-suffix branch (`schema.go:50-62`). BAML embeds its own
output schema in the prompt, so the manual `.agentfield_schema.json` path is
obsolete. This is a code deletion, not a runtime flag — no open item remains.

**Go parser-only entry (named):** `b.Parse.ExtractDynamic(text, ...)` (the Go
`Parse` global, added BAML v0.204.0) — the LLM-free parse path, Go equivalent of
Python's `b.parse.Fn`. For runtime schemas use `b.NewTypeBuilder()` +
`b.WithTypeBuilder(tb)`.

**Fixture schema (`testdata/messy_cli_output.txt`):** identical embedded JSON to
the Python fixture — ` ```json\n{"name":"Ada Lovelace","birth_year":1815,"field":"mathematics"}\n``` `
in prose; the independent oracle is the hand-authored expected struct
`Person{Name:"Ada Lovelace", BirthYear:1815, Field:"mathematics"}`.

**Property**: for any fixture whose embedded struct JSON validates, parsed struct
equals the hand-authored expected struct.

**Files touched**: `sdk/go/harness/schema.go` (swap parse internals to
`b.Parse.*`; delete the large-schema gate + its two consts), `sdk/go/harness/runner.go`
(keep `FailureType`/retry orchestration; map BAML parse error → `FailureSchema`),
`sdk/go/harness/schema_test.go`,
`sdk/go/harness/testdata/messy_cli_output.txt` (new).

### TDD Cycle
🔴 **Red** — deterministic `go test` feeds the fixture to the harness parse entry,
asserts `== Person{...}`. Fails.
🟢 **Green** — replace parse internals with `b.Parse.*` BAML parser-only. Passes.
🔵 **Refactor** — delete the dead reflection/large-schema-gate code; confirm
unparseable input still yields `FailureSchema`.

---

## Implementation Order (vertical slices)

1. **B-A1** — Python BAML codegen foundation → *visible:* importable `baml_client`.
2. **B-B1** — real static structured call → *visible:* correct typed object printed.
3. **B-C1** — Python harness parser-only → *visible:* messy fixture parsed; ladder gone.
4. **B-B2** — reroute `agent.ai(schema=)` via TypeBuilder → *visible:* real `.ai()` typed object; old ladder deleted.
5. **B-F1** — build/CI drift gate → *visible:* `make build` regenerates; gate green.
6. **B-D1** — Go client + real call → *visible:* `go test` typed struct; reflection replaced.
7. **B-D2** — Go harness parser-only → *visible:* fixture parsed; `schema.go` path replaced.

## Risks & Open Items

All previously-deferred decisions are now **resolved** (see Resolved Decisions):
commit-vs-generate (§1), Go `WithSchema`/`JSON` fate (§2), large-schema-on-disk
(§3 — dropped), timeout (§4), sync/async (§5), exception/parse-retry contract
(§6), build/test venv + py-version (§7). The four findings from the 2026-06-08
plan review are folded in: **C1** (env/venv — §7, B-A1 precondition), **I1**
(exception contract — §6, B-B2 `deserialize`), **I2** (B-C1 returns a
`schema`-typed instance via `deserialize`), **I3** (line drifts corrected:
`request.go:456`, `runner.go:410`, `_runner.py:355/471`), plus the **Go toolchain
gate** (go.mod 1.21 → 1.24, B-D1). Remaining genuine risks:

- **TypeBuilder coverage** for nested/optional/union/enum/recursive Pydantic
  fields — de-risked by the standalone deterministic `test_baml_bridge.py` (B-B2)
  using `b.parse` (no network) and the mapping table as the independent oracle.
- **BAML Go pre-stability** — `v0.x`, C-FFI runtime, Linux/macOS-only, no formal
  goroutine-safety guarantee. Mitigated by the B-D1 concurrency smoke test and the
  documented platform limit. A future BAML API break is possible (pinned `0.222.0`).
- **Codegen determinism** — drift gate (B-F1) assumes `baml-cli generate` is
  deterministic; asserted empirically in B-A1 (generate-twice → empty diff). If a
  BAML release breaks this, the gate false-positives and must be revisited.
- **Per-call timeout** — BAML has no per-call timeout option (issue #1630); honored
  via the `asyncio.wait_for` wall-clock guard + `clients.baml` `request_timeout_ms`
  (Resolved Decisions §4), not via BAML call options.
- **Real-LLM test cost/flakiness** — gated by API key; assert on stably-extractable
  facts only (`name == "Ada Lovelace"`, `birth_year == 1815`).

> **Provenance:** every BAML symbol in this plan was verified against
> `baml-py==0.222.0` and BoundaryML docs/generated-client source on 2026-06-07;
> every codebase `file:line` was verified against the working tree (drifted cites
> corrected per the 2026-06-08 review — see I3). **Caveat (Resolved Decisions §7):**
> the original symbol check ran in the repo-root `.venv` (py3.14); B-A1
> **re-verifies in `sdk/python/.venv` (py3.12)** — the SDK's real build/test env —
> before relying on any symbol. See the review-response notes and the
> `…-REVIEW.md` alongside this plan.

## References

- Research: `thoughts/searchable/shared/research/2026-06-07-08-32-baml-structured-output-integration-surface.md`
- Plan review: `2026-06-07-09-04-tdd-baml-structured-output/2026-06-07-09-04-tdd-baml-structured-output-REVIEW.md`
- BAML docs: https://docs.boundaryml.com (TypeBuilder, generator, ClientRegistry refs); `baml-py==0.222.0` (PyPI), `github.com/boundaryml/baml v0.222.0` (Go runtime)
- `sdk/python/agentfield/agent_ai.py:257,539-557,673,729,733-741`; parse ladder `:788-810` (deleted), `ValueError` raise `:808`, retry loop `:814-827` (kept — §6)
- `sdk/python/agentfield/harness/_schema.py:209`, `_runner.py:46,59,332,496` (`try_parse_from_text` call sites at `:355`,`:471` — `:348` is `parse_and_validate`), `_result.py` (`FailureType`)
- `sdk/python/agentfield/rate_limiter.py:55-98,209` (`StatelessRateLimiter.execute_with_retry`)
- `sdk/go/ai/request.go:258,286-291,390,456` (`goTypeToJSONType` at `:456`; `:455` is its doc comment), `sdk/go/ai/response.go:100,109`, `sdk/go/harness/schema.go:15,16,50`, `sdk/go/harness/runner.go:410` (`FailureSchema`; not `:401`), `sdk/go/go.mod` (`go 1.21` → bump to `1.24`)
- `sdk/python/pyproject.toml:31-32,77-81`
