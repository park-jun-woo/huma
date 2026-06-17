package humaquest

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/park-jun-woo/huma/internal/analyzer"
	"github.com/park-jun-woo/huma/internal/config"
	"github.com/park-jun-woo/huma/internal/hurlcheck"
	"github.com/park-jun-woo/huma/internal/scanner"
	"github.com/park-jun-woo/reins/pkg/gate"
	"github.com/park-jun-woo/reins/pkg/quest"
)

// ---------------------------------------------------------------------------
// unverifiedVerdict
// ---------------------------------------------------------------------------

func TestUnverifiedVerdict(t *testing.T) {
	v := unverifiedVerdict("GET /x", "an oracle", "no oracle", "fix it")
	if v.Outcome != quest.OutReview {
		t.Errorf("Outcome = %q, want REVIEW", v.Outcome)
	}
	if v.RootCause != "C-01" {
		t.Errorf("RootCause = %q, want C-01", v.RootCause)
	}
	if len(v.Facts) != 1 || v.Facts[0].Rule != "C-01" || v.Facts[0].Where != "GET /x" {
		t.Errorf("Facts wrong: %+v", v.Facts)
	}
	if v.Facts[0].Expected != "an oracle" || v.Facts[0].Actual != "no oracle" {
		t.Errorf("Fact content wrong: %+v", v.Facts[0])
	}
	if v.Feedback != "fix it" {
		t.Errorf("Feedback = %q", v.Feedback)
	}
}

// ---------------------------------------------------------------------------
// passVerdict
// ---------------------------------------------------------------------------

func TestPassVerdict(t *testing.T) {
	branches := []analyzer.ResponseBranch{br(200, 1), br(404, 2)}
	v := passVerdict(3, 2, "both", branches)
	if v.Outcome != quest.OutPass {
		t.Errorf("Outcome = %q, want PASS", v.Outcome)
	}
	if len(v.Facts) != 0 {
		t.Errorf("PASS carries no Facts, got %+v", v.Facts)
	}
	for _, want := range []string{"COVERED", "CRI 3", "client branches 2", "A=2", "both"} {
		if !strings.Contains(v.Feedback, want) {
			t.Errorf("Feedback %q missing %q", v.Feedback, want)
		}
	}
}

// ---------------------------------------------------------------------------
// improveFeedback
// ---------------------------------------------------------------------------

func TestImproveFeedback(t *testing.T) {
	uncovered := []analyzer.ResponseBranch{br(404, 88)}

	t.Run("no prior signal", func(t *testing.T) {
		got := improveFeedback(uncovered, 50, 0)
		if !strings.Contains(got, "coverage 50%") || strings.Contains(got, "stalled") {
			t.Errorf("got %q", got)
		}
		if !strings.Contains(got, "404@L88") {
			t.Errorf("missing uncovered list: %q", got)
		}
	})

	t.Run("stalled (not rising)", func(t *testing.T) {
		got := improveFeedback(uncovered, 50, 50)
		if !strings.Contains(got, "stalled at 50% (was 50%)") {
			t.Errorf("want stalled note, got %q", got)
		}
	})

	t.Run("improving", func(t *testing.T) {
		got := improveFeedback(uncovered, 75, 50)
		if !strings.Contains(got, "up from 50%") {
			t.Errorf("want up-from note, got %q", got)
		}
	})
}

// ---------------------------------------------------------------------------
// boundaryNoReason
// ---------------------------------------------------------------------------

