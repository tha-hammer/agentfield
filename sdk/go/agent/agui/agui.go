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
	"fmt"

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
