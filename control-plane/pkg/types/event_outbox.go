package types

import "time"

// EventOutboxRecord is the durable, storage-neutral shape of an execution event
// in the outbox. It lives in pkg/types so both the storage layer (which writes
// it) and the events layer (whose durable bus appends it) can share it without
// an import cycle. Seq and CreatedAt are assigned by the store on append.
type EventOutboxRecord struct {
	Seq         int64
	EventType   string
	ExecutionID string
	WorkflowID  string
	AgentNodeID string
	Payload     string
	CreatedAt   time.Time
}
