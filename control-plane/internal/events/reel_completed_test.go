package events_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/Agent-Field/agentfield/control-plane/internal/events"
	"github.com/Agent-Field/agentfield/control-plane/pkg/types"

	"github.com/stretchr/testify/require"
)

// terminalReelExecution mirrors a succeeded reel-af execution whose result carries the dict
// returned by the reel_research_to_reel reasoner: source_execution_id / duration_s / beat_count
// plus a video body that must never leak into the event (C-Notification).
func terminalReelExecution() *types.Execution {
	return &types.Execution{
		ExecutionID: "exec_reel_1",
		RunID:       "run_reel_1",
		AgentNodeID: events.ReelAgentNodeID,
		Status:      string(types.ExecutionStatusSucceeded),
		ResultPayload: json.RawMessage(`{
			"source_execution_id": "exec_abc123",
			"duration_s": 12.5,
			"beat_count": 5,
			"video_path": "/out/reel.mp4 — must never appear in the event",
			"narration": "the spoken body — must never leak"
		}`),
	}
}

// ─────────────────────────── envelope shape ───────────────────────────

func TestBuildReelCompletedEvent_EnvelopeTypeAndSource(t *testing.T) {
	event, err := events.BuildReelCompletedEvent(terminalReelExecution(), "ce-fixed", time.Now())
	require.NoError(t, err)
	require.Equal(t, events.ReelCompletedEventType, event.Type)
	require.Equal(t, events.ReelCompletedSource, event.Source)
}

func TestBuildReelCompletedEvent_SubjectIsExecutionID(t *testing.T) {
	exec := terminalReelExecution()
	event, err := events.BuildReelCompletedEvent(exec, "ce-fixed", time.Now())
	require.NoError(t, err)
	require.Equal(t, exec.ExecutionID, event.Subject) // C-Correlation
}

func TestBuildReelCompletedEvent_InjectedIDAndTimeAreUsed(t *testing.T) {
	fixedTime := time.Date(2026, 7, 13, 18, 0, 0, 0, time.UTC)
	event, err := events.BuildReelCompletedEvent(terminalReelExecution(), "ce-fixed", fixedTime)
	require.NoError(t, err)
	require.Equal(t, "ce-fixed", event.ID)
	require.Equal(t, "2026-07-13T18:00:00Z", event.Time)
}

func TestNewReelCompletedEventID_PresentUniquePrefixed(t *testing.T) {
	id1 := events.NewReelCompletedEventID()
	id2 := events.NewReelCompletedEventID()
	require.NotEmpty(t, id1)
	require.True(t, len(id1) > len("ce-"))
	require.Contains(t, id1, "ce-")
	require.NotEqual(t, id1, id2)
}

// ─────────────────────────── DTO field sources ───────────────────────────

func TestBuildReelCompletedEvent_DataHasExactlyTheDTOKeys(t *testing.T) {
	event, err := events.BuildReelCompletedEvent(terminalReelExecution(), "ce-fixed", time.Now())
	require.NoError(t, err)

	raw, err := json.Marshal(event.Data)
	require.NoError(t, err)
	var asMap map[string]json.RawMessage
	require.NoError(t, json.Unmarshal(raw, &asMap))

	expectedKeys := []string{"run_id", "status", "reel_ref", "source_execution_id", "duration_s", "beat_count"}
	require.Len(t, asMap, len(expectedKeys))
	for _, key := range expectedKeys {
		_, ok := asMap[key]
		require.True(t, ok, "expected data key %q to be present", key)
	}
}

func TestBuildReelCompletedEvent_SourceExecutionIDFromResult(t *testing.T) {
	event, err := events.BuildReelCompletedEvent(terminalReelExecution(), "ce-fixed", time.Now())
	require.NoError(t, err)
	require.Equal(t, "exec_abc123", event.Data.SourceExecutionID)
}

func TestBuildReelCompletedEvent_DurationAndBeatCountFromResult(t *testing.T) {
	event, err := events.BuildReelCompletedEvent(terminalReelExecution(), "ce-fixed", time.Now())
	require.NoError(t, err)
	require.Equal(t, 12.5, event.Data.DurationS)
	require.Equal(t, 5, event.Data.BeatCount)
}

func TestBuildReelCompletedEvent_ReelRefIsExecutionIDKeyed(t *testing.T) {
	exec := terminalReelExecution()
	event, err := events.BuildReelCompletedEvent(exec, "ce-fixed", time.Now())
	require.NoError(t, err)
	require.Equal(t, "cp-execution://"+exec.ExecutionID+"/result", event.Data.ReelRef)
}

func TestBuildReelCompletedEvent_RunIDAndStatusFromRun(t *testing.T) {
	exec := terminalReelExecution()
	event, err := events.BuildReelCompletedEvent(exec, "ce-fixed", time.Now())
	require.NoError(t, err)
	require.Equal(t, exec.RunID, event.Data.RunID)
	require.Equal(t, "succeeded", event.Data.Status)
}

// ─────────────────────────── C-Notification: no body ───────────────────────────

