package prompt

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/park-jun-woo/huma/internal/analyzer"
	"github.com/park-jun-woo/huma/internal/config"
	"github.com/park-jun-woo/huma/internal/scanner"
)

func TestFormatBranchLine(t *testing.T) {
	// with Code
	got := formatBranchLine(analyzer.ResponseBranch{Status: 200, File: "/x/h.go", Line: 4, Code: "c.JSON(...)"})
	if !strings.Contains(got, "200 — h.go:4  c.JSON(...)") {
		t.Errorf("with code: %q", got)
	}
	// without Code
	got = formatBranchLine(analyzer.ResponseBranch{Status: 404, File: "/x/h.go", Line: 5})
	if !strings.Contains(got, "404 — h.go:5") || strings.Contains(got, "  c.") {
		t.Errorf("without code: %q", got)
	}
}

func TestCollectBranchSection(t *testing.T) {
	lines, statusList := collectBranchSection([]analyzer.ResponseBranch{
		{Status: 200, File: "h.go", Line: 4},
		{Status: 404, File: "h.go", Line: 5},
	})
	if !strings.Contains(lines, "200") || !strings.Contains(lines, "404") {
		t.Errorf("lines = %q", lines)
	}
	if statusList != "200, 404" {
		t.Errorf("statusList = %q, want '200, 404'", statusList)
	}
	// empty
	lines, statusList = collectBranchSection(nil)
	if lines != "" || statusList != "" {
		t.Errorf("empty → (%q,%q)", lines, statusList)
	}
}

func TestSetupPrompt(t *testing.T) {
	out := SetupPrompt()
	for _, want := range []string{"SETUP", "manifest.yaml", "apiVersion: yongol/v1", "huma next"} {
		if !strings.Contains(out, want) {
			t.Errorf("missing %q", want)
		}
	}
}

func TestStartPrompt(t *testing.T) {
	cfg := &config.Config{
		BaseURL: "http://localhost:8080",
		Server:  config.ServerConfig{Ready: "/health", Start: "./server"},
		Deps:    config.DepsConfig{Up: "docker compose up -d"},
	}
	out := StartPrompt(cfg)
	for _, want := range []string{"START", "http://localhost:8080/health", "docker compose up -d", "./server"} {
		if !strings.Contains(out, want) {
			t.Errorf("missing %q in %q", want, out)
		}
	}
	// no deps.up → that line omitted
	cfg.Deps.Up = ""
	out = StartPrompt(cfg)
	if strings.Contains(out, "docker compose") {
		t.Errorf("deps.up empty should omit line: %q", out)
	}
}

func TestUnverifiedPrompt(t *testing.T) {
	ep := &scanner.Endpoint{Method: "GET", Path: "/x"}
	// unlinked + static mode
	out := UnverifiedPrompt(ep, &config.Config{})
	if !strings.Contains(out, "source unlinked") || !strings.Contains(out, "static mode") {
		t.Errorf("unlinked/static: %q", out)
	}
	// linked + instrumented
	ep2 := &scanner.Endpoint{Method: "GET", Path: "/x", Source: "h.go"}
	out = UnverifiedPrompt(ep2, &config.Config{Server: config.ServerConfig{Start: "./s"}})
	if !strings.Contains(out, "source: h.go") || !strings.Contains(out, "Total==0") {
		t.Errorf("linked/instrumented: %q", out)
	}
	// nil cfg → uninstrumented branch (cfg!=nil guard is false)
	out = UnverifiedPrompt(ep, nil)
	if !strings.Contains(out, "Total==0") {
		t.Errorf("nil cfg → uninstrumented: %q", out)
	}
}

func TestStaticTodoPrompt_WithBranches(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "h.go")
	os.WriteFile(src, []byte("package main\n\nfunc GetUser(c interface{}) {\n\t_ = c\n}\n"), 0o644)
	ep := &scanner.Endpoint{Method: "GET", Path: "/users/:id", Handler: "GetUser", Source: src, Line: 3}
	branches := []analyzer.ResponseBranch{{Status: 200, File: src, Line: 4}, {Status: 404, File: src, Line: 5}}
	out := StaticTodoPrompt(ep, "hurl", "host", branches)
	for _, want := range []string{"# TODO  GET /users/:id", "Handler source", "Expected responses", "status 200, 404", "huma next"} {
		if !strings.Contains(out, want) {
			t.Errorf("missing %q in %q", want, out)
		}
	}
}

func TestStaticTodoPrompt_NoBranches(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "h.go")
	os.WriteFile(src, []byte("package main\n\nfunc GetUser(c interface{}) {\n\t_ = c\n}\n"), 0o644)
	ep := &scanner.Endpoint{Method: "GET", Path: "/users/:id", Handler: "GetUser", Source: src, Line: 3}
	out := StaticTodoPrompt(ep, "hurl", "host", nil)
	if !strings.Contains(out, "Hurl example") {
		t.Errorf("no branches should show hurl example: %q", out)
	}
}
