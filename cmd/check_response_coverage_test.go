package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/park-jun-woo/huma/internal/scanner"
)

func TestCheckResponseCoverage_WithResponses(t *testing.T) {
	dir := t.TempDir()
	hurlFile := filepath.Join(dir, "test.hurl")
	os.WriteFile(hurlFile, []byte("POST {{host}}/users\nHTTP 201\n"), 0o644)

	ep := &scanner.Endpoint{
		ID: "ep1", Method: "POST", Path: "/users",
		Handler: "CreateUser", Source: "handler.go", Line: 10,
		Responses: json.RawMessage(`[{"status": 201, "line": 20}, {"status": 400, "line": 25}]`),
	}

	result := checkResponseCoverage(ep, hurlFile)
	if result == nil {
		t.Fatal("expected non-nil result")
	}
	if result.Total != 2 {
		t.Fatalf("expected 2 total, got %d", result.Total)
	}
	if result.Covered != 1 {
		t.Fatalf("expected 1 covered, got %d", result.Covered)
	}
	if len(result.Missing) != 1 {
		t.Fatalf("expected 1 missing, got %d", len(result.Missing))
	}
	if result.Missing[0].Status != 400 {
		t.Fatalf("expected 400 missing, got %d", result.Missing[0].Status)
	}
}

func TestCheckResponseCoverage_WithSourceFile(t *testing.T) {
	dir := t.TempDir()
	hurlFile := filepath.Join(dir, "test.hurl")
	os.WriteFile(hurlFile, []byte("GET {{host}}/users\nHTTP 200\n"), 0o644)

	srcFile := filepath.Join(dir, "handler.go")
	src := `package main

import "net/http"

func GetUsers(c interface{}) {
	c.JSON(http.StatusOK, nil)
	c.JSON(http.StatusInternalServerError, nil)
}
`
	os.WriteFile(srcFile, []byte(src), 0o644)

	ep := &scanner.Endpoint{
		ID: "ep1", Method: "GET", Path: "/users",
		Handler: "GetUsers", Source: srcFile, Line: 5,
	}

	result := checkResponseCoverage(ep, hurlFile)
	if result == nil {
		t.Fatal("expected non-nil result")
	}
	if result.Total != 2 {
		t.Fatalf("expected 2 total, got %d", result.Total)
	}
	if result.Covered != 1 {
		t.Fatalf("expected 1 covered, got %d", result.Covered)
	}
	if len(result.Missing) != 1 {
		t.Fatalf("expected 1 missing, got %d", len(result.Missing))
	}
}

func TestCheckResponseCoverage_NoSource(t *testing.T) {
	ep := &scanner.Endpoint{
		ID: "ep1", Method: "GET", Path: "/test",
		Handler: "H", Source: "", Line: 0,
	}

	result := checkResponseCoverage(ep, "nonexistent.hurl")
	if result != nil {
		t.Fatal("expected nil for no source")
	}
}

func TestCheckResponseCoverage_InvalidHurlFile(t *testing.T) {
	dir := t.TempDir()
	srcFile := filepath.Join(dir, "handler.go")
	src := `package main

import "net/http"

func H(c interface{}) {
	c.JSON(http.StatusOK, nil)
}
`
	os.WriteFile(srcFile, []byte(src), 0o644)

	ep := &scanner.Endpoint{
		ID: "ep1", Method: "GET", Path: "/test",
		Handler: "H", Source: srcFile, Line: 5,
	}

	result := checkResponseCoverage(ep, "/nonexistent/test.hurl")
	if result != nil {
		t.Fatal("expected nil for invalid hurl file")
	}
}

func TestCheckResponseCoverage_AllCovered(t *testing.T) {
	dir := t.TempDir()
	hurlFile := filepath.Join(dir, "test.hurl")
	os.WriteFile(hurlFile, []byte("GET {{host}}/users\nHTTP 200\n"), 0o644)

	srcFile := filepath.Join(dir, "handler.go")
	src := `package main

import "net/http"

func H(c interface{}) {
	c.JSON(http.StatusOK, nil)
}
`
	os.WriteFile(srcFile, []byte(src), 0o644)

	ep := &scanner.Endpoint{
		ID: "ep1", Method: "GET", Path: "/test",
		Handler: "H", Source: srcFile, Line: 5,
	}

	result := checkResponseCoverage(ep, hurlFile)
	if result == nil {
		t.Fatal("expected non-nil result")
	}
	if len(result.Missing) != 0 {
		t.Fatalf("expected 0 missing, got %d", len(result.Missing))
	}
}