func TestBuildReelCompletedEvent_VideoBodyAbsentFromData(t *testing.T) {
	event, err := events.BuildReelCompletedEvent(terminalReelExecution(), "ce-fixed", time.Now())
	require.NoError(t, err)

	raw, err := json.Marshal(event)
	require.NoError(t, err)
	require.NotContains(t, string(raw), "video_path")
	require.NotContains(t, string(raw), "narration")
	require.NotContains(t, string(raw), "must never")
}

func TestBuildReelCompletedEvent_NilExecution(t *testing.T) {
	_, err := events.BuildReelCompletedEvent(nil, "ce-fixed", time.Now())
	require.Error(t, err)
}

// ─────────────────────────── succeeded-only emit ───────────────────────────

func TestBuildReelCompletedEvent_EmittedForSucceededOnly(t *testing.T) {
	for _, badStatus := range []string{"failed", "cancelled", "running", ""} {
		t.Run(badStatus, func(t *testing.T) {
			exec := terminalReelExecution()
			exec.Status = badStatus
			_, err := events.BuildReelCompletedEvent(exec, "ce-fixed", time.Now())
			require.Error(t, err)
		})
	}
}

func TestBuildReelCompletedEvent_ValidWhenResultMissingEntirely(t *testing.T) {
	exec := terminalReelExecution()
	exec.ResultPayload = nil
	event, err := events.BuildReelCompletedEvent(exec, "ce-fixed", time.Now())
	require.NoError(t, err)
	require.Equal(t, exec.ExecutionID, event.Subject)
	require.Equal(t, "", event.Data.SourceExecutionID)
	require.Equal(t, 0, event.Data.BeatCount)
}

// ─────────────────────────── gate: IsReelCompletionCandidate ───────────────────────────

func TestIsReelCompletionCandidate(t *testing.T) {
	require.True(t, events.IsReelCompletionCandidate("reel-af"))
	require.False(t, events.IsReelCompletionCandidate("meta_deep_research"))
	require.False(t, events.IsReelCompletionCandidate(""))
}

// ─────────────────────────── BuildReelCompletedOutboxRecord (call-site helper) ───────────────────────────

func TestBuildReelCompletedOutboxRecord_EmitsForReelNodeSucceeded(t *testing.T) {
	rec, shouldAppend, err := events.BuildReelCompletedOutboxRecord(terminalReelExecution())
	require.NoError(t, err)
	require.True(t, shouldAppend)
	require.NotNil(t, rec)
	require.Equal(t, events.ReelCompletedEventType, rec.EventType)
	require.Equal(t, "exec_reel_1", rec.ExecutionID)
	require.Equal(t, "reel-af", rec.AgentNodeID)
	require.Contains(t, rec.Payload, `"type":"com.silmari.reel.completed.v1"`)
	require.NotContains(t, rec.Payload, "video_path")
}

func TestBuildReelCompletedOutboxRecord_SkipsNonReelNode(t *testing.T) {
	exec := terminalReelExecution()
	exec.AgentNodeID = "meta_deep_research"
	rec, shouldAppend, err := events.BuildReelCompletedOutboxRecord(exec)
	require.NoError(t, err)
	require.False(t, shouldAppend)
	require.Nil(t, rec)
}

func TestBuildReelCompletedOutboxRecord_SkipsNonSucceeded(t *testing.T) {
	exec := terminalReelExecution()
	exec.Status = string(types.ExecutionStatusFailed)
	rec, shouldAppend, err := events.BuildReelCompletedOutboxRecord(exec)
	require.NoError(t, err)
	require.False(t, shouldAppend)
	require.Nil(t, rec)
}

func TestBuildReelCompletedOutboxRecord_NilExecution(t *testing.T) {
	rec, shouldAppend, err := events.BuildReelCompletedOutboxRecord(nil)
	require.NoError(t, err)
	require.False(t, shouldAppend)
	require.Nil(t, rec)
}

// ─────────────────────────── dispatcher: builder-select by AgentNodeID ───────────────────────────

func TestBuildCompletedOutboxRecord_RoutesReelNode(t *testing.T) {
	rec, shouldAppend, err := events.BuildCompletedOutboxRecord(terminalReelExecution())
	require.NoError(t, err)
	require.True(t, shouldAppend)
	require.Equal(t, events.ReelCompletedEventType, rec.EventType)
}

func TestBuildCompletedOutboxRecord_RoutesResearchNode(t *testing.T) {
	rec, shouldAppend, err := events.BuildCompletedOutboxRecord(terminalResearchExecution())
	require.NoError(t, err)
	require.True(t, shouldAppend)
	require.Equal(t, events.ResearchCompletedEventType, rec.EventType)
}

func TestBuildCompletedOutboxRecord_SkipsUnrelatedAgent(t *testing.T) {
	exec := terminalReelExecution()
	exec.AgentNodeID = "some-other-agent"
	rec, shouldAppend, err := events.BuildCompletedOutboxRecord(exec)
	require.NoError(t, err)
	require.False(t, shouldAppend)
	require.Nil(t, rec)
}

func TestBuildCompletedOutboxRecord_SkipsNonSucceededReel(t *testing.T) {
	exec := terminalReelExecution()
	exec.Status = string(types.ExecutionStatusFailed)
	_, shouldAppend, err := events.BuildCompletedOutboxRecord(exec)
	require.NoError(t, err)
	require.False(t, shouldAppend)
}
