package adapter

import "testing"

func TestDenoReset_NoOp(t *testing.T) {
	a := &DenoAdapter{}
	err := a.Reset()
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}
