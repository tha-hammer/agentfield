package storage

import (
	"context"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/Agent-Field/agentfield/control-plane/pkg/types"

	"github.com/stretchr/testify/require"
)

// INT-02 Behavior 1 (Go half) — C-Outbox: the research.completed append must land in the SAME
// transaction as the terminal state write. These tests exercise a real, on-disk, ACID-transactional
// store (SQLite via database/sql — the same substrate BeginTx/Commit/Rollback semantics apply to
// postgres in production), not a mock, so rollback here is a genuine transaction rollback — the
// same mechanism guaranteeing atomicity in postgres. Mirrors internal/events/durable_bus_test.go's
// newOutboxStore pattern (real store, not a fake).
func newOutboxTxTestStorage(t *testing.T) (*LocalStorage, context.Context) {
	t.Helper()
	ctx := context.Background()
	dir := t.TempDir()
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
			t.Skip("sqlite3 compiled without FTS5; skipping outbox tx tests")
		}
		require.NoError(t, err)
	}
	t.Cleanup(func() { _ = ls.Close(ctx) })
	return ls, ctx
}

func seedRunningExecution(t *testing.T, ls *LocalStorage, ctx context.Context, executionID, agentNodeID string) {
	t.Helper()
	now := time.Now().UTC()
	require.NoError(t, ls.CreateExecutionRecord(ctx, &types.Execution{
		ExecutionID: executionID,
		RunID:       "run-" + executionID,
		AgentNodeID: agentNodeID,
		ReasonerID:  "execute_deep_research",
		NodeID:      agentNodeID,
		Status:      string(types.ExecutionStatusRunning),
		StartedAt:   now,
		CreatedAt:   now,
		UpdatedAt:   now,
	}))
}

func succeededUpdater(result []byte) func(*types.Execution) (*types.Execution, error) {
	return func(current *types.Execution) (*types.Execution, error) {
		current.Status = string(types.ExecutionStatusSucceeded)
		current.ResultPayload = result
		return current, nil
	}
}

// TestUpdateExecutionRecordWithOutbox_AppendsInSameTx is the GREEN case: the terminal succeeded
// write and the outbox append both land, in the same tx, with a monotonic Seq.
func TestUpdateExecutionRecordWithOutbox_AppendsInSameTx(t *testing.T) {
	ls, ctx := newOutboxTxTestStorage(t)
	seedRunningExecution(t, ls, ctx, "exec-outbox-1", "meta_deep_research")

	updated, err := ls.UpdateExecutionRecordWithOutbox(ctx, "exec-outbox-1",
		succeededUpdater([]byte(`{"metadata":{"query":"Q","title":"T"}}`)),
		func(updated *types.Execution) (*types.EventOutboxRecord, bool, error) {
			return &types.EventOutboxRecord{
				EventType:   "research.completed",
				ExecutionID: updated.ExecutionID,
				WorkflowID:  updated.RunID,
				AgentNodeID: updated.AgentNodeID,
				Payload:     `{"type":"research.completed"}`,
			}, true, nil
		},
	)
	require.NoError(t, err)
	require.Equal(t, string(types.ExecutionStatusSucceeded), updated.Status)

	rows, err := ls.ReadEventOutboxAfter(ctx, 0, 100)
	require.NoError(t, err)
	require.Len(t, rows, 1)
	require.Equal(t, "research.completed", rows[0].EventType)
	require.Equal(t, "exec-outbox-1", rows[0].ExecutionID)
	require.Greater(t, rows[0].Seq, int64(0)) // monotonic seq assigned
}

