---
date: 2026-06-07T08:32:48-04:00
researcher: maceo
git_commit: fb8559d5d88d42e8b15043ff746611ee482152ba
branch: main
repository: agentfield
topic: "Current structured-output, schema-validation, and prompt-construction surfaces (BAML integration target)"
tags: [research, codebase, structured-output, schema-validation, zod, prompt-construction, python-sdk, go-sdk, web-ui, harness]
status: complete
last_updated: 2026-06-07
last_updated_by: maceo
---

# Research: Current structured-output / schema-validation / prompt-construction surfaces

**Date**: 2026-06-07 08:32:48 EDT
**Researcher**: maceo
**Git Commit**: fb8559d5d88d42e8b15043ff746611ee482152ba
**Branch**: main
**Repository**: agentfield

## Research Question

We want to add BAML to enhance structured output generation, reduce zod errors,
reduce parse errors, and give the code "prompt as function." This research
documents — *as the code exists today* — every place the codebase generates
structured LLM output, builds prompts, parses/validates model responses, and
uses zod. It is a map of the surface BAML would touch, not a proposal for how to
integrate it.

> Documentarian note: this report describes what IS. It does not evaluate the
> current design or recommend changes. Integration questions are collected in
> **Open Questions**.

## Summary

There are **four** distinct structured-output / schema-validation surfaces in
the repo. The example cited in the research request (`try_parse_from_text`,
`DEFAULT_SCHEMA_RETRIES`, `FailureType`, `opencode.py`) is **surface #2 — the
Python SDK harness** — confirmed present in this repo.

| # | Surface | Language | Mechanism today | Where |
|---|---------|----------|-----------------|-------|
| 1 | Inline AI calls (`agent.ai(... schema=)`) | Python | Pydantic model → `model_json_schema()` → LiteLLM `response_format: json_schema` → `json.loads` + `schema(**data)`, 2 parse retries | `sdk/python/agentfield/agent_ai.py` |
| 2 | Agent-runner **harness** (CLI providers: opencode/codex/claude/gemini) | Python | **Post-hoc text parsing**, not constrained decoding: agent writes JSON to a file or stdout → cosmetic repair → Pydantic validate → AI repair call → harness retry | `sdk/python/agentfield/harness/` |
| 3 | Go SDK AI + Go harness | Go | Struct → JSON schema via reflection → OpenAI/OpenRouter `response_format` → `json.Unmarshal`; Go harness mirrors #2 (schema file on disk) | `sdk/go/ai/`, `sdk/go/harness/` |
| 4 | Web UI form validation | TS/React | **zod** built at runtime from JSON Schema via `jsonSchemaToZod`, used by AutoForm + `validateFormData` | `control-plane/web/client/src/utils/` |

