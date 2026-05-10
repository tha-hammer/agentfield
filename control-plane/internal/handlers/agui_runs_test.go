package handlers

import (
	"bufio"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/Agent-Field/agentfield/control-plane/pkg/types"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

// aguiFrame is a parsed SSE frame: just the JSON object decoded from the
// `data:` line. The canonical AG-UI encoder emits frames as `data: <json>\n\n`
// only — no `event:` line — so the JSON `type` field is the sole discriminator.
type aguiFrame struct {
	Data map[string]any
}

func (f aguiFrame) Type() string {
	t, _ := f.Data["type"].(string)
	return t
}

// parseAGUIStream splits an SSE response body into one frame per AG-UI event.
// Strict on shape: every frame must be `data: <json>\n\n`. We assert against
// the strictness because that's exactly what the AG-UI spec guarantees and
// what the reference encoders emit (see ag-ui-protocol/ag-ui encoder.ts /
// encoder.py).
func parseAGUIStream(t *testing.T, body string) []aguiFrame {
	t.Helper()
	var frames []aguiFrame
	scanner := bufio.NewScanner(strings.NewReader(body))
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)

	var curData string
	flush := func() {
		if curData == "" {
			return
		}
		var decoded map[string]any
		require.NoError(t, json.Unmarshal([]byte(curData), &decoded), "data line is not JSON: %s", curData)
		frames = append(frames, aguiFrame{Data: decoded})
		curData = ""
	}

	for scanner.Scan() {
		line := scanner.Text()
		switch {
		case line == "":
			flush()
		case strings.HasPrefix(line, "event:"):
			t.Fatalf("AG-UI frames must not include an `event:` line; got: %q", line)
		case strings.HasPrefix(line, "data: "):
			curData = strings.TrimPrefix(line, "data: ")
		}
	}
	flush()
	return frames
}

func mountAGUIRouter(t *testing.T, store *reasonerTestStorage) *gin.Engine {
	t.Helper()
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/api/v1/agui/runs/:node_id/:reasoner_name", AGUIRunHandler(store))
	return router
}

// runAgentInputBody returns a canonical RunAgentInputSchema-shaped body. The
// vanilla @ag-ui/client HttpAgent — and therefore CopilotKit's runtime that
// wraps it — POSTs exactly this shape. Tests should always go through this
// helper so the assertion about "we accept the canonical shape" is real.
func runAgentInputBody(t *testing.T, threadID, runID, prompt string) string {
	t.Helper()
	body := map[string]any{
		"threadId": threadID,
		"runId":    runID,
		"messages": []map[string]any{
			{"id": "u1", "role": "user", "content": prompt},
		},
		"tools":          []any{},
		"context":        []any{},
		"state":          map[string]any{},
		"forwardedProps": map[string]any{},
	}
	b, err := json.Marshal(body)
	require.NoError(t, err)
	return string(b)
}

