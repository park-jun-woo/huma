package humaquest

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/park-jun-woo/huma/internal/adapter"
	"github.com/park-jun-woo/huma/internal/config"
)

// fakeAdapter is a spy implementing adapter.Adapter for the setup-lifecycle tests.
// It records which lifecycle hooks ran and can inject errors at each stage.
type fakeAdapter struct {
	buildErr error
	startErr error
	waitErr  error
	stopErr  error

	built, started, waited, stopped bool
}

func (f *fakeAdapter) Build() error { f.built = true; return f.buildErr }
func (f *fakeAdapter) Start() error { f.started = true; return f.startErr }
func (f *fakeAdapter) WaitReady() error {
	f.waited = true
	return f.waitErr
}
func (f *fakeAdapter) Stop() error { f.stopped = true; return f.stopErr }
func (f *fakeAdapter) Collect(string, int, int) (*adapter.CoverageResult, error) {
	return nil, nil
}
func (f *fakeAdapter) Reset() error { return nil }

func TestCaptureSetup_NoSetup(t *testing.T) {
	f := &fakeAdapter{}
	cfg := &config.Config{} // Setup.Hurl == ""
	vars, err := captureSetup(f, cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(vars) != 0 {
		t.Fatalf("expected empty map, got %v", vars)
	}
	if f.built || f.started {
		t.Fatal("server lifecycle must not run when no setup is configured")
	}
}

func TestCaptureSetup_BuildError(t *testing.T) {
	f := &fakeAdapter{buildErr: errors.New("boom")}
	cfg := &config.Config{}
	cfg.Setup.Hurl = "setup.hurl"
	_, err := captureSetup(f, cfg)
	if err == nil || !strings.Contains(err.Error(), "build server") {
		t.Fatalf("expected build error, got %v", err)
	}
}

func TestCaptureSetup_StartError(t *testing.T) {
	f := &fakeAdapter{startErr: errors.New("boom")}
	cfg := &config.Config{}
	cfg.Setup.Hurl = "setup.hurl"
	_, err := captureSetup(f, cfg)
	if err == nil || !strings.Contains(err.Error(), "start server") {
		t.Fatalf("expected start error, got %v", err)
	}
}

func TestCaptureSetup_WaitReadyError(t *testing.T) {
	f := &fakeAdapter{waitErr: errors.New("boom")}
	cfg := &config.Config{}
	cfg.Setup.Hurl = "setup.hurl"
	_, err := captureSetup(f, cfg)
	if err == nil || !strings.Contains(err.Error(), "wait ready") {
		t.Fatalf("expected wait-ready error, got %v", err)
	}
	if !f.stopped {
		t.Fatal("Stop (deferred) must run after Start succeeds")
	}
}

func TestCaptureSetup_RunJSONError(t *testing.T) {
	// Build/Start/WaitReady all succeed; RunJSON fails because the .hurl file
	// does not exist. This reaches the final error branch and exercises Stop().
	f := &fakeAdapter{}
	cfg := &config.Config{}
	cfg.Setup.Hurl = "/nonexistent/setup.hurl"
	_, err := captureSetup(f, cfg)
	if err == nil || !strings.Contains(err.Error(), "setup capture") {
		t.Fatalf("expected setup-capture error, got %v", err)
	}
	if !f.stopped {
		t.Fatal("Stop (deferred) must run")
	}
}

func TestSetupVars_NoSetupNoAuth(t *testing.T) {
	f := &fakeAdapter{}
	cfg := &config.Config{} // neither Setup.Hurl nor Auth.SecretEnv
	var buf bytes.Buffer
	vars := setupVars(f, cfg, &buf)
	if len(vars) != 0 {
		t.Fatalf("expected empty map, got %v", vars)
	}
	if buf.Len() != 0 {
		t.Fatalf("expected no log output, got %q", buf.String())
	}
}

func TestSetupVars_MintPath(t *testing.T) {
	f := &fakeAdapter{}
	cfg := &config.Config{}
	cfg.Auth.SecretEnv = "HUMA_TEST_SETUPVARS_SECRET"
	cfg.Auth.Claims = map[string]string{"sub": "bob"}
	t.Setenv("HUMA_TEST_SETUPVARS_SECRET", "secret-value")

	var buf bytes.Buffer
	vars := setupVars(f, cfg, &buf)
	if vars["token"] == "" {
		t.Fatalf("expected minted token, got %v", vars)
	}
	// Logs KEYS (token) but never the value.
	out := buf.String()
	if !strings.Contains(out, "mint(testing.auth)") || !strings.Contains(out, "token") {
		t.Fatalf("expected key-only log, got %q", out)
	}
	if strings.Contains(out, vars["token"]) {
		t.Fatal("token VALUE must never be logged")
	}
}

func TestSetupVars_MintError(t *testing.T) {
	// Auth.SecretEnv set but env var empty -> mintToken errors ->
	// setupVars warns and continues token-less (empty map).
	f := &fakeAdapter{}
	cfg := &config.Config{}
	cfg.Auth.SecretEnv = "HUMA_TEST_SETUPVARS_MISSING"
	t.Setenv("HUMA_TEST_SETUPVARS_MISSING", "")

	var buf bytes.Buffer
	vars := setupVars(f, cfg, &buf)
	if len(vars) != 0 {
		t.Fatalf("expected empty map after mint failure, got %v", vars)
	}
	if !strings.Contains(buf.String(), "warning:") {
		t.Fatalf("expected loud warning, got %q", buf.String())
	}
}

func TestSetupVars_CapturePath(t *testing.T) {
	// Setup.Hurl set takes precedence over Auth. captureSetup will fail here
	// (missing file) -> setupVars warns and returns empty map. This proves the
	// capture branch is selected (path label "capture(testing.setup)").
	f := &fakeAdapter{}
	cfg := &config.Config{}
	cfg.Setup.Hurl = "/nonexistent/setup.hurl"
	cfg.Auth.SecretEnv = "HUMA_TEST_SETUPVARS_SECRET2" // must be ignored

	var buf bytes.Buffer
	vars := setupVars(f, cfg, &buf)
	if len(vars) != 0 {
		t.Fatalf("expected empty map, got %v", vars)
	}
	out := buf.String()
	if !strings.Contains(out, "capture(testing.setup)") {
		t.Fatalf("expected capture path selected, got %q", out)
	}
}
