package adapter

import "testing"

func TestDenoCollect_ReturnsNil(t *testing.T) {
	a := &DenoAdapter{}
	result, err := a.Collect("handler.ts", 1, 10)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if result != nil {
		t.Fatal("expected nil result")
	}
}
