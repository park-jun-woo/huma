package prompt

import "testing"

func TestHasParam2_True(t *testing.T) {
	if !hasParam("/users/:id") {
		t.Fatal("expected true")
	}
}

func TestHasParam2_False(t *testing.T) {
	if hasParam("/users") {
		t.Fatal("expected false")
	}
}
