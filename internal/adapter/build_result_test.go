package adapter

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/park-jun-woo/huma/internal/coverage"
)

func TestBuildResult_FullCoverage(t *testing.T) {
	tmpDir := t.TempDir()
	f := filepath.Join(tmpDir, "handler.go")
	os.WriteFile(f, []byte("line1\nline2\nline3\n"), 0o644)

	blocks := []coverage.Block{
		{File: f, StartLine: 1, EndLine: 3, Count: 1},
	}

	result, err := buildResult(blocks, f, 1, 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Percent != 100 {
		t.Fatalf("expected 100%%, got %f", result.Percent)
	}
	if len(result.Uncovered) != 0 {
		t.Fatalf("expected 0 uncovered, got %d", len(result.Uncovered))
	}
}

func TestBuildResult_PreservesCoveredLines(t *testing.T) {
	tmpDir := t.TempDir()
	f := filepath.Join(tmpDir, "handler.go")
	os.WriteFile(f, []byte("l1\nl2\nl3\nl4\n"), 0o644)

	blocks := []coverage.Block{
		{File: f, StartLine: 1, EndLine: 2, Count: 1},
		{File: f, StartLine: 3, EndLine: 4, Count: 0},
	}
	result, err := buildResult(blocks, f, 1, 4)
	if err != nil {
		t.Fatal(err)
	}
	if !result.IsLineCovered(1) || !result.IsLineCovered(2) {
		t.Fatal("expected lines 1,2 in covered set")
	}
	if result.IsLineCovered(3) {
		t.Fatal("expected line 3 not covered")
	}
}

func TestBuildResult_PartialCoverage(t *testing.T) {
	tmpDir := t.TempDir()
	f := filepath.Join(tmpDir, "handler.go")
	os.WriteFile(f, []byte("line1\nline2\nline3\nline4\n"), 0o644)

	blocks := []coverage.Block{
		{File: f, StartLine: 1, EndLine: 2, Count: 1},
		{File: f, StartLine: 3, EndLine: 4, Count: 0},
	}

	result, err := buildResult(blocks, f, 1, 4)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Covered >= result.Total {
		t.Fatalf("expected partial coverage, got covered=%d total=%d", result.Covered, result.Total)
	}
	if len(result.Uncovered) == 0 {
		t.Fatal("expected uncovered lines")
	}
}

func TestBuildResult_EmptyBlocks(t *testing.T) {
	tmpDir := t.TempDir()
	f := filepath.Join(tmpDir, "handler.go")
	os.WriteFile(f, []byte("line1\n"), 0o644)

	result, err := buildResult(nil, f, 1, 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Total != 0 {
		t.Fatalf("expected 0 total, got %d", result.Total)
	}
	if result.Percent != 0 {
		t.Fatalf("expected 0%%, got %f", result.Percent)
	}
}

func TestBuildResult_ReadUncoveredError(t *testing.T) {
	// Non-existent file causes readUncoveredLines to fail.
	// buildResult should still return a result with nil uncovered lines.
	blocks := []coverage.Block{
		{File: "/nonexistent/handler.go", StartLine: 1, EndLine: 3, Count: 0},
	}

	result, err := buildResult(blocks, "/nonexistent/handler.go", 1, 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Uncovered != nil {
		t.Fatalf("expected nil uncovered on error, got %v", result.Uncovered)
	}
	if result.Total == 0 {
		t.Fatal("expected non-zero total from blocks")
	}
}
