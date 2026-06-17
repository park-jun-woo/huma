package prompt

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/park-jun-woo/huma/internal/analyzer"
	"github.com/park-jun-woo/huma/internal/hurlcheck"
	"github.com/park-jun-woo/huma/internal/scanner"
)

func TestTodoPrompt_WithHandler(t *testing.T) {
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
	if !strings.Contains(result, "## Instructions") {
		t.Fatal("expected instructions section")
	}
}

func TestTodoPrompt_WithoutHandler(t *testing.T) {
	ep := &scanner.Endpoint{
		ID: "ep1", Method: "POST", Path: "/users",
		Handler: "CreateUser", Source: "/nonexistent/file.go", Line: 1,
	}

	result := TodoPrompt(ep, "hurl", "base_url")
	if !strings.Contains(result, "# TODO  POST /users") {
		t.Fatal("expected TODO header")
	}
	// Should not have handler source since file doesn't exist
	if strings.Contains(result, "## Handler source") {
		t.Fatal("should not have handler source for missing file")
	}
}

func TestFailPrompt(t *testing.T) {
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

func TestPassPrompt(t *testing.T) {
	ep := &scanner.Endpoint{
		ID: "ep1", Method: "DELETE", Path: "/users/:id",
		Handler: "DeleteUser", Source: "handler.go", Line: 1,
	}

	result := PassPrompt(ep)
	if result != "# PASS  DELETE /users/:id\n" {
		t.Fatalf("unexpected: %q", result)
	}
}

func TestAllComplete(t *testing.T) {
	result := AllComplete(10, 10)
	if !strings.Contains(result, "All endpoints complete!") {
		t.Fatal("expected completion message")
	}
	if !strings.Contains(result, "PASS: 10 (100%)") {
		t.Fatal("expected 100% pass")
	}
}

func TestAllComplete_ZeroTotal(t *testing.T) {
	result := AllComplete(0, 0)
	if !strings.Contains(result, "PASS: 0 (0%)") {
		t.Fatal("expected 0%")
	}
}

func TestImprovePrompt(t *testing.T) {
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
	if !strings.Contains(result, "handler.go:15") {
		t.Fatal("expected uncovered line reference from reason")
	}
}

func TestPercent(t *testing.T) {
	if percent(0, 0) != 0 {
		t.Fatal("expected 0 for zero total")
	}
	if percent(5, 10) != 50 {
		t.Fatal("expected 50")
	}
	if percent(10, 10) != 100 {
		t.Fatal("expected 100")
	}
	if percent(1, 3) != 33 {
		t.Fatal("expected 33")
	}
}

func TestHurlExample_GET(t *testing.T) {
	result := hurlExample("GET", "/users", "base_url")
	if !strings.Contains(result, "GET {{base_url}}/users") {
		t.Fatal("expected GET template")
	}
	if !strings.Contains(result, "count > 0") {
		t.Fatal("expected list assertion")
	}
}

func TestHurlExample_GETWithParam(t *testing.T) {
	result := hurlExample("GET", "/users/:id", "base_url")
	if !strings.Contains(result, "GET {{base_url}}/users/1") {
		t.Fatal("expected GET with param replaced")
	}
	if !strings.Contains(result, "$.id") {
		t.Fatal("expected id assertion")
	}
}

func TestHurlExample_POST(t *testing.T) {
	result := hurlExample("POST", "/users", "base_url")
	if !strings.Contains(result, "POST {{base_url}}/users") {
		t.Fatal("expected POST template")
	}
	if !strings.Contains(result, "HTTP 201") {
		t.Fatal("expected 201 status")
	}
}

func TestHurlExample_PUT(t *testing.T) {
	result := hurlExample("PUT", "/users/:id", "base_url")
	if !strings.Contains(result, "PUT {{base_url}}/users/1") {
		t.Fatal("expected PUT template")
	}
}

func TestHurlExample_PATCH(t *testing.T) {
	result := hurlExample("PATCH", "/users/:id", "base_url")
	if !strings.Contains(result, "PATCH {{base_url}}/users/1") {
		t.Fatal("expected PATCH template")
	}
}

func TestHurlExample_DELETE(t *testing.T) {
	result := hurlExample("DELETE", "/users/:id", "base_url")
	if !strings.Contains(result, "DELETE {{base_url}}/users/1") {
		t.Fatal("expected DELETE template")
	}
	if !strings.Contains(result, "HTTP 204") {
		t.Fatal("expected 204 status")
	}
}

func TestHurlExample_HEAD(t *testing.T) {
	result := hurlExample("HEAD", "/health", "base_url")
	if !strings.Contains(result, "HEAD {{base_url}}/health") {
		t.Fatal("expected HEAD template")
	}
	if !strings.Contains(result, "HTTP 200") {
		t.Fatal("expected 200 status")
	}
}

func TestHasParam(t *testing.T) {
	if !hasParam("/users/:id") {
		t.Fatal("expected true")
	}
	if hasParam("/users") {
		t.Fatal("expected false")
	}
}

func TestReplaceParams(t *testing.T) {
	result := replaceParams("/users/:id/posts/:postId")
	if result != "/users/1/posts/1" {
		t.Fatalf("expected /users/1/posts/1, got %s", result)
	}
}

func TestImprovePrompt_NoUncovered(t *testing.T) {
	ep := &scanner.Endpoint{
		ID: "ep1", Method: "GET", Path: "/test",
		Handler: "H", Source: "h.go", Line: 1,
	}
	result := ImprovePrompt(ep, "test.hurl", "")
	if strings.Contains(result, "Previous attempt fell short") {
		t.Fatal("empty reason should omit shortfall section")
	}
}

func TestResponseImprovePrompt(t *testing.T) {
	ep := &scanner.Endpoint{
		ID: "ep1", Method: "POST", Path: "/users",
		Handler: "CreateUser", Source: "handler.go", Line: 10,
	}
	respResult := &hurlcheck.ResponseCoverageResult{
		Covered: 2,
		Total:   4,
		Percent: 50,
		Missing: []analyzer.ResponseBranch{
			{Status: 409, File: "handler.go", Line: 62, Code: "c.JSON(http.StatusConflict, ...)"},
			{Status: 500, File: "handler.go", Line: 66},
		},
	}

	result := ResponseImprovePrompt(ep, "hurl/post_users.hurl", respResult)
	if !strings.Contains(result, "# IMPROVE  POST /users") {
		t.Fatal("expected IMPROVE header")
	}
	if !strings.Contains(result, "Response coverage: 2/4 (50%)") {
		t.Fatal("expected response coverage line")
	}
	if !strings.Contains(result, "MISSING") {
		t.Fatal("expected MISSING section")
	}
	if !strings.Contains(result, "409") {
		t.Fatal("expected 409 in missing")
	}
	if !strings.Contains(result, "500") {
		t.Fatal("expected 500 in missing")
	}
	if !strings.Contains(result, "handler.go:62") {
		t.Fatal("expected handler.go:62")
	}
}

func TestResponseImprovePrompt_NoMissing(t *testing.T) {
	ep := &scanner.Endpoint{
		ID: "ep1", Method: "GET", Path: "/test",
		Handler: "H", Source: "h.go", Line: 1,
	}
	respResult := &hurlcheck.ResponseCoverageResult{
		Covered: 2,
		Total:   2,
		Percent: 100,
	}

	result := ResponseImprovePrompt(ep, "test.hurl", respResult)
	if strings.Contains(result, "MISSING") {
		t.Fatal("should not have missing section")
	}
}

func TestFormatMissingBranch_WithCode(t *testing.T) {
	branch := analyzer.ResponseBranch{
		Status: 409, File: "internal/api/handler.go", Line: 62, Code: "c.JSON(...)",
	}
	result := formatMissingBranch(branch)
	if !strings.Contains(result, "409") {
		t.Fatal("expected 409")
	}
	if !strings.Contains(result, "handler.go:62") {
		t.Fatal("expected handler.go:62")
	}
	if !strings.Contains(result, "c.JSON(...)") {
		t.Fatal("expected code")
	}
}

func TestFormatMissingBranch_WithoutCode(t *testing.T) {
	branch := analyzer.ResponseBranch{
		Status: 500, File: "handler.go", Line: 66,
	}
	result := formatMissingBranch(branch)
	if !strings.Contains(result, "500") {
		t.Fatal("expected 500")
	}
	if !strings.Contains(result, "handler.go:66") {
		t.Fatal("expected handler.go:66")
	}
}