// TestAGUIRunHandler_HappyPath_EmitsCanonicalEventSequence is the core
// assertion: a successful run produces RUN_STARTED → TEXT_MESSAGE_START →
// TEXT_MESSAGE_CONTENT → TEXT_MESSAGE_END → MESSAGES_SNAPSHOT → RUN_FINISHED,
// in that order. Thread/run IDs propagate from the request to RUN_FINISHED.
// The reasoner sees the AG-UI envelope (prompt extracted from the trailing
// user message) — proving the body-shape change wired up correctly.
func TestAGUIRunHandler_HappyPath_EmitsCanonicalEventSequence(t *testing.T) {
	var seenInput map[string]any
	agentServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/reasoners/echo", r.URL.Path)
		require.Equal(t, http.MethodPost, r.Method)
		raw, _ := io.ReadAll(r.Body)
		require.NoError(t, json.Unmarshal(raw, &seenInput))
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"result":"hello world"}`))
	}))
	defer agentServer.Close()

	store := &reasonerTestStorage{agent: &types.AgentNode{
		ID:              "node-1",
		BaseURL:         agentServer.URL,
		HealthStatus:    types.HealthStatusActive,
		LifecycleStatus: types.AgentStatusReady,
		Reasoners:       []types.ReasonerDefinition{{ID: "echo"}},
	}}
	router := mountAGUIRouter(t, store)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/agui/runs/node-1/echo",
		strings.NewReader(runAgentInputBody(t, "thread-test", "run-test", "hi")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code, "response: %s", w.Body.String())
	require.Equal(t, "text/event-stream", w.Header().Get("Content-Type"))

	// The reasoner received the canonical AG-UI envelope, plus the `prompt`
	// convenience extracted from the trailing user message.
	require.Equal(t, "hi", seenInput["prompt"])
	require.Equal(t, "thread-test", seenInput["threadId"])
	require.Equal(t, "run-test", seenInput["runId"])
	gotMessages, _ := seenInput["messages"].([]any)
	require.Len(t, gotMessages, 1)

	frames := parseAGUIStream(t, w.Body.String())
	wantSequence := []string{
		"RUN_STARTED",
		"TEXT_MESSAGE_START",
		"TEXT_MESSAGE_CONTENT",
		"TEXT_MESSAGE_END",
		"MESSAGES_SNAPSHOT",
		"RUN_FINISHED",
	}
	require.Len(t, frames, len(wantSequence), "frames: %+v", frames)
	for i, want := range wantSequence {
		require.Equal(t, want, frames[i].Type(), "frame %d: %v", i, frames[i].Data)
	}

	require.Equal(t, "thread-test", frames[0].Data["threadId"])
	require.Equal(t, "run-test", frames[0].Data["runId"])
	require.NotContains(t, frames[0].Data, "input",
		"input must be omitted; the spec types it as RunAgentInput, not freeform")

	msgID, _ := frames[1].Data["messageId"].(string)
	require.NotEmpty(t, msgID)
	require.Equal(t, "assistant", frames[1].Data["role"])
	require.Equal(t, msgID, frames[2].Data["messageId"])
	require.Equal(t, "hello world", frames[2].Data["delta"])
	require.Equal(t, msgID, frames[3].Data["messageId"])

	// MESSAGES_SNAPSHOT carries inbound history + the new assistant turn,
	// and the assistant's content matches the delta we emitted.
	snapMsgs, _ := frames[4].Data["messages"].([]any)
	require.Len(t, snapMsgs, 2, "snapshot should have 1 user + 1 assistant message")
	last, _ := snapMsgs[1].(map[string]any)
	require.Equal(t, "assistant", last["role"])
	require.Equal(t, "hello world", last["content"])
	require.Equal(t, msgID, last["id"])

	// RUN_FINISHED carries threadId/runId, success outcome, and the parsed
	// agent JSON.
	require.Equal(t, "thread-test", frames[5].Data["threadId"])
	require.Equal(t, "run-test", frames[5].Data["runId"])
	outcome, _ := frames[5].Data["outcome"].(map[string]any)
	require.Equal(t, "success", outcome["type"])
	require.Equal(t, map[string]any{"result": "hello world"}, frames[5].Data["result"])

	if ts, ok := frames[0].Data["timestamp"]; ok {
		_, isFloat := ts.(float64)
		require.True(t, isFloat, "timestamp must be a number, got %T", ts)
	}
}

// TestAGUIRunHandler_GeneratesIDsWhenAbsent confirms that omitted threadId
// and runId are auto-populated rather than left empty — clients shouldn't
// have to mint IDs themselves to get a valid AG-UI stream.
func TestAGUIRunHandler_GeneratesIDsWhenAbsent(t *testing.T) {
	agentServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"result":"ok"}`))
	}))
	defer agentServer.Close()

	store := &reasonerTestStorage{agent: &types.AgentNode{
		ID:              "node-1",
		BaseURL:         agentServer.URL,
		HealthStatus:    types.HealthStatusActive,
		LifecycleStatus: types.AgentStatusReady,
		Reasoners:       []types.ReasonerDefinition{{ID: "echo"}},
	}}
	router := mountAGUIRouter(t, store)

	// Omit threadId and runId — vanilla HttpAgent always sends them, but a
	// test client may not.
	req := httptest.NewRequest(http.MethodPost, "/api/v1/agui/runs/node-1/echo",
		strings.NewReader(`{"messages":[{"role":"user","content":"hi"}]}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code, w.Body.String())
	frames := parseAGUIStream(t, w.Body.String())
	require.NotEmpty(t, frames)
	require.Equal(t, "RUN_STARTED", frames[0].Type())
	threadID, _ := frames[0].Data["threadId"].(string)
	runID, _ := frames[0].Data["runId"].(string)
	require.NotEmpty(t, threadID, "threadId should be auto-generated")
	require.NotEmpty(t, runID, "runId should be auto-generated")

	last := frames[len(frames)-1]
	require.Equal(t, "RUN_FINISHED", last.Type())
	require.Equal(t, threadID, last.Data["threadId"])
	require.Equal(t, runID, last.Data["runId"])
}

// TestAGUIRunHandler_AgentFailureEmitsRunError confirms the streaming-side
// error path: once SSE is open, downstream agent failure must surface as a
// terminal RUN_ERROR frame, never as a partial happy-path-shaped sequence.
func TestAGUIRunHandler_AgentFailureEmitsRunError(t *testing.T) {
	agentServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"upstream blew up"}`))
	}))
	defer agentServer.Close()

	store := &reasonerTestStorage{agent: &types.AgentNode{
		ID:              "node-1",
		BaseURL:         agentServer.URL,
		HealthStatus:    types.HealthStatusActive,
		LifecycleStatus: types.AgentStatusReady,
		Reasoners:       []types.ReasonerDefinition{{ID: "boom"}},
	}}
	router := mountAGUIRouter(t, store)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/agui/runs/node-1/boom",
		strings.NewReader(runAgentInputBody(t, "t", "r", "x")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code, w.Body.String())
	frames := parseAGUIStream(t, w.Body.String())
	require.GreaterOrEqual(t, len(frames), 2)
	require.Equal(t, "RUN_STARTED", frames[0].Type())

	last := frames[len(frames)-1]
	require.Equal(t, "RUN_ERROR", last.Type())
	require.NotEmpty(t, last.Data["message"])
	require.Equal(t, "ERR_AGENT_CALL", last.Data["code"])

	for _, f := range frames[1:] {
		require.NotContains(t,
			[]string{"TEXT_MESSAGE_START", "TEXT_MESSAGE_CONTENT", "TEXT_MESSAGE_END", "MESSAGES_SNAPSHOT", "RUN_FINISHED"},
			f.Type(), "unexpected post-error frame: %s", f.Type())
	}
}

