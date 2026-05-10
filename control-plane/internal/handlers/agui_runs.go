package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/Agent-Field/agentfield/control-plane/internal/agui"
	"github.com/Agent-Field/agentfield/control-plane/internal/storage"
	"github.com/Agent-Field/agentfield/control-plane/internal/utils"
	"github.com/Agent-Field/agentfield/control-plane/pkg/types"

	"github.com/gin-gonic/gin"
)

// AGUIHeartbeatInterval is how often we emit an SSE comment (`: keep-alive`)
// while waiting for a slow reasoner. AG-UI clients silently drop comment
// lines per the SSE spec, but proxies (nginx, ALBs) see the bytes and don't
// idle out the connection. 15s leaves comfortable headroom under the 60s
// nginx default. Exposed for tests.
var AGUIHeartbeatInterval = 15 * time.Second

// agentInvocation is the result of calling the agent reasoner endpoint:
// either a fully buffered body (for traditional reasoners that return a
// single JSON object) or a live io.ReadCloser carrying NDJSON chunks
// (for streaming reasoners that yield events as they happen). Exactly
// one of Body / Stream is non-nil; ContentType disambiguates.
type agentInvocation struct {
	Body        []byte
	Stream      io.ReadCloser
	ContentType string
}

// IsStreaming reports whether the reasoner returned an NDJSON stream
// (Content-Type: application/x-ndjson) the handler should consume
// chunk-by-chunk and forward as live AG-UI events.
func (r *agentInvocation) IsStreaming() bool {
	return r != nil && r.Stream != nil
}

// agentInvoker abstracts the outbound HTTP call to the agent's reasoner so
// tests can stub behavior without spinning up a real server. The default
// implementation (httpAgentInvoker) buffers the body for non-NDJSON
// responses and hands back the live stream for NDJSON.
type agentInvoker interface {
	Invoke(ctx context.Context, agent *types.AgentNode, reasonerName string, input []byte) (*agentInvocation, error)
}

type httpAgentInvoker struct{ client *http.Client }

