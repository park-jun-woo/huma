package cmd

import (
	"reflect"
	"strings"
	"testing"

	"github.com/park-jun-woo/huma/internal/adapter"
	"github.com/park-jun-woo/huma/internal/analyzer"
	"github.com/park-jun-woo/huma/internal/config"
	"github.com/park-jun-woo/huma/internal/scanner"
	"github.com/park-jun-woo/huma/internal/session"
)

func TestBranchProvenanceString(t *testing.T) {
	tests := []struct {
		name string
		p    BranchProvenance
		want string
	}{
		{"both", BranchProvenance{HasSource: true, HasDeclared: true}, "both"},
		{"source", BranchProvenance{HasSource: true}, "source"},
		{"declared", BranchProvenance{HasDeclared: true}, "declared"},
		{"none", BranchProvenance{}, "none"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.String(); got != tt.want {
				t.Errorf("String() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestBranchStatuses(t *testing.T) {
	got := branchStatuses([]analyzer.ResponseBranch{{Status: 200}, {Status: 404}})
	if !reflect.DeepEqual(got, []int{200, 404}) {
		t.Errorf("branchStatuses = %v, want [200 404]", got)
	}
	if got := branchStatuses(nil); len(got) != 0 {
		t.Errorf("branchStatuses(nil) = %v, want empty", got)
	}
}

func TestConcatBranches(t *testing.T) {
	src := []analyzer.ResponseBranch{{Status: 200, Line: 5}}
	decl := []analyzer.ResponseBranch{{Status: 200, Line: 0}, {Status: 400}}
	got := concatBranches(src, decl)
	if len(got) != 3 {
		t.Fatalf("len = %d, want 3", len(got))
	}
	// source first → its Line wins on later dedup
	if got[0].Line != 5 {
		t.Errorf("first branch Line = %d, want 5 (source first)", got[0].Line)
	}
}

func TestCountKeywordOverlap(t *testing.T) {
	tests := []struct {
		name string
		a, b []string
		want int
	}{
		{"two overlap", []string{"users", "profiles", "x"}, []string{"users", "profiles"}, 2},
		{"none", []string{"a"}, []string{"b"}, 0},
		{"empty", nil, nil, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := countKeywordOverlap(tt.a, tt.b); got != tt.want {
				t.Errorf("countKeywordOverlap = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestDedupByStatus(t *testing.T) {
	in := []analyzer.ResponseBranch{
		{Status: 200, Line: 5},
		{Status: 200, Line: 9},
		{Status: 400, Line: 7},
	}
	got := dedupByStatus(in)
	if len(got) != 2 {
		t.Fatalf("len = %d, want 2", len(got))
	}
	if got[0].Line != 5 {
		t.Errorf("kept Line = %d, want 5 (first occurrence)", got[0].Line)
	}
}

func TestDeclaredBranches(t *testing.T) {
	ep := &scanner.Endpoint{Source: "x.go"}
	if got := declaredBranches(ep); got != nil {
		t.Errorf("empty Responses should yield nil, got %v", got)
	}
	ep.Responses = []byte(`[{"status":200},{"status":404}]`)
	got := declaredBranches(ep)
	if len(got) != 2 {
		t.Fatalf("len = %d, want 2", len(got))
	}
}

func TestExtractFilenameKeywords(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want []string
	}{
		{"strips method and id", "get_users_id.hurl", []string{"users"}},
		{"multiple", "post_orgs_teams.hurl", []string{"orgs", "teams"}},
		{"too few parts", "get.hurl", nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractFilenameKeywords(tt.in)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExtractPathKeywords(t *testing.T) {
	got := extractPathKeywords("/api/v1/users/:user_id/profiles")
	want := []string{"api", "v1", "users", "profiles"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
	if got := extractPathKeywords("/"); got != nil {
		t.Errorf("root path should yield nil, got %v", got)
	}
}

func TestHasOracle(t *testing.T) {
	if !hasOracle(&scanner.Endpoint{Source: "x.go"}, nil) {
		t.Error("source-linked endpoint should have oracle")
	}
	if hasOracle(&scanner.Endpoint{}, nil) {
		t.Error("no source, no cov should have no oracle")
	}
	if hasOracle(&scanner.Endpoint{}, &adapter.CoverageResult{Total: 0}) {
		t.Error("Total=0 should have no oracle")
	}
	if !hasOracle(&scanner.Endpoint{}, &adapter.CoverageResult{Total: 5}) {
		t.Error("instrumented (Total>0) should have oracle")
	}
}

func TestLineCohesion(t *testing.T) {
	got := lineCohesion([]analyzer.ResponseBranch{{Line: 5}, {Line: 5}, {Line: 9}})
	if got[5] != 2 || got[9] != 1 {
		t.Errorf("cohesion = %v, want {5:2, 9:1}", got)
	}
}

func TestMatchLineMethod(t *testing.T) {
	m, u, ok := matchLineMethod("GET {{host}}/x")
	if !ok || m != "GET" || u != "{{host}}/x" {
		t.Errorf("got (%q,%q,%v)", m, u, ok)
	}
	if _, _, ok := matchLineMethod("HTTP 200"); ok {
		t.Error("non-method line should not match")
	}
}

func TestProvenanceLabel(t *testing.T) {
	if provenanceLabel("") != "n/a" {
		t.Error("empty → n/a")
	}
	if provenanceLabel("both") != "both" {
		t.Error("passthrough")
	}
}

func TestProvenanceOf(t *testing.T) {
	src := []analyzer.ResponseBranch{{Status: 200}}
	p := provenanceOf(src, nil)
	if !p.HasSource || p.HasDeclared {
		t.Errorf("got %+v, want source only", p)
	}
	p = provenanceOf(nil, src)
	if p.HasSource || !p.HasDeclared {
		t.Errorf("got %+v, want declared only", p)
	}
}

func TestResolveGate(t *testing.T) {
	if resolveGate(&config.Config{}, 3) != 3 {
		t.Error("no RequireCRI → reachableMax")
	}
	cri := 1
	if resolveGate(&config.Config{RequireCRI: &cri}, 3) != 1 {
		t.Error("explicit RequireCRI wins")
	}
}

func TestSourceBranches(t *testing.T) {
	if got := sourceBranches(&scanner.Endpoint{}, "go"); got != nil {
		t.Errorf("no source → nil, got %v", got)
	}
}

func TestSplitClientBranches(t *testing.T) {
	in := []analyzer.ResponseBranch{{Status: 200}, {Status: 404}, {Status: 500}, {Status: 503}}
	client, advisory := splitClientBranches(in)
	if len(client) != 2 {
		t.Errorf("client = %v, want 2 (200,404)", client)
	}
	if len(advisory) != 2 {
		t.Errorf("advisory = %v, want 2 (500,503)", advisory)
	}
}

func TestStripHurlURL(t *testing.T) {
	tests := []struct {
		in, want string
	}{
		{"{{host}}/api/users/{{user_id}}", "/api/users/:user_id"},
		{"/api/health", "/api/health"},
	}
	for _, tt := range tests {
		if got := stripHurlURL(tt.in); got != tt.want {
			t.Errorf("stripHurlURL(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

func TestUncoveredBranches(t *testing.T) {
	branches := []analyzer.ResponseBranch{{Status: 200, Line: 4}, {Status: 400, Line: 5}, {Status: 404, Line: 0}}
	cov := &adapter.CoverageResult{CoveredLines: map[int]bool{4: true}}
	got := uncoveredBranches(branches, cov)
	// 200 (line 4) covered; 400 (line 5) uncovered; 404 (line 0) always uncovered
	if len(got) != 2 {
		t.Fatalf("got %d uncovered, want 2", len(got))
	}
	// nil cov → all uncovered
	if got := uncoveredBranches(branches, nil); len(got) != 3 {
		t.Errorf("nil cov → all uncovered, got %d", len(got))
	}
}

func TestBranchBound(t *testing.T) {
	cov := &adapter.CoverageResult{CoveredLines: map[int]bool{4: true}}
	cohesion := map[int]int{4: 1, 5: 2}
	// covered, unique line → bound
	if !branchBound(analyzer.ResponseBranch{Status: 200, Line: 4}, cov, nil, "GET /x", cohesion) {
		t.Error("covered unique line should be bound")
	}
	// line 0 → not bound
	if branchBound(analyzer.ResponseBranch{Status: 200, Line: 0}, cov, nil, "GET /x", cohesion) {
		t.Error("line 0 should not be bound")
	}
	// shared line (cohesion>1) → not bound
	if branchBound(analyzer.ResponseBranch{Status: 200, Line: 5}, cov, nil, "GET /x", cohesion) {
		t.Error("cohesive line should not be bound")
	}
	// exempt → bound regardless
	ex := []config.UnreachableEntry{{Endpoint: "GET /x", Status: 404, Reason: "r", Evidence: "e"}}
	if !branchBound(analyzer.ResponseBranch{Status: 404, Line: 0}, cov, ex, "GET /x", cohesion) {
		t.Error("exempt branch should be bound")
	}
}

func TestBranchBindingOK(t *testing.T) {
	ep := &scanner.Endpoint{Method: "GET", Path: "/x"}
	// no branches → trivially OK
	if !branchBindingOK(ep, nil, &adapter.CoverageResult{}) {
		t.Error("no branches → OK")
	}
	cov := &adapter.CoverageResult{CoveredLines: map[int]bool{4: true}}
	// uncovered branch → not OK
	if branchBindingOK(ep, []analyzer.ResponseBranch{{Status: 200, Line: 9}}, cov) {
		t.Error("uncovered branch → not bound")
	}
	// all covered unique → OK
	if !branchBindingOK(ep, []analyzer.ResponseBranch{{Status: 200, Line: 4}}, cov) {
		t.Error("covered branch → bound")
	}
}

func TestAllExempt(t *testing.T) {
	ex := []config.UnreachableEntry{{Endpoint: "GET /x", Status: 404, Reason: "r", Evidence: "e"}}
	branches := []analyzer.ResponseBranch{{Status: 404}}
	if !allExempt(branches, ex, "GET /x") {
		t.Error("404 exempt → all exempt")
	}
	branches = append(branches, analyzer.ResponseBranch{Status: 400})
	if allExempt(branches, ex, "GET /x") {
		t.Error("400 not exempt → not all exempt")
	}
}

func TestStalled(t *testing.T) {
	if stalled(nil, &adapter.CoverageResult{Percent: 50}) {
		t.Error("nil entry → not stalled")
	}
	e := &session.Entry{ImproveCount: 0, PrevCoverage: 50}
	if stalled(e, &adapter.CoverageResult{Percent: 50}) {
		t.Error("ImproveCount 0 → not stalled")
	}
	e = &session.Entry{ImproveCount: 1, PrevCoverage: 50}
	if !stalled(e, &adapter.CoverageResult{Percent: 50}) {
		t.Error("retried, no improvement → stalled")
	}
	if stalled(e, &adapter.CoverageResult{Percent: 60}) {
		t.Error("improved → not stalled")
	}
}

func TestStatusLine(t *testing.T) {
	tests := []struct {
		name     string
		entry    session.Entry
		contains string
	}{
		{"unverified", session.Entry{Status: session.StatusUnverified, Endpoint: scanner.Endpoint{Method: "GET", Path: "/x"}}, "UNVERIFIED"},
		{"pass", session.Entry{Status: session.StatusPass, CRI: 3, AGrade: 2, Endpoint: scanner.Endpoint{Method: "GET", Path: "/x"}}, "CRI 3"},
		{"improve", session.Entry{Status: session.StatusImprove, Coverage: 50, Endpoint: scanner.Endpoint{Method: "GET", Path: "/x"}}, "IMPROVE"},
		{"todo", session.Entry{Status: session.StatusTodo, Endpoint: scanner.Endpoint{Method: "GET", Path: "/x"}}, "TODO"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := statusLine(tt.entry); !strings.Contains(got, tt.contains) {
				t.Errorf("statusLine = %q, want contains %q", got, tt.contains)
			}
		})
	}
}

func TestDistribution(t *testing.T) {
	tests := []struct {
		name string
		r    scanner.LinkResult
		want string
	}{
		{"empty lang known", scanner.LinkResult{LangKnown: true}, "none"},
		{"empty lang unknown", scanner.LinkResult{LangKnown: false}, "lang=unknown, low-confidence"},
		{"with exts", scanner.LinkResult{LangKnown: true, ByExt: map[string]int{".go": 142, ".py": 3}}, "go: 142, py: 3"},
		{"unknown with exts", scanner.LinkResult{LangKnown: false, ByExt: map[string]int{".go": 1}}, "go: 1 (lang=unknown, low-confidence)"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := distribution(tt.r); got != tt.want {
				t.Errorf("distribution = %q, want %q", got, tt.want)
			}
		})
	}
}
