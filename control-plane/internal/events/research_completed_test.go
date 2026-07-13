package events_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/Agent-Field/agentfield/control-plane/internal/events"
	"github.com/Agent-Field/agentfield/control-plane/pkg/types"

	"github.com/stretchr/testify/require"
)

// terminalResearchExecution mirrors the Python producer test fixture's
// `terminal_execution` (tests/producer/conftest.py): a succeeded deep-research
// execution whose result carries metadata.query/title and a research_package body
// that must never leak into the event.
func terminalResearchExecution() *types.Execution {
	return &types.Execution{
		ExecutionID: "exec_abc123",
		RunID:       "3f6d1a90-0000-4000-8000-000000000abc",
		AgentNodeID: events.ResearchAgentNodeID,
		Status:      string(types.ExecutionStatusSucceeded),
		ResultPayload: json.RawMessage(`{
			"metadata": {
				"query": "How do short-form reels convert viewers into subscribers?",
				"title": "How short-form reels convert viewers"
			},
			"research_package": {"document": "the full body — must never appear in the event"}
		}`),
	}
}

// ─────────────────────────── envelope shape ───────────────────────────

func TestBuildResearchCompletedEvent_EnvelopeTypeAndSource(t *testing.T) {
	event, err := events.BuildResearchCompletedEvent(terminalResearchExecution(), "ce-fixed", time.Now())
	require.NoError(t, err)
	require.Equal(t, events.ResearchCompletedEventType, event.Type)
	require.Equal(t, events.ResearchCompletedSource, event.Source)
}

func TestBuildResearchCompletedEvent_SubjectIsExecutionID(t *testing.T) {
	exec := terminalResearchExecution()
	event, err := events.BuildResearchCompletedEvent(exec, "ce-fixed", time.Now())
	require.NoError(t, err)
	require.Equal(t, exec.ExecutionID, event.Subject) // C-Correlation
}

func TestBuildResearchCompletedEvent_InjectedIDAndTimeAreUsed(t *testing.T) {
	fixedTime := time.Date(2026, 7, 12, 18, 0, 0, 0, time.UTC)
	event, err := events.BuildResearchCompletedEvent(terminalResearchExecution(), "ce-fixed", fixedTime)
	require.NoError(t, err)
	require.Equal(t, "ce-fixed", event.ID)
	require.Equal(t, "2026-07-12T18:00:00Z", event.Time)
}

func TestNewResearchCompletedEventID_PresentUniquePrefixed(t *testing.T) {
	id1 := events.NewResearchCompletedEventID()
	id2 := events.NewResearchCompletedEventID()
	require.NotEmpty(t, id1)
	require.True(t, len(id1) > len("ce-"))
	require.Contains(t, id1, "ce-")
	require.NotEqual(t, id1, id2)
}

// ─────────────────────────── DTO field sources ───────────────────────────

func TestBuildResearchCompletedEvent_DataHasExactlyTheDTOKeys(t *testing.T) {
	event, err := events.BuildResearchCompletedEvent(terminalResearchExecution(), "ce-fixed", time.Now())
	require.NoError(t, err)

	raw, err := json.Marshal(event.Data)
	require.NoError(t, err)
	var asMap map[string]json.RawMessage
	require.NoError(t, json.Unmarshal(raw, &asMap))

	expectedKeys := []string{"run_id", "status", "title", "result_ref", "research_prompt", "research_document_id"}
	require.Len(t, asMap, len(expectedKeys))
	for _, key := range expectedKeys {
		_, ok := asMap[key]
		require.True(t, ok, "expected data key %q to be present", key)
	}
}

func TestBuildResearchCompletedEvent_ResearchPromptFromMetadataQuery(t *testing.T) {
	exec := terminalResearchExecution()
	event, err := events.BuildResearchCompletedEvent(exec, "ce-fixed", time.Now())
	require.NoError(t, err)
	require.NotNil(t, event.Data.ResearchPrompt)
	require.Equal(t, "How do short-form reels convert viewers into subscribers?", *event.Data.ResearchPrompt)
}

func TestBuildResearchCompletedEvent_TitleFromMetadataTitle(t *testing.T) {
	exec := terminalResearchExecution()
	event, err := events.BuildResearchCompletedEvent(exec, "ce-fixed", time.Now())
	require.NoError(t, err)
	require.NotNil(t, event.Data.Title)
	require.Equal(t, "How short-form reels convert viewers", *event.Data.Title)
}

func TestBuildResearchCompletedEvent_ResearchDocumentIDIsExecutionID(t *testing.T) {
	exec := terminalResearchExecution()
	event, err := events.BuildResearchCompletedEvent(exec, "ce-fixed", time.Now())
	require.NoError(t, err)
	require.Equal(t, exec.ExecutionID, event.Data.ResearchDocumentID)
}

func TestBuildResearchCompletedEvent_ResultRefIsExecutionIDKeyed(t *testing.T) {
	exec := terminalResearchExecution()
	event, err := events.BuildResearchCompletedEvent(exec, "ce-fixed", time.Now())
	require.NoError(t, err)
	require.Equal(t, "cp-execution://"+exec.ExecutionID+"/result", event.Data.ResultRef)
}

