package scanner

import "testing"

func TestParseRawEndpoint_Valid(t *testing.T) {
	r := rawEndpoint{Method: "GET", Path: "/users", Handler: "GetUsers", File: "handler.go", Line: 10}
	ep := parseRawEndpoint(r)
	if ep == nil {
		t.Fatal("expected non-nil endpoint")
	}
	if ep.Method != "GET" || ep.Path != "/users" {
		t.Fatalf("unexpected: %+v", ep)
	}
	if ep.Handler != "GetUsers" {
		t.Fatalf("expected GetUsers, got %s", ep.Handler)
	}
}

func TestParseRawEndpoint_EmptyMethod(t *testing.T) {
	r := rawEndpoint{Method: "", Path: "/users"}
	ep := parseRawEndpoint(r)
	if ep != nil {
		t.Fatal("expected nil for empty method")
	}
}

func TestParseRawEndpoint_EmptyPath(t *testing.T) {
	r := rawEndpoint{Method: "GET", Path: ""}
	ep := parseRawEndpoint(r)
	if ep != nil {
		t.Fatal("expected nil for empty path")
	}
}

func TestParseRawEndpoint_HumaFormat(t *testing.T) {
	r := rawEndpoint{Method: "GET", Path: "/users", Handler: "handler.go:GetUsers"}
	ep := parseRawEndpoint(r)
	if ep == nil {
		t.Fatal("expected non-nil endpoint")
	}
	if ep.Handler != "GetUsers" {
		t.Fatalf("expected GetUsers, got %s", ep.Handler)
	}
	if ep.Source != "handler.go" {
		t.Fatalf("expected handler.go, got %s", ep.Source)
	}
}

func TestParseRawEndpoint_HumaFormatWithFile(t *testing.T) {
	r := rawEndpoint{Method: "GET", Path: "/users", Handler: "other.go:GetUsers", File: "original.go"}
	ep := parseRawEndpoint(r)
	if ep == nil {
		t.Fatal("expected non-nil endpoint")
	}
	if ep.Source != "original.go" {
		t.Fatalf("expected original.go (file takes precedence), got %s", ep.Source)
	}
}
