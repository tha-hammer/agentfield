package handlers

import (
	"encoding/json"
	"testing"

	"github.com/Agent-Field/agentfield/control-plane/internal/agui"

	"github.com/stretchr/testify/require"
)

// captureWriter returns a writer fn that records every emitted event,
// optionally returning false on the Nth write to exercise short-circuit paths.
func captureWriter(failOn int) (writer func(agui.Event) bool, events *[]agui.Event) {
	collected := make([]agui.Event, 0, 16)
	count := 0
	return func(ev agui.Event) bool {
		count++
		collected = append(collected, ev)
		if failOn > 0 && count >= failOn {
			return false
		}
		return true
	}, &collected
}

func eventTypes(events []agui.Event) []string {
	out := make([]string, 0, len(events))
	for _, ev := range events {
		out = append(out, ev.Type())
	}
	return out
}

// TestDispatchChunk_GuardEarlyReturns walks the early-return guards on
// every chunk type that has one. None of these should write any events
// or short-circuit the loop.
func TestDispatchChunk_GuardEarlyReturns(t *testing.T) {
	cases := []struct {
		name string
		ch   streamingChunk
	}{
		{"empty text delta", streamingChunk{Type: "text"}},
		{"empty reasoning delta", streamingChunk{Type: "reasoning"}},
		{"tool_call_start missing id", streamingChunk{Type: "tool_call_start", Name: "x"}},
		{"tool_call_start missing name", streamingChunk{Type: "tool_call_start", ID: "tc1"}},
		{"tool_call_args missing id", streamingChunk{Type: "tool_call_args", Delta: "x"}},
		{"tool_call_args missing delta", streamingChunk{Type: "tool_call_args", ID: "tc1"}},
		{"tool_call_end missing id", streamingChunk{Type: "tool_call_end"}},
		{"tool_call_result missing id", streamingChunk{Type: "tool_call_result", Content: "x"}},
		{"state_delta empty ops", streamingChunk{Type: "state_delta", Ops: nil}},
		{"step_started missing name", streamingChunk{Type: "step_started"}},
		{"step_finished missing name", streamingChunk{Type: "step_finished"}},
		{"custom missing name", streamingChunk{Type: "custom", Value: 1}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			write, events := captureWriter(0)
			st := &streamingState{messageID: "msg-1"}
			require.True(t, dispatchChunk(write, st, tc.ch), "guard branch must keep stream alive")
			require.Empty(t, *events, "guard branch must not emit events")
		})
	}
}

// TestDispatchChunk_ToolCallLifecycle covers tool_call_start (with and
// without inline arguments), tool_call_args appending to an in-flight
// call, tool_call_end, and tool_call_result with both default and
// explicit role.
func TestDispatchChunk_ToolCallLifecycle(t *testing.T) {
	write, events := captureWriter(0)
	st := &streamingState{messageID: "msg-1"}

	require.True(t, dispatchChunk(write, st, streamingChunk{
		Type:      "tool_call_start",
		ID:        "tc1",
		Name:      "showFlightCard",
		Arguments: json.RawMessage(`{"flight":"AA-12"}`),
	}))
	// Start without inline args (a parent message should default to st.messageID).
	require.True(t, dispatchChunk(write, st, streamingChunk{Type: "tool_call_start", ID: "tc2", Name: "ping"}))
	require.True(t, dispatchChunk(write, st, streamingChunk{Type: "tool_call_args", ID: "tc2", Delta: `{"x":1}`}))
	require.True(t, dispatchChunk(write, st, streamingChunk{Type: "tool_call_end", ID: "tc2"}))
	require.True(t, dispatchChunk(write, st, streamingChunk{Type: "tool_call_result", ID: "tc2", Content: "done", Role: "system"}))
	require.True(t, dispatchChunk(write, st, streamingChunk{Type: "tool_call_result", ID: "tc1", Content: "ok"}))

	require.Equal(t, []string{
		"TOOL_CALL_START",  // tc1
		"TOOL_CALL_ARGS",   // tc1 inline args
		"TOOL_CALL_START",  // tc2
		"TOOL_CALL_ARGS",   // tc2 streamed delta
		"TOOL_CALL_END",    // tc2
		"TOOL_CALL_RESULT", // tc2 explicit role
		"TOOL_CALL_RESULT", // tc1 default role
	}, eventTypes(*events))

	require.Len(t, st.toolCalls, 2)
	require.Equal(t, `{"x":1}`, st.toolCalls[1].Function.Arguments,
		"tool_call_args should append to the in-flight call's arguments")
	tcResultExplicit := (*events)[5].(agui.ToolCallResult)
	require.Equal(t, "system", tcResultExplicit.Role)
	tcResultDefault := (*events)[6].(agui.ToolCallResult)
	require.Equal(t, "tool", tcResultDefault.Role)
}

