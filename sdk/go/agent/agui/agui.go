// Package agui provides helpers for AgentField Go reasoners that want to
// surface AG-UI / CopilotKit-compatible Generative UI events through the
// control plane's POST /api/v1/agui/runs/<node>/<reasoner> adapter.
//
// Reasoners opt into the richer event types by returning specific fields
// in their response map; this package builds those fields in the canonical
// shape the control plane expects, so authors don't have to memorize the
// wire contract.
//
// Wire contract (mirrors the Python agentfield.agui module):
//
//   - "result": the human-facing assistant text (used as the
//     TEXT_MESSAGE_CONTENT delta).
//   - "toolCalls": []map{id, name, arguments, result?} — surfaced as
//     TOOL_CALL_START/_ARGS/_END (and _RESULT if `result` is set).
//   - "state": full agent state — emitted as STATE_SNAPSHOT.
//   - "stateDelta": []map{op, path, value} (RFC 6902) — emitted as
//     STATE_DELTA after the snapshot.
//
// See https://docs.ag-ui.com/concepts/events for the upstream protocol.
package agui

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/Agent-Field/agentfield/sdk/go/ai"
	"github.com/Agent-Field/agentfield/sdk/go/harness"
)

// ToolCall builds a single AG-UI tool-call entry. The control plane
// translates each entry into a TOOL_CALL_START/_ARGS/_END triad. If
// `result` is non-nil, TOOL_CALL_RESULT is also emitted so already-executed
// traces (e.g. from ai.ExecuteToolCallLoopResult) render as completed in
// the UI.
//
// `id` may be empty; the control plane synthesizes a stable ID per call.
// Pass an explicit id when correlating with a follow-up tool message
// from a frontend handler.
func ToolCall(id, name string, arguments map[string]any, result any) map[string]any {
	if name == "" {
		// Names are required by the AG-UI schema; an empty name will be
		// silently dropped by the control plane. Surface the bug eagerly.
		return nil
	}
	entry := map[string]any{"name": name}
	if id != "" {
		entry["id"] = id
	}
	if arguments == nil {
		entry["arguments"] = map[string]any{}
	} else {
		entry["arguments"] = arguments
	}
	if result != nil {
		entry["result"] = result
	}
	return entry
}

// ToolCallsFromTrace converts an ai.ToolCallTrace from
// Client.ExecuteToolCallLoopResult into the AG-UI toolCalls list shape.
// Each record becomes an entry with its arguments and the executed
// result (or an {"error":"..."} object if the call failed). Nil or
// empty traces return an empty slice so callers can splat the result
// safely:
//
//	return map[string]any{
//	    "result":    res.Text(),
//	    "toolCalls": agui.ToolCallsFromTrace(res.Trace),
//	}, nil
func ToolCallsFromTrace(trace *ai.ToolCallTrace) []map[string]any {
	if trace == nil || len(trace.Calls) == 0 {
		return []map[string]any{}
	}
	out := make([]map[string]any, 0, len(trace.Calls))
	for i, rec := range trace.Calls {
		entry := map[string]any{
			"id":        fmt.Sprintf("tc-trace-%d", i),
			"name":      rec.ToolName,
			"arguments": rec.Arguments,
		}
		if rec.Arguments == nil {
			entry["arguments"] = map[string]any{}
		}
		switch {
		case rec.Error != "":
			entry["result"] = map[string]any{"error": rec.Error}
		case rec.Result != nil:
			entry["result"] = rec.Result
		}
		out = append(out, entry)
	}
	return out
}

// StateDeltaReplace builds a single RFC 6902 "replace" patch op for a
// stateDelta array. Path must start with "/".
func StateDeltaReplace(path string, value any) (map[string]any, error) {
	if len(path) == 0 || path[0] != '/' {
		return nil, fmt.Errorf("RFC 6902 paths must start with '/' (got %q)", path)
	}
	return map[string]any{"op": "replace", "path": path, "value": value}, nil
}