func (i httpAgentInvoker) Invoke(ctx context.Context, agent *types.AgentNode, reasonerName string, input []byte) (*agentInvocation, error) {
	url := fmt.Sprintf("%s/reasoners/%s", agent.BaseURL, reasonerName)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(input))
	if err != nil {
		return nil, fmt.Errorf("create agent request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	// Tell the reasoner we accept either a plain JSON response (buffered
	// path) or NDJSON (streaming path). Reasoners that opted into
	// streaming can switch on Accept; reasoners that didn't ignore it.
	req.Header.Set("Accept", "application/x-ndjson, application/json")

	client := i.client
	if client == nil {
		client = http.DefaultClient
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("agent call failed: %w", err)
	}

	ct := resp.Header.Get("Content-Type")
	// Streaming response: hand the body straight back. Caller is
	// responsible for closing it; we don't read it here.
	if resp.StatusCode < http.StatusBadRequest && strings.HasPrefix(ct, "application/x-ndjson") {
		return &agentInvocation{Stream: resp.Body, ContentType: ct}, nil
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read agent response: %w", err)
	}
	if resp.StatusCode >= http.StatusBadRequest {
		return &agentInvocation{Body: body, ContentType: ct},
			fmt.Errorf("agent returned %d: %s", resp.StatusCode, truncateForLog(body))
	}
	return &agentInvocation{Body: body, ContentType: ct}, nil
}

// AGUIRunHandler handles POST /api/v1/agui/runs/:node_id/:reasoner_name.
//
// It is the AG-UI protocol adapter: clients (CopilotKit's CopilotRuntime
// proxying through @ag-ui/client's HttpAgent, or any other AG-UI consumer)
// post a canonical RunAgentInput body, the handler invokes the named
// reasoner, and the response is a Server-Sent Events stream of AG-UI events.
//
// Capabilities (see https://docs.ag-ui.com/concepts/events):
//
//   - Lifecycle: RUN_STARTED / RUN_FINISHED / RUN_ERROR.
//   - Text messages: TEXT_MESSAGE_START / _CONTENT / _END for the
//     assistant turn. The single TEXT_MESSAGE_CONTENT carries the
//     reasoner's full result; token-level streaming is a follow-up.
//   - Tool calls: if the reasoner result contains a `toolCalls` array
//     (one per `useCopilotAction`-style render), TOOL_CALL_START /
//     TOOL_CALL_ARGS / TOOL_CALL_END frames are emitted before the
//     text turn closes. CopilotKit's frontend pattern-matches
//     `toolCallName` against registered actions to drive Generative UI.
//   - State: if the reasoner result contains a `state` object,
//     STATE_SNAPSHOT is emitted before RUN_FINISHED — the value
//     `useCoAgent({ state })` reads on the client.
//   - MESSAGES_SNAPSHOT closes every successful run with the canonical
//     conversation history, so multi-turn clients can persist it.
func AGUIRunHandler(storageProvider storage.StorageProvider) gin.HandlerFunc {
	return aguiRunHandler(storageProvider, httpAgentInvoker{})
}

func aguiRunHandler(storageProvider storage.StorageProvider, invoker agentInvoker) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		nodeID := strings.TrimSpace(c.Param("node_id"))
		reasonerName := strings.TrimSpace(c.Param("reasoner_name"))
		if nodeID == "" || reasonerName == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "node_id and reasoner_name are required"})
			return
		}

		var req agui.RunAgentInput
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		agent, err := storageProvider.GetAgent(ctx, nodeID)
		if err != nil || agent == nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error": fmt.Sprintf("node '%s' not found", nodeID),
			})
			return
		}
		if !reasonerExists(agent, reasonerName) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": fmt.Sprintf("reasoner '%s' not found on node '%s'", reasonerName, nodeID),
			})
			return
		}

		// Validation passed — switch to streaming mode. From here on we
		// report failures via RunError frames instead of HTTP error
		// responses, since the SSE stream is already open.
		threadID := req.ThreadID
		if threadID == "" {
			threadID = "thread-" + utils.GenerateExecutionID()
		}
		runID := req.RunID
		if runID == "" {
			runID = "run-" + utils.GenerateExecutionID()
		}

		c.Header("Content-Type", "text/event-stream")
		c.Header("Cache-Control", "no-cache")
		c.Header("Connection", "keep-alive")
		c.Header("X-Accel-Buffering", "no")

		flush := func() {
			if f, ok := c.Writer.(http.Flusher); ok {
				f.Flush()
			}
		}

		write := func(ev agui.Event) bool {
			if err := agui.WriteSSE(c.Writer, ev); err != nil {
				return false
			}
			flush()
			return true
		}

		if !write(agui.RunStarted{
			ThreadID:  threadID,
			RunID:     runID,
			Timestamp: agui.NowMillis(),
		}) {
			return
		}

		reasonerInput := buildReasonerInput(req)
		inputJSON, err := json.Marshal(reasonerInput)
		if err != nil {
			write(agui.RunError{
				Message:   fmt.Sprintf("failed to marshal input: %v", err),
				Code:      "ERR_INPUT_MARSHAL",
				Timestamp: agui.NowMillis(),
			})
			return
		}

		// Run the agent invocation in a goroutine so the main loop can
		// emit SSE keep-alive comments while we wait for the first byte.
		// (Once the body starts streaming, that's its own activity.) AG-UI
		// has no heartbeat event, but `:` comment frames are valid SSE
		// that clients ignore and proxies see as activity.
		type invokeResultT struct {
			res *agentInvocation
			err error
		}
		resultCh := make(chan invokeResultT, 1)
		go func() {
			r, e := invoker.Invoke(ctx, agent, reasonerName, inputJSON)
			resultCh <- invokeResultT{res: r, err: e}
		}()

		ticker := time.NewTicker(AGUIHeartbeatInterval)
		defer ticker.Stop()

		var (
			invocation *agentInvocation
			invokeErr  error
		)
	waitLoop:
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				if _, err := fmt.Fprint(c.Writer, ": keep-alive\n\n"); err != nil {
					return
				}
				flush()
			case r := <-resultCh:
				invocation, invokeErr = r.res, r.err
				break waitLoop
			}
		}

		if invokeErr != nil {
			write(agui.RunError{
				Message:   invokeErr.Error(),
				Code:      "ERR_AGENT_CALL",
				Timestamp: agui.NowMillis(),
			})
			return
		}

		messageID := "msg-" + utils.GenerateExecutionID()

		// Streaming reasoner — Content-Type is application/x-ndjson and
		// the body is a live chunk stream. Drain it here, dispatching
		// each tagged event to its AG-UI counterpart immediately, then
		// wrap up with MESSAGES_SNAPSHOT + RUN_FINISHED. This is the
		// path that makes "Generative UI" feel live instead of stuttery.
		if invocation.IsStreaming() {
			runStreamingDispatch(ctx, c, write, invocation.Stream, req, threadID, runID, messageID)
			return
		}

		// Buffered reasoner — the rest of this function processes the
		// fully-buffered JSON body the way it did before streaming
		// support landed.
		body := invocation.Body
		// Decode the agent response so we can surface the structured pieces
		// CopilotKit understands: tool calls, state, and the assistant text.
		parsed, parsedOK := decodeReasonerResponse(body)

		// Reasoning segments first — frontends render these in a
		// collapsible "Thinking…" pane above the user-facing answer, so
		// emitting them before tool calls / text matches the UX flow.
		if reasoning := extractReasoning(parsed); len(reasoning) > 0 {
			reasoningContextID := "reasoning-" + utils.GenerateExecutionID()
			if !write(agui.ReasoningStart{
				MessageID: reasoningContextID,
				Timestamp: agui.NowMillis(),
			}) {
				return
			}
			for i, seg := range reasoning {
				segID := seg.ID
				if segID == "" {
					segID = fmt.Sprintf("%s-seg-%d", reasoningContextID, i)
				}
				if !write(agui.ReasoningMessageStart{
					MessageID: segID,
					Role:      "reasoning",
					Timestamp: agui.NowMillis(),
				}) {
					return
				}
				for _, chunk := range chunkText(seg.Content, AGUITextChunkSize) {
					if !write(agui.ReasoningMessageContent{
						MessageID: segID,
						Delta:     chunk,
						Timestamp: agui.NowMillis(),
					}) {
						return
					}
				}
				if !write(agui.ReasoningMessageEnd{
					MessageID: segID,
					Timestamp: agui.NowMillis(),
				}) {
					return
				}
			}
			if !write(agui.ReasoningEnd{
				MessageID: reasoningContextID,
				Timestamp: agui.NowMillis(),
			}) {
				return
			}
		}

		// Tool calls next so the frontend can dispatch render handlers
		// (useCopilotAction) before the text turn closes. The text turn
		// then carries any textual answer the reasoner produced.
		toolCalls := extractToolCalls(parsed)
		assistantToolCalls := make([]agui.ToolCall, 0, len(toolCalls))
		for _, tc := range toolCalls {
			argsJSON, _ := json.Marshal(tc.Arguments)
			argsStr := string(argsJSON)
			if !write(agui.ToolCallStart{
				ToolCallID:      tc.ID,
				ToolCallName:    tc.Name,
				ParentMessageID: messageID,
				Timestamp:       agui.NowMillis(),
			}) {
				return
			}
			if !write(agui.ToolCallArgs{
				ToolCallID: tc.ID,
				Delta:      argsStr,
				Timestamp:  agui.NowMillis(),
			}) {
				return
			}
			if !write(agui.ToolCallEnd{
				ToolCallID: tc.ID,
				Timestamp:  agui.NowMillis(),
			}) {
				return
			}
			// If the reasoner already executed the tool server-side and
			// gave us a result (e.g. a .ai(tools=...) trace), emit
			// TOOL_CALL_RESULT so the trace renders in the same place the
			// frontend would expect a tool message to live.
			if tc.HasResult {
				if !write(agui.ToolCallResult{
					MessageID:  "msg-toolresult-" + tc.ID,
					ToolCallID: tc.ID,
					Content:    stringifyResult(tc.Result),
					Role:       "tool",
					Timestamp:  agui.NowMillis(),
				}) {
					return
				}
			}
			assistantToolCalls = append(assistantToolCalls, agui.ToolCall{
				ID:   tc.ID,
				Type: "function",
				Function: agui.ToolCallFunction{
					Name:      tc.Name,
					Arguments: argsStr,
				},
			})
		}

		// Text turn. Assembled even when empty so clients see a complete
		// triad — schema permits empty delta. Long replies are chunked
		// across multiple TEXT_MESSAGE_CONTENT frames so frontends can
		// paint progressively even though the reasoner is synchronous.
		assistantText := extractAssistantText(parsed, parsedOK, body)
		if !write(agui.TextMessageStart{
			MessageID: messageID,
			Role:      "assistant",
			Timestamp: agui.NowMillis(),
		}) {
			return
		}
		for _, chunk := range chunkText(assistantText, AGUITextChunkSize) {
			if !write(agui.TextMessageContent{
				MessageID: messageID,
				Delta:     chunk,
				Timestamp: agui.NowMillis(),
			}) {
				return
			}
		}
		if !write(agui.TextMessageEnd{
			MessageID: messageID,
			Timestamp: agui.NowMillis(),
		}) {
			return
		}

		// State snapshot first (if reasoner returned full state), then
		// any RFC 6902 patches the reasoner emits via `stateDelta`.
		// Snapshot before MESSAGES_SNAPSHOT so the client correlates the
		// new state with the new turn.
		if state, hasState := extractState(parsed); hasState {
			if !write(agui.StateSnapshot{
				Snapshot:  state,
				Timestamp: agui.NowMillis(),
			}) {
				return
			}
		}
		if delta := extractStateDelta(parsed); delta != nil {
			if !write(agui.StateDelta{
				Delta:     delta,
				Timestamp: agui.NowMillis(),
			}) {
				return
			}
		}

		// Canonical history snapshot: inbound messages + the assistant turn
		// we just produced.
		assistant := agui.Message{
			ID:        messageID,
			Role:      "assistant",
			Content:   assistantText,
			ToolCalls: assistantToolCalls,
		}
		full := append([]agui.Message{}, req.Messages...)
		full = append(full, assistant)
		if !write(agui.MessagesSnapshot{
			Messages:  full,
			Timestamp: agui.NowMillis(),
		}) {
			return
		}

		write(agui.RunFinished{
			ThreadID:  threadID,
			RunID:     runID,
			Outcome:   &agui.Outcome{Type: "success"},
			Result:    parsed,
			Timestamp: agui.NowMillis(),
		})
	}
}

