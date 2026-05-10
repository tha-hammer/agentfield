"""AG-UI protocol helpers for AgentField reasoners.

This module also exposes a *streaming* reasoner contract — see
``serialize_stream`` and the chunk builders (``text_chunk``,
``reasoning_chunk``, ``tool_call_start_chunk`` …) — for live
per-token AG-UI events. A streaming reasoner is a normal FastAPI
endpoint that returns a ``StreamingResponse`` with content-type
``application/x-ndjson``; the AgentField control plane sniffs the
content-type and dispatches each line as a live AG-UI event.



Reasoners reach the AG-UI / CopilotKit frontend via the control plane's
``POST /api/v1/agui/runs/<node>/<reasoner>`` adapter. The adapter expects
a small set of optional fields in the reasoner's response to drive the
richer AG-UI events (tool calls, shared state, RFC 6902 patches).

This module is the documented contract for those fields. Reasoner authors
opt into Generative UI / shared state by returning the values these
helpers build:

.. code-block:: python

    @app.reasoner()
    async def book_flight(prompt: str = "", state: dict | None = None):
        return {
            "result": "Picking flight options.",
            "toolCalls": [
                tool_call(name="showFlightCard", arguments={"from": "SFO", "to": "JFK"}),
            ],
            "state": {"counter": (state or {}).get("counter", 0) + 1},
        }

When a reasoner uses ``await app.ai(..., tools=[...])`` and wants the
LLM's tool-calling trace to surface in the UI, pass the returned
``ToolCallResponse.trace`` into :func:`tool_calls_from_trace`:

.. code-block:: python

    result = await app.ai("help the user", tools="discover")
    return {
        "result": result.text,
        "toolCalls": tool_calls_from_trace(result.trace),
    }

Wire shape mirrors the canonical AG-UI ``TOOL_CALL_*`` events
(https://docs.ag-ui.com/concepts/events).
"""

from __future__ import annotations

from typing import Any, Iterable, List, Mapping, Optional

from .tool_calling import ToolCallRecord, ToolCallTrace

__all__ = [
    "tool_call",
    "tool_calls_from_trace",
    "state_delta_replace",
    "state_delta_from_diff",
    "reasoning",
    "reasoning_segment",
    # Streaming chunk builders + serializer for the live AG-UI path.
    "text_chunk",
    "reasoning_chunk",
    "reasoning_end_chunk",
    "tool_call_start_chunk",
    "tool_call_args_chunk",
    "tool_call_end_chunk",
    "tool_call_result_chunk",
    "state_chunk",
    "state_delta_chunk",
    "step_started_chunk",
    "step_finished_chunk",
    "raw_chunk",
    "custom_chunk",
    "final_chunk",
    "error_chunk",
    "serialize_stream",
    "relay_harness_stream",
    "STREAMING_CONTENT_TYPE",
]

STREAMING_CONTENT_TYPE = "application/x-ndjson"


def tool_call(
    name: str,
    arguments: Optional[Mapping[str, Any]] = None,
    *,
    id: Optional[str] = None,
    result: Any = None,
    has_result: bool = False,
) -> dict:
    """Build a single AG-UI tool-call entry.

    The control plane translates each entry into a
    ``TOOL_CALL_START`` / ``TOOL_CALL_ARGS`` / ``TOOL_CALL_END`` triad.
    When ``has_result=True`` (or ``result`` is non-None), it also emits
    ``TOOL_CALL_RESULT`` so a server-side trace renders in the UI.

    Args:
        name: The tool name. CopilotKit pattern-matches this against
            ``useCopilotAction({name, render})`` registrations to drive
            Generative UI.
        arguments: A JSON-serializable mapping of arguments.
        id: Optional stable ID. If omitted, the control plane synthesizes
            one (which works for one-shot calls but breaks correlation
            with follow-up tool messages).
        result: Optional result. Set this when the tool was already
            executed server-side (e.g. inside ``app.ai(tools=...)``).
        has_result: Pass True to force ``result=None`` to be treated as
            an explicit "executed and returned null" instead of "not
            executed yet". Defaults to True if ``result`` is non-None.
    """
    entry: dict = {"name": name, "arguments": dict(arguments or {})}
    if id is not None:
        entry["id"] = id
    if result is not None or has_result:
        entry["result"] = result
    return entry


