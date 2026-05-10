"""Tests for the agentfield.agui helpers — the documented contract for
opt-in Generative UI / shared state through the control plane's AG-UI
adapter."""

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
