package session

import (
	"testing"

	"github.com/park-jun-woo/huma/internal/scanner"
)

func TestCurrentEntry_ReturnsTodo(t *testing.T) {
	sess := New()
	sess.Merge([]scanner.Endpoint{{ID: "ep1", Method: "GET", Path: "/a"}})
	entry := sess.CurrentEntry()
	if entry == nil {
		t.Fatal("expected non-nil")
	}
	if entry.Status != StatusTodo {
		t.Fatalf("expected TODO, got %s", entry.Status)
	}
}

func TestCurrentEntry_ReturnsImprove(t *testing.T) {
	sess := New()
	sess.Merge([]scanner.Endpoint{{ID: "ep1", Method: "GET", Path: "/a"}})
	sess.MarkImprove("ep1", 50)
	entry := sess.CurrentEntry()
	if entry == nil || entry.Status != StatusImprove {
		t.Fatal("expected IMPROVE entry")
	}
}

func TestCurrentEntry_AllComplete(t *testing.T) {
	sess := New()
	sess.Merge([]scanner.Endpoint{{ID: "ep1", Method: "GET", Path: "/a"}})
	sess.MarkPass("ep1")
	entry := sess.CurrentEntry()
	if entry != nil {
		t.Fatal("expected nil when all complete")
	}
}
