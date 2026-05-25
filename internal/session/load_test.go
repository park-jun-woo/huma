package session

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/park-jun-woo/huma/internal/scanner"
)

func TestLoad2_Success(t *testing.T) {
	tmpDir := t.TempDir()
	orig, _ := os.Getwd()
	t.Cleanup(func() { os.Chdir(orig) })
	os.Chdir(tmpDir)

	sess := New()
	sess.Merge([]scanner.Endpoint{{ID: "ep1", Method: "GET", Path: "/a"}})
	sess.Save()

	loaded, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(loaded.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(loaded.Entries))
	}
}

func TestLoad2_FileNotFound(t *testing.T) {
	tmpDir := t.TempDir()
	orig, _ := os.Getwd()
	t.Cleanup(func() { os.Chdir(orig) })
	os.Chdir(tmpDir)

	_, err := Load()
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestLoad2_InvalidJSON(t *testing.T) {
	tmpDir := t.TempDir()
	orig, _ := os.Getwd()
	t.Cleanup(func() { os.Chdir(orig) })
	os.Chdir(tmpDir)

	os.MkdirAll(filepath.Join(tmpDir, ".huma"), 0o755)
	os.WriteFile(filepath.Join(tmpDir, ".huma", "session.json"), []byte("INVALID"), 0o644)

	_, err := Load()
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}