// buildReasonerInput translates a canonical AG-UI RunAgentInput into the
// dict shape AgentField reasoners receive. We pass the full envelope (so
// reasoners that care can inspect tools/state/messages/context) plus a
// `prompt` convenience extracted from the trailing user message.
func buildReasonerInput(req agui.RunAgentInput) map[string]any {
	input := map[string]any{
		"prompt":   req.LastUserMessageText(),
		"messages": req.Messages,
		"tools":    req.Tools,
		"context":  req.Context,
		"threadId": req.ThreadID,
		"runId":    req.RunID,
	}
	if len(req.State) > 0 {
		var state any
		if err := json.Unmarshal(req.State, &state); err == nil {
			input["state"] = state
		}
	}
	if len(req.ForwardedProps) > 0 {
		var fp any
		if err := json.Unmarshal(req.ForwardedProps, &fp); err == nil {
			input["forwardedProps"] = fp
		}
	}
	return input
}

// decodeReasonerResponse json-decodes the agent body. Returns the parsed
// value and whether decoding succeeded; non-JSON responses fall through to
// the raw-body path in extractAssistantText.
func decodeReasonerResponse(body []byte) (any, bool) {
	var parsed any
	if err := json.Unmarshal(body, &parsed); err != nil {
		return nil, false
	}
	return parsed, true
}

