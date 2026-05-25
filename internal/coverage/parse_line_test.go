package coverage

import "testing"

func TestParseLine2_Valid(t *testing.T) {
	b, err := parseLine("github.com/user/repo/handler.go:10.2,12.4 1 5")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if b.Count != 5 {
		t.Fatalf("expected count 5, got %d", b.Count)
	}
}

func TestParseLine2_NoSpace(t *testing.T) {
	_, err := parseLine("nospaces")
	if err == nil {
		t.Fatal("expected error for no space")
	}
}

func TestParseLine2_NoSecondSpace(t *testing.T) {
	_, err := parseLine("onlyonespace 5")
	if err == nil {
		t.Fatal("expected error for no second space")
	}
}

func TestParseLine2_InvalidCount(t *testing.T) {
	_, err := parseLine("file:1.0,2.0 1 notanumber")
	if err == nil {
		t.Fatal("expected error for invalid count")
	}
}

func TestParseLine2_ZeroCount(t *testing.T) {
	b, err := parseLine("github.com/user/repo/handler.go:10.2,12.4 1 0")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if b.Count != 0 {
		t.Fatalf("expected count 0, got %d", b.Count)
	}
}
