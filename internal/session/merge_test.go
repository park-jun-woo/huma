package session

import (
	"testing"

	"github.com/park-jun-woo/huma/internal/scanner"
)

func TestMerge2_NewEndpoints(t *testing.T) {
	sess := New()
	sess.Merge([]scanner.Endpoint{
		{ID: "ep1", Method: "GET", Path: "/a"},
		{ID: "ep2", Method: "POST", Path: "/b"},
	})
	if len(sess.Entries) != 2 {
		t.Fatalf("expected 2, got %d", len(sess.Entries))
	}
	if sess.Entries[0].Status != StatusTodo {
		t.Fatal("expected TODO")
	}
}

func TestMerge2_PreservesStatus(t *testing.T) {
	sess := New()
	sess.Merge([]scanner.Endpoint{{ID: "ep1", Method: "GET", Path: "/a"}})
	sess.MarkPass("ep1")

	sess.Merge([]scanner.Endpoint{{ID: "ep1", Method: "GET", Path: "/a"}})
	if sess.Entries[0].Status != StatusPass {
		t.Fatalf("expected PASS preserved, got %s", sess.Entries[0].Status)
	}
}

func TestMerge2_AddsNew(t *testing.T) {
	sess := New()
	sess.Merge([]scanner.Endpoint{{ID: "ep1", Method: "GET", Path: "/a"}})
	sess.MarkPass("ep1")

	sess.Merge([]scanner.Endpoint{
		{ID: "ep1", Method: "GET", Path: "/a"},
		{ID: "ep2", Method: "POST", Path: "/b"},
	})
	if len(sess.Entries) != 2 {
		t.Fatalf("expected 2, got %d", len(sess.Entries))
	}
	if sess.Entries[1].Status != StatusTodo {
		t.Fatal("expected new entry as TODO")
	}
}
