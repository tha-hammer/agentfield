package events_test

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/Agent-Field/agentfield/control-plane/internal/events"
	"github.com/stretchr/testify/require"
)

// TestCompletionEnvelopeRoundTripsThroughOutbox asserts that a CloudEvents-shaped
// Data payload survives Publish → ReadEventOutboxAfter with ExecutionID preserved
// as the correlation key and the whole envelope in Payload.
func TestCompletionEnvelopeRoundTripsThroughOutbox(t *testing.T) {
	ls, ctx := newOutboxStore(t)
	bus := events.NewDurableExecutionBus(ls)

	ceData := map[string]any{
		"specversion": "1.0",
		"id":          "ce-test-id",
		"source":      "deep-research",
		"type":        "com.silmari.research.completed",
		"subject":     "exec_9",
		"data": map[string]any{
			"result_ref":  "pkg://x",
			"prompt":      "why FDO",
			"document_id": "doc_42",
		},
	}

	ev := events.ExecutionEvent{
		Type:        "com.silmari.research.completed",
		ExecutionID: "exec_9",
		AgentNodeID: "deep-research",
		Timestamp:   time.Now().UTC(),
		Data:        ceData,
	}

	require.NoError(t, bus.Publish(ctx, ev))

	recs, err := ls.ReadEventOutboxAfter(ctx, 0, 10)
	require.NoError(t, err)
	require.Len(t, recs, 1)
	require.Equal(t, "exec_9", recs[0].ExecutionID)

	// The Payload is the JSON-marshaled ExecutionEvent; Data inside carries
	// the CloudEvents envelope. Unmarshal and verify the envelope survived.
	var stored events.ExecutionEvent
	require.NoError(t, json.Unmarshal([]byte(recs[0].Payload), &stored))
	require.Equal(t, "exec_9", stored.ExecutionID)

	dataBytes, err := json.Marshal(stored.Data)
	require.NoError(t, err)

	var roundTripped map[string]any
	require.NoError(t, json.Unmarshal(dataBytes, &roundTripped))
	require.Equal(t, "exec_9", roundTripped["subject"])
	require.Equal(t, "com.silmari.research.completed", roundTripped["type"])
	require.Equal(t, "1.0", roundTripped["specversion"])

	innerData, ok := roundTripped["data"].(map[string]any)
	require.True(t, ok, "data should be a map")
	require.Equal(t, "pkg://x", innerData["result_ref"])
}

// readGolden reads and parses the shared golden JSON fixture.
func readGolden(t *testing.T, relPath string) map[string]any {
	t.Helper()
	absPath := filepath.Join(filepath.Dir("."), relPath)
	raw, err := os.ReadFile(absPath)
	require.NoError(t, err, "golden fixture not found at %s", absPath)
	var m map[string]any
	require.NoError(t, json.Unmarshal(raw, &m))
	return m
}

// TestGoldenEnvelopeSurvivesOutboxRoundTrip asserts the SAME golden envelope
// the Python test asserts survives the Go outbox path — cross-language shape parity.
func TestGoldenEnvelopeSurvivesOutboxRoundTrip(t *testing.T) {
	golden := readGolden(t, "../../../testdata/completion_envelope.golden.json")
	ls, ctx := newOutboxStore(t)
	bus := events.NewDurableExecutionBus(ls)

	ev := events.ExecutionEvent{
		Type:        events.ExecutionEventType(golden["type"].(string)),
		ExecutionID: golden["subject"].(string),
		AgentNodeID: golden["source"].(string),
		Timestamp:   time.Now().UTC(),
		Data:        golden,
	}

	require.NoError(t, bus.Publish(context.Background(), ev))

	recs, err := ls.ReadEventOutboxAfter(ctx, 0, 10)
	require.NoError(t, err)
	require.Len(t, recs, 1)

	// Unmarshal the stored Payload → ExecutionEvent, then extract Data
	var stored events.ExecutionEvent
	require.NoError(t, json.Unmarshal([]byte(recs[0].Payload), &stored))

	dataBytes, err := json.Marshal(stored.Data)
	require.NoError(t, err)

	var roundTripped map[string]any
	require.NoError(t, json.Unmarshal(dataBytes, &roundTripped))

	// The golden's "id" is "PLACEHOLDER" — we only care about shape parity,
	// not the uuid value. Compare all fields except id.
	goldenCopy := make(map[string]any)
	for k, v := range golden {
		goldenCopy[k] = v
	}
	delete(goldenCopy, "id")
	delete(roundTripped, "id")

	require.Equal(t, goldenCopy["specversion"], roundTripped["specversion"])
	require.Equal(t, goldenCopy["source"], roundTripped["source"])
	require.Equal(t, goldenCopy["type"], roundTripped["type"])
	require.Equal(t, goldenCopy["subject"], roundTripped["subject"])
	require.Equal(t, goldenCopy["dataschema"], roundTripped["dataschema"])

	// Deep-compare the data payload
	goldenData, _ := json.Marshal(goldenCopy["data"])
	rtData, _ := json.Marshal(roundTripped["data"])
	require.JSONEq(t, string(goldenData), string(rtData), "CloudEvents data payload drifted through outbox")
}
