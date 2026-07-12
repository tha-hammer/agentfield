package events

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/Agent-Field/agentfield/control-plane/pkg/types"
	"github.com/google/uuid"
)

// Producer half of INT-02 Behavior 1 (cross-app hand-off, C-Outbox/C-Notification/
// C-Correlation — specs/cross-app-handoff.pattern.md §4). Mirrors the shape built by
// silmari-af-deep-research/research_completed_event.py's build_research_completed,
// which owns the canonical envelope this Go emit must match.

const (
	// ResearchCompletedEventType is the CloudEvents "type" for the completion announcement.
	ResearchCompletedEventType = "research.completed"
	// ResearchCompletedSource is the CloudEvents "source" — the producing app.
	ResearchCompletedSource = "silmari-af-deep-research"
	// ResearchAgentNodeID is the AgentField node_id that owns research.completed (C-Own).
	ResearchAgentNodeID = "meta_deep_research"
	// researchCompletedSucceededStatus is the one execution status this emits for.
	researchCompletedSucceededStatus = string(types.ExecutionStatusSucceeded)
	// researchResultRefScheme is the by-reference result address scheme (C-Notification —
	// the event never carries the document body, only a reference to fetch it on demand).
	researchResultRefScheme = "cp-execution"
)

// ResearchCompletedData is the small owner DTO (C-Notification): ids + primitives + a small
// snapshot only — never the mutable research_package body.
type ResearchCompletedData struct {
	RunID              string  `json:"run_id"`
	Status             string  `json:"status"`
	Title              *string `json:"title"`
	ResultRef          string  `json:"result_ref"`
	ResearchPrompt     *string `json:"research_prompt"`
	ResearchDocumentID string  `json:"research_document_id"`
}

// ResearchCompletedEvent is the CloudEvents envelope for research.completed.
type ResearchCompletedEvent struct {
	ID      string                `json:"id"`
	Source  string                `json:"source"`
	Type    string                `json:"type"`
	Subject string                `json:"subject"`
	Time    string                `json:"time"`
	Data    ResearchCompletedData `json:"data"`
}

// researchResultSnapshot is the subset of an execution's result_payload this builder reads.
// The document body (research_package) is deliberately not modeled here — it must never be
// copied into the event (C-Notification).
type researchResultSnapshot struct {
	Metadata struct {
		Query *string `json:"query"`
		Title *string `json:"title"`
	} `json:"metadata"`
}

// IsResearchCompletionCandidate reports whether agentNodeID is the deep-research node that
// owns research.completed (C-Own). Callers gate emission on this AND a succeeded status.
func IsResearchCompletionCandidate(agentNodeID string) bool {
	return agentNodeID == ResearchAgentNodeID
}

// NewResearchCompletedEventID mints a fresh, unique-per-emit CloudEvents id.
func NewResearchCompletedEventID() string {
	return "ce-" + uuid.NewString()
}

// researchResultRef is the by-reference result address keyed by executionID.
func researchResultRef(executionID string) string {
	return fmt.Sprintf("%s://%s/result", researchResultRefScheme, executionID)
}

// BuildResearchCompletedEvent builds the research.completed CloudEvent for a terminal
// succeeded execution. eventID and eventTime are the non-deterministic envelope fields,
// injectable so callers (and tests) can pin them. Returns an error unless exec reached
// terminal succeeded — this builder emits for succeeded only.
func BuildResearchCompletedEvent(exec *types.Execution, eventID string, eventTime time.Time) (*ResearchCompletedEvent, error) {
	if exec == nil {
		return nil, fmt.Errorf("nil execution")
	}
	if exec.Status != researchCompletedSucceededStatus {
		return nil, fmt.Errorf("research.completed is emitted for %q only, got status=%q", researchCompletedSucceededStatus, exec.Status)
	}

	var snapshot researchResultSnapshot
	if len(exec.ResultPayload) > 0 {
		// A malformed or absent result degrades to nil title/research_prompt — the
		// correlation keys (subject, execution_id) remain intact regardless.
		_ = json.Unmarshal(exec.ResultPayload, &snapshot)
	}

	return &ResearchCompletedEvent{
		ID:      eventID,
		Source:  ResearchCompletedSource,
		Type:    ResearchCompletedEventType,
		Subject: exec.ExecutionID,
		Time:    eventTime.UTC().Format(time.RFC3339),
		Data: ResearchCompletedData{
			RunID:              exec.RunID,
			Status:             researchCompletedSucceededStatus,
			Title:              snapshot.Metadata.Title,
			ResultRef:          researchResultRef(exec.ExecutionID),
			ResearchPrompt:     snapshot.Metadata.Query,
			ResearchDocumentID: exec.ExecutionID,
		},
	}, nil
}

// BuildResearchCompletedOutboxRecord is the C-Outbox call-site helper: it gates on
// IsResearchCompletionCandidate + succeeded, builds the envelope, and marshals it into an
// EventOutboxRecord ready to append inside the terminal state-write transaction. shouldAppend
// is false (with a nil record and nil error) when exec doesn't qualify — this is the normal,
// expected outcome for every non-research or non-terminal execution, not an error case.
func BuildResearchCompletedOutboxRecord(exec *types.Execution) (rec *types.EventOutboxRecord, shouldAppend bool, err error) {
	if exec == nil || !IsResearchCompletionCandidate(exec.AgentNodeID) || exec.Status != researchCompletedSucceededStatus {
		return nil, false, nil
	}

	now := time.Now().UTC()
	event, err := BuildResearchCompletedEvent(exec, NewResearchCompletedEventID(), now)
	if err != nil {
		return nil, false, err
	}

	payload, err := json.Marshal(event)
	if err != nil {
		return nil, false, fmt.Errorf("marshal research.completed event: %w", err)
	}

	return &types.EventOutboxRecord{
		EventType:   ResearchCompletedEventType,
		ExecutionID: exec.ExecutionID,
		WorkflowID:  exec.RunID,
		AgentNodeID: exec.AgentNodeID,
		Payload:     string(payload),
		CreatedAt:   now,
	}, true, nil
}