// TestUpdateExecutionRecordWithOutbox_SkipsWhenShouldAppendFalse proves the gate is real: a
// builder that declines to append leaves the outbox untouched while the state write still commits.
func TestUpdateExecutionRecordWithOutbox_SkipsWhenShouldAppendFalse(t *testing.T) {
	ls, ctx := newOutboxTxTestStorage(t)
	seedRunningExecution(t, ls, ctx, "exec-outbox-2", "some-other-agent")

	updated, err := ls.UpdateExecutionRecordWithOutbox(ctx, "exec-outbox-2",
		succeededUpdater([]byte(`{}`)),
		func(updated *types.Execution) (*types.EventOutboxRecord, bool, error) {
			return nil, false, nil // not a research node — the real gate this mirrors
		},
	)
	require.NoError(t, err)
	require.Equal(t, string(types.ExecutionStatusSucceeded), updated.Status)

	rows, err := ls.ReadEventOutboxAfter(ctx, 0, 100)
	require.NoError(t, err)
	require.Empty(t, rows)
}

// TestUpdateExecutionRecordWithOutbox_OutboxBuilderErrorRollsBackStateWrite is the atomicity
// proof (rollback direction): if the outbox side fails, the terminal state write rolls back with
// it — the run is never left marked-succeeded-without-its-event.
func TestUpdateExecutionRecordWithOutbox_OutboxBuilderErrorRollsBackStateWrite(t *testing.T) {
	ls, ctx := newOutboxTxTestStorage(t)
	seedRunningExecution(t, ls, ctx, "exec-outbox-3", "meta_deep_research")

	_, err := ls.UpdateExecutionRecordWithOutbox(ctx, "exec-outbox-3",
		succeededUpdater([]byte(`{}`)),
		func(updated *types.Execution) (*types.EventOutboxRecord, bool, error) {
			return nil, false, errAppendBoom
		},
	)
	require.Error(t, err)

	// The state write rolled back with the failed outbox build: status is still "running".
	record, err := ls.GetExecutionRecord(ctx, "exec-outbox-3")
	require.NoError(t, err)
	require.Equal(t, string(types.ExecutionStatusRunning), record.Status)

	rows, err := ls.ReadEventOutboxAfter(ctx, 0, 100)
	require.NoError(t, err)
	require.Empty(t, rows)
}

// TestAppendEventOutboxTx_RollbackDiscardsAppend is the direct atomicity primitive proof: an
// append on a tx that is rolled back (rather than committed) leaves no trace — the SAME mechanism
// UpdateExecutionRecordWithOutbox relies on when the surrounding state write fails.
func TestAppendEventOutboxTx_RollbackDiscardsAppend(t *testing.T) {
	ls, ctx := newOutboxTxTestStorage(t)

	tx, err := ls.requireSQLDB().BeginTx(ctx, nil)
	require.NoError(t, err)

	seq, err := ls.AppendEventOutboxTx(ctx, tx, EventOutboxRecord{
		EventType:   "research.completed",
		ExecutionID: "exec-rollback",
		Payload:     `{"type":"research.completed"}`,
	})
	require.NoError(t, err)
	require.Greater(t, seq, int64(0))

	require.NoError(t, tx.Rollback())

	count, err := ls.CountEventOutbox(ctx)
	require.NoError(t, err)
	require.Equal(t, int64(0), count) // RED: rollback discards the append, exactly like the state write
}

// TestUpdateExecutionRecord_RedAtSeam_PlainPathNeverAppendsToOutbox is the red-at-seam baseline:
// before this behavior, calling the plain UpdateExecutionRecord (no outbox hook) on a qualifying
// research execution never appends research.completed — proving the gap this behavior closes.
func TestUpdateExecutionRecord_RedAtSeam_PlainPathNeverAppendsToOutbox(t *testing.T) {
	ls, ctx := newOutboxTxTestStorage(t)
	seedRunningExecution(t, ls, ctx, "exec-outbox-4", "meta_deep_research")

	updated, err := ls.UpdateExecutionRecord(ctx, "exec-outbox-4", succeededUpdater([]byte(`{}`)))
	require.NoError(t, err)
	require.Equal(t, string(types.ExecutionStatusSucceeded), updated.Status)

	rows, err := ls.ReadEventOutboxAfter(ctx, 0, 100)
	require.NoError(t, err)
	require.Empty(t, rows) // the connector is "disabled" on this path — nothing appended
}

