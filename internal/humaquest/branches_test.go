package humaquest

import (
	"encoding/json"
	"testing"

	"github.com/park-jun-woo/huma/internal/analyzer"
	"github.com/park-jun-woo/huma/internal/hurlcheck"
	"github.com/park-jun-woo/huma/internal/scanner"
)

// goHandlerSrc is a real Go handler with three distinct response branches on
// separate lines, so the analyzer yields unambiguous (Status, Line) pairs.
const goHandlerSrc = `package main

import "net/http"

func CreateUser(c interface{}) {
	c.JSON(http.StatusOK, nil)
	c.JSON(404, nil)
	c.AbortWithStatusJSON(http.StatusInternalServerError, nil)
}
`

// writeGoHandler writes goHandlerSrc to a temp file and returns its path.
func writeGoHandler(t *testing.T) string {
	t.Helper()
	return writeTempFile(t, "handler.go", goHandlerSrc)
}

// ---------------------------------------------------------------------------
// analyzeBranches
// ---------------------------------------------------------------------------

func TestAnalyzeBranches(t *testing.T) {
	src := writeGoHandler(t)

	t.Run("no source → nil", func(t *testing.T) {
		if got := analyzeBranches(&scanner.Endpoint{}, "go"); got != nil {
			t.Errorf("no source want nil, got %+v", got)
		}
	})

	t.Run("unsupported lang → nil", func(t *testing.T) {
		ep := &scanner.Endpoint{Source: src, Handler: "CreateUser"}
		if got := analyzeBranches(ep, "cobol"); got != nil {
			t.Errorf("nil analyzer want nil, got %+v", got)
		}
	})

	t.Run("analysis failure → nil", func(t *testing.T) {
		ep := &scanner.Endpoint{Source: "/no/such/file.go", Handler: "CreateUser"}
		if got := analyzeBranches(ep, "go"); got != nil {
			t.Errorf("analyze error want nil, got %+v", got)
		}
	})

	t.Run("success extracts branches", func(t *testing.T) {
		ep := &scanner.Endpoint{Source: src, Handler: "CreateUser"}
		got := analyzeBranches(ep, "go")
		want := []int{200, 404, 500}
		if !equalIntSet(sortedStatuses(got), want) {
			t.Errorf("statuses = %v, want %v", sortedStatuses(got), want)
		}
		for _, b := range got {
			if b.Line == 0 {
				t.Errorf("branch %d missing source line", b.Status)
			}
		}
	})
}

// ---------------------------------------------------------------------------
// sourceBranches
// ---------------------------------------------------------------------------

func TestSourceBranches(t *testing.T) {
	if got := sourceBranches(&scanner.Endpoint{}, "go"); got != nil {
		t.Errorf("no source link want nil, got %+v", got)
	}
	ep := &scanner.Endpoint{Source: writeGoHandler(t), Handler: "CreateUser"}
	if got := sourceBranches(ep, "go"); len(got) != 3 {
		t.Errorf("source-linked want 3 branches, got %d", len(got))
	}
}

// ---------------------------------------------------------------------------
// declaredBranches
// ---------------------------------------------------------------------------

func TestDeclaredBranches(t *testing.T) {
	if got := declaredBranches(&scanner.Endpoint{}); got != nil {
		t.Errorf("no responses want nil, got %+v", got)
	}
	ep := &scanner.Endpoint{
		Source:    "h.go",
		Responses: json.RawMessage(`[{"status":200,"line":1},{"status":404,"line":2}]`),
	}
	got := declaredBranches(ep)
	if !equalIntSet(sortedStatuses(got), []int{200, 404}) {
		t.Errorf("declared statuses = %v", sortedStatuses(got))
	}
}

// ---------------------------------------------------------------------------
// responseBranches — the monotonic union
// ---------------------------------------------------------------------------

func TestResponseBranches_Union(t *testing.T) {
	src := writeGoHandler(t)
	// Source yields 200/404/500; OpenAPI additionally declares 409.
	ep := &scanner.Endpoint{
		Source:    src,
		Handler:   "CreateUser",
		Responses: json.RawMessage(`[{"status":200,"line":1},{"status":409,"line":2}]`),
	}
	client, prov := responseBranches(ep, "go")

	// 500 is advisory (server) → excluded from client set; union of client
	// branches is {200, 404, 409}.
	if !equalIntSet(sortedStatuses(client), []int{200, 404, 409}) {
		t.Errorf("client union = %v, want [200 404 409]", sortedStatuses(client))
	}
	// The source-derived 200 Line (non-zero) must win the dedup over the
	// declared 200 (line 1 from OpenAPI is overridden by source order).
	for _, b := range client {
		if b.Status == 200 && b.Line == 1 {
			t.Error("declared 200 line should not override source 200")
		}
	}
	if prov.String() != "both" {
		t.Errorf("provenance = %q, want both", prov.String())
	}
}

func TestResponseBranches_SourceOnly(t *testing.T) {
	ep := &scanner.Endpoint{Source: writeGoHandler(t), Handler: "CreateUser"}
	_, prov := responseBranches(ep, "go")
	if prov.String() != "source" {
		t.Errorf("provenance = %q, want source", prov.String())
	}
}

func TestResponseBranches_DeclaredOnly(t *testing.T) {
	ep := &scanner.Endpoint{Responses: json.RawMessage(`[{"status":200,"line":1}]`)}
	client, prov := responseBranches(ep, "go")
	if prov.String() != "declared" {
		t.Errorf("provenance = %q, want declared", prov.String())
	}
	if !equalIntSet(sortedStatuses(client), []int{200}) {
		t.Errorf("client = %v", sortedStatuses(client))
	}
}

func TestResponseBranches_None(t *testing.T) {
	client, prov := responseBranches(&scanner.Endpoint{}, "go")
	if len(client) != 0 || prov.String() != "none" {
		t.Errorf("empty endpoint → no branches/none provenance, got %d / %q", len(client), prov.String())
	}
}

// ---------------------------------------------------------------------------
// staticAGrade
// ---------------------------------------------------------------------------

func TestStaticAGrade(t *testing.T) {
	branches := []analyzer.ResponseBranch{br(200, 1), br(404, 2)}

	t.Run("min across statuses", func(t *testing.T) {
		entries := []hurlcheck.HurlEntry{
			{Status: 200, Grade: 3},
			{Status: 404, Grade: 1},
		}
		if got := staticAGrade(entries, branches); got != 1 {
			t.Errorf("A-grade = %d, want 1 (min)", got)
		}
	})

	t.Run("uncovered status contributes 0", func(t *testing.T) {
		entries := []hurlcheck.HurlEntry{{Status: 200, Grade: 3}}
		if got := staticAGrade(entries, branches); got != 0 {
			t.Errorf("missing 404 assert → 0, got %d", got)
		}
	})

	t.Run("no branches → 0", func(t *testing.T) {
		entries := []hurlcheck.HurlEntry{{Status: 200, Grade: 3}}
		if got := staticAGrade(entries, nil); got != 0 {
			t.Errorf("empty statuses → 0, got %d", got)
		}
	})
}

// equalIntSet compares two int slices for equality (order-sensitive after sort).
func equalIntSet(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
