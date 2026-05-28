package adapter

import "testing"

func TestRustCollect_ReturnsNil(t *testing.T) {
	a := &RustAdapter{}
	result, err := a.Collect("handler.rs", 1, 10)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if result != nil {
		t.Fatal("expected nil result")
	}
}