// reasonerToolCall is the synthetic shape AgentField reasoners use to
// declare tool calls. Reasoners return
//
//	{"toolCalls": [{"id", "name", "arguments", "result"?}, ...]}
//
// to drive frontend useCopilotAction renders. The optional `result` field,
// when present, indicates the call was already executed server-side and
// causes us to emit TOOL_CALL_RESULT after TOOL_CALL_END — so the trace
// (e.g. from .ai(tools=...) ToolCallTrace) shows up in the UI alongside
// the live calls.
type reasonerToolCall struct {
	ID        string
	Name      string
	Arguments any
	Result    any
	HasResult bool
}

// extractToolCalls reads a `toolCalls` array from the reasoner response,
// if present. Each entry needs at least a name; id and arguments are
// optional and synthesized when missing. `result` is optional.
func extractToolCalls(parsed any) []reasonerToolCall {
	obj, ok := parsed.(map[string]any)
	if !ok {
		return nil
	}
	raw, ok := obj["toolCalls"].([]any)
	if !ok {
		return nil
	}
	out := make([]reasonerToolCall, 0, len(raw))
	for i, entry := range raw {
		m, ok := entry.(map[string]any)
		if !ok {
			continue
		}
		name, _ := m["name"].(string)
		if name == "" {
			continue
		}
		id, _ := m["id"].(string)
		if id == "" {
			id = fmt.Sprintf("toolcall-%d-%s", i, utils.GenerateExecutionID())
		}
		args := m["arguments"]
		if args == nil {
			args = map[string]any{}
		}
		result, hasResult := m["result"]
		out = append(out, reasonerToolCall{
			ID:        id,
			Name:      name,
			Arguments: args,
			Result:    result,
			HasResult: hasResult,
		})
	}
	return out
}

