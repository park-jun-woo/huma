package humaquest

import (
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"testing"

	"github.com/park-jun-woo/huma/internal/adapter"
	"github.com/park-jun-woo/huma/internal/analyzer"
	"github.com/park-jun-woo/huma/internal/config"
	"github.com/park-jun-woo/huma/internal/scanner"
)

// br is a compact constructor for a response branch with a status and source line.
func br(status, line int) analyzer.ResponseBranch {
	return analyzer.ResponseBranch{Status: status, Line: line}
}

// covWith builds a CoverageResult whose CoveredLines set marks the given lines hit.
func covWith(total int, pct float64, lines ...int) *adapter.CoverageResult {
	set := make(map[int]bool, len(lines))
	for _, l := range lines {
		set[l] = true
	}
	return &adapter.CoverageResult{Total: total, Percent: pct, CoveredLines: set}
}

// writeUnreachable writes a .huma/unreachable.yaml under the current dir.
func writeUnreachable(t *testing.T, dir, body string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Join(dir, ".huma"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, ".huma", "unreachable.yaml"), []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
}

// ---------------------------------------------------------------------------
// dedupByStatus
// ---------------------------------------------------------------------------

func TestDedupByStatus(t *testing.T) {
	in := []analyzer.ResponseBranch{br(200, 10), br(404, 20), br(200, 99), br(409, 30)}
	got := dedupByStatus(in)
	if len(got) != 3 {
		t.Fatalf("want 3 unique branches, got %d", len(got))
	}
	// First occurrence (source-preferred Line) wins for 200.
	if got[0].Status != 200 || got[0].Line != 10 {
		t.Errorf("first 200 not kept: %+v", got[0])
	}
	statuses := []int{got[0].Status, got[1].Status, got[2].Status}
	if !reflect.DeepEqual(statuses, []int{200, 404, 409}) {
		t.Errorf("order/statuses = %v", statuses)
	}
}

func TestDedupByStatus_Empty(t *testing.T) {
	if got := dedupByStatus(nil); len(got) != 0 {
		t.Errorf("want empty, got %v", got)
	}
}

// ---------------------------------------------------------------------------
// concatBranches
// ---------------------------------------------------------------------------

func TestConcatBranches(t *testing.T) {
	src := []analyzer.ResponseBranch{br(200, 1)}
	decl := []analyzer.ResponseBranch{br(404, 2), br(500, 3)}
	got := concatBranches(src, decl)
	if len(got) != 3 {
		t.Fatalf("want 3, got %d", len(got))
	}
	// Source first so source Line wins a later dedup.
	if got[0].Status != 200 || got[1].Status != 404 || got[2].Status != 500 {
		t.Errorf("concat order wrong: %+v", got)
	}
}

// ---------------------------------------------------------------------------
// splitClientBranches
// ---------------------------------------------------------------------------

func TestSplitClientBranches(t *testing.T) {
	in := []analyzer.ResponseBranch{br(200, 1), br(500, 2), br(404, 3), br(503, 4)}
	client, advisory := splitClientBranches(in)
	if len(client) != 2 || client[0].Status != 200 || client[1].Status != 404 {
		t.Errorf("client set wrong: %+v", client)
	}
	if len(advisory) != 2 || advisory[0].Status != 500 || advisory[1].Status != 503 {
		t.Errorf("advisory set wrong: %+v", advisory)
	}
}

func TestSplitClientBranches_Empty(t *testing.T) {
	client, advisory := splitClientBranches(nil)
	if client != nil || advisory != nil {
		t.Errorf("empty input should yield nil slices, got %v / %v", client, advisory)
	}
}

// ---------------------------------------------------------------------------
// branchStatuses
// ---------------------------------------------------------------------------

func TestBranchStatuses(t *testing.T) {
	got := branchStatuses([]analyzer.ResponseBranch{br(201, 1), br(409, 2)})
	if !reflect.DeepEqual(got, []int{201, 409}) {
		t.Errorf("got %v", got)
	}
	if got := branchStatuses(nil); len(got) != 0 {
		t.Errorf("empty want 0, got %v", got)
	}
}

// ---------------------------------------------------------------------------
// provenanceOf + BranchProvenance.String
// ---------------------------------------------------------------------------

