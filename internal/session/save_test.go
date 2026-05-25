package session

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/park-jun-woo/huma/internal/scanner"
)

func TestSave2_Success(t *testing.T) {
	tmpDir := t.TempDir()
	orig, _ := os.Getwd()
	t.Cleanup(func() { os.Chdir(orig) })
	os.Chdir(tmpDir)

	sess := New()
	sess.Merge([]scanner.Endpoint{{ID: "ep1", Method: "GET", Path: "/a"}})
	err := sess.Save()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, err := os.ReadFile(filepath.Join(tmpDir, ".huma", "session.json"))
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	if len(data) == 0 {
		t.Fatal("expected non-empty file")
	}
}

func TestSave2_WriteError(t *testing.T) {
	tmpDir := t.TempDir()
	orig, _ := os.Getwd()
	t.Cleanup(func() { os.Chdir(orig) })
	os.Chdir(tmpDir)

	// Create .huma dir and make it read-only so WriteFile fails
	dir := filepath.Join(tmpDir, ".huma")
	os.MkdirAll(dir, 0o755)
	os.Chmod(dir, 0o555)
	t.Cleanup(func() { os.Chmod(dir, 0o755) })

	sess := New()
	err := sess.Save()
	if err == nil {
		t.Fatal("expected write error")
	}
}

func TestSave2_MkdirError(t *testing.T) {
	tmpDir := t.TempDir()
	orig, _ := os.Getwd()
	t.Cleanup(func() { os.Chdir(orig) })
	os.Chdir(tmpDir)

	os.WriteFile(filepath.Join(tmpDir, ".huma"), []byte("block"), 0o644)

	sess := New()
	err := sess.Save()
	if err == nil {
		t.Fatal("expected error from MkdirAll")
	}
}
