// Package agui implements a minimal subset of the AG-UI protocol
// (https://docs.ag-ui.com/concepts/events) so the control plane can emit an
// AG-UI-compatible Server-Sent Events stream that frontends like CopilotKit
// can consume.
//
// Wire format and field shapes are kept faithful to the reference TypeScript
// and Python SDKs at https://github.com/ag-ui-protocol/ag-ui:
//
//   - SSE frames are `data: <json>\n\n` only — no `event:` line. The TS
//     EventEncoder.encodeSSE and the Python EventEncoder._encode_sse both
//     emit exactly this; the discriminator lives in the JSON `type` field.
//   - Event type values are UPPER_SNAKE_CASE (RUN_STARTED, TEXT_MESSAGE_CONTENT, …),
//     matching the EventType enum the reference clients validate against.
//   - `timestamp` is an optional Unix-millisecond integer.
//   - Optional fields are omitted when empty (mirrors `exclude_none=True`).
//
// This is the POC subset — lifecycle + a single TextMessage carrying the
// reasoner's final result. Token-level streaming, tool-call frames, and
// state deltas land in subsequent iterations.
package agui

import (
	"encoding/json"
	"fmt"
	"io"
	"time"
)

// Event is implemented by every AG-UI event payload. The Type method returns
// the canonical AG-UI event name used in the JSON `type` field (e.g.
// "RUN_STARTED"). It is exposed so the SSE writer can name the frame in
// errors and logs without re-marshaling.
type Event interface {
	Type() string
}

// RunStarted signals the beginning of an agent run.
//
// The `input` field is intentionally omitted from this struct: the reference
// schema types it as RunAgentInput (threadId/runId/state/messages/tools/
// context/forwardedProps), not a freeform map. Until we plumb that structured
// shape through, we surface `threadId` and `runId` only — strict clients
// validating against RunAgentInputSchema would reject a freeform map here.
type RunStarted struct {
	ThreadID    string `json:"threadId"`
	RunID       string `json:"runId"`
	ParentRunID string `json:"parentRunId,omitempty"`
	Timestamp   int64  `json:"timestamp,omitempty"`
}

func (RunStarted) Type() string { return "RUN_STARTED" }

func (e RunStarted) MarshalJSON() ([]byte, error) {
	type alias RunStarted
	return json.Marshal(struct {
		Type string `json:"type"`
		alias
	}{Type: e.Type(), alias: alias(e)})
}

// RunFinished signals a successful (or interrupted) run completion.
// Per the reference schema both threadId and runId are required.
type RunFinished struct {
	ThreadID  string   `json:"threadId"`
	RunID     string   `json:"runId"`
	Outcome   *Outcome `json:"outcome,omitempty"`
	Result    any      `json:"result,omitempty"`
	Timestamp int64    `json:"timestamp,omitempty"`
}

func (RunFinished) Type() string { return "RUN_FINISHED" }

func (e RunFinished) MarshalJSON() ([]byte, error) {
	type alias RunFinished
	return json.Marshal(struct {
		Type string `json:"type"`
		alias
	}{Type: e.Type(), alias: alias(e)})
}

// Outcome is a discriminated union: {type: "success"} | {type: "interrupt", interrupts: [...]}.
type Outcome struct {
	Type       string      `json:"type"`
	Interrupts []Interrupt `json:"interrupts,omitempty"`
}

// Interrupt represents a pause point requiring external resolution
// (e.g. human approval). Reserved for HITL flows; not used by the POC.
type Interrupt struct {
	ID     string `json:"id"`
	Reason string `json:"reason,omitempty"`
}

// RunError signals an unrecoverable failure. Terminates the stream.
type RunError struct {
	Message   string `json:"message"`
	Code      string `json:"code,omitempty"`
	Timestamp int64  `json:"timestamp,omitempty"`
}

func (RunError) Type() string { return "RUN_ERROR" }

