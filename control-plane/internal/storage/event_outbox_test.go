package storage

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// openLocalStorageAt builds a real on-disk SQLite+BoltDB store under an explicit
// dir (a setupLocalStorage variant) so durability across reopen can be tested.
func openLocalStorageAt(t *testing.T, dir string) (*LocalStorage, context.Context) {
	t.Helper()
	ctx := context.Background()
	cfg := StorageConfig{
		Mode: "local",
		Local: LocalStorageConfig{
			DatabasePath: filepath.Join(dir, "agentfield.db"),
			KVStorePath:  filepath.Join(dir, "agentfield.bolt"),
		},
	}
	ls := NewLocalStorage(LocalStorageConfig{})
	if err := ls.Initialize(ctx, cfg); err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "fts5") {
			t.Skip("sqlite3 compiled without FTS5; skipping outbox tests")
		}
		require.NoError(t, err)
	}
	return ls, ctx
}

func mustCountOutbox(t *testing.T, ls *LocalStorage, ctx context.Context) int64 {
	t.Helper()
	n, err := ls.CountEventOutbox(ctx)
	require.NoError(t, err)
	return n
}

// B1: Outbox append assigns a durable, monotonic seq.
func TestEventOutbox_Append_AssignsMonotonicSeq(t *testing.T) {
	ls, ctx := setupLocalStorage(t)

	seq1, err := ls.AppendEventOutbox(ctx, EventOutboxRecord{EventType: "execution_created", ExecutionID: "e1"})
	require.NoError(t, err)
	require.GreaterOrEqual(t, seq1, int64(1))

	seq2, err := ls.AppendEventOutbox(ctx, EventOutboxRecord{EventType: "execution_completed", ExecutionID: "e1"})
	require.NoError(t, err)
	require.Greater(t, seq2, seq1)

	got, err := ls.ReadEventOutboxAfter(ctx, 0, 10)
	require.NoError(t, err)
	require.Len(t, got, 2)
	require.Equal(t, "execution_created", got[0].EventType)
	require.Equal(t, "e1", got[0].ExecutionID)
}

// B1 property/concurrency: N concurrent appends yield strictly monotonic, unique seqs.
func TestEventOutbox_Append_ConcurrentSeqUnique(t *testing.T) {
	ls, ctx := setupLocalStorage(t)

	const workers, per = 8, 25
	var wg sync.WaitGroup
	var mu sync.Mutex
	seqs := make([]int64, 0, workers*per)

	for w := 0; w < workers; w++ {
		wg.Add(1)
		go func(w int) {
			defer wg.Done()
			for i := 0; i < per; i++ {
				s, err := ls.AppendEventOutbox(ctx, EventOutboxRecord{EventType: "e", ExecutionID: fmt.Sprintf("w%d", w)})
				if err != nil {
					t.Errorf("append: %v", err)
					return
				}
				mu.Lock()
				seqs = append(seqs, s)
				mu.Unlock()
			}
		}(w)
	}
	wg.Wait()

	require.Len(t, seqs, workers*per)
	seen := make(map[int64]struct{}, len(seqs))
	for _, s := range seqs {
		_, dup := seen[s]
		require.Falsef(t, dup, "duplicate seq %d", s)
		seen[s] = struct{}{}
	}
}

// B2: Cursor read returns events after a seq, in order, bounded by limit.
func TestEventOutbox_ReadAfter_ReturnsOrderedTail(t *testing.T) {
	ls, ctx := setupLocalStorage(t)
	for i := 0; i < 5; i++ {
		_, err := ls.AppendEventOutbox(ctx, EventOutboxRecord{EventType: "e", ExecutionID: "x"})
		require.NoError(t, err)
	}

	got, err := ls.ReadEventOutboxAfter(ctx, 2, 10)
	require.NoError(t, err)
	require.Len(t, got, 3)
	require.Equal(t, int64(3), got[0].Seq)
	// ascending
	for i := 1; i < len(got); i++ {
		require.Greater(t, got[i].Seq, got[i-1].Seq)
	}

	empty, err := ls.ReadEventOutboxAfter(ctx, 100, 10)
	require.NoError(t, err)
	require.Empty(t, empty)

	// limit caps the batch
	capped, err := ls.ReadEventOutboxAfter(ctx, 0, 2)
	require.NoError(t, err)
	require.Len(t, capped, 2)

	// negative k treated as 0
	fromNeg, err := ls.ReadEventOutboxAfter(ctx, -5, 10)
	require.NoError(t, err)
	require.Len(t, fromNeg, 5)
}

