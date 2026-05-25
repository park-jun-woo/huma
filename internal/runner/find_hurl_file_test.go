package runner

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/park-jun-woo/huma/internal/scanner"
)

func TestFindHurlFile2_InHurlDir(t *testing.T) {
	tmpDir := t.TempDir()
	hurlDir := filepath.Join(tmpDir, "hurl")
	os.MkdirAll(hurlDir, 0o755)

	ep := &scanner.Endpoint{Method: "GET", Path: "/users"}
	name := hurlFileName(ep)
	os.WriteFile(filepath.Join(hurlDir, name), []byte("GET /users\nHTTP 200\n"), 0o644)

	result := FindHurlFile(ep, hurlDir)
	if result == "" {
		t.Fatal("expected to find hurl file")
	}
}

func TestFindHurlFile2_InCurrentDir(t *testing.T) {
	tmpDir := t.TempDir()
	orig, _ := os.Getwd()
	t.Cleanup(func() { os.Chdir(orig) })
	os.Chdir(tmpDir)

	ep := &scanner.Endpoint{Method: "GET", Path: "/users"}
	name := hurlFileName(ep)
	os.WriteFile(filepath.Join(tmpDir, name), []byte("GET /users\nHTTP 200\n"), 0o644)

	result := FindHurlFile(ep, "nonexistent_dir")
	if result == "" {
		t.Fatal("expected to find hurl file in current dir")
	}
}

func TestFindHurlFile2_NotFound(t *testing.T) {
	ep := &scanner.Endpoint{Method: "GET", Path: "/users"}
	result := FindHurlFile(ep, "/nonexistent")
	if result != "" {
		t.Fatalf("expected empty, got %s", result)
	}
}
