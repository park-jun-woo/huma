package adapter

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

type lifecycleMock struct {
	resetErr   error
	startErr   error
	waitErr    error
	stopErr    error
	collectRes *CoverageResult
	collectErr error
}

func (m *lifecycleMock) Build() error { return nil }
func (m *lifecycleMock) Reset() error { return m.resetErr }
func (m *lifecycleMock) Start() error { return m.startErr }
func (m *lifecycleMock) WaitReady() error { return m.waitErr }
func (m *lifecycleMock) Stop() error  { return m.stopErr }
func (m *lifecycleMock) Collect(string, int, int) (*CoverageResult, error) {
	return m.collectRes, m.collectErr
}

func TestRunWithCoverage_ResetError(t *testing.T) {
	m := &lifecycleMock{resetErr: errors.New("reset failed")}
	_, _, err := RunWithCoverage(m, "test.hurl", map[string]string{"base_url": "http://localhost:8080"}, "handler.go", "Handler")
	if err == nil {
		t.Fatal("expected error")
	}
	if err.Error() != "reset coverage: reset failed" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRunWithCoverage_StartError(t *testing.T) {
	m := &lifecycleMock{startErr: errors.New("start failed")}
	_, _, err := RunWithCoverage(m, "test.hurl", map[string]string{"base_url": "http://localhost:8080"}, "handler.go", "Handler")
	if err == nil {
		t.Fatal("expected error")
	}
	if err.Error() != "start server: start failed" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRunWithCoverage_WaitReadyError(t *testing.T) {
	m := &lifecycleMock{waitErr: errors.New("not ready")}
	_, _, err := RunWithCoverage(m, "test.hurl", map[string]string{"base_url": "http://localhost:8080"}, "handler.go", "Handler")
	if err == nil {
		t.Fatal("expected error")
	}
	if err.Error() != "wait ready: not ready" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRunWithCoverage_HurlFail(t *testing.T) {
	// Hurl with nonexistent file returns a fail result (not an error)
	m := &lifecycleMock{}
	result, covResult, err := RunWithCoverage(m, "/nonexistent/test.hurl", map[string]string{"base_url": "http://localhost:99999"}, "handler.go", "Handler")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result == nil {
		t.Fatal("expected non-nil result")
	}
	if result.Pass {
		t.Fatal("expected pass=false")
	}
	if covResult != nil {
		t.Fatal("expected nil covResult for failed hurl")
	}
}

func TestRunWithCoverage_StopError(t *testing.T) {
	// Hurl returns a fail result (err=nil), but Stop returns an error
	m := &lifecycleMock{stopErr: errors.New("stop failed")}
	_, _, err := RunWithCoverage(m, "/nonexistent/test.hurl", map[string]string{"base_url": "http://localhost:19999"}, "handler.go", "Handler")
	if err == nil {
		t.Fatal("expected stop error")
	}
	if err.Error() != "stop server: stop failed" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRunWithCoverage_HurlPassReadHandlerError(t *testing.T) {
	// Start a real HTTP test server so hurl can pass
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"ok": true}`))
	}))
	defer ts.Close()

	tmpDir := t.TempDir()
	hurlFile := filepath.Join(tmpDir, "test.hurl")
	// Simple hurl test that should pass
	hurlContent := fmt.Sprintf("GET %s/test\nHTTP 200\n", ts.URL)
	os.WriteFile(hurlFile, []byte(hurlContent), 0o644)

	m := &lifecycleMock{}
	// handler file doesn't exist, so ReadHandler will fail => skip coverage
	result, covResult, err := RunWithCoverage(m, hurlFile, nil, "/nonexistent/handler.go", "Handler")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result == nil || !result.Pass {
		t.Fatal("expected pass")
	}
	if covResult != nil {
		t.Fatal("expected nil covResult when ReadHandler fails")
	}
}

func TestRunWithCoverage_HurlPassCollectError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	tmpDir := t.TempDir()
	hurlFile := filepath.Join(tmpDir, "test.hurl")
	hurlContent := fmt.Sprintf("GET %s/test\nHTTP 200\n", ts.URL)
	os.WriteFile(hurlFile, []byte(hurlContent), 0o644)

	// Create a source file with a handler function
	handlerFile := filepath.Join(tmpDir, "handler.go")
	handlerContent := `package main

func Handler(c interface{}) {
	// handler body
}
`
	os.WriteFile(handlerFile, []byte(handlerContent), 0o644)

	m := &lifecycleMock{collectErr: errors.New("collect failed")}
	_, _, err := RunWithCoverage(m, hurlFile, nil, handlerFile, "Handler")
	if err == nil {
		t.Fatal("expected collect error")
	}
	if err.Error() != "collect coverage: collect failed" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRunWithCoverage_HurlPassCollectSuccess(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	tmpDir := t.TempDir()
	hurlFile := filepath.Join(tmpDir, "test.hurl")
	hurlContent := fmt.Sprintf("GET %s/test\nHTTP 200\n", ts.URL)
	os.WriteFile(hurlFile, []byte(hurlContent), 0o644)

	handlerFile := filepath.Join(tmpDir, "handler.go")
	handlerContent := `package main

func Handler(c interface{}) {
	// handler body
}
`
	os.WriteFile(handlerFile, []byte(handlerContent), 0o644)

	covRes := &CoverageResult{Covered: 3, Total: 4, Percent: 75}
	m := &lifecycleMock{collectRes: covRes}
	result, gotCov, err := RunWithCoverage(m, hurlFile, nil, handlerFile, "Handler")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result == nil || !result.Pass {
		t.Fatal("expected pass")
	}
	if gotCov == nil {
		t.Fatal("expected coverage result")
	}
	if gotCov.Percent != 75 {
		t.Fatalf("expected 75%%, got %.0f%%", gotCov.Percent)
	}
}
