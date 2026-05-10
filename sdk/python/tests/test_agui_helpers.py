"""Tests for the agentfield.agui helpers — the documented contract for
opt-in Generative UI / shared state through the control plane's AG-UI
adapter."""

import json

import pytest

from agentfield import agui
from agentfield.tool_calling import ToolCallRecord, ToolCallTrace


class TestToolCall:
    def test_minimal(self):
        e = agui.tool_call(name="showFlightCard")
        assert e == {"name": "showFlightCard", "arguments": {}}
        assert "result" not in e
        assert "id" not in e

    def test_with_arguments_and_id(self):
        e = agui.tool_call(name="x", arguments={"a": 1, "b": "z"}, id="tc-1")
        assert e == {"name": "x", "arguments": {"a": 1, "b": "z"}, "id": "tc-1"}

    def test_with_result_attaches_for_executed_calls(self):
        e = agui.tool_call(name="getWeather", result={"temp": 62})
        assert e["result"] == {"temp": 62}

    def test_explicit_null_result(self):
        e = agui.tool_call(name="x", has_result=True)
        assert "result" in e
        assert e["result"] is None


class TestToolCallsFromTrace:
    def test_none_trace_returns_empty(self):
        assert agui.tool_calls_from_trace(None) == []

    def test_empty_calls_returns_empty(self):
        trace = ToolCallTrace(calls=[])
        assert agui.tool_calls_from_trace(trace) == []

    def test_records_become_entries_with_results(self):
        trace = ToolCallTrace(
            calls=[
                ToolCallRecord(
                    tool_name="getWeather",
                    arguments={"city": "SF"},
                    result={"temp": 62},
                ),
                ToolCallRecord(
                    tool_name="lookup",
                    arguments={"q": "gates"},
                    error="api timeout",
                ),
            ]
        )
        out = agui.tool_calls_from_trace(trace)
        assert len(out) == 2
        assert out[0]["name"] == "getWeather"
        assert out[0]["arguments"] == {"city": "SF"}
        assert out[0]["result"] == {"temp": 62}
        # id is synthesized so the control plane can correlate frames
        # without colliding across calls in the same trace.
        assert out[0]["id"] == "tc-trace-0"

        assert out[1]["name"] == "lookup"
        assert out[1]["result"] == {"error": "api timeout"}
        assert out[1]["id"] == "tc-trace-1"

    def test_trace_with_no_result_or_error_omits_result_field(self):
        trace = ToolCallTrace(calls=[ToolCallRecord(tool_name="x", arguments={})])
        out = agui.tool_calls_from_trace(trace)
        assert "result" not in out[0]


class TestStateDeltaHelpers:
    def test_replace_op(self):
        assert agui.state_delta_replace("/counter", 2) == {
            "op": "replace",
            "path": "/counter",
            "value": 2,
        }

    def test_replace_rejects_invalid_path(self):
        with pytest.raises(ValueError):
            agui.state_delta_replace("counter", 2)  # missing leading slash

    def test_diff_emits_replace_for_changed(self):
        ops = agui.state_delta_from_diff({"a": 1, "b": 2}, {"a": 1, "b": 3})
        assert ops == [{"op": "replace", "path": "/b", "value": 3}]

    def test_diff_emits_add_for_new_keys(self):
        ops = agui.state_delta_from_diff({}, {"x": 1})
        assert ops == [{"op": "add", "path": "/x", "value": 1}]

    def test_diff_emits_remove_for_dropped_keys(self):
        ops = agui.state_delta_from_diff({"x": 1}, {})
        assert ops == [{"op": "remove", "path": "/x"}]

    def test_diff_no_ops_when_identical(self):
        assert agui.state_delta_from_diff({"a": 1}, {"a": 1}) == []


class TestReasoningHelpers:
    def test_segment_minimal(self):
        seg = agui.reasoning_segment("thinking")
        assert seg == {"content": "thinking"}

    def test_segment_with_id(self):
        seg = agui.reasoning_segment("thinking", id="r-1")
        assert seg == {"content": "thinking", "id": "r-1"}

    def test_reasoning_strings_pass_through(self):
        assert agui.reasoning("step 1", "step 2") == ["step 1", "step 2"]

    def test_reasoning_drops_empty_strings(self):
        assert agui.reasoning("step 1", "", "step 2") == ["step 1", "step 2"]

    def test_reasoning_accepts_segment_dicts(self):
        out = agui.reasoning("step 1", agui.reasoning_segment("step 2", id="r-2"))
        assert out == ["step 1", {"content": "step 2", "id": "r-2"}]

    def test_reasoning_rejects_unknown_types(self):
        with pytest.raises(TypeError):
            agui.reasoning(42)