// TestAGUIRunHandler_EmitsHeartbeatWhileReasonerIsSlow confirms long-running
// reasoners produce SSE comment frames (`: keep-alive`) so proxies don't
// idle-time-out the connection. Comments are invisible to AG-UI clients but
// keep intermediaries happy.
func TestAGUIRunHandler_EmitsHeartbeatWhileReasonerIsSlow(t *testing.T) {
	prev := AGUIHeartbeatInterval
	AGUIHeartbeatInterval = 50 * time.Millisecond
	defer func() { AGUIHeartbeatInterval = prev }()

	agentServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		time.Sleep(250 * time.Millisecond)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"result":"finally"}`))
	}))
	defer agentServer.Close()

	store := &reasonerTestStorage{agent: &types.AgentNode{
		ID:              "node-1",
		BaseURL:         agentServer.URL,
		HealthStatus:    types.HealthStatusActive,
		LifecycleStatus: types.AgentStatusReady,
		Reasoners:       []types.ReasonerDefinition{{ID: "slow"}},
	}}
	router := mountAGUIRouter(t, store)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/agui/runs/node-1/slow",
		strings.NewReader(runAgentInputBody(t, "t", "r", "x")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code, w.Body.String())
	body := w.Body.String()
	require.Contains(t, body, ": keep-alive",
		"expected at least one SSE comment heartbeat in:\n%s", body)

	frames := parseAGUIStream(t, body)
	require.Equal(t, "RUN_STARTED", frames[0].Type())
	require.Equal(t, "RUN_FINISHED", frames[len(frames)-1].Type())
}

// TestAGUIRunHandler_AgentBodyWithoutResultKey covers the fallthrough in
// extractAssistantText: when the agent returns a JSON object that doesn't
// have `result` or `content`, internal-only keys (toolCalls, state) are
// stripped and the rest is JSON-encoded as the delta.
func TestAGUIRunHandler_AgentBodyWithoutResultKey(t *testing.T) {
	agentServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"status":"ok","count":3}`))
	}))
	defer agentServer.Close()

	store := &reasonerTestStorage{agent: &types.AgentNode{
		ID:              "node-1",
		BaseURL:         agentServer.URL,
		HealthStatus:    types.HealthStatusActive,
		LifecycleStatus: types.AgentStatusReady,
		Reasoners:       []types.ReasonerDefinition{{ID: "ping"}},
	}}
	router := mountAGUIRouter(t, store)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/agui/runs/node-1/ping",
		strings.NewReader(runAgentInputBody(t, "t", "r", "x")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code, w.Body.String())
	frames := parseAGUIStream(t, w.Body.String())
	// content frame is index 2 in the canonical sequence.
	require.Equal(t, `{"count":3,"status":"ok"}`, frames[2].Data["delta"])
}

// TestStringifyResult_BranchCoverage covers the cheap branches of the
// helper directly: string passthrough, nil, and arbitrary value JSON-encode.
func TestStringifyResult_BranchCoverage(t *testing.T) {
	require.Equal(t, "hello", stringifyResult("hello"))
	require.Equal(t, "", stringifyResult(nil))
	require.Equal(t, `[1,2,3]`, stringifyResult([]any{1, 2, 3}))
	require.Equal(t, `{"a":1}`, stringifyResult(map[string]any{"a": 1}))
}

// TestExtractAssistantText_AllBranches exercises the helper directly so
// every priority rung is covered: result key, content key, top-level
// string, top-level non-map non-string (number), filtered-empty map, and
// the non-JSON raw-body fallthrough.
func TestExtractAssistantText_AllBranches(t *testing.T) {
	require.Equal(t, "raw bytes", extractAssistantText(nil, false, []byte("raw bytes")),
		"non-JSON falls through to raw body")
	require.Equal(t, "answer", extractAssistantText(map[string]any{"result": "answer"}, true, nil),
		"`result` key wins")
	require.Equal(t, "alt", extractAssistantText(map[string]any{"content": "alt"}, true, nil),
		"`content` key is the second priority")
	require.Equal(t, "just-a-string", extractAssistantText("just-a-string", true, nil),
		"top-level JSON string passes through")
	require.Equal(t, "42", extractAssistantText(float64(42), true, nil),
		"top-level non-map non-string is JSON-encoded")
	require.Equal(t, "", extractAssistantText(map[string]any{"toolCalls": []any{}, "state": map[string]any{}}, true, nil),
		"a body containing only internal-only fields collapses to empty delta")
}

// TestExtractToolCalls_NonMapInput covers the non-map branch (e.g. the
// reasoner returned a top-level string or array — no toolCalls possible).
func TestExtractToolCalls_NonMapInput(t *testing.T) {
	require.Nil(t, extractToolCalls("just a string"))
	require.Nil(t, extractToolCalls([]any{1, 2, 3}))
	require.Nil(t, extractToolCalls(nil))
	// Map without a `toolCalls` array also returns nil.
	require.Nil(t, extractToolCalls(map[string]any{"result": "x"}))
}

// TestExtractState_NonMapAndAbsent covers both the non-map and the
// missing-key paths.
func TestExtractState_NonMapAndAbsent(t *testing.T) {
	_, ok := extractState("not a map")
	require.False(t, ok)
	_, ok = extractState(map[string]any{"result": "x"})
	require.False(t, ok, "absent state key returns ok=false")
	v, ok := extractState(map[string]any{"state": nil})
	require.True(t, ok, "explicit null state still returns ok=true")
	require.Nil(t, v)
}

