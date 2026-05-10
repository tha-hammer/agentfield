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
