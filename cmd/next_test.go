package cmd

import (
	"bytes"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/park-jun-woo/huma/internal/adapter"
	"github.com/park-jun-woo/huma/internal/config"
	"github.com/park-jun-woo/huma/internal/hurlcheck"
	"github.com/park-jun-woo/huma/internal/runner"
	"github.com/park-jun-woo/huma/internal/scanner"
	"github.com/park-jun-woo/huma/internal/session"
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

// withProbeCheck replaces probeCheckFn for tests and restores it after.
func withProbeCheck(t *testing.T, result bool) {
	t.Helper()
	orig := probeCheckFn
	t.Cleanup(func() { probeCheckFn = orig })
	probeCheckFn = func(url string) bool { return result }
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
	want := "[A-02] Server build command failed\n  ▶ build failed"
	if err.Error() != want {
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

// withReadOnlySessionDir creates a file at .huma so session.Save() fails with MkdirAll error.
func withReadOnlySessionDir(t *testing.T) string {
	t.Helper()
	tmpDir := t.TempDir()
	orig, _ := os.Getwd()
	t.Cleanup(func() { os.Chdir(orig) })
	os.Chdir(tmpDir)
	// Create a regular file at .huma so MkdirAll fails
	os.WriteFile(filepath.Join(tmpDir, ".huma"), []byte("block"), 0o444)
	return tmpDir
}

func TestRunWithCoverage_PassSaveError(t *testing.T) {
	tmpDir := withReadOnlySessionDir(t)
	ma := &mockAdapter{}
	withMockAdapter(t, ma)
	withMockRunFn(t, func(a adapter.Adapter, hurl string, vars map[string]string, src, handler string) (*runner.Result, *adapter.CoverageResult, error) {
		return &runner.Result{Pass: true}, nil, nil
	})

	cfg := makeCfg(tmpDir)
	ep := makeEndpoint()
	sess := makeSession(ep)

	err := runWithCoverage(cfg, sess, &ep, "test.hurl")
	if err == nil {
		t.Fatal("expected save error, got nil")
	}
}

func TestRunWithCoverage_StalledSaveError(t *testing.T) {
	tmpDir := withReadOnlySessionDir(t)
	ma := &mockAdapter{}
	withMockAdapter(t, ma)
	withMockRunFn(t, func(a adapter.Adapter, hurl string, vars map[string]string, src, handler string) (*runner.Result, *adapter.CoverageResult, error) {
		return &runner.Result{Pass: true}, &adapter.CoverageResult{Covered: 5, Total: 10, Percent: 50}, nil
	})

	cfg := makeCfg(tmpDir)
	ep := makeEndpoint()
	sess := makeSession(ep)
	sess.MarkImprove("ep1", 50)
	sess.MarkImprove("ep1", 50)

	err := runWithCoverage(cfg, sess, &ep, "test.hurl")
	if err == nil {
		t.Fatal("expected save error, got nil")
	}
}

func TestRunWithCoverage_ImproveSaveError(t *testing.T) {
	tmpDir := withReadOnlySessionDir(t)
	ma := &mockAdapter{}
	withMockAdapter(t, ma)
	withMockRunFn(t, func(a adapter.Adapter, hurl string, vars map[string]string, src, handler string) (*runner.Result, *adapter.CoverageResult, error) {
		return &runner.Result{Pass: true}, &adapter.CoverageResult{Covered: 7, Total: 10, Percent: 70}, nil
	})

	cfg := makeCfg(tmpDir)
	ep := makeEndpoint()
	sess := makeSession(ep)

	err := runWithCoverage(cfg, sess, &ep, "test.hurl")
	if err == nil {
		t.Fatal("expected save error, got nil")
	}
}

// withCheckResponseCoverage replaces checkResponseCoverageFn for tests and restores it after.
func withCheckResponseCoverage(t *testing.T, fn func(*scanner.Endpoint, string, string) *hurlcheck.ResponseCoverageResult) {
	t.Helper()
	orig := checkResponseCoverageFn
	t.Cleanup(func() { checkResponseCoverageFn = orig })
	checkResponseCoverageFn = fn
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

// captureStdout captures os.Stdout during fn execution and returns the output.
func captureStdout(t *testing.T, fn func()) string {
	t.Helper()
	origStdout := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe: %v", err)
	}
	os.Stdout = w

	fn()

	w.Close()
	os.Stdout = origStdout

	var buf bytes.Buffer
	io.Copy(&buf, r)
	r.Close()
	return buf.String()
}

func TestStaticMode_PassThenNextTodoIncludesStaticPrompt(t *testing.T) {
	tmpDir := withSessionDir(t)

	// 1. Create manifest.yaml for static mode (no server.start)
	manifestContent := `apiVersion: v1
kind: manifest
backend:
  lang: go
testing:
  base_url: http://localhost:8080
  hurl_dir: hurl
`
	os.WriteFile(filepath.Join(tmpDir, "manifest.yaml"), []byte(manifestContent), 0o644)

	// 2. Create a real analyzable source so ep1 has an oracle (source link)
	//    and a non-empty branch denominator — required to reach SCAFFOLDED.
	alphaSrc := filepath.Join(tmpDir, "alpha.go")
	os.WriteFile(alphaSrc, []byte(`package main

import "net/http"

func AlphaHandler(c interface{}) {
	c.JSON(http.StatusOK, nil)
}
`), 0o644)

	// 3. Create session with two endpoints
	ep1 := scanner.Endpoint{ID: "ep1", Method: "GET", Path: "/alpha", Handler: "AlphaHandler", Source: alphaSrc, Line: 5}
	ep2 := scanner.Endpoint{ID: "ep2", Method: "POST", Path: "/beta", Handler: "BetaHandler", Source: "beta.go", Line: 1}
	sess := makeSession(ep1, ep2)
	if err := sess.Save(); err != nil {
		t.Fatalf("save session: %v", err)
	}

	// 4. Create hurl file for ep1 so it gets picked up by FindHurlFile,
	//    covering the 200 branch (non-vacuous status assertion).
	hurlDir := filepath.Join(tmpDir, "hurl")
	os.MkdirAll(hurlDir, 0o755)
	hurlFile := filepath.Join(hurlDir, "get_alpha.hurl")
	os.WriteFile(hurlFile, []byte("GET http://localhost:8080/alpha\nHTTP 200\n"), 0o644)

	// 5. Capture stdout and run the next command
	output := captureStdout(t, func() {
		err := nextCmd.RunE(nextCmd, nil)
		if err != nil {
			t.Errorf("nextCmd.RunE: %v", err)
		}
	})

	// 6. Verify ep1 was marked PASS
	sess2, err := session.Load()
	if err != nil {
		t.Fatalf("load session after run: %v", err)
	}
	for _, e := range sess2.Entries {
		if e.ID == "ep1" && e.Status != session.StatusPass {
			t.Fatalf("expected ep1 PASS, got %s", e.Status)
		}
	}

	// 7. Verify output contains ep2 info (StaticTodoPrompt output)
	if !strings.Contains(output, "POST /beta") {
		t.Fatalf("expected output to contain ep2 info 'POST /beta', got:\n%s", output)
	}
	if !strings.Contains(output, "TODO") {
		t.Fatalf("expected output to contain 'TODO', got:\n%s", output)
	}
}
