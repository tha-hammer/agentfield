package handlers

import (
	"context"
	"testing"
	"time"

	"github.com/Agent-Field/agentfield/control-plane/internal/events"
	"github.com/Agent-Field/agentfield/control-plane/pkg/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// MW Phase 3 B4 (Option B) closure test: completeExecution appends reel.completed to the outbox
// IN THE SAME terminal state-write tx for a reel-af node (mirrors the research.completed wiring
// test), and appends nothing for an unrelated agent. Reuses outboxCapableTestStorage from
// execute_outbox_wiring_test.go — that stand-in's append fires only inside
// UpdateExecutionRecordWithOutbox, proving the record rides the completion transaction.

func TestCompleteExecution_EmitsReelCompletedForReelNode(t *testing.T) {
	agent := &types.AgentNode{
		ID:              events.ReelAgentNodeID,
		BaseURL:         "https://example.com",
		HealthStatus:    types.HealthStatusActive,
		LifecycleStatus: types.AgentStatusReady,
	}
	store := &outboxCapableTestStorage{testExecutionStorage: newTestExecutionStorage(agent)}
	controller := newExecutionController(store, nil, nil, time.Second, "")

	now := time.Now().UTC()
	require.NoError(t, store.CreateExecutionRecord(context.Background(), &types.Execution{
		ExecutionID: "exec-reel-1",
		RunID:       "run-reel-1",
		AgentNodeID: agent.ID,
		Status:      types.ExecutionStatusRunning,
		CreatedAt:   now,
		StartedAt:   now,
		UpdatedAt:   now,
	}))

	plan := &preparedExecution{
		exec: &types.Execution{
			ExecutionID: "exec-reel-1",
			RunID:       "run-reel-1",
		},
		agent:  agent,
		target: &parsedTarget{NodeID: agent.ID, TargetName: "reel_research_to_reel"},
	}

	result := []byte(`{"source_execution_id":"exec_abc123","duration_s":12.5,"beat_count":5,"video_path":"never leaks"}`)
	require.NoError(t, controller.completeExecution(context.Background(), plan, result, 10*time.Millisecond))

	require.Equal(t, 1, store.appendCount)
	require.NotNil(t, store.lastRecord)
	assert.Equal(t, events.ReelCompletedEventType, store.lastRecord.EventType)
	assert.Equal(t, "exec-reel-1", store.lastRecord.ExecutionID)
	assert.Contains(t, store.lastRecord.Payload, `"subject":"exec-reel-1"`)
	assert.Contains(t, store.lastRecord.Payload, `"source_execution_id":"exec_abc123"`)
	assert.NotContains(t, store.lastRecord.Payload, "never leaks") // C-Notification: body absent
}

func TestCompleteExecution_SkipsReelCompletedForUnrelatedNode(t *testing.T) {
	agent := &types.AgentNode{
		ID:              "some-other-agent",
		BaseURL:         "https://example.com",
		HealthStatus:    types.HealthStatusActive,
		LifecycleStatus: types.AgentStatusReady,
	}
	store := &outboxCapableTestStorage{testExecutionStorage: newTestExecutionStorage(agent)}
	controller := newExecutionController(store, nil, nil, time.Second, "")

	now := time.Now().UTC()
	require.NoError(t, store.CreateExecutionRecord(context.Background(), &types.Execution{
		ExecutionID: "exec-neither-1",
		RunID:       "run-neither-1",
		AgentNodeID: agent.ID,
		Status:      types.ExecutionStatusRunning,
		CreatedAt:   now,
		StartedAt:   now,
		UpdatedAt:   now,
	}))

	plan := &preparedExecution{
		exec:   &types.Execution{ExecutionID: "exec-neither-1", RunID: "run-neither-1"},
		agent:  agent,
		target: &parsedTarget{NodeID: agent.ID, TargetName: "some_reasoner"},
	}

	require.NoError(t, controller.completeExecution(context.Background(), plan, []byte(`{"ok":true}`), 10*time.Millisecond))

	assert.Equal(t, 0, store.appendCount)
	assert.Nil(t, store.lastRecord)
}
