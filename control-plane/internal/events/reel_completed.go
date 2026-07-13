package events

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/Agent-Field/agentfield/control-plane/pkg/types"
	"github.com/google/uuid"
)

// Producer half of MW Phase 3 Behavior 4 (Option B — cross-app hand-off,
// C-Outbox/C-Notification/C-Correlation, specs/cross-app-handoff.pattern.md §4).
// Mirrors research_completed.go EXACTLY: reel-af's completion rides the identical
// same-tx durable-outbox path as research.completed. The reel-af SDK announce()
// is the reference/test surface; THIS Go builder is the production producer.

const (
	// ReelCompletedEventType is the CloudEvents "type" for the reel completion announcement.
	// Version lives in the type string per cross-app-handoff.pattern.md Decision 2.
	ReelCompletedEventType = "com.silmari.reel.completed.v1"
	// ReelCompletedSource is the CloudEvents "source" — the producing app.
	ReelCompletedSource = "reel-af"
	// ReelAgentNodeID is the AgentField node_id that owns reel.completed (C-Own).
	ReelAgentNodeID = "reel-af"
	// reelCompletedSucceededStatus is the one execution status this emits for.
	reelCompletedSucceededStatus = string(types.ExecutionStatusSucceeded)
	// reelResultRefScheme is the by-reference reel address scheme (C-Notification — the event
	// never carries the video body, only a reference to fetch it on demand).
	reelResultRefScheme = "cp-execution"
)

// ReelCompletedData is the small owner DTO (C-Notification): ids + primitives + a small
// snapshot only — never the mutable reel artifact body. Frozen contract v1:
// {run_id, status, reel_ref, source_execution_id, duration_s, beat_count}.
type ReelCompletedData struct {
	RunID             string  `json:"run_id"`
	Status            string  `json:"status"`
	ReelRef           string  `json:"reel_ref"`
	SourceExecutionID string  `json:"source_execution_id"`
	DurationS         float64 `json:"duration_s"`
	BeatCount         int     `json:"beat_count"`
}

// ReelCompletedEvent is the CloudEvents envelope for reel.completed.
type ReelCompletedEvent struct {
	ID      string            `json:"id"`
	Source  string            `json:"source"`
	Type    string            `json:"type"`
	Subject string            `json:"subject"`
	Time    string            `json:"time"`
	Data    ReelCompletedData `json:"data"`
}

// reelResultSnapshot is the subset of a reel execution's result_payload this builder reads
// (the dict returned by reel-af's reel_research_to_reel reasoner). The video body / narration
// are deliberately not modeled here — they must never be copied into the event (C-Notification).
type reelResultSnapshot struct {
	SourceExecutionID string  `json:"source_execution_id"`
	DurationS         float64 `json:"duration_s"`
	BeatCount         int     `json:"beat_count"`
}

// IsReelCompletionCandidate reports whether agentNodeID is the reel-af node that owns
// reel.completed (C-Own). Callers gate emission on this AND a succeeded status.
func IsReelCompletionCandidate(agentNodeID string) bool {
	return agentNodeID == ReelAgentNodeID
}

// NewReelCompletedEventID mints a fresh, unique-per-emit CloudEvents id.
func NewReelCompletedEventID() string {
	return "ce-" + uuid.NewString()
}

// reelResultRef is the by-reference reel address keyed by executionID.
func reelResultRef(executionID string) string {
	return fmt.Sprintf("%s://%s/result", reelResultRefScheme, executionID)
}

// BuildReelCompletedEvent builds the reel.completed CloudEvent for a terminal succeeded reel
// execution. eventID and eventTime are the non-deterministic envelope fields, injectable so
// callers (and tests) can pin them. Returns an error unless exec reached terminal succeeded —
// this builder emits for succeeded only.
func BuildReelCompletedEvent(exec *types.Execution, eventID string, eventTime time.Time) (*ReelCompletedEvent, error) {
	if exec == nil {
		return nil, fmt.Errorf("nil execution")
	}
	if exec.Status != reelCompletedSucceededStatus {
		return nil, fmt.Errorf("reel.completed is emitted for %q only, got status=%q", reelCompletedSucceededStatus, exec.Status)
	}

	var snapshot reelResultSnapshot
	if len(exec.ResultPayload) > 0 {
		// A malformed or absent result degrades to zero-value source/duration/beats — the
		// correlation keys (subject, execution_id) remain intact regardless.
		_ = json.Unmarshal(exec.ResultPayload, &snapshot)
	}

	return &ReelCompletedEvent{
		ID:      eventID,
		Source:  ReelCompletedSource,
		Type:    ReelCompletedEventType,
		Subject: exec.ExecutionID,
		Time:    eventTime.UTC().Format(time.RFC3339),
		Data: ReelCompletedData{
			RunID:             exec.RunID,
			Status:            reelCompletedSucceededStatus,
			ReelRef:           reelResultRef(exec.ExecutionID),
			SourceExecutionID: snapshot.SourceExecutionID,
			DurationS:         snapshot.DurationS,
			BeatCount:         snapshot.BeatCount,
		},
	}, nil
}

// BuildReelCompletedOutboxRecord is the C-Outbox call-site helper: it gates on
// IsReelCompletionCandidate, builds the envelope, and marshals it into an EventOutboxRecord
// ready to append inside the terminal state-write transaction. shouldAppend is false (with a nil
// record and nil error) when exec doesn't qualify — the normal, expected outcome for every
// non-reel or non-succeeded execution, not an error case. The succeeded-only rule lives solely
// in BuildReelCompletedEvent (single source of truth); its "not succeeded" error is swallowed
// here rather than duplicating the status check. Mirrors BuildResearchCompletedOutboxRecord.
func BuildReelCompletedOutboxRecord(exec *types.Execution) (rec *types.EventOutboxRecord, shouldAppend bool, err error) {
	if exec == nil || !IsReelCompletionCandidate(exec.AgentNodeID) {
		return nil, false, nil
	}

	now := time.Now().UTC()
	event, err := BuildReelCompletedEvent(exec, NewReelCompletedEventID(), now)
	if err != nil {
		return nil, false, nil
	}

	payload, err := json.Marshal(event)
	if err != nil {
		return nil, false, fmt.Errorf("marshal reel.completed event: %w", err)
	}

	return &types.EventOutboxRecord{
		EventType:   ReelCompletedEventType,
		ExecutionID: exec.ExecutionID,
		WorkflowID:  exec.RunID,
		AgentNodeID: exec.AgentNodeID,
		Payload:     string(payload),
		CreatedAt:   now,
	}, true, nil
}
