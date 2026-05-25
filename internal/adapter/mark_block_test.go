package adapter

import (
	"testing"

	"github.com/park-jun-woo/huma/internal/coverage"
)

func TestMarkBlock_FullyInRange(t *testing.T) {
	covered := make(map[int]bool)
	total := make(map[int]bool)
	b := coverage.Block{StartLine: 5, EndLine: 8, Count: 1}

	markBlock(b, 1, 10, covered, total)

	for line := 5; line <= 8; line++ {
		if !total[line] {
			t.Fatalf("line %d should be in total", line)
		}
		if !covered[line] {
			t.Fatalf("line %d should be covered", line)
		}
	}
}

func TestMarkBlock_NotCovered(t *testing.T) {
	covered := make(map[int]bool)
	total := make(map[int]bool)
	b := coverage.Block{StartLine: 5, EndLine: 8, Count: 0}

	markBlock(b, 1, 10, covered, total)

	for line := 5; line <= 8; line++ {
		if !total[line] {
			t.Fatalf("line %d should be in total", line)
		}
		if covered[line] {
			t.Fatalf("line %d should NOT be covered", line)
		}
	}
}

func TestMarkBlock_ClampedToRange(t *testing.T) {
	covered := make(map[int]bool)
	total := make(map[int]bool)
	b := coverage.Block{StartLine: 3, EndLine: 12, Count: 1}

	markBlock(b, 5, 10, covered, total)

	// Lines outside [5,10] should not be marked
	if total[3] || total[4] || total[11] || total[12] {
		t.Fatal("lines outside range should not be in total")
	}
	for line := 5; line <= 10; line++ {
		if !total[line] {
			t.Fatalf("line %d should be in total", line)
		}
		if !covered[line] {
			t.Fatalf("line %d should be covered", line)
		}
	}
}

func TestMarkBlock_StartBeforeRange(t *testing.T) {
	covered := make(map[int]bool)
	total := make(map[int]bool)
	b := coverage.Block{StartLine: 1, EndLine: 5, Count: 1}

	markBlock(b, 3, 10, covered, total)

	if total[1] || total[2] {
		t.Fatal("lines before startLine should not be in total")
	}
	for line := 3; line <= 5; line++ {
		if !total[line] {
			t.Fatalf("line %d should be in total", line)
		}
	}
}

func TestMarkBlock_EndAfterRange(t *testing.T) {
	covered := make(map[int]bool)
	total := make(map[int]bool)
	b := coverage.Block{StartLine: 8, EndLine: 15, Count: 1}

	markBlock(b, 5, 10, covered, total)

	if total[11] || total[15] {
		t.Fatal("lines after endLine should not be in total")
	}
	for line := 8; line <= 10; line++ {
		if !total[line] {
			t.Fatalf("line %d should be in total", line)
		}
	}
}

func TestMarkBlock_OutOfRange(t *testing.T) {
	covered := make(map[int]bool)
	total := make(map[int]bool)
	b := coverage.Block{StartLine: 20, EndLine: 30, Count: 1}

	markBlock(b, 5, 10, covered, total)

	if len(total) != 0 {
		t.Fatal("no lines should be marked for out-of-range block")
	}
}