// extractReasoning reads a chain-of-thought from the reasoner response.
// Reasoners that want to surface model thinking in CopilotKit's "Thinking…"
// pane return either:
//
//	{"reasoning": "the agent's chain-of-thought as a single string"}
//
// or a list of per-step strings:
//
//	{"reasoning": ["step 1...", "step 2..."]}
//
// In either case the handler emits REASONING_START → one or more
// REASONING_MESSAGE_START / _CONTENT / _END pairs → REASONING_END.
// Reasoners that already structured the trace can pass an explicit list
// of segment dicts:
//
//	{"reasoning": [{"id": "r-0", "content": "..."}, ...]}
func extractReasoning(parsed any) []reasoningSegment {
	obj, ok := parsed.(map[string]any)
	if !ok {
		return nil
	}
	raw, has := obj["reasoning"]
	if !has || raw == nil {
		return nil
	}
	switch v := raw.(type) {
	case string:
		if v == "" {
			return nil
		}
		return []reasoningSegment{{Content: v}}
	case []any:
		out := make([]reasoningSegment, 0, len(v))
		for i, entry := range v {
			switch s := entry.(type) {
			case string:
				if s == "" {
					continue
				}
				out = append(out, reasoningSegment{Content: s})
			case map[string]any:
				content, _ := s["content"].(string)
				if content == "" {
					continue
				}
				id, _ := s["id"].(string)
				if id == "" {
					id = fmt.Sprintf("r-%d-%s", i, utils.GenerateExecutionID())
				}
				out = append(out, reasoningSegment{ID: id, Content: content})
			}
		}
		if len(out) == 0 {
			return nil
		}
		return out
	}
	return nil
}

type reasoningSegment struct {
	ID      string
	Content string
}

