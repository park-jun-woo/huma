package session

import "testing"

func TestNew2(t *testing.T) {
	sess := New()
	if sess == nil {
		t.Fatal("expected non-nil")
	}
	if len(sess.Entries) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(sess.Entries))
	}
}
