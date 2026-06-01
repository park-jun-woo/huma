package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/park-jun-woo/huma/internal/adapter"
	"github.com/park-jun-woo/huma/internal/config"
	"github.com/park-jun-woo/huma/internal/hurlcheck"
	"github.com/park-jun-woo/huma/internal/runner"
	"github.com/park-jun-woo/huma/internal/scanner"
	"github.com/park-jun-woo/huma/internal/session"
)

func TestPrintLinkResult(t *testing.T) {
	out := captureStdout(t, func() {
		printLinkResult(scanner.LinkResult{
			Linked: 2, Skipped: 1, ExtMismatch: 1, LangKnown: true,
			ByExt:        map[string]int{".go": 2},
			SkipMessages: []string{"ep3 skipped: ext mismatch"},
		}, 3, "/root")
	})
	for _, want := range []string{"Linked 2/3", "go: 2", "Skipped 1", "ext mismatch"} {
		if !contains(out, want) {
			t.Errorf("output missing %q: %q", want, out)
		}
	}
}

func TestPrintLiveNonPass(t *testing.T) {
	ep := &scanner.Endpoint{Method: "GET", Path: "/x"}
	cov := &adapter.CoverageResult{Total: 1, Covered: 0, Percent: 0}
	cfg := &config.Config{}

	if handled := captureBool(t, func() bool { return printLiveNonPass(outcomeImprove, ep, "t.hurl", cov, cfg) }); !handled {
		t.Error("improve should be handled")
	}
	if handled := captureBool(t, func() bool { return printLiveNonPass(outcomeUnverified, ep, "t.hurl", cov, cfg) }); !handled {
		t.Error("unverified should be handled")
	}
	if handled := captureBool(t, func() bool { return printLiveNonPass(outcomePass, ep, "t.hurl", cov, cfg) }); handled {
		t.Error("pass should not be handled")
	}
}

func TestPrintStaticNonPass(t *testing.T) {
	ep := &scanner.Endpoint{Method: "GET", Path: "/x"}
	res := &hurlcheck.ResponseCoverageResult{Total: 1, Percent: 0}
	cfg := &config.Config{}

	if handled := captureBool(t, func() bool { return printStaticNonPass(outcomeImprove, ep, "t.hurl", res, cfg) }); !handled {
		t.Error("improve should be handled")
	}
	if handled := captureBool(t, func() bool { return printStaticNonPass(outcomeUnverified, ep, "t.hurl", res, cfg) }); !handled {
		t.Error("unverified should be handled")
	}
	if handled := captureBool(t, func() bool { return printStaticNonPass(outcomeDone, ep, "t.hurl", res, cfg) }); handled {
		t.Error("done should not be handled")
	}
}

func TestPrecheckEndpoints(t *testing.T) {
	// nil cfg is a no-op
	sess := session.New()
	precheckEndpoints(sess, nil)

	dir := withSessionDir(t)
	src := writeGoHandler(t, dir, "StatusOK")
	hurlDir := filepath.Join(dir, "hurl")
	os.MkdirAll(hurlDir, 0o755)

	ep := scanner.Endpoint{ID: "ep1", Method: "GET", Path: "/x", Handler: "H", Source: src, Line: 3}
	// place a hurl file with the expected name so FindHurlFile picks it up
	name := filepath.Base(runner.HurlFileName(&ep, hurlDir))
	os.WriteFile(filepath.Join(hurlDir, name), []byte("GET {{host}}/x\nHTTP 200\n[Asserts]\njsonpath \"$.id\" exists\n"), 0o644)

	sess = session.New()
	sess.Merge([]scanner.Endpoint{ep})
	cfg := &config.Config{HurlDir: hurlDir, Scan: config.ScanConfig{Lang: "go"}}
	precheckEndpoints(sess, cfg)

	if sess.Entries[0].Status != session.StatusPass {
		t.Fatalf("expected SCAFFOLDED PASS, got %s", sess.Entries[0].Status)
	}
	if sess.Entries[0].CRI != 1 {
		t.Fatalf("expected CRI 1, got %d", sess.Entries[0].CRI)
	}
}

func TestPrecheckEntry_NoHurl(t *testing.T) {
	withSessionDir(t)
	ep := scanner.Endpoint{ID: "ep1", Method: "GET", Path: "/x"}
	sess := session.New()
	sess.Merge([]scanner.Endpoint{ep})
	// no hurl file → returns early, stays TODO
	cfg := &config.Config{HurlDir: "/nonexistent", Scan: config.ScanConfig{Lang: "go"}}
	precheckEntry(sess, &sess.Entries[0], cfg)
	if sess.Entries[0].Status != session.StatusTodo {
		t.Fatalf("expected TODO (no hurl), got %s", sess.Entries[0].Status)
	}
}

