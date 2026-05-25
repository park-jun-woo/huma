package runner

import (
	"testing"

	"github.com/park-jun-woo/huma/internal/scanner"
)

func TestHurlFileName2_Simple(t *testing.T) {
	ep := &scanner.Endpoint{Method: "GET", Path: "/users"}
	result := HurlFileName(ep, "hurl")
	if result != "hurl/get_users.hurl" {
		t.Fatalf("expected hurl/get_users.hurl, got %s", result)
	}
}

func TestHurlFileName2_WithParam(t *testing.T) {
	ep := &scanner.Endpoint{Method: "DELETE", Path: "/users/:id"}
	result := HurlFileName(ep, "tests/hurl")
	if result != "tests/hurl/delete_users_id.hurl" {
		t.Fatalf("expected tests/hurl/delete_users_id.hurl, got %s", result)
	}
}
