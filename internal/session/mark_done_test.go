package session

import (
	"testing"

	"github.com/park-jun-woo/huma/internal/scanner"
)

func TestMarkDone_SetsCoverageAndStatus(t *testing.T) {
	sess := New()
	sess.Merge([]scanner.Endpoint{{ID: "ep1", Method: "GET", Path: "/a"}})
	sess.MarkDone("ep1", 85.5)

	for _, e := range sess.Entries {
		if e.ID == "ep1" {
			if e.Status != StatusDone {
				t.Fatalf("expected DONE, got %s", e.Status)
			}
			if e.Coverage != 85.5 {
				t.Fatalf("expected 85.5, got %f", e.Coverage)
			}
			return
		}
	}
	t.Fatal("entry not found")
}

func TestMarkDone_NotFound(t *testing.T) {
	sess := New()
	sess.Merge([]scanner.Endpoint{{ID: "ep1", Method: "GET", Path: "/a"}})
	sess.MarkDone("nonexistent", 50)
	// Should not panic, just no-op
	if sess.Entries[0].Status != StatusTodo {
		t.Fatalf("expected entry unchanged, got %s", sess.Entries[0].Status)
	}
}
