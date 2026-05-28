package analyzer

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGoAnalyzer_BasicJSON(t *testing.T) {
	src := `package main

import "net/http"

func CreateUser(c interface{}) {
	c.JSON(http.StatusCreated, nil)
	c.JSON(400, nil)
	c.AbortWithStatusJSON(http.StatusConflict, nil)
}
`
	dir := t.TempDir()
	file := filepath.Join(dir, "handler.go")
	os.WriteFile(file, []byte(src), 0o644)

	a := &GoAnalyzer{}
	branches, err := a.Analyze(file, "CreateUser", 0, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(branches) != 3 {
		t.Fatalf("expected 3 branches, got %d", len(branches))
	}

	statuses := map[int]bool{}
	for _, b := range branches {
		statuses[b.Status] = true
	}
	if !statuses[201] {
		t.Fatal("expected 201")
	}
	if !statuses[400] {
		t.Fatal("expected 400")
	}
	if !statuses[409] {
		t.Fatal("expected 409")
	}
}

func TestGoAnalyzer_StatusAndAbortWithStatus(t *testing.T) {
	src := `package main

import "net/http"

func NoContent(c interface{}) {
	c.Status(204)
	c.AbortWithStatus(http.StatusForbidden)
}
`
	dir := t.TempDir()
	file := filepath.Join(dir, "handler.go")
	os.WriteFile(file, []byte(src), 0o644)

	a := &GoAnalyzer{}
	branches, err := a.Analyze(file, "NoContent", 0, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(branches) != 2 {
		t.Fatalf("expected 2 branches, got %d", len(branches))
	}
	if branches[0].Status != 204 {
		t.Fatalf("expected 204, got %d", branches[0].Status)
	}
	if branches[1].Status != 403 {
		t.Fatalf("expected 403, got %d", branches[1].Status)
	}
}

func TestGoAnalyzer_HandlerNotFound(t *testing.T) {
	src := `package main

func OtherFunc() {}
`
	dir := t.TempDir()
	file := filepath.Join(dir, "handler.go")
	os.WriteFile(file, []byte(src), 0o644)

	a := &GoAnalyzer{}
	branches, err := a.Analyze(file, "NonExistent", 0, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if branches != nil {
		t.Fatalf("expected nil, got %v", branches)
	}
}

func TestGoAnalyzer_InvalidFile(t *testing.T) {
	a := &GoAnalyzer{}
	_, err := a.Analyze("/nonexistent/file.go", "Func", 0, 0)
	if err == nil {
		t.Fatal("expected error for nonexistent file")
	}
}

func TestGoAnalyzer_HelperCallTracing(t *testing.T) {
	src := `package main

import "net/http"

func errorResponse(c interface{}, status int, msg string) {
	c.JSON(status, msg)
}

func Handler(c interface{}) {
	c.JSON(http.StatusOK, nil)
	errorResponse(c, http.StatusBadRequest, "invalid")
}
`
	dir := t.TempDir()
	file := filepath.Join(dir, "handler.go")
	os.WriteFile(file, []byte(src), 0o644)

	a := &GoAnalyzer{}
	branches, err := a.Analyze(file, "Handler", 0, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(branches) != 2 {
		t.Fatalf("expected 2 branches, got %d", len(branches))
	}

	statuses := map[int]bool{}
	for _, b := range branches {
		statuses[b.Status] = true
	}
	if !statuses[200] {
		t.Fatal("expected 200")
	}
	if !statuses[400] {
		t.Fatal("expected 400")
	}
}

func TestGoAnalyzer_VariableStatusSkipped(t *testing.T) {
	src := `package main

func Handler(c interface{}) {
	status := getStatus()
	c.JSON(status, nil)
}

func getStatus() int { return 200 }
`
	dir := t.TempDir()
	file := filepath.Join(dir, "handler.go")
	os.WriteFile(file, []byte(src), 0o644)

	a := &GoAnalyzer{}
	branches, err := a.Analyze(file, "Handler", 0, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Variable status should be skipped
	if len(branches) != 0 {
		t.Fatalf("expected 0 branches (variable skipped), got %d", len(branches))
	}
}

func TestGoAnalyzer_ParseError(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "bad.go")
	os.WriteFile(file, []byte("not valid go code"), 0o644)

	a := &GoAnalyzer{}
	_, err := a.Analyze(file, "Handler", 0, 0)
	if err == nil {
		t.Fatal("expected error for invalid Go file")
	}
}

func TestGoAnalyzer_FiberSendStatus(t *testing.T) {
	src := `package main

func Handler(c interface{}) {
	c.SendStatus(204)
}
`
	dir := t.TempDir()
	file := filepath.Join(dir, "handler.go")
	os.WriteFile(file, []byte(src), 0o644)

	a := &GoAnalyzer{}
	branches, err := a.Analyze(file, "Handler", 0, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(branches) != 1 {
		t.Fatalf("expected 1 branch, got %d", len(branches))
	}
	if branches[0].Status != 204 {
		t.Fatalf("expected 204, got %d", branches[0].Status)
	}
}

func TestGoAnalyzer_FiberNewError(t *testing.T) {
	src := `package main

import "github.com/gofiber/fiber/v2"

func Handler(c interface{}) error {
	return fiber.NewError(400, "bad request")
}
`
	dir := t.TempDir()
	file := filepath.Join(dir, "handler.go")
	os.WriteFile(file, []byte(src), 0o644)

	a := &GoAnalyzer{}
	branches, err := a.Analyze(file, "Handler", 0, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(branches) != 1 {
		t.Fatalf("expected 1 branch, got %d", len(branches))
	}
	if branches[0].Status != 400 {
		t.Fatalf("expected 400, got %d", branches[0].Status)
	}
}

func TestGoAnalyzer_FiberErrConstant(t *testing.T) {
	src := `package main

import "github.com/gofiber/fiber/v2"

func Handler(c interface{}) error {
	return fiber.ErrNotFound
}
`
	dir := t.TempDir()
	file := filepath.Join(dir, "handler.go")
	os.WriteFile(file, []byte(src), 0o644)

	a := &GoAnalyzer{}
	branches, err := a.Analyze(file, "Handler", 0, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(branches) != 1 {
		t.Fatalf("expected 1 branch, got %d", len(branches))
	}
	if branches[0].Status != 404 {
		t.Fatalf("expected 404, got %d", branches[0].Status)
	}
}

func TestGoAnalyzer_FiberStatusChain(t *testing.T) {
	src := `package main

func Handler(c interface{}) error {
	return c.Status(201).JSON(data)
}
`
	dir := t.TempDir()
	file := filepath.Join(dir, "handler.go")
	os.WriteFile(file, []byte(src), 0o644)

	a := &GoAnalyzer{}
	branches, err := a.Analyze(file, "Handler", 0, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Status(201) should be extracted
	found201 := false
	for _, b := range branches {
		if b.Status == 201 {
			found201 = true
		}
	}
	if !found201 {
		t.Fatalf("expected 201 from Status chain, got branches: %v", branches)
	}
}

func TestGoAnalyzer_EchoNoContent(t *testing.T) {
	src := `package main

func Handler(c interface{}) error {
	return c.NoContent(204)
}
`
	dir := t.TempDir()
	file := filepath.Join(dir, "handler.go")
	os.WriteFile(file, []byte(src), 0o644)

	a := &GoAnalyzer{}
	branches, err := a.Analyze(file, "Handler", 0, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(branches) != 1 {
		t.Fatalf("expected 1 branch, got %d", len(branches))
	}
	if branches[0].Status != 204 {
		t.Fatalf("expected 204, got %d", branches[0].Status)
	}
}

func TestGoAnalyzer_EchoStringHTML(t *testing.T) {
	src := `package main

import "net/http"

func Handler(c interface{}) error {
	c.String(http.StatusOK, "hello")
	return c.HTML(200, "<h1>hi</h1>")
}
`
	dir := t.TempDir()
	file := filepath.Join(dir, "handler.go")
	os.WriteFile(file, []byte(src), 0o644)

	a := &GoAnalyzer{}
	branches, err := a.Analyze(file, "Handler", 0, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(branches) != 2 {
		t.Fatalf("expected 2 branches, got %d", len(branches))
	}
	statuses := map[int]bool{}
	for _, b := range branches {
		statuses[b.Status] = true
	}
	if !statuses[200] {
		t.Fatal("expected 200")
	}
}

func TestGoAnalyzer_EchoRedirect(t *testing.T) {
	src := `package main

func Handler(c interface{}) error {
	return c.Redirect(301, "/new")
}
`
	dir := t.TempDir()
	file := filepath.Join(dir, "handler.go")
	os.WriteFile(file, []byte(src), 0o644)

	a := &GoAnalyzer{}
	branches, err := a.Analyze(file, "Handler", 0, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(branches) != 1 {
		t.Fatalf("expected 1 branch, got %d", len(branches))
	}
	if branches[0].Status != 301 {
		t.Fatalf("expected 301, got %d", branches[0].Status)
	}
}

func TestGoAnalyzer_EchoNewHTTPError(t *testing.T) {
	src := `package main

import "github.com/labstack/echo/v4"

func Handler(c interface{}) error {
	return echo.NewHTTPError(400, "bad request")
}
`
	dir := t.TempDir()
	file := filepath.Join(dir, "handler.go")
	os.WriteFile(file, []byte(src), 0o644)

	a := &GoAnalyzer{}
	branches, err := a.Analyze(file, "Handler", 0, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(branches) != 1 {
		t.Fatalf("expected 1 branch, got %d", len(branches))
	}
	if branches[0].Status != 400 {
		t.Fatalf("expected 400, got %d", branches[0].Status)
	}
}

func TestGoAnalyzer_EchoErrConstant(t *testing.T) {
	src := `package main

import "github.com/labstack/echo/v4"

func Handler(c interface{}) error {
	return echo.ErrNotFound
}
`
	dir := t.TempDir()
	file := filepath.Join(dir, "handler.go")
	os.WriteFile(file, []byte(src), 0o644)

	a := &GoAnalyzer{}
	branches, err := a.Analyze(file, "Handler", 0, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(branches) != 1 {
		t.Fatalf("expected 1 branch, got %d", len(branches))
	}
	if branches[0].Status != 404 {
		t.Fatalf("expected 404, got %d", branches[0].Status)
	}
}

func TestGoAnalyzer_EchoXMLBlobPatterns(t *testing.T) {
	src := `package main

func Handler(c interface{}) error {
	c.XML(200, data)
	c.JSONBlob(201, blob)
	return c.Blob(202, "application/octet-stream", blob)
}
`
	dir := t.TempDir()
	file := filepath.Join(dir, "handler.go")
	os.WriteFile(file, []byte(src), 0o644)

	a := &GoAnalyzer{}
	branches, err := a.Analyze(file, "Handler", 0, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(branches) != 3 {
		t.Fatalf("expected 3 branches, got %d", len(branches))
	}
	statuses := map[int]bool{}
	for _, b := range branches {
		statuses[b.Status] = true
	}
	if !statuses[200] {
		t.Fatal("expected 200 from XML")
	}
	if !statuses[201] {
		t.Fatal("expected 201 from JSONBlob")
	}
	if !statuses[202] {
		t.Fatal("expected 202 from Blob")
	}
}

func TestGoAnalyzer_EchoMultiplePatterns(t *testing.T) {
	src := `package main

import "github.com/labstack/echo/v4"

func Handler(c interface{}) error {
	if invalid {
		return echo.ErrBadRequest
	}
	if notfound {
		return echo.NewHTTPError(404, "not found")
	}
	return c.NoContent(204)
}
`
	dir := t.TempDir()
	file := filepath.Join(dir, "handler.go")
	os.WriteFile(file, []byte(src), 0o644)

	a := &GoAnalyzer{}
	branches, err := a.Analyze(file, "Handler", 0, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(branches) != 3 {
		t.Fatalf("expected 3 branches, got %d", len(branches))
	}

	statuses := map[int]bool{}
	for _, b := range branches {
		statuses[b.Status] = true
	}
	if !statuses[400] {
		t.Fatal("expected 400 from echo.ErrBadRequest")
	}
	if !statuses[404] {
		t.Fatal("expected 404 from echo.NewHTTPError")
	}
	if !statuses[204] {
		t.Fatal("expected 204 from NoContent")
	}
}

func TestGoAnalyzer_FiberMultiplePatterns(t *testing.T) {
	src := `package main

import "github.com/gofiber/fiber/v2"

func Handler(c interface{}) error {
	if invalid {
		return fiber.ErrBadRequest
	}
	if notfound {
		return fiber.NewError(404, "not found")
	}
	c.SendStatus(204)
	return nil
}
`
	dir := t.TempDir()
	file := filepath.Join(dir, "handler.go")
	os.WriteFile(file, []byte(src), 0o644)

	a := &GoAnalyzer{}
	branches, err := a.Analyze(file, "Handler", 0, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(branches) != 3 {
		t.Fatalf("expected 3 branches, got %d", len(branches))
	}

	statuses := map[int]bool{}
	for _, b := range branches {
		statuses[b.Status] = true
	}
	if !statuses[400] {
		t.Fatal("expected 400 from fiber.ErrBadRequest")
	}
	if !statuses[404] {
		t.Fatal("expected 404 from fiber.NewError")
	}
	if !statuses[204] {
		t.Fatal("expected 204 from SendStatus")
	}
}
