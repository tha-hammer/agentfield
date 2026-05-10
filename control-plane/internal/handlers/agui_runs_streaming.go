package handlers

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/Agent-Field/agentfield/control-plane/internal/agui"
	"github.com/Agent-Field/agentfield/control-plane/internal/utils"

	"github.com/gin-gonic/gin"
)

// AGUIStreamingMaxLineBytes caps the size of any one NDJSON chunk the
// reasoner can send. Without this, a misbehaving reasoner could stream
// an unbounded line and exhaust handler memory. 1 MiB is generous for
// per-token deltas while still bounding the worst case. Exposed for tests.
var AGUIStreamingMaxLineBytes = 1 << 20

// streamingChunk is the wire shape between an AgentField streaming
// reasoner and the AG-UI handler. Reasoners emit one JSON object per
// line on stdout (NDJSON); this struct decodes them. All fields are
// optional — `Type` selects the variant.
//
// Recognized variants and their AG-UI translation:
//
//	{"type":"text",    "delta":"hello"}                 -> TEXT_MESSAGE_CONTENT
//	{"type":"reasoning","delta":"thinking..."}          -> REASONING_MESSAGE_CONTENT
//	{"type":"tool_call_start","id":"tc1","name":"x", "arguments":{...}, "parentMessageId":"..."}
//	                                                    -> TOOL_CALL_START + (single
//	                                                       TOOL_CALL_ARGS if arguments
//	                                                       supplied)
//	{"type":"tool_call_args", "id":"tc1","delta":"..."} -> TOOL_CALL_ARGS
//	{"type":"tool_call_end",  "id":"tc1"}               -> TOOL_CALL_END
//	{"type":"tool_call_result","id":"tc1","content":"..."} -> TOOL_CALL_RESULT
//	{"type":"state",          "snapshot":{...}}         -> STATE_SNAPSHOT
//	{"type":"state_delta",    "ops":[...]}              -> STATE_DELTA (RFC 6902)
//	{"type":"step_started",   "name":"plan"}            -> STEP_STARTED
//	{"type":"step_finished",  "name":"plan"}            -> STEP_FINISHED
//	{"type":"raw",            "event":..., "source":"x"} -> RAW
//	{"type":"custom",         "name":"...","value":...} -> CUSTOM
//	{"type":"final",          "data":{<buffered-shape>}} -> applies any
//	   leftover toolCalls / state / stateDelta / reasoning the reasoner
//	   wants to send at the end of the stream, plus closes any open text
//	   or reasoning sessions.
//	{"type":"error",          "message":"...","code":"..."} -> RUN_ERROR (terminal)
//
// Unknown types are skipped silently with a debug log so reasoner authors
// can iterate without forcing a control-plane upgrade.
type streamingChunk struct {
	Type string `json:"type"`

	// text / reasoning / tool_call_args
	Delta string `json:"delta,omitempty"`

	// reasoning / tool_call_*
	ID              string          `json:"id,omitempty"`
	Name            string          `json:"name,omitempty"`
	ParentMessageID string          `json:"parentMessageId,omitempty"`
	Arguments       json.RawMessage `json:"arguments,omitempty"`

	// tool_call_result
	Content string `json:"content,omitempty"`
	Role    string `json:"role,omitempty"`

	// state / state_delta
	Snapshot any   `json:"snapshot,omitempty"`
	Ops      []any `json:"ops,omitempty"`

	// raw
	Event  any    `json:"event,omitempty"`
	Source string `json:"source,omitempty"`

	// custom
	Value any `json:"value,omitempty"`

	// final
	Data map[string]any `json:"data,omitempty"`

	// error
	Message string `json:"message,omitempty"`
	Code    string `json:"code,omitempty"`
}

// streamingState holds the bookkeeping the dispatcher needs across
// chunks: which text/reasoning sessions are currently open, what tool
// calls have been declared, what assistant message is being built up.
type streamingState struct {
	messageID    string
	textOpen     bool
	textBuf      []byte // accumulates text deltas for the assistant message
	reasoningCtx string // empty if no reasoning context is open
	reasoningSeg string // empty if no reasoning message is open
	toolCalls    []agui.ToolCall
	stateSet     bool
	state        any
}