// ReasoningSegment builds one REASONING_MESSAGE segment for buffered-mode
// emission. Reasoners surface chain-of-thought to CopilotKit's
// "Thinking…" pane by returning a "reasoning" field whose value is a
// list of segments (or plain strings). Each segment becomes a
// REASONING_MESSAGE_START / _CONTENT / _END triad inside a
// REASONING_START / _END boundary.
//
//	return map[string]any{
//	    "result":    "Booked AA-12.",
//	    "reasoning": []any{
//	        agui.ReasoningSegment("Looking up flights..."),
//	        agui.ReasoningSegment("AA-12 is the cheapest non-stop."),
//	    },
//	}, nil
//
// Pass id="" to let the control plane synthesize one.
func ReasoningSegment(content, id string) map[string]any {
	out := map[string]any{"content": content}
	if id != "" {
		out["id"] = id
	}
	return out
}

// Reasoning builds a "reasoning" field value from a mix of plain strings
// and segment maps. Strings are passed through verbatim; mappings are
// shallow-copied. Returns an []any so it slots straight into the
// reasoner response map.
func Reasoning(segments ...any) ([]any, error) {
	out := make([]any, 0, len(segments))
	for _, s := range segments {
		switch v := s.(type) {
		case string:
			if v != "" {
				out = append(out, v)
			}
		case map[string]any:
			cp := make(map[string]any, len(v))
			for k, val := range v {
				cp[k] = val
			}
			out = append(out, cp)
		default:
			return nil, fmt.Errorf("agui.Reasoning: segments must be string or map[string]any (got %T)", s)
		}
	}
	return out, nil
}

// ----------------------------------------------------------------------------
// Streaming chunk builders + serializer.
//
// Reasoners that want live AG-UI events return chunks (built with these
// helpers) from a goroutine and pipe them through SerializeStream into an
// http.ResponseWriter with Content-Type "application/x-ndjson". The
// AgentField control plane sniffs the content-type and dispatches each
// line as a live AG-UI event (see internal/handlers/agui_runs_streaming.go).
// ----------------------------------------------------------------------------

// StreamingContentType is the response content-type a streaming reasoner
// must set so the control plane recognizes it as a live stream.
const StreamingContentType = "application/x-ndjson"

// TextChunk is one piece of streaming assistant text.
func TextChunk(delta string) map[string]any {
	return map[string]any{"type": "text", "delta": delta}
}

// ReasoningChunk is one piece of chain-of-thought rendered in
// CopilotKit's "Thinking…" pane.
func ReasoningChunk(delta string) map[string]any {
	return map[string]any{"type": "reasoning", "delta": delta}
}

// ReasoningEndChunk closes the current reasoning segment so the next
// ReasoningChunk opens a fresh one.
func ReasoningEndChunk() map[string]any {
	return map[string]any{"type": "reasoning_end"}
}

// ToolCallStartChunk opens a tool call. Pass arguments inline if you
// have them all up front; otherwise stream them with ToolCallArgsChunk.
func ToolCallStartChunk(id, name string, arguments map[string]any, parentMessageID string) map[string]any {
	out := map[string]any{"type": "tool_call_start", "id": id, "name": name}
	if arguments != nil {
		out["arguments"] = arguments
	}
	if parentMessageID != "" {
		out["parentMessageId"] = parentMessageID
	}
	return out
}

// ToolCallArgsChunk streams a piece of the tool-call arguments JSON.
func ToolCallArgsChunk(id, delta string) map[string]any {
	return map[string]any{"type": "tool_call_args", "id": id, "delta": delta}
}

// ToolCallEndChunk closes a tool call.
func ToolCallEndChunk(id string) map[string]any {
	return map[string]any{"type": "tool_call_end", "id": id}
}