func TestProvenanceOfAndString(t *testing.T) {
	some := []analyzer.ResponseBranch{br(200, 1)}
	tests := []struct {
		name string
		src  []analyzer.ResponseBranch
		decl []analyzer.ResponseBranch
		want string
	}{
		{"both", some, some, "both"},
		{"source only", some, nil, "source"},
		{"declared only", nil, some, "declared"},
		{"none", nil, nil, "none"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := provenanceOf(tt.src, tt.decl)
			if p.HasSource != (len(tt.src) > 0) || p.HasDeclared != (len(tt.decl) > 0) {
				t.Errorf("provenance flags wrong: %+v", p)
			}
			if got := p.String(); got != tt.want {
				t.Errorf("String() = %q, want %q", got, tt.want)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// criLabel
// ---------------------------------------------------------------------------

func TestCriLabel(t *testing.T) {
	tests := map[int]string{
		3:  "COVERED",
		2:  "SMOKE",
		1:  "SCAFFOLDED",
		0:  "UNVERIFIED",
		-1: "UNVERIFIED",
	}
	for tier, want := range tests {
		if got := criLabel(tier); got != want {
			t.Errorf("criLabel(%d) = %q, want %q", tier, got, want)
		}
	}
}

// ---------------------------------------------------------------------------
// formatBranches
// ---------------------------------------------------------------------------

func TestFormatBranches(t *testing.T) {
	if got := formatBranches(nil); got != "none" {
		t.Errorf("empty want \"none\", got %q", got)
	}
	got := formatBranches([]analyzer.ResponseBranch{br(404, 88), br(409, 0)})
	if got != "404@L88, 409@L0" {
		t.Errorf("formatBranches = %q", got)
	}
}

// ---------------------------------------------------------------------------
// lineCohesion
// ---------------------------------------------------------------------------

func TestLineCohesion(t *testing.T) {
	got := lineCohesion([]analyzer.ResponseBranch{br(200, 5), br(404, 5), br(409, 7)})
	if got[5] != 2 || got[7] != 1 {
		t.Errorf("cohesion map wrong: %v", got)
	}
	if got := lineCohesion(nil); len(got) != 0 {
		t.Errorf("empty want empty map, got %v", got)
	}
}

// ---------------------------------------------------------------------------
// uncoveredBranches
// ---------------------------------------------------------------------------

func TestUncoveredBranches(t *testing.T) {
	branches := []analyzer.ResponseBranch{br(200, 10), br(404, 20), br(409, 0)}
	cov := covWith(2, 50, 10) // line 10 hit, 20 not, 0 never bindable

	got := uncoveredBranches(branches, cov)
	// 200 covered → excluded; 404 (line not hit) and 409 (line 0) remain.
	if len(got) != 2 || got[0].Status != 404 || got[1].Status != 409 {
		t.Errorf("uncovered = %+v", got)
	}
}

func TestUncoveredBranches_NilCoverage(t *testing.T) {
	branches := []analyzer.ResponseBranch{br(200, 10)}
	got := uncoveredBranches(branches, nil)
	if len(got) != 1 {
		t.Errorf("nil cov → all uncovered, got %+v", got)
	}
}

// ---------------------------------------------------------------------------
// hasOracle
// ---------------------------------------------------------------------------

func TestHasOracle(t *testing.T) {
	tests := []struct {
		name string
		ep   scanner.Endpoint
		cov  *adapter.CoverageResult
		want bool
	}{
		{"source linked", scanner.Endpoint{Source: "h.go"}, nil, true},
		{"instrumented runtime", scanner.Endpoint{}, covWith(5, 100, 1), true},
		{"nil cov no source", scanner.Endpoint{}, nil, false},
		{"cov total zero no source", scanner.Endpoint{}, covWith(0, 0), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := hasOracle(&tt.ep, tt.cov); got != tt.want {
				t.Errorf("hasOracle = %v, want %v", got, tt.want)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// branchBound
// ---------------------------------------------------------------------------

func TestBranchBound(t *testing.T) {
	key := "GET /x"
	exemptions := []config.UnreachableEntry{{Endpoint: key, Status: 409, Reason: "r", Evidence: "e"}}
	cov := covWith(3, 100, 10)

	tests := []struct {
		name     string
		b        analyzer.ResponseBranch
		cohesion map[int]int
		want     bool
	}{
		{"exempt status", br(409, 0), map[int]int{}, true},
		{"line zero not exempt", br(404, 0), map[int]int{}, false},
		{"cohesive line", br(404, 10), map[int]int{10: 2}, false},
		{"unique line covered", br(200, 10), map[int]int{10: 1}, true},
		{"unique line uncovered", br(404, 20), map[int]int{20: 1}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := branchBound(tt.b, cov, exemptions, key, tt.cohesion); got != tt.want {
				t.Errorf("branchBound = %v, want %v", got, tt.want)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// branchBindingOK
// ---------------------------------------------------------------------------

func TestBranchBindingOK(t *testing.T) {
	ep := &scanner.Endpoint{Method: "GET", Path: "/x"}

	t.Run("no branches trivially bound", func(t *testing.T) {
		chdir(t, t.TempDir())
		if !branchBindingOK(ep, nil, covWith(0, 0)) {
			t.Error("empty branch set should be bound")
		}
	})

	t.Run("all bound", func(t *testing.T) {
		chdir(t, t.TempDir())
		branches := []analyzer.ResponseBranch{br(200, 10), br(404, 20)}
		cov := covWith(2, 100, 10, 20)
		if !branchBindingOK(ep, branches, cov) {
			t.Error("all unique covered lines should bind")
		}
	})

	t.Run("one unbound fails", func(t *testing.T) {
		chdir(t, t.TempDir())
		branches := []analyzer.ResponseBranch{br(200, 10), br(404, 20)}
		cov := covWith(2, 50, 10) // 20 not hit
		if branchBindingOK(ep, branches, cov) {
			t.Error("an uncovered branch line should fail binding")
		}
	})

	t.Run("unbound but exempt via unreachable.yaml", func(t *testing.T) {
		dir := t.TempDir()
		chdir(t, dir)
		writeUnreachable(t, dir, "- endpoint: GET /x\n  status: 404\n  reason: dead\n  evidence: h.go:1\n")
		branches := []analyzer.ResponseBranch{br(200, 10), br(404, 20)}
		cov := covWith(2, 50, 10) // 20 not hit, but exempt
		if !branchBindingOK(ep, branches, cov) {
			t.Error("exempt uncovered branch should still bind")
		}
	})
}

// ---------------------------------------------------------------------------
// allExempt
// ---------------------------------------------------------------------------

func TestAllExempt(t *testing.T) {
	key := "GET /x"
	exemptions := []config.UnreachableEntry{
		{Endpoint: key, Status: 200, Reason: "r", Evidence: "e"},
		{Endpoint: key, Status: 404, Reason: "r", Evidence: "e"},
	}

	if !allExempt(nil, exemptions, key) {
		t.Error("empty set is vacuously all-exempt")
	}
	if !allExempt([]analyzer.ResponseBranch{br(200, 1), br(404, 2)}, exemptions, key) {
		t.Error("both statuses exempt → true")
	}
	if allExempt([]analyzer.ResponseBranch{br(200, 1), br(409, 2)}, exemptions, key) {
		t.Error("409 not exempt → false")
	}
}

// ---------------------------------------------------------------------------
// reasonsCover
// ---------------------------------------------------------------------------

func TestReasonsCover(t *testing.T) {
	ep := &scanner.Endpoint{Method: "GET", Path: "/x"}

	t.Run("empty uncovered trivially satisfied", func(t *testing.T) {
		chdir(t, t.TempDir())
		if !reasonsCover(ep, nil) {
			t.Error("no uncovered branches → satisfied")
		}
	})

	t.Run("no unreachable file → false", func(t *testing.T) {
		chdir(t, t.TempDir())
		if reasonsCover(ep, []analyzer.ResponseBranch{br(404, 1)}) {
			t.Error("uncovered with no reasons → false")
		}
	})

	t.Run("malformed unreachable file → false", func(t *testing.T) {
		dir := t.TempDir()
		chdir(t, dir)
		writeUnreachable(t, dir, "::: not yaml :::\n  - broken")
		if reasonsCover(ep, []analyzer.ResponseBranch{br(404, 1)}) {
			t.Error("load error → false")
		}
	})

	t.Run("all uncovered exempt → true", func(t *testing.T) {
		dir := t.TempDir()
		chdir(t, dir)
		writeUnreachable(t, dir, "- endpoint: GET /x\n  status: 404\n  reason: dead\n  evidence: h.go:1\n")
		if !reasonsCover(ep, []analyzer.ResponseBranch{br(404, 1)}) {
			t.Error("exempt uncovered → true")
		}
	})

	t.Run("partially exempt → false", func(t *testing.T) {
		dir := t.TempDir()
		chdir(t, dir)
		writeUnreachable(t, dir, "- endpoint: GET /x\n  status: 404\n  reason: dead\n  evidence: h.go:1\n")
		if reasonsCover(ep, []analyzer.ResponseBranch{br(404, 1), br(409, 2)}) {
			t.Error("409 unexempted → false")
		}
	})
}

// ---------------------------------------------------------------------------
// doneReasonsSatisfied
// ---------------------------------------------------------------------------

func TestDoneReasonsSatisfied(t *testing.T) {
	ep := &scanner.Endpoint{Method: "GET", Path: "/x"}
	branches := []analyzer.ResponseBranch{br(200, 10), br(404, 20)}

	t.Run("all lines covered → satisfied", func(t *testing.T) {
		chdir(t, t.TempDir())
		cov := covWith(2, 100, 10, 20)
		if !doneReasonsSatisfied(ep, branches, cov) {
			t.Error("no uncovered branches → satisfied")
		}
	})

	t.Run("uncovered without reason → not satisfied", func(t *testing.T) {
		chdir(t, t.TempDir())
		cov := covWith(2, 50, 10) // 20 uncovered, no reasons
		if doneReasonsSatisfied(ep, branches, cov) {
			t.Error("uncovered no-reason → not satisfied")
		}
	})

	t.Run("uncovered with reason → satisfied", func(t *testing.T) {
		dir := t.TempDir()
		chdir(t, dir)
		writeUnreachable(t, dir, "- endpoint: GET /x\n  status: 404\n  reason: dead\n  evidence: h.go:1\n")
		cov := covWith(2, 50, 10)
		if !doneReasonsSatisfied(ep, branches, cov) {
			t.Error("exempt uncovered → satisfied")
		}
	})
}

// ---------------------------------------------------------------------------
// resolveGate
// ---------------------------------------------------------------------------

func TestResolveGate(t *testing.T) {
	if got := resolveGate(&config.Config{}, 2); got != 2 {
		t.Errorf("no require_cri → reachable max, got %d", got)
	}
	three := 3
	if got := resolveGate(&config.Config{RequireCRI: &three}, 1); got != 3 {
		t.Errorf("explicit require_cri wins, got %d", got)
	}
}

// ---------------------------------------------------------------------------
// decodeCoverage
// ---------------------------------------------------------------------------

func TestDecodeCoverage(t *testing.T) {
	if cov, present := decodeCoverage(""); cov != nil || present {
		t.Error("empty raw → (nil,false)")
	}
	if cov, present := decodeCoverage("{not json"); cov != nil || present {
		t.Error("malformed → (nil,false)")
	}
	cov, present := decodeCoverage(`{"Total":4,"Percent":75}`)
	if !present || cov == nil {
		t.Fatal("valid coverage → (cov,true)")
	}
	if cov.Total != 4 || cov.Percent != 75 {
		t.Errorf("decoded wrong: %+v", cov)
	}
}

// ---------------------------------------------------------------------------
// computeCRI
// ---------------------------------------------------------------------------

func TestComputeCRI(t *testing.T) {
	srcEP := &scanner.Endpoint{Source: "h.go", Method: "GET", Path: "/x"}
	noOracleEP := &scanner.Endpoint{Method: "GET", Path: "/x"}
	branches := []analyzer.ResponseBranch{br(200, 10)}

	tests := []struct {
		name       string
		ep         *scanner.Endpoint
		branches   []analyzer.ResponseBranch
		cov        *adapter.CoverageResult
		covPresent bool
		want       int
	}{
		{"no oracle → 0", noOracleEP, branches, nil, false, 0},
		{"no branches → 0", srcEP, nil, nil, false, 0},
		{"static only → 1", srcEP, branches, nil, false, 1},
		{"covPresent nil cov → 2", srcEP, branches, nil, true, 2},
		{"covPresent total zero → 2", srcEP, branches, covWith(0, 0), true, 2},
		{"100% bound → 3", srcEP, branches, covWith(1, 100, 10), true, 3},
		{"100% but unbound → 2", srcEP, []analyzer.ResponseBranch{br(404, 20)}, covWith(1, 100, 10), true, 2},
		{"below 100% → 2", srcEP, branches, covWith(2, 50, 10), true, 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			chdir(t, t.TempDir()) // branchBindingOK reads cwd unreachable.yaml
			if got := computeCRI(tt.ep, tt.branches, tt.cov, tt.covPresent); got != tt.want {
				t.Errorf("computeCRI = %d, want %d", got, tt.want)
			}
		})
	}
}

// sortedStatuses is a tiny utility some assertions use to compare unordered sets.
func sortedStatuses(b []analyzer.ResponseBranch) []int {
	out := branchStatuses(b)
	sort.Ints(out)
	return out
}
