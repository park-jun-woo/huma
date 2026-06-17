package humaquest

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/park-jun-woo/huma/internal/scanner"
)

const scanOpenAPI = `
openapi: '3.0.0'
info:
  title: Scan Test API
paths:
  /api/v1/login:
    post:
      operationId: Login
  /api/v1/me:
    get:
      operationId: GetMe
`

func epByMethodPath(eps []scanner.Endpoint) map[string]scanner.Endpoint {
	m := make(map[string]scanner.Endpoint, len(eps))
	for _, ep := range eps {
		m[ep.Method+" "+ep.Path] = ep
	}
	return m
}

func TestScanEndpoints_OpenAPIFile(t *testing.T) {
	path := writeTempFile(t, "openapi.yaml", scanOpenAPI)

	eps, err := scanEndpoints(path)
	if err != nil {
		t.Fatalf("scanEndpoints: %v", err)
	}
	if len(eps) != 2 {
		t.Fatalf("want 2 endpoints, got %d", len(eps))
	}
	by := epByMethodPath(eps)
	if _, ok := by["POST /api/v1/login"]; !ok {
		t.Errorf("missing POST /api/v1/login; got %v", by)
	}
	if _, ok := by["GET /api/v1/me"]; !ok {
		t.Errorf("missing GET /api/v1/me; got %v", by)
	}
	for _, ep := range eps {
		if ep.ID == "" {
			t.Errorf("endpoint %s %s has empty ID", ep.Method, ep.Path)
		}
	}
}

func TestScanEndpoints_Stdin(t *testing.T) {
	// Exercise the "-" branch by replacing os.Stdin with a pipe carrying OpenAPI.
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	go func() {
		_, _ = w.Write([]byte(scanOpenAPI))
		_ = w.Close()
	}()

	orig := os.Stdin
	os.Stdin = r
	defer func() { os.Stdin = orig }()

	eps, err := scanEndpoints("-")
	if err != nil {
		t.Fatalf("scanEndpoints(-): %v", err)
	}
	if len(eps) != 2 {
		t.Fatalf("want 2 endpoints from stdin, got %d", len(eps))
	}
}