// ToolCallResultChunk reports a server-side tool result. Use when the
// reasoner already executed the tool and wants the trace to render as
// completed in the UI.
func ToolCallResultChunk(id, content, role string) map[string]any {
	if role == "" {
		role = "tool"
	}
	return map[string]any{"type": "tool_call_result", "id": id, "content": content, "role": role}
}

// StateChunk publishes a full agent state snapshot.
func StateChunk(snapshot any) map[string]any {
	return map[string]any{"type": "state", "snapshot": snapshot}
}

// StateDeltaChunk publishes RFC 6902 patch ops applied incrementally on
// top of the last snapshot.
func StateDeltaChunk(ops []any) map[string]any {
	return map[string]any{"type": "state_delta", "ops": ops}
}

// StepStartedChunk / StepFinishedChunk mark named-step boundaries inside
// the run.
func StepStartedChunk(name string) map[string]any {
	return map[string]any{"type": "step_started", "name": name}
}

func StepFinishedChunk(name string) map[string]any {
	return map[string]any{"type": "step_finished", "name": name}
}

// RawChunk passes a foreign-system event through verbatim.
func RawChunk(event any, source string) map[string]any {
	out := map[string]any{"type": "raw", "event": event}
	if source != "" {
		out["source"] = source
	}
	return out
}

// CustomChunk emits an application-defined event with a name and value.
func CustomChunk(name string, value any) map[string]any {
	out := map[string]any{"type": "custom", "name": name}
	if value != nil {
		out["value"] = value
	}
	return out
}

// FinalChunk packages a trailing buffered envelope. The dispatcher
// applies any toolCalls / state / stateDelta / reasoning / result fields
// in `data` as if from a non-streaming reasoner — useful when the
// reasoner can stream text live but only knows the structured fields at
// the end.
func FinalChunk(data map[string]any) map[string]any {
	return map[string]any{"type": "final", "data": data}
}

// ErrorChunk is a terminal error. The dispatcher emits RUN_ERROR and
// stops the run; later chunks are ignored.
func ErrorChunk(message, code string) map[string]any {
	out := map[string]any{"type": "error", "message": message}
	if code != "" {
		out["code"] = code
	}
	return out
}

// RelayHarnessResult translates a buffered Claude Agent harness result
// (the messages slice on harness.Result) into AG-UI streaming chunks,
// message-by-message. Mirrors the Python SDK's relay_harness_stream.
//
// The Go harness is buffered (it returns a Result after the run finishes)
// so this helper is itself buffered: it walks res.Messages once and
// returns the equivalent chunk slice. Reasoners that want to stream the
// chunks live can either feed the slice into a channel and call
// SerializeStream, or interleave their own custom chunks.
//
// Recognized message shapes (matching the dict form of the Python and
// JS Claude Agent SDK message stream):
//
//   - {type:"assistant", message:{content:[{type:"text", text:"..."}]}}
//     → one TextChunk per text block
//   - {type:"assistant", message:{content:[{type:"thinking", thinking:"..."}]}}
//     → one ReasoningChunk per thinking block
//   - {type:"assistant", message:{content:[{type:"tool_use", id, name, input}]}}
//     → ToolCallStartChunk + ToolCallEndChunk per tool_use block
//   - {type:"user", message:{content:[{type:"tool_result", tool_use_id, content}]}}
//     → ToolCallResultChunk per tool_result block
//   - {type:"result", ...} → skipped (the dispatcher's stream-end logic
//     synthesizes MESSAGES_SNAPSHOT + RUN_FINISHED)
//   - Anything unrecognized is wrapped as a RawChunk.
//
// Note: the Claude Agent SDK buffers per-message, not per-token. True
// per-token streaming requires the raw Anthropic streaming API.
func RelayHarnessResult(res *harness.Result) []map[string]any {
	if res == nil || len(res.Messages) == 0 {
		return nil
	}
	out := make([]map[string]any, 0, len(res.Messages)*2)
	for _, msg := range res.Messages {
		out = append(out, relayHarnessMessage(msg)...)
	}
	return out
}

