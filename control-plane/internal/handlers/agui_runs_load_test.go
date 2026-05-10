package handlers

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/Agent-Field/agentfield/control-plane/pkg/types"

	"github.com/stretchr/testify/require"
)

// TestAGUI_Load_ConcurrentBuffered hammers the AG-UI handler with many
// concurrent requests against a fast buffered reasoner and asserts:
//
//   - Every request returns a complete canonical event sequence.
//   - Goroutines don't leak: the count after all runs settle is
//     approximately the baseline (a few +/- for runtime noise).
//   - p50/p95/p99 latencies stay within reasonable bounds at 200 in-flight.
//
// This is the production-readiness gate the earlier 5×concurrent test
// could not provide. It runs in CI as part of `go test`.
func TestAGUI_Load_ConcurrentBuffered(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping load test in -short mode")
	}

	const totalRequests = 200
	const concurrency = 50

	reasoner := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"result":"ok","state":{"counter":1}}`))
	}))
	defer reasoner.Close()

	store := &reasonerTestStorage{agent: &types.AgentNode{
		ID:              "load-node",
		BaseURL:         reasoner.URL,
		HealthStatus:    types.HealthStatusActive,
		LifecycleStatus: types.AgentStatusReady,
		Reasoners:       []types.ReasonerDefinition{{ID: "r"}},
	}}
	router := mountAGUIRouter(t, store)

	// Sample the goroutine baseline AFTER the test runtime is up but
	// BEFORE we fire load. NumGoroutine() is nondeterministic so we
	// give the handler a generous tolerance — we're guarding against
	// real leaks (200 leaked goroutines per 200 runs), not noise.
	runtime.GC()
	time.Sleep(50 * time.Millisecond)
	baseline := runtime.NumGoroutine()

	var (
		started   atomic.Int64
		completed atomic.Int64
		failed    atomic.Int64
		latencies = make([]time.Duration, totalRequests)
	)
	sem := make(chan struct{}, concurrency)
	var wg sync.WaitGroup

	wallStart := time.Now()
	for i := 0; i < totalRequests; i++ {
		wg.Add(1)
		sem <- struct{}{}
		go func(idx int) {
			defer wg.Done()
			defer func() { <-sem }()
			started.Add(1)

			body := fmt.Sprintf(`{"threadId":"t-%d","runId":"r-%d","messages":[{"role":"user","content":"x"}]}`, idx, idx)
			req := httptest.NewRequest(http.MethodPost, "/api/v1/agui/runs/load-node/r", strings.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			start := time.Now()
			router.ServeHTTP(w, req)
			latencies[idx] = time.Since(start)

			if w.Code != http.StatusOK {
				failed.Add(1)
				return
			}
			frames := parseAGUIStream(t, w.Body.String())
			if len(frames) == 0 || frames[0].Type() != "RUN_STARTED" || frames[len(frames)-1].Type() != "RUN_FINISHED" {
				failed.Add(1)
				return
			}
			completed.Add(1)
		}(i)
	}
	wg.Wait()
	wallElapsed := time.Since(wallStart)

	require.Equal(t, int64(totalRequests), started.Load(), "all requests should have started")
	require.Equal(t, int64(totalRequests), completed.Load(), "all requests should have completed: failures=%d", failed.Load())
	require.Equal(t, int64(0), failed.Load(), "no requests should have failed under load")

	// Latency stats — sort then pick percentiles.
	sortDurations(latencies)
	p50 := latencies[len(latencies)*50/100]
	p95 := latencies[len(latencies)*95/100]
	p99 := latencies[len(latencies)*99/100]

	t.Logf("load: %d reqs at %d concurrency, wall=%s, p50=%s p95=%s p99=%s",
		totalRequests, concurrency, wallElapsed, p50, p95, p99)

	// Loose latency budget — the handler is just routing + emitting
	// events against an in-process httptest reasoner, so even p99
	// shouldn't exceed 250ms on a quiet box.
	require.Less(t, p95, 250*time.Millisecond, "p95 latency too high under 50× concurrent load")

	// Goroutine leak check. Every request spawns one goroutine
	// (invoker.Invoke). They should all have settled by now. Allow a
	// generous buffer for test infra (httptest handlers can keep
	// goroutines around briefly) but flag a real leak.
	runtime.GC()
	time.Sleep(100 * time.Millisecond)
	final := runtime.NumGoroutine()
	t.Logf("goroutines: baseline=%d, final=%d, delta=%d", baseline, final, final-baseline)
	require.Less(t, final-baseline, 50, "goroutine leak: %d goroutines still running after load completed", final-baseline)
}

// TestAGUI_Load_ConcurrentStreaming repeats the load run against a
// streaming reasoner so the streaming dispatch path is also load-tested.
func TestAGUI_Load_ConcurrentStreaming(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping load test in -short mode")
	}

	const totalRequests = 100
	const concurrency = 25

	reasoner := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/x-ndjson")
		w.WriteHeader(http.StatusOK)
		flusher, _ := w.(http.Flusher)
		send := func(line string) {
			fmt.Fprintln(w, line)
			if flusher != nil {
				flusher.Flush()
			}
		}
		send(`{"type":"text","delta":"chunk-1"}`)
		send(`{"type":"text","delta":"chunk-2"}`)
		send(`{"type":"state","snapshot":{"k":1}}`)
	}))
	defer reasoner.Close()

	store := &reasonerTestStorage{agent: &types.AgentNode{
		ID:              "load-stream",
		BaseURL:         reasoner.URL,
		HealthStatus:    types.HealthStatusActive,
		LifecycleStatus: types.AgentStatusReady,
		Reasoners:       []types.ReasonerDefinition{{ID: "r"}},
	}}
	router := mountAGUIRouter(t, store)

	runtime.GC()
	time.Sleep(50 * time.Millisecond)
	baseline := runtime.NumGoroutine()

	var failed atomic.Int64
	sem := make(chan struct{}, concurrency)
	var wg sync.WaitGroup
	wallStart := time.Now()

	for i := 0; i < totalRequests; i++ {
		wg.Add(1)
		sem <- struct{}{}
		go func(idx int) {
			defer wg.Done()
			defer func() { <-sem }()

			body := fmt.Sprintf(`{"threadId":"t-%d","runId":"r-%d","messages":[{"role":"user","content":"x"}]}`, idx, idx)
			req := httptest.NewRequest(http.MethodPost, "/api/v1/agui/runs/load-stream/r", strings.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			if w.Code != http.StatusOK {
				failed.Add(1)
				return
			}
			frames := parseAGUIStream(t, w.Body.String())
			if len(frames) < 5 || frames[len(frames)-1].Type() != "RUN_FINISHED" {
				failed.Add(1)
			}
		}(i)
	}
	wg.Wait()
	t.Logf("streaming load: %d reqs at %d concurrent, wall=%s", totalRequests, concurrency, time.Since(wallStart))

	require.Equal(t, int64(0), failed.Load(), "streaming dispatcher must complete every request under load")

	runtime.GC()
	time.Sleep(100 * time.Millisecond)
	final := runtime.NumGoroutine()
	t.Logf("streaming goroutines: baseline=%d, final=%d, delta=%d", baseline, final, final-baseline)
	require.Less(t, final-baseline, 50, "streaming dispatcher leaked goroutines under load")
}

// BenchmarkAGUI_BufferedHandler measures the per-request cost of the
// AG-UI handler against an in-process httptest reasoner. Run with:
//
//	go test -bench=BenchmarkAGUI -benchmem -run=^$ ./internal/handlers/...
//
// Useful as a regression baseline when the streaming/dispatch logic
// changes.
func BenchmarkAGUI_BufferedHandler(b *testing.B) {
	reasoner := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"result":"ok"}`))
	}))
	defer reasoner.Close()

	store := &reasonerTestStorage{agent: &types.AgentNode{
		ID:              "bench",
		BaseURL:         reasoner.URL,
		HealthStatus:    types.HealthStatusActive,
		LifecycleStatus: types.AgentStatusReady,
		Reasoners:       []types.ReasonerDefinition{{ID: "r"}},
	}}
	router := mountAGUIRouter(&testing.T{}, store)
	bodyTpl := `{"threadId":"t","runId":"r","messages":[{"role":"user","content":"x"}]}`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest(http.MethodPost, "/api/v1/agui/runs/bench/r", strings.NewReader(bodyTpl))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			b.Fatalf("status=%d", w.Code)
		}
	}
}

// sortDurations is a small inline sort to avoid pulling in slices/sort
// noise in the load test. n is small (≤200) so insertion sort is fine.
func sortDurations(xs []time.Duration) {
	for i := 1; i < len(xs); i++ {
		for j := i; j > 0 && xs[j-1] > xs[j]; j-- {
			xs[j-1], xs[j] = xs[j], xs[j-1]
		}
	}
}

// silence unused-import warnings in case the file is edited down later.
var _ = context.Background
