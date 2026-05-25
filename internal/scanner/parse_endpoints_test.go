package scanner

import "testing"

func TestParseEndpoints_Valid(t *testing.T) {
	data := []byte(`[
		{"method": "GET", "path": "/users", "handler": "GetUsers", "file": "handler.go", "line": 10},
		{"method": "POST", "path": "/users", "handler": "CreateUser", "file": "handler.go", "line": 20}
	]`)

	eps, err := ParseEndpoints(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(eps) != 2 {
		t.Fatalf("expected 2, got %d", len(eps))
	}
	if eps[0].Method != "GET" || eps[0].Path != "/users" {
		t.Fatalf("unexpected first endpoint: %+v", eps[0])
	}
}

func TestParseEndpoints_BadJSON(t *testing.T) {
	_, err := ParseEndpoints([]byte("INVALID"))
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestParseEndpoints_Empty(t *testing.T) {
	eps, err := ParseEndpoints([]byte("[]"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(eps) != 0 {
		t.Fatalf("expected 0, got %d", len(eps))
	}
}
