package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/park-jun-woo/huma/internal/adapter"
	"github.com/park-jun-woo/huma/internal/analyzer"
	"github.com/park-jun-woo/huma/internal/config"
	"github.com/park-jun-woo/huma/internal/scanner"
	"github.com/park-jun-woo/huma/internal/session"
)

// BUG-001: thin OpenAPI (status 1 each) + no source link + no instrumentation
// must NOT collapse to 100% PASS. It must be UNVERIFIED.
func TestStaticVerdict_ThinOpenAPI_Unverified(t *testing.T) {
	dir := t.TempDir()
	hurl := filepath.Join(dir, "t.hurl")
	os.WriteFile(hurl, []byte("GET {{host}}/x\nHTTP 200\n"), 0o644)

	ep := scanner.Endpoint{
		ID: "ep1", Method: "GET", Path: "/x",
		Source:    "", // unlinked (OpenAPI-driven scan)
		Responses: json.RawMessage(`[{"status":200}]`),
	}
	sess := session.New()
	sess.Merge([]scanner.Endpoint{ep})
	cfg := &config.Config{Scan: config.ScanConfig{Lang: "go"}}

	oc, _ := staticVerdict(cfg, sess, &ep, hurl)
	if oc != outcomeUnverified {
		t.Fatalf("expected UNVERIFIED, got %v", oc)
	}
	if sess.Entries[0].Status != session.StatusUnverified {
		t.Fatalf("expected status UNVERIFIED, got %s", sess.Entries[0].Status)
	}
}

// Source-linked endpoint with all branches covered → SCAFFOLDED PASS in static.
func TestStaticVerdict_SourceLinked_Scaffolded(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "h.go")
	os.WriteFile(src, []byte(`package main
import "net/http"
func H(c interface{}) {
	c.JSON(http.StatusOK, nil)
}
`), 0o644)
	hurl := filepath.Join(dir, "t.hurl")
	os.WriteFile(hurl, []byte("GET {{host}}/x\nHTTP 200\n"), 0o644)

	ep := scanner.Endpoint{ID: "ep1", Method: "GET", Path: "/x", Handler: "H", Source: src, Line: 3}
	sess := session.New()
	sess.Merge([]scanner.Endpoint{ep})
	cfg := &config.Config{Scan: config.ScanConfig{Lang: "go"}}

	oc, _ := staticVerdict(cfg, sess, &ep, hurl)
	if oc != outcomePass {
		t.Fatalf("expected PASS, got %v", oc)
	}
	if sess.Entries[0].CRI != 1 {
		t.Fatalf("expected CRI 1 (SCAFFOLDED), got %d", sess.Entries[0].CRI)
	}
}

// C-02: OpenAPI declaring fewer statuses than source must not shrink the
// denominator — source branches are the floor.
func TestResponseBranches_UnionMonotonic(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "h.go")
	os.WriteFile(src, []byte(`package main
import "net/http"
func H(c interface{}) {
	c.JSON(http.StatusOK, nil)
	c.JSON(http.StatusBadRequest, nil)
}
`), 0o644)
	ep := scanner.Endpoint{
		ID: "ep1", Method: "GET", Path: "/x", Handler: "H", Source: src, Line: 3,
		Responses: json.RawMessage(`[{"status":200}]`), // thin declaration
	}
	branches, prov := responseBranches(&ep, "go")
	statuses := map[int]bool{}
	for _, b := range branches {
		statuses[b.Status] = true
	}
	if !statuses[200] || !statuses[400] {
		t.Fatalf("expected union to contain 200 and 400, got %v", statuses)
	}
	if !prov.HasSource || !prov.HasDeclared {
		t.Fatalf("expected provenance both, got %s", prov.String())
	}
}

// Live: instrumented, status matched but branch line not hit → not COVERED → UNVERIFIED.
func TestLiveVerdict_BranchLineUnbound_Unverified(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "h.go")
	os.WriteFile(src, []byte(`package main
import "net/http"
func H(c interface{}) {
	c.JSON(http.StatusOK, nil)
	c.JSON(http.StatusBadRequest, nil)
}
`), 0o644)
	hurl := filepath.Join(dir, "t.hurl")
	os.WriteFile(hurl, []byte("GET {{host}}/x\nHTTP 200\nHTTP 400\n"), 0o644)

	ep := scanner.Endpoint{ID: "ep1", Method: "GET", Path: "/x", Handler: "H", Source: src, Line: 3}
	sess := session.New()
	sess.Merge([]scanner.Endpoint{ep})
	cfg := &config.Config{Scan: config.ScanConfig{Lang: "go"}, Server: config.ServerConfig{Start: "x"}}

	// Branches are at lines 4 (200) and 5 (400). Only line 4 covered → 400 unbound.
	cov := &adapter.CoverageResult{
		Covered: 1, Total: 2, Percent: 50,
		CoveredLines: map[int]bool{4: true},
	}
	entry := sess.CurrentEntry()
	oc := liveVerdict(cfg, sess, &ep, hurl, cov, entry)
	if oc != outcomeImprove {
		t.Fatalf("expected IMPROVE on first <100%% attempt, got %v", oc)
	}
}

