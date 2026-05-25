package coverage

import "testing"

func TestFilterBlocks_MatchesFileAndRange(t *testing.T) {
	blocks := []Block{
		{File: "pkg/handler.go", StartLine: 5, EndLine: 10, Count: 1},
		{File: "pkg/handler.go", StartLine: 15, EndLine: 20, Count: 1},
		{File: "pkg/other.go", StartLine: 5, EndLine: 10, Count: 1},
	}

	filtered := FilterBlocks(blocks, "handler.go", 1, 12)
	if len(filtered) != 1 {
		t.Fatalf("expected 1, got %d", len(filtered))
	}
	if filtered[0].StartLine != 5 {
		t.Fatalf("expected start 5, got %d", filtered[0].StartLine)
	}
}

func TestFilterBlocks_NoMatch(t *testing.T) {
	blocks := []Block{
		{File: "pkg/handler.go", StartLine: 20, EndLine: 30, Count: 1},
	}

	filtered := FilterBlocks(blocks, "handler.go", 1, 10)
	if len(filtered) != 0 {
		t.Fatalf("expected 0, got %d", len(filtered))
	}
}

func TestFilterBlocks_EmptyBlocks(t *testing.T) {
	filtered := FilterBlocks(nil, "handler.go", 1, 10)
	if len(filtered) != 0 {
		t.Fatalf("expected 0, got %d", len(filtered))
	}
}

func TestFilterBlocks_PartialOverlap(t *testing.T) {
	blocks := []Block{
		{File: "handler.go", StartLine: 8, EndLine: 15, Count: 1},
	}
	filtered := FilterBlocks(blocks, "handler.go", 10, 20)
	if len(filtered) != 1 {
		t.Fatalf("expected 1, got %d", len(filtered))
	}
}
