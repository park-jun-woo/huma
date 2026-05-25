package prompt

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/park-jun-woo/huma/internal/scanner"
)

func TestTodoPrompt2_WithHandler(t *testing.T) {
	tmpDir := t.TempDir()
	handlerFile := filepath.Join(tmpDir, "handler.go")
	content := `package main

func GetUser(c interface{}) {
	// handler body
}
`
	os.WriteFile(handlerFile, []byte(content), 0o644)

	ep := &scanner.Endpoint{
		ID: "ep1", Method: "GET", Path: "/users/:id",
		Handler: "GetUser", Source: handlerFile, Line: 3,
	}

	result := TodoPrompt(ep, "hurl", "base_url")
	if !strings.Contains(result, "# TODO  GET /users/:id") {
		t.Fatal("expected TODO header")
	}
	if !strings.Contains(result, "## Handler source") {
		t.Fatal("expected handler source section")
	}
	if !strings.Contains(result, "## Hurl example") {
		t.Fatal("expected hurl example section")
	}
}

func TestTodoPrompt2_WithoutHandler(t *testing.T) {
	ep := &scanner.Endpoint{
		ID: "ep1", Method: "POST", Path: "/users",
		Handler: "CreateUser", Source: "/nonexistent/file.go", Line: 1,
	}

	result := TodoPrompt(ep, "hurl", "base_url")
	if !strings.Contains(result, "# TODO  POST /users") {
		t.Fatal("expected TODO header")
	}
	if strings.Contains(result, "## Handler source") {
		t.Fatal("should not have handler source for missing file")
	}
}
