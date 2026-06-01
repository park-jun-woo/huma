package scanner

import (
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"testing"
)

func TestAllowedExts(t *testing.T) {
	exts, known := allowedExts("go")
	if !known || !exts[".go"] || exts[".py"] {
		t.Errorf("go → %v known=%v", exts, known)
	}
	// unknown lang → full fallback, false
	exts, known = allowedExts("cobol")
	if known {
		t.Error("unknown lang should be known=false")
	}
	if !exts[".go"] || !exts[".py"] || !exts[".rs"] {
		t.Errorf("fallback should contain all source exts, got %v", exts)
	}
	// empty lang → fallback
	_, known = allowedExts("")
	if known {
		t.Error("empty lang → known=false")
	}
}

func TestLangLabel(t *testing.T) {
	if langLabel("") != "backend.lang=unknown" {
		t.Error("empty → unknown")
	}
	if langLabel("go") != "backend.lang=go" {
		t.Error("go label")
	}
}

func TestAmbiguousMessage(t *testing.T) {
	matches := []handlerMatch{{File: "a.go", Line: 1}, {File: "b.go", Line: 2}}
	msg := ambiguousMessage("GetUser", matches, "/root", "go")
	for _, want := range []string{"GetUser", "2 candidates", "/root", "backend.lang=go", "a.go, b.go"} {
		if !strings.Contains(msg, want) {
			t.Errorf("missing %q in %q", want, msg)
		}
	}
}

func TestExtMismatchMessage(t *testing.T) {
	msg := extMismatchMessage("GetUser", "h.ts", "/root", "go")
	for _, want := range []string{"GetUser", ".ts", "/root", "backend.lang=go"} {
		if !strings.Contains(msg, want) {
			t.Errorf("missing %q in %q", want, msg)
		}
	}
}

func TestClassifyMatches(t *testing.T) {
	extSet := map[string]bool{".go": true}
	// not found
	ep := &Endpoint{Handler: "H"}
	if oc, _ := classifyMatches(ep, nil, "/r", "go", extSet); oc != outcomeNotFound {
		t.Error("no matches → not found")
	}
	// ambiguous
	if oc, msg := classifyMatches(ep, []handlerMatch{{File: "a.go"}, {File: "b.go"}}, "/r", "go", extSet); oc != outcomeAmbiguous || msg == "" {
		t.Error("2 matches → ambiguous")
	}
	// ext mismatch
	if oc, msg := classifyMatches(ep, []handlerMatch{{File: "a.ts", Line: 3}}, "/r", "go", extSet); oc != outcomeExtMismatch || msg == "" {
		t.Error("ts under go → ext mismatch")
	}
	// linked → sets Source/Line
	ep = &Endpoint{Handler: "H"}
	if oc, _ := classifyMatches(ep, []handlerMatch{{File: "a.go", Line: 7}}, "/r", "go", extSet); oc != outcomeLinked {
		t.Error("single allowed → linked")
	}
	if ep.Source != "a.go" || ep.Line != 7 {
		t.Errorf("linked should set Source/Line, got %s:%d", ep.Source, ep.Line)
	}
}

func TestLinkEndpoint(t *testing.T) {
	extSet := map[string]bool{".go": true}
	// already linked → noop
	if oc, _ := linkEndpoint(&Endpoint{Source: "x.go", Handler: "H"}, nil, "/r", "go", extSet); oc != outcomeNoop {
		t.Error("already linked → noop")
	}
	// no handler → noop
	if oc, _ := linkEndpoint(&Endpoint{}, nil, "/r", "go", extSet); oc != outcomeNoop {
		t.Error("no handler → noop")
	}
	// handler, no files → not found
	if oc, _ := linkEndpoint(&Endpoint{Handler: "H"}, nil, "/r", "go", extSet); oc != outcomeNotFound {
		t.Error("no files → not found")
	}
}

func TestCollectRawEndpoints(t *testing.T) {
	raw := []rawEndpoint{
		{Method: "GET", Path: "/a", Handler: "A"},
		{Method: "", Path: ""}, // likely dropped by parseRawEndpoint
	}
	eps := collectRawEndpoints(raw)
	if len(eps) < 1 {
		t.Fatalf("expected at least 1 endpoint, got %d", len(eps))
	}
	if eps[0].Method != "GET" {
		t.Errorf("first endpoint method = %s", eps[0].Method)
	}
}

func TestCollectResponseCodes(t *testing.T) {
	codes := collectResponseCodes(map[string]interface{}{"200": nil, "404": nil, "default": nil})
	got := map[int]bool{}
	for _, c := range codes {
		got[c] = true
	}
	if !got[200] || !got[404] {
		t.Errorf("expected 200,404, got %v", codes)
	}
	// nil/unknown type → nil
	if c := collectResponseCodes(42); c != nil {
		t.Errorf("non-map → nil, got %v", c)
	}
}

