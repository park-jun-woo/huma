package scanner

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// writePolyglotRepo creates a polyglot repo mirroring BUG-002: a Go backend
// handler in PascalCase under cmd/ and a same-named camelCase JS client under
// frontend/. Returns the repo root.
func writePolyglotRepo(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	os.MkdirAll(filepath.Join(root, "cmd"), 0o755)
	os.MkdirAll(filepath.Join(root, "frontend", "src", "api"), 0o755)
	os.WriteFile(filepath.Join(root, "cmd", "subscribers.go"), []byte(`package main

func (a *App) CreateSubscriber(c echo.Context) error { return nil }
func (a *App) GetSettings(c echo.Context) error { return nil }
`), 0o644)
	os.WriteFile(filepath.Join(root, "frontend", "src", "api", "index.js"), []byte(`export default {
  createSubscriber: (data) => http.post('/api/subscribers', data),
  getSettings: () => http.get('/api/settings'),
}
`), 0o644)
	return root
}

func TestLinkSource_BUG002_LinksGoNotFrontendJS(t *testing.T) {
	root := writePolyglotRepo(t)
	eps := []Endpoint{
		{ID: "1", Method: "POST", Path: "/subscribers", Handler: "createSubscriber"},
		{ID: "2", Method: "GET", Path: "/settings", Handler: "getSettings"},
	}
	res := LinkSource(eps, root, "go")
	if res.Linked != 2 {
		t.Fatalf("expected 2 linked, got %d (%+v)", res.Linked, res)
	}
	for _, ep := range eps {
		if !strings.HasSuffix(ep.Source, ".go") {
			t.Fatalf("expected Go link, got %s for %s", ep.Source, ep.Handler)
		}
		if strings.Contains(ep.Source, "frontend") {
			t.Fatalf("frontend JS mislink for %s: %s", ep.Handler, ep.Source)
		}
	}
	if res.ByExt[".go"] != 2 {
		t.Fatalf("expected go distribution 2, got %+v", res.ByExt)
	}
}

func TestLinkSource_GoLang_ExcludesJSCandidates(t *testing.T) {
	root := t.TempDir()
	os.WriteFile(filepath.Join(root, "client.js"), []byte(`export default { createSubscriber: () => {} }`), 0o644)
	eps := []Endpoint{{ID: "1", Handler: "createSubscriber"}}
	res := LinkSource(eps, root, "go")
	// .js is not a candidate under lang=go, so nothing matches; honest UNVERIFIED.
	if res.Linked != 0 || res.Skipped != 0 {
		t.Fatalf("expected 0 linked/0 skipped (js excluded), got %+v", res)
	}
	if eps[0].Source != "" {
		t.Fatalf("expected unlinked, got %s", eps[0].Source)
	}
}

func TestLinkSource_ObjectLiteralKeyNotMatchedForGo(t *testing.T) {
	root := t.TempDir()
	// A .go file that happens to contain an object-literal-like `name:` must not
	// match a Go handler def (only `func ... Name(` matches).
	os.WriteFile(filepath.Join(root, "x.go"), []byte(`package main

var m = map[string]any{
	"createSubscriber": nil,
}
`), 0o644)
	eps := []Endpoint{{ID: "1", Handler: "createSubscriber"}}
	res := LinkSource(eps, root, "go")
	if res.Linked != 0 {
		t.Fatalf("expected object-literal key NOT matched, got %+v", res)
	}
}

func TestLinkSource_AmbiguousRejected(t *testing.T) {
	root := t.TempDir()
	os.WriteFile(filepath.Join(root, "a.go"), []byte("package main\nfunc CreateUser() {}\n"), 0o644)
	os.WriteFile(filepath.Join(root, "b.go"), []byte("package main\nfunc CreateUser() {}\n"), 0o644)
	eps := []Endpoint{{ID: "1", Handler: "createUser"}}
	res := LinkSource(eps, root, "go")
	if res.Linked != 0 || res.Ambiguous != 1 || res.Skipped != 1 {
		t.Fatalf("expected ambiguous rejection, got %+v", res)
	}
	if eps[0].Source != "" {
		t.Fatalf("ambiguous must stay UNVERIFIED, got %s", eps[0].Source)
	}
	if len(res.SkipMessages) != 1 || !strings.Contains(res.SkipMessages[0], "UNVERIFIED") {
		t.Fatalf("expected skip reason exposed, got %+v", res.SkipMessages)
	}
}

func TestLinkSource_UnknownLangFallbackAndExtMismatch(t *testing.T) {
	root := t.TempDir()
	// Only a JS client exists; lang unknown -> full-ext fallback collects it,
	// but ext-mismatch safety net still rejects because .js is not in any known
	// lang's langExts? No: unknown lang uses sourceExts as extSet, so .js IS
	// allowed and it links (low-confidence). Verify low-confidence linking.
	os.WriteFile(filepath.Join(root, "client.js"), []byte(`export default { createSubscriber: () => {} }`), 0o644)
	eps := []Endpoint{{ID: "1", Handler: "createSubscriber"}}
	res := LinkSource(eps, root, "")
	if res.LangKnown {
		t.Fatalf("expected lang unknown")
	}
	// Under unknown lang the object-literal key is allowed (JS last-resort),
	// so it links as low-confidence; this preserves pre-Phase046 behavior.
	if res.Linked != 1 {
		t.Fatalf("expected 1 low-confidence link under unknown lang, got %+v", res)
	}
}
