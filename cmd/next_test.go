package cmd

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/park-jun-woo/hurlfill/internal/adapter"
	"github.com/park-jun-woo/hurlfill/internal/config"
	"github.com/park-jun-woo/hurlfill/internal/runner"
	"github.com/park-jun-woo/hurlfill/internal/scanner"
	"github.com/park-jun-woo/hurlfill/internal/session"
)

// mockAdapter implements adapter.Adapter for testing.
type mockAdapter struct {
	buildErr error
}

func (m *mockAdapter) Build() error                                        { return m.buildErr }
func (m *mockAdapter) Reset() error                                        { return nil }
func (m *mockAdapter) Start() error                                        { return nil }
func (m *mockAdapter) WaitReady() error                                    { return nil }
func (m *mockAdapter) Stop() error                                         { return nil }
func (m *mockAdapter) Collect(string, int, int) (*adapter.CoverageResult, error) {
	return nil, nil
}

func makeEndpoint() scanner.Endpoint {
	return scanner.Endpoint{
		ID: "ep1", Method: "GET", Path: "/test",
		Handler: "TestHandler", Source: "test.go", Line: 1,
	}
}

func makeCfg(tmpDir string) *config.Config {
	return &config.Config{
		BaseURL: "http://localhost:9999",
		HurlDir: filepath.Join(tmpDir, "hurl"),
		Server:  config.ServerConfig{Build: "true"},
	}
}

func makeSession(eps ...scanner.Endpoint) *session.Session {
	sess := session.New()
	sess.Merge(eps)
	return sess
}

// withSessionDir sets up a temp directory and changes to it so session.Save works.
func withSessionDir(t *testing.T) string {
	t.Helper()
	tmpDir := t.TempDir()
	orig, _ := os.Getwd()
	t.Cleanup(func() { os.Chdir(orig) })
	os.Chdir(tmpDir)
	return tmpDir
}

// withMockAdapter replaces newAdapterFn for tests and restores it after.
func withMockAdapter(t *testing.T, ma adapter.Adapter) {
	t.Helper()
	orig := newAdapterFn
	t.Cleanup(func() { newAdapterFn = orig })
	newAdapterFn = func(cfg *config.Config) adapter.Adapter { return ma }
}

// withMockRunFn replaces adapterRunFn for tests and restores it after.
func withMockRunFn(t *testing.T, fn func(adapter.Adapter, string, map[string]string, string, string) (*runner.Result, *adapter.CoverageResult, error)) {
	t.Helper()
	orig := adapterRunFn
	t.Cleanup(func() { adapterRunFn = orig })
	adapterRunFn = fn
}