// TestDispatchChunk_StateAndSteps covers state, state_delta, step_started,
// step_finished, raw, and custom chunks on the happy path.
func TestDispatchChunk_StateAndSteps(t *testing.T) {
	write, events := captureWriter(0)
	st := &streamingState{messageID: "msg-1"}

	require.True(t, dispatchChunk(write, st, streamingChunk{Type: "state", Snapshot: map[string]any{"k": 1}}))
	require.True(t, st.stateSet)
	require.True(t, dispatchChunk(write, st, streamingChunk{Type: "state_delta", Ops: []any{
		map[string]any{"op": "replace", "path": "/k", "value": 2},
	}}))
	require.True(t, dispatchChunk(write, st, streamingChunk{Type: "step_started", Name: "plan"}))
	require.True(t, dispatchChunk(write, st, streamingChunk{Type: "step_finished", Name: "plan"}))
	require.True(t, dispatchChunk(write, st, streamingChunk{Type: "raw", Event: map[string]any{"k": 1}, Source: "ext"}))
	require.True(t, dispatchChunk(write, st, streamingChunk{Type: "custom", Name: "ack", Value: true}))
	require.True(t, dispatchChunk(write, st, streamingChunk{Type: "unknown_kind"}))

	require.Equal(t, []string{
		"STATE_SNAPSHOT",
		"STATE_DELTA",
		"STEP_STARTED",
		"STEP_FINISHED",
		"RAW",
		"CUSTOM",
		"RAW", // unknown chunk falls into default → emits RAW
	}, eventTypes(*events))
}

// TestDispatchChunk_ErrorChunkTerminates verifies the error chunk emits
// RUN_ERROR and returns false to short-circuit the dispatch loop.
func TestDispatchChunk_ErrorChunkTerminates(t *testing.T) {
	write, events := captureWriter(0)
	st := &streamingState{messageID: "msg-1"}

	require.False(t, dispatchChunk(write, st, streamingChunk{
		Type:    "error",
		Message: "boom",
		Code:    "E_BOOM",
	}), "error chunk must short-circuit the dispatch loop")
	require.Equal(t, []string{"RUN_ERROR"}, eventTypes(*events))

	runErr := (*events)[0].(agui.RunError)
	require.Equal(t, "boom", runErr.Message)
	require.Equal(t, "E_BOOM", runErr.Code)
}

// TestDispatchChunk_ReasoningEndIdempotent confirms reasoning_end is a
// no-op when no reasoning segment is open and emits the End frame when
// one is.
func TestDispatchChunk_ReasoningEndIdempotent(t *testing.T) {
	write, events := captureWriter(0)
	st := &streamingState{messageID: "msg-1"}

	require.True(t, dispatchChunk(write, st, streamingChunk{Type: "reasoning_end"}))
	require.Empty(t, *events, "reasoning_end is a no-op without an open segment")

	require.True(t, dispatchChunk(write, st, streamingChunk{Type: "reasoning", Delta: "thinking..."}))
	require.True(t, dispatchChunk(write, st, streamingChunk{Type: "reasoning_end"}))
	require.Equal(t, []string{
		"REASONING_START",
		"REASONING_MESSAGE_START",
		"REASONING_MESSAGE_CONTENT",
		"REASONING_MESSAGE_END",
	}, eventTypes(*events))
	require.Empty(t, st.reasoningSeg, "reasoning_end clears the open segment id")
	require.NotEmpty(t, st.reasoningCtx, "reasoning_end leaves the outer context open")
}

// TestDispatchChunk_TextClosesReasoning ensures a text chunk closes any
// open reasoning session before opening the assistant text turn.
func TestDispatchChunk_TextClosesReasoning(t *testing.T) {
	write, events := captureWriter(0)
	st := &streamingState{messageID: "msg-1"}

	require.True(t, dispatchChunk(write, st, streamingChunk{Type: "reasoning", Delta: "thought"}))
	require.True(t, dispatchChunk(write, st, streamingChunk{Type: "text", Delta: "hello"}))

	types := eventTypes(*events)
	// REASONING_MESSAGE_END + REASONING_END must precede TEXT_MESSAGE_START.
	require.Contains(t, types, "REASONING_MESSAGE_END")
	require.Contains(t, types, "REASONING_END")
	require.Contains(t, types, "TEXT_MESSAGE_START")
	require.Contains(t, types, "TEXT_MESSAGE_CONTENT")
	require.Empty(t, st.reasoningCtx)
	require.Empty(t, st.reasoningSeg)
	require.True(t, st.textOpen)
}