func TestBoundaryNoReason(t *testing.T) {
	ep := &scanner.Endpoint{Method: "GET", Path: "/x"}
	uncovered := []analyzer.ResponseBranch{br(404, 1)}

	t.Run("nil item → false", func(t *testing.T) {
		chdir(t, t.TempDir())
		if boundaryNoReason(gate.Context{}, ep, uncovered) {
			t.Error("nil item is not the boundary")
		}
	})

	t.Run("before boundary → false", func(t *testing.T) {
		chdir(t, t.TempDir())
		ctx := gate.Context{Item: &quest.Item{Tries: 0}}
		if boundaryNoReason(ctx, ep, uncovered) {
			t.Error("Tries 0 is not the last attempt")
		}
	})

	t.Run("at boundary no reason → true", func(t *testing.T) {
		chdir(t, t.TempDir())
		ctx := gate.Context{Item: &quest.Item{Tries: quest.MaxTries - 1}}
		if !boundaryNoReason(ctx, ep, uncovered) {
			t.Error("last attempt without reason → true")
		}
	})

	t.Run("at boundary with reason → false", func(t *testing.T) {
		dir := t.TempDir()
		chdir(t, dir)
		writeUnreachable(t, dir, "- endpoint: GET /x\n  status: 404\n  reason: dead\n  evidence: h.go:1\n")
		ctx := gate.Context{Item: &quest.Item{Tries: quest.MaxTries - 1}}
		if boundaryNoReason(ctx, ep, uncovered) {
			t.Error("reason-backed boundary → false")
		}
	})
}

// ---------------------------------------------------------------------------
// boundaryReviewVerdict
// ---------------------------------------------------------------------------

func TestBoundaryReviewVerdict(t *testing.T) {
	ep := &scanner.Endpoint{Method: "GET", Path: "/x"}
	uncovered := []analyzer.ResponseBranch{br(404, 88)}
	v := boundaryReviewVerdict(ep, uncovered, 50)
	if v.Outcome != quest.OutReview {
		t.Errorf("Outcome = %q, want REVIEW", v.Outcome)
	}
	if v.RootCause != "C-04" {
		t.Errorf("RootCause = %q, want C-04", v.RootCause)
	}
	if len(v.Facts) != 1 || v.Facts[0].Rule != "C-04" || v.Facts[0].Where != "GET /x" {
		t.Errorf("Facts wrong: %+v", v.Facts)
	}
	if !strings.Contains(v.Feedback, "404@L88") || !strings.Contains(v.Feedback, "retry limit") {
		t.Errorf("Feedback = %q", v.Feedback)
	}
}

// ---------------------------------------------------------------------------
// improveVerdict
// ---------------------------------------------------------------------------

func TestImproveVerdict(t *testing.T) {
	ep := &scanner.Endpoint{Method: "GET", Path: "/x"}
	branches := []analyzer.ResponseBranch{br(200, 10), br(404, 20)}
	uncovered := []analyzer.ResponseBranch{br(404, 20)}

	t.Run("retry → OutFail C-03", func(t *testing.T) {
		chdir(t, t.TempDir())
		ctx := gate.Context{Item: &quest.Item{Tries: 0}}
		v := improveVerdict(ctx, ep, branches, uncovered, 50, 0)
		if v.Outcome != quest.OutFail {
			t.Errorf("Outcome = %q, want FAIL", v.Outcome)
		}
		if v.RootCause != "C-03" {
			t.Errorf("RootCause = %q, want C-03", v.RootCause)
		}
		if len(v.Facts) != 1 || v.Facts[0].Rule != "C-03" {
			t.Errorf("Facts wrong: %+v", v.Facts)
		}
		if !strings.Contains(v.Facts[0].Expected, "all 2 client branch") {
			t.Errorf("Expected text = %q", v.Facts[0].Expected)
		}
		if !strings.Contains(v.Facts[0].Actual, "404@L20") {
			t.Errorf("Actual text = %q", v.Facts[0].Actual)
		}
	})

	t.Run("boundary no reason → OutReview C-04", func(t *testing.T) {
		chdir(t, t.TempDir())
		ctx := gate.Context{Item: &quest.Item{Tries: quest.MaxTries - 1}}
		v := improveVerdict(ctx, ep, branches, uncovered, 50, 0)
		if v.Outcome != quest.OutReview {
			t.Errorf("Outcome = %q, want REVIEW at boundary", v.Outcome)
		}
		if v.RootCause != "C-04" {
			t.Errorf("RootCause = %q, want C-04", v.RootCause)
		}
	})
}

// ---------------------------------------------------------------------------
// assertionImproveVerdict / assertionImproveFeedback
// ---------------------------------------------------------------------------