// TestExtractStateDelta covers presence, non-map, and empty cases.
func TestExtractStateDelta(t *testing.T) {
	require.Nil(t, extractStateDelta("not a map"))
	require.Nil(t, extractStateDelta(map[string]any{}), "absent stateDelta key")
	require.Nil(t, extractStateDelta(map[string]any{"stateDelta": []any{}}),
		"empty stateDelta is treated as absent")
	d := extractStateDelta(map[string]any{"stateDelta": []any{
		map[string]any{"op": "replace", "path": "/x", "value": 1},
	}})
	require.Len(t, d, 1)
}

// TestChunkText covers the token-streaming chunker: rune boundaries,
// empty input, oversize input, exact boundary.
func TestChunkText(t *testing.T) {
	require.Equal(t, []string{""}, chunkText("", 4))
	require.Equal(t, []string{"abc"}, chunkText("abc", 4))
	require.Equal(t, []string{"abcd", "ef"}, chunkText("abcdef", 4))
	require.Equal(t, []string{"hello"}, chunkText("hello", -1), "non-positive size returns input unchanged")
	// Multi-byte runes (emoji) must split on rune boundaries.
	emoji := "🤖🤖🤖"
	chunks := chunkText(emoji, 4)
	for _, c := range chunks {
		require.Equal(t, "🤖", c, "each chunk should hold exactly one emoji at size=4")
	}
	require.Equal(t, 3, len(chunks))
}

// TestAGUIRunHandler_ToolCalls_EmitsResultEventForServerSideCalls covers
// the .ai(tools=...) trace surfacing path: when a reasoner reports a tool
// call as already-executed by including a `result` field, the handler
// emits TOOL_CALL_RESULT after TOOL_CALL_END.
func TestAGUIRunHandler_ToolCalls_EmitsResultEventForServerSideCalls(t *testing.T) {
	agentServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"result":"queried weather",
			"toolCalls":[{
				"id":"tc-w1","name":"getWeather",
				"arguments":{"city":"SF"},
				"result":{"temp":62,"summary":"foggy"}
			}]
		}`))
	}))
	defer agentServer.Close()

	store := &reasonerTestStorage{agent: &types.AgentNode{
		ID:              "node-1",
		BaseURL:         agentServer.URL,
		HealthStatus:    types.HealthStatusActive,
		LifecycleStatus: types.AgentStatusReady,
		Reasoners:       []types.ReasonerDefinition{{ID: "weather"}},
	}}
	router := mountAGUIRouter(t, store)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/agui/runs/node-1/weather",
		strings.NewReader(runAgentInputBody(t, "t", "r", "weather in SF?")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code, w.Body.String())
	frames := parseAGUIStream(t, w.Body.String())

	// TOOL_CALL_RESULT must come immediately after TOOL_CALL_END for the
	// same toolCallId.
	idx := func(typ string) int {
		for i, f := range frames {
			if f.Type() == typ {
				return i
			}
		}
		return -1
	}
	require.Less(t, idx("TOOL_CALL_END"), idx("TOOL_CALL_RESULT"),
		"TOOL_CALL_RESULT must follow TOOL_CALL_END")
	resFrame := frames[idx("TOOL_CALL_RESULT")]
	require.Equal(t, "tc-w1", resFrame.Data["toolCallId"])
	require.Equal(t, "tool", resFrame.Data["role"])
	require.JSONEq(t, `{"summary":"foggy","temp":62}`, resFrame.Data["content"].(string))
}

// TestAGUIRunHandler_StateDelta covers Tier 3's incremental-patch path:
// when the reasoner returns `stateDelta` (RFC 6902), STATE_DELTA is
// emitted alongside (or instead of) STATE_SNAPSHOT.
func TestAGUIRunHandler_StateDelta(t *testing.T) {
	agentServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"result":"updated",
			"state":{"counter":1},
			"stateDelta":[{"op":"replace","path":"/counter","value":2}]
		}`))
	}))
	defer agentServer.Close()

	store := &reasonerTestStorage{agent: &types.AgentNode{
		ID:              "node-1",
		BaseURL:         agentServer.URL,
		HealthStatus:    types.HealthStatusActive,
		LifecycleStatus: types.AgentStatusReady,
		Reasoners:       []types.ReasonerDefinition{{ID: "tick"}},
	}}
	router := mountAGUIRouter(t, store)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/agui/runs/node-1/tick",
		strings.NewReader(runAgentInputBody(t, "t", "r", "tick")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	frames := parseAGUIStream(t, w.Body.String())

	// Both forms emitted, snapshot first.
	idx := func(typ string) int {
		for i, f := range frames {
			if f.Type() == typ {
				return i
			}
		}
		return -1
	}
	require.NotEqual(t, -1, idx("STATE_SNAPSHOT"))
	require.NotEqual(t, -1, idx("STATE_DELTA"))
	require.Less(t, idx("STATE_SNAPSHOT"), idx("STATE_DELTA"))
	delta, _ := frames[idx("STATE_DELTA")].Data["delta"].([]any)
	require.Len(t, delta, 1)
	op, _ := delta[0].(map[string]any)
	require.Equal(t, "replace", op["op"])
	require.Equal(t, "/counter", op["path"])
}

