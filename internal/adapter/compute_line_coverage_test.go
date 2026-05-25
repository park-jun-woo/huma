package adapter

import (
	"testing"

	"github.com/park-jun-woo/huma/internal/coverage"
)

func TestComputeLineCoverage_Empty(t *testing.T) {
	covered, total := computeLineCoverage(nil, 1, 10)
	if len(covered) != 0 {
		t.Fatalf("expected 0 covered, got %d", len(covered))
	}
	if len(total) != 0 {
		t.Fatalf("expected 0 total, got %d", len(total))
	}
}

func TestComputeLineCoverage_SingleBlock(t *testing.T) {
	blocks := []coverage.Block{
		{File: "a.go", StartLine: 5, EndLine: 8, Count: 1},
	}
	covered, total := computeLineCoverage(blocks, 1, 10)
	if len(total) == 0 {
		t.Fatal("expected non-empty total")
	}
	if len(covered) == 0 {
		t.Fatal("expected non-empty covered")
	}
}

func TestComputeLineCoverage_MultipleBlocks(t *testing.T) {
	blocks := []coverage.Block{
		{File: "a.go", StartLine: 2, EndLine: 4, Count: 1},
		{File: "a.go", StartLine: 5, EndLine: 7, Count: 0},
	}
	covered, total := computeLineCoverage(blocks, 1, 10)
	// All lines from both blocks should be in total
	if len(total) == 0 {
		t.Fatal("expected non-empty total")
	}
	// Only lines from the first block (Count>0) should be covered
	for line := range covered {
		if line < 2 || line > 4 {
			t.Fatalf("unexpected covered line: %d", line)
		}
	}
}
