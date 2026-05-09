package agui

import "encoding/json"

// RunAgentInput mirrors the canonical RunAgentInputSchema from the AG-UI
// reference SDK (sdks/typescript/packages/core/src/types.ts). The vanilla
// @ag-ui/client HttpAgent — and the CopilotRuntime that wraps it — POSTs
// this exact shape to backends.
//
// We keep all fields permissive (json.RawMessage / any) so unrecognized
// or evolving sub-fields pass through without forcing schema bumps.
type RunAgentInput struct {
	ThreadID       string            `json:"threadId"`
	RunID          string            `json:"runId"`
	ParentRunID    string            `json:"parentRunId,omitempty"`
	State          json.RawMessage   `json:"state,omitempty"`
	Messages       []Message         `json:"messages,omitempty"`
	Tools          []Tool            `json:"tools,omitempty"`
	Context        []ContextItem     `json:"context,omitempty"`
	ForwardedProps json.RawMessage   `json:"forwardedProps,omitempty"`
	Resume         []json.RawMessage `json:"resume,omitempty"`
}

// Message is the canonical AG-UI message envelope (MessageSchema). Role
// drives the discriminated union; we keep optional fields so user/assistant/
// tool messages all round-trip through the same struct.
type Message struct {
	ID         string     `json:"id,omitempty"`
	Role       string     `json:"role"`
	Content    string     `json:"content,omitempty"`
	Name       string     `json:"name,omitempty"`
	ToolCallID string     `json:"toolCallId,omitempty"`
	ToolCalls  []ToolCall `json:"toolCalls,omitempty"`
}

// ToolCall is an assistant-message-attached tool invocation, matching
// ToolCallSchema. The function arguments are a JSON string per OpenAI
// convention.
type ToolCall struct {
	ID       string           `json:"id"`
	Type     string           `json:"type"` // always "function" today
	Function ToolCallFunction `json:"function"`
}

type ToolCallFunction struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
}

// Tool describes a tool the frontend has registered (e.g. via
// useCopilotAction). Reasoners can choose to invoke these by emitting a
// matching TOOL_CALL_* sequence.
type Tool struct {
	Name        string          `json:"name"`
	Description string          `json:"description,omitempty"`
	Parameters  json.RawMessage `json:"parameters,omitempty"` // JSON Schema
	Metadata    json.RawMessage `json:"metadata,omitempty"`
}

// ContextItem is one (description, value) pair from the readables stream
// (e.g. useCopilotReadable). Value is freeform JSON.
type ContextItem struct {
	Description string          `json:"description,omitempty"`
	Value       json.RawMessage `json:"value,omitempty"`
}

// LastUserMessageText returns the trailing user-role message's content,
// which is the conventional "prompt" for chat-style agents. Empty string
// if the trailing message is not user-role or messages is empty.
func (r RunAgentInput) LastUserMessageText() string {
	for i := len(r.Messages) - 1; i >= 0; i-- {
		if r.Messages[i].Role == "user" {
			return r.Messages[i].Content
		}
	}
	return ""
}
