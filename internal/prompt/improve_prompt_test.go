package prompt

import (
	"strings"
	"testing"

	"github.com/park-jun-woo/huma/internal/scanner"
)

func TestImprovePrompt2(t *testing.T) {
	ep := &scanner.Endpoint{
		ID: "ep1", Method: "GET", Path: "/users",
		Handler: "GetUsers", Source: "handler.go", Line: 1,
	}
	reason := "R-COV: coverage 70% (7/10)\nuncovered handler.go:15  if err != nil {"

	result := ImprovePrompt(ep, "hurl/get_users.hurl", reason)
	if !strings.Contains(result, "# IMPROVE  GET /users") {
		t.Fatal("expected IMPROVE header")
	}
	if !strings.Contains(result, "Previous attempt fell short") {
		t.Fatal("expected shortfall section")
	}
	if !strings.Contains(result, "coverage 70% (7/10)") {
		t.Fatal("expected echoed reason / coverage")
	}
	if !strings.Contains(result, "handler.go:15") {
		t.Fatal("expected uncovered line reference from reason")
	}
	if !strings.Contains(result, "hurl/get_users.hurl") {
		t.Fatal("expected hurl file in instructions")
	}
}

func TestImprovePrompt2_NoReason(t *testing.T) {
	ep := &scanner.Endpoint{
		ID: "ep1", Method: "GET", Path: "/test",
	}
	result := ImprovePrompt(ep, "test.hurl", "")
	if strings.Contains(result, "Previous attempt fell short") {
		t.Fatal("empty reason should omit shortfall section")
	}
	if !strings.Contains(result, "test.hurl") {
		t.Fatal("expected hurl file in instructions")
	}
}