func TestRunWithCoverage_BuildFailure(t *testing.T) {
	tmpDir := withSessionDir(t)
	ma := &mockAdapter{buildErr: errors.New("build failed")}
	withMockAdapter(t, ma)

	cfg := makeCfg(tmpDir)
	ep := makeEndpoint()
	sess := makeSession(ep)

	err := runWithCoverage(cfg, sess, &ep, "test.hurl")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err.Error() != "build: build failed" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRunWithCoverage_RunError(t *testing.T) {
	tmpDir := withSessionDir(t)
	ma := &mockAdapter{}
	withMockAdapter(t, ma)
	withMockRunFn(t, func(a adapter.Adapter, hurl string, vars map[string]string, src, handler string) (*runner.Result, *adapter.CoverageResult, error) {
		return nil, nil, errors.New("run failed")
	})

	cfg := makeCfg(tmpDir)
	ep := makeEndpoint()
	sess := makeSession(ep)

	err := runWithCoverage(cfg, sess, &ep, "test.hurl")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err.Error() != "run failed" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRunWithCoverage_HurlFail(t *testing.T) {
	tmpDir := withSessionDir(t)
	ma := &mockAdapter{}
	withMockAdapter(t, ma)
	withMockRunFn(t, func(a adapter.Adapter, hurl string, vars map[string]string, src, handler string) (*runner.Result, *adapter.CoverageResult, error) {
		return &runner.Result{Pass: false, Feedback: "assertion failed"}, nil, nil
	})

	cfg := makeCfg(tmpDir)
	ep := makeEndpoint()
	sess := makeSession(ep)

	err := runWithCoverage(cfg, sess, &ep, "test.hurl")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Status should still be TODO since hurl failed
	entry := sess.CurrentEntry()
	if entry == nil || entry.Status != session.StatusTodo {
		t.Fatal("expected entry to remain TODO")
	}
}

func TestRunWithCoverage_PassWithNilCoverage(t *testing.T) {
	tmpDir := withSessionDir(t)
	ma := &mockAdapter{}
	withMockAdapter(t, ma)
	withMockRunFn(t, func(a adapter.Adapter, hurl string, vars map[string]string, src, handler string) (*runner.Result, *adapter.CoverageResult, error) {
		return &runner.Result{Pass: true}, nil, nil
	})

	cfg := makeCfg(tmpDir)
	ep := makeEndpoint()
	sess := makeSession(ep)

	err := runWithCoverage(cfg, sess, &ep, "test.hurl")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Should be marked PASS since covResult is nil
	for _, e := range sess.Entries {
		if e.ID == "ep1" && e.Status != session.StatusPass {
			t.Fatalf("expected PASS, got %s", e.Status)
		}
	}
}

func TestRunWithCoverage_PassWith100Coverage(t *testing.T) {
	tmpDir := withSessionDir(t)
	ma := &mockAdapter{}
	withMockAdapter(t, ma)
	withMockRunFn(t, func(a adapter.Adapter, hurl string, vars map[string]string, src, handler string) (*runner.Result, *adapter.CoverageResult, error) {
		return &runner.Result{Pass: true}, &adapter.CoverageResult{Covered: 10, Total: 10, Percent: 100}, nil
	})

	cfg := makeCfg(tmpDir)
	ep := makeEndpoint()
	sess := makeSession(ep)

	err := runWithCoverage(cfg, sess, &ep, "test.hurl")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, e := range sess.Entries {
		if e.ID == "ep1" && e.Status != session.StatusPass {
			t.Fatalf("expected PASS, got %s", e.Status)
		}
	}
}

func TestRunWithCoverage_PassWithZeroTotal(t *testing.T) {
	tmpDir := withSessionDir(t)
	ma := &mockAdapter{}
	withMockAdapter(t, ma)
	withMockRunFn(t, func(a adapter.Adapter, hurl string, vars map[string]string, src, handler string) (*runner.Result, *adapter.CoverageResult, error) {
		return &runner.Result{Pass: true}, &adapter.CoverageResult{Covered: 0, Total: 0, Percent: 0}, nil
	})

	cfg := makeCfg(tmpDir)
	ep := makeEndpoint()
	sess := makeSession(ep)

	err := runWithCoverage(cfg, sess, &ep, "test.hurl")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, e := range sess.Entries {
		if e.ID == "ep1" && e.Status != session.StatusPass {
			t.Fatalf("expected PASS, got %s", e.Status)
		}
	}
}

func TestRunWithCoverage_StalledImprovement(t *testing.T) {
	tmpDir := withSessionDir(t)
	ma := &mockAdapter{}
	withMockAdapter(t, ma)
	// Return 50% coverage — same as PrevCoverage, so stalled
	withMockRunFn(t, func(a adapter.Adapter, hurl string, vars map[string]string, src, handler string) (*runner.Result, *adapter.CoverageResult, error) {
		return &runner.Result{Pass: true}, &adapter.CoverageResult{Covered: 5, Total: 10, Percent: 50}, nil
	})

	cfg := makeCfg(tmpDir)
	ep := makeEndpoint()
	sess := makeSession(ep)
	// After MarkImprove("ep1", 50): PrevCoverage=0, Coverage=50, ImproveCount=1
	// The stall check: ImproveCount >= 1 && covResult.Percent <= entry.PrevCoverage
	// We need PrevCoverage to be >= 50. So we need a second MarkImprove.
	// After MarkImprove("ep1", 50) again: PrevCoverage=50, Coverage=50, ImproveCount=2
	// Now check: 2 >= 1 && 50 <= 50 => true => DONE
	sess.MarkImprove("ep1", 50)
	sess.MarkImprove("ep1", 50)

	err := runWithCoverage(cfg, sess, &ep, "test.hurl")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Should be marked DONE since coverage didn't improve
	for _, e := range sess.Entries {
		if e.ID == "ep1" {
			if e.Status != session.StatusDone {
				t.Fatalf("expected DONE, got %s", e.Status)
			}
			return
		}
	}
	t.Fatal("entry not found")
}

func TestRunWithCoverage_ImprovedCoverage(t *testing.T) {
	tmpDir := withSessionDir(t)
	ma := &mockAdapter{}
	withMockAdapter(t, ma)
	withMockRunFn(t, func(a adapter.Adapter, hurl string, vars map[string]string, src, handler string) (*runner.Result, *adapter.CoverageResult, error) {
		return &runner.Result{Pass: true}, &adapter.CoverageResult{Covered: 7, Total: 10, Percent: 70}, nil
	})

	cfg := makeCfg(tmpDir)
	ep := makeEndpoint()
	sess := makeSession(ep)

	err := runWithCoverage(cfg, sess, &ep, "test.hurl")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Should be marked IMPROVE since this is first attempt
	for _, e := range sess.Entries {
		if e.ID == "ep1" {
			if e.Status != session.StatusImprove {
				t.Fatalf("expected IMPROVE, got %s", e.Status)
			}
			return
		}
	}
	t.Fatal("entry not found")
}

func TestRunWithCoverage_PassAllComplete(t *testing.T) {
	tmpDir := withSessionDir(t)
	ma := &mockAdapter{}
	withMockAdapter(t, ma)
	withMockRunFn(t, func(a adapter.Adapter, hurl string, vars map[string]string, src, handler string) (*runner.Result, *adapter.CoverageResult, error) {
		return &runner.Result{Pass: true}, nil, nil
	})

	cfg := makeCfg(tmpDir)
	ep := makeEndpoint()
	// Only one endpoint - after marking PASS, next is nil (all complete)
	sess := makeSession(ep)

	err := runWithCoverage(cfg, sess, &ep, "test.hurl")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if next := sess.Current(); next != nil {
		t.Fatal("expected all complete (Current() == nil)")
	}
}

func TestRunWithCoverage_PassThenNextTodo(t *testing.T) {
	tmpDir := withSessionDir(t)
	ma := &mockAdapter{}
	withMockAdapter(t, ma)
	withMockRunFn(t, func(a adapter.Adapter, hurl string, vars map[string]string, src, handler string) (*runner.Result, *adapter.CoverageResult, error) {
		return &runner.Result{Pass: true}, nil, nil
	})

	cfg := makeCfg(tmpDir)
	ep1 := scanner.Endpoint{ID: "ep1", Method: "GET", Path: "/test1", Handler: "H1", Source: "a.go", Line: 1}
	ep2 := scanner.Endpoint{ID: "ep2", Method: "POST", Path: "/test2", Handler: "H2", Source: "b.go", Line: 1}
	sess := makeSession(ep1, ep2)

	err := runWithCoverage(cfg, sess, &ep1, "test.hurl")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	next := sess.Current()
	if next == nil {
		t.Fatal("expected next endpoint")
	}
	if next.ID != "ep2" {
		t.Fatalf("expected ep2, got %s", next.ID)
	}
}

func TestRunWithCoverage_StalledThenAllComplete(t *testing.T) {
	tmpDir := withSessionDir(t)
	ma := &mockAdapter{}
	withMockAdapter(t, ma)
	withMockRunFn(t, func(a adapter.Adapter, hurl string, vars map[string]string, src, handler string) (*runner.Result, *adapter.CoverageResult, error) {
		return &runner.Result{Pass: true}, &adapter.CoverageResult{Covered: 5, Total: 10, Percent: 50}, nil
	})

	cfg := makeCfg(tmpDir)
	ep := makeEndpoint()
	sess := makeSession(ep)
	// Two improve calls so PrevCoverage=50, ImproveCount=2
	sess.MarkImprove("ep1", 50)
	sess.MarkImprove("ep1", 50)

	err := runWithCoverage(cfg, sess, &ep, "test.hurl")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Should be DONE and all complete
	if next := sess.Current(); next != nil {
		t.Fatal("expected all complete")
	}
}

func TestRunWithCoverage_StalledThenNextTodo(t *testing.T) {
	tmpDir := withSessionDir(t)
	ma := &mockAdapter{}
	withMockAdapter(t, ma)
	withMockRunFn(t, func(a adapter.Adapter, hurl string, vars map[string]string, src, handler string) (*runner.Result, *adapter.CoverageResult, error) {
		return &runner.Result{Pass: true}, &adapter.CoverageResult{Covered: 5, Total: 10, Percent: 50}, nil
	})

	cfg := makeCfg(tmpDir)
	ep1 := scanner.Endpoint{ID: "ep1", Method: "GET", Path: "/a", Handler: "H1", Source: "a.go", Line: 1}
	ep2 := scanner.Endpoint{ID: "ep2", Method: "GET", Path: "/b", Handler: "H2", Source: "b.go", Line: 1}
	sess := makeSession(ep1, ep2)
	// Two improve calls so PrevCoverage=50, ImproveCount=2
	sess.MarkImprove("ep1", 50)
	sess.MarkImprove("ep1", 50)

	err := runWithCoverage(cfg, sess, &ep1, "test.hurl")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	next := sess.Current()
	if next == nil || next.ID != "ep2" {
		t.Fatal("expected ep2 as next")
	}
}