func TestScanEndpoints_EdgeFunctionsDir(t *testing.T) {
	dir := t.TempDir()
	mkEdge := func(name, body string) {
		p := filepath.Join(dir, name, "index.ts")
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	mkEdge("hello", `if (req.method === 'GET') { return ok() }`)
	mkEdge("create", `if (req.method === 'POST') { return create() }`)

	eps, err := scanEndpoints(dir)
	if err != nil {
		t.Fatalf("scanEndpoints(dir): %v", err)
	}
	if len(eps) != 2 {
		t.Fatalf("want 2 edge endpoints, got %d", len(eps))
	}
	by := epByMethodPath(eps)
	if _, ok := by["GET /functions/v1/hello"]; !ok {
		t.Errorf("missing GET /functions/v1/hello; got %v", by)
	}
	if _, ok := by["POST /functions/v1/create"]; !ok {
		t.Errorf("missing POST /functions/v1/create; got %v", by)
	}
}

func TestScanEndpoints_LinkSource(t *testing.T) {
	// Build a temp source tree plus an endpoint-list JSON (handlers, no source),
	// then pass a link-source root so handlers are resolved to file:line in place.
	root := t.TempDir()
	src := filepath.Join(root, "handlers.go")
	if err := os.WriteFile(src, []byte("package main\n\nfunc Login(c interface{}) {}\nfunc GetMe(c interface{}) {}\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	// Endpoint-list JSON (not OpenAPI) carrying handler names without source.
	list := `[
  {"method":"POST","path":"/login","handler":"Login"},
  {"method":"GET","path":"/me","handler":"GetMe"}
]`
	epPath := filepath.Join(root, "endpoints.json")
	if err := os.WriteFile(epPath, []byte(list), 0o644); err != nil {
		t.Fatal(err)
	}

	eps, err := scanEndpoints(epPath, root)
	if err != nil {
		t.Fatalf("scanEndpoints with link-source: %v", err)
	}
	if len(eps) != 2 {
		t.Fatalf("want 2 endpoints, got %d", len(eps))
	}

	linked := 0
	for _, ep := range eps {
		if ep.Source != "" && ep.Line > 0 {
			linked++
		}
	}
	if linked != 2 {
		t.Fatalf("expected both handlers linked to source, got %d linked: %+v", linked, eps)
	}
}

func TestScanEndpoints_NoSourceErrorsE01(t *testing.T) {
	// Empty args + no openapi.yaml in cwd (the package dir) → E-01.
	if _, err := os.Stat("openapi.yaml"); err == nil {
		t.Skip("unexpected openapi.yaml in package dir")
	}
	_, err := scanEndpoints()
	if err == nil {
		t.Fatal("expected E-01 error when no source given and none discoverable")
	}
	if !strings.Contains(err.Error(), "E-01") {
		t.Errorf("error %q does not mention E-01", err)
	}
}

func TestScanEndpoints_ReadError(t *testing.T) {
	_, err := scanEndpoints(filepath.Join(t.TempDir(), "nope.yaml"))
	if err == nil {
		t.Fatal("expected read error for missing file")
	}
	if !strings.Contains(err.Error(), "read input") {
		t.Errorf("error %q does not mention read input", err)
	}
}

func TestScanEndpoints_ParseError(t *testing.T) {
	// Existing file whose content is neither OpenAPI nor a valid endpoint list.
	path := writeTempFile(t, "bad.yaml", "}{ not json not yaml }{")
	_, err := scanEndpoints(path)
	if err == nil {
		t.Fatal("expected parse error for malformed input")
	}
	if !strings.Contains(err.Error(), "parse endpoints") {
		t.Errorf("error %q does not mention parse endpoints", err)
	}
}

// chdirTemp switches into a fresh temp dir for the duration of the test.
func chdirTemp(t *testing.T) string {
	t.Helper()
	orig, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	dir := t.TempDir()
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(orig) })
	return dir
}

func TestScanEndpoints_FindsOpenAPIInCwd(t *testing.T) {
	// Empty args + an openapi.yaml in cwd → FindOpenAPIFile resolves it.
	dir := chdirTemp(t)
	if err := os.WriteFile(filepath.Join(dir, "openapi.yaml"), []byte(scanOpenAPI), 0o644); err != nil {
		t.Fatal(err)
	}
	eps, err := scanEndpoints()
	if err != nil {
		t.Fatalf("scanEndpoints(): %v", err)
	}
	if len(eps) != 2 {
		t.Fatalf("want 2 endpoints, got %d", len(eps))
	}
}

func TestScanEndpoints_LoadConfigError(t *testing.T) {
	// A corrupt manifest.yaml in cwd makes config.Load fail with a non-ErrNoManifest
	// error, which scanEndpoints surfaces as "load config".
	dir := chdirTemp(t)
	if err := os.WriteFile(filepath.Join(dir, "manifest.yaml"), []byte("{{invalid yaml"), 0o644); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(dir, "openapi.yaml")
	if err := os.WriteFile(path, []byte(scanOpenAPI), 0o644); err != nil {
		t.Fatal(err)
	}
	_, err := scanEndpoints(path)
	if err == nil {
		t.Fatal("expected load config error from corrupt manifest")
	}
	if !strings.Contains(err.Error(), "load config") {
		t.Errorf("error %q does not mention load config", err)
	}
}

func TestScanEndpoints_LinkSourceUsesManifestLang(t *testing.T) {
	// With a valid manifest present, the link-source path reads backend.lang from
	// config (the cfg != nil branch).
	dir := chdirTemp(t)
	manifest := "apiVersion: v1\nbackend:\n  lang: go\n"
	if err := os.WriteFile(filepath.Join(dir, "manifest.yaml"), []byte(manifest), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "handlers.go"),
		[]byte("package main\n\nfunc Login(c interface{}) {}\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	list := `[{"method":"POST","path":"/login","handler":"Login"}]`
	if err := os.WriteFile(filepath.Join(dir, "endpoints.json"), []byte(list), 0o644); err != nil {
		t.Fatal(err)
	}
	eps, err := scanEndpoints("endpoints.json", ".")
	if err != nil {
		t.Fatalf("scanEndpoints: %v", err)
	}
	if len(eps) != 1 {
		t.Fatalf("want 1 endpoint, got %d", len(eps))
	}
	if eps[0].Source == "" || eps[0].Line == 0 {
		t.Fatalf("expected Login linked via manifest lang, got %+v", eps[0])
	}
}

func TestScanEndpoints_EdgeFunctionsError(t *testing.T) {
	// A directory that stats fine but cannot be read (perm 000) makes
	// ParseEdgeFunctions fail; scanEndpoints wraps it as "scan edge functions".
	dir := filepath.Join(t.TempDir(), "locked")
	if err := os.Mkdir(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(dir, 0o000); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chmod(dir, 0o755) })

	// Skip if running with privileges that bypass the permission bits.
	if f, err := os.ReadDir(dir); err == nil {
		_ = f
		t.Skip("directory readable despite perm 000 (privileged env)")
	}

	_, err := scanEndpoints(dir)
	if err == nil {
		t.Fatal("expected scan edge functions error")
	}
	if !strings.Contains(err.Error(), "scan edge functions") {
		t.Errorf("error %q does not mention scan edge functions", err)
	}
}
