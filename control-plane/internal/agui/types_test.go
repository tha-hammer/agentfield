package agui

import "testing"

// TestLastUserMessageText covers the trailing-user-message extractor that
// the handler uses to populate the reasoner's `prompt` convenience field.
func TestLastUserMessageText(t *testing.T) {
	cases := []struct {
		name string
		in   RunAgentInput
		want string
	}{
		{
			name: "empty messages",
			in:   RunAgentInput{},
			want: "",
		},
		{
			name: "single user message",
			in: RunAgentInput{Messages: []Message{
				{Role: "user", Content: "hi"},
			}},
			want: "hi",
		},
		{
			name: "skips trailing assistant turn — picks last user message",
			in: RunAgentInput{Messages: []Message{
				{Role: "user", Content: "first"},
				{Role: "assistant", Content: "ack"},
				{Role: "user", Content: "second"},
				{Role: "assistant", Content: "ack2"},
			}},
			want: "second",
		},
		{
			name: "skips trailing tool message",
			in: RunAgentInput{Messages: []Message{
				{Role: "user", Content: "kick off"},
				{Role: "tool", ToolCallID: "tc1", Content: "tool-output"},
			}},
			want: "kick off",
		},
		{
			name: "no user messages",
			in: RunAgentInput{Messages: []Message{
				{Role: "system", Content: "you are helpful"},
				{Role: "assistant", Content: "ok"},
			}},
			want: "",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.in.LastUserMessageText(); got != tc.want {
				t.Fatalf("LastUserMessageText() = %q, want %q", got, tc.want)
			}
		})
	}
}