// TestAGUIRunHandler_ChunkedTextStreaming verifies that long assistant
// replies are split across multiple TEXT_MESSAGE_CONTENT deltas (so the
// frontend can paint progressively) while the start/end frames stay
// singletons.
func TestAGUIRunHandler_ChunkedTextStreaming(t *testing.T) {
	prev := AGUITextChunkSize
	AGUITextChunkSize = 8 // tiny chunks so we can assert multi-frame easily
	defer func() { AGUITextChunkSize = prev }()

	long := strings.Repeat("a", 25) // 25 / 8 = 4 chunks (8+8+8+1)
	agentServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"result":"` + long + `"}`))
	}))
	defer agentServer.Close()

	store := &reasonerTestStorage{agent: &types.AgentNode{
		ID:              "node-1",
		BaseURL:         agentServer.URL,
		HealthStatus:    types.HealthStatusActive,
		LifecycleStatus: types.AgentStatusReady,
		Reasoners:       []types.ReasonerDefinition{{ID: "long"}},
	}}
	router := mountAGUIRouter(t, store)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/agui/runs/node-1/long",
		strings.NewReader(runAgentInputBody(t, "t", "r", "x")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	frames := parseAGUIStream(t, w.Body.String())

	starts, contents, ends := 0, 0, 0
	concatenated := ""
	var msgID string
	for _, f := range frames {
		switch f.Type() {
		case "TEXT_MESSAGE_START":
			starts++
			msgID, _ = f.Data["messageId"].(string)
		case "TEXT_MESSAGE_CONTENT":
			contents++
			require.Equal(t, msgID, f.Data["messageId"], "all content frames must share the same messageId")
			d, _ := f.Data["delta"].(string)
			concatenated += d
		case "TEXT_MESSAGE_END":
			ends++
		}
	}
	require.Equal(t, 1, starts, "exactly one START frame")
	require.Equal(t, 1, ends, "exactly one END frame")
	require.GreaterOrEqual(t, contents, 4, "expected long reply to be split into ≥4 chunks (got %d)", contents)
	require.Equal(t, long, concatenated, "concatenated deltas must equal the full reply")
}

// TestAGUIRunHandler_AgentReturnsNonJSON falls through to the
// `string(body)` branch when the agent's response isn't valid JSON.
func TestAGUIRunHandler_AgentReturnsNonJSON(t *testing.T) {
	agentServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		_, _ = w.Write([]byte(`plain text answer`))
	}))
	defer agentServer.Close()

	store := &reasonerTestStorage{agent: &types.AgentNode{
		ID:              "node-1",
		BaseURL:         agentServer.URL,
		HealthStatus:    types.HealthStatusActive,
		LifecycleStatus: types.AgentStatusReady,
		Reasoners:       []types.ReasonerDefinition{{ID: "raw"}},
	}}
	router := mountAGUIRouter(t, store)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/agui/runs/node-1/raw",
		strings.NewReader(runAgentInputBody(t, "t", "r", "x")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code, w.Body.String())
	frames := parseAGUIStream(t, w.Body.String())
	require.Equal(t, "plain text answer", frames[2].Data["delta"])
}

// TestAGUIRunHandler_ContextCancelMidFlight covers the <-ctx.Done() branch
// in the wait loop: client cancellation during a slow reasoner must return
// cleanly without emitting any post-RUN_STARTED frames.
func TestAGUIRunHandler_ContextCancelMidFlight(t *testing.T) {
	prev := AGUIHeartbeatInterval
	AGUIHeartbeatInterval = time.Hour
	defer func() { AGUIHeartbeatInterval = prev }()

	released := make(chan struct{})
	agentServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		select {
		case <-released:
		case <-r.Context().Done():
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"result":"too late"}`))
	}))
	defer func() { close(released); agentServer.Close() }()

	store := &reasonerTestStorage{agent: &types.AgentNode{
		ID:              "node-1",
		BaseURL:         agentServer.URL,
		HealthStatus:    types.HealthStatusActive,
		LifecycleStatus: types.AgentStatusReady,
		Reasoners:       []types.ReasonerDefinition{{ID: "hang"}},
	}}
	router := mountAGUIRouter(t, store)

	ctx, cancel := context.WithCancel(context.Background())
	req := httptest.NewRequest(http.MethodPost, "/api/v1/agui/runs/node-1/hang",
		strings.NewReader(runAgentInputBody(t, "t", "r", "x"))).WithContext(ctx)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	done := make(chan struct{})
	go func() {
		router.ServeHTTP(w, req)
		close(done)
	}()

	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		if strings.Contains(w.Body.String(), `"type":"RUN_STARTED"`) {
			break
		}
		time.Sleep(5 * time.Millisecond)
	}
	require.Contains(t, w.Body.String(), `"type":"RUN_STARTED"`, "RUN_STARTED should arrive before cancel")
	cancel()

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("handler did not return within 2s of context cancel")
	}

	body := w.Body.String()
	require.NotContains(t, body, "TEXT_MESSAGE_START")
	require.NotContains(t, body, "RUN_FINISHED")
}

// TestAGUIRunHandler_RejectsMalformedJSON covers the c.ShouldBindJSON error
// branch — completely invalid request bodies must be rejected as 400 before
// any of the agent lookup or stream-opening logic runs.
func TestAGUIRunHandler_RejectsMalformedJSON(t *testing.T) {
	store := &reasonerTestStorage{agent: &types.AgentNode{
		ID:              "node-1",
		BaseURL:         "http://unused",
		HealthStatus:    types.HealthStatusActive,
		LifecycleStatus: types.AgentStatusReady,
		Reasoners:       []types.ReasonerDefinition{{ID: "echo"}},
	}}
	router := mountAGUIRouter(t, store)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/agui/runs/node-1/echo", strings.NewReader("not-json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code, w.Body.String())
	require.NotEqual(t, "text/event-stream", w.Header().Get("Content-Type"))
}