def tool_calls_from_trace(trace: Optional[ToolCallTrace]) -> List[dict]:
    """Convert a ``ToolCallTrace`` from ``app.ai(tools=...)`` into the
    AG-UI ``toolCalls`` list shape.

    Each :class:`ToolCallRecord` becomes a tool-call entry with its
    arguments, and the executed result (or error) attached so the UI can
    render the trace as a sequence of completed tool calls. Empty traces
    return ``[]`` so callers can splat the result safely:

    .. code-block:: python

        return {"result": text, "toolCalls": tool_calls_from_trace(trace)}

    Args:
        trace: A trace from :class:`ToolCallResponse`, or None.

    Returns:
        A list of dicts in AG-UI ``toolCalls`` format. Empty if ``trace``
        is None or has no calls.
    """
    if trace is None or not getattr(trace, "calls", None):
        return []
    out: List[dict] = []
    for i, rec in enumerate(trace.calls):
        out.append(_record_to_entry(rec, i))
    return out


def _record_to_entry(rec: ToolCallRecord, index: int) -> dict:
    """Translate one ``ToolCallRecord`` into an AG-UI tool-call entry."""
    entry: dict = {
        "id": f"tc-trace-{index}",
        "name": rec.tool_name,
        "arguments": dict(rec.arguments or {}),
    }
    # The trace records either a result or an error; surface either as
    # the AG-UI tool-call result so frontend renderers can show a final
    # state instead of a perpetually "running" placeholder.
    if rec.error is not None:
        entry["result"] = {"error": rec.error}
    elif rec.result is not None:
        entry["result"] = rec.result
    return entry


def state_delta_replace(path: str, value: Any) -> dict:
    """Build a single RFC 6902 ``replace`` patch op for ``stateDelta``.

    .. code-block:: python

        return {
            "result": "...",
            "stateDelta": [
                state_delta_replace("/counter", 2),
                state_delta_replace("/lastUpdated", "2026-05-09"),
            ],
        }

    The control plane re-emits the array as a ``STATE_DELTA`` event,
    which CopilotKit's ``useCoAgent`` applies on top of the previously
    snapshot-emitted state.
    """
    if not path.startswith("/"):
        raise ValueError("RFC 6902 paths must start with '/'")
    return {"op": "replace", "path": path, "value": value}


def reasoning_segment(content: str, *, id: Optional[str] = None) -> dict:
    """Build a single REASONING_MESSAGE segment.

    Reasoners surface chain-of-thought to CopilotKit's "Thinking…" pane
    by returning either a plain string or a list of these segments under
    the ``reasoning`` field of their response:

    .. code-block:: python

        return {
            "result": "Booked AA-12.",
            "reasoning": [
                agui.reasoning_segment("Looking up flights for SFO->JFK..."),
                agui.reasoning_segment("AA-12 is the cheapest non-stop."),
            ],
        }

    Each segment becomes a REASONING_MESSAGE_START / _CONTENT / _END
    triad inside a REASONING_START / END boundary. Long content is
    auto-chunked across multiple REASONING_MESSAGE_CONTENT deltas.
    """
    out: dict = {"content": content}
    if id is not None:
        out["id"] = id
    return out


def reasoning(*segments: Any) -> List[Any]:
    """Build a ``reasoning`` field value from a mix of strings and segments.

    Convenience wrapper so reasoners can write::

        return {"result": text, "reasoning": agui.reasoning("step 1", "step 2")}

    instead of constructing the list manually.
    """
    out: List[Any] = []
    for s in segments:
        if isinstance(s, str):
            if s:
                out.append(s)
        elif isinstance(s, Mapping):
            out.append(dict(s))
        else:
            raise TypeError(
                f"reasoning() segments must be str or mapping; got {type(s).__name__}"
            )
    return out