func TestAssertionImproveVerdict(t *testing.T) {
	ep := &scanner.Endpoint{Method: "GET", Path: "/x", Source: "h.go"}
	branches := []analyzer.ResponseBranch{br(200, 10), br(404, 20)}

	t.Run("retry → OutFail C-03 with cap transparency", func(t *testing.T) {
		chdir(t, t.TempDir())
		ctx := gate.Context{Item: &quest.Item{Tries: 0}}
		v := assertionImproveVerdict(ctx, ep, branches, 3, 1)
		if v.Outcome != quest.OutFail || v.RootCause != "C-03" {
			t.Errorf("want OutFail/C-03, got %+v", v)
		}
		if len(v.Facts) != 1 || v.Facts[0].Rule != "C-03" {
			t.Errorf("Facts wrong: %+v", v.Facts)
		}
		if !strings.Contains(v.Facts[0].Actual, "A=1 caps CRI to 1 (staged 3)") {
			t.Errorf("cap transparency missing in Fact: %q", v.Facts[0].Actual)
		}
		if !strings.Contains(v.Feedback, "A=1 caps CRI to 1 (staged 3)") {
			t.Errorf("cap transparency missing in Feedback: %q", v.Feedback)
		}
	})

	t.Run("boundary no reason → OutReview C-04", func(t *testing.T) {
		chdir(t, t.TempDir())
		ctx := gate.Context{Item: &quest.Item{Tries: quest.MaxTries - 1}}
		v := assertionImproveVerdict(ctx, ep, branches, 3, 1)
		if v.Outcome != quest.OutReview || v.RootCause != "C-04" {
			t.Errorf("boundary → OutReview/C-04, got %+v", v)
		}
	})
}

func TestAssertionImproveFeedback(t *testing.T) {
	t.Run("no response fields → generic hint", func(t *testing.T) {
		ep := &scanner.Endpoint{Method: "GET", Path: "/x"}
		got := assertionImproveFeedback(ep, 3, 1)
		if !strings.Contains(got, "A=1 caps CRI to 1 (staged 3)") {
			t.Errorf("missing cap line: %q", got)
		}
		if strings.Contains(got, "toward A=3:") {
			t.Errorf("no fields → should not list targets: %q", got)
		}
	})

	t.Run("response fields → lists up to three targets", func(t *testing.T) {
		ep := &scanner.Endpoint{Method: "GET", Path: "/x", ResponseFields: []scanner.ResponseField{
			{Path: "id", Type: "int"},
			{Path: "name", Type: "string"},
			{Path: "email", Type: "string"},
			{Path: "extra", Type: "string"},
		}}
		got := assertionImproveFeedback(ep, 3, 1)
		if !strings.Contains(got, "id, name, email") || strings.Contains(got, "extra") {
			t.Errorf("want first three field paths, got %q", got)
		}
		if !strings.Contains(got, "assert response body fields toward A=3:") {
			t.Errorf("missing field-listing prefix: %q", got)
		}
		if !strings.HasSuffix(got, ").") {
			t.Errorf("expected close paren + period: %q", got)
		}
	})

	t.Run("skips fields with empty path", func(t *testing.T) {
		ep := &scanner.Endpoint{Method: "GET", Path: "/x", ResponseFields: []scanner.ResponseField{
			{Path: "", Type: "string"},
			{Path: "token", Type: "string"},
		}}
		got := assertionImproveFeedback(ep, 4, 2)
		if !strings.Contains(got, "toward A=3: token)") {
			t.Errorf("want only non-empty path listed, got %q", got)
		}
	})

	t.Run("all empty paths → falls back to generic hint", func(t *testing.T) {
		ep := &scanner.Endpoint{Method: "GET", Path: "/x", ResponseFields: []scanner.ResponseField{
			{Path: "", Type: "string"},
		}}
		got := assertionImproveFeedback(ep, 4, 2)
		if strings.Contains(got, "toward A=3:") {
			t.Errorf("all paths empty → should fall back to generic hint: %q", got)
		}
		if !strings.HasSuffix(got, ".") {
			t.Errorf("expected trailing period: %q", got)
		}
	})
}

// ---------------------------------------------------------------------------
// gateVerdict
// ---------------------------------------------------------------------------