// TestAGUIRunHandler_ValidationErrorsReturnJSON: pre-stream validation
// errors come back as plain JSON 4xx, never as an SSE stream. Once we emit
// RUN_STARTED the contract becomes "you'll see RUN_ERROR on failure" — but
// until the first frame, conventional REST errors win.
func TestAGUIRunHandler_ValidationErrorsReturnJSON(t *testing.T) {
	store := &reasonerTestStorage{agent: &types.AgentNode{
		ID:              "node-1",
		BaseURL:         "http://unused",
		HealthStatus:    types.HealthStatusActive,
		LifecycleStatus: types.AgentStatusReady,
		Reasoners:       []types.ReasonerDefinition{{ID: "echo"}},
	}}
	router := mountAGUIRouter(t, store)

	cases := []struct {
		name     string
		path     string
		wantCode int
		wantMsg  string
	}{
		{"unknown node", "/api/v1/agui/runs/missing-node/echo", http.StatusNotFound, "node 'missing-node' not found"},
		{"unknown reasoner on known node", "/api/v1/agui/runs/node-1/does-not-exist", http.StatusNotFound, "reasoner 'does-not-exist' not found"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, tc.path, strings.NewReader(runAgentInputBody(t, "t", "r", "x")))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			require.Equal(t, tc.wantCode, w.Code, w.Body.String())
			require.NotEqual(t, "text/event-stream", w.Header().Get("Content-Type"),
				"validation errors must not open the SSE stream")
			require.Contains(t, w.Body.String(), tc.wantMsg)
		})
	}
}

// TestAGUIRunHandler_ToolCalls_EmitsTriadAndAttachesToAssistantSnapshot
// covers Tier 2: when the reasoner declares a tool call (synthetic shape
// `{"toolCalls":[{id,name,arguments}]}`), the handler must emit
// TOOL_CALL_START → _ARGS → _END (BEFORE TEXT_MESSAGE_*) and attach the
// tool calls to the assistant turn in MESSAGES_SNAPSHOT — the wire shape
// CopilotKit's frontend pattern-matches against `useCopilotAction` to drive
// Generative UI.
func TestAGUIRunHandler_ToolCalls_EmitsTriadAndAttachesToAssistantSnapshot(t *testing.T) {
	agentServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"result":"booking your flight",
			"toolCalls":[
				{"id":"tc1","name":"showFlightCard","arguments":{"from":"SFO","to":"JFK"}}
			]
		}`))
	}))
	defer agentServer.Close()

	store := &reasonerTestStorage{agent: &types.AgentNode{
		ID:              "node-1",
		BaseURL:         agentServer.URL,
		HealthStatus:    types.HealthStatusActive,
		LifecycleStatus: types.AgentStatusReady,
		Reasoners:       []types.ReasonerDefinition{{ID: "agent"}},
	}}
	router := mountAGUIRouter(t, store)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/agui/runs/node-1/agent",
		strings.NewReader(runAgentInputBody(t, "t", "r", "book me SFO->JFK")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code, w.Body.String())
	frames := parseAGUIStream(t, w.Body.String())

	wantTypes := []string{
		"RUN_STARTED",
		"TOOL_CALL_START",
		"TOOL_CALL_ARGS",
		"TOOL_CALL_END",
		"TEXT_MESSAGE_START",
		"TEXT_MESSAGE_CONTENT",
		"TEXT_MESSAGE_END",
		"MESSAGES_SNAPSHOT",
		"RUN_FINISHED",
	}
	require.Len(t, frames, len(wantTypes))
	for i, want := range wantTypes {
		require.Equal(t, want, frames[i].Type(), "frame %d: %v", i, frames[i].Data)
	}

	require.Equal(t, "tc1", frames[1].Data["toolCallId"])
	require.Equal(t, "showFlightCard", frames[1].Data["toolCallName"])
	// parentMessageId stitches the tool call into the assistant turn.
	require.NotEmpty(t, frames[1].Data["parentMessageId"])
	require.Equal(t, frames[1].Data["parentMessageId"], frames[4].Data["messageId"])

	require.Equal(t, "tc1", frames[2].Data["toolCallId"])
	require.JSONEq(t, `{"from":"SFO","to":"JFK"}`, frames[2].Data["delta"].(string))
	require.Equal(t, "tc1", frames[3].Data["toolCallId"])

	require.Equal(t, "booking your flight", frames[5].Data["delta"])

	// MESSAGES_SNAPSHOT carries the tool-call attached to the assistant turn.
	snap, _ := frames[7].Data["messages"].([]any)
	require.Len(t, snap, 2)
	assistant, _ := snap[1].(map[string]any)
	require.Equal(t, "assistant", assistant["role"])
	tcs, _ := assistant["toolCalls"].([]any)
	require.Len(t, tcs, 1)
	tc, _ := tcs[0].(map[string]any)
	require.Equal(t, "tc1", tc["id"])
	require.Equal(t, "function", tc["type"])
	fn, _ := tc["function"].(map[string]any)
	require.Equal(t, "showFlightCard", fn["name"])
	require.JSONEq(t, `{"from":"SFO","to":"JFK"}`, fn["arguments"].(string))
}

// TestAGUIRunHandler_ToolCalls_AutoIDIfMissing covers the synthetic-id
// fallback in extractToolCalls.
func TestAGUIRunHandler_ToolCalls_AutoIDIfMissing(t *testing.T) {
	agentServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"result":"ok","toolCalls":[{"name":"alpha"}]}`))
	}))
	defer agentServer.Close()

	store := &reasonerTestStorage{agent: &types.AgentNode{
		ID:              "node-1",
		BaseURL:         agentServer.URL,
		HealthStatus:    types.HealthStatusActive,
		LifecycleStatus: types.AgentStatusReady,
		Reasoners:       []types.ReasonerDefinition{{ID: "a"}},
	}}
	router := mountAGUIRouter(t, store)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/agui/runs/node-1/a",
		strings.NewReader(runAgentInputBody(t, "t", "r", "x")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	frames := parseAGUIStream(t, w.Body.String())
	require.Equal(t, "TOOL_CALL_START", frames[1].Type())
	id, _ := frames[1].Data["toolCallId"].(string)
	require.NotEmpty(t, id, "tool-call id must be auto-generated when missing")
	// Same id must propagate through the triad.
	require.Equal(t, id, frames[2].Data["toolCallId"])
	require.Equal(t, id, frames[3].Data["toolCallId"])
}