func relayHarnessMessage(msg map[string]any) []map[string]any {
	if msg == nil {
		return nil
	}
	mtype, _ := msg["type"].(string)
	if mtype == "result" {
		return nil
	}
	if mtype == "system" {
		return []map[string]any{RawChunk(msg, "harness")}
	}
	if mtype != "assistant" && mtype != "user" {
		return []map[string]any{RawChunk(msg, "harness")}
	}

	content := harnessMessageContent(msg)
	if content == nil {
		return []map[string]any{RawChunk(msg, "harness")}
	}
	if s, ok := content.(string); ok {
		if mtype == "assistant" && s != "" {
			return []map[string]any{TextChunk(s)}
		}
		return nil
	}
	blocks, ok := content.([]any)
	if !ok {
		return []map[string]any{RawChunk(msg, "harness")}
	}

	out := make([]map[string]any, 0, len(blocks))
	for _, raw := range blocks {
		block, ok := raw.(map[string]any)
		if !ok {
			continue
		}
		btype, _ := block["type"].(string)
		switch btype {
		case "text":
			text, _ := block["text"].(string)
			if text != "" {
				out = append(out, TextChunk(text))
			}
		case "thinking":
			thinking, _ := block["thinking"].(string)
			if thinking != "" {
				out = append(out, ReasoningChunk(thinking))
			}
		case "tool_use":
			id, _ := block["id"].(string)
			name, _ := block["name"].(string)
			if id == "" || name == "" {
				continue
			}
			input, _ := block["input"].(map[string]any)
			out = append(out, ToolCallStartChunk(id, name, input, ""))
			out = append(out, ToolCallEndChunk(id))
		case "tool_result":
			id, _ := block["tool_use_id"].(string)
			if id == "" {
				continue
			}
			inner := harnessToolResultContent(block["content"])
			out = append(out, ToolCallResultChunk(id, inner, "tool"))
		default:
			out = append(out, RawChunk(block, "harness"))
		}
	}
	return out
}

func harnessMessageContent(msg map[string]any) any {
	if v, ok := msg["content"]; ok {
		return v
	}
	inner, ok := msg["message"].(map[string]any)
	if !ok {
		return nil
	}
	return inner["content"]
}

func harnessToolResultContent(v any) string {
	switch t := v.(type) {
	case string:
		return t
	case []any:
		var b []byte
		for _, item := range t {
			m, ok := item.(map[string]any)
			if !ok {
				continue
			}
			s, _ := m["text"].(string)
			b = append(b, s...)
		}
		return string(b)
	case nil:
		return ""
	default:
		return fmt.Sprintf("%v", t)
	}
}

// SerializeStream consumes a chunks channel (closed by the producer when
// done) and writes one NDJSON line per chunk to w, flushing after each.
// `w` should be an http.ResponseWriter with Content-Type set to
// StreamingContentType. Returns the first write or encode error
// encountered, or nil when the channel closes cleanly.
//
// Typical usage in an HTTP reasoner endpoint:
//
//	w.Header().Set("Content-Type", agui.StreamingContentType)
//	w.WriteHeader(http.StatusOK)
//	chunks := make(chan map[string]any, 8)
//	go produceChunks(ctx, chunks)  // closes chunks when done
//	if err := agui.SerializeStream(ctx, w, chunks); err != nil { ... }
func SerializeStream(ctx context.Context, w io.Writer, chunks <-chan map[string]any) error {
	flusher, _ := w.(interface{ Flush() })
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case ch, ok := <-chunks:
			if !ok {
				return nil
			}
			if err := enc.Encode(ch); err != nil {
				return fmt.Errorf("encode chunk: %w", err)
			}
			if flusher != nil {
				flusher.Flush()
			}
		}
	}
}