// runStreamingDispatch consumes the reasoner's NDJSON stream and emits
// AG-UI events as they arrive. Closes the stream when done. Wraps the
// run with TEXT_MESSAGE_START/_END (synthesized lazily on first text
// chunk) and finishes with MESSAGES_SNAPSHOT + RUN_FINISHED — the same
// closing shape buffered reasoners produce, so frontends don't have to
// branch on streaming-vs-buffered.
func runStreamingDispatch(
	ctx context.Context,
	c *gin.Context,
	write func(agui.Event) bool,
	stream io.ReadCloser,
	req agui.RunAgentInput,
	threadID, runID, messageID string,
) {
	defer stream.Close()
	st := &streamingState{messageID: messageID}

	scanner := bufio.NewScanner(stream)
	scanner.Buffer(make([]byte, 0, 64*1024), AGUIStreamingMaxLineBytes)

	for scanner.Scan() {
		select {
		case <-ctx.Done():
			return
		default:
		}
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}
		var ch streamingChunk
		if err := json.Unmarshal(line, &ch); err != nil {
			// One bad chunk shouldn't blow up the run. Surface it as
			// RAW so the frontend at least sees that something garbled
			// went past, and keep going.
			write(agui.RawEvent{
				Event:     map[string]any{"raw": string(line), "decode_error": err.Error()},
				Source:    "agentfield-streaming",
				Timestamp: agui.NowMillis(),
			})
			continue
		}
		if !dispatchChunk(write, st, ch) {
			return
		}
	}
	if err := scanner.Err(); err != nil {
		write(agui.RunError{
			Message:   fmt.Sprintf("read streaming reasoner: %v", err),
			Code:      "ERR_AGENT_STREAM",
			Timestamp: agui.NowMillis(),
		})
		return
	}

	// Stream ended — close any open text/reasoning sessions, emit the
	// canonical close-frames the buffered path would have emitted, and
	// finish the run.
	closeTextSession(write, st)
	closeReasoningSession(write, st)

	assistant := agui.Message{
		ID:        st.messageID,
		Role:      "assistant",
		Content:   string(st.textBuf),
		ToolCalls: st.toolCalls,
	}
	full := append([]agui.Message{}, req.Messages...)
	full = append(full, assistant)
	if !write(agui.MessagesSnapshot{
		Messages:  full,
		Timestamp: agui.NowMillis(),
	}) {
		return
	}

	finished := agui.RunFinished{
		ThreadID:  threadID,
		RunID:     runID,
		Outcome:   &agui.Outcome{Type: "success"},
		Timestamp: agui.NowMillis(),
	}
	if st.stateSet {
		finished.Result = map[string]any{"state": st.state}
	}
	write(finished)
}

