package storage

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/Agent-Field/agentfield/control-plane/internal/logger"
	"github.com/Agent-Field/agentfield/control-plane/pkg/types"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// EventOutboxRecord is the storage-facing shape of a durable execution event.
// It aliases the neutral pkg/types definition so the events layer can append
// records without importing storage. Seq and CreatedAt are set on append.
type EventOutboxRecord = types.EventOutboxRecord

// PruneResult reports what a rotation pass removed. OverflowUnread counts pruned
// rows above the consumer low-water mark, so unread loss is never silent.
type PruneResult struct {
	Deleted        int
	OverflowUnread int
}

const outboxReadDefaultLimit = 100

// eventOutboxModelToRecord is the single row->record mapper shared by every read
// path (mirrors modelToInboundEvent in triggers.go).
func eventOutboxModelToRecord(m EventOutboxModel) EventOutboxRecord {
	return EventOutboxRecord{
		Seq:         m.Seq,
		EventType:   m.EventType,
		ExecutionID: m.ExecutionID,
		WorkflowID:  m.WorkflowID,
		AgentNodeID: m.AgentNodeID,
		Payload:     m.Payload,
		CreatedAt:   m.CreatedAt,
	}
}

// AppendEventOutbox persists an event and returns its assigned monotonic seq.
// The seq is a dedicated autoincrement PK, so cursor reads never depend on
// timestamp ordering.
func (ls *LocalStorage) AppendEventOutbox(ctx context.Context, rec EventOutboxRecord) (int64, error) {
	gormDB, err := ls.gormWithContext(ctx)
	if err != nil {
		return 0, err
	}
	model := EventOutboxModel{
		EventType:   rec.EventType,
		ExecutionID: rec.ExecutionID,
		WorkflowID:  rec.WorkflowID,
		AgentNodeID: rec.AgentNodeID,
		Payload:     rec.Payload,
		CreatedAt:   rec.CreatedAt,
	}
	if model.Payload == "" {
		model.Payload = "{}"
	}
	if model.CreatedAt.IsZero() {
		model.CreatedAt = time.Now().UTC()
	}
	if err := gormDB.Create(&model).Error; err != nil {
		return 0, err
	}
	return model.Seq, nil
}

// ReadEventOutboxAfter returns events with seq strictly greater than afterSeq,
// ascending, at most limit rows. A negative afterSeq is treated as 0; an
// afterSeq at or beyond the max seq yields an empty slice, not an error.
func (ls *LocalStorage) ReadEventOutboxAfter(ctx context.Context, afterSeq int64, limit int) ([]EventOutboxRecord, error) {
	if afterSeq < 0 {
		afterSeq = 0
	}
	if limit <= 0 || limit > 1000 {
		limit = outboxReadDefaultLimit
	}
	gormDB, err := ls.gormWithContext(ctx)
	if err != nil {
		return nil, err
	}
	var rows []EventOutboxModel
	if err := gormDB.Where("seq > ?", afterSeq).
		Order("seq ASC").
		Limit(limit).
		Find(&rows).Error; err != nil {
		return nil, err
	}
	out := make([]EventOutboxRecord, 0, len(rows))
	for _, r := range rows {
		out = append(out, eventOutboxModelToRecord(r))
	}
	return out, nil
}

// CountEventOutbox returns the number of rows currently in the outbox.
func (ls *LocalStorage) CountEventOutbox(ctx context.Context) (int64, error) {
	gormDB, err := ls.gormWithContext(ctx)
	if err != nil {
		return 0, err
	}
	var n int64
	if err := gormDB.Model(&EventOutboxModel{}).Count(&n).Error; err != nil {
		return 0, err
	}
	return n, nil
}

// HasTerminalOutboxEvent reports whether the durable outbox holds a terminal
// (completed/failed) event for executionID. The sync awaiter uses this on
// timeout to recover a terminal event whose live delivery was dropped. Uses the
// execution_id index.
func (ls *LocalStorage) HasTerminalOutboxEvent(ctx context.Context, executionID string) (bool, error) {
	gormDB, err := ls.gormWithContext(ctx)
	if err != nil {
		return false, err
	}
	var count int64
	if err := gormDB.Model(&EventOutboxModel{}).
		Where("execution_id = ? AND event_type IN ?", executionID, []string{"execution_completed", "execution_failed"}).
		Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// GetOutboxCursor returns how far consumerID has read the outbox, or 0 if the
// consumer has no cursor row yet.
func (ls *LocalStorage) GetOutboxCursor(ctx context.Context, consumerID string) (int64, error) {
	gormDB, err := ls.gormWithContext(ctx)
	if err != nil {
		return 0, err
	}
	var row EventOutboxCursorModel
	err = gormDB.Where("consumer_id = ?", consumerID).First(&row).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}
	return row.LastSeq, nil
}