class TestStreamingChunkBuilders:
    def test_text_chunk(self):
        assert agui.text_chunk("hello") == {"type": "text", "delta": "hello"}

    def test_reasoning_chunks(self):
        assert agui.reasoning_chunk("think") == {"type": "reasoning", "delta": "think"}
        assert agui.reasoning_end_chunk() == {"type": "reasoning_end"}

    def test_tool_call_start_with_args_inline(self):
        c = agui.tool_call_start_chunk("tc1", "showCard", arguments={"x": 1})
        assert c == {"type": "tool_call_start", "id": "tc1", "name": "showCard", "arguments": {"x": 1}}

    def test_tool_call_start_with_parent(self):
        c = agui.tool_call_start_chunk("tc1", "x", parent_message_id="m1")
        assert c["parentMessageId"] == "m1"

    def test_tool_call_args_stream(self):
        assert agui.tool_call_args_chunk("tc1", '{"x') == {
            "type": "tool_call_args",
            "id": "tc1",
            "delta": '{"x',
        }

    def test_tool_call_end_and_result(self):
        assert agui.tool_call_end_chunk("tc1") == {"type": "tool_call_end", "id": "tc1"}
        r = agui.tool_call_result_chunk("tc1", "ok", role="tool")
        assert r == {"type": "tool_call_result", "id": "tc1", "content": "ok", "role": "tool"}

    def test_state_chunks(self):
        assert agui.state_chunk({"counter": 1}) == {"type": "state", "snapshot": {"counter": 1}}
        assert agui.state_delta_chunk([{"op": "replace", "path": "/x", "value": 1}]) == {
            "type": "state_delta",
            "ops": [{"op": "replace", "path": "/x", "value": 1}],
        }

    def test_step_chunks(self):
        assert agui.step_started_chunk("plan") == {"type": "step_started", "name": "plan"}
        assert agui.step_finished_chunk("plan") == {"type": "step_finished", "name": "plan"}

    def test_raw_and_custom(self):
        assert agui.raw_chunk({"x": 1}, source="harness") == {
            "type": "raw",
            "event": {"x": 1},
            "source": "harness",
        }
        assert agui.custom_chunk("ack", value={"ok": True}) == {
            "type": "custom",
            "name": "ack",
            "value": {"ok": True},
        }

    def test_final_chunk(self):
        c = agui.final_chunk({"toolCalls": [{"name": "x"}]})
        assert c == {"type": "final", "data": {"toolCalls": [{"name": "x"}]}}

    def test_error_chunk(self):
        assert agui.error_chunk("boom", code="E1") == {
            "type": "error",
            "message": "boom",
            "code": "E1",
        }


class TestSerializeStream:
    @pytest.mark.asyncio
    async def test_yields_ndjson_lines(self):
        async def gen():
            yield agui.text_chunk("hello ")
            yield agui.text_chunk("world")
            yield agui.tool_call_start_chunk("tc1", "x")

        lines = []
        async for chunk in agui.serialize_stream(gen()):
            assert isinstance(chunk, bytes)
            assert chunk.endswith(b"\n")
            lines.append(chunk.decode("utf-8").rstrip("\n"))
        assert len(lines) == 3
        assert json.loads(lines[0]) == {"type": "text", "delta": "hello "}
        assert json.loads(lines[1]) == {"type": "text", "delta": "world"}
        assert json.loads(lines[2])["type"] == "tool_call_start"

    @pytest.mark.asyncio
    async def test_bare_string_wraps_as_text_chunk(self):
        async def gen():
            yield "ergonomic"

        out = []
        async for chunk in agui.serialize_stream(gen()):
            out.append(json.loads(chunk))
        assert out == [{"type": "text", "delta": "ergonomic"}]

    @pytest.mark.asyncio
    async def test_invalid_yield_raises_typeerror(self):
        async def gen():
            yield 42  # not str / not dict

        with pytest.raises(TypeError):
            async for _ in agui.serialize_stream(gen()):
                pass