# ---------------------------------------------------------------------------
# Streaming chunk builders
#
# Each function returns a small dict in the wire shape the control plane's
# streaming dispatcher consumes (see internal/handlers/agui_runs_streaming.go).
# The reasoner author yields these from an async generator; serialize_stream
# turns each yield into one NDJSON line for the FastAPI StreamingResponse.
# ---------------------------------------------------------------------------


def text_chunk(delta: str) -> dict:
    """One chunk of assistant text. Concatenated client-side."""
    return {"type": "text", "delta": delta}


def reasoning_chunk(delta: str) -> dict:
    """One chunk of chain-of-thought, rendered in CopilotKit's
    "Thinking…" pane. Yield multiple in a row for a single thought,
    then ``reasoning_end_chunk()`` to start a new thought segment."""
    return {"type": "reasoning", "delta": delta}


def reasoning_end_chunk() -> dict:
    """Closes the current reasoning segment so the next ``reasoning_chunk``
    opens a fresh one. The outer reasoning context auto-closes at stream
    end or when the first text/tool-call chunk arrives."""
    return {"type": "reasoning_end"}


def tool_call_start_chunk(
    id: str,
    name: str,
    *,
    arguments: Optional[Mapping[str, Any]] = None,
    parent_message_id: Optional[str] = None,
) -> dict:
    """Open a tool call. If you already have the full ``arguments``,
    pass them here and the dispatcher emits one TOOL_CALL_ARGS frame
    immediately; otherwise stream them with ``tool_call_args_chunk``."""
    out: dict = {"type": "tool_call_start", "id": id, "name": name}
    if arguments is not None:
        out["arguments"] = dict(arguments)
    if parent_message_id is not None:
        out["parentMessageId"] = parent_message_id
    return out


def tool_call_args_chunk(id: str, delta: str) -> dict:
    """One chunk of streaming tool-call arguments. ``delta`` is a
    string — typically a piece of the JSON-encoded arguments object as
    the LLM emits it. Concatenated client-side into the final args JSON."""
    return {"type": "tool_call_args", "id": id, "delta": delta}


def tool_call_end_chunk(id: str) -> dict:
    """Close a tool call."""
    return {"type": "tool_call_end", "id": id}


def tool_call_result_chunk(id: str, content: str, *, role: str = "tool") -> dict:
    """Server-side tool result. Use when the reasoner already executed
    the tool (e.g. via ``app.ai(tools=...)``) and wants the trace to
    render as completed in the UI."""
    return {"type": "tool_call_result", "id": id, "content": content, "role": role}


def state_chunk(snapshot: Any) -> dict:
    """Full agent state snapshot — the value ``useCoAgent({state})``
    reads on the frontend."""
    return {"type": "state", "snapshot": snapshot}


def state_delta_chunk(ops: List[dict]) -> dict:
    """RFC 6902 patch ops applied incrementally on top of the last
    snapshot the client received. Cheaper than re-emitting full state
    every turn."""
    return {"type": "state_delta", "ops": list(ops)}


def step_started_chunk(name: str) -> dict:
    """Mark the start of a named step inside the run. Useful for
    multi-stage agents where a frontend wants to render a progress UI."""
    return {"type": "step_started", "name": name}


def step_finished_chunk(name: str) -> dict:
    """Mark a step finished."""
    return {"type": "step_finished", "name": name}


def raw_chunk(event: Any, *, source: Optional[str] = None) -> dict:
    """Pass a foreign-system event through verbatim. Frontends that
    subscribed via ``onRawEvent`` see it; others ignore it."""
    out: dict = {"type": "raw", "event": event}
    if source is not None:
        out["source"] = source
    return out


def custom_chunk(name: str, value: Any = None) -> dict:
    """Application-defined event. Frontends subscribe by ``name``."""
    out: dict = {"type": "custom", "name": name}
    if value is not None:
        out["value"] = value
    return out


