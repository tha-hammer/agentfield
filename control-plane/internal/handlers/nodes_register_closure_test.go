package handlers_test

// Closure test for bsb slice 1 (restart-resilience: recoverable orphaning).
//
// It proves the full production chain end-to-end through mounted code, with no
// shortcuts:
//
//   seed running step (owned by node X) + agent_nodes.instance_id=old
//     -> POST /register  [RegisterNodeHandler]  (instance_id=new)
//        -> MarkAgentExecutionsOrphaned  (flips X's in-flight rows)
//     -> GET /runs/:run_id  [GetWorkflowRunDetailHandler]  (production read)
//        -> run is NOT failed; the orphaned step is recoverable `pending`.
//
// This file is `package handlers_test` (external) on purpose: the run-detail
// handler lives in package `ui`, which imports package `handlers`, so an
// internal `package handlers` test importing `ui` would create an import cycle.
// Every symbol used here is exported, so no unexported access is needed.
//
// RED-AT-SEAM: against the pre-slice code (the reap wrote terminal `failed`),
// the seeded step — owned by the re-registering node — is flipped to `failed`,
// so the OBSERVE assertions (`run.status != "failed"`, step `pending`) go RED.
// The connector under test is the `pending`-not-`failed` write in
// MarkAgentExecutionsOrphaned; disabling it (writing `failed`) makes this test
// red for the right reason.

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"
	"time"

	handlers "github.com/Agent-Field/agentfield/control-plane/internal/handlers"
	ui "github.com/Agent-Field/agentfield/control-plane/internal/handlers/ui"
	"github.com/Agent-Field/agentfield/control-plane/internal/storage"
	"github.com/Agent-Field/agentfield/control-plane/pkg/types"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

// newClosureStore stands up a real LocalStorage on an ephemeral sqlite/bolt DB.
// It fails-closed: if storage cannot initialize, the test errors (it never
// silently skips to green), except for the environmental FTS5-missing case that
// the storage-package harness also skips on.
func newClosureStore(t *testing.T) (*storage.LocalStorage, context.Context) {
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
			t.Skip("sqlite3 compiled without FTS5; skipping closure test")
		}
		require.NoError(t, err, "initialize local storage")
	}
	t.Cleanup(func() { _ = ls.Close(ctx) })
	return ls, ctx
}

func TestRegisterNodeHandler_RestartMakesInFlightStepsRecoverable_Closure(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ls, ctx := newClosureStore(t)
	now := time.Now().UTC()

	const (
		nodeID   = "github-buddy"
		runID    = "run-exec-1"
		execID   = "exec-1"
		baseURL  = "http://10.0.0.5:8080"
		oldInst  = "old-instance"
		newInst  = "new-instance"
		reasoner = "test.reasoner"
	)

	// SOURCE (seed only): the OLD agent instance, so re-registration detects a
	// genuine instance flip and triggers the reap.
	require.NoError(t, ls.RegisterAgent(ctx, &types.AgentNode{
		ID:              nodeID,
		BaseURL:         baseURL,
		InstanceID:      oldInst,
		Version:         "", // empty → GetAgent finds this row on re-register
		HealthStatus:    types.HealthStatusActive,
		LifecycleStatus: types.AgentStatusReady,
	}))

	// SOURCE (seed only): an in-flight step OWNED BY the re-registering node —
	// exactly the work a restart orphans. workflow_executions is the reap's
	// source of truth; the paired executions row is what the run-detail read
	// returns (it queries `executions` by run_id).
	require.NoError(t, ls.StoreWorkflowExecution(ctx, &types.WorkflowExecution{
		WorkflowID:          "wf-" + execID,
		ExecutionID:         execID,
		AgentFieldRequestID: "req-" + execID,
		AgentNodeID:         nodeID,
		ReasonerID:          reasoner,
		Status:              "running",
		StartedAt:           now,
		CreatedAt:           now,
		UpdatedAt:           now,
		WorkflowTags:        []string{},
		InputData:           json.RawMessage("{}"),
		OutputData:          json.RawMessage("{}"),
	}))
	require.NoError(t, ls.CreateExecutionRecord(ctx, &types.Execution{
		ExecutionID: execID,
		RunID:       runID,
		AgentNodeID: nodeID,
		ReasonerID:  reasoner,
		NodeID:      nodeID,
		Status:      "running",
		StartedAt:   now,
	}))

	// TRIGGER (start at the production entrypoint = highest_new_connector): the
	// re-register call with a NEW instance_id, through the real handler.
	registerRouter := gin.New()
	registerRouter.POST("/register", handlers.RegisterNodeHandler(ls, nil, nil, nil, nil, nil))
	body := `{"id":"` + nodeID + `","base_url":"` + baseURL + `","instance_id":"` + newInst + `",` +
		`"callback_discovery":{"mode":"manual","preferred":"` + baseURL + `"}}`
	regReq := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(body))
	regReq.Header.Set("Content-Type", "application/json")
	regRec := httptest.NewRecorder()
	registerRouter.ServeHTTP(regRec, regReq)
	require.Equal(t, http.StatusCreated, regRec.Code,
		"re-registration with a new instance must succeed; the instance flip is signal, not error")

	// OBSERVE (assert via the production read path only — never a raw table read
	// and never a direct MarkAgentExecutionsOrphaned call).
	detailRouter := gin.New()
	detailRouter.GET("/runs/:run_id", ui.NewWorkflowRunHandler(ls).GetWorkflowRunDetailHandler)
	detReq := httptest.NewRequest(http.MethodGet, "/runs/"+runID, nil)
	detRec := httptest.NewRecorder()
	detailRouter.ServeHTTP(detRec, detReq)
	require.Equal(t, http.StatusOK, detRec.Code, "run detail must be found for the seeded run")

	var detail struct {
		Run struct {
			Status      string `json:"status"`
			FailedSteps int    `json:"failed_steps"`
		} `json:"run"`
		Executions []struct {
			ExecutionID string `json:"execution_id"`
			Status      string `json:"status"`
		} `json:"executions"`
	}
	require.NoError(t, json.Unmarshal(detRec.Body.Bytes(), &detail))

	// The run is recoverable, not dead.
	require.NotEqual(t, "failed", detail.Run.Status,
		"a single restart must NOT mark the whole run failed — it is recoverable")
	require.Zero(t, detail.Run.FailedSteps,
		"the orphaned step must be excluded from failed_steps")

	// The orphaned in-flight step is retryable `pending`, not terminal `failed`.
	require.Len(t, detail.Executions, 1, "the seeded step must be present in the run detail")
	require.Equal(t, execID, detail.Executions[0].ExecutionID)
	require.Equal(t, "pending", detail.Executions[0].Status,
		"orphaned in-flight step must be recoverable pending, not terminal failed")
}
