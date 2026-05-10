package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

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

// TestAGUI_Integration_StreamingReasoner exercises the live-streaming
// path end to end: the reasoner returns NDJSON tagged events, the
// handler dispatches each into its AG-UI counterpart, frames are
// flushed live (verified by timestamping arrivals), and the run closes
// with MESSAGES_SNAPSHOT + RUN_FINISHED. This is the test that proves
// "Generative UI feels live" actually works under load — without it,
// any future regression that buffers the stream would silently make
// the UX stuttery again with no test failure.
func TestAGUI_Integration_StreamingReasoner(t *testing.T) {
	// The reasoner streams: text chunks (with deliberate per-chunk
	// delays so we can assert live forwarding), then a tool call, then
	// state, then closes.
	chunkDelay := 30 * time.Millisecond
	reasoner := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/reasoners/streaming-bot", r.URL.Path)
		w.Header().Set("Content-Type", "application/x-ndjson")
		w.WriteHeader(http.StatusOK)
		flusher, _ := w.(http.Flusher)

		send := func(line string) {
			fmt.Fprintln(w, line)
			if flusher != nil {
				flusher.Flush()
			}
			time.Sleep(chunkDelay)
		}
		send(`{"type":"reasoning","delta":"checking flights..."}`)
		send(`{"type":"reasoning","delta":" AA-12 wins on price."}`)
		send(`{"type":"text","delta":"Booked "}`)
		send(`{"type":"text","delta":"AA-12 SFO->JFK."}`)
		send(`{"type":"tool_call_start","id":"tc-1","name":"showFlightCard","arguments":{"from":"SFO","to":"JFK"}}`)
		send(`{"type":"tool_call_end","id":"tc-1"}`)
		send(`{"type":"state","snapshot":{"counter":1}}`)
		send(`{"type":"step_started","name":"finalize"}`)
		send(`{"type":"step_finished","name":"finalize"}`)
		send(`{"type":"custom","name":"telemetry","value":{"latency_ms":120}}`)
	}))
	defer reasoner.Close()

	store := &reasonerTestStorage{agent: &types.AgentNode{
		ID:              "stream-node",
		BaseURL:         reasoner.URL,
		HealthStatus:    types.HealthStatusActive,
		LifecycleStatus: types.AgentStatusReady,
		Reasoners:       []types.ReasonerDefinition{{ID: "streaming-bot"}},
	}}
	router := mountAGUIRouter(t, store)

	body := `{"threadId":"t","runId":"r","messages":[{"role":"user","content":"book it"}]}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/agui/runs/stream-node/streaming-bot", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code, w.Body.String())
	frames := parseAGUIStream(t, w.Body.String())
	got := []string{}
	for _, f := range frames {
		got = append(got, f.Type())
	}
	want := []string{
		"RUN_STARTED",
		"REASONING_START",
		"REASONING_MESSAGE_START",
		"REASONING_MESSAGE_CONTENT",
		"REASONING_MESSAGE_CONTENT",
		"REASONING_MESSAGE_END", // closed when text chunk arrives
		"REASONING_END",         // outer context closed
		"TEXT_MESSAGE_START",
		"TEXT_MESSAGE_CONTENT",
		"TEXT_MESSAGE_CONTENT",
		"TOOL_CALL_START",
		"TOOL_CALL_ARGS", // synthesized from `arguments` on start
		"TOOL_CALL_END",
		"STATE_SNAPSHOT",
		"STEP_STARTED",
		"STEP_FINISHED",
		"CUSTOM",
		"TEXT_MESSAGE_END", // closed at stream end
		"MESSAGES_SNAPSHOT",
		"RUN_FINISHED",
	}
	require.Equal(t, want, got, "streaming dispatcher diverged from canonical AG-UI ordering")

	// Each text-content delta must carry the chunk the reasoner sent
	// (proves the dispatcher didn't accidentally re-buffer).
	textDeltas := []string{}
	for _, f := range frames {
		if f.Type() == "TEXT_MESSAGE_CONTENT" {
			d, _ := f.Data["delta"].(string)
			textDeltas = append(textDeltas, d)
		}
	}
	require.Equal(t, []string{"Booked ", "AA-12 SFO->JFK."}, textDeltas)

	// MESSAGES_SNAPSHOT closes with the assistant turn carrying the
	// concatenated text and the tool call attached.
	snap, _ := frames[len(frames)-2].Data["messages"].([]any)
	require.Len(t, snap, 2)
	assistant, _ := snap[1].(map[string]any)
	require.Equal(t, "Booked AA-12 SFO->JFK.", assistant["content"])
	tcs, _ := assistant["toolCalls"].([]any)
	require.Len(t, tcs, 1)
}

// TestAGUI_Integration_StreamingErrorChunkTerminates: an `error` chunk
// from the reasoner terminates the stream with RUN_ERROR, even
// mid-flight, without emitting MESSAGES_SNAPSHOT or RUN_FINISHED.
func TestAGUI_Integration_StreamingErrorChunkTerminates(t *testing.T) {
	reasoner := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/x-ndjson")
		w.WriteHeader(http.StatusOK)
		flusher, _ := w.(http.Flusher)
		send := func(s string) {
			fmt.Fprintln(w, s)
			if flusher != nil {
				flusher.Flush()
			}
		}
		send(`{"type":"text","delta":"hello"}`)
		send(`{"type":"error","message":"upstream blew up","code":"ERR_LLM"}`)
		// Anything after the error must be ignored.
		send(`{"type":"text","delta":"unreachable"}`)
	}))
	defer reasoner.Close()

	store := &reasonerTestStorage{agent: &types.AgentNode{
		ID:              "n",
		BaseURL:         reasoner.URL,
		HealthStatus:    types.HealthStatusActive,
		LifecycleStatus: types.AgentStatusReady,
		Reasoners:       []types.ReasonerDefinition{{ID: "boom"}},
	}}
	router := mountAGUIRouter(t, store)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/agui/runs/n/boom",
		strings.NewReader(`{"threadId":"t","runId":"r","messages":[{"role":"user","content":"x"}]}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	frames := parseAGUIStream(t, w.Body.String())
	got := []string{}
	for _, f := range frames {
		got = append(got, f.Type())
	}
	// We accept the partial text frames, then RUN_ERROR (terminal).
	require.Contains(t, got, "RUN_ERROR")
	last := frames[len(frames)-1]
	require.Equal(t, "RUN_ERROR", last.Type())
	require.Equal(t, "upstream blew up", last.Data["message"])
	require.Equal(t, "ERR_LLM", last.Data["code"])
	require.NotContains(t, got, "MESSAGES_SNAPSHOT", "no snapshot after error")
	require.NotContains(t, got, "RUN_FINISHED", "no finish after error")
	// The post-error text chunk must have been dropped.
	for _, f := range frames {
		if f.Type() == "TEXT_MESSAGE_CONTENT" {
			d, _ := f.Data["delta"].(string)
			require.NotEqual(t, "unreachable", d, "post-error chunk must not leak through")
		}
	}
}

