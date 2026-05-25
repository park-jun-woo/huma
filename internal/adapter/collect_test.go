package adapter

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCollect_CovdataError2(t *testing.T) {
	a := &GoAdapter{
		coverDir: "/nonexistent/path/to/coverdata",
	}
	_, err := a.Collect("handler.go", 10, 20)
	if err == nil {
		t.Fatal("expected error from covdata textfmt")
	}
}

func TestCollect_EmptyCoverage2(t *testing.T) {
	tmpDir := t.TempDir()
	orig, _ := os.Getwd()
	t.Cleanup(func() { os.Chdir(orig) })
	os.Chdir(tmpDir)

	covDir := filepath.Join(tmpDir, ".huma", "coverdata")
	os.MkdirAll(covDir, 0o755)

	a := &GoAdapter{coverDir: covDir}
	result, err := a.Collect("handler.go", 10, 20)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Total != 0 {
		t.Fatalf("expected 0 total, got %d", result.Total)
	}
}

func TestCollect_ParseCoverageFileError(t *testing.T) {
	tmpDir := t.TempDir()
	orig, _ := os.Getwd()
	t.Cleanup(func() { os.Chdir(orig) })
	os.Chdir(tmpDir)

	covDir := filepath.Join(tmpDir, ".huma", "coverdata")
	os.MkdirAll(covDir, 0o755)

	// Create .huma/coverage.out as a directory so covdata textfmt writes
	// the mode line but ParseCoverageFile cannot open it as a regular file.
	// Actually, covdata -o=dir will fail. Instead, create it as a directory
	// so that covdata textfmt fails to write to it.
	// The trick: make the .huma directory writable but make coverage.out
	// a directory, which covdata -o cannot overwrite.
	coverOutPath := filepath.Join(tmpDir, ".huma", "coverage.out")
	os.MkdirAll(coverOutPath, 0o755) // create coverage.out as a directory

	a := &GoAdapter{coverDir: covDir}
	_, err := a.Collect("handler.go", 1, 5)
	// covdata will fail trying to write to a directory, so we expect an error.
	// The error may come from covdata or from ParseCoverageFile.
	if err == nil {
		t.Fatal("expected error when coverage.out is a directory")
	}
}
