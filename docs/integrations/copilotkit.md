# CopilotKit / AG-UI integration

AgentField speaks the [AG-UI protocol](https://docs.ag-ui.com) so any
AG-UI-compatible frontend — most notably [CopilotKit](https://docs.copilotkit.ai) —
can use AgentField as the agent backend with no custom adapter.

This page is the contract. If you're writing a reasoner that should
drive Generative UI or shared state in CopilotKit, the fields below are
how you opt in.

## Topology

```
Browser ──▶ <CopilotChat> / useCoAgent / useCopilotAction
            ──▶ CopilotRuntime (Next.js /api/copilotkit)
            ──▶ @ag-ui/client HttpAgent
            ──▶ POST /api/v1/agui/runs/<node_id>/<reasoner_name>
            ──▶ AgentField reasoner
```

CopilotKit posts a canonical `RunAgentInput` body. The control plane
forwards the same envelope to your reasoner and translates the response
into AG-UI Server-Sent Events.

## Endpoint

```
POST /api/v1/agui/runs/:node_id/:reasoner_name
Content-Type: application/json
```

Body shape (see `RunAgentInputSchema` in `@ag-ui/core`):

```json
{
  "threadId": "string",
  "runId": "string",
  "messages": [{ "role": "user|assistant|tool|system", "content": "...", "toolCalls": [...] }],
  "tools": [{ "name": "...", "description": "...", "parameters": { ... } }],
  "context": [{ "description": "...", "value": ... }],
  "state": { ... },
  "forwardedProps": { ... }
}
```

The control plane fans this into the reasoner input map under the same
keys, plus a `prompt` convenience extracted from the trailing user
message.

Response: an SSE stream of AG-UI events.

## Reasoner contract

Reasoners can return a flat result, or a structured map opting into any
of these AG-UI surfaces:

| Reasoner field | Emitted as | Used by |
|---|---|---|
| `result` (string or anything) | `TEXT_MESSAGE_CONTENT` | `<CopilotChat>` assistant bubble |
| `content` (alias for `result`) | `TEXT_MESSAGE_CONTENT` | same |
| `toolCalls: [{id, name, arguments, result?}]` | `TOOL_CALL_START` → `_ARGS` → `_END` (and `_RESULT` if `result` set) | `useCopilotAction({name, render})` |
| `state: {...}` | `STATE_SNAPSHOT` | `useCoAgent({state})` |
| `stateDelta: [...]` (RFC 6902 ops) | `STATE_DELTA` (after snapshot) | `useCoAgent({state})` |

If none of `result`/`content` is present, the control plane stringifies
the rest of the body (minus `toolCalls`/`state` internals) so you still
see something.

Long `result` values are auto-chunked across multiple
`TEXT_MESSAGE_CONTENT` deltas (default 256 chars each) so the frontend
can paint progressively even though the reasoner is synchronous. Each
delta carries the same `messageId`; concatenation reproduces the full
text.

### Python example

```python
from agentfield import Agent, agui

app = Agent(node_id="my-app")

@app.reasoner()
async def book_flight(prompt: str = "", state: dict | None = None):
    counter = (state or {}).get("counter", 0) + 1
    return {
        "result": "Pulling up flight options.",
        "toolCalls": [
            agui.tool_call(
                name="showFlightCard",
                arguments={"from": "SFO", "to": "JFK", "depart": "2026-06-01"},
                id="tc-flight-1",
            ),
        ],
        "state": {"counter": counter, "lastBooking": "AA-12"},
        "stateDelta": [
            agui.state_delta_replace("/counter", counter),
        ],
    }
```

If your reasoner uses `app.ai(tools=...)` and you want the LLM's
tool-calling trace to surface in the UI, hand the trace to
`agui.tool_calls_from_trace`:

```python
@app.reasoner()
async def smart_chat(prompt: str = ""):
    result = await app.ai(prompt, tools="discover")
    return {
        "result": result.text,
        "toolCalls": agui.tool_calls_from_trace(result.trace),
    }
```

Each entry in the trace becomes a TOOL_CALL_*/_RESULT triad — the UI
shows a completed-tool indicator instead of a perpetually-pending
placeholder.

### Go example

```go
import (
    "context"
    "github.com/Agent-Field/agentfield/sdk/go/agent"
    "github.com/Agent-Field/agentfield/sdk/go/agent/agui"
)

a, _ := agent.New(agent.Config{NodeID: "my-app"})
a.RegisterReasoner("book_flight", func(ctx context.Context, in map[string]any) (any, error) {
    return map[string]any{
        "result": "Pulling up flight options.",
        "toolCalls": []map[string]any{
            agui.ToolCall("tc-1", "showFlightCard", map[string]any{
                "from": "SFO", "to": "JFK",
            }, nil),
        },
        "state": map[string]any{"lastBooking": "AA-12"},
    }, nil
})
```

For a Go reasoner using the AI tool-call loop:

```go
res, _ := aiClient.ExecuteToolCallLoopResult(ctx, prompt, tools, callFn)
return map[string]any{
    "result":    res.Text(),
    "toolCalls": agui.ToolCallsFromTrace(res.Trace),
}, nil
```

## Frontend wiring

Standard CopilotKit App Router setup, with one `HttpAgent` per reasoner:

```ts
// app/api/copilotkit/route.ts
import { CopilotRuntime, copilotRuntimeNextJSAppRouterEndpoint } from "@copilotkit/runtime";
import { HttpAgent } from "@ag-ui/client";

const BASE = "http://your-control-plane/api/v1/agui/runs/your-node";

const runtime = new CopilotRuntime({
  agents: {
    chat:        new HttpAgent({ url: `${BASE}/chat` }),
    book_flight: new HttpAgent({ url: `${BASE}/book_flight` }),
  },
});

export const POST = async (req: Request) => {
  const { handleRequest } = copilotRuntimeNextJSAppRouterEndpoint({
    runtime, endpoint: "/api/copilotkit",
  });
  return handleRequest(req);
};
```

```tsx
// app/page.tsx
"use client";
import { CopilotKit, useCopilotAction } from "@copilotkit/react-core";
import { CopilotChat } from "@copilotkit/react-ui";
import "@copilotkit/react-ui/styles.css";

function FlightCard({ from, to, depart }: any) {
  return <div className="flight-card">{from} → {to} ({depart})</div>;
}

function Page() {
  // Render-only: the agent emits a TOOL_CALL_*; the UI just visualizes it.
  // `available: "frontend"` is required for render-only actions in
  // CopilotKit v1.57+.
  useCopilotAction({
    name: "showFlightCard",
    available: "frontend",
    parameters: [
      { name: "from", type: "string" },
      { name: "to", type: "string" },
      { name: "depart", type: "string" },
    ],
    render: ({ args }) => <FlightCard {...(args as any)} />,
  });

  return (
    <CopilotKit runtimeUrl="/api/copilotkit" agent="book_flight">
      <CopilotChat />
    </CopilotKit>
  );
}
```

For round-trip frontend tools (the agent calls a tool, the user
interacts, the tool returns a result that loops back to the agent on
the next turn), use `available: "enabled"` with a `handler` instead of
`render`. CopilotKit posts the tool's return value as a
`role: "tool"` message in the next run — the control plane forwards it
intact to the reasoner.

## Auth

The endpoint sits behind the same DID/VC permission middleware as
`/execute`. When `AGENTFIELD_FEATURES_DID_AUTHORIZATION_ENABLED=true`,
callers must include a valid DID-signed request just like for direct
reasoner invocations.

## Live streaming (per-token + per-tool-arg deltas)

The reasoner contract above buffers a full response and returns it as a
single dict. For live UX — text appearing token-by-token,
`TOOL_CALL_ARGS` streaming as the LLM emits them, `REASONING_*` events
flowing as the model thinks — return an NDJSON stream instead. The
control plane's streaming dispatcher (see
`control-plane/internal/handlers/agui_runs_streaming.go`) detects
`Content-Type: application/x-ndjson` and translates each line into the
matching AG-UI event in real time.

### Python streaming reasoner

```python
from fastapi import Request
from fastapi.responses import StreamingResponse
from agentfield import Agent, agui

app = Agent(node_id="my-app")

@app.post("/reasoners/chat")
async def chat(request: Request):
    body = await request.json()
    return StreamingResponse(
        agui.serialize_stream(_chunks(body)),
        media_type=agui.STREAMING_CONTENT_TYPE,
    )

async def _chunks(body):
    # Reasoning shows up in CopilotKit's "Thinking…" pane.
    yield agui.reasoning_chunk("Looking up flights...")
    yield agui.reasoning_end_chunk()
    # Text chunks paint progressively in <CopilotChat>.
    async for token in llm.stream(body["prompt"]):
        yield agui.text_chunk(token)
    # Tool calls drive useCopilotAction renders.
    yield agui.tool_call_start_chunk("tc-1", "showFlightCard",
                                     arguments={"from": "SFO", "to": "JFK"})
    yield agui.tool_call_end_chunk("tc-1")
    # Shared state lands in useCoAgent.
    yield agui.state_chunk({"counter": 1})
```

The control plane wraps the stream with `RUN_STARTED` / `RUN_FINISHED`,
manages text and reasoning open/close lifecycle automatically, and emits
`MESSAGES_SNAPSHOT` at stream end.

### Go streaming reasoner

```go
import (
    "net/http"
    "github.com/Agent-Field/agentfield/sdk/go/agent/agui"
)

func chat(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", agui.StreamingContentType)
    w.WriteHeader(http.StatusOK)

    chunks := make(chan map[string]any, 8)
    go func() {
        defer close(chunks)
        chunks <- agui.ReasoningChunk("Looking up flights...")
        chunks <- agui.ReasoningEndChunk()
        for _, tok := range []string{"Booked ", "AA-12."} {
            chunks <- agui.TextChunk(tok)
        }
        chunks <- agui.ToolCallStartChunk("tc-1", "showFlightCard",
            map[string]any{"from": "SFO", "to": "JFK"}, "")
        chunks <- agui.ToolCallEndChunk("tc-1")
    }()
    _ = agui.SerializeStream(r.Context(), w, chunks)
}
```

### `.harness()` relay

The Anthropic Claude harness already produces a streaming async iterator
of messages. Pipe it straight to AG-UI:

```python
from claude_agent_sdk import query, ClaudeAgentOptions
from agentfield import agui

async def _chunks(body):
    opts = ClaudeAgentOptions(...)
    async for chunk in agui.relay_harness_stream(
        query(prompt=body["prompt"], options=opts)
    ):
        yield chunk
```

`relay_harness_stream` translates Claude SDK message types into the
right AG-UI chunks: `text` blocks → `TEXT_MESSAGE_CONTENT`,
`thinking` blocks → `REASONING_*`, `tool_use` blocks → `TOOL_CALL_*`,
`tool_result` blocks → `TOOL_CALL_RESULT`. Note: the harness streams
per-message, not per-token, so this path delivers message-level
streaming. True per-token streaming requires the raw Anthropic API.

## Reasoner contract — full chunk reference

When using the streaming path, each NDJSON line is one of these tagged
chunks (built by helpers in `agentfield.agui` / `sdk/go/agent/agui`):

| Chunk `type` | Maps to | Notes |
|---|---|---|
| `text` | `TEXT_MESSAGE_CONTENT` | `START`/`END` synthesized lazily on first/last text chunk |
| `reasoning` | `REASONING_MESSAGE_CONTENT` | Outer `REASONING_START`/`END` synthesized; emit `reasoning_end` to start a new segment within the same context |
| `tool_call_start` | `TOOL_CALL_START` (+ `_ARGS` if `arguments` provided inline) | |
| `tool_call_args` | `TOOL_CALL_ARGS` | Streamed as the LLM emits arg JSON |
| `tool_call_end` | `TOOL_CALL_END` | |
| `tool_call_result` | `TOOL_CALL_RESULT` | For server-side tools |
| `state` | `STATE_SNAPSHOT` | |
| `state_delta` | `STATE_DELTA` (RFC 6902 patches) | |
| `step_started` / `step_finished` | `STEP_STARTED` / `STEP_FINISHED` | CopilotKit ignores; useful for other AG-UI consumers |
| `raw` | `RAW` | Foreign-system passthrough |
| `custom` | `CUSTOM` | App-specific event with `name` + `value` |
| `final` | Applies a buffered-shape envelope | Use to send trailing `toolCalls` / `state` / etc. without re-implementing buffered logic |
| `error` | `RUN_ERROR` (terminal) | Subsequent chunks are ignored |

## Performance

Load tested at 50× concurrent buffered requests and 25× concurrent
streaming requests in CI (`internal/handlers/agui_runs_load_test.go`):

- Buffered: 200 reqs in ~90 ms wall, p50 ≈ 4 ms, p95 ≈ 75 ms, p99 ≈ 77 ms
- Streaming dispatcher: 100 reqs in ~18 ms wall, no goroutine leaks
- Per-request benchmark (`go test -bench=BenchmarkAGUI`): ~389 µs/op, 26 KB/op

## What we don't yet do

- **Per-token streaming via the buffered reasoner contract.** Reasoners
  using `@app.reasoner()` still buffer; the streaming path requires the
  separate FastAPI / chunk-channel pattern shown above. We auto-chunk
  buffered responses on emission so the UX is acceptable, but the
  source of truth is still a synchronous return.
- **Bidirectional cancellation propagation into the streaming reasoner.**
  Client disconnect aborts the streaming HTTP read on our end, but the
  reasoner needs its own context plumbing to actually stop work.
