package events

import (
	"context"
	"encoding/json"
	"sync/atomic"

	"github.com/Agent-Field/agentfield/control-plane/internal/logger"
	"github.com/Agent-Field/agentfield/control-plane/pkg/types"
)

// OutboxStore is the durable backstop the DurableExecutionBus writes through.
// It is declared here (the consuming package) per Go convention and is
// satisfied by *storage.LocalStorage. Only the methods the bus needs appear.
type OutboxStore interface {
	AppendEventOutbox(ctx context.Context, rec types.EventOutboxRecord) (int64, error)
	GetOutboxCursor(ctx context.Context, consumerID string) (int64, error)
	AdvanceOutboxCursor(ctx context.Context, consumerID string, seq int64) error
}

// outboxAppendFailedTotal counts publishes whose durable append failed. Exposed
// via OutboxAppendFailures for tests and observability (metric
// outbox_append_failed_total).
var outboxAppendFailedTotal atomic.Int64

// OutboxAppendFailures returns the number of publishes whose durable append
// failed since process start.
func OutboxAppendFailures() int64 { return outboxAppendFailedTotal.Load() }

// DurableExecutionBus persists every event before the best-effort live fan-out,
// so a full subscriber buffer or a process restart never loses an event. It
// composes ExecutionEventBus for the live path rather than forking its drop
// logic — durability is hidden in Publish (a deep module: Publish/Subscribe).
type DurableExecutionBus struct {
	store OutboxStore
	live  *ExecutionEventBus
}

// NewDurableExecutionBus wraps store with a fresh live fan-out bus.
func NewDurableExecutionBus(store OutboxStore) *DurableExecutionBus {
	return &DurableExecutionBus{store: store, live: NewExecutionEventBus()}
}

// Subscribe registers a live subscriber. It also registers the consumer's
// durable cursor at 0 (best-effort) so rotation's overflow accounting can see
// an active-but-behind consumer's lag rather than under-counting it.
func (b *DurableExecutionBus) Subscribe(subscriberID string) chan ExecutionEvent {
	if err := b.store.AdvanceOutboxCursor(context.Background(), subscriberID, 0); err != nil {
		logger.Logger.Debug().Err(err).Msgf("[DurableExecutionBus] cursor register failed for %s", subscriberID)
	}
	return b.live.Subscribe(subscriberID)
}

// Unsubscribe removes the live subscriber. The durable cursor is retained so a
// reconnecting consumer resumes where it left off.
func (b *DurableExecutionBus) Unsubscribe(subscriberID string) {
	b.live.Unsubscribe(subscriberID)
}

// Publish persists the event first. If the durable append fails it increments
// outbox_append_failed_total and RETURNS the error WITHOUT doing a live-only
// send, so a store failure is never masked as a successful publish. Only on a
// successful append does it perform the existing non-blocking live fan-out.
func (b *DurableExecutionBus) Publish(ctx context.Context, event ExecutionEvent) error {
	payload, err := json.Marshal(event)
	if err != nil {
		return err
	}
	rec := types.EventOutboxRecord{
		EventType:   string(event.Type),
		ExecutionID: event.ExecutionID,
		WorkflowID:  event.WorkflowID,
		AgentNodeID: event.AgentNodeID,
		Payload:     string(payload),
		CreatedAt:   event.Timestamp,
	}
	if _, err := b.store.AppendEventOutbox(ctx, rec); err != nil {
		outboxAppendFailedTotal.Add(1)
		logger.Logger.Error().Err(err).
			Str("execution_id", event.ExecutionID).
			Str("event_type", string(event.Type)).
			Msg("[DurableExecutionBus] durable append failed; event not published")
		return err
	}
	b.live.Publish(event)
	return nil
}

// Live exposes the underlying live bus for consumers that only need the
// best-effort channel fan-out (e.g. existing SSE subscribers).
func (b *DurableExecutionBus) Live() *ExecutionEventBus { return b.live }
