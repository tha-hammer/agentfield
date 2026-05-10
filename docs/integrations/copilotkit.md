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

## What we don't yet do

- **Live token streaming.** The reasoner returns a complete result; we
  chunk it on emission, but per-token streaming requires reasoner-side
  streaming, which is the next iteration. The `agentInvoker` interface
  in the handler is the seam where that will plug in.
- **Live tool-argument streaming.** `TOOL_CALL_ARGS` carries the full
  arguments JSON in one delta today, not progressive token chunks.
- **`STEP_*` / `RAW` / `CUSTOM` events.** CopilotKit ignores `STEP_*`
  per their `GOTCHAS.md`; the others are app-specific listener territory.
- **`.harness()` provider relay.** The Anthropic SDK already streams
  messages from the harness subprocess, but the current provider
  buffers them. Plumbing those out as nested `TEXT_MESSAGE_*` /
  `TOOL_CALL_*` is per-provider work.
