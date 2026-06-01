package hurlcheck

import "testing"

func TestStartEntry(t *testing.T) {
	a := &entryAccumulator{}
	a.startEntry("GET {{host}}/users")
	if a.cur == nil || a.cur.Method != "GET" || a.cur.URL != "{{host}}/users" {
		t.Fatalf("expected GET entry, got %+v", a.cur)
	}

	// starting a new entry flushes the previous one into entries
	a.cur.Status = 200
	a.startEntry("POST {{host}}/orders")
	if len(a.entries) != 1 {
		t.Fatalf("expected previous entry flushed, got %d", len(a.entries))
	}
	if a.entries[0].Method != "GET" {
		t.Errorf("flushed entry method = %s, want GET", a.entries[0].Method)
	}
	if a.cur.Method != "POST" {
		t.Errorf("new cur method = %s, want POST", a.cur.Method)
	}
}
