-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS event_outbox (
    seq          BIGSERIAL PRIMARY KEY,
    event_type   TEXT NOT NULL DEFAULT '',
    execution_id TEXT NOT NULL DEFAULT '',
    workflow_id  TEXT NOT NULL DEFAULT '',
    agent_node_id TEXT NOT NULL DEFAULT '',
    payload      TEXT NOT NULL DEFAULT '{}',
    created_at   TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_event_outbox_execution_id ON event_outbox(execution_id);
CREATE INDEX IF NOT EXISTS idx_event_outbox_created_at ON event_outbox(created_at);

CREATE TABLE IF NOT EXISTS event_outbox_cursor (
    consumer_id TEXT PRIMARY KEY,
    last_seq    BIGINT NOT NULL DEFAULT 0
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS event_outbox_cursor;
DROP INDEX IF EXISTS idx_event_outbox_created_at;
DROP INDEX IF EXISTS idx_event_outbox_execution_id;
DROP TABLE IF EXISTS event_outbox;
-- +goose StatementEnd
