package prompt

import (
	"testing"

	"github.com/park-jun-woo/huma/internal/scanner"
)

func TestPassPrompt2(t *testing.T) {
	ep := &scanner.Endpoint{Method: "DELETE", Path: "/users/:id"}
	result := PassPrompt(ep)
	if result != "# PASS  DELETE /users/:id\n" {
		t.Fatalf("unexpected: %q", result)
	}
}
