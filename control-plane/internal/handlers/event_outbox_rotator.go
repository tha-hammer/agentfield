package handlers

import (
	"context"
	"sync"
	"time"

	"github.com/Agent-Field/agentfield/control-plane/internal/config"
	"github.com/Agent-Field/agentfield/control-plane/internal/logger"
	"github.com/Agent-Field/agentfield/control-plane/internal/storage"
)

// OutboxPruner is the storage surface the rotator needs. It is declared here
// (the consuming package) per Go convention and is satisfied by
// *storage.LocalStorage, so the rotator never depends on the full
// StorageProvider interface.
type OutboxPruner interface {
	PruneEventOutbox(ctx context.Context, olderThan time.Time, maxRows int) (storage.PruneResult, error)
	CountEventOutbox(ctx context.Context) (int64, error)
}

// EventOutboxRotator periodically prunes the durable event outbox to keep it
// capped by age and row count. It mirrors ExecutionCleanupService's lifecycle
// (stopChan + WaitGroup, Enabled-guarded Start, clean Stop).
type EventOutboxRotator struct {
	store     OutboxPruner
	config    config.EventOutboxConfig
	stopChan  chan struct{}
	wg        sync.WaitGroup
	isRunning bool
	mu        sync.Mutex
}

// NewEventOutboxRotator constructs a rotator with normalized config (a safe
// non-zero prune interval).
func NewEventOutboxRotator(store OutboxPruner, cfg config.EventOutboxConfig) *EventOutboxRotator {
	return &EventOutboxRotator{
		store:    store,
		config:   config.EffectiveEventOutbox(cfg),
		stopChan: make(chan struct{}),
	}
}

// Start begins the background rotation loop. It is a no-op when disabled.
func (r *EventOutboxRotator) Start(ctx context.Context) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.isRunning {
		return nil
	}
	if !r.config.Enabled {
		logger.Logger.Warn().Msg("Event outbox rotation is disabled")
		return nil
	}

	logger.Logger.Debug().
		Dur("retention_max_age", r.config.RetentionMaxAge).
		Int("retention_max_rows", r.config.RetentionMaxRows).
		Dur("prune_interval", r.config.PruneInterval).
		Msg("Starting event outbox rotator")

	r.isRunning = true
	r.wg.Add(1)
	go r.loop(ctx)
	return nil
}

// Stop halts the rotation loop and waits for the goroutine to drain.
func (r *EventOutboxRotator) Stop() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if !r.isRunning {
		return nil
	}
	close(r.stopChan)
	r.wg.Wait()
	r.isRunning = false
	logger.Logger.Debug().Msg("Event outbox rotator stopped")
	return nil
}

func (r *EventOutboxRotator) loop(ctx context.Context) {
	defer r.wg.Done()

	ticker := time.NewTicker(r.config.PruneInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-r.stopChan:
			return
		case <-ticker.C:
			if err := r.pruneOnce(ctx); err != nil {
				logger.Logger.Error().Err(err).Msg("event outbox rotation failed")
			}
		}
	}
}

// pruneOnce enforces the age+count caps exactly once. It is the only place caps
// are applied; the ticker loop is a thin wrapper around this directly-callable
// tested seam.
func (r *EventOutboxRotator) pruneOnce(ctx context.Context) error {
	var olderThan time.Time
	if r.config.RetentionMaxAge > 0 {
		olderThan = time.Now().Add(-r.config.RetentionMaxAge)
	}
	res, err := r.store.PruneEventOutbox(ctx, olderThan, r.config.RetentionMaxRows)
	if err != nil {
		return err
	}
	if res.Deleted > 0 {
		logger.Logger.Debug().
			Int("deleted", res.Deleted).
			Int("overflow_unread", res.OverflowUnread).
			Msg("event outbox rotated")
	}
	return nil
}
