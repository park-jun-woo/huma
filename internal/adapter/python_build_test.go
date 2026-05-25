package adapter

import "testing"

func TestPythonBuild_NoOp(t *testing.T) {
	a := &PythonAdapter{}
	err := a.Build()
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}
