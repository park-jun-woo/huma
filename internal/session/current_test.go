package session

import (
	"testing"

	"github.com/park-jun-woo/huma/internal/scanner"
)

func TestCurrent2_ReturnsTodo(t *testing.T) {
	sess := New()
	sess.Merge([]scanner.Endpoint{
		{ID: "ep1", Method: "GET", Path: "/a"},
	})
	ep := sess.Current()
	if ep == nil {
		t.Fatal("expected non-nil")
	}
	if ep.ID != "ep1" {
		t.Fatalf("expected ep1, got %s", ep.ID)
	}
}

func TestCurrent_SkipsPass(t *testing.T) {
	sess := New()
	sess.Merge([]scanner.Endpoint{
		{ID: "ep1", Method: "GET", Path: "/a"},
		{ID: "ep2", Method: "POST", Path: "/b"},
	})
	sess.MarkPass("ep1")
	ep := sess.Current()
	if ep == nil || ep.ID != "ep2" {
		t.Fatal("expected ep2")
	}
}

func TestCurrent2_ReturnsImprove(t *testing.T) {
	sess := New()
	sess.Merge([]scanner.Endpoint{
		{ID: "ep1", Method: "GET", Path: "/a"},
	})
	sess.MarkImprove("ep1", 50)
	ep := sess.Current()
	if ep == nil || ep.ID != "ep1" {
		t.Fatal("expected ep1 with IMPROVE status")
	}
}

func TestCurrent_AllComplete(t *testing.T) {
	sess := New()
	sess.Merge([]scanner.Endpoint{
		{ID: "ep1", Method: "GET", Path: "/a"},
	})
	sess.MarkPass("ep1")
	ep := sess.Current()
	if ep != nil {
		t.Fatal("expected nil when all complete")
	}
}