Key facts:
- **zod is used only in the web UI**, and only for **form input validation** —
  not for validating LLM/API responses. UI API/SSE/NDJSON responses are parsed
  with bare `JSON.parse` + type guards / `as` casts (no zod). (Surface #4.)
- **The control plane (Go) does not call LLMs** and has no `ai` package usage;
  it stores `InputSchema`/`OutputSchema` as `map[string]interface{}` metadata.
- Structured output in the **harness** (#2/#3) is explicitly post-hoc text
  parsing with a multi-tier repair/retry ladder, classified by a `FailureType`
  enum.

## Detailed Findings

### Surface #1 — Python inline AI (`agent_ai.py`)

The primary "call an LLM, get a typed object" path.

- **Request, schema attach** — `sdk/python/agentfield/agent_ai.py:546-557`: when a
  Pydantic `schema` is passed, it is converted and attached as a strict JSON
  schema response format (verified):
  ```python
  if schema:
      litellm_params["response_format"] = {
          "type": "json_schema",
          "json_schema": {
              "schema": schema.model_json_schema(),
              "name": schema.__name__,
              "strict": True,
          },
      }
  ```
- **Prompt / schema instruction assembly** — `agent_ai.py:402-427`: a natural-language
  schema-adherence instruction block is prepended to the system prompt when a
  schema is present ("You must exactly adhere to the output schema…"). Schema
  also rendered via `schema.model_json_schema()`.
- **`ai()` entry signature** — `agent_ai.py:257-282`, with `response_format`
  param documented at `agent_ai.py:267,304` and applied at `agent_ai.py:389-391`.
- **Parse + validate** — `agent_ai.py:788-810`: `json.loads(...)` then
  `schema(**json_data)`; catches `JSONDecodeError`/`ValueError`/`ValidationError`;
  a regex fallback (`r"\{.*\}"`) attempts to extract a JSON object from prose.
- **Parse retry loop** — `agent_ai.py:729` `max_parse_retries = 2` (verified),
  loop at `agent_ai.py:816` (`for attempt in range(max_parse_retries + 1)`),
  retry gated on the string `"Could not parse structured response"`
  (`agent_ai.py:809,820`).
- **LiteLLM param assembly / provider patches** — `sdk/python/agentfield/types.py:634-682`
  (`AIConfig.get_litellm_params()`); OpenAI rename `max_tokens`→`max_completion_tokens`.
  Provider patching also in `sdk/python/agentfield/litellm_adapters.py`
  (`get_provider_from_model:19`, `apply_openai_patches:46`, `apply_provider_patches:78`).
- **Tool-calling branch** — also parses structured responses
  (`agent_ai.py:653-660`), rate-limited via `rate_limiter.py`
  (`execute_with_retry` at `sdk/python/agentfield/rate_limiter.py:209-280`).
- **Pydantic helpers** — `sdk/python/agentfield/pydantic_utils.py:62-96`
  (`convert_dict_to_model`), auto-conversion of function args at `:150-154`.
- **Schema-from-typehints** — `sdk/python/agentfield/agent_schema.py:3-65`
  (`_type_to_json_schema`, calls `model_json_schema()` for Pydantic models).

### Surface #2 — Python SDK harness (the cited example)

Located at `sdk/python/agentfield/harness/`. Files:
`_runner.py`, `_schema.py`, `_result.py`, `_cli.py`, and
`providers/` (`claude.py`, `codex.py`, `gemini.py`, `opencode.py`, `_base.py`,
`_factory.py`). This is the surface the request's example quotes.

- **`FailureType` enum** — `harness/_result.py:8-25` (verified): values
  `NONE, CRASH, TIMEOUT, API_ERROR, NO_OUTPUT, SCHEMA`. `SCHEMA` = "output file
  exists but fails schema validation." Carried on `RawResult.failure_type`
  (`_result.py:45`) and `HarnessResult` (`_result.py:49-55`, with `parsed: Any`).
- **`DEFAULT_SCHEMA_RETRIES = 2`** — `harness/_runner.py:46` (verified);
  `_SCHEMA_REPAIR_TIMEOUT_SECONDS = 90` at `:51`.
- **Post-hoc parsing strategies** — `harness/_schema.py`:
  - `parse_and_validate(file_path, schema)` — `:185` (file → cosmetic repair → validate).
  - `try_parse_from_text(text, schema)` — `:209` (verified): 3 strategies —
    (1) ```` ```json ```` fences, (2) largest top-level `{ ... }` block,
    (3) cosmetic repair of whole text.
  - `read_repair_and_parse` `:155`, `validate_against_schema` `:168`.
- **Retry ladder** — `harness/_runner.py`:
  - `_ai_schema_repair(...)` `:59` — one tool-less LLM call to reformat malformed
    JSON, `response_format={"type":"json_object"}`, 90s timeout; falls through to
    `try_parse_from_text` on success.
  - `_handle_schema_with_retry(...)` `:332`; `schema_max_retries` read from
    options at `:342-343` (default `DEFAULT_SCHEMA_RETRIES`); retry loop `:414`.
  - Final failure message `"Schema validation failed after {n} retry attempt(s)…"`
    at `_runner.py:496`, returning `parsed=None`, `FailureType.SCHEMA`.
- **Final-text extraction across CLI tools** — `harness/_cli.py:68-110`
  (`extract_final_text`) handles `type:"result"`, Codex `item.completed`,
  OpenCode `type:"text"` part, last assistant message.
- **OpenCode provider** — `harness/providers/opencode.py` (~348 lines):
  invokes `opencode run` subprocess; `execute()` `:140`, `_execute_impl()` `:150`;
  uses `extract_final_text` (`:16,265`). Structured output = post-hoc parse of
  file/stdout, not constrained decoding.

### Surface #3 — Go SDK AI + Go harness

**Go SDK `sdk/go/ai/`** (mirror of #1 in Go):
- LLM calls — `sdk/go/ai/client.go:38` `Complete()`, `:66` `CompleteWithMessages()`,
  `:84-147` `doRequest()` (POST `/chat/completions`), `:151-256` `StreamComplete()`.
- Provider config — `sdk/go/ai/config.go:44-51` (OpenAI default / OpenRouter),
  `IsOpenRouter()` `:88-91`; default model `gpt-4o`.
- Prompt/message assembly — `sdk/go/ai/request.go:11-42` (Request),
  `:70-83` (Message/ContentPart), `WithSystem()` `:191-204`,
  custom `MarshalJSON` `:104-137`.
- **Schema from Go struct** — `sdk/go/ai/request.go:258` `WithSchema()` (verified;
  accepts struct / `json.RawMessage` / string / bytes); `structToJSONSchema()`
  `:390` (verified, reflection over `json:` tags + `description` tag);
  `goTypeToJSONType()` `:456`; `Strict: true` at `:290`.
- Response parse — `sdk/go/ai/response.go:99-111` (`JSON()`/`Into()` → `json.Unmarshal`),
  `StructuredAI()` convenience `sdk/go/ai/client.go:329-342`.
- Tool calling — `sdk/go/ai/tool_calling.go` (`ExecuteToolCallLoop` `:126`,
  arg parse `:197-200`, no retry on arg-parse failure — empty map fallback).
- Error typing — `sdk/go/ai/response.go:56-65` (`ErrorResponse`); otherwise plain
  `error` wrapping. **No LLM-call retry** in `doRequest`.

**Go harness `sdk/go/harness/`** (mirror of #2 in Go) — present and parallel to the
Python harness: `runner.go`, `schema.go`, `result.go`, `provider.go`, `cli.go`,
and providers `claudecode.go`, `codex.go`, `gemini.go`, `opencode.go`.
- `sdk/go/harness/schema.go`: `schemaFilename = ".agentfield_schema.json"` `:15`,
  `SchemaPath()` `:24-26`, `json.MarshalIndent` of schema `:38`,
  `largeSchemaTokenThreshold` gate `:50` (large schemas written to disk).
  `StructToJSONSchema`, `BuildPromptSuffix`, `BuildFollowupPrompt` referenced from
  tests (`sdk/go/harness/coverage_helpers_test.go:136,145,220`).

**Control plane (Go, `control-plane/internal/`)** — does **not** call LLMs; no `ai`
package usage. Stores schema metadata only: `InputSchema`/`OutputSchema` as
`map[string]interface{}` in `control-plane/internal/handlers/nodes_register.go:867-868`;
discovery flags `include_input_schema`/`include_output_schema` in
`control-plane/internal/handlers/discovery_test.go:71`. DB-level retry (not LLM) in
`control-plane/internal/handlers/retry.go:9-36`. No `control-plane/internal/mcp/`
directory exists.

### Surface #4 — Web UI zod usage (`control-plane/web/client/`)

zod is present at `package.json` (`"zod": "^4.1.12"`, `"@autoform/zod": "^5.0.0"`)
and used **only for form validation**.

- **JSON Schema → zod at runtime** — `control-plane/web/client/src/utils/jsonSchemaToZod.ts`
  (verified): `import { z, type ZodTypeAny } from "zod"` `:1`; `buildField()` `:14`
  maps types to `z.enum`/`z.number`/`z.boolean`/`z.array`/`z.string`/object with
  `.nullable()`/`.optional()`; `jsonSchemaToZodObject()` `:53` builds
  `z.object(shape).passthrough()`.
- **Validation + ZodError handling** — `control-plane/web/client/src/utils/schemaUtils.ts`:
  `import { ZodError } from "zod"` `:2`; `validateFormData()` `:307`;
  `zodSchema.safeParse(input)` `:319` (verified); `ZodError` formatting at
  `:335-344` (maps `.issues` to `path: message`); permissive fallback `:346`.
  Tests in `schemaUtils.test.ts:241-256`.
- **Form integration** — `control-plane/web/client/src/.../ExecutionForm.tsx`:
  `ZodProvider`/`jsonSchemaToZodObject` from `@autoform/zod` (`:9,19`), provider
  built `:321-332` with try/catch `console.warn` on conversion failure `:328-331`;
  raw-JSON editor `JSON.parse` with "Invalid JSON" message `:406`.
- **Non-zod parsing of API/stream data** (representative; all use `JSON.parse` +
  try/catch / type guards, **no zod**):
  - NDJSON logs — `src/services/api.ts:436,505,518` (silent skip).
  - VC documents — `src/services/vcApi.ts:45,384,495,544,694`.
  - SSE — `src/hooks/useSSE.ts:150`, `src/services/reasonersApi.ts:402`,
    `src/components/execution/ExecutionObservabilityPanel.tsx:142`.
  - Error/structured panels & dialogs —
    `src/components/execution/RedesignedErrorPanel.tsx:20`,
    `ExecutionApprovalPanel.tsx:69`, `triggers/NewTriggerDialog.tsx:190`,
    `nodes/NodeProcessLogsPanel.tsx:134`, `WorkflowDAG/index.tsx:1231`,
    `pages/PlaygroundPage.tsx:259`, `pages/VerifyProvenancePage.tsx:74`,
    API error bodies in `services/observabilityWebhookApi.ts:122`,
    `services/configurationApi.ts:53`.
- **API/structured types are TS interfaces, not zod** — `src/types/execution.ts`
  (`ExecutionRequest:1`, `ExecutionResponse:22`, `JsonSchema:87-112`,
  `FormField:114-134`), plus `types/executions.ts`, `types/workflows.ts`.
  No `z.infer` types found.

## Code References

- `sdk/python/agentfield/agent_ai.py:546-557` — Pydantic schema → LiteLLM `response_format: json_schema`.
- `sdk/python/agentfield/agent_ai.py:729,816` — inline parse retry (`max_parse_retries = 2`).
- `sdk/python/agentfield/harness/_result.py:8-25` — `FailureType` enum.
- `sdk/python/agentfield/harness/_runner.py:46,59,332,496` — schema retries const, AI repair, retry handler, failure message.
- `sdk/python/agentfield/harness/_schema.py:155,168,185,209` — repair/validate/parse functions.
- `sdk/python/agentfield/harness/_cli.py:68-110` — `extract_final_text`.
- `sdk/python/agentfield/harness/providers/opencode.py:140,150,265` — OpenCode provider.
- `sdk/python/agentfield/litellm_adapters.py:19,46,78` — provider patching.
- `sdk/go/ai/request.go:258,290,390,456` — `WithSchema`, strict, struct→schema, type map.
- `sdk/go/ai/client.go:38,84-147,329-342` — Complete / doRequest / StructuredAI.
- `sdk/go/ai/response.go:99-111` — `Into()`/`JSON()` unmarshal.
- `sdk/go/harness/schema.go:15,24,38,50` — on-disk schema file, large-schema gate.
- `control-plane/internal/handlers/nodes_register.go:867-868` — schema stored as `map[string]interface{}`.
- `control-plane/web/client/src/utils/jsonSchemaToZod.ts:1,14,53` — JSON Schema → zod.
- `control-plane/web/client/src/utils/schemaUtils.ts:2,307,319,335-344` — `safeParse`, ZodError formatting.
- `control-plane/web/client/package.json` — `zod ^4.1.12`, `@autoform/zod ^5.0.0`.

## Architecture Documentation

- **Schema source of truth flows JSON-Schema-shaped end to end.** Python emits
  `model_json_schema()`, Go emits reflection-based JSON schema, the control plane
  stores `Input/OutputSchema` as opaque JSON maps, and the UI converts that JSON
  Schema into zod at runtime (`jsonSchemaToZod`). The same JSON-Schema shape
  appears at all four surfaces.
- **Two structured-output strategies coexist.** Inline AI (#1/#3 SDK `ai`) uses
  provider-native `response_format` (constrained-ish) + a small parse retry. The
  harness (#2/#3 harness) deliberately uses **post-hoc text parsing** of an agent
  CLI's file/stdout output, with a 4-tier ladder: file parse → stdout extract →
  AI repair call → full agent re-run, classified by `FailureType`.
- **zod is isolated to UI form validation.** It is not on the LLM-response path
  and not on the API-response path.
- **Python and Go harnesses are mirror implementations** (same concepts:
  schema file on disk, prompt suffix, followup prompt, failure typing).

## Historical Context (from thoughts/)

- `thoughts/searchable/shared/research/2026-06-04-19-49-agentfield-ai-backend-claim.md`
  — prior research on the agentfield AI backend (does not cover structured
  output / parsing / zod specifically).
- `specs/sdk-python.md`, `specs/sdk-go.md`, `specs/architecture-overview.md`
  contain only brief table-row references pointing at `agent_ai.py`
  ("structured outputs"), `multimodal_response.py`, `litellm_adapters.py`,
  `agent_schema.py`, and Go `harness/schema.go` — no design/research docs on the
  topic.

## Related Research

- `thoughts/searchable/shared/research/2026-06-04-19-49-agentfield-ai-backend-claim.md`

## Open Questions

These are integration-scoping questions surfaced by the map (not answered here —
they would belong to a follow-up plan, not this documentarian research):

- BAML generates per-language clients (Python/Go/TS). The repo has structured
  output at all three. Which surface(s) are in scope — inline `ai` (#1/#3), the
  harness (#2/#3), the UI (#4), or all?
- The harness (#2) does post-hoc text parsing of *third-party agent CLIs*
  (opencode/codex/claude/gemini), where the model output is not under direct
  `response_format` control. How (or whether) BAML's parser applies there is
  unresolved.
- zod errors today come only from UI form validation built off JSON Schema. The
  relationship between "reduce zod errors" and a BAML-generated TS client would
  need definition.
- The control plane stores schemas as opaque JSON maps and does not call LLMs —
  whether BAML has any role server-side is open.
