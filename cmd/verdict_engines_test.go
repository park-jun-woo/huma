package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/park-jun-woo/huma/internal/adapter"
	"github.com/park-jun-woo/huma/internal/analyzer"
	"github.com/park-jun-woo/huma/internal/config"
	"github.com/park-jun-woo/huma/internal/scanner"
	"github.com/park-jun-woo/huma/internal/session"
)

// writeGoHandler writes a handler source returning the supplied statuses on
// consecutive lines starting at line 4, and returns its path.
func writeGoHandler(t *testing.T, dir string, statuses ...string) string {
	t.Helper()
	body := "package main\nimport \"net/http\"\nfunc H(c interface{}) {\n"
	for _, s := range statuses {
		body += "\tc.JSON(http." + s + ", nil)\n"
	}
	body += "}\n"
	p := filepath.Join(dir, "h.go")
	if err := os.WriteFile(p, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestSmokeVerdict(t *testing.T) {
	dir := withSessionDir(t)
	src := writeGoHandler(t, dir, "StatusOK")
	ep := scanner.Endpoint{ID: "ep1", Method: "GET", Path: "/x", Handler: "H", Source: src, Line: 3}
	sess := session.New()
	sess.Merge([]scanner.Endpoint{ep})

	// no gate → SMOKE PASS at CRI 2
	oc := smokeVerdict(&config.Config{}, sess, &ep, 1, "source")
	if oc != outcomePass {
		t.Fatalf("expected PASS, got %v", oc)
	}
	if sess.Entries[0].CRI != 2 {
		t.Fatalf("expected CRI 2, got %d", sess.Entries[0].CRI)
	}

	// require_cri 3 > reachable 2 → UNVERIFIED
	cri := 3
	oc = smokeVerdict(&config.Config{RequireCRI: &cri}, sess, &ep, 1, "source")
	if oc != outcomeUnverified {
		t.Fatalf("expected UNVERIFIED with gate, got %v", oc)
	}
}

func TestCoveredVerdict_Pass(t *testing.T) {
	dir := withSessionDir(t)
	src := writeGoHandler(t, dir, "StatusOK")
	ep := scanner.Endpoint{ID: "ep1", Method: "GET", Path: "/x", Handler: "H", Source: src, Line: 3}
	sess := session.New()
	sess.Merge([]scanner.Endpoint{ep})
	branches := []analyzer.ResponseBranch{{Status: 200, Line: 4}}
	cov := &adapter.CoverageResult{Percent: 100, Total: 1, Covered: 1, CoveredLines: map[int]bool{4: true}}

	oc := coveredVerdict(&config.Config{}, sess, &ep, branches, cov, sess.CurrentEntry(), 2, "source")
	if oc != outcomePass {
		t.Fatalf("expected PASS, got %v", oc)
	}
	if sess.Entries[0].CRI != 3 {
		t.Fatalf("expected CRI 3, got %d", sess.Entries[0].CRI)
	}
}

func TestCoveredVerdict_GatedUnverified(t *testing.T) {
	dir := withSessionDir(t)
	src := writeGoHandler(t, dir, "StatusOK")
	ep := scanner.Endpoint{ID: "ep1", Method: "GET", Path: "/x", Handler: "H", Source: src, Line: 3}
	sess := session.New()
	sess.Merge([]scanner.Endpoint{ep})
	branches := []analyzer.ResponseBranch{{Status: 200, Line: 4}}

	// First sub-100% attempt routes to IMPROVE.
	cov2 := &adapter.CoverageResult{Percent: 50, Total: 2, Covered: 1, CoveredLines: map[int]bool{4: true}}
	oc := coveredVerdict(&config.Config{}, sess, &ep, branches, cov2, sess.CurrentEntry(), 2, "source")
	if oc != outcomeImprove {
		t.Fatalf("expected IMPROVE on first sub-100%%, got %v", oc)
	}
}

func TestCoveredVerdict_Stalled(t *testing.T) {
	dir := withSessionDir(t)
	src := writeGoHandler(t, dir, "StatusOK", "StatusNotFound")
	ep := scanner.Endpoint{ID: "ep1", Method: "GET", Path: "/x", Handler: "H", Source: src, Line: 3}
	sess := session.New()
	sess.Merge([]scanner.Endpoint{ep})
	branches := []analyzer.ResponseBranch{{Status: 200, Line: 4}, {Status: 404, Line: 5}}
	cov := &adapter.CoverageResult{Percent: 50, Total: 2, Covered: 1, CoveredLines: map[int]bool{4: true}}

	// stall it
	sess.MarkImprove("ep1", 50)
	sess.MarkImprove("ep1", 50)

	// no unreachable.yaml → stalledVerdict → UNVERIFIED
	oc := coveredVerdict(&config.Config{}, sess, &ep, branches, cov, sess.CurrentEntry(), 2, "source")
	if oc != outcomeUnverified {
		t.Fatalf("expected UNVERIFIED (stalled, no reason), got %v", oc)
	}
}

func TestStalledVerdict_DoneWithReason(t *testing.T) {
	dir := withSessionDir(t)
	src := writeGoHandler(t, dir, "StatusOK", "StatusNotFound")
	ep := scanner.Endpoint{ID: "ep1", Method: "GET", Path: "/x", Handler: "H", Source: src, Line: 3}
	sess := session.New()
	sess.Merge([]scanner.Endpoint{ep})
	branches := []analyzer.ResponseBranch{{Status: 200, Line: 4}, {Status: 404, Line: 5}}
	cov := &adapter.CoverageResult{Percent: 50, Total: 2, Covered: 1, CoveredLines: map[int]bool{4: true}}

	os.MkdirAll(filepath.Join(dir, ".huma"), 0o755)
	os.WriteFile(filepath.Join(dir, ".huma", "unreachable.yaml"), []byte(
		"- endpoint: GET /x\n  status: 404\n  reason: cannot trigger\n  evidence: h.go:5\n"), 0o644)

	oc := stalledVerdict(sess, &ep, branches, cov, 2, "source")
	if oc != outcomeDone {
		t.Fatalf("expected DONE with reason, got %v", oc)
	}
	if sess.Entries[0].Status != session.StatusDone {
		t.Fatalf("expected status DONE, got %s", sess.Entries[0].Status)
	}
}

func TestDoneReasonsSatisfied(t *testing.T) {
	dir := withSessionDir(t)
	ep := &scanner.Endpoint{Method: "GET", Path: "/x"}
	cov := &adapter.CoverageResult{CoveredLines: map[int]bool{4: true}}

	// nothing uncovered → satisfied
	covered := []analyzer.ResponseBranch{{Status: 200, Line: 4}}
	if !doneReasonsSatisfied(ep, covered, cov) {
		t.Error("all covered → satisfied")
	}

	// uncovered, no unreachable.yaml → not satisfied
	branches := []analyzer.ResponseBranch{{Status: 200, Line: 4}, {Status: 404, Line: 5}}
	if doneReasonsSatisfied(ep, branches, cov) {
		t.Error("uncovered without reason → not satisfied")
	}

	// add reason → satisfied
	os.MkdirAll(filepath.Join(dir, ".huma"), 0o755)
	os.WriteFile(filepath.Join(dir, ".huma", "unreachable.yaml"), []byte(
		"- endpoint: GET /x\n  status: 404\n  reason: r\n  evidence: e\n"), 0o644)
	if !doneReasonsSatisfied(ep, branches, cov) {
		t.Error("uncovered with reason → satisfied")
	}
}

func TestStaticAGrade(t *testing.T) {
	dir := t.TempDir()
	hurl := filepath.Join(dir, "t.hurl")
	os.WriteFile(hurl, []byte("GET {{host}}/x\nHTTP 200\n[Asserts]\njsonpath \"$.id\" exists\n"), 0o644)
	branches := []analyzer.ResponseBranch{{Status: 200}}
	g := staticAGrade(hurl, branches)
	if g < 1 {
		t.Errorf("expected positive A-grade, got %d", g)
	}
	// missing file → 0
	if g := staticAGrade(filepath.Join(dir, "nope.hurl"), branches); g != 0 {
		t.Errorf("missing file → 0, got %d", g)
	}
}