// TestApplyFinal_FullEnvelope drives applyFinal with reasoning,
// toolCalls (with and without result), state, stateDelta, and result
// fields all populated.
func TestApplyFinal_FullEnvelope(t *testing.T) {
	write, events := captureWriter(0)
	st := &streamingState{messageID: "msg-1"}

	applyFinal(write, st, map[string]any{
		"reasoning": []any{"step 1", map[string]any{"content": "step 2", "id": "r-1"}},
		"toolCalls": []any{
			map[string]any{"id": "tc1", "name": "x", "arguments": map[string]any{"a": 1}, "result": "ok"},
			map[string]any{"id": "tc2", "name": "y", "arguments": map[string]any{}},
		},
		"state":      map[string]any{"counter": 7},
		"stateDelta": []any{map[string]any{"op": "replace", "path": "/counter", "value": 8}},
		"result":     "Done.",
	})

	types := eventTypes(*events)
	require.Contains(t, types, "REASONING_START")
	require.Contains(t, types, "REASONING_MESSAGE_START")
	require.Contains(t, types, "REASONING_MESSAGE_CONTENT")
	require.Contains(t, types, "REASONING_MESSAGE_END")
	require.Contains(t, types, "TOOL_CALL_START")
	require.Contains(t, types, "TOOL_CALL_ARGS")
	require.Contains(t, types, "TOOL_CALL_END")
	require.Contains(t, types, "TOOL_CALL_RESULT")
	require.Contains(t, types, "STATE_SNAPSHOT")
	require.Contains(t, types, "STATE_DELTA")
	require.Contains(t, types, "TEXT_MESSAGE_START")
	require.Contains(t, types, "TEXT_MESSAGE_CONTENT")

	require.True(t, st.textOpen, "final result text leaves the text session open for stream-end to close")
	require.True(t, st.stateSet)
	require.Len(t, st.toolCalls, 2)
}

// TestApplyFinal_NilDataIsNoOp confirms a nil data map is silently
// dropped — the reasoner can emit a final chunk without any structured
// fields.
func TestApplyFinal_NilDataIsNoOp(t *testing.T) {
	write, events := captureWriter(0)
	st := &streamingState{messageID: "msg-1"}
	applyFinal(write, st, nil)
	require.Empty(t, *events)
	require.False(t, st.textOpen)
}

// TestApplyFinal_ReusesOpenReasoningContext verifies that when a
// reasoning context is already open, applyFinal appends segments inside
// it instead of opening a new outer context.
func TestApplyFinal_ReusesOpenReasoningContext(t *testing.T) {
	write, events := captureWriter(0)
	st := &streamingState{messageID: "msg-1"}

	require.True(t, dispatchChunk(write, st, streamingChunk{Type: "reasoning", Delta: "first"}))
	priorReasoningStarts := 0
	for _, ev := range *events {
		if ev.Type() == "REASONING_START" {
			priorReasoningStarts++
		}
	}
	require.Equal(t, 1, priorReasoningStarts)

	applyFinal(write, st, map[string]any{"reasoning": []any{"another"}})

	totalStarts := 0
	for _, ev := range *events {
		if ev.Type() == "REASONING_START" {
			totalStarts++
		}
	}
	require.Equal(t, 1, totalStarts, "applyFinal must reuse the already-open reasoning context")
}

// TestCloseSessions_NoOpWhenIdle covers the early-return branches in
// closeTextSession and closeReasoningSession when no session is open.
func TestCloseSessions_NoOpWhenIdle(t *testing.T) {
	write, events := captureWriter(0)
	st := &streamingState{messageID: "msg-1"}
	require.True(t, closeTextSession(write, st))
	require.True(t, closeReasoningSession(write, st))
	require.Empty(t, *events)
}

// TestCloseSessions_WriteFailureShortCircuits covers the rare case where
// the writer returns false mid-close (client disconnect): the close
// helpers must propagate the failure so the dispatch loop can stop.
func TestCloseSessions_WriteFailureShortCircuits(t *testing.T) {
	st := &streamingState{
		messageID:    "msg-1",
		textOpen:     true,
		reasoningSeg: "seg-1",
		reasoningCtx: "ctx-1",
	}
	failingWrite := func(agui.Event) bool { return false }
	require.False(t, closeTextSession(failingWrite, st))
	require.False(t, closeReasoningSession(failingWrite, st))
}
