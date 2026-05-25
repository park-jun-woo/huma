package coverage

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseCoverageFile_Valid(t *testing.T) {
	tmpDir := t.TempDir()
	f := filepath.Join(tmpDir, "coverage.out")
	content := `mode: set
github.com/user/repo/handler.go:10.30,12.16 2 1
github.com/user/repo/handler.go:12.16,14.3 1 0
`
	os.WriteFile(f, []byte(content), 0o644)

	blocks, err := ParseCoverageFile(f)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(blocks) != 2 {
		t.Fatalf("expected 2 blocks, got %d", len(blocks))
	}
}

func TestParseCoverageFile_MissingFile(t *testing.T) {
	_, err := ParseCoverageFile("/nonexistent/coverage.out")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestParseCoverageFile_EmptyFile(t *testing.T) {
	tmpDir := t.TempDir()
	f := filepath.Join(tmpDir, "coverage.out")
	os.WriteFile(f, []byte("mode: set\n"), 0o644)

	blocks, err := ParseCoverageFile(f)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(blocks) != 0 {
		t.Fatalf("expected 0 blocks, got %d", len(blocks))
	}
}

func TestParseCoverageFile_LongLineScanError(t *testing.T) {
	tmpDir := t.TempDir()
	f := filepath.Join(tmpDir, "coverage.out")
	// Create a line longer than bufio.Scanner's max token size (64KB)
	longLine := make([]byte, 1024*1024)
	for i := range longLine {
		longLine[i] = 'x'
	}
	content := append([]byte("mode: set\n"), longLine...)
	os.WriteFile(f, content, 0o644)

	_, err := ParseCoverageFile(f)
	if err == nil {
		t.Fatal("expected scanner error for long line")
	}
}

func TestParseCoverageFile_InvalidLines(t *testing.T) {
	tmpDir := t.TempDir()
	f := filepath.Join(tmpDir, "coverage.out")
	content := "mode: set\ninvalid line\n"
	os.WriteFile(f, []byte(content), 0o644)

	blocks, err := ParseCoverageFile(f)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Invalid lines are skipped
	if len(blocks) != 0 {
		t.Fatalf("expected 0 blocks, got %d", len(blocks))
	}
}