// TestAGUIRunHandler_ToolCalls_SkipsMalformedEntries — a tool-call with no
// name is silently dropped (rather than failing the whole turn). Mirrors the
// extractToolCalls guards.
func TestAGUIRunHandler_ToolCalls_SkipsMalformedEntries(t *testing.T) {
	agentServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"result":"ok","toolCalls":[
			{"id":"x","name":""},
			"not-an-object",
			{"id":"y","name":"good","arguments":{}}
		]}`))
	}))
	defer agentServer.Close()

	store := &reasonerTestStorage{agent: &types.AgentNode{
		ID:              "node-1",
		BaseURL:         agentServer.URL,
		HealthStatus:    types.HealthStatusActive,
		LifecycleStatus: types.AgentStatusReady,
		Reasoners:       []types.ReasonerDefinition{{ID: "a"}},
	}}
	router := mountAGUIRouter(t, store)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/agui/runs/node-1/a",
		strings.NewReader(runAgentInputBody(t, "t", "r", "x")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	frames := parseAGUIStream(t, w.Body.String())
	starts := 0
	for _, f := range frames {
		if f.Type() == "TOOL_CALL_START" {
			starts++
			require.Equal(t, "good", f.Data["toolCallName"])
		}
	}
	require.Equal(t, 1, starts, "only the well-formed tool call should be emitted")
}

// TestAGUIRunHandler_State_EmitsSnapshotAndForwardsInbound covers Tier 3:
// the inbound `state` field on RunAgentInput must reach the reasoner, and a
// reasoner-returned `state` field must be re-emitted as a STATE_SNAPSHOT
// before MESSAGES_SNAPSHOT and RUN_FINISHED.
func TestAGUIRunHandler_State_EmitsSnapshotAndForwardsInbound(t *testing.T) {
	var seenInput map[string]any
	agentServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		raw, _ := io.ReadAll(r.Body)
		require.NoError(t, json.Unmarshal(raw, &seenInput))
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"result":"counter incremented","state":{"counter":2}}`))
	}))
	defer agentServer.Close()

	store := &reasonerTestStorage{agent: &types.AgentNode{
		ID:              "node-1",
		BaseURL:         agentServer.URL,
		HealthStatus:    types.HealthStatusActive,
		LifecycleStatus: types.AgentStatusReady,
		Reasoners:       []types.ReasonerDefinition{{ID: "stateful"}},
	}}
	router := mountAGUIRouter(t, store)

	body := `{
		"threadId":"t","runId":"r",
		"messages":[{"role":"user","content":"increment"}],
		"state":{"counter":1}
	}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/agui/runs/node-1/stateful", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code, w.Body.String())

	// Inbound state landed on the reasoner.
	gotState, _ := seenInput["state"].(map[string]any)
	require.EqualValues(t, 1, gotState["counter"])

	frames := parseAGUIStream(t, w.Body.String())
	// Find STATE_SNAPSHOT in the stream and verify it carries the new value.
	var snap aguiFrame
	for _, f := range frames {
		if f.Type() == "STATE_SNAPSHOT" {
			snap = f
			break
		}
	}
	require.NotEmpty(t, snap.Data, "STATE_SNAPSHOT must be emitted when reasoner returns state")
	snapVal, _ := snap.Data["snapshot"].(map[string]any)
	require.EqualValues(t, 2, snapVal["counter"])

	// Order: STATE_SNAPSHOT after TEXT_MESSAGE_END but before MESSAGES_SNAPSHOT.
	idx := func(typ string) int {
		for i, f := range frames {
			if f.Type() == typ {
				return i
			}
		}
		return -1
	}
	require.Less(t, idx("TEXT_MESSAGE_END"), idx("STATE_SNAPSHOT"))
	require.Less(t, idx("STATE_SNAPSHOT"), idx("MESSAGES_SNAPSHOT"))
	require.Less(t, idx("MESSAGES_SNAPSHOT"), idx("RUN_FINISHED"))
}

// TestAGUIRunHandler_State_OmittedWhenReasonerDoesNotReturnIt — Tier 3
// doesn't synthesize a STATE_SNAPSHOT for stateless reasoners; we only emit
// when the reasoner opts in via a top-level `state` field.
func TestAGUIRunHandler_State_OmittedWhenReasonerDoesNotReturnIt(t *testing.T) {
	agentServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"result":"plain"}`))
	}))
	defer agentServer.Close()

	store := &reasonerTestStorage{agent: &types.AgentNode{
		ID:              "node-1",
		BaseURL:         agentServer.URL,
		HealthStatus:    types.HealthStatusActive,
		LifecycleStatus: types.AgentStatusReady,
		Reasoners:       []types.ReasonerDefinition{{ID: "plain"}},
	}}
	router := mountAGUIRouter(t, store)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/agui/runs/node-1/plain",
		strings.NewReader(runAgentInputBody(t, "t", "r", "x")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	frames := parseAGUIStream(t, w.Body.String())
	for _, f := range frames {
		require.NotEqual(t, "STATE_SNAPSHOT", f.Type(),
			"STATE_SNAPSHOT must not be emitted unless the reasoner opts in")
	}
}