func TestBuildResearchCompletedEvent_RunIDAndStatusFromRun(t *testing.T) {
	exec := terminalResearchExecution()
	event, err := events.BuildResearchCompletedEvent(exec, "ce-fixed", time.Now())
	require.NoError(t, err)
	require.Equal(t, exec.RunID, event.Data.RunID)
	require.Equal(t, "succeeded", event.Data.Status)
}

// ─────────────────────────── C-Notification: no body ───────────────────────────

func TestBuildResearchCompletedEvent_ResearchPackageBodyAbsentFromData(t *testing.T) {
	event, err := events.BuildResearchCompletedEvent(terminalResearchExecution(), "ce-fixed", time.Now())
	require.NoError(t, err)

	raw, err := json.Marshal(event)
	require.NoError(t, err)
	require.NotContains(t, string(raw), "research_package")
	require.NotContains(t, string(raw), "the full body")
}

func TestBuildResearchCompletedEvent_NilExecution(t *testing.T) {
	_, err := events.BuildResearchCompletedEvent(nil, "ce-fixed", time.Now())
	require.Error(t, err)
}

// ─────────────────────────── succeeded-only emit ───────────────────────────

func TestBuildResearchCompletedEvent_EmittedForSucceededOnly(t *testing.T) {
	for _, badStatus := range []string{"failed", "cancelled", "running", ""} {
		t.Run(badStatus, func(t *testing.T) {
			exec := terminalResearchExecution()
			exec.Status = badStatus
			_, err := events.BuildResearchCompletedEvent(exec, "ce-fixed", time.Now())
			require.Error(t, err)
		})
	}
}

// ─────────────────────────── nullable snapshot fields ───────────────────────────

func TestBuildResearchCompletedEvent_TitleNoneWhenMetadataTitleAbsent(t *testing.T) {
	exec := terminalResearchExecution()
	exec.ResultPayload = json.RawMessage(`{"metadata": {"query": "Q"}}`)
	event, err := events.BuildResearchCompletedEvent(exec, "ce-fixed", time.Now())
	require.NoError(t, err)
	require.Nil(t, event.Data.Title)
}

func TestBuildResearchCompletedEvent_ResearchPromptNoneWhenMetadataQueryAbsent(t *testing.T) {
	exec := terminalResearchExecution()
	exec.ResultPayload = json.RawMessage(`{"metadata": {"title": "T"}}`)
	event, err := events.BuildResearchCompletedEvent(exec, "ce-fixed", time.Now())
	require.NoError(t, err)
	require.Nil(t, event.Data.ResearchPrompt)
}

func TestBuildResearchCompletedEvent_ValidWhenResultMissingEntirely(t *testing.T) {
	exec := terminalResearchExecution()
	exec.ResultPayload = nil
	event, err := events.BuildResearchCompletedEvent(exec, "ce-fixed", time.Now())
	require.NoError(t, err)
	require.Equal(t, exec.ExecutionID, event.Subject)
	require.Nil(t, event.Data.Title)
	require.Nil(t, event.Data.ResearchPrompt)
}

// ─────────────────────────── gate: IsResearchCompletionCandidate ───────────────────────────

func TestIsResearchCompletionCandidate(t *testing.T) {
	require.True(t, events.IsResearchCompletionCandidate("meta_deep_research"))
	require.False(t, events.IsResearchCompletionCandidate("some-other-agent"))
	require.False(t, events.IsResearchCompletionCandidate(""))
}

// ─────────────────────────── BuildResearchCompletedOutboxRecord (the call-site helper) ───────────────────────────

func TestBuildResearchCompletedOutboxRecord_EmitsForResearchNodeSucceeded(t *testing.T) {
	rec, shouldAppend, err := events.BuildResearchCompletedOutboxRecord(terminalResearchExecution())
	require.NoError(t, err)
	require.True(t, shouldAppend)
	require.NotNil(t, rec)
	require.Equal(t, events.ResearchCompletedEventType, rec.EventType)
	require.Equal(t, "exec_abc123", rec.ExecutionID)
	require.Equal(t, "meta_deep_research", rec.AgentNodeID)
	require.Contains(t, rec.Payload, `"type":"com.silmari.research.completed.v1"`)
	require.NotContains(t, rec.Payload, "research_package")
}

func TestBuildResearchCompletedOutboxRecord_SkipsNonResearchNode(t *testing.T) {
	exec := terminalResearchExecution()
	exec.AgentNodeID = "some-other-agent"
	rec, shouldAppend, err := events.BuildResearchCompletedOutboxRecord(exec)
	require.NoError(t, err)
	require.False(t, shouldAppend)
	require.Nil(t, rec)
}

func TestBuildResearchCompletedOutboxRecord_SkipsNonSucceeded(t *testing.T) {
	exec := terminalResearchExecution()
	exec.Status = string(types.ExecutionStatusFailed)
	rec, shouldAppend, err := events.BuildResearchCompletedOutboxRecord(exec)
	require.NoError(t, err)
	require.False(t, shouldAppend)
	require.Nil(t, rec)
}

func TestBuildResearchCompletedOutboxRecord_NilExecution(t *testing.T) {
	rec, shouldAppend, err := events.BuildResearchCompletedOutboxRecord(nil)
	require.NoError(t, err)
	require.False(t, shouldAppend)
	require.Nil(t, rec)
}