func (e RunError) MarshalJSON() ([]byte, error) {
	type alias RunError
	return json.Marshal(struct {
		Type string `json:"type"`
		alias
	}{Type: e.Type(), alias: alias(e)})
}

// TextMessageStart opens an assistant text message. Subsequent
// TextMessageContent events with the same messageId carry the body.
type TextMessageStart struct {
	MessageID string `json:"messageId"`
	Role      string `json:"role,omitempty"` // defaults to "assistant" client-side when omitted
	Timestamp int64  `json:"timestamp,omitempty"`
}

func (TextMessageStart) Type() string { return "TEXT_MESSAGE_START" }

func (e TextMessageStart) MarshalJSON() ([]byte, error) {
	type alias TextMessageStart
	return json.Marshal(struct {
		Type string `json:"type"`
		alias
	}{Type: e.Type(), alias: alias(e)})
}

// TextMessageContent carries one chunk of the assistant message body.
// The POC emits a single content event with the full reasoner result;
// once reasoner-side streaming lands, this will be emitted per token chunk.
type TextMessageContent struct {
	MessageID string `json:"messageId"`
	Delta     string `json:"delta"`
	Timestamp int64  `json:"timestamp,omitempty"`
}

func (TextMessageContent) Type() string { return "TEXT_MESSAGE_CONTENT" }

func (e TextMessageContent) MarshalJSON() ([]byte, error) {
	type alias TextMessageContent
	return json.Marshal(struct {
		Type string `json:"type"`
		alias
	}{Type: e.Type(), alias: alias(e)})
}

// TextMessageEnd closes a text message.
type TextMessageEnd struct {
	MessageID string `json:"messageId"`
	Timestamp int64  `json:"timestamp,omitempty"`
}

func (TextMessageEnd) Type() string { return "TEXT_MESSAGE_END" }

func (e TextMessageEnd) MarshalJSON() ([]byte, error) {
	type alias TextMessageEnd
	return json.Marshal(struct {
		Type string `json:"type"`
		alias
	}{Type: e.Type(), alias: alias(e)})
}

// ToolCallStart opens a tool-call frame. CopilotKit pattern-matches
// `toolCallName` against `useCopilotAction({name, render})` registrations
// to drive Generative UI — there is no separate "render" event.
type ToolCallStart struct {
	ToolCallID      string `json:"toolCallId"`
	ToolCallName    string `json:"toolCallName"`
	ParentMessageID string `json:"parentMessageId,omitempty"`
	Timestamp       int64  `json:"timestamp,omitempty"`
}

func (ToolCallStart) Type() string { return "TOOL_CALL_START" }

func (e ToolCallStart) MarshalJSON() ([]byte, error) {
	type alias ToolCallStart
	return json.Marshal(struct {
		Type string `json:"type"`
		alias
	}{Type: e.Type(), alias: alias(e)})
}

// ToolCallArgs streams a chunk of the tool-call arguments JSON. Frontends
// concatenate deltas to assemble the full arguments object before invoking
// the action handler.
type ToolCallArgs struct {
	ToolCallID string `json:"toolCallId"`
	Delta      string `json:"delta"`
	Timestamp  int64  `json:"timestamp,omitempty"`
}

func (ToolCallArgs) Type() string { return "TOOL_CALL_ARGS" }

func (e ToolCallArgs) MarshalJSON() ([]byte, error) {
	type alias ToolCallArgs
	return json.Marshal(struct {
		Type string `json:"type"`
		alias
	}{Type: e.Type(), alias: alias(e)})
}

// ToolCallEnd closes a tool-call frame.
type ToolCallEnd struct {
	ToolCallID string `json:"toolCallId"`
	Timestamp  int64  `json:"timestamp,omitempty"`
}

func (ToolCallEnd) Type() string { return "TOOL_CALL_END" }

