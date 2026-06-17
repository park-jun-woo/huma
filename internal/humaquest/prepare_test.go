package humaquest

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/park-jun-woo/huma/internal/hurlcheck"
	"github.com/park-jun-woo/huma/internal/runner"
	"github.com/park-jun-woo/huma/internal/scanner"
	"github.com/park-jun-woo/reins/pkg/quest"
)

// validHurl is a conforming .hurl body: a single GET entry referencing the
// {{base_url}} host template var with two assertions (A-grade 3).
const validHurl = `GET {{base_url}}/api/v1/users
HTTP 200
jsonpath "$.id" exists
jsonpath "$.name" exists
`

// prepEndpoint is the endpoint under test in Prepare/locateHurl scenarios. Its
// conventional .hurl path is hurl/get_api_v1_users.hurl.
func prepEndpoint() scanner.Endpoint {
	return scanner.Endpoint{
		ID:     "GET_/api/v1/users",
		Method: "GET",
		Path:   "/api/v1/users",
	}
}

// prepItem wraps an endpoint payload in an Item ready for Prepare.
func prepItem(t *testing.T, ep scanner.Endpoint) *quest.Item {
	t.Helper()
	it := &quest.Item{Key: ep.ID, State: quest.TODO}
	if err := it.SetPayload(&ep); err != nil {
		t.Fatalf("SetPayload: %v", err)
	}
	return it
}

// ---------------------------------------------------------------------------
// Prepare
// ---------------------------------------------------------------------------

