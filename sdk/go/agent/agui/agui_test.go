package agui

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"strings"
	"testing"
	"time"

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

func TestStreamingChunkBuilders(t *testing.T) {
	require.Equal(t, map[string]any{"type": "text", "delta": "hi"}, TextChunk("hi"))
	require.Equal(t, map[string]any{"type": "reasoning", "delta": "think"}, ReasoningChunk("think"))
	require.Equal(t, map[string]any{"type": "reasoning_end"}, ReasoningEndChunk())

	tcStart := ToolCallStartChunk("tc1", "x", map[string]any{"a": 1}, "msg-1")
	require.Equal(t, "tool_call_start", tcStart["type"])
	require.Equal(t, "tc1", tcStart["id"])
	require.Equal(t, "x", tcStart["name"])
	require.Equal(t, map[string]any{"a": 1}, tcStart["arguments"])
	require.Equal(t, "msg-1", tcStart["parentMessageId"])

	tcStartNoExtras := ToolCallStartChunk("tc2", "x", nil, "")
	require.NotContains(t, tcStartNoExtras, "arguments")
	require.NotContains(t, tcStartNoExtras, "parentMessageId")

	require.Equal(t, map[string]any{"type": "tool_call_args", "id": "tc1", "delta": "{\"x"}, ToolCallArgsChunk("tc1", "{\"x"))
	require.Equal(t, map[string]any{"type": "tool_call_end", "id": "tc1"}, ToolCallEndChunk("tc1"))

	res := ToolCallResultChunk("tc1", "ok", "")
	require.Equal(t, "tool", res["role"], "default role is 'tool'")
	require.Equal(t, "ok", res["content"])

	require.Equal(t, map[string]any{"type": "state", "snapshot": map[string]any{"a": 1}}, StateChunk(map[string]any{"a": 1}))
	require.Equal(t, "state_delta", StateDeltaChunk([]any{map[string]any{"op": "replace"}})["type"])

	require.Equal(t, "step_started", StepStartedChunk("plan")["type"])
	require.Equal(t, "step_finished", StepFinishedChunk("plan")["type"])

	raw := RawChunk(map[string]any{"x": 1}, "harness")
	require.Equal(t, "raw", raw["type"])
	require.Equal(t, "harness", raw["source"])

	rawNoSrc := RawChunk(map[string]any{"x": 1}, "")
	require.NotContains(t, rawNoSrc, "source")

	custom := CustomChunk("ack", map[string]any{"ok": true})
	require.Equal(t, "custom", custom["type"])
	require.Equal(t, "ack", custom["name"])

	customNil := CustomChunk("ack", nil)
	require.NotContains(t, customNil, "value")

	final := FinalChunk(map[string]any{"toolCalls": []any{}})
	require.Equal(t, "final", final["type"])

	errCh := ErrorChunk("boom", "E1")
	require.Equal(t, "error", errCh["type"])
	require.Equal(t, "boom", errCh["message"])
	require.Equal(t, "E1", errCh["code"])

	errChNoCode := ErrorChunk("boom", "")
	require.NotContains(t, errChNoCode, "code")
}

func TestSerializeStream(t *testing.T) {
	ch := make(chan map[string]any, 4)
	ch <- TextChunk("hello ")
	ch <- TextChunk("world")
	ch <- StateChunk(map[string]any{"counter": 1})
	close(ch)

	var buf bytes.Buffer
	require.NoError(t, SerializeStream(context.Background(), &buf, ch))
	lines := strings.Split(strings.TrimRight(buf.String(), "\n"), "\n")
	require.Len(t, lines, 3)

	var first map[string]any
	require.NoError(t, json.Unmarshal([]byte(lines[0]), &first))
	require.Equal(t, "text", first["type"])
	require.Equal(t, "hello ", first["delta"])

	var third map[string]any
	require.NoError(t, json.Unmarshal([]byte(lines[2]), &third))
	require.Equal(t, "state", third["type"])
}

func TestSerializeStream_RespectsContext(t *testing.T) {
	ch := make(chan map[string]any) // never closed; never sends
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() { done <- SerializeStream(ctx, io.Discard, ch) }()
	cancel()
	select {
	case err := <-done:
		require.ErrorIs(t, err, context.Canceled)
	case <-time.After(time.Second):
		t.Fatal("SerializeStream did not honor context cancellation")
	}
}
