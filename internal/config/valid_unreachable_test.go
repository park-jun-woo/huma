package config

import "testing"

func TestValidUnreachable(t *testing.T) {
	raw := []UnreachableEntry{
		{Endpoint: "GET /a", Status: 404, Reason: "r", Evidence: "e"},   // valid
		{Endpoint: "GET /b", Status: 500, Reason: "", Evidence: "e"},    // no reason
		{Endpoint: "GET /c", Status: 503, Reason: "r", Evidence: ""},    // no evidence
		{Endpoint: "GET /d", Status: 400, Reason: "r2", Evidence: "e2"}, // valid
	}
	got := validUnreachable(raw)
	if len(got) != 2 {
		t.Fatalf("expected 2 valid entries, got %d: %+v", len(got), got)
	}
	if got[0].Endpoint != "GET /a" || got[1].Endpoint != "GET /d" {
		t.Errorf("unexpected kept entries: %+v", got)
	}
	// all invalid → empty
	if g := validUnreachable([]UnreachableEntry{{Reason: "r"}}); len(g) != 0 {
		t.Errorf("invalid only → empty, got %+v", g)
	}
}
