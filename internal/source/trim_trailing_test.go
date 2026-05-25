package source

import "testing"

func TestTrimTrailing_RemovesBlanks(t *testing.T) {
	lines := []string{"code", "", "  ", ""}
	result := trimTrailing(lines)
	if len(result) != 1 || result[0] != "code" {
		t.Fatalf("expected [code], got %v", result)
	}
}

func TestTrimTrailing_RemovesComments(t *testing.T) {
	lines := []string{"code", "// comment", ""}
	result := trimTrailing(lines)
	if len(result) != 1 || result[0] != "code" {
		t.Fatalf("expected [code], got %v", result)
	}
}

func TestTrimTrailing_NoTrimNeeded(t *testing.T) {
	lines := []string{"line1", "line2"}
	result := trimTrailing(lines)
	if len(result) != 2 {
		t.Fatalf("expected 2, got %d", len(result))
	}
}

func TestTrimTrailing_Empty(t *testing.T) {
	result := trimTrailing(nil)
	if len(result) != 0 {
		t.Fatalf("expected 0, got %d", len(result))
	}
}

func TestTrimTrailing_AllBlanks(t *testing.T) {
	lines := []string{"", "  ", "// comment"}
	result := trimTrailing(lines)
	if len(result) != 0 {
		t.Fatalf("expected 0, got %d", len(result))
	}
}
