package scanner

import (
	"testing"
)

func TestParseEndpoints_Basic(t *testing.T) {
	data := []byte(`[
		{"method": "GET", "path": "/users", "handler": "ListUsers", "file": "handler.go", "line": 10},
		{"method": "POST", "path": "/users", "handler": "CreateUser", "file": "handler.go", "line": 20}
	]`)

	endpoints, err := ParseEndpoints(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(endpoints) != 2 {
		t.Fatalf("expected 2 endpoints, got %d", len(endpoints))
	}
	if endpoints[0].Method != "GET" || endpoints[0].Path != "/users" {
		t.Fatalf("unexpected first endpoint: %+v", endpoints[0])
	}
	if endpoints[0].Handler != "ListUsers" {
		t.Fatalf("expected ListUsers, got %s", endpoints[0].Handler)
	}
	if endpoints[0].Source != "handler.go" {
		t.Fatalf("expected handler.go, got %s", endpoints[0].Source)
	}
	if endpoints[0].Line != 10 {
		t.Fatalf("expected line 10, got %d", endpoints[0].Line)
	}
	if endpoints[0].ID == "" {
		t.Fatal("expected non-empty ID")
	}
}

func TestParseEndpoints_HumaFormat(t *testing.T) {
	data := []byte(`[
		{"method": "GET", "path": "/api/v1/buildings", "handler": "internal/api/building/handler.go:ListBuildings", "file": "internal/api/building/handler.go", "line": 15}
	]`)

	endpoints, err := ParseEndpoints(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(endpoints) != 1 {
		t.Fatalf("expected 1 endpoint, got %d", len(endpoints))
	}
	if endpoints[0].Handler != "ListBuildings" {
		t.Fatalf("expected ListBuildings, got %s", endpoints[0].Handler)
	}
	if endpoints[0].Source != "internal/api/building/handler.go" {
		t.Fatalf("expected file from handler split, got %s", endpoints[0].Source)
	}
}

func TestParseEndpoints_HumaFormatNoFile(t *testing.T) {
	data := []byte(`[
		{"method": "POST", "path": "/login", "handler": "auth/handler.go:Login"}
	]`)

	endpoints, err := ParseEndpoints(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(endpoints) != 1 {
		t.Fatalf("expected 1 endpoint, got %d", len(endpoints))
	}
	if endpoints[0].Handler != "Login" {
		t.Fatalf("expected Login, got %s", endpoints[0].Handler)
	}
	if endpoints[0].Source != "auth/handler.go" {
		t.Fatalf("expected auth/handler.go, got %s", endpoints[0].Source)
	}
}

func TestParseEndpoints_SkipMissingMethodOrPath(t *testing.T) {
	data := []byte(`[
		{"method": "GET", "path": ""},
		{"method": "", "path": "/test"},
		{"method": "GET", "path": "/valid"}
	]`)

	endpoints, err := ParseEndpoints(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(endpoints) != 1 {
		t.Fatalf("expected 1 endpoint, got %d", len(endpoints))
	}
	if endpoints[0].Path != "/valid" {
		t.Fatalf("expected /valid, got %s", endpoints[0].Path)
	}
}

func TestParseEndpoints_InvalidJSON(t *testing.T) {
	data := []byte(`not json`)

	_, err := ParseEndpoints(data)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestParseEndpoints_EmptyArray(t *testing.T) {
	data := []byte(`[]`)

	endpoints, err := ParseEndpoints(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(endpoints) != 0 {
		t.Fatalf("expected 0 endpoints, got %d", len(endpoints))
	}
}

func TestParseEndpoints_MinimalFields(t *testing.T) {
	data := []byte(`[{"method": "DELETE", "path": "/items/:id"}]`)

	endpoints, err := ParseEndpoints(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(endpoints) != 1 {
		t.Fatalf("expected 1 endpoint, got %d", len(endpoints))
	}
	if endpoints[0].Handler != "" {
		t.Fatalf("expected empty handler, got %s", endpoints[0].Handler)
	}
	if endpoints[0].Source != "" {
		t.Fatalf("expected empty source, got %s", endpoints[0].Source)
	}
}

func TestMakeID(t *testing.T) {
	id1 := makeID("GET", "/users")
	id2 := makeID("POST", "/users")
	id3 := makeID("GET", "/users")

	if id1 == id2 {
		t.Fatal("different method+path should produce different IDs")
	}
	if id1 != id3 {
		t.Fatal("same method+path should produce same ID")
	}
	if len(id1) != 16 { // 8 bytes = 16 hex chars
		t.Fatalf("expected 16 hex chars, got %d", len(id1))
	}
}
