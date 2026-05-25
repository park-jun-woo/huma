package session

import (
	"testing"

	"github.com/park-jun-woo/huma/internal/scanner"
)

func TestMarkImprove_SetsStatusAndCoverage(t *testing.T) {
	sess := New()
	sess.Merge([]scanner.Endpoint{{ID: "ep1", Method: "GET", Path: "/a"}})
	sess.MarkImprove("ep1", 70)

	for _, e := range sess.Entries {
		if e.ID == "ep1" {
			if e.Status != StatusImprove {
				t.Fatalf("expected IMPROVE, got %s", e.Status)
			}
			if e.Coverage != 70 {
				t.Fatalf("expected 70, got %f", e.Coverage)
			}
			if e.ImproveCount != 1 {
				t.Fatalf("expected 1, got %d", e.ImproveCount)
			}
			return
		}
	}
	t.Fatal("entry not found")
}

func TestMarkImprove_TracksPrevCoverage(t *testing.T) {
	sess := New()
	sess.Merge([]scanner.Endpoint{{ID: "ep1", Method: "GET", Path: "/a"}})
	sess.MarkImprove("ep1", 50)
	sess.MarkImprove("ep1", 70)

	for _, e := range sess.Entries {
		if e.ID == "ep1" {
			if e.PrevCoverage != 50 {
				t.Fatalf("expected PrevCoverage 50, got %f", e.PrevCoverage)
			}
			if e.Coverage != 70 {
				t.Fatalf("expected Coverage 70, got %f", e.Coverage)
			}
			if e.ImproveCount != 2 {
				t.Fatalf("expected 2, got %d", e.ImproveCount)
			}
			return
		}
	}
	t.Fatal("entry not found")
}

func TestMarkImprove_NotFound(t *testing.T) {
	sess := New()
	sess.Merge([]scanner.Endpoint{{ID: "ep1", Method: "GET", Path: "/a"}})
	sess.MarkImprove("nonexistent", 50)
	if sess.Entries[0].Status != StatusTodo {
		t.Fatal("expected unchanged")
	}
}
