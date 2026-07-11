package handlers

import (
	"context"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/Agent-Field/agentfield/control-plane/internal/config"
	"github.com/Agent-Field/agentfield/control-plane/internal/storage"
	"github.com/Agent-Field/agentfield/control-plane/pkg/types"

	"github.com/stretchr/testify/require"
)

func newRotatorStore(t *testing.T) (*storage.LocalStorage, context.Context) {
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
			t.Skip("sqlite3 compiled without FTS5; skipping rotator tests")
		}
		require.NoError(t, err)
	}
	t.Cleanup(func() { _ = ls.Close(ctx) })
	return ls, ctx
}

func seedOutbox(t *testing.T, ls *storage.LocalStorage, ctx context.Context, n int) {
	t.Helper()
	for i := 0; i < n; i++ {
		_, err := ls.AppendEventOutbox(ctx, types.EventOutboxRecord{EventType: "e", ExecutionID: "x"})
		require.NoError(t, err)
	}
}

func countOutbox(t *testing.T, ls *storage.LocalStorage, ctx context.Context) int64 {
	t.Helper()
	n, err := ls.CountEventOutbox(ctx)
	require.NoError(t, err)
	return n
}

// B10: pruneOnce enforces caps (directly-callable seam, no ticker).
func TestEventOutboxRotator_PruneOnce_EnforcesCaps(t *testing.T) {
	ls, ctx := newRotatorStore(t)
	seedOutbox(t, ls, ctx, 20)

	r := NewEventOutboxRotator(ls, config.EventOutboxConfig{Enabled: true, RetentionMaxRows: 5, PruneInterval: time.Minute})
	require.NoError(t, r.pruneOnce(ctx))
	require.Equal(t, int64(5), countOutbox(t, ls, ctx))
}

// B10: Enabled=false -> Start is a no-op; nothing pruned.
func TestEventOutboxRotator_Disabled_NoOp(t *testing.T) {
	ls, ctx := newRotatorStore(t)
	seedOutbox(t, ls, ctx, 10)

	r := NewEventOutboxRotator(ls, config.EventOutboxConfig{Enabled: false, RetentionMaxRows: 2, PruneInterval: time.Minute})
	require.NoError(t, r.Start(ctx))
	require.NoError(t, r.Stop())
	require.Equal(t, int64(10), countOutbox(t, ls, ctx)) // untouched
}

// B10: Start then Stop returns without leaking a goroutine (WaitGroup drained).
func TestEventOutboxRotator_Stop_NoLeak(t *testing.T) {
	ls, ctx := newRotatorStore(t)
	r := NewEventOutboxRotator(ls, config.EventOutboxConfig{Enabled: true, RetentionMaxRows: 5, PruneInterval: time.Minute})
	require.NoError(t, r.Start(ctx))
	done := make(chan struct{})
	go func() { _ = r.Stop(); close(done) }()
	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("Stop did not return; goroutine leaked")
	}
}
