package runner

import (
	"testing"

	"github.com/park-jun-woo/huma/internal/scanner"
)

func TestHurlFileName_Simple(t *testing.T) {
	ep := &scanner.Endpoint{Method: "GET", Path: "/users"}
	result := hurlFileName(ep)
	if result != "get_users.hurl" {
		t.Fatalf("expected get_users.hurl, got %s", result)
	}
}

func TestHurlFileName_WithParam(t *testing.T) {
	ep := &scanner.Endpoint{Method: "POST", Path: "/users/:id/posts"}
	result := hurlFileName(ep)
	if result != "post_users_id_posts.hurl" {
		t.Fatalf("expected post_users_id_posts.hurl, got %s", result)
	}
}
