package session

import (
	"testing"

	"github.com/park-jun-woo/huma/internal/scanner"
)

func TestMarkUnverified(t *testing.T) {
	s := New()
	s.Merge([]scanner.Endpoint{{ID: "ep1", Method: "GET", Path: "/a"}})
	s.MarkUnverified("ep1")
	if s.Entries[0].Status != StatusUnverified {
		t.Fatalf("expected UNVERIFIED, got %s", s.Entries[0].Status)
	}
	if s.Entries[0].CRI != 0 {
		t.Fatalf("expected CRI 0, got %d", s.Entries[0].CRI)
	}
}

func TestSetVerdict(t *testing.T) {
	s := New()
	s.Merge([]scanner.Endpoint{{ID: "ep1", Method: "GET", Path: "/a"}})
	s.SetVerdict("ep1", 3, 2, "both")
	e := s.Entries[0]
	if e.CRI != 3 || e.AGrade != 2 || e.Provenance != "both" {
		t.Fatalf("unexpected verdict fields: %+v", e)
	}
}

func TestUnverifiedCount(t *testing.T) {
	s := New()
	s.Merge([]scanner.Endpoint{
		{ID: "ep1", Method: "GET", Path: "/a"},
		{ID: "ep2", Method: "GET", Path: "/b"},
	})
	s.MarkUnverified("ep1")
	if n := s.Unverified(); n != 1 {
		t.Fatalf("expected 1 unverified, got %d", n)
	}
	// UNVERIFIED counts as todo, never pass.
	_, pass, todo := s.Stats()
	if pass != 0 || todo != 2 {
		t.Fatalf("expected pass=0 todo=2, got pass=%d todo=%d", pass, todo)
	}
}

func TestCRILabel(t *testing.T) {
	cases := map[int]string{0: "UNVERIFIED", 1: "SCAFFOLDED", 2: "SMOKE", 3: "COVERED"}
	for cri, want := range cases {
		if got := CRILabel(cri); got != want {
			t.Errorf("CRILabel(%d)=%s, want %s", cri, got, want)
		}
	}
}

func TestCurrent_IncludesUnverified(t *testing.T) {
	s := New()
	s.Merge([]scanner.Endpoint{{ID: "ep1", Method: "GET", Path: "/a"}})
	s.MarkUnverified("ep1")
	if s.Current() == nil {
		t.Fatal("expected Current to return UNVERIFIED entry as needing work")
	}
}
