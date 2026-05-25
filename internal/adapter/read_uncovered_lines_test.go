package adapter

import (
	"os"
	"path/filepath"
	"testing"
)

func TestReadUncoveredLines_Basic(t *testing.T) {
	tmpDir := t.TempDir()
	f := filepath.Join(tmpDir, "handler.go")
	os.WriteFile(f, []byte("line1\nline2\nline3\nline4\nline5\n"), 0o644)

	covered := map[int]bool{1: true, 3: true}
	total := map[int]bool{1: true, 2: true, 3: true, 4: true}

	result, err := readUncoveredLines(f, covered, total)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 2 {
		t.Fatalf("expected 2 uncovered lines, got %d", len(result))
	}
	if result[0].Line != 2 || result[0].Code != "line2" {
		t.Fatalf("unexpected first uncovered line: %+v", result[0])
	}
	if result[1].Line != 4 || result[1].Code != "line4" {
		t.Fatalf("unexpected second uncovered line: %+v", result[1])
	}
}

func TestReadUncoveredLines_FileNotFound2(t *testing.T) {
	_, err := readUncoveredLines("/nonexistent/file.go", nil, nil)
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestReadUncoveredLines_AllCovered2(t *testing.T) {
	tmpDir := t.TempDir()
	f := filepath.Join(tmpDir, "handler.go")
	os.WriteFile(f, []byte("line1\nline2\n"), 0o644)

	covered := map[int]bool{1: true, 2: true}
	total := map[int]bool{1: true, 2: true}

	result, err := readUncoveredLines(f, covered, total)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 0 {
		t.Fatalf("expected 0 uncovered lines, got %d", len(result))
	}
}

func TestReadUncoveredLines_Trimming(t *testing.T) {
	tmpDir := t.TempDir()
	f := filepath.Join(tmpDir, "handler.go")
	os.WriteFile(f, []byte("code  \t\n"), 0o644)

	covered := map[int]bool{}
	total := map[int]bool{1: true}

	result, err := readUncoveredLines(f, covered, total)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 1 || result[0].Code != "code" {
		t.Fatalf("expected trimmed code 'code', got %q", result[0].Code)
	}
}
