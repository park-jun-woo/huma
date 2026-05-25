package coverage

import "testing"

func TestFilterMissingLines_InRange(t *testing.T) {
	result := filterMissingLines([]int{1, 5, 10, 15, 20}, 5, 15)
	if len(result) != 3 {
		t.Fatalf("expected 3, got %d", len(result))
	}
}

func TestFilterMissingLines_NoneInRange(t *testing.T) {
	result := filterMissingLines([]int{1, 2, 3}, 10, 20)
	if len(result) != 0 {
		t.Fatalf("expected 0, got %d", len(result))
	}
}

func TestFilterMissingLines_Empty(t *testing.T) {
	result := filterMissingLines(nil, 1, 10)
	if len(result) != 0 {
		t.Fatalf("expected 0, got %d", len(result))
	}
}

func TestFilterMissingLines_BoundaryValues(t *testing.T) {
	result := filterMissingLines([]int{5, 10}, 5, 10)
	if len(result) != 2 {
		t.Fatalf("expected 2, got %d", len(result))
	}
}
