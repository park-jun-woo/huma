package cmd

import (
	"errors"
	"os"
	"testing"

	"github.com/park-jun-woo/hurlfill/internal/adapter"
	"github.com/park-jun-woo/hurlfill/internal/runner"
	"github.com/park-jun-woo/hurlfill/internal/session"
)

func TestVerifyWithCoverage_BuildFailure(t *testing.T) {
	tmpDir := withSessionDir(t)
	ma := &mockAdapter{buildErr: errors.New("build failed")}
	withMockAdapter(t, ma)

	cfg := makeCfg(tmpDir)
	ep := makeEndpoint()
	sess := makeSession(ep)

	err := verifyWithCoverage(cfg, sess, &ep, "test.hurl")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err.Error() != "build: build failed" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestVerifyWithCoverage_RunError(t *testing.T) {
	tmpDir := withSessionDir(t)
	ma := &mockAdapter{}
	withMockAdapter(t, ma)
	withMockRunFn(t, func(a adapter.Adapter, hurl, base, src, handler string) (*runner.Result, *adapter.CoverageResult, error) {
		return nil, nil, errors.New("run failed")
	})

	cfg := makeCfg(tmpDir)
	ep := makeEndpoint()
	sess := makeSession(ep)

	err := verifyWithCoverage(cfg, sess, &ep, "test.hurl")
	if err == nil || err.Error() != "run failed" {
		t.Fatalf("expected 'run failed', got %v", err)
	}
}

func TestVerifyWithCoverage_HurlFail(t *testing.T) {
	tmpDir := withSessionDir(t)
	ma := &mockAdapter{}
	withMockAdapter(t, ma)
	withMockRunFn(t, func(a adapter.Adapter, hurl, base, src, handler string) (*runner.Result, *adapter.CoverageResult, error) {
		return &runner.Result{Pass: false, Feedback: "assertion failed"}, nil, nil
	})

	cfg := makeCfg(tmpDir)
	ep := makeEndpoint()
	sess := makeSession(ep)

	err := verifyWithCoverage(cfg, sess, &ep, "test.hurl")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	entry := sess.CurrentEntry()
	if entry == nil || entry.Status != session.StatusTodo {
		t.Fatal("expected entry to remain TODO")
	}
}

func TestVerifyWithCoverage_PassNilCoverage(t *testing.T) {
	tmpDir := withSessionDir(t)
	ma := &mockAdapter{}
	withMockAdapter(t, ma)
	withMockRunFn(t, func(a adapter.Adapter, hurl, base, src, handler string) (*runner.Result, *adapter.CoverageResult, error) {
		return &runner.Result{Pass: true}, nil, nil
	})

	cfg := makeCfg(tmpDir)
	ep := makeEndpoint()
	sess := makeSession(ep)

	err := verifyWithCoverage(cfg, sess, &ep, "test.hurl")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, e := range sess.Entries {
		if e.ID == "ep1" && e.Status != session.StatusPass {
			t.Fatalf("expected PASS, got %s", e.Status)
		}
	}
}

func TestVerifyWithCoverage_Pass100Pct(t *testing.T) {
	tmpDir := withSessionDir(t)
	ma := &mockAdapter{}
	withMockAdapter(t, ma)
	withMockRunFn(t, func(a adapter.Adapter, hurl, base, src, handler string) (*runner.Result, *adapter.CoverageResult, error) {
		return &runner.Result{Pass: true}, &adapter.CoverageResult{Covered: 10, Total: 10, Percent: 100}, nil
	})

	cfg := makeCfg(tmpDir)
	ep := makeEndpoint()
	sess := makeSession(ep)

	err := verifyWithCoverage(cfg, sess, &ep, "test.hurl")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, e := range sess.Entries {
		if e.ID == "ep1" && e.Status != session.StatusPass {
			t.Fatalf("expected PASS, got %s", e.Status)
		}
	}
}

func TestVerifyWithCoverage_PassZeroTotal(t *testing.T) {
	tmpDir := withSessionDir(t)
	ma := &mockAdapter{}
	withMockAdapter(t, ma)
	withMockRunFn(t, func(a adapter.Adapter, hurl, base, src, handler string) (*runner.Result, *adapter.CoverageResult, error) {
		return &runner.Result{Pass: true}, &adapter.CoverageResult{Covered: 0, Total: 0, Percent: 0}, nil
	})

	cfg := makeCfg(tmpDir)
	ep := makeEndpoint()
	sess := makeSession(ep)

	err := verifyWithCoverage(cfg, sess, &ep, "test.hurl")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, e := range sess.Entries {
		if e.ID == "ep1" && e.Status != session.StatusPass {
			t.Fatalf("expected PASS, got %s", e.Status)
		}
	}
}

func TestVerifyWithCoverage_Stalled(t *testing.T) {
	tmpDir := withSessionDir(t)
	ma := &mockAdapter{}
	withMockAdapter(t, ma)
	withMockRunFn(t, func(a adapter.Adapter, hurl, base, src, handler string) (*runner.Result, *adapter.CoverageResult, error) {
		return &runner.Result{Pass: true}, &adapter.CoverageResult{Covered: 5, Total: 10, Percent: 50}, nil
	})

	cfg := makeCfg(tmpDir)
	ep := makeEndpoint()
	sess := makeSession(ep)
	sess.MarkImprove("ep1", 50)
	sess.MarkImprove("ep1", 50)

	err := verifyWithCoverage(cfg, sess, &ep, "test.hurl")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
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

func TestVerifyWithCoverage_Improve(t *testing.T) {
	tmpDir := withSessionDir(t)
	ma := &mockAdapter{}
	withMockAdapter(t, ma)
	withMockRunFn(t, func(a adapter.Adapter, hurl, base, src, handler string) (*runner.Result, *adapter.CoverageResult, error) {
		return &runner.Result{Pass: true}, &adapter.CoverageResult{Covered: 7, Total: 10, Percent: 70}, nil
	})

	cfg := makeCfg(tmpDir)
	ep := makeEndpoint()
	sess := makeSession(ep)

	err := verifyWithCoverage(cfg, sess, &ep, "test.hurl")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
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

func TestVerifyWithCoverage_SaveErrorOnPass(t *testing.T) {
	withSessionDir(t)
	ma := &mockAdapter{}
	withMockAdapter(t, ma)
	withMockRunFn(t, func(a adapter.Adapter, hurl, base, src, handler string) (*runner.Result, *adapter.CoverageResult, error) {
		return &runner.Result{Pass: true}, nil, nil
	})

	cfg := makeCfg("")
	ep := makeEndpoint()
	sess := makeSession(ep)

	// Make .hurlfill directory read-only to cause Save() to fail
	os.MkdirAll(".hurlfill", 0o755)
	os.Chmod(".hurlfill", 0o444)
	t.Cleanup(func() { os.Chmod(".hurlfill", 0o755) })

	err := verifyWithCoverage(cfg, sess, &ep, "test.hurl")
	if err == nil {
		t.Fatal("expected save error, got nil")
	}
}

func TestVerifyWithCoverage_SaveErrorOnStalled(t *testing.T) {
	withSessionDir(t)
	ma := &mockAdapter{}
	withMockAdapter(t, ma)
	withMockRunFn(t, func(a adapter.Adapter, hurl, base, src, handler string) (*runner.Result, *adapter.CoverageResult, error) {
		return &runner.Result{Pass: true}, &adapter.CoverageResult{Covered: 5, Total: 10, Percent: 50}, nil
	})

	cfg := makeCfg("")
	ep := makeEndpoint()
	sess := makeSession(ep)
	sess.MarkImprove("ep1", 50)
	sess.MarkImprove("ep1", 50)

	os.MkdirAll(".hurlfill", 0o755)
	os.Chmod(".hurlfill", 0o444)
	t.Cleanup(func() { os.Chmod(".hurlfill", 0o755) })

	err := verifyWithCoverage(cfg, sess, &ep, "test.hurl")
	if err == nil {
		t.Fatal("expected save error, got nil")
	}
}

func TestVerifyWithCoverage_SaveErrorOnImprove(t *testing.T) {
	withSessionDir(t)
	ma := &mockAdapter{}
	withMockAdapter(t, ma)
	withMockRunFn(t, func(a adapter.Adapter, hurl, base, src, handler string) (*runner.Result, *adapter.CoverageResult, error) {
		return &runner.Result{Pass: true}, &adapter.CoverageResult{Covered: 7, Total: 10, Percent: 70}, nil
	})

	cfg := makeCfg("")
	ep := makeEndpoint()
	sess := makeSession(ep)

	os.MkdirAll(".hurlfill", 0o755)
	os.Chmod(".hurlfill", 0o444)
	t.Cleanup(func() { os.Chmod(".hurlfill", 0o755) })

	err := verifyWithCoverage(cfg, sess, &ep, "test.hurl")
	if err == nil {
		t.Fatal("expected save error, got nil")
	}
}
