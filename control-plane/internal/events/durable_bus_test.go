package events_test

import (
	"context"
	"errors"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/Agent-Field/agentfield/control-plane/internal/events"
	"github.com/Agent-Field/agentfield/control-plane/internal/storage"
	"github.com/Agent-Field/agentfield/control-plane/pkg/types"

	"github.com/stretchr/testify/require"
)

// newOutboxStore builds a real on-disk SQLite+BoltDB store via storage's public
// API (setupLocalStorage is unexported), so these bus tests exercise the real
// durable path, not a fake.
func newOutboxStore(t *testing.T) (*storage.LocalStorage, context.Context) {
	t.Helper()
	ctx := context.Background()
	dir := t.TempDir()
	cfg := storage.StorageConfig{
		Mode: "local",
		Local: storage.LocalStorageConfig{
			DatabasePath: filepath.Join(dir, "agentfield.db"),
			KVStorePath:  filepath.Join(dir, "agentfield.bolt"),
		},
	}
	ls := storage.NewLocalStorage(storage.LocalStorageConfig{})
	if err := ls.Initialize(ctx, cfg); err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "fts5") {
			t.Skip("sqlite3 compiled without FTS5; skipping durable bus tests")
		}
		require.NoError(t, err)
	}
	t.Cleanup(func() { _ = ls.Close(ctx) })
	return ls, ctx
}

func execEvt(id string) events.ExecutionEvent {
	return events.ExecutionEvent{
		Type:        events.ExecutionCompleted,
		ExecutionID: id,
		Timestamp:   time.Now().UTC(),
	}
}

// B4: Publish is durable-first — event survives a full subscriber buffer.
func TestDurableBus_FullBuffer_DropsLiveButPersists(t *testing.T) {
	ls, ctx := newOutboxStore(t)
	bus := events.NewDurableExecutionBus(ls)
	_ = bus.Subscribe("slow") // never drained -> its 100-slot buffer fills

	for i := 0; i < 100; i++ {
		require.NoError(t, bus.Publish(ctx, execEvt("e")))
	}
	require.NoError(t, bus.Publish(ctx, execEvt("e"))) // 101st: live-dropped, still persisted

	got, err := ls.ReadEventOutboxAfter(ctx, 0, 200)
	require.NoError(t, err)
	require.Len(t, got, 101) // ALL persisted, including the live-dropped one
}

// B4: no subscribers at all -> event still persisted.
func TestDurableBus_NoSubscribers_StillPersists(t *testing.T) {
	ls, ctx := newOutboxStore(t)
	bus := events.NewDurableExecutionBus(ls)
	require.NoError(t, bus.Publish(ctx, execEvt("e1")))
	got, err := ls.ReadEventOutboxAfter(ctx, 0, 10)
	require.NoError(t, err)
	require.Len(t, got, 1)
}

// failingStore forces AppendEventOutbox to error, to prove the durable-first
// contract: Publish returns the error and does NOT do a live-only send.
type failingStore struct{ err error }

func (f failingStore) AppendEventOutbox(ctx context.Context, rec types.EventOutboxRecord) (int64, error) {
	return 0, f.err
}
func (f failingStore) GetOutboxCursor(ctx context.Context, consumerID string) (int64, error) {
	return 0, nil
}
func (f failingStore) AdvanceOutboxCursor(ctx context.Context, consumerID string, seq int64) error {
	return nil
}

// B4: publish error path (store down) surfaces, does not panic, counter climbs.
func TestDurableBus_AppendError_SurfacesAndCounts(t *testing.T) {
	before := events.OutboxAppendFailures()
	bus := events.NewDurableExecutionBus(failingStore{err: errors.New("store down")})
	err := bus.Publish(context.Background(), execEvt("e1"))
	require.Error(t, err)
	require.Equal(t, before+1, events.OutboxAppendFailures())
}