// extractStateDelta reads a `stateDelta` array from the reasoner response,
// if present. Reasoners that prefer to emit incremental RFC 6902 patches
// instead of (or in addition to) full snapshots return:
//
//	{"stateDelta": [{"op":"replace","path":"/counter","value":2}, ...]}
//
// The handler emits this as a STATE_DELTA event. Both forms can coexist:
// emit STATE_SNAPSHOT first to establish a baseline, then STATE_DELTA for
// fine-grained updates.
func extractStateDelta(parsed any) []any {
	obj, ok := parsed.(map[string]any)
	if !ok {
		return nil
	}
	raw, ok := obj["stateDelta"].([]any)
	if !ok || len(raw) == 0 {
		return nil
	}
	return raw
}

// AGUITextChunkSize is the maximum size of a single TEXT_MESSAGE_CONTENT
// delta. Long reasoner responses are split into multiple deltas so the
// frontend can begin painting before the full reply lands. 256 chars is
// the sweet spot: small enough that long replies render progressively,
// large enough that short replies fit in one frame and don't pay extra
// SSE overhead. Exposed for tests.
var AGUITextChunkSize = 256

// chunkText splits a string into pieces of up to size bytes. For empty
// input, returns a single empty chunk so callers always emit one
// TEXT_MESSAGE_CONTENT delta (the schema permits empty deltas, and a
// missing content frame would break clients that expect the full triad).
// Splits on rune boundaries so multi-byte UTF-8 sequences (emoji, CJK)
// don't get cut mid-byte.
func chunkText(s string, size int) []string {
	if size <= 0 {
		return []string{s}
	}
	if s == "" {
		return []string{""}
	}
	out := make([]string, 0, (len(s)/size)+1)
	current := make([]rune, 0, size)
	currentBytes := 0
	for _, r := range s {
		rb := len(string(r))
		if currentBytes+rb > size && len(current) > 0 {
			out = append(out, string(current))
			current = current[:0]
			currentBytes = 0
		}
		current = append(current, r)
		currentBytes += rb
	}
	if len(current) > 0 {
		out = append(out, string(current))
	}
	return out
}

// extractState returns the reasoner's top-level `state` field if any,
// for emission as STATE_SNAPSHOT.
func extractState(parsed any) (any, bool) {
	obj, ok := parsed.(map[string]any)
	if !ok {
		return nil, false
	}
	state, has := obj["state"]
	return state, has
}

// extractAssistantText picks the human-facing answer for the assistant
// turn. Priority:
//  1. Reasoner returned a top-level `result` field — stringify it.
//  2. Reasoner returned a top-level `content` field — stringify it.
//  3. Reasoner returned a string body — use it verbatim.
//  4. Otherwise return the JSON-encoded body with `toolCalls` and `state`
//     stripped, so the user sees something sensible if they didn't follow
//     the `result` / `content` convention.
//  5. If the body wasn't JSON at all, return it raw.
func extractAssistantText(parsed any, parsedOK bool, rawBody []byte) string {
	if !parsedOK {
		return string(rawBody)
	}
	if obj, ok := parsed.(map[string]any); ok {
		if r, has := obj["result"]; has {
			return stringifyResult(r)
		}
		if r, has := obj["content"]; has {
			return stringifyResult(r)
		}
		filtered := make(map[string]any, len(obj))
		for k, v := range obj {
			if k == "toolCalls" || k == "state" {
				continue
			}
			filtered[k] = v
		}
		if len(filtered) == 0 {
			return ""
		}
		return stringifyResult(filtered)
	}
	if s, ok := parsed.(string); ok {
		return s
	}
	return stringifyResult(parsed)
}

func reasonerExists(agent *types.AgentNode, name string) bool {
	for _, r := range agent.Reasoners {
		if r.ID == name {
			return true
		}
	}
	return false
}

// stringifyResult renders an arbitrary JSON value as a text chunk suitable
// for the AG-UI TextMessageContent delta. Strings pass through verbatim;
// everything else is JSON-encoded.
func stringifyResult(v any) string {
	if s, ok := v.(string); ok {
		return s
	}
	if v == nil {
		return ""
	}
	encoded, err := json.Marshal(v)
	if err != nil {
		return fmt.Sprintf("%v", v)
	}
	return string(encoded)
}
