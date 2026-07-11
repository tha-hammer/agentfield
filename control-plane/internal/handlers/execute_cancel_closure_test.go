package handlers

import (
	"context"
	"testing"
	"time"

	"github.com/Agent-Field/agentfield/control-plane/internal/events"
	"github.com/Agent-Field/agentfield/control-plane/pkg/types"

	"github.com/stretchr/testify/require"
)

// B6 [BLOCKING — closure]: a cancel published through the production durable
// path reaches an injected-side consumer (the awaiter subscribes to the injected
// bus). Previously cancel went to the global bus ONLY, so injected consumers
// missed it. OBSERVE = the injected subscriber's channel.
func TestClosure_Cancel_ReachesInjectedConsumer(t *testing.T) {
	store := newTestExecutionStorage(&types.AgentNode{ID: "n1"})
	injected := store.GetExecutionEventBus().Subscribe("injected-awaiter")

	publishExecutionEventDurable(context.Background(), store, events.ExecutionEvent{
		Type:        events.ExecutionCancelledEvent,
		ExecutionID: "e1",
		Status:      types.ExecutionStatusCancelled,
		Timestamp:   time.Now(),
	})

	select {
	case ev := <-injected:
		require.Equal(t, events.ExecutionCancelledEvent, ev.Type)
		require.Equal(t, "e1", ev.ExecutionID)
	case <-time.After(time.Second):
		t.Fatal("injected consumer did not observe the cancel")
	}
}

// B6 red-at-seam: with the split-brain intact (cancel published to the global
// bus ONLY), the injected-side consumer never receives it.
func TestClosure_Cancel_RedWhenGlobalOnly(t *testing.T) {
	store := newTestExecutionStorage(&types.AgentNode{ID: "n1"})
	injected := store.GetExecutionEventBus().Subscribe("injected-awaiter")

	// OLD behavior: global-only publish.
	events.GlobalExecutionEventBus.Publish(events.ExecutionEvent{
		Type:        events.ExecutionCancelledEvent,
		ExecutionID: "e1",
		Status:      types.ExecutionStatusCancelled,
		Timestamp:   time.Now(),
	})

	select {
	case <-injected:
		t.Fatal("injected consumer should NOT see a global-only publish (red-at-seam)")
	case <-time.After(200 * time.Millisecond):
		// expected: injected side never receives -> the pre-fix red state
	}
}
