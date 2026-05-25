package adapter

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGoReset_Success(t *testing.T) {
	tmpDir := t.TempDir()
	covDir := filepath.Join(tmpDir, "coverdata")
	os.MkdirAll(covDir, 0o755)
	os.WriteFile(filepath.Join(covDir, "old.dat"), []byte("old"), 0o644)

	a := &GoAdapter{coverDir: covDir}
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

func TestGoReset_RemoveAllError(t *testing.T) {
	tmpDir := t.TempDir()
	parentDir := filepath.Join(tmpDir, "parent")
	covDir := filepath.Join(parentDir, "coverdata")
	os.MkdirAll(covDir, 0o755)
	os.WriteFile(filepath.Join(covDir, "file.dat"), []byte("data"), 0o644)
	os.Chmod(parentDir, 0o555)
	t.Cleanup(func() { os.Chmod(parentDir, 0o755) })

	a := &GoAdapter{coverDir: covDir}
	err := a.Reset()
	if err == nil {
		t.Fatal("expected error from RemoveAll")
	}
}

func TestGoReset_MkdirAllError(t *testing.T) {
	tmpDir := t.TempDir()
	parentDir := filepath.Join(tmpDir, "parent")
	covDir := filepath.Join(parentDir, "coverdata")
	os.MkdirAll(covDir, 0o755)
	os.RemoveAll(covDir)
	os.Chmod(parentDir, 0o555)
	t.Cleanup(func() { os.Chmod(parentDir, 0o755) })

	a := &GoAdapter{coverDir: covDir}
	err := a.Reset()
	if err == nil {
		t.Fatal("expected error from MkdirAll")
	}
}

func TestGoReset_CreateDir(t *testing.T) {
	tmpDir := t.TempDir()
	covDir := filepath.Join(tmpDir, "new", "coverdata")

	a := &GoAdapter{coverDir: covDir}
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