// dispatchChunk emits the AG-UI events corresponding to one NDJSON
// chunk. Returns false on a write failure (so the caller stops the loop).
func dispatchChunk(write func(agui.Event) bool, st *streamingState, ch streamingChunk) bool {
	switch ch.Type {
	case "text":
		if ch.Delta == "" {
			return true
		}
		// Reasoning sessions close before the text turn opens — frontends
		// don't expect text chunks interleaved with reasoning.
		if !closeReasoningSession(write, st) {
			return false
		}
		if !st.textOpen {
			if !write(agui.TextMessageStart{
				MessageID: st.messageID,
				Role:      "assistant",
				Timestamp: agui.NowMillis(),
			}) {
				return false
			}
			st.textOpen = true
		}
		st.textBuf = append(st.textBuf, ch.Delta...)
		return write(agui.TextMessageContent{
			MessageID: st.messageID,
			Delta:     ch.Delta,
			Timestamp: agui.NowMillis(),
		})

	case "reasoning":
		if ch.Delta == "" {
			return true
		}
		// Open the outer reasoning context lazily on first chunk.
		if st.reasoningCtx == "" {
			st.reasoningCtx = "reasoning-" + utils.GenerateExecutionID()
			if !write(agui.ReasoningStart{
				MessageID: st.reasoningCtx,
				Timestamp: agui.NowMillis(),
			}) {
				return false
			}
		}
		// Open a per-segment message lazily — the reasoner can send a
		// `reasoning_end` chunk between segments to close one and start
		// the next, but for the simple case (single contiguous thinking
		// block) we batch all deltas into one message.
		if st.reasoningSeg == "" {
			st.reasoningSeg = st.reasoningCtx + "-seg-" + utils.GenerateExecutionID()
			if !write(agui.ReasoningMessageStart{
				MessageID: st.reasoningSeg,
				Role:      "reasoning",
				Timestamp: agui.NowMillis(),
			}) {
				return false
			}
		}
		return write(agui.ReasoningMessageContent{
			MessageID: st.reasoningSeg,
			Delta:     ch.Delta,
			Timestamp: agui.NowMillis(),
		})

	case "reasoning_end":
		// Ends the current reasoning segment (so the next "reasoning"
		// chunk opens a fresh one). Doesn't close the outer context;
		// that happens at stream end or when a "text"/"final" chunk
		// arrives.
		if st.reasoningSeg != "" {
			if !write(agui.ReasoningMessageEnd{
				MessageID: st.reasoningSeg,
				Timestamp: agui.NowMillis(),
			}) {
				return false
			}
			st.reasoningSeg = ""
		}
		return true

	case "tool_call_start":
		if ch.ID == "" || ch.Name == "" {
			return true
		}
		parent := ch.ParentMessageID
		if parent == "" {
			parent = st.messageID
		}
		if !write(agui.ToolCallStart{
			ToolCallID:      ch.ID,
			ToolCallName:    ch.Name,
			ParentMessageID: parent,
			Timestamp:       agui.NowMillis(),
		}) {
			return false
		}
		// Convenience: if the reasoner already has the full arguments
		// at start time (non-streaming-args reasoner), pre-emit them.
		argsStr := ""
		if len(ch.Arguments) > 0 {
			argsStr = string(ch.Arguments)
			if !write(agui.ToolCallArgs{
				ToolCallID: ch.ID,
				Delta:      argsStr,
				Timestamp:  agui.NowMillis(),
			}) {
				return false
			}
		}
		st.toolCalls = append(st.toolCalls, agui.ToolCall{
			ID:   ch.ID,
			Type: "function",
			Function: agui.ToolCallFunction{
				Name:      ch.Name,
				Arguments: argsStr,
			},
		})
		return true

	case "tool_call_args":
		if ch.ID == "" || ch.Delta == "" {
			return true
		}
		// Append to whichever ToolCall.Function.Arguments matches.
		for i := range st.toolCalls {
			if st.toolCalls[i].ID == ch.ID {
				st.toolCalls[i].Function.Arguments += ch.Delta
				break
			}
		}
		return write(agui.ToolCallArgs{
			ToolCallID: ch.ID,
			Delta:      ch.Delta,
			Timestamp:  agui.NowMillis(),
		})

	case "tool_call_end":
		if ch.ID == "" {
			return true
		}
		return write(agui.ToolCallEnd{
			ToolCallID: ch.ID,
			Timestamp:  agui.NowMillis(),
		})

	case "tool_call_result":
		if ch.ID == "" {
			return true
		}
		role := ch.Role
		if role == "" {
			role = "tool"
		}
		return write(agui.ToolCallResult{
			MessageID:  "msg-toolresult-" + ch.ID,
			ToolCallID: ch.ID,
			Content:    ch.Content,
			Role:       role,
			Timestamp:  agui.NowMillis(),
		})

	case "state":
		st.stateSet = true
		st.state = ch.Snapshot
		return write(agui.StateSnapshot{
			Snapshot:  ch.Snapshot,
			Timestamp: agui.NowMillis(),
		})

	case "state_delta":
		if len(ch.Ops) == 0 {
			return true
		}
		return write(agui.StateDelta{
			Delta:     ch.Ops,
			Timestamp: agui.NowMillis(),
		})

	case "step_started":
		if ch.Name == "" {
			return true
		}
		return write(agui.StepStarted{StepName: ch.Name, Timestamp: agui.NowMillis()})

	case "step_finished":
		if ch.Name == "" {
			return true
		}
		return write(agui.StepFinished{StepName: ch.Name, Timestamp: agui.NowMillis()})

	case "raw":
		return write(agui.RawEvent{
			Event:     ch.Event,
			Source:    ch.Source,
			Timestamp: agui.NowMillis(),
		})

	case "custom":
		if ch.Name == "" {
			return true
		}
		return write(agui.CustomEvent{
			Name:      ch.Name,
			Value:     ch.Value,
			Timestamp: agui.NowMillis(),
		})

	case "error":
		// Terminal — emit RUN_ERROR and return false to short-circuit.
		write(agui.RunError{
			Message:   ch.Message,
			Code:      ch.Code,
			Timestamp: agui.NowMillis(),
		})
		return false

	case "final":
		// Treat the data field as a buffered-mode response: extract any
		// not-yet-sent reasoning / tool calls / state / stateDelta and
		// emit them. This lets a streaming reasoner shovel structured
		// trailing fields without re-implementing the buffered logic.
		applyFinal(write, st, ch.Data)
		return true

	default:
		// Unknown chunk type — surface as RAW with a hint so the
		// frontend has visibility, then continue.
		write(agui.RawEvent{
			Event:     map[string]any{"unknown_chunk_type": ch.Type},
			Source:    "agentfield-streaming",
			Timestamp: agui.NowMillis(),
		})
		return true
	}
}

