package runner

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/park-jun-woo/hurlfill/internal/scanner"
)

func TestFindHurlFile_InHurlDir(t *testing.T) {
	tmpDir := t.TempDir()
	hurlDir := filepath.Join(tmpDir, "hurl")
	os.MkdirAll(hurlDir, 0o755)

	ep := &scanner.Endpoint{Method: "GET", Path: "/users"}
	name := hurlFileName(ep)
	os.WriteFile(filepath.Join(hurlDir, name), []byte("GET /test\nHTTP 200\n"), 0o644)

	found := FindHurlFile(ep, hurlDir)
	if found == "" {
		t.Fatal("expected to find hurl file")
	}
}

func TestFindHurlFile_InCurrentDir(t *testing.T) {
	tmpDir := t.TempDir()
	orig, _ := os.Getwd()
	t.Cleanup(func() { os.Chdir(orig) })
	os.Chdir(tmpDir)

	ep := &scanner.Endpoint{Method: "POST", Path: "/users"}
	name := hurlFileName(ep)
	os.WriteFile(filepath.Join(tmpDir, name), []byte("POST /test\nHTTP 201\n"), 0o644)

	found := FindHurlFile(ep, "nonexistent_hurl_dir")
	if found == "" {
		t.Fatal("expected to find hurl file in current dir")
	}
}

func TestFindHurlFile_NotFound(t *testing.T) {
	ep := &scanner.Endpoint{Method: "DELETE", Path: "/users/:id"}
	found := FindHurlFile(ep, "/nonexistent/dir")
	if found != "" {
		t.Fatalf("expected empty, got %s", found)
	}
}

func TestRun_Pass(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"ok": true}`))
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
		t.Fatalf("expected pass, feedback: %s", result.Feedback)
	}
}

func TestRun_Fail(t *testing.T) {
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
	if result.Feedback == "" {
		t.Fatal("expected feedback")
	}
}

func TestRun_HurlNotFound(t *testing.T) {
	// Temporarily remove hurl from PATH to trigger execution error
	origPath := os.Getenv("PATH")
	os.Setenv("PATH", "/nonexistent")
	t.Cleanup(func() { os.Setenv("PATH", origPath) })

	_, err := Run("test.hurl", nil)
	if err == nil {
		t.Fatal("expected error when hurl is not in PATH")
	}
}

func TestHurlFileName_Public(t *testing.T) {
	ep := &scanner.Endpoint{Method: "GET", Path: "/users/:id"}
	result := HurlFileName(ep, "hurl")
	expected := filepath.Join("hurl", "get_users_id.hurl")
	if result != expected {
		t.Fatalf("expected %s, got %s", expected, result)
	}
}

func TestHurlFileName_Private(t *testing.T) {
	tests := []struct {
		method string
		path   string
		want   string
	}{
		{"GET", "/users", "get_users.hurl"},
		{"POST", "/users", "post_users.hurl"},
		{"PUT", "/users/:id", "put_users_id.hurl"},
		{"DELETE", "/api/v1/items/:id", "delete_api_v1_items_id.hurl"},
	}
	for _, tt := range tests {
		ep := &scanner.Endpoint{Method: tt.method, Path: tt.path}
		got := hurlFileName(ep)
		if got != tt.want {
			t.Errorf("hurlFileName(%s %s) = %s, want %s", tt.method, tt.path, got, tt.want)
		}
	}
}