// B3: Events survive a store restart (durability).
func TestEventOutbox_SurvivesReopen(t *testing.T) {
	dir := t.TempDir()
	ls1, ctx := openLocalStorageAt(t, dir)
	_, err := ls1.AppendEventOutbox(ctx, EventOutboxRecord{EventType: "execution_completed", ExecutionID: "e1"})
	require.NoError(t, err)
	require.NoError(t, ls1.Close(ctx))

	ls2, ctx2 := openLocalStorageAt(t, dir)
	t.Cleanup(func() { _ = ls2.Close(ctx2) })
	got, err := ls2.ReadEventOutboxAfter(ctx2, 0, 10)
	require.NoError(t, err)
	require.Len(t, got, 1)
	// reopen preserves max seq: next append continues, no reuse of seq 1
	next, err := ls2.AppendEventOutbox(ctx2, EventOutboxRecord{EventType: "e", ExecutionID: "e2"})
	require.NoError(t, err)
	require.Greater(t, next, got[0].Seq)
}

// backdateOutboxCreatedAt back-dates a row's created_at by age (mirrors
// backdateExecutionUpdatedAt) so age-based rotation is testable without sleeps.
func backdateOutboxCreatedAt(t *testing.T, ls *LocalStorage, seq int64, age time.Duration) {
	t.Helper()
	gormDB, err := ls.gormWithContext(context.Background())
	require.NoError(t, err)
	require.NoError(t, gormDB.Model(&EventOutboxModel{}).
		Where("seq = ?", seq).
		Update("created_at", time.Now().UTC().Add(-age)).Error)
}

// B5 (storage): cursor get returns 0 for unknown, advance upserts monotonically.
func TestEventOutbox_Cursor_GetAdvance(t *testing.T) {
	ls, ctx := setupLocalStorage(t)

	cur, err := ls.GetOutboxCursor(ctx, "tracer")
	require.NoError(t, err)
	require.Equal(t, int64(0), cur) // unknown consumer -> 0

	require.NoError(t, ls.AdvanceOutboxCursor(ctx, "tracer", 5))
	cur, err = ls.GetOutboxCursor(ctx, "tracer")
	require.NoError(t, err)
	require.Equal(t, int64(5), cur)

	// never regresses (register-at-0 after advancing must not reset)
	require.NoError(t, ls.AdvanceOutboxCursor(ctx, "tracer", 0))
	cur, err = ls.GetOutboxCursor(ctx, "tracer")
	require.NoError(t, err)
	require.Equal(t, int64(5), cur)
}

// B7: Rotation prunes by age (idempotent).
func TestEventOutbox_PruneByAge(t *testing.T) {
	ls, ctx := setupLocalStorage(t)
	old, err := ls.AppendEventOutbox(ctx, EventOutboxRecord{EventType: "e", ExecutionID: "x"})
	require.NoError(t, err)
	backdateOutboxCreatedAt(t, ls, old, 72*time.Hour)
	_, err = ls.AppendEventOutbox(ctx, EventOutboxRecord{EventType: "e", ExecutionID: "y"})
	require.NoError(t, err)

	res, err := ls.PruneEventOutbox(ctx, time.Now().Add(-48*time.Hour), 0)
	require.NoError(t, err)
	require.Equal(t, 1, res.Deleted)

	res2, err := ls.PruneEventOutbox(ctx, time.Now().Add(-48*time.Hour), 0)
	require.NoError(t, err)
	require.Equal(t, 0, res2.Deleted) // idempotent
}

// B8: Rotation caps by row count (keep newest N).
func TestEventOutbox_PruneByCount_KeepsNewest(t *testing.T) {
	ls, ctx := setupLocalStorage(t)
	var seqs []int64
	for i := 0; i < 10; i++ {
		s, err := ls.AppendEventOutbox(ctx, EventOutboxRecord{EventType: "e", ExecutionID: "x"})
		require.NoError(t, err)
		seqs = append(seqs, s)
	}
	_, err := ls.PruneEventOutbox(ctx, time.Time{}, 4)
	require.NoError(t, err)

	got, err := ls.ReadEventOutboxAfter(ctx, 0, 100)
	require.NoError(t, err)
	require.Len(t, got, 4)
	require.Equal(t, seqs[6], got[0].Seq) // oldest survivor is the 7th
}

// B9: Hard cap is loud — pruning unread events is bounded and counted.
func TestClosure_Rotation_UnreadPrune_IsCountedNotSilent(t *testing.T) {
	ls, ctx := setupLocalStorage(t)
	for i := 0; i < 10; i++ {
		_, err := ls.AppendEventOutbox(ctx, EventOutboxRecord{EventType: "e", ExecutionID: "x"})
		require.NoError(t, err)
	}
	require.NoError(t, ls.AdvanceOutboxCursor(ctx, "tracer", 3)) // read only up to seq 3

	res, err := ls.PruneEventOutbox(ctx, time.Time{}, 4) // cap 4 forces dropping unread seq 4..6
	require.NoError(t, err)
	require.Equal(t, int64(4), mustCountOutbox(t, ls, ctx)) // hard cap honored
	require.Equal(t, 3, res.OverflowUnread)                 // seq 4,5,6 unread & pruned -> counted
}

var _ = time.Now
