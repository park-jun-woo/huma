package prompt

import (
	"strings"
	"testing"

	"github.com/park-jun-woo/huma/internal/scanner"
)

func TestFailPrompt2(t *testing.T) {
	ep := &scanner.Endpoint{
		ID: "ep1", Method: "GET", Path: "/users",
		Handler: "GetUsers", Source: "handler.go", Line: 1,
	}

	result := FailPrompt(ep, "hurl/get_users.hurl", "assertion failed at line 5")
	if !strings.Contains(result, "# FAIL  GET /users") {
		t.Fatal("expected FAIL header")
	}
	if !strings.Contains(result, "assertion failed at line 5") {
		t.Fatal("expected feedback")
	}
	if !strings.Contains(result, "hurl/get_users.hurl") {
		t.Fatal("expected hurl file path")
	}
}