// closeTextSession emits TEXT_MESSAGE_END if a text session is open.
// Returns false on write failure.
func closeTextSession(write func(agui.Event) bool, st *streamingState) bool {
	if !st.textOpen {
		return true
	}
	st.textOpen = false
	return write(agui.TextMessageEnd{
		MessageID: st.messageID,
		Timestamp: agui.NowMillis(),
	})
}

// closeReasoningSession closes any open reasoning message and the outer
// reasoning context. No-op if neither is open.
func closeReasoningSession(write func(agui.Event) bool, st *streamingState) bool {
	if st.reasoningSeg != "" {
		if !write(agui.ReasoningMessageEnd{
			MessageID: st.reasoningSeg,
			Timestamp: agui.NowMillis(),
		}) {
			return false
		}
		st.reasoningSeg = ""
	}
	if st.reasoningCtx != "" {
		if !write(agui.ReasoningEnd{
			MessageID: st.reasoningCtx,
			Timestamp: agui.NowMillis(),
		}) {
			return false
		}
		st.reasoningCtx = ""
	}
	return true
}

// applyFinal lets a streaming reasoner emit one trailing buffered-shape
// envelope to ship any structured fields it didn't send chunk-by-chunk.
// Honors the same field names the buffered path recognizes.
func applyFinal(write func(agui.Event) bool, st *streamingState, data map[string]any) {
	if data == nil {
		return
	}
	// Reasoning (string or list).
	if reasoning := extractReasoning(data); len(reasoning) > 0 {
		// Open a fresh reasoning context if none is open; reuse the
		// open one otherwise.
		ctxID := st.reasoningCtx
		if ctxID == "" {
			ctxID = "reasoning-" + utils.GenerateExecutionID()
			if !write(agui.ReasoningStart{MessageID: ctxID, Timestamp: agui.NowMillis()}) {
				return
			}
			st.reasoningCtx = ctxID
		}
		for i, seg := range reasoning {
			segID := seg.ID
			if segID == "" {
				segID = fmt.Sprintf("%s-final-%d", ctxID, i)
			}
			write(agui.ReasoningMessageStart{MessageID: segID, Role: "reasoning", Timestamp: agui.NowMillis()})
			for _, chunk := range chunkText(seg.Content, AGUITextChunkSize) {
				write(agui.ReasoningMessageContent{MessageID: segID, Delta: chunk, Timestamp: agui.NowMillis()})
			}
			write(agui.ReasoningMessageEnd{MessageID: segID, Timestamp: agui.NowMillis()})
		}
	}
	// Tool calls.
	for _, tc := range extractToolCalls(data) {
		argsJSON, _ := json.Marshal(tc.Arguments)
		argsStr := string(argsJSON)
		write(agui.ToolCallStart{ToolCallID: tc.ID, ToolCallName: tc.Name, ParentMessageID: st.messageID, Timestamp: agui.NowMillis()})
		write(agui.ToolCallArgs{ToolCallID: tc.ID, Delta: argsStr, Timestamp: agui.NowMillis()})
		write(agui.ToolCallEnd{ToolCallID: tc.ID, Timestamp: agui.NowMillis()})
		if tc.HasResult {
			write(agui.ToolCallResult{
				MessageID:  "msg-toolresult-" + tc.ID,
				ToolCallID: tc.ID,
				Content:    stringifyResult(tc.Result),
				Role:       "tool",
				Timestamp:  agui.NowMillis(),
			})
		}
		st.toolCalls = append(st.toolCalls, agui.ToolCall{
			ID:       tc.ID,
			Type:     "function",
			Function: agui.ToolCallFunction{Name: tc.Name, Arguments: argsStr},
		})
	}
	// State.
	if state, ok := extractState(data); ok {
		st.stateSet = true
		st.state = state
		write(agui.StateSnapshot{Snapshot: state, Timestamp: agui.NowMillis()})
	}
	if delta := extractStateDelta(data); delta != nil {
		write(agui.StateDelta{Delta: delta, Timestamp: agui.NowMillis()})
	}
	// Trailing text in `result` — append to any open text turn or open one.
	if r, has := data["result"]; has {
		text := stringifyResult(r)
		if text != "" {
			if !st.textOpen {
				write(agui.TextMessageStart{MessageID: st.messageID, Role: "assistant", Timestamp: agui.NowMillis()})
				st.textOpen = true
			}
			for _, chunk := range chunkText(text, AGUITextChunkSize) {
				write(agui.TextMessageContent{MessageID: st.messageID, Delta: chunk, Timestamp: agui.NowMillis()})
			}
			st.textBuf = append(st.textBuf, text...)
		}
	}
}
