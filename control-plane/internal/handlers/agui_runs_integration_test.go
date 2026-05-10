package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Agent-Field/agentfield/control-plane/pkg/types"

	"github.com/stretchr/testify/require"
)

// agentField-style reasoner stub: mimics the wire shape an AgentField
// Python or Go SDK reasoner produces — JSON object with at least a
// `result` field, optionally `toolCalls` / `state` / `stateDelta` fields
// — so this test guards against the integration contract drifting.
//
// Without this, a future SDK rename of `prompt` -> `userPrompt` (or any
// similar tweak) would silently break Generative UI / shared state
// without failing any unit test, because the unit tests stub the
// agentInvoker interface and never inspect the reasoner-side input
// shape.

// TestAGUI_Integration_FullSequence runs the full AG-UI handler against
// a live httptest reasoner that returns the same shape a real .ai()
// reasoner would when authors use agentfield.agui helpers. Asserts:
//
//   - the reasoner received the canonical AG-UI envelope (prompt,
//     messages, tools, state, context, threadId, runId)
//   - the SSE stream carries lifecycle + tool calls (with TOOL_CALL_RESULT
//     for executed traces) + state snapshot + state delta + chunked text +
//     messages snapshot, in canonical order
//   - the assistant turn in MESSAGES_SNAPSHOT carries the tool calls
//     stitched onto it
func TestAGUI_Integration_FullSequence(t *testing.T) {
	var seenInput map[string]any
	reasoner := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/reasoners/integ", r.URL.Path)
		require.Equal(t, "application/json", r.Header.Get("Content-Type"))
		raw, _ := io.ReadAll(r.Body)
		require.NoError(t, json.Unmarshal(raw, &seenInput))

		// Mimic an SDK reasoner that used app.ai(tools=...) and returned
		// the trace via agentfield.agui.tool_calls_from_trace, plus a
		// fresh state and a single delta op.
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"result": "Booked SFO to JFK. Counter is now 2.",
			"toolCalls": [{
				"id": "tc-trace-0",
				"name": "showFlightCard",
				"arguments": {"from":"SFO","to":"JFK"},
				"result": {"flightId":"AA-12","status":"booked"}
			}],
			"state": {"counter": 2, "lastBooking": "AA-12"},
			"stateDelta": [
				{"op":"replace","path":"/counter","value":2}
			]
		}`))
	}))
	defer reasoner.Close()

	store := &reasonerTestStorage{agent: &types.AgentNode{
		ID:              "integ-node",
		BaseURL:         reasoner.URL,
		HealthStatus:    types.HealthStatusActive,
		LifecycleStatus: types.AgentStatusReady,
		Reasoners:       []types.ReasonerDefinition{{ID: "integ"}},
	}}
	router := mountAGUIRouter(t, store)

	// Build a canonical RunAgentInput that exercises every surface the
	// reasoner is supposed to receive: prompt + multi-message history +
	// tools + state + context + forwardedProps.
	body := `{
		"threadId": "thread-int", "runId": "run-int",
		"messages": [
			{"role":"system","content":"you are helpful"},
			{"role":"user","content":"book SFO->JFK"}
		],
		"tools": [{"name":"showFlightCard","description":"render a flight card"}],
		"context": [{"description":"user prefs","value":{"seat":"aisle"}}],
		"state": {"counter": 1},
		"forwardedProps": {"locale":"en-US"}
	}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/agui/runs/integ-node/integ", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code, w.Body.String())

	// 1. Reasoner saw the canonical envelope.
	require.Equal(t, "book SFO->JFK", seenInput["prompt"])
	require.Equal(t, "thread-int", seenInput["threadId"])
	require.Equal(t, "run-int", seenInput["runId"])
	gotMessages, _ := seenInput["messages"].([]any)
	require.Len(t, gotMessages, 2)
	gotTools, _ := seenInput["tools"].([]any)
	require.Len(t, gotTools, 1)
	gotContext, _ := seenInput["context"].([]any)
	require.Len(t, gotContext, 1)
	gotState, _ := seenInput["state"].(map[string]any)
	require.EqualValues(t, 1, gotState["counter"])
	gotFP, _ := seenInput["forwardedProps"].(map[string]any)
	require.Equal(t, "en-US", gotFP["locale"])

	// 2. Wire output: full canonical sequence.
	frames := parseAGUIStream(t, w.Body.String())
	types := []string{}
	for _, f := range frames {
		types = append(types, f.Type())
	}
	want := []string{
		"RUN_STARTED",
		"TOOL_CALL_START",
		"TOOL_CALL_ARGS",
		"TOOL_CALL_END",
		"TOOL_CALL_RESULT",
		"TEXT_MESSAGE_START",
		"TEXT_MESSAGE_CONTENT",
		"TEXT_MESSAGE_END",
		"STATE_SNAPSHOT",
		"STATE_DELTA",
		"MESSAGES_SNAPSHOT",
		"RUN_FINISHED",
	}
	require.Equal(t, want, types, "frame sequence diverged from canonical AG-UI order")

	// 3. TOOL_CALL_RESULT carries the executed trace's result.
	resFrame := frames[4]
	require.Equal(t, "tc-trace-0", resFrame.Data["toolCallId"])
	require.Equal(t, "tool", resFrame.Data["role"])
	require.JSONEq(t, `{"flightId":"AA-12","status":"booked"}`, resFrame.Data["content"].(string))

	// 4. STATE_SNAPSHOT carries new value; STATE_DELTA carries the patch.
	snap, _ := frames[8].Data["snapshot"].(map[string]any)
	require.EqualValues(t, 2, snap["counter"])
	require.Equal(t, "AA-12", snap["lastBooking"])
	delta, _ := frames[9].Data["delta"].([]any)
	require.Len(t, delta, 1)
	op, _ := delta[0].(map[string]any)
	require.Equal(t, "replace", op["op"])

	// 5. MESSAGES_SNAPSHOT — assistant turn carries tool calls.
	msgs, _ := frames[10].Data["messages"].([]any)
	require.Len(t, msgs, 3, "should be 2 inbound + 1 new assistant")
	assistant, _ := msgs[2].(map[string]any)
	require.Equal(t, "assistant", assistant["role"])
	tcs, _ := assistant["toolCalls"].([]any)
	require.Len(t, tcs, 1)
	tc, _ := tcs[0].(map[string]any)
	require.Equal(t, "tc-trace-0", tc["id"])
	fn, _ := tc["function"].(map[string]any)
	require.Equal(t, "showFlightCard", fn["name"])
}