// TestAGUIRunHandler_PassesToolMessagesThrough — when the inbound history
// contains a `role:"tool"` message (CopilotKit posts these on the next run
// after a frontend useCopilotAction completes), it must reach the reasoner
// intact.
func TestAGUIRunHandler_PassesToolMessagesThrough(t *testing.T) {
	var seenInput map[string]any
	agentServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		raw, _ := io.ReadAll(r.Body)
		require.NoError(t, json.Unmarshal(raw, &seenInput))
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"result":"thanks"}`))
	}))
	defer agentServer.Close()

	store := &reasonerTestStorage{agent: &types.AgentNode{
		ID:              "node-1",
		BaseURL:         agentServer.URL,
		HealthStatus:    types.HealthStatusActive,
		LifecycleStatus: types.AgentStatusReady,
		Reasoners:       []types.ReasonerDefinition{{ID: "echo"}},
	}}
	router := mountAGUIRouter(t, store)

	body := `{
		"threadId":"t","runId":"r2",
		"messages":[
			{"role":"user","content":"book SFO->JFK"},
			{"role":"assistant","toolCalls":[{"id":"tc1","type":"function","function":{"name":"showFlightCard","arguments":"{\"from\":\"SFO\"}"}}]},
			{"role":"tool","toolCallId":"tc1","content":"user clicked confirm"},
			{"role":"user","content":"now book the return"}
		]
	}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/agui/runs/node-1/echo", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code, w.Body.String())
	require.Equal(t, "now book the return", seenInput["prompt"])
	msgs, _ := seenInput["messages"].([]any)
	require.Len(t, msgs, 4)
	toolMsg, _ := msgs[2].(map[string]any)
	require.Equal(t, "tool", toolMsg["role"])
	require.Equal(t, "tc1", toolMsg["toolCallId"])
	require.Equal(t, "user clicked confirm", toolMsg["content"])
}

// TestHTTPAgentInvoker_HappyPath exercises the real httpAgentInvoker
// against a stub agent server — handler tests use an interface stub so this
// concrete path otherwise goes uncovered.
func TestHTTPAgentInvoker_HappyPath(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/reasoners/ping", r.URL.Path)
		require.Equal(t, "application/json", r.Header.Get("Content-Type"))
		got, _ := io.ReadAll(r.Body)
		require.JSONEq(t, `{"k":1}`, string(got))
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true}`))
	}))
	defer server.Close()

	body, err := httpAgentInvoker{}.Invoke(context.Background(),
		&types.AgentNode{BaseURL: server.URL}, "ping", []byte(`{"k":1}`))
	require.NoError(t, err)
	require.JSONEq(t, `{"ok":true}`, string(body))
}

// TestHTTPAgentInvoker_4xxBubblesUpAsError covers the resp.StatusCode >= 400
// branch — the body is still returned but as a callError so the handler can
// turn it into a RUN_ERROR.
func TestHTTPAgentInvoker_4xxBubblesUpAsError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"oops":"server"}`))
	}))
	defer server.Close()

	body, err := httpAgentInvoker{}.Invoke(context.Background(),
		&types.AgentNode{BaseURL: server.URL}, "boom", []byte(`{}`))
	require.Error(t, err)
	require.Contains(t, err.Error(), "agent returned 500")
	require.Contains(t, string(body), "oops")
}

// TestHTTPAgentInvoker_DialFailureSurfacesError covers the client.Do error
// branch by pointing the invoker at a closed listener.
func TestHTTPAgentInvoker_DialFailureSurfacesError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}))
	addr := server.URL
	server.Close()

	_, err := httpAgentInvoker{}.Invoke(context.Background(),
		&types.AgentNode{BaseURL: addr}, "ping", []byte(`{}`))
	require.Error(t, err)
	require.Contains(t, err.Error(), "agent call failed")
}

// TestHTTPAgentInvoker_BadURLFailsRequestConstruction covers the
// http.NewRequestWithContext error branch — an invalid URL never makes it
// to a dial.
func TestHTTPAgentInvoker_BadURLFailsRequestConstruction(t *testing.T) {
	_, err := httpAgentInvoker{}.Invoke(context.Background(),
		&types.AgentNode{BaseURL: "http://bad\nhost"}, "ping", []byte(`{}`))
	require.Error(t, err)
	require.Contains(t, err.Error(), "create agent request")
}
