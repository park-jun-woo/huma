package scanner

import "testing"

func TestMakeID_Deterministic(t *testing.T) {
	id1 := makeID("GET", "/users")
	id2 := makeID("GET", "/users")
	if id1 != id2 {
		t.Fatal("expected same ID for same input")
	}
}

func TestMakeID_DifferentMethods(t *testing.T) {
	id1 := makeID("GET", "/users")
	id2 := makeID("POST", "/users")
	if id1 == id2 {
		t.Fatal("expected different IDs for different methods")
	}
}

func TestMakeID_DifferentPaths(t *testing.T) {
	id1 := makeID("GET", "/users")
	id2 := makeID("GET", "/posts")
	if id1 == id2 {
		t.Fatal("expected different IDs for different paths")
	}
}

func TestMakeID_Length(t *testing.T) {
	id := makeID("GET", "/users")
	if len(id) != 16 {
		t.Fatalf("expected 16 hex chars, got %d", len(id))
	}
}
