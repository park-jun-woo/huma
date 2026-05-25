package adapter

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNodeReset_Success(t *testing.T) {
	tmpDir := t.TempDir()
	covDir := filepath.Join(tmpDir, "v8cov")
	os.MkdirAll(covDir, 0o755)
	os.WriteFile(filepath.Join(covDir, "old.json"), []byte("data"), 0o644)

	a := &NodeAdapter{coverDir: covDir}
	err := a.Reset()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	entries, err := os.ReadDir(covDir)
	if err != nil {
		t.Fatalf("read dir: %v", err)
	}
	if len(entries) != 0 {
		t.Fatalf("expected empty dir, got %d entries", len(entries))
	}
}

func TestNodeReset_CreateDir(t *testing.T) {
	tmpDir := t.TempDir()
	covDir := filepath.Join(tmpDir, "new", "v8cov")

	a := &NodeAdapter{coverDir: covDir}
	err := a.Reset()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	info, err := os.Stat(covDir)
	if err != nil {
		t.Fatalf("stat: %v", err)
	}
	if !info.IsDir() {
		t.Fatal("expected directory")
	}
}

func TestNodeReset_RemoveAllError(t *testing.T) {
	tmpDir := t.TempDir()
	parentDir := filepath.Join(tmpDir, "parent")
	covDir := filepath.Join(parentDir, "v8cov")
	os.MkdirAll(covDir, 0o755)
	os.WriteFile(filepath.Join(covDir, "file.dat"), []byte("data"), 0o644)
	os.Chmod(parentDir, 0o555)
	t.Cleanup(func() { os.Chmod(parentDir, 0o755) })

	a := &NodeAdapter{coverDir: covDir}
	err := a.Reset()
	if err == nil {
		t.Fatal("expected error from RemoveAll")
	}
}

func TestNodeReset_MkdirAllError(t *testing.T) {
	tmpDir := t.TempDir()
	parentDir := filepath.Join(tmpDir, "parent")
	covDir := filepath.Join(parentDir, "v8cov")
	os.MkdirAll(covDir, 0o755)

	// After RemoveAll(covDir) succeeds, make parent read-only so MkdirAll(covDir) fails
	// We need to remove covDir first, then make parent read-only
	os.RemoveAll(covDir) // remove so RemoveAll in Reset is a no-op
	os.Chmod(parentDir, 0o555)
	t.Cleanup(func() { os.Chmod(parentDir, 0o755) })

	a := &NodeAdapter{coverDir: covDir}
	err := a.Reset()
	if err == nil {
		t.Fatal("expected error from MkdirAll")
	}
}