// TestAGUI_Integration_FollowupTurnWithToolMessage verifies the second
// half of the CopilotKit "user clicked confirm" loop: when the next
// run's inbound history includes a role:"tool" message, the reasoner
// receives it intact so it can produce a follow-up response.
func TestAGUI_Integration_FollowupTurnWithToolMessage(t *testing.T) {
	var seenInput map[string]any
	reasoner := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		raw, _ := io.ReadAll(r.Body)
		require.NoError(t, json.Unmarshal(raw, &seenInput))
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"result":"Booking confirmed."}`))
	}))
	defer reasoner.Close()

	store := &reasonerTestStorage{agent: &types.AgentNode{
		ID:              "n",
		BaseURL:         reasoner.URL,
		HealthStatus:    types.HealthStatusActive,
		LifecycleStatus: types.AgentStatusReady,
		Reasoners:       []types.ReasonerDefinition{{ID: "f"}},
	}}
	router := mountAGUIRouter(t, store)

	body := `{
		"threadId":"t","runId":"r2",
		"messages":[
			{"role":"user","content":"book SFO->JFK"},
			{"role":"assistant","toolCalls":[{
				"id":"tc1","type":"function",
				"function":{"name":"showFlightCard","arguments":"{\"from\":\"SFO\"}"}
			}]},
			{"role":"tool","toolCallId":"tc1","content":"user confirmed"},
			{"role":"user","content":"now book the return"}
		]
	}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/agui/runs/n/f", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code, w.Body.String())
	// The reasoner must see the tool message verbatim — that's what
	// closes the click-confirm loop. Without this, the agent has no way
	// of knowing the tool ran.
	require.Equal(t, "now book the return", seenInput["prompt"])
	msgs, _ := seenInput["messages"].([]any)
	require.Len(t, msgs, 4)
	tool, _ := msgs[2].(map[string]any)
	require.Equal(t, "tool", tool["role"])
	require.Equal(t, "tc1", tool["toolCallId"])
	require.Equal(t, "user confirmed", tool["content"])
}
