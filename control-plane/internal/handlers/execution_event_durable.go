package handlers

import (
	"context"
	"encoding/json"

	"github.com/Agent-Field/agentfield/control-plane/internal/events"
	"github.com/Agent-Field/agentfield/control-plane/internal/logger"
	"github.com/Agent-Field/agentfield/control-plane/pkg/types"
)

// outboxWriter is the durable-append surface, satisfied by *storage.LocalStorage.
// It is declared here (consumer side) so handlers never depend on the full
// StorageProvider interface for the outbox.
type outboxWriter interface {
	AppendEventOutbox(ctx context.Context, rec types.EventOutboxRecord) (int64, error)
}

// outboxTerminalChecker reports whether a durable terminal event exists for an
// execution. The sync awaiter uses it to recover a dropped terminal signal.
type outboxTerminalChecker interface {
	HasTerminalOutboxEvent(ctx context.Context, executionID string) (bool, error)
}

// persistExecutionEvent appends event to the durable outbox. A failure is logged
// loudly (never silent) but does not block the notification path — the outbox is
// the recovery backstop, and stores that don't implement it (non-LocalStorage)
// simply skip persistence.
func persistExecutionEvent(ctx context.Context, store any, event events.ExecutionEvent) {
	w, ok := store.(outboxWriter)
	if !ok {
		return
	}
	payload, err := json.Marshal(event)
	if err != nil {
		payload = []byte("{}")
	}
	if _, err := w.AppendEventOutbox(ctx, types.EventOutboxRecord{
		EventType:   string(event.Type),
		ExecutionID: event.ExecutionID,
		WorkflowID:  event.WorkflowID,
		AgentNodeID: event.AgentNodeID,
		Payload:     string(payload),
		CreatedAt:   event.Timestamp,
	}); err != nil {
		logger.Logger.Error().Err(err).
			Str("execution_id", event.ExecutionID).
			Str("event_type", string(event.Type)).
			Msg("durable outbox append failed; event not persisted")
	}
}

// publishExecutionEventDurable persists the event and publishes it on BOTH the
// injected per-storage bus (so injected-side consumers such as the sync awaiter
// observe it) and the global bus (so global consumers such as the OTel tracer,
// observability forwarder, telemetry, and cancel dispatcher still observe it).
//
// This closes the split-brain for events that previously reached the global bus
// ONLY (cancel / pause / resume / execution-log): injected-side consumers now
// observe them, and every such event is durably persisted. The global bus is
// retained deliberately — four live consumers depend on it — so this is an
// additive unification rather than a risky removal.
func publishExecutionEventDurable(ctx context.Context, store ExecutionStore, event events.ExecutionEvent) {
	persistExecutionEvent(ctx, store, event)
	if bus := store.GetExecutionEventBus(); bus != nil {
		bus.Publish(event)
	}
	events.GlobalExecutionEventBus.Publish(event)
}