// AdvanceOutboxCursor upserts consumerID's cursor. It is monotonic: last_seq
// never regresses, so registering a consumer at seq 0 (on subscribe) cannot
// reset a cursor that has already advanced.
func (ls *LocalStorage) AdvanceOutboxCursor(ctx context.Context, consumerID string, seq int64) error {
	gormDB, err := ls.gormWithContext(ctx)
	if err != nil {
		return err
	}
	current, err := ls.GetOutboxCursor(ctx, consumerID)
	if err != nil {
		return err
	}
	if current > seq {
		seq = current
	}
	return gormDB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "consumer_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"last_seq"}),
	}).Create(&EventOutboxCursorModel{ConsumerID: consumerID, LastSeq: seq}).Error
}

// outboxLowWater returns the minimum cursor position across all registered
// consumers (0 when none are registered). Rows above the low-water mark are
// unread by at least one registered consumer.
func (ls *LocalStorage) outboxLowWater(gormDB *gorm.DB) (int64, error) {
	var lw sql.NullInt64
	if err := gormDB.Model(&EventOutboxCursorModel{}).Select("MIN(last_seq)").Row().Scan(&lw); err != nil {
		return 0, err
	}
	if !lw.Valid {
		return 0, nil
	}
	return lw.Int64, nil
}

// PruneEventOutbox enforces the retention caps and returns what it removed.
// olderThan (non-zero) deletes rows created before the cutoff; maxRows (>0)
// keeps only the newest maxRows by seq. OverflowUnread counts pruned rows above
// the consumer low-water mark — unread loss is bounded, counted, and logged,
// never silent. Both branches are idempotent.
func (ls *LocalStorage) PruneEventOutbox(ctx context.Context, olderThan time.Time, maxRows int) (PruneResult, error) {
	var res PruneResult
	gormDB, err := ls.gormWithContext(ctx)
	if err != nil {
		return res, err
	}

	lowWater, err := ls.outboxLowWater(gormDB)
	if err != nil {
		return res, err
	}

	// Age branch: delete rows older than the cutoff.
	if !olderThan.IsZero() {
		var overflow int64
		if err := gormDB.Model(&EventOutboxModel{}).
			Where("created_at < ? AND seq > ?", olderThan, lowWater).
			Count(&overflow).Error; err != nil {
			return res, err
		}
		r := gormDB.Where("created_at < ?", olderThan).Delete(&EventOutboxModel{})
		if r.Error != nil {
			return res, r.Error
		}
		res.Deleted += int(r.RowsAffected)
		res.OverflowUnread += int(overflow)
	}

	// Count branch: keep only the newest maxRows by seq.
	if maxRows > 0 {
		var total int64
		if err := gormDB.Model(&EventOutboxModel{}).Count(&total).Error; err != nil {
			return res, err
		}
		if total > int64(maxRows) {
			// keepMin = smallest seq among the newest maxRows rows.
			var keepMin int64
			if err := gormDB.Model(&EventOutboxModel{}).
				Select("seq").Order("seq DESC").
				Offset(maxRows - 1).Limit(1).
				Row().Scan(&keepMin); err != nil {
				return res, err
			}
			var overflow int64
			if err := gormDB.Model(&EventOutboxModel{}).
				Where("seq < ? AND seq > ?", keepMin, lowWater).
				Count(&overflow).Error; err != nil {
				return res, err
			}
			r := gormDB.Where("seq < ?", keepMin).Delete(&EventOutboxModel{})
			if r.Error != nil {
				return res, r.Error
			}
			res.Deleted += int(r.RowsAffected)
			res.OverflowUnread += int(overflow)
		}
	}

	if res.OverflowUnread > 0 {
		logger.Logger.Warn().
			Int("overflow_unread", res.OverflowUnread).
			Int("deleted", res.Deleted).
			Int64("low_water", lowWater).
			Msg("event outbox rotation pruned unread events to honor the cap")
	}

	return res, nil
}