func TestPrecheckEntry_UnverifiedNoOracle(t *testing.T) {
	dir := withSessionDir(t)
	hurlDir := filepath.Join(dir, "hurl")
	os.MkdirAll(hurlDir, 0o755)
	// unlinked endpoint, no responses → no oracle / no branches → UNVERIFIED
	ep := scanner.Endpoint{ID: "ep1", Method: "GET", Path: "/x"}
	name := filepath.Base(runner.HurlFileName(&ep, hurlDir))
	os.WriteFile(filepath.Join(hurlDir, name), []byte("GET {{host}}/x\nHTTP 200\n"), 0o644)
	sess := session.New()
	sess.Merge([]scanner.Endpoint{ep})
	cfg := &config.Config{HurlDir: hurlDir, Scan: config.ScanConfig{Lang: "go"}}
	precheckEntry(sess, &sess.Entries[0], cfg)
	if sess.Entries[0].Status != session.StatusUnverified {
		t.Fatalf("expected UNVERIFIED, got %s", sess.Entries[0].Status)
	}
}

func TestFindContentMatch(t *testing.T) {
	dir := t.TempDir()
	eps := []scanner.Endpoint{
		{Method: "GET", Path: "/api/users/:id"},
	}
	// a hurl declaring GET /api/users/{{id}} but named wrong → finds expected
	got := findContentMatch("GET", "/api/users/:id", "wrong_name.hurl", eps, dir)
	if got == "" {
		t.Error("expected a match for matching method+path")
	}
	// no method match
	if got := findContentMatch("POST", "/api/users/:id", "x.hurl", eps, dir); got != "" {
		t.Errorf("expected no match for POST, got %q", got)
	}
}

func TestFindKeywordMatch(t *testing.T) {
	dir := t.TempDir()
	eps := []scanner.Endpoint{
		{Method: "GET", Path: "/api/v1/users/profiles"},
		{Method: "POST", Path: "/api/v1/orders"},
	}
	got := findKeywordMatch("get_users_profiles.hurl", eps, dir)
	if got == "" {
		t.Error("expected keyword match for users/profiles")
	}
	// no keywords → empty
	if got := findKeywordMatch("get.hurl", eps, dir); got != "" {
		t.Errorf("no keywords → empty, got %q", got)
	}
}

func TestFindLikelyMatch(t *testing.T) {
	dir := t.TempDir()
	eps := []scanner.Endpoint{{Method: "GET", Path: "/api/users/:id"}}
	// content parse path: write a hurl with parseable request
	hurlPath := filepath.Join(dir, "misnamed.hurl")
	os.WriteFile(hurlPath, []byte("GET {{host}}/api/users/{{id}}\nHTTP 200\n"), 0o644)
	got := findLikelyMatch("misnamed.hurl", eps, dir)
	if got == "" {
		t.Error("expected likely match via content")
	}
}

func TestWarnMismatchedHurlFiles(t *testing.T) {
	dir := t.TempDir()
	eps := []scanner.Endpoint{{Method: "GET", Path: "/api/users/:id"}}
	// expected name for the endpoint
	expected := filepath.Base(runner.HurlFileName(&eps[0], dir))
	// write a correctly named file (no warning) and a mismatched one
	os.WriteFile(filepath.Join(dir, expected), []byte("GET {{host}}/api/users/{{id}}\nHTTP 200\n"), 0o644)
	os.WriteFile(filepath.Join(dir, "stray_file.hurl"), []byte("GET {{host}}/other\nHTTP 200\n"), 0o644)

	// a misnamed file whose content DOES match an endpoint → expected != "" branch
	os.WriteFile(filepath.Join(dir, "misnamed_users.hurl"), []byte("GET {{host}}/api/users/{{id}}\nHTTP 200\n"), 0o644)
	// a non-.hurl file and a subdirectory → skipped (line 35 continue branch)
	os.WriteFile(filepath.Join(dir, "notes.txt"), []byte("ignore"), 0o644)
	os.MkdirAll(filepath.Join(dir, "subdir.hurl"), 0o755)

	out := captureStdout(t, func() { warnMismatchedHurlFiles(dir, eps) })
	if !contains(out, "WARNING") || !contains(out, "stray_file.hurl") {
		t.Errorf("expected mismatch warning, got %q", out)
	}
	if !contains(out, "expected:") {
		t.Errorf("expected an 'expected:' suggestion line, got %q", out)
	}
	if !contains(out, "no matching endpoint") {
		t.Errorf("expected a 'no matching endpoint' line, got %q", out)
	}

	// no mismatches → no output
	dir2 := t.TempDir()
	os.WriteFile(filepath.Join(dir2, filepath.Base(runner.HurlFileName(&eps[0], dir2))), []byte("x"), 0o644)
	out = captureStdout(t, func() { warnMismatchedHurlFiles(dir2, eps) })
	if out != "" {
		t.Errorf("expected no warning, got %q", out)
	}

	// unreadable dir → silent return
	captureStdout(t, func() { warnMismatchedHurlFiles("/nonexistent-dir-xyz", eps) })
}

// --- test helpers ---

func captureBool(t *testing.T, fn func() bool) bool {
	t.Helper()
	var result bool
	captureStdout(t, func() { result = fn() })
	return result
}

func contains(s, sub string) bool {
	return strings.Contains(s, sub)
}