func TestPrepare_HappyPath(t *testing.T) {
	dir := t.TempDir()
	chdir(t, dir)

	// Conventional .hurl at hurl/get_api_v1_users.hurl.
	hurlPath := filepath.Join(dir, "hurl", "get_api_v1_users.hurl")
	if err := os.MkdirAll(filepath.Dir(hurlPath), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(hurlPath, []byte(validHurl), 0o644); err != nil {
		t.Fatal(err)
	}

	// A real handler source so gate.Context.Source is populated.
	if err := os.WriteFile(filepath.Join(dir, "handlers.go"),
		[]byte("package main\n\nfunc ListUsers(c interface{}) {\n\treturn\n}\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	ep := prepEndpoint()
	ep.Source = "handlers.go"
	ep.Handler = "ListUsers"
	it := prepItem(t, ep)

	ctx, verdict, err := (humaDef{}).Prepare(nil, it, nil)
	if err != nil {
		t.Fatalf("Prepare: %v", err)
	}
	if verdict != nil {
		t.Fatalf("happy path should not short-circuit, got verdict %+v", verdict)
	}

	hi, ok := ctx.Submission.(*hurlInfo)
	if !ok {
		t.Fatalf("Submission is %T, want *hurlInfo", ctx.Submission)
	}
	if len(hi.Entries) != 1 {
		t.Fatalf("want 1 parsed entry, got %d", len(hi.Entries))
	}
	if hi.AGrade != 3 {
		t.Errorf("AGrade = %d, want 3 (status + 2 asserts)", hi.AGrade)
	}
	if !hi.NamingOK {
		t.Error("NamingOK = false for conventional path, want true")
	}
	if !hi.HostVarOK {
		t.Error("HostVarOK = false despite {{base_url}} URLs, want true")
	}
	if filepath.Base(hi.HurlPath) != "get_api_v1_users.hurl" {
		t.Errorf("HurlPath = %q, want conventional name", hi.HurlPath)
	}
	if ctx.Item != it {
		t.Error("ctx.Item not wired to the input item")
	}
	if !strings.Contains(ctx.Source, "ListUsers") {
		t.Errorf("ctx.Source missing handler body, got %q", ctx.Source)
	}
}

func TestPrepare_RawPathVerbatim(t *testing.T) {
	// raw set → used verbatim (no manifest, no conventional layout needed).
	chdir(t, t.TempDir())
	hurlPath := writeTempFile(t, "custom.hurl", validHurl)

	it := prepItem(t, prepEndpoint())
	ctx, verdict, err := (humaDef{}).Prepare(nil, it, []byte("  "+hurlPath+"  "))
	if err != nil {
		t.Fatalf("Prepare: %v", err)
	}
	if verdict != nil {
		t.Fatalf("unexpected short-circuit verdict %+v", verdict)
	}
	hi := ctx.Submission.(*hurlInfo)
	if hi.HurlPath != hurlPath {
		t.Errorf("HurlPath = %q, want trimmed %q", hi.HurlPath, hurlPath)
	}
	// A hand-named path that does not match the convention → NamingOK false.
	if hi.NamingOK {
		t.Error("NamingOK = true for custom.hurl, want false")
	}
}

func TestPrepare_Exempt_ShortCircuitsSkip(t *testing.T) {
	dir := t.TempDir()
	chdir(t, dir)

	// unreachable.yaml exempting the single declared 200 branch.
	if err := os.MkdirAll(filepath.Join(dir, ".huma"), 0o755); err != nil {
		t.Fatal(err)
	}
	unreach := "- endpoint: GET /api/v1/users\n  status: 200\n  reason: dead code\n  evidence: handlers.go:1\n"
	if err := os.WriteFile(filepath.Join(dir, ".huma", "unreachable.yaml"), []byte(unreach), 0o644); err != nil {
		t.Fatal(err)
	}

	ep := prepEndpoint()
	ep.Responses = json.RawMessage(`[{"status":200,"line":1}]`)
	it := prepItem(t, ep)

	ctx, verdict, err := (humaDef{}).Prepare(nil, it, nil)
	if err != nil {
		t.Fatalf("Prepare: %v", err)
	}
	if verdict == nil || verdict.Outcome != quest.OutSkip {
		t.Fatalf("want OutSkip verdict, got verdict=%+v", verdict)
	}
	if ctx.Submission != nil {
		t.Errorf("exempt path should return zero Context, got Submission %v", ctx.Submission)
	}
	if !strings.Contains(verdict.Feedback, "exempt") {
		t.Errorf("skip feedback missing reason: %q", verdict.Feedback)
	}
}

func TestPrepare_MissingHurl_FailsH01(t *testing.T) {
	chdir(t, t.TempDir())

	it := prepItem(t, prepEndpoint())
	ctx, verdict, err := (humaDef{}).Prepare(nil, it, nil)
	if err != nil {
		t.Fatalf("Prepare: %v", err)
	}
	if verdict == nil || verdict.Outcome != quest.OutFail {
		t.Fatalf("want OutFail verdict, got %+v", verdict)
	}
	if verdict.RootCause != "H-01" {
		t.Errorf("RootCause = %q, want H-01", verdict.RootCause)
	}
	if !strings.Contains(verdict.Feedback, "get_api_v1_users.hurl") {
		t.Errorf("feedback should name the expected path, got %q", verdict.Feedback)
	}
	if ctx.Submission != nil {
		t.Error("missing-hurl path should return zero Context")
	}
}

func TestPrepare_DecodePayloadError(t *testing.T) {
	chdir(t, t.TempDir())
	it := &quest.Item{Key: "bad", State: quest.TODO, Payload: json.RawMessage(`{"id": 12345`)}
	if _, _, err := (humaDef{}).Prepare(nil, it, nil); err == nil {
		t.Fatal("expected DecodePayload error, got nil")
	}
}

// TestPrepare_ParseErrorPropagates drives the locateHurl err branch through
// Prepare: a present-but-unparseable .hurl (an over-long line) surfaces as a
// non-nil error rather than a verdict.
func TestPrepare_ParseErrorPropagates(t *testing.T) {
	chdir(t, t.TempDir())
	hurlPath := writeTempFile(t, "huge.hurl", strings.Repeat("A", 2*1024*1024))

	it := prepItem(t, prepEndpoint())
	_, _, err := (humaDef{}).Prepare(nil, it, []byte(hurlPath))
	if err == nil {
		t.Fatal("expected parse error to propagate from Prepare")
	}
}

// ---------------------------------------------------------------------------
// locateHurl
// ---------------------------------------------------------------------------

func TestLocateHurl_DerivedPathMissing(t *testing.T) {
	ep := prepEndpoint()
	path, entries, found, err := locateHurl(nil, ep, "hurl")
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if found {
		t.Error("found = true for nonexistent derived path")
	}
	if entries != nil {
		t.Error("entries should be nil when not found")
	}
	want := runner.HurlFileName(&ep, "hurl")
	if path != want {
		t.Errorf("derived path = %q, want %q", path, want)
	}
}

func TestLocateHurl_VerbatimPathFound(t *testing.T) {
	hurlPath := writeTempFile(t, "custom.hurl", validHurl)
	path, entries, found, err := locateHurl([]byte("  "+hurlPath+"\n"), prepEndpoint(), "hurl")
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if !found {
		t.Fatal("found = false for existing file")
	}
	if path != hurlPath {
		t.Errorf("path = %q, want trimmed verbatim %q", path, hurlPath)
	}
	if len(entries) != 1 {
		t.Errorf("want 1 entry, got %d", len(entries))
	}
}

func TestLocateHurl_VerbatimMissing(t *testing.T) {
	missing := filepath.Join(t.TempDir(), "nope.hurl")
	path, _, found, err := locateHurl([]byte(missing), prepEndpoint(), "hurl")
	if err != nil {
		t.Fatalf("missing file must not be an error, got %v", err)
	}
	if found {
		t.Error("found = true for missing verbatim path")
	}
	if path != missing {
		t.Errorf("path = %q, want %q", path, missing)
	}
}

func TestLocateHurl_ParseError(t *testing.T) {
	// A present file with a line longer than the scanner buffer → parse error.
	hurlPath := writeTempFile(t, "huge.hurl", strings.Repeat("A", 2*1024*1024))
	path, entries, found, err := locateHurl([]byte(hurlPath), prepEndpoint(), "hurl")
	if err == nil {
		t.Fatal("expected parse error on oversized line")
	}
	if !found {
		t.Error("found should be true (file exists) even on parse error")
	}
	if entries != nil {
		t.Error("entries should be nil on parse error")
	}
	if path != hurlPath {
		t.Errorf("path = %q, want %q", path, hurlPath)
	}
}

// ---------------------------------------------------------------------------
// isExempt
// ---------------------------------------------------------------------------

func TestIsExempt_NoResponses(t *testing.T) {
	chdir(t, t.TempDir())
	if exempt, _ := isExempt(prepEndpoint()); exempt {
		t.Error("endpoint with no declared responses must not be exempt")
	}
}

func TestIsExempt_NoUnreachableFile(t *testing.T) {
	chdir(t, t.TempDir())
	ep := prepEndpoint()
	ep.Responses = json.RawMessage(`[{"status":200,"line":1}]`)
	if exempt, _ := isExempt(ep); exempt {
		t.Error("no unreachable.yaml → must not be exempt")
	}
}

func TestIsExempt_AllBranchesExempt(t *testing.T) {
	dir := t.TempDir()
	chdir(t, dir)
	if err := os.MkdirAll(filepath.Join(dir, ".huma"), 0o755); err != nil {
		t.Fatal(err)
	}
	unreach := "" +
		"- endpoint: GET /api/v1/users\n  status: 200\n  reason: r\n  evidence: e\n" +
		"- endpoint: GET /api/v1/users\n  status: 404\n  reason: r\n  evidence: e\n"
	if err := os.WriteFile(filepath.Join(dir, ".huma", "unreachable.yaml"), []byte(unreach), 0o644); err != nil {
		t.Fatal(err)
	}
	ep := prepEndpoint()
	ep.Responses = json.RawMessage(`[{"status":200,"line":1},{"status":404,"line":2}]`)
	exempt, why := isExempt(ep)
	if !exempt {
		t.Fatal("all branches exempt → expected exempt=true")
	}
	if !strings.Contains(why, "GET /api/v1/users") {
		t.Errorf("reason should name the endpoint, got %q", why)
	}
}

func TestIsExempt_PartialNotExempt(t *testing.T) {
	dir := t.TempDir()
	chdir(t, dir)
	if err := os.MkdirAll(filepath.Join(dir, ".huma"), 0o755); err != nil {
		t.Fatal(err)
	}
	// Only the 200 branch is exempt; 404 remains → not all exempt.
	unreach := "- endpoint: GET /api/v1/users\n  status: 200\n  reason: r\n  evidence: e\n"
	if err := os.WriteFile(filepath.Join(dir, ".huma", "unreachable.yaml"), []byte(unreach), 0o644); err != nil {
		t.Fatal(err)
	}
	ep := prepEndpoint()
	ep.Responses = json.RawMessage(`[{"status":200,"line":1},{"status":404,"line":2}]`)
	if exempt, _ := isExempt(ep); exempt {
		t.Error("one unexempt branch → must not be exempt")
	}
}

// ---------------------------------------------------------------------------
// namingOK
// ---------------------------------------------------------------------------

func TestNamingOK(t *testing.T) {
	ep := prepEndpoint()
	conventional := runner.HurlFileName(&ep, "hurl")
	tests := []struct {
		name string
		path string
		want bool
	}{
		{"derived conventional path", conventional, true},
		{"conventional name under other dir", filepath.Join("anywhere", "get_api_v1_users.hurl"), true},
		{"hand-named file", filepath.Join("hurl", "users.hurl"), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := namingOK(tt.path, ep, "hurl"); got != tt.want {
				t.Errorf("namingOK(%q) = %v, want %v", tt.path, got, tt.want)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// hostVarOK
// ---------------------------------------------------------------------------

func TestHostVarOK(t *testing.T) {
	tests := []struct {
		name    string
		entries []hurlcheck.HurlEntry
		want    bool
	}{
		{"empty set", nil, false},
		{
			"all reference host var",
			[]hurlcheck.HurlEntry{
				{URL: "{{base_url}}/a"},
				{URL: "{{base_url}}/b"},
			},
			true,
		},
		{
			"one hardcoded URL",
			[]hurlcheck.HurlEntry{
				{URL: "{{base_url}}/a"},
				{URL: "http://localhost:8080/b"},
			},
			false,
		},
		{
			"skip and empty-URL entries ignored",
			[]hurlcheck.HurlEntry{
				{URL: "http://hardcoded/x", Skip: true}, // skipped → ignored
				{URL: ""},                               // no URL → ignored
				{URL: "{{base_url}}/c"},
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := hostVarOK(tt.entries, "base_url"); got != tt.want {
				t.Errorf("hostVarOK = %v, want %v", got, tt.want)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// readHandlerSource
// ---------------------------------------------------------------------------

func TestReadHandlerSource(t *testing.T) {
	dir := t.TempDir()
	srcPath := filepath.Join(dir, "handlers.go")
	if err := os.WriteFile(srcPath,
		[]byte("package main\n\nfunc ListUsers(c interface{}) {\n\treturn\n}\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	t.Run("real source and handler", func(t *testing.T) {
		ep := prepEndpoint()
		ep.Source = srcPath
		ep.Handler = "ListUsers"
		got := readHandlerSource(ep)
		if !strings.Contains(got, "ListUsers") {
			t.Errorf("expected handler body, got %q", got)
		}
	})

	t.Run("empty source field", func(t *testing.T) {
		ep := prepEndpoint()
		ep.Handler = "ListUsers"
		if got := readHandlerSource(ep); got != "" {
			t.Errorf("empty Source → want \"\", got %q", got)
		}
	})

	t.Run("empty handler field", func(t *testing.T) {
		ep := prepEndpoint()
		ep.Source = srcPath
		if got := readHandlerSource(ep); got != "" {
			t.Errorf("empty Handler → want \"\", got %q", got)
		}
	})

	t.Run("missing source file", func(t *testing.T) {
		ep := prepEndpoint()
		ep.Source = filepath.Join(dir, "nope.go")
		ep.Handler = "ListUsers"
		if got := readHandlerSource(ep); got != "" {
			t.Errorf("missing file → want \"\", got %q", got)
		}
	})

	t.Run("handler not in file", func(t *testing.T) {
		ep := prepEndpoint()
		ep.Source = srcPath
		ep.Handler = "DoesNotExist"
		if got := readHandlerSource(ep); got != "" {
			t.Errorf("unlocatable handler → want \"\", got %q", got)
		}
	})
}