var errAppendBoom = errBoom{}

type errBoom struct{}

func (errBoom) Error() string { return "boom: outbox build failed" }

// TestAppendEventOutboxTx_DefaultsEmptyPayload proves an empty Payload defaults to "{}" rather
// than persisting an empty string (event_outbox.payload is NOT NULL DEFAULT '{}').
func TestAppendEventOutboxTx_DefaultsEmptyPayload(t *testing.T) {
	ls, ctx := newOutboxTxTestStorage(t)

	tx, err := ls.requireSQLDB().BeginTx(ctx, nil)
	require.NoError(t, err)

	_, err = ls.AppendEventOutboxTx(ctx, tx, EventOutboxRecord{
		EventType:   "research.completed",
		ExecutionID: "exec-empty-payload",
	})
	require.NoError(t, err)
	require.NoError(t, tx.Commit())

	rows, err := ls.ReadEventOutboxAfter(ctx, 0, 100)
	require.NoError(t, err)
	require.Len(t, rows, 1)
	require.Equal(t, "{}", rows[0].Payload)
}

// TestAppendEventOutboxTx_PostgresRETURNINGBranch drives the postgres code path (INSERT ...
// RETURNING seq) directly: mattn/go-sqlite3 (v1.14+) supports RETURNING, so wrapping the SAME
// real, on-disk *sql.Tx with mode="postgres" exercises that branch's actual SQL + Scan behavior
// without needing a live postgres connection — the branch is selected purely on tx.mode, so this
// is the real code, not a stand-in.
func TestAppendEventOutboxTx_PostgresRETURNINGBranch(t *testing.T) {
	ls, ctx := newOutboxTxTestStorage(t)

	rawTx, err := ls.requireSQLDB().BeginTx(ctx, nil)
	require.NoError(t, err)
	pgTx := newSQLTx(rawTx.Tx, "postgres")

	seq, err := ls.AppendEventOutboxTx(ctx, pgTx, EventOutboxRecord{
		EventType:   "research.completed",
		ExecutionID: "exec-pg-branch",
		Payload:     `{"type":"research.completed"}`,
	})
	require.NoError(t, err)
	require.Greater(t, seq, int64(0))
	require.NoError(t, pgTx.Commit())

	rows, err := ls.ReadEventOutboxAfter(ctx, 0, 100)
	require.NoError(t, err)
	require.Len(t, rows, 1)
	require.Equal(t, "exec-pg-branch", rows[0].ExecutionID)
}

// TestAppendEventOutboxTx_ErrorsOnDeadTx proves both the postgres (Scan) and local (ExecContext)
// error-wrapping branches surface a real error: a tx already committed can no longer accept
// writes, so the underlying driver call fails and AppendEventOutboxTx wraps it rather than
// panicking or silently dropping it.
func TestAppendEventOutboxTx_ErrorsOnDeadTx(t *testing.T) {
	ls, ctx := newOutboxTxTestStorage(t)

	t.Run("local mode", func(t *testing.T) {
		tx, err := ls.requireSQLDB().BeginTx(ctx, nil)
		require.NoError(t, err)
		require.NoError(t, tx.Commit()) // tx is now dead

		_, err = ls.AppendEventOutboxTx(ctx, tx, EventOutboxRecord{
			EventType:   "research.completed",
			ExecutionID: "exec-dead-tx-local",
		})
		require.Error(t, err)
	})

	t.Run("postgres mode", func(t *testing.T) {
		rawTx, err := ls.requireSQLDB().BeginTx(ctx, nil)
		require.NoError(t, err)
		require.NoError(t, rawTx.Commit()) // tx is now dead
		pgTx := newSQLTx(rawTx.Tx, "postgres")

		_, err = ls.AppendEventOutboxTx(ctx, pgTx, EventOutboxRecord{
			EventType:   "research.completed",
			ExecutionID: "exec-dead-tx-pg",
		})
		require.Error(t, err)
	})
}
