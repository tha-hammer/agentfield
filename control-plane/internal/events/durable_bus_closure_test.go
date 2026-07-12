package events_test

import (
	"context"
	"testing"

	"github.com/Agent-Field/agentfield/control-plane/internal/events"
	"github.com/Agent-Field/agentfield/control-plane/pkg/types"

	"github.com/stretchr/testify/require"
)

// B5 [BLOCKING — closure]: a consumer that was disconnected sees every event it
// missed. SOURCE: events published through the production bus while the consumer
// was offline. OBSERVE: the consumer's production cursor read path
// (GetOutboxCursor + ReadEventOutboxAfter). No direct outbox seeding.
func TestClosure_DisconnectedConsumerCatchesUp(t *testing.T) {
	ls, ctx := newOutboxStore(t)
	bus := events.NewDurableExecutionBus(ls)

	// consumer "tracer" is offline; events published through production path:
	require.NoError(t, bus.Publish(ctx, execEvt("e1")))
	require.NoError(t, bus.Publish(ctx, execEvt("e2")))

	cur, err := ls.GetOutboxCursor(ctx, "tracer") // 0 — never connected
	require.NoError(t, err)
	missed, err := ls.ReadEventOutboxAfter(ctx, cur, 100)
	require.NoError(t, err)
	require.Len(t, missed, 2) // both missed events recovered

	require.NoError(t, ls.AdvanceOutboxCursor(ctx, "tracer", missed[len(missed)-1].Seq))
	cur2, err := ls.GetOutboxCursor(ctx, "tracer")
	require.NoError(t, err)
	again, err := ls.ReadEventOutboxAfter(ctx, cur2, 100)
	require.NoError(t, err)
	require.Empty(t, again) // caught up; nothing re-delivered
}

// nonPersistingStore is the RED-AT-SEAM: it accepts appends but persists
// nothing, simulating the durable write being disabled. With it, the same
// catch-up read returns empty — proving the durable write is what closes the gap.
type nonPersistingStore struct{ seq int64 }

func (s *nonPersistingStore) AppendEventOutbox(ctx context.Context, rec types.EventOutboxRecord) (int64, error) {
	s.seq++
	return s.seq, nil // returns a seq but stores nothing
}
func (s *nonPersistingStore) GetOutboxCursor(ctx context.Context, consumerID string) (int64, error) {
	return 0, nil
}
func (s *nonPersistingStore) AdvanceOutboxCursor(ctx context.Context, consumerID string, seq int64) error {
	return nil
}

// readableNonPersistingStore adds a ReadEventOutboxAfter that always returns
// empty, standing in for the disabled durable path's observable result.
func (s *nonPersistingStore) ReadEventOutboxAfter(ctx context.Context, afterSeq int64, limit int) ([]types.EventOutboxRecord, error) {
	return nil, nil
}

// B5 red-at-seam: with the durable write disabled, the reconnecting consumer
// sees nothing — the closure test would be red without the persisted outbox.
func TestClosure_DisconnectedConsumer_RedWhenDurableWriteDisabled(t *testing.T) {
	store := &nonPersistingStore{}
	bus := events.NewDurableExecutionBus(store)
	require.NoError(t, bus.Publish(context.Background(), execEvt("e1")))
	require.NoError(t, bus.Publish(context.Background(), execEvt("e2")))

	missed, err := store.ReadEventOutboxAfter(context.Background(), 0, 100)
	require.NoError(t, err)
	require.Empty(t, missed) // RED: durable write off -> catch-up recovers nothing
}
