package session

import (
	"testing"

	"github.com/park-jun-woo/huma/internal/scanner"
)

func TestStats2_AllTodo(t *testing.T) {
	sess := New()
	sess.Merge([]scanner.Endpoint{
		{ID: "ep1", Method: "GET", Path: "/a"},
		{ID: "ep2", Method: "POST", Path: "/b"},
	})
	total, pass, todo := sess.Stats()
	if total != 2 || pass != 0 || todo != 2 {
		t.Fatalf("unexpected: total=%d pass=%d todo=%d", total, pass, todo)
	}
}

func TestStats2_Mixed(t *testing.T) {
	sess := New()
	sess.Merge([]scanner.Endpoint{
		{ID: "ep1", Method: "GET", Path: "/a"},
		{ID: "ep2", Method: "POST", Path: "/b"},
		{ID: "ep3", Method: "PUT", Path: "/c"},
	})
	sess.MarkPass("ep1")
	sess.MarkDone("ep2", 80)

	total, pass, todo := sess.Stats()
	if total != 3 || pass != 2 || todo != 1 {
		t.Fatalf("unexpected: total=%d pass=%d todo=%d", total, pass, todo)
	}
}

func TestStats2_Empty(t *testing.T) {
	sess := New()
	total, pass, todo := sess.Stats()
	if total != 0 || pass != 0 || todo != 0 {
		t.Fatalf("unexpected: total=%d pass=%d todo=%d", total, pass, todo)
	}
}

func TestStats2_ImproveIsTodo(t *testing.T) {
	sess := New()
	sess.Merge([]scanner.Endpoint{{ID: "ep1", Method: "GET", Path: "/a"}})
	sess.MarkImprove("ep1", 50)

	_, _, todo := sess.Stats()
	if todo != 1 {
		t.Fatalf("expected IMPROVE counted as todo, got %d", todo)
	}
}