func TestGateVerdict(t *testing.T) {
	ep := &scanner.Endpoint{Method: "GET", Path: "/x"}
	branches := []analyzer.ResponseBranch{br(200, 1)}

	t.Run("tier meets gate → PASS", func(t *testing.T) {
		v := gateVerdict(&config.Config{}, 2, ep, branches, 3, "source")
		if v.Outcome != quest.OutPass {
			t.Errorf("Outcome = %q, want PASS", v.Outcome)
		}
		if !strings.Contains(v.Feedback, "SMOKE") {
			t.Errorf("Feedback = %q", v.Feedback)
		}
	})

	t.Run("tier below explicit require_cri → UNVERIFIED", func(t *testing.T) {
		three := 3
		v := gateVerdict(&config.Config{RequireCRI: &three}, 2, ep, branches, 3, "source")
		if v.Outcome != quest.OutReview {
			t.Errorf("Outcome = %q, want REVIEW", v.Outcome)
		}
		if v.RootCause != "C-01" {
			t.Errorf("RootCause = %q, want C-01", v.RootCause)
		}
		if !strings.Contains(v.Feedback, "require_cri=3") {
			t.Errorf("Feedback = %q", v.Feedback)
		}
	})
}

// ---------------------------------------------------------------------------
// verdictFromCRI
// ---------------------------------------------------------------------------