class TestHarnessRelay:
    """Coverage for relay_harness_stream — the bridge that turns a Claude
    Agent SDK / harness async iterator into AG-UI streaming chunks."""

    @pytest.mark.asyncio
    async def test_assistant_text_block_becomes_text_chunk(self):
        async def fake_harness():
            yield {
                "type": "assistant",
                "message": {"content": [{"type": "text", "text": "Hello!"}]},
            }
            yield {"type": "result", "subtype": "success", "result": "Hello!"}

        chunks = [c async for c in agui.relay_harness_stream(fake_harness())]
        # result message yields nothing; only the text chunk survives.
        assert chunks == [{"type": "text", "delta": "Hello!"}]

    @pytest.mark.asyncio
    async def test_thinking_block_becomes_reasoning_chunk(self):
        async def fake():
            yield {
                "type": "assistant",
                "message": {"content": [
                    {"type": "thinking", "thinking": "Let me think..."},
                    {"type": "text", "text": "Done."},
                ]},
            }

        chunks = [c async for c in agui.relay_harness_stream(fake())]
        assert chunks[0] == {"type": "reasoning", "delta": "Let me think..."}
        assert chunks[1] == {"type": "text", "delta": "Done."}

    @pytest.mark.asyncio
    async def test_tool_use_emits_start_and_end(self):
        async def fake():
            yield {
                "type": "assistant",
                "message": {"content": [{
                    "type": "tool_use",
                    "id": "tu-1",
                    "name": "get_weather",
                    "input": {"city": "SF"},
                }]},
            }

        chunks = [c async for c in agui.relay_harness_stream(fake())]
        assert chunks[0]["type"] == "tool_call_start"
        assert chunks[0]["id"] == "tu-1"
        assert chunks[0]["name"] == "get_weather"
        assert chunks[0]["arguments"] == {"city": "SF"}
        assert chunks[1] == {"type": "tool_call_end", "id": "tu-1"}

    @pytest.mark.asyncio
    async def test_tool_result_emits_result_chunk(self):
        async def fake():
            yield {
                "type": "user",
                "message": {"content": [{
                    "type": "tool_result",
                    "tool_use_id": "tu-1",
                    "content": "62°F, foggy",
                }]},
            }

        chunks = [c async for c in agui.relay_harness_stream(fake())]
        assert chunks[0]["type"] == "tool_call_result"
        assert chunks[0]["id"] == "tu-1"
        assert chunks[0]["content"] == "62°F, foggy"

    @pytest.mark.asyncio
    async def test_tool_result_with_block_list_stitches_text(self):
        async def fake():
            yield {
                "type": "user",
                "message": {"content": [{
                    "type": "tool_result",
                    "tool_use_id": "tu-1",
                    "content": [
                        {"type": "text", "text": "part 1 "},
                        {"type": "text", "text": "part 2"},
                    ],
                }]},
            }

        chunks = [c async for c in agui.relay_harness_stream(fake())]
        assert chunks[0]["content"] == "part 1 part 2"

    @pytest.mark.asyncio
    async def test_unknown_block_falls_back_to_raw(self):
        async def fake():
            yield {
                "type": "assistant",
                "message": {"content": [{"type": "weird-thing", "data": 42}]},
            }

        chunks = [c async for c in agui.relay_harness_stream(fake())]
        assert chunks[0]["type"] == "raw"
        assert chunks[0]["source"] == "harness"

    @pytest.mark.asyncio
    async def test_unknown_message_type_becomes_raw(self):
        async def fake():
            yield {"type": "system", "info": "starting"}
            yield {"type": "totally_unknown", "x": 1}

        chunks = [c async for c in agui.relay_harness_stream(fake())]
        assert all(c["type"] == "raw" for c in chunks)

    @pytest.mark.asyncio
    async def test_result_message_yields_nothing(self):
        async def fake():
            yield {"type": "result", "subtype": "success", "result": "done"}

        chunks = [c async for c in agui.relay_harness_stream(fake())]
        assert chunks == []


class TestStreamingFastAPIRoundTrip:
    """End-to-end: a FastAPI app using StreamingResponse + serialize_stream
    must produce exactly the NDJSON bytes the control plane's streaming
    dispatcher consumes. This is the SDK-side test of the wire contract."""

    @pytest.mark.asyncio
    async def test_streaming_endpoint_returns_ndjson(self):
        from fastapi import FastAPI
        from fastapi.responses import StreamingResponse
        from httpx import ASGITransport, AsyncClient

        app = FastAPI()

        async def chunks():
            yield agui.reasoning_chunk("checking flights...")
            yield agui.text_chunk("Booked ")
            yield agui.text_chunk("AA-12.")
            yield agui.tool_call_start_chunk("tc1", "showFlightCard", arguments={"flight": "AA-12"})
            yield agui.tool_call_end_chunk("tc1")
            yield agui.state_chunk({"counter": 1})

        @app.post("/reasoners/chat")
        async def chat():
            return StreamingResponse(
                agui.serialize_stream(chunks()),
                media_type=agui.STREAMING_CONTENT_TYPE,
            )

        async with AsyncClient(transport=ASGITransport(app=app), base_url="http://test") as client:
            resp = await client.post("/reasoners/chat")
            assert resp.status_code == 200
            assert resp.headers["content-type"].startswith("application/x-ndjson")
            lines = [line for line in resp.text.split("\n") if line]
            assert len(lines) == 6
            decoded = [json.loads(line) for line in lines]
            assert decoded[0]["type"] == "reasoning"
            assert decoded[1]["type"] == "text"
            assert decoded[1]["delta"] == "Booked "
            assert decoded[3]["type"] == "tool_call_start"
            assert decoded[3]["arguments"] == {"flight": "AA-12"}
            assert decoded[5]["type"] == "state"
