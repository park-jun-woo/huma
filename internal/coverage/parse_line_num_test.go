package coverage

import "testing"

func TestParseLineNum2_WithDot(t *testing.T) {
	n, err := parseLineNum("42.5")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 42 {
		t.Fatalf("expected 42, got %d", n)
	}
}

func TestParseLineNum2_NoDot(t *testing.T) {
	n, err := parseLineNum("42")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 42 {
		t.Fatalf("expected 42, got %d", n)
	}
}

func TestParseLineNum2_Invalid(t *testing.T) {
	_, err := parseLineNum("abc")
	if err == nil {
		t.Fatal("expected error for invalid input")
	}
}
