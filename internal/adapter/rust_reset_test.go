package adapter

import "testing"

func TestRustReset_NoOp(t *testing.T) {
	a := &RustAdapter{}
	err := a.Reset()
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}
