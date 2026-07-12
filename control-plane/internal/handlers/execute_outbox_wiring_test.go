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

// outboxCapableTestStorage adds ExecutionOutboxUpdater to testExecutionStorage so these tests
// can prove completeExecution builds and passes the research.completed outbox record correctly
// gated (research node + succeeded only), without needing a real transactional DB — that
// atomicity guarantee is proven separately by internal/storage's closure tests.
type outboxCapableTestStorage struct {
	*testExecutionStorage
	appendCount int
	lastRecord  *types.EventOutboxRecord
}

func (s *outboxCapableTestStorage) UpdateExecutionRecordWithOutbox(
	ctx context.Context,
	executionID string,
	updater func(*types.Execution) (*types.Execution, error),
	outboxBuilder func(updated *types.Execution) (*types.EventOutboxRecord, bool, error),
) (*types.Execution, error) {
	updated, err := s.testExecutionStorage.UpdateExecutionRecord(ctx, executionID, updater)
	if err != nil {
		return nil, err
	}
	if outboxBuilder != nil {
		rec, shouldAppend, err := outboxBuilder(updated)
		if err != nil {
			return nil, err
		}
		if shouldAppend {
			s.appendCount++
			s.lastRecord = rec
		}
	}
	return updated, nil
}

func TestCompleteExecution_EmitsResearchCompletedForResearchNode(t *testing.T) {
	agent := &types.AgentNode{
		ID:              events.ResearchAgentNodeID,
		BaseURL:         "https://example.com",
		HealthStatus:    types.HealthStatusActive,
		LifecycleStatus: types.AgentStatusReady,
	}
	store := &outboxCapableTestStorage{testExecutionStorage: newTestExecutionStorage(agent)}
	controller := newExecutionController(store, nil, nil, time.Second, "")

	now := time.Now().UTC()
	require.NoError(t, store.CreateExecutionRecord(context.Background(), &types.Execution{
		ExecutionID: "exec-research-1",
		RunID:       "run-research-1",
		AgentNodeID: agent.ID,
		Status:      types.ExecutionStatusRunning,
		CreatedAt:   now,
		StartedAt:   now,
		UpdatedAt:   now,
	}))

	plan := &preparedExecution{
		exec: &types.Execution{
			ExecutionID: "exec-research-1",
			RunID:       "run-research-1",
		},
		agent:  agent,
		target: &parsedTarget{NodeID: agent.ID, TargetName: "execute_deep_research"},
	}

	result := []byte(`{"metadata":{"query":"Q","title":"T"},"research_package":{"body":"never leaks"}}`)
	require.NoError(t, controller.completeExecution(context.Background(), plan, result, 10*time.Millisecond))

	require.Equal(t, 1, store.appendCount)
	require.NotNil(t, store.lastRecord)
	assert.Equal(t, events.ResearchCompletedEventType, store.lastRecord.EventType)
	assert.Equal(t, "exec-research-1", store.lastRecord.ExecutionID)
	assert.Contains(t, store.lastRecord.Payload, `"subject":"exec-research-1"`)
	assert.NotContains(t, store.lastRecord.Payload, "never leaks") // C-Notification: body absent
}

func TestCompleteExecution_SkipsResearchCompletedForOtherNodes(t *testing.T) {
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
		ExecutionID: "exec-other-1",
		RunID:       "run-other-1",
		AgentNodeID: agent.ID,
		Status:      types.ExecutionStatusRunning,
		CreatedAt:   now,
		StartedAt:   now,
		UpdatedAt:   now,
	}))

	plan := &preparedExecution{
		exec: &types.Execution{
			ExecutionID: "exec-other-1",
			RunID:       "run-other-1",
		},
		agent:  agent,
		target: &parsedTarget{NodeID: agent.ID, TargetName: "some_reasoner"},
	}

	require.NoError(t, controller.completeExecution(context.Background(), plan, []byte(`{"ok":true}`), 10*time.Millisecond))

	assert.Equal(t, 0, store.appendCount)
	assert.Nil(t, store.lastRecord)
}