func TestVerdictFromCRI(t *testing.T) {
	ep := &scanner.Endpoint{Method: "GET", Path: "/x", Source: "h.go"}
	branches := []analyzer.ResponseBranch{br(200, 10), br(404, 20)}
	cfg := &config.Config{}

	t.Run("tier 0 → UNVERIFIED", func(t *testing.T) {
		chdir(t, t.TempDir())
		sub := &hurlInfo{}
		ctx := gate.Context{Item: &quest.Item{}}
		v := verdictFromCRI(ctx, cfg, 0, ep, sub, branches, nil, "source", 0, 0)
		if v.Outcome != quest.OutReview || v.RootCause != "C-01" {
			t.Errorf("tier 0 → UNVERIFIED, got %+v", v)
		}
		if !strings.Contains(v.Feedback, "no independent oracle") {
			t.Errorf("Feedback = %q", v.Feedback)
		}
	})

	t.Run("tier 3 + A=3 → PASS via gate", func(t *testing.T) {
		chdir(t, t.TempDir())
		sub := &hurlInfo{}
		ctx := gate.Context{Item: &quest.Item{}}
		// A=3 == staged 3 → A is not the limiting axis → COVERED PASS (no regression).
		v := verdictFromCRI(ctx, cfg, 3, ep, sub, branches, covWith(2, 100, 10, 20), "source", 3, 0)
		if v.Outcome != quest.OutPass {
			t.Errorf("tier 3 → PASS, got %+v", v)
		}
	})

	t.Run("tier 3 + A=1 → assertion IMPROVE (CV-4)", func(t *testing.T) {
		chdir(t, t.TempDir())
		sub := &hurlInfo{}
		ctx := gate.Context{Item: &quest.Item{Tries: 0}}
		// staged 3 (100% + bound) but A=1 (status only) → A-branch preempts → C-03.
		v := verdictFromCRI(ctx, cfg, 3, ep, sub, branches, covWith(2, 100, 10, 20), "source", 1, 0)
		if v.Outcome != quest.OutFail || v.RootCause != "C-03" {
			t.Errorf("A=1 caps COVERED → assertion IMPROVE, got %+v", v)
		}
		if !strings.Contains(v.Feedback, "A=1 caps CRI to 1 (staged 3)") {
			t.Errorf("cap transparency missing: %q", v.Feedback)
		}
	})

	t.Run("cri>0 + A=0 with oracle → assertion IMPROVE not UNVERIFIED", func(t *testing.T) {
		chdir(t, t.TempDir())
		sub := &hurlInfo{}
		ctx := gate.Context{Item: &quest.Item{Tries: 0}}
		// A=0 but staged>0 (oracle present) → C-03 IMPROVE, NOT C-01/UNVERIFIED.
		v := verdictFromCRI(ctx, cfg, 2, ep, sub, branches, covWith(0, 0), "source", 0, 0)
		if v.Outcome != quest.OutFail || v.RootCause != "C-03" {
			t.Errorf("A=0 with oracle → C-03 IMPROVE, got %+v", v)
		}
	})

	t.Run("tier 2 instrumented incomplete → IMPROVE", func(t *testing.T) {
		chdir(t, t.TempDir())
		sub := &hurlInfo{}
		ctx := gate.Context{Item: &quest.Item{Tries: 0}}
		cov := covWith(2, 50, 10) // 404 uncovered
		// A=2 == staged 2 → not the limiting axis → exercises the coverage IMPROVE path.
		v := verdictFromCRI(ctx, cfg, 2, ep, sub, branches, cov, "source", 2, 0)
		if v.Outcome != quest.OutFail || v.RootCause != "C-03" {
			t.Errorf("tier 2 instrumented → IMPROVE, got %+v", v)
		}
	})

	t.Run("tier 2 uninstrumented green → SMOKE gate PASS", func(t *testing.T) {
		chdir(t, t.TempDir())
		sub := &hurlInfo{}
		ctx := gate.Context{Item: &quest.Item{}}
		// cov.Total == 0 → uninstrumented path → gateVerdict(2).
		v := verdictFromCRI(ctx, cfg, 2, ep, sub, branches, covWith(0, 0), "source", 2, 0)
		if v.Outcome != quest.OutPass || !strings.Contains(v.Feedback, "SMOKE") {
			t.Errorf("tier 2 uninstrumented → SMOKE PASS, got %+v", v)
		}
	})

	t.Run("tier 1 static all covered → SCAFFOLDED gate PASS", func(t *testing.T) {
		chdir(t, t.TempDir())
		// Entries assert both 200 and 404 with non-vacuous grades → respResult
		// Missing is empty → gateVerdict(1).
		sub := &hurlInfo{Entries: []hurlcheck.HurlEntry{
			{Status: 200, Grade: 3},
			{Status: 404, Grade: 3},
		}}
		ctx := gate.Context{Item: &quest.Item{}}
		v := verdictFromCRI(ctx, cfg, 1, ep, sub, branches, nil, "source", 3, 0)
		if v.Outcome != quest.OutPass || !strings.Contains(v.Feedback, "SCAFFOLDED") {
			t.Errorf("tier 1 covered → SCAFFOLDED PASS, got %+v", v)
		}
	})

	t.Run("tier 1 static missing statuses → IMPROVE", func(t *testing.T) {
		chdir(t, t.TempDir())
		// Only 200 asserted → 404 missing → improveVerdict.
		sub := &hurlInfo{Entries: []hurlcheck.HurlEntry{{Status: 200, Grade: 3}}}
		ctx := gate.Context{Item: &quest.Item{Tries: 0}}
		// A=1 == staged 1 → not the limiting axis → exercises the static missing-status IMPROVE path.
		v := verdictFromCRI(ctx, cfg, 1, ep, sub, branches, nil, "source", 1, 0)
		if v.Outcome != quest.OutFail || v.RootCause != "C-03" {
			t.Errorf("tier 1 missing → IMPROVE, got %+v", v)
		}
	})
}

// ---------------------------------------------------------------------------
// Evaluate (integration)
// ---------------------------------------------------------------------------

// goManifest is a minimal manifest.yaml so config.Load resolves Scan.Lang to go
// (without it Evaluate falls back to an empty lang and the analyzer is nil).
const goManifest = `apiVersion: yongol/v1
kind: Project
metadata:
  name: eval-test
backend:
  lang: go
  framework: gin
  module: github.com/test/test
testing:
  base_url: http://localhost:8080
  hurl_dir: hurl
`

// writeManifest writes goManifest into dir.
func writeManifest(t *testing.T, dir string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(dir, "manifest.yaml"), []byte(goManifest), 0o644); err != nil {
		t.Fatal(err)
	}
}

// evalItem builds an Item carrying the endpoint payload (with PrevCoverage).
func evalItem(t *testing.T, ep scanner.Endpoint, tries int, prev float64) *quest.Item {
	t.Helper()
	it := &quest.Item{Key: ep.ID, State: quest.TODO, Tries: tries}
	ps := payloadState{Endpoint: ep, PrevCoverage: prev}
	if err := it.SetPayload(&ps); err != nil {
		t.Fatalf("SetPayload: %v", err)
	}
	return it
}

