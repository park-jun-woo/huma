package session

import (
	"testing"

	"github.com/park-jun-woo/huma/internal/scanner"
)

func TestMarkPass2_SetsStatus(t *testing.T) {
	sess := New()
	sess.Merge([]scanner.Endpoint{{ID: "ep1", Method: "GET", Path: "/a"}})
	sess.MarkPass("ep1")

	for _, e := range sess.Entries {
		if e.ID == "ep1" {
			if e.Status != StatusPass {
				t.Fatalf("expected PASS, got %s", e.Status)
			}
			return
		}
	}
	t.Fatal("entry not found")
}

func TestMarkPass2_NotFound(t *testing.T) {
	sess := New()
	sess.Merge([]scanner.Endpoint{{ID: "ep1", Method: "GET", Path: "/a"}})
	sess.MarkPass("nonexistent")
	if sess.Entries[0].Status != StatusTodo {
		t.Fatal("expected unchanged")
	}
}
