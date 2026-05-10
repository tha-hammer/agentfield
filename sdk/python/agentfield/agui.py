"""AG-UI protocol helpers for AgentField reasoners.

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
]


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