def final_chunk(data: Mapping[str, Any]) -> dict:
    """Trailing buffered envelope — the dispatcher applies any
    ``toolCalls`` / ``state`` / ``stateDelta`` / ``reasoning`` /
    ``result`` fields here as if from a non-streaming reasoner. Useful
    when the reasoner can stream text live but only knows the
    structured fields at the end."""
    return {"type": "final", "data": dict(data)}


def error_chunk(message: str, *, code: Optional[str] = None) -> dict:
    """Terminal error. The dispatcher emits RUN_ERROR and stops the run;
    any subsequent chunks the reasoner sends are ignored."""
    out: dict = {"type": "error", "message": message}
    if code is not None:
        out["code"] = code
    return out


async def relay_harness_stream(harness_iter: Any) -> Any:
    """Relay a Claude Agent SDK / harness async-iterator of messages
    into AG-UI streaming chunks, message-by-message.

    The Claude Agent SDK yields one Python dict (or message object) per
    turn — assistant text blocks, tool-use blocks, tool-result blocks,
    a final ``result`` envelope. This function translates each into the
    smallest sensible AG-UI chunk(s) so a reasoner can pipe a harness
    run straight to the AG-UI stream::

        from claude_agent_sdk import query, ClaudeAgentOptions
        from agentfield import agui

        async def _chunks(body):
            opts = ClaudeAgentOptions(...)
            async for ch in agui.relay_harness_stream(
                query(prompt=body["prompt"], options=opts)
            ):
                yield ch

    Recognized message shapes (matches the dict form
    ``HarnessResult.messages`` records):

      - ``{"type":"assistant","message":{"content":[{"type":"text","text":"..."}, ...]}}``
        → one ``text`` chunk per text block
      - ``{"type":"assistant","message":{"content":[{"type":"tool_use","id":"...","name":"...","input":{...}}, ...]}}``
        → ``tool_call_start`` + ``tool_call_end`` per tool_use block
      - ``{"type":"user","message":{"content":[{"type":"tool_result","tool_use_id":"...","content":"..."}, ...]}}``
        → ``tool_call_result`` per tool_result block
      - ``{"type":"result","subtype":"success","result":"..."}`` →
        terminal — yields nothing (the dispatcher's stream-end logic
        wraps the run with MESSAGES_SNAPSHOT + RUN_FINISHED).
      - Anything unrecognized is wrapped as a ``raw`` chunk so the trace
        is preserved without us inventing ad-hoc event types.

    Note: the Claude Agent SDK buffers per-message rather than per-token,
    so this path streams at message granularity. True per-token streaming
    requires the raw Anthropic streaming API, not the harness."""
    async for raw in harness_iter:
        if isinstance(raw, dict):
            msg = raw
        elif hasattr(raw, "__dict__"):
            msg = dict(raw.__dict__)
        else:
            yield raw_chunk({"raw": str(raw)}, source="harness")
            continue

        msg_type = str(msg.get("type", ""))
        if msg_type == "result":
            # The harness's result message holds the final aggregated
            # text; the AG-UI stream's MESSAGES_SNAPSHOT / RUN_FINISHED
            # frames will be synthesized by the control-plane dispatcher
            # at stream end, so we don't need to emit anything here.
            continue
        if msg_type == "system":
            yield raw_chunk(msg, source="harness")
            continue

        if msg_type in ("assistant", "user"):
            content = _harness_message_content(msg)
            if content is None:
                yield raw_chunk(msg, source="harness")
                continue
            if isinstance(content, str):
                if msg_type == "assistant" and content:
                    yield text_chunk(content)
                continue
            if isinstance(content, list):
                for block in content:
                    if not isinstance(block, dict):
                        continue
                    btype = block.get("type")
                    if btype == "text":
                        text = block.get("text", "")
                        if text:
                            yield text_chunk(text)
                    elif btype == "thinking":
                        # Anthropic extended-thinking blocks render as
                        # REASONING_* events — exactly the "Thinking…"
                        # pane CopilotKit shows.
                        thinking = block.get("thinking", "")
                        if thinking:
                            yield reasoning_chunk(thinking)
                    elif btype == "tool_use":
                        tcid = str(block.get("id", ""))
                        name = str(block.get("name", ""))
                        if tcid and name:
                            inp = block.get("input")
                            if not isinstance(inp, Mapping):
                                inp = {}
                            yield tool_call_start_chunk(tcid, name, arguments=inp)
                            yield tool_call_end_chunk(tcid)
                    elif btype == "tool_result":
                        tcid = str(block.get("tool_use_id", ""))
                        if tcid:
                            inner = block.get("content", "")
                            if isinstance(inner, list):
                                # tool_result content may itself be a
                                # block list — stitch text blocks.
                                inner = "".join(
                                    str(b.get("text", "")) for b in inner if isinstance(b, dict)
                                )
                            elif not isinstance(inner, str):
                                inner = str(inner)
                            yield tool_call_result_chunk(tcid, inner, role="tool")
                    else:
                        yield raw_chunk(block, source="harness")
            continue

        # Unknown top-level message — preserve as raw.
        yield raw_chunk(msg, source="harness")


