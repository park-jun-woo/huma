package scanner

import (
	"os"
	"path/filepath"
	"sort"
	"testing"
)

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestParseEdgeFunctions_Basic(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "func-a", "index.ts"),
		`if (req.method === 'GET') { return getHandler() }`)
	writeFile(t, filepath.Join(dir, "func-b", "index.ts"),
		`if (req.method === 'POST') { return postHandler() }`)

	eps, err := ParseEdgeFunctions(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(eps) != 2 {
		t.Fatalf("want 2 endpoints, got %d", len(eps))
	}

	sort.Slice(eps, func(i, j int) bool { return eps[i].Path < eps[j].Path })

	if eps[0].Method != "GET" || eps[0].Path != "/functions/v1/func-a" {
		t.Errorf("eps[0] = %s %s", eps[0].Method, eps[0].Path)
	}
	if eps[1].Method != "POST" || eps[1].Path != "/functions/v1/func-b" {
		t.Errorf("eps[1] = %s %s", eps[1].Method, eps[1].Path)
	}
	if eps[0].Handler != "func-a" {
		t.Errorf("handler = %s, want func-a", eps[0].Handler)
	}
}

func TestParseEdgeFunctions_SkipUnderscore(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "_shared", "cors.ts"),
		`export const cors = {}`)
	writeFile(t, filepath.Join(dir, "hello", "index.ts"),
		`if (req.method === 'GET') {}`)

	eps, err := ParseEdgeFunctions(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(eps) != 1 {
		t.Fatalf("want 1 endpoint, got %d", len(eps))
	}
	if eps[0].Path != "/functions/v1/hello" {
		t.Errorf("path = %s", eps[0].Path)
	}
}

func TestParseEdgeFunctions_DefaultPost(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "webhook", "index.ts"),
		`const data = await req.json()
return new Response(JSON.stringify({ ok: true }))`)

	eps, err := ParseEdgeFunctions(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(eps) != 1 {
		t.Fatalf("want 1 endpoint, got %d", len(eps))
	}
	if eps[0].Method != "POST" {
		t.Errorf("method = %s, want POST", eps[0].Method)
	}
}

func TestParseEdgeFunctions_MultiMethod(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "items", "index.ts"),
		`if (req.method === 'GET') { return list() }
if (req.method === 'POST') { return create() }`)

	eps, err := ParseEdgeFunctions(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(eps) != 2 {
		t.Fatalf("want 2 endpoints, got %d", len(eps))
	}

	methods := map[string]bool{}
	for _, ep := range eps {
		methods[ep.Method] = true
	}
	if !methods["GET"] || !methods["POST"] {
		t.Errorf("methods = %v", methods)
	}
}

func TestParseEdgeFunctions_SwitchCase(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "resource", "index.ts"),
		`switch (req.method) {
  case 'GET':
    return handleGet()
  case 'POST':
    return handlePost()
}`)

	eps, err := ParseEdgeFunctions(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(eps) != 2 {
		t.Fatalf("want 2 endpoints, got %d", len(eps))
	}

	methods := map[string]bool{}
	for _, ep := range eps {
		methods[ep.Method] = true
	}
	if !methods["GET"] || !methods["POST"] {
		t.Errorf("methods = %v", methods)
	}
}

func TestParseEdgeFunctions_SkipOptions(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "api", "index.ts"),
		`if (req.method === 'OPTIONS') { return corsResponse() }
if (req.method === 'GET') { return getData() }`)

	eps, err := ParseEdgeFunctions(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(eps) != 1 {
		t.Fatalf("want 1 endpoint, got %d", len(eps))
	}
	if eps[0].Method != "GET" {
		t.Errorf("method = %s, want GET", eps[0].Method)
	}
}

func TestParseEdgeFunctions_SkipNonMethod(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "status", "index.ts"),
		`switch (record.status) {
  case 'active':
    break
  case 'inactive':
    break
}
if (req.method === 'GET') { return ok() }`)

	eps, err := ParseEdgeFunctions(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(eps) != 1 {
		t.Fatalf("want 1 endpoint, got %d", len(eps))
	}
	if eps[0].Method != "GET" {
		t.Errorf("method = %s, want GET", eps[0].Method)
	}
}

func TestParseEdgeFunctions_EmptyDir(t *testing.T) {
	dir := t.TempDir()

	eps, err := ParseEdgeFunctions(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(eps) != 0 {
		t.Fatalf("want 0 endpoints, got %d", len(eps))
	}
}

func TestParseEdgeFunctions_NoEntryFile(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "lib", "utils.ts"),
		`export function helper() {}`)

	eps, err := ParseEdgeFunctions(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(eps) != 0 {
		t.Fatalf("want 0 endpoints, got %d", len(eps))
	}
}

func TestParseEdgeFunctions_NegativeComparison(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "check", "index.ts"),
		`if (req.method !== 'POST') { return methodNotAllowed() }
const data = await req.json()`)

	eps, err := ParseEdgeFunctions(dir)
	if err != nil {
		t.Fatal(err)
	}
	// Negative comparison should not extract the method; default POST applies.
	if len(eps) != 1 {
		t.Fatalf("want 1 endpoint, got %d", len(eps))
	}
	if eps[0].Method != "POST" {
		t.Errorf("method = %s, want POST (default)", eps[0].Method)
	}
}