func TestResponseMapKeys(t *testing.T) {
	if k := responseMapKeys(map[string]interface{}{"a": 1, "b": 2}); len(k) != 2 {
		t.Errorf("string map → 2 keys, got %d", len(k))
	}
	if k := responseMapKeys(map[interface{}]interface{}{1: "a"}); len(k) != 1 {
		t.Errorf("iface map → 1 key, got %d", len(k))
	}
	if k := responseMapKeys("x"); k != nil {
		t.Errorf("non-map → nil, got %v", k)
	}
}

func TestToStatusCode(t *testing.T) {
	cases := []struct {
		in   interface{}
		want int
	}{
		{200, 200},
		{float64(404), 404},
		{"201", 201},
		{"notanumber", 0},
		{true, 0},
	}
	for _, c := range cases {
		if got := toStatusCode(c.in); got != c.want {
			t.Errorf("toStatusCode(%v) = %d, want %d", c.in, got, c.want)
		}
	}
}

func TestLookupResponse(t *testing.T) {
	strMap := map[string]interface{}{"200": map[string]interface{}{"description": "ok"}}
	if obj, ok := lookupResponse(strMap, 200); !ok || obj["description"] != "ok" {
		t.Error("string-keyed lookup 200")
	}
	if _, ok := lookupResponse(strMap, 404); ok {
		t.Error("missing code → false")
	}
	ifaceMap := map[interface{}]interface{}{200: map[string]interface{}{"x": 1}}
	if _, ok := lookupResponse(ifaceMap, 200); !ok {
		t.Error("iface-keyed lookup 200")
	}
	if _, ok := lookupResponse(42, 200); ok {
		t.Error("non-map → false")
	}
}

func TestLookupResponseIface(t *testing.T) {
	m := map[interface{}]interface{}{
		200:        map[string]interface{}{"x": 1},
		"bad":      "notamap",
	}
	if _, ok := lookupResponseIface(m, 200); !ok {
		t.Error("expected 200 found")
	}
	if _, ok := lookupResponseIface(m, 500); ok {
		t.Error("500 not present → false")
	}
}

func TestIntField(t *testing.T) {
	m := map[string]interface{}{"i": 5, "f": float64(7), "s": "x"}
	if intField(m, "i") != 5 {
		t.Error("int")
	}
	if intField(m, "f") != 7 {
		t.Error("float64")
	}
	if intField(m, "s") != 0 {
		t.Error("wrong type → 0")
	}
	if intField(m, "missing") != 0 {
		t.Error("missing → 0")
	}
}

func TestStringField(t *testing.T) {
	m := map[string]interface{}{"s": "hi", "i": 5}
	if stringField(m, "s") != "hi" {
		t.Error("string")
	}
	if stringField(m, "i") != "" {
		t.Error("wrong type → empty")
	}
	if stringField(m, "missing") != "" {
		t.Error("missing → empty")
	}
}

