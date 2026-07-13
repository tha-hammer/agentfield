package events

import "github.com/Agent-Field/agentfield/control-plane/pkg/types"

// completedOutboxBuilder builds the terminal-completion outbox record for one agent's
// completion event. Each builder gates on its own AgentNodeID and returns shouldAppend=false
// for executions it does not own.
type completedOutboxBuilder func(*types.Execution) (*types.EventOutboxRecord, bool, error)

// completedOutboxBuilders is the ordered builder registry selected by AgentNodeID
// (MW Phase 3 B4 — Option B, approved 2026-07-13). research.completed and reel.completed ride
// the IDENTICAL same-tx durable-outbox path; onboarding another app's completion event is one
// more entry here — no change to the completeExecution call site.
var completedOutboxBuilders = []completedOutboxBuilder{
	BuildResearchCompletedOutboxRecord,
	BuildReelCompletedOutboxRecord,
}

// BuildCompletedOutboxRecord selects the outbox builder for exec by its AgentNodeID and returns
// the record to append inside the terminal state-write transaction (C-Outbox). shouldAppend is
// false when no builder claims exec — the normal case for agents that emit no completion event.
// The first builder that claims exec wins; builders are mutually exclusive by AgentNodeID, so
// order is not significant beyond that. This is the single call-site helper invoked from
// completeExecution, replacing the direct research-only builder call.
func BuildCompletedOutboxRecord(exec *types.Execution) (*types.EventOutboxRecord, bool, error) {
	for _, build := range completedOutboxBuilders {
		rec, shouldAppend, err := build(exec)
		if err != nil {
			return nil, false, err
		}
		if shouldAppend {
			return rec, true, nil
		}
	}
	return nil, false, nil
}
