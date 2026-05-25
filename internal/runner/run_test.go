package runner

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestRun2_Pass(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	tmpDir := t.TempDir()
	hurlFile := filepath.Join(tmpDir, "test.hurl")
	content := fmt.Sprintf("GET %s/test\nHTTP 200\n", ts.URL)
	os.WriteFile(hurlFile, []byte(content), 0o644)

	result, err := Run(hurlFile, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Pass {
		t.Fatal("expected pass")
	}
}

func TestRun2_Fail(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer ts.Close()

	tmpDir := t.TempDir()
	hurlFile := filepath.Join(tmpDir, "test.hurl")
	content := fmt.Sprintf("GET %s/test\nHTTP 200\n", ts.URL)
	os.WriteFile(hurlFile, []byte(content), 0o644)

	result, err := Run(hurlFile, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Pass {
		t.Fatal("expected fail")
	}
}

func TestRun2_WithVariables(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	tmpDir := t.TempDir()
	hurlFile := filepath.Join(tmpDir, "test.hurl")
	content := fmt.Sprintf("GET %s/test\nHTTP 200\n", ts.URL)
	os.WriteFile(hurlFile, []byte(content), 0o644)

	vars := map[string]string{"base_url": ts.URL}
	result, err := Run(hurlFile, vars)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Pass {
		t.Fatal("expected pass")
	}
}

func TestRun2_HurlNotFound(t *testing.T) {
	origPath := os.Getenv("PATH")
	os.Setenv("PATH", "")
	t.Cleanup(func() { os.Setenv("PATH", origPath) })

	_, err := Run("test.hurl", nil)
	if err == nil {
		t.Fatal("expected error when hurl binary not found")
	}
}