func TestEvaluate_BadSubmission(t *testing.T) {
	chdir(t, t.TempDir())
	v := (humaDef{}).Evaluate(gate.Context{Submission: "not a hurlInfo"})
	if v.Outcome != quest.OutReview || v.RootCause != "C-01" {
		t.Errorf("bad submission → UNVERIFIED, got %+v", v)
	}
}

func TestEvaluate_NilSubmission(t *testing.T) {
	chdir(t, t.TempDir())
	var sub *hurlInfo
	v := (humaDef{}).Evaluate(gate.Context{Submission: sub})
	if v.Outcome != quest.OutReview {
		t.Errorf("nil submission → UNVERIFIED, got %+v", v)
	}
}

func TestEvaluate_NoOracleUnverified(t *testing.T) {
	chdir(t, t.TempDir())
	// No source, no coverage → no oracle → tier 0 → UNVERIFIED.
	ep := scanner.Endpoint{Method: "GET", Path: "/x", Responses: json.RawMessage(`[{"status":200,"line":1}]`)}
	sub := &hurlInfo{Endpoint: ep}
	ctx := gate.Context{Submission: sub, Item: evalItem(t, ep, 0, 0)}
	v := (humaDef{}).Evaluate(ctx)
	if v.Outcome != quest.OutReview || v.RootCause != "C-01" {
		t.Errorf("no oracle → UNVERIFIED, got %+v", v)
	}
}

func TestEvaluate_FullCoverageCovered(t *testing.T) {
	dir := t.TempDir()
	chdir(t, dir)
	writeManifest(t, dir)
	src := writeGoHandler(t) // 200/404/500 on distinct lines
	ep := scanner.Endpoint{Method: "GET", Path: "/x", Source: src, Handler: "CreateUser"}
	// A=3 on both client branches so the A axis does not cap COVERED (Phase 007).
	sub := &hurlInfo{Endpoint: ep, Entries: []hurlcheck.HurlEntry{{Status: 200, Grade: 3}, {Status: 404, Grade: 3}}}

	// Determine the source lines for the client branches (200, 404) so coverage
	// can mark them hit at 100%.
	client, _ := responseBranches(&ep, "go")
	lines := make([]int, 0, len(client))
	for _, b := range client {
		lines = append(lines, b.Line)
	}
	cov := covWith(len(lines), 100, lines...)
	covJSON, _ := json.Marshal(cov)

	ctx := gate.Context{
		Submission: sub,
		Item:       evalItem(t, ep, 0, 0),
		Grounds:    map[string]string{"coverage": string(covJSON)},
	}
	v := (humaDef{}).Evaluate(ctx)
	if v.Outcome != quest.OutPass {
		t.Fatalf("full coverage → PASS, got %+v", v)
	}
	if !strings.Contains(v.Feedback, "COVERED") {
		t.Errorf("Feedback = %q, want COVERED", v.Feedback)
	}
}

func TestEvaluate_FullCoverageShallowAssertionImprove(t *testing.T) {
	// CV-4: status-only assertions (A=1) with 100% line coverage must NOT be
	// COVERED — the A axis caps CRI and routes to the assertion-depth IMPROVE.
	dir := t.TempDir()
	chdir(t, dir)
	writeManifest(t, dir)
	src := writeGoHandler(t)
	ep := scanner.Endpoint{Method: "GET", Path: "/x", Source: src, Handler: "CreateUser"}
	// Status asserted (A=1) but no body shape/invariants on the client branches.
	sub := &hurlInfo{Endpoint: ep, Entries: []hurlcheck.HurlEntry{{Status: 200, Grade: 1}, {Status: 404, Grade: 1}}}

	client, _ := responseBranches(&ep, "go")
	lines := make([]int, 0, len(client))
	for _, b := range client {
		lines = append(lines, b.Line)
	}
	cov := covWith(len(lines), 100, lines...)
	covJSON, _ := json.Marshal(cov)

	ctx := gate.Context{
		Submission: sub,
		Item:       evalItem(t, ep, 0, 0),
		Grounds:    map[string]string{"coverage": string(covJSON)},
	}
	v := (humaDef{}).Evaluate(ctx)
	if v.Outcome != quest.OutFail || v.RootCause != "C-03" {
		t.Fatalf("A=1 + 100%% coverage → assertion IMPROVE(C-03), got %+v", v)
	}
	if !strings.Contains(v.Feedback, "caps CRI") {
		t.Errorf("cap transparency missing from feedback: %q", v.Feedback)
	}
}

