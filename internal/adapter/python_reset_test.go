package adapter

import (
	"os"
	"path/filepath"
	"testing"
)

func TestPythonReset_FileExists(t *testing.T) {
	tmpDir := t.TempDir()
	orig, _ := os.Getwd()
	t.Cleanup(func() { os.Chdir(orig) })
	os.Chdir(tmpDir)

	os.WriteFile(filepath.Join(tmpDir, ".coverage"), []byte("data"), 0o644)

	a := &PythonAdapter{}
	err := a.Reset()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, err := os.Stat(filepath.Join(tmpDir, ".coverage")); !os.IsNotExist(err) {
		t.Fatal("expected .coverage to be removed")
	}
}

func TestPythonReset_FileNotExists(t *testing.T) {
	tmpDir := t.TempDir()
	orig, _ := os.Getwd()
	t.Cleanup(func() { os.Chdir(orig) })
	os.Chdir(tmpDir)

	a := &PythonAdapter{}
	err := a.Reset()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestPythonReset_PermissionError(t *testing.T) {
	tmpDir := t.TempDir()
	orig, _ := os.Getwd()
	t.Cleanup(func() { os.Chdir(orig) })
	os.Chdir(tmpDir)

	os.WriteFile(filepath.Join(tmpDir, ".coverage"), []byte("data"), 0o644)
	// Make directory read-only so remove fails
	os.Chmod(tmpDir, 0o555)
	t.Cleanup(func() { os.Chmod(tmpDir, 0o755) })

	a := &PythonAdapter{}
	err := a.Reset()
	if err == nil {
		t.Fatal("expected permission error")
	}
}