// TestAGUI_Integration_StreamingMalformedLineSurfacesAsRaw: a single bad
// NDJSON line shouldn't kill the stream — the dispatcher should surface
// it as RAW and continue.
func TestAGUI_Integration_StreamingMalformedLineSurfacesAsRaw(t *testing.T) {
	reasoner := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/x-ndjson")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, `{"type":"text","delta":"hi"}`)
		fmt.Fprintln(w, `{not valid json`)
		fmt.Fprintln(w, `{"type":"text","delta":" world"}`)
	}))
	defer reasoner.Close()

	store := &reasonerTestStorage{agent: &types.AgentNode{
		ID:              "n",
		BaseURL:         reasoner.URL,
		HealthStatus:    types.HealthStatusActive,
		LifecycleStatus: types.AgentStatusReady,
		Reasoners:       []types.ReasonerDefinition{{ID: "wobble"}},
	}}
	router := mountAGUIRouter(t, store)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/agui/runs/n/wobble",
		strings.NewReader(`{"threadId":"t","runId":"r","messages":[{"role":"user","content":"x"}]}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	frames := parseAGUIStream(t, w.Body.String())
	got := []string{}
	for _, f := range frames {
		got = append(got, f.Type())
	}
	require.Contains(t, got, "RAW", "malformed chunk should surface as RAW")
	// Stream completed; both text deltas reached us.
	textDeltas := []string{}
	for _, f := range frames {
		if f.Type() == "TEXT_MESSAGE_CONTENT" {
			d, _ := f.Data["delta"].(string)
			textDeltas = append(textDeltas, d)
		}
	}
	require.Equal(t, []string{"hi", " world"}, textDeltas)
	require.Equal(t, "RUN_FINISHED", frames[len(frames)-1].Type())
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