func TestEvaluate_PartialCoverageImprove(t *testing.T) {
	dir := t.TempDir()
	chdir(t, dir)
	writeManifest(t, dir)
	src := writeGoHandler(t)
	ep := scanner.Endpoint{Method: "GET", Path: "/x", Source: src, Handler: "CreateUser"}
	// A=3 so the coverage IMPROVE path is exercised, not the A-branch (Phase 007).
	sub := &hurlInfo{Endpoint: ep, Entries: []hurlcheck.HurlEntry{{Status: 200, Grade: 3}, {Status: 404, Grade: 3}}}

	// Cover only the first client branch line; leave another uncovered → <100%.
	client, _ := responseBranches(&ep, "go")
	cov := covWith(len(client), 50, client[0].Line)
	covJSON, _ := json.Marshal(cov)

	ctx := gate.Context{
		Submission: sub,
		Item:       evalItem(t, ep, 0, 0),
		Grounds:    map[string]string{"coverage": string(covJSON)},
	}
	v := (humaDef{}).Evaluate(ctx)
	if v.Outcome != quest.OutFail || v.RootCause != "C-03" {
		t.Errorf("partial coverage → IMPROVE, got %+v", v)
	}
}

func TestEvaluate_LastTryNoReasonReview(t *testing.T) {
	dir := t.TempDir()
	chdir(t, dir)
	writeManifest(t, dir)
	src := writeGoHandler(t)
	ep := scanner.Endpoint{Method: "GET", Path: "/x", Source: src, Handler: "CreateUser"}
	// A=3 so the coverage boundary path is exercised, not the A-branch (Phase 007).
	sub := &hurlInfo{Endpoint: ep, Entries: []hurlcheck.HurlEntry{{Status: 200, Grade: 3}, {Status: 404, Grade: 3}}}

	client, _ := responseBranches(&ep, "go")
	cov := covWith(len(client), 50, client[0].Line) // partial, uncovered remain
	covJSON, _ := json.Marshal(cov)

	ctx := gate.Context{
		Submission: sub,
		Item:       evalItem(t, ep, quest.MaxTries-1, 40), // last attempt, no unreachable.yaml
		Grounds:    map[string]string{"coverage": string(covJSON)},
	}
	v := (humaDef{}).Evaluate(ctx)
	if v.Outcome != quest.OutReview || v.RootCause != "C-04" {
		t.Errorf("last try no reason → UNVERIFIED(C-04), got %+v", v)
	}
}

// ---------------------------------------------------------------------------
// priorCoverage
// ---------------------------------------------------------------------------

func TestPriorCoverage(t *testing.T) {
	t.Run("nil item → 0", func(t *testing.T) {
		if got := priorCoverage(nil); got != 0 {
			t.Errorf("nil item want 0, got %v", got)
		}
	})

	t.Run("decode error → 0", func(t *testing.T) {
		it := &quest.Item{Payload: json.RawMessage(`{bad`)}
		if got := priorCoverage(it); got != 0 {
			t.Errorf("decode error want 0, got %v", got)
		}
	})

	t.Run("reads PrevCoverage", func(t *testing.T) {
		ep := scanner.Endpoint{Method: "GET", Path: "/x"}
		it := evalItem(t, ep, 0, 73.5)
		if got := priorCoverage(it); got != 73.5 {
			t.Errorf("want 73.5, got %v", got)
		}
	})

	t.Run("bare endpoint payload → 0 neutral", func(t *testing.T) {
		ep := scanner.Endpoint{Method: "GET", Path: "/x"}
		it := &quest.Item{}
		if err := it.SetPayload(&ep); err != nil {
			t.Fatal(err)
		}
		if got := priorCoverage(it); got != 0 {
			t.Errorf("bare endpoint → 0, got %v", got)
		}
	})
}
