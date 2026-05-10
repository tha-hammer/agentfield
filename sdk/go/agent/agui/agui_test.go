package agui

import (
	"testing"

	"github.com/Agent-Field/agentfield/sdk/go/ai"

	"github.com/stretchr/testify/require"
)

func TestToolCall_MinimalAndFull(t *testing.T) {
	require.Nil(t, ToolCall("", "", nil, nil), "empty name returns nil so caller surfaces the bug")

	minimal := ToolCall("", "showFlightCard", nil, nil)
	require.Equal(t, "showFlightCard", minimal["name"])
	require.NotContains(t, minimal, "id", "id only present when caller supplies one")
	require.Equal(t, map[string]any{}, minimal["arguments"])
	require.NotContains(t, minimal, "result")

	full := ToolCall("tc-1", "x", map[string]any{"a": 1}, map[string]any{"ok": true})
	require.Equal(t, "tc-1", full["id"])
	require.Equal(t, map[string]any{"a": 1}, full["arguments"])
	require.Equal(t, map[string]any{"ok": true}, full["result"])
}

func TestToolCallsFromTrace(t *testing.T) {
	require.Empty(t, ToolCallsFromTrace(nil))
	require.Empty(t, ToolCallsFromTrace(&ai.ToolCallTrace{}))

	trace := &ai.ToolCallTrace{
		Calls: []ai.ToolCallRecord{
			{ToolName: "getWeather", Arguments: map[string]any{"city": "SF"}, Result: map[string]any{"temp": 62.0}},
			{ToolName: "lookup", Arguments: map[string]any{"q": "x"}, Error: "timeout"},
			{ToolName: "noargs"},
		},
	}
	out := ToolCallsFromTrace(trace)
	require.Len(t, out, 3)

	require.Equal(t, "tc-trace-0", out[0]["id"])
	require.Equal(t, "getWeather", out[0]["name"])
	require.Equal(t, map[string]any{"temp": 62.0}, out[0]["result"])

	require.Equal(t, "tc-trace-1", out[1]["id"])
	require.Equal(t, map[string]any{"error": "timeout"}, out[1]["result"], "errors surface as {error:...}")

	require.Equal(t, "tc-trace-2", out[2]["id"])
	require.Equal(t, map[string]any{}, out[2]["arguments"], "nil arguments default to empty map")
	require.NotContains(t, out[2], "result", "no result and no error means omit the field")
}

func TestStateDeltaReplace(t *testing.T) {
	op, err := StateDeltaReplace("/counter", 2)
	require.NoError(t, err)
	require.Equal(t, map[string]any{"op": "replace", "path": "/counter", "value": 2}, op)

	_, err = StateDeltaReplace("counter", 2)
	require.Error(t, err, "path without leading slash is invalid")
	_, err = StateDeltaReplace("", 2)
	require.Error(t, err, "empty path is invalid")
}