def _harness_message_content(msg: Mapping[str, Any]) -> Any:
    """Reach into the harness message envelope for the content list,
    handling both the bare ``content`` shape and the ``message.content``
    shape the Claude Agent SDK uses."""
    if "content" in msg:
        return msg["content"]
    inner = msg.get("message")
    if isinstance(inner, Mapping):
        return inner.get("content")
    return None


async def serialize_stream(generator: Any) -> Any:
    """Serialize an async generator of chunk dicts (or strings — strings
    are wrapped as text chunks) into an async iterator of NDJSON-encoded
    ``bytes``, suitable for ``fastapi.StreamingResponse``::

        from fastapi import Request
        from fastapi.responses import StreamingResponse
        from agentfield import agui

        @app.post("/reasoners/chat")
        async def chat(request: Request):
            body = await request.json()
            return StreamingResponse(
                agui.serialize_stream(_chat_chunks(body)),
                media_type=agui.STREAMING_CONTENT_TYPE,
            )

        async def _chat_chunks(body):
            async for token in llm.stream(body["prompt"]):
                yield agui.text_chunk(token)

    Bare strings yielded by the generator are auto-wrapped as text
    chunks for ergonomics. Anything else must be a dict produced by one
    of the chunk builders above (or a hand-rolled equivalent)."""
    import json as _json

    async for item in generator:
        if isinstance(item, str):
            payload = text_chunk(item)
        elif isinstance(item, Mapping):
            payload = dict(item)
        else:
            raise TypeError(
                "streaming reasoner yielded non-str/non-dict value of type "
                f"{type(item).__name__}; use one of agui.*_chunk(...)"
            )
        # No spaces — these are machine-to-machine; keep lines compact.
        yield (_json.dumps(payload, separators=(",", ":")) + "\n").encode("utf-8")


def state_delta_from_diff(
    before: Mapping[str, Any],
    after: Mapping[str, Any],
) -> List[dict]:
    """Compute a minimal RFC 6902 patch list for top-level keys that
    differ between ``before`` and ``after``.

    This is a deliberately shallow utility — it only walks the top level
    of the mapping and emits ``replace``/``add``/``remove`` ops as
    needed. Reasoners with nested state should construct patches
    explicitly (or just emit a full ``state`` snapshot).
    """
    ops: List[dict] = []
    keys: Iterable[str] = sorted(set(before.keys()) | set(after.keys()))
    for k in keys:
        path = f"/{k}"
        if k in before and k in after:
            if before[k] != after[k]:
                ops.append({"op": "replace", "path": path, "value": after[k]})
        elif k in after:
            ops.append({"op": "add", "path": path, "value": after[k]})
        else:
            ops.append({"op": "remove", "path": path})
    return ops