// Vacuous green: skip:true entry must not count as covered.
func TestStaticVerdict_VacuousSkip_NotCovered(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "h.go")
	os.WriteFile(src, []byte(`package main
import "net/http"
func H(c interface{}) {
	c.JSON(http.StatusOK, nil)
}
`), 0o644)
	hurl := filepath.Join(dir, "t.hurl")
	// status 200 present but skipped → vacuous, must not cover the branch.
	os.WriteFile(hurl, []byte("GET {{host}}/x\n[Options]\nskip: true\nHTTP 200\n"), 0o644)

	ep := scanner.Endpoint{ID: "ep1", Method: "GET", Path: "/x", Handler: "H", Source: src, Line: 3}
	sess := session.New()
	sess.Merge([]scanner.Endpoint{ep})
	cfg := &config.Config{Scan: config.ScanConfig{Lang: "go"}}

	oc, res := staticVerdict(cfg, sess, &ep, hurl)
	if oc != outcomeImprove {
		t.Fatalf("expected IMPROVE (200 uncovered due to skip), got %v", oc)
	}
	if res == nil || len(res.Missing) != 1 {
		t.Fatalf("expected 200 missing, got %+v", res)
	}
}

// Live instrumented, all branches runtime-bound at 100% → COVERED PASS.
func TestLiveVerdict_Covered(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "h.go")
	os.WriteFile(src, []byte(`package main
import "net/http"
func H(c interface{}) {
	c.JSON(http.StatusOK, nil)
}
`), 0o644)
	hurl := filepath.Join(dir, "t.hurl")
	os.WriteFile(hurl, []byte("GET {{host}}/x\nHTTP 200\n[Asserts]\njsonpath \"$.id\" exists\n"), 0o644)

	ep := scanner.Endpoint{ID: "ep1", Method: "GET", Path: "/x", Handler: "H", Source: src, Line: 3}
	sess := session.New()
	sess.Merge([]scanner.Endpoint{ep})
	cfg := &config.Config{Scan: config.ScanConfig{Lang: "go"}, Server: config.ServerConfig{Start: "x"}}

	cov := &adapter.CoverageResult{
		Covered: 1, Total: 1, Percent: 100,
		CoveredLines: map[int]bool{4: true}, // 200 branch line
	}
	oc := liveVerdict(cfg, sess, &ep, hurl, cov, sess.CurrentEntry())
	if oc != outcomePass {
		t.Fatalf("expected PASS, got %v", oc)
	}
	if sess.Entries[0].CRI != 3 {
		t.Fatalf("expected CRI 3 (COVERED), got %d", sess.Entries[0].CRI)
	}
}

// DONE requires an unreachable.yaml reason for the uncovered branch (C-04).
func TestLiveVerdict_DoneRequiresReason(t *testing.T) {
	dir := withSessionDir(t)
	src := filepath.Join(dir, "h.go")
	os.WriteFile(src, []byte(`package main
import "net/http"
func H(c interface{}) {
	c.JSON(http.StatusOK, nil)
	c.JSON(http.StatusServiceUnavailable, nil)
}
`), 0o644)
	hurl := filepath.Join(dir, "t.hurl")
	os.WriteFile(hurl, []byte("GET {{host}}/x\nHTTP 200\n"), 0o644)

	ep := scanner.Endpoint{ID: "ep1", Method: "GET", Path: "/x", Handler: "H", Source: src, Line: 3}
	sess := session.New()
	sess.Merge([]scanner.Endpoint{ep})
	cfg := &config.Config{Scan: config.ScanConfig{Lang: "go"}, Server: config.ServerConfig{Start: "x"}}

	// 503 is a server branch (>=500) → advisory, excluded from client gate.
	// Use a 404 to keep it client-side instead.
	os.WriteFile(src, []byte(`package main
import "net/http"
func H(c interface{}) {
	c.JSON(http.StatusOK, nil)
	c.JSON(http.StatusNotFound, nil)
}
`), 0o644)

	// Stalled: improveCount>=1 and percent <= prev. Only line 4 (200) hit.
	cov := &adapter.CoverageResult{Covered: 1, Total: 2, Percent: 50, CoveredLines: map[int]bool{4: true}}
	sess.MarkImprove("ep1", 50)
	sess.MarkImprove("ep1", 50)

	// Without unreachable.yaml → UNVERIFIED (not DONE).
	oc := liveVerdict(cfg, sess, &ep, hurl, cov, sess.CurrentEntry())
	if oc != outcomeUnverified {
		t.Fatalf("expected UNVERIFIED without reason, got %v", oc)
	}

	// With a valid reason for the uncovered 404 branch → DONE.
	os.MkdirAll(".huma", 0o755)
	os.WriteFile(filepath.Join(".huma", "unreachable.yaml"), []byte(
		"- endpoint: GET /x\n  status: 404\n  reason: cannot trigger in unit env\n  evidence: h.go:5\n"), 0o644)
	sess.MarkImprove("ep1", 50)
	sess.MarkImprove("ep1", 50)
	oc = liveVerdict(cfg, sess, &ep, hurl, cov, sess.CurrentEntry())
	if oc != outcomeDone {
		t.Fatalf("expected DONE with reason, got %v", oc)
	}
}

var _ = analyzer.ResponseBranch{}
