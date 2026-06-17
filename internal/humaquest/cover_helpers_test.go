package humaquest

import (
	"math"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/park-jun-woo/huma/internal/adapter"
	"github.com/park-jun-woo/reins/pkg/quest"
)

// ---------------------------------------------------------------------------
// coverageGround  (round-trips with decodeCoverage)
// ---------------------------------------------------------------------------

func TestCoverageGround_RoundTrip(t *testing.T) {
	cov := covWith(10, 80, 1, 5, 9)
	cov.Covered = 8
	ground, err := coverageGround(cov)
	if err != nil {
		t.Fatalf("coverageGround: %v", err)
	}
	got, present := decodeCoverage(ground)
	if !present {
		t.Fatal("decodeCoverage: not present")
	}
	if got.Total != 10 || got.Percent != 80 || got.Covered != 8 {
		t.Errorf("round-trip mismatch: %+v", got)
	}
	if !got.CoveredLines[5] {
		t.Error("CoveredLines not preserved")
	}
}

func TestCoverageGround_MarshalError(t *testing.T) {
	// A non-finite float is not representable in JSON → json.Marshal errors,
	// driving coverageGround's defensive error branch.
	_, err := coverageGround(&adapter.CoverageResult{Percent: math.Inf(1)})
	if err == nil {
		t.Fatal("want marshal error on non-finite Percent")
	}
}

func TestCoverageGround_Nil(t *testing.T) {
	ground, err := coverageGround(nil)
	if err != nil {
		t.Fatalf("coverageGround(nil): %v", err)
	}
	if ground != "null" {
		t.Errorf("ground = %q, want null", ground)
	}
}

// ---------------------------------------------------------------------------
// coveragePercent
// ---------------------------------------------------------------------------

func TestCoveragePercent(t *testing.T) {
	if got := coveragePercent(nil); got != 0 {
		t.Errorf("nil → %v, want 0", got)
	}
	if got := coveragePercent(&adapter.CoverageResult{Percent: 42.5}); got != 42.5 {
		t.Errorf("got %v, want 42.5", got)
	}
}

// ---------------------------------------------------------------------------
// loadCoverSession
// ---------------------------------------------------------------------------

func TestLoadCoverSession_Missing(t *testing.T) {
	_, err := loadCoverSession(filepath.Join(t.TempDir(), "absent.json"))
	if err == nil {
		t.Fatal("want error for missing session")
	}
	if !strings.Contains(err.Error(), "huma scan") {
		t.Errorf("error = %v, want actionable scan hint", err)
	}
}

func TestLoadCoverSession_Loads(t *testing.T) {
	dir := t.TempDir()
	path := seedCoverSession(t, dir, coverEP("a"), coverEP("b"))
	s, err := loadCoverSession(path)
	if err != nil {
		t.Fatalf("loadCoverSession: %v", err)
	}
	if len(s.Items) != 2 {
		t.Errorf("items = %d, want 2", len(s.Items))
	}
}

func TestLoadCoverSession_CorruptError(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "session.json")
	if err := os.WriteFile(path, []byte("{not json"), 0o644); err != nil {
		t.Fatal(err)
	}
	_, err := loadCoverSession(path)
	if err == nil {
		t.Fatal("want error for corrupt session")
	}
	// Not the missing-file message — a genuine parse error.
	if strings.Contains(err.Error(), "huma scan") {
		t.Errorf("corrupt file misreported as missing: %v", err)
	}
}

// ---------------------------------------------------------------------------
// newJSONLSink / jsonlSink.Emit
// ---------------------------------------------------------------------------

func TestNewJSONLSink_CreatesParentDir(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "nested", "deep", "out.jsonl")
	sink, err := newJSONLSink(path)
	if err != nil {
		t.Fatalf("newJSONLSink: %v", err)
	}
	if sink.path != path {
		t.Errorf("path = %q, want %q", sink.path, path)
	}
	if fi, err := os.Stat(filepath.Dir(path)); err != nil || !fi.IsDir() {
		t.Errorf("parent dir not created: %v", err)
	}
}

func TestNewJSONLSink_BareFilename(t *testing.T) {
	chdir(t, t.TempDir())
	sink, err := newJSONLSink("out.jsonl") // dir == "." → no mkdir
	if err != nil {
		t.Fatalf("newJSONLSink: %v", err)
	}
	if sink.path != "out.jsonl" {
		t.Errorf("path = %q", sink.path)
	}
}

func TestNewJSONLSink_MkdirError(t *testing.T) {
	dir := t.TempDir()
	// Create a FILE where the sink wants a parent directory.
	blocker := filepath.Join(dir, "blocker")
	if err := os.WriteFile(blocker, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	_, err := newJSONLSink(filepath.Join(blocker, "out.jsonl"))
	if err == nil {
		t.Fatal("want mkdir error when parent path is a file")
	}
}

func TestJSONLSink_Emit_AppendsLines(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "out.jsonl")
	sink, _ := newJSONLSink(path)

	it1 := &quest.Item{Key: "a", State: quest.PASS}
	it2 := &quest.Item{Key: "b", State: quest.DONE}
	if err := sink.Emit(it1); err != nil {
		t.Fatalf("Emit: %v", err)
	}
	if err := sink.Emit(it2); err != nil {
		t.Fatalf("Emit: %v", err)
	}

	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	lines := strings.Split(strings.TrimRight(string(b), "\n"), "\n")
	if len(lines) != 2 {
		t.Fatalf("lines = %d, want 2", len(lines))
	}
	if !strings.Contains(lines[0], `"key":"a"`) || !strings.Contains(lines[1], `"key":"b"`) {
		t.Errorf("content wrong: %v", lines)
	}
}

func TestJSONLSink_Emit_MarshalError(t *testing.T) {
	dir := t.TempDir()
	sink, _ := newJSONLSink(filepath.Join(dir, "out.jsonl"))
	// An Item carrying an invalid RawMessage Payload opens fine but fails to
	// json.Marshal, driving Emit's marshal-error branch.
	it := &quest.Item{Key: "a", Payload: []byte("{not json")}
	if err := sink.Emit(it); err == nil {
		t.Fatal("want marshal error for invalid payload")
	}
}

func TestJSONLSink_Emit_OpenError(t *testing.T) {
	dir := t.TempDir()
	// path is a directory → OpenFile for write fails.
	sink := &jsonlSink{path: dir}
	if err := sink.Emit(&quest.Item{Key: "a"}); err == nil {
		t.Fatal("want open error when path is a directory")
	}
}