func (e ToolCallEnd) MarshalJSON() ([]byte, error) {
	type alias ToolCallEnd
	return json.Marshal(struct {
		Type string `json:"type"`
		alias
	}{Type: e.Type(), alias: alias(e)})
}

// ToolCallResult delivers the outcome of a server-side tool call. For
// frontend-handled tools (via useCopilotAction), the result instead arrives
// as the next inbound POST's trailing tool-role message — no TOOL_CALL_RESULT
// event is emitted by the backend in that flow.
type ToolCallResult struct {
	MessageID  string `json:"messageId"`
	ToolCallID string `json:"toolCallId"`
	Content    string `json:"content"`
	Role       string `json:"role,omitempty"`
	Timestamp  int64  `json:"timestamp,omitempty"`
}

func (ToolCallResult) Type() string { return "TOOL_CALL_RESULT" }

func (e ToolCallResult) MarshalJSON() ([]byte, error) {
	type alias ToolCallResult
	return json.Marshal(struct {
		Type string `json:"type"`
		alias
	}{Type: e.Type(), alias: alias(e)})
}

// MessagesSnapshot publishes the full conversation after a turn so clients
// can refresh their canonical thread state. CopilotKit's in-memory runtime
// derives persisted history from the trailing snapshot.
type MessagesSnapshot struct {
	Messages  []Message `json:"messages"`
	Timestamp int64     `json:"timestamp,omitempty"`
}

func (MessagesSnapshot) Type() string { return "MESSAGES_SNAPSHOT" }

func (e MessagesSnapshot) MarshalJSON() ([]byte, error) {
	type alias MessagesSnapshot
	return json.Marshal(struct {
		Type string `json:"type"`
		alias
	}{Type: e.Type(), alias: alias(e)})
}

// StateSnapshot publishes the agent's full shared state — the value
// `useCoAgent({ state })` reads on the frontend. Reasoners opt in by
// returning a top-level `state` field.
type StateSnapshot struct {
	Snapshot  any   `json:"snapshot"`
	Timestamp int64 `json:"timestamp,omitempty"`
}

func (StateSnapshot) Type() string { return "STATE_SNAPSHOT" }

func (e StateSnapshot) MarshalJSON() ([]byte, error) {
	type alias StateSnapshot
	return json.Marshal(struct {
		Type string `json:"type"`
		alias
	}{Type: e.Type(), alias: alias(e)})
}

// StateDelta carries an RFC 6902 JSON Patch document applied incrementally
// to the previously-emitted snapshot. Optional alternative to repeatedly
// emitting full snapshots.
type StateDelta struct {
	Delta     []any `json:"delta"`
	Timestamp int64 `json:"timestamp,omitempty"`
}

func (StateDelta) Type() string { return "STATE_DELTA" }

func (e StateDelta) MarshalJSON() ([]byte, error) {
	type alias StateDelta
	return json.Marshal(struct {
		Type string `json:"type"`
		alias
	}{Type: e.Type(), alias: alias(e)})
}

// NowMillis returns the current Unix time in milliseconds. Wrapped so tests
// can replace it. Milliseconds match the JS `Date.now()` convention that
// AG-UI clients are most likely to interpret correctly.
var NowMillis = func() int64 { return time.Now().UnixMilli() }

// WriteSSE writes one AG-UI event to w in the canonical wire format used by
// the reference TS and Python encoders:
//
//	data: <json>
//
// (followed by a blank line). The discriminator is in the JSON `type` field,
// not in an SSE `event:` line — clients dispatch on the JSON `type`. Caller
// is responsible for flushing.
func WriteSSE(w io.Writer, ev Event) error {
	payload, err := json.Marshal(ev)
	if err != nil {
		return fmt.Errorf("marshal %s: %w", ev.Type(), err)
	}
	if _, err := fmt.Fprintf(w, "data: %s\n\n", payload); err != nil {
		return fmt.Errorf("write %s: %w", ev.Type(), err)
	}
	return nil
}