func TestNormalizeSymbol(t *testing.T) {
	cases := map[string]string{
		"createSubscriber":  "createsubscriber",
		"CreateSubscriber":  "createsubscriber",
		"create_subscriber": "createsubscriber",
		"Get-User-123":      "getuser123",
		"":                  "",
	}
	for in, want := range cases {
		if got := normalizeSymbol(in); got != want {
			t.Errorf("normalizeSymbol(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestGenerateOperationID(t *testing.T) {
	cases := []struct {
		method, path, want string
	}{
		{"GET", "/users/{id}", "GET_users_id"},
		{"POST", "/orgs", "POST_orgs"},
	}
	for _, c := range cases {
		if got := generateOperationID(c.method, c.path); got != c.want {
			t.Errorf("generateOperationID(%q,%q) = %q, want %q", c.method, c.path, got, c.want)
		}
	}
}

func TestHandlerDefPatterns(t *testing.T) {
	for ext, wantLen := range map[string]int{
		".go": 1, ".py": 1, ".rs": 1, ".java": 1, ".cs": 1, ".php": 1,
	} {
		if got := handlerDefPatterns("x" + ext); len(got) != wantLen {
			t.Errorf("%s → %d patterns, want %d", ext, len(got), wantLen)
		}
	}
	if got := handlerDefPatterns("x.js"); len(got) == 0 {
		t.Error(".js should have patterns")
	}
	if got := handlerDefPatterns("x.txt"); got != nil {
		t.Errorf("unknown ext → nil, got %v", got)
	}
}

func TestFindEntryFile(t *testing.T) {
	dir := t.TempDir()
	if findEntryFile(dir) != "" {
		t.Error("empty dir → empty")
	}
	os.WriteFile(filepath.Join(dir, "index.ts"), []byte("x"), 0o644)
	if findEntryFile(dir) != "index.ts" {
		t.Error("index.ts present → index.ts")
	}
}

func TestExtractMethods(t *testing.T) {
	// no methods → default POST
	if got := extractMethods("const x = 1;"); !reflect.DeepEqual(got, []string{"POST"}) {
		t.Errorf("no match → [POST], got %v", got)
	}
	content := `if (req.method === "GET") {}
switch (m) { case "DELETE": break; }`
	got := extractMethods(content)
	set := map[string]bool{}
	for _, m := range got {
		set[m] = true
	}
	if !set["GET"] || !set["DELETE"] {
		t.Errorf("expected GET and DELETE, got %v", got)
	}
}

func TestCollectLineMethods(t *testing.T) {
	re := regexp.MustCompile(`method\s*===?\s*['"](\w+)['"]`)
	seen := map[string]bool{}
	var methods []string
	collectLineMethods(`req.method === "GET"`, re, seen, &methods)
	collectLineMethods(`req.method === "GET"`, re, seen, &methods) // dup → ignored
	collectLineMethods(`req.method === "FOObar"`, re, seen, &methods) // not allowed → ignored
	if !reflect.DeepEqual(methods, []string{"GET"}) {
		t.Errorf("got %v, want [GET]", methods)
	}
}

func TestFindHandlerDef(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "h.go")
	os.WriteFile(f, []byte("package main\nfunc GetUser(c C) {}\n"), 0o644)
	matches := findHandlerDef([]string{f}, "getUser")
	if len(matches) != 1 || matches[0].Line != 2 {
		t.Errorf("expected match at line 2, got %v", matches)
	}
	// empty handler name → nil
	if got := findHandlerDef([]string{f}, ""); got != nil {
		t.Errorf("empty target → nil, got %v", got)
	}
	// no match
	if got := findHandlerDef([]string{f}, "Nonexistent"); got != nil {
		t.Errorf("no match → nil, got %v", got)
	}
}

func TestScanForHandler(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "h.go")
	os.WriteFile(f, []byte("package main\n\nfunc CreateUser(c C) {}\n"), 0o644)
	if got := scanForHandler(f, normalizeSymbol("create_user")); got != 3 {
		t.Errorf("expected line 3, got %d", got)
	}
	// unknown extension → 0
	txt := filepath.Join(dir, "h.txt")
	os.WriteFile(txt, []byte("func CreateUser"), 0o644)
	if got := scanForHandler(txt, "createuser"); got != 0 {
		t.Errorf("unknown ext → 0, got %d", got)
	}
	// missing file → 0
	if got := scanForHandler(filepath.Join(dir, "nope.go"), "x"); got != 0 {
		t.Errorf("missing file → 0, got %d", got)
	}
}

func TestScanForPattern(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "x.txt")
	os.WriteFile(f, []byte("line one\ntarget here\nline three\n"), 0o644)
	re := regexp.MustCompile(`target`)
	if got := scanForPattern(f, re); got != 2 {
		t.Errorf("expected line 2, got %d", got)
	}
	if got := scanForPattern(filepath.Join(dir, "nope.txt"), re); got != 0 {
		t.Errorf("missing file → 0, got %d", got)
	}
	noMatch := regexp.MustCompile(`zzz`)
	if got := scanForPattern(f, noMatch); got != 0 {
		t.Errorf("no match → 0, got %d", got)
	}
}

func TestCollectSourceFiles(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "a.go"), []byte("x"), 0o644)
	os.WriteFile(filepath.Join(dir, "b.py"), []byte("x"), 0o644)
	os.MkdirAll(filepath.Join(dir, "sub"), 0o755)
	os.WriteFile(filepath.Join(dir, "sub", "c.go"), []byte("x"), 0o644)

	// lang=go → only .go files (recursive)
	files := collectSourceFiles(dir, "go")
	if len(files) != 2 {
		t.Errorf("expected 2 .go files, got %d: %v", len(files), files)
	}
	for _, f := range files {
		if filepath.Ext(f) != ".go" {
			t.Errorf("non-go file collected: %s", f)
		}
	}
	// unknown lang → all source files (3)
	all := collectSourceFiles(dir, "")
	if len(all) != 3 {
		t.Errorf("unknown lang → 3 files, got %d", len(all))
	}
}

func TestExtractChildFields(t *testing.T) {
	childProps := map[string]interface{}{
		"name": map[string]interface{}{"type": "string"},
		"age":  map[string]interface{}{"type": "integer"},
		"bad":  "notamap",
	}
	fields := extractChildFields("user", childProps, nil)
	if len(fields) != 2 {
		t.Fatalf("expected 2 fields (bad skipped), got %d", len(fields))
	}
	paths := map[string]string{}
	for _, f := range fields {
		paths[f.Path] = f.Type
	}
	if paths["$.user.name"] != "string" || paths["$.user.age"] != "integer" {
		t.Errorf("unexpected fields: %v", paths)
	}
}
