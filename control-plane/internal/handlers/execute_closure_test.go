package handlers

import (
	"context"
	"encoding/json"
	"sync/atomic"
	"testing"
	"time"

	"github.com/Agent-Field/agentfield/control-plane/pkg/types"

	"github.com/stretchr/testify/require"
)

// b5bStore embeds the fake execution store and overrides GetExecutionRecord to
// return "running" on the pre-subscribe check and "succeeded" thereafter,
// simulating a completion that lands AFTER the awaiter subscribed but whose live
// terminal event was dropped. hasTerminal models the durable outbox signal.
type b5bStore struct {
	*testExecutionStorage
	calls         int32
	terminalAfter int32
	hasTerminal   bool
}

func (s *b5bStore) GetExecutionRecord(ctx context.Context, id string) (*types.Execution, error) {
	n := atomic.AddInt32(&s.calls, 1)
	status := string(types.ExecutionStatusRunning)
	if n >= s.terminalAfter {
		status = string(types.ExecutionStatusSucceeded)
	}
	return &types.Execution{
		ExecutionID:   id,
		Status:        status,
		ResultPayload: json.RawMessage(`{"result":"ok"}`),
	}, nil
}

func (s *b5bStore) HasTerminalOutboxEvent(ctx context.Context, id string) (bool, error) {
	return s.hasTerminal, nil
}

// B5b [BLOCKING — closure]: the sync execute awaiter returns the real result via
// the durable outbox after a dropped terminal event, instead of a false timeout.
// OBSERVE = the awaiter's production return value.
func TestClosure_SyncAwaiter_RecoversDroppedTerminalViaCursor(t *testing.T) {
	store := &b5bStore{
		testExecutionStorage: newTestExecutionStorage(&types.AgentNode{ID: "n1"}),
		terminalAfter:        2, // pre-check(1)=running, timeout-recovery(2)=succeeded
		hasTerminal:          true,
	}
	controller := newExecutionController(store, nil, nil, 90*time.Second, "")

	// No live publish: the awaiter's channel stays silent, its timer fires, and
	// it must recover the terminal from the durable outbox.
	exec, err := controller.waitForExecutionCompletion(context.Background(), "e1", 100*time.Millisecond)
	require.NoError(t, err) // NOT a false timeout
	require.NotNil(t, exec)
	require.Equal(t, string(types.ExecutionStatusSucceeded), exec.Status)
}

// B5b red-at-seam: with no durable terminal event (durable write disabled), the
// dropped live event yields a false timeout — the bug the outbox fixes.
func TestClosure_SyncAwaiter_RedWhenNoDurableTerminal(t *testing.T) {
	store := &b5bStore{
		testExecutionStorage: newTestExecutionStorage(&types.AgentNode{ID: "n1"}),
		terminalAfter:        2,
		hasTerminal:          false, // durable seam off
	}
	controller := newExecutionController(store, nil, nil, 90*time.Second, "")

	_, err := controller.waitForExecutionCompletion(context.Background(), "e1", 100*time.Millisecond)
	require.Error(t, err) // RED: false timeout without the durable recovery
}
