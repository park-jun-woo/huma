package session

import (
	"math"
	"os"
	"path/filepath"
	"testing"

	"github.com/park-jun-woo/huma/internal/scanner"
)

func TestNew(t *testing.T) {
	s := New()
	if s == nil {
		t.Fatal("expected non-nil session")
	}
	if len(s.Entries) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(s.Entries))
	}
}

func TestLoadAndSave(t *testing.T) {
	tmpDir := t.TempDir()
	orig, _ := os.Getwd()
	t.Cleanup(func() { os.Chdir(orig) })
	os.Chdir(tmpDir)

	s := New()
	s.Merge([]scanner.Endpoint{
		{ID: "ep1", Method: "GET", Path: "/test"},
	})

	if err := s.Save(); err != nil {
		t.Fatalf("save: %v", err)
	}

	loaded, err := Load()
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if len(loaded.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(loaded.Entries))
	}
	if loaded.Entries[0].ID != "ep1" {
		t.Fatalf("expected ep1, got %s", loaded.Entries[0].ID)
	}
}

func TestLoad_NoSession(t *testing.T) {
	tmpDir := t.TempDir()
	orig, _ := os.Getwd()
	t.Cleanup(func() { os.Chdir(orig) })
	os.Chdir(tmpDir)

	_, err := Load()
	if err == nil {
		t.Fatal("expected error for missing session")
	}
}

func TestLoad_InvalidJSON(t *testing.T) {
	tmpDir := t.TempDir()
	orig, _ := os.Getwd()
	t.Cleanup(func() { os.Chdir(orig) })
	os.Chdir(tmpDir)

	sessDir := filepath.Join(tmpDir, ".huma")
	os.MkdirAll(sessDir, 0o755)
	os.WriteFile(filepath.Join(sessDir, "session.json"), []byte("invalid json"), 0o644)

	_, err := Load()
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestSave_Error(t *testing.T) {
	tmpDir := t.TempDir()
	orig, _ := os.Getwd()
	t.Cleanup(func() { os.Chdir(orig) })
	os.Chdir(tmpDir)

	// Create a file where .huma should be a directory
	os.WriteFile(filepath.Join(tmpDir, ".huma"), []byte("blocker"), 0o644)

	s := New()
	err := s.Save()
	if err == nil {
		t.Fatal("expected error when .huma is a file")
	}
}

func TestMerge_NewEndpoints(t *testing.T) {
	s := New()
	eps := []scanner.Endpoint{
		{ID: "ep1", Method: "GET", Path: "/a"},
		{ID: "ep2", Method: "POST", Path: "/b"},
	}
	s.Merge(eps)

	if len(s.Entries) != 2 {
		t.Fatalf("expected 2, got %d", len(s.Entries))
	}
	if s.Entries[0].Status != StatusTodo {
		t.Fatalf("expected TODO, got %s", s.Entries[0].Status)
	}
}

func TestMerge_PreservesStatus(t *testing.T) {
	s := New()
	s.Merge([]scanner.Endpoint{{ID: "ep1", Method: "GET", Path: "/a"}})
	s.MarkPass("ep1")

	// Re-merge with updated endpoint
	s.Merge([]scanner.Endpoint{{ID: "ep1", Method: "GET", Path: "/a/updated"}})

	if len(s.Entries) != 1 {
		t.Fatalf("expected 1, got %d", len(s.Entries))
	}
	if s.Entries[0].Status != StatusPass {
		t.Fatalf("expected PASS preserved, got %s", s.Entries[0].Status)
	}
	if s.Entries[0].Path != "/a/updated" {
		t.Fatalf("expected updated path, got %s", s.Entries[0].Path)
	}
}

func TestMerge_RemovesDeleted(t *testing.T) {
	s := New()
	s.Merge([]scanner.Endpoint{
		{ID: "ep1", Method: "GET", Path: "/a"},
		{ID: "ep2", Method: "POST", Path: "/b"},
	})

	// Re-merge with only ep1
	s.Merge([]scanner.Endpoint{{ID: "ep1", Method: "GET", Path: "/a"}})

	if len(s.Entries) != 1 {
		t.Fatalf("expected 1, got %d", len(s.Entries))
	}
}

func TestCurrent_ReturnsTodo(t *testing.T) {
	s := New()
	s.Merge([]scanner.Endpoint{
		{ID: "ep1", Method: "GET", Path: "/a"},
	})

	ep := s.Current()
	if ep == nil {
		t.Fatal("expected non-nil")
	}
	if ep.ID != "ep1" {
		t.Fatalf("expected ep1, got %s", ep.ID)
	}
}

func TestCurrent_ReturnsImprove(t *testing.T) {
	s := New()
	s.Merge([]scanner.Endpoint{
		{ID: "ep1", Method: "GET", Path: "/a"},
	})
	s.MarkImprove("ep1", 50)

	ep := s.Current()
	if ep == nil || ep.ID != "ep1" {
		t.Fatal("expected ep1 in IMPROVE status")
	}
}

func TestCurrent_SkipsPassAndDone(t *testing.T) {
	s := New()
	s.Merge([]scanner.Endpoint{
		{ID: "ep1", Method: "GET", Path: "/a"},
		{ID: "ep2", Method: "POST", Path: "/b"},
		{ID: "ep3", Method: "PUT", Path: "/c"},
	})
	s.MarkPass("ep1")
	s.MarkDone("ep2", 80)

	ep := s.Current()
	if ep == nil || ep.ID != "ep3" {
		t.Fatal("expected ep3")
	}
}

func TestCurrent_AllDone(t *testing.T) {
	s := New()
	s.Merge([]scanner.Endpoint{
		{ID: "ep1", Method: "GET", Path: "/a"},
	})
	s.MarkPass("ep1")

	if s.Current() != nil {
		t.Fatal("expected nil when all done")
	}
}

func TestCurrentEntry(t *testing.T) {
	s := New()
	s.Merge([]scanner.Endpoint{
		{ID: "ep1", Method: "GET", Path: "/a"},
	})

	entry := s.CurrentEntry()
	if entry == nil {
		t.Fatal("expected non-nil entry")
	}
	if entry.ID != "ep1" {
		t.Fatalf("expected ep1, got %s", entry.ID)
	}
}

func TestCurrentEntry_AllDone(t *testing.T) {
	s := New()
	s.Merge([]scanner.Endpoint{{ID: "ep1", Method: "GET", Path: "/a"}})
	s.MarkPass("ep1")

	if s.CurrentEntry() != nil {
		t.Fatal("expected nil")
	}
}

func TestMarkImprove(t *testing.T) {
	s := New()
	s.Merge([]scanner.Endpoint{{ID: "ep1", Method: "GET", Path: "/a"}})

	s.MarkImprove("ep1", 50)
	e := s.Entries[0]
	if e.Status != StatusImprove {
		t.Fatalf("expected IMPROVE, got %s", e.Status)
	}
	if e.Coverage != 50 {
		t.Fatalf("expected 50, got %f", e.Coverage)
	}
	if e.ImproveCount != 1 {
		t.Fatalf("expected 1, got %d", e.ImproveCount)
	}

	s.MarkImprove("ep1", 70)
	e = s.Entries[0]
	if e.PrevCoverage != 50 {
		t.Fatalf("expected prev 50, got %f", e.PrevCoverage)
	}
	if e.Coverage != 70 {
		t.Fatalf("expected 70, got %f", e.Coverage)
	}
	if e.ImproveCount != 2 {
		t.Fatalf("expected 2, got %d", e.ImproveCount)
	}
}

func TestMarkDone(t *testing.T) {
	s := New()
	s.Merge([]scanner.Endpoint{{ID: "ep1", Method: "GET", Path: "/a"}})

	s.MarkDone("ep1", 80)
	if s.Entries[0].Status != StatusDone {
		t.Fatalf("expected DONE, got %s", s.Entries[0].Status)
	}
	if s.Entries[0].Coverage != 80 {
		t.Fatalf("expected 80, got %f", s.Entries[0].Coverage)
	}
}

func TestMarkPass(t *testing.T) {
	s := New()
	s.Merge([]scanner.Endpoint{{ID: "ep1", Method: "GET", Path: "/a"}})

	s.MarkPass("ep1")
	if s.Entries[0].Status != StatusPass {
		t.Fatalf("expected PASS, got %s", s.Entries[0].Status)
	}
}

func TestSave_MarshalError(t *testing.T) {
	tmpDir := t.TempDir()
	orig, _ := os.Getwd()
	t.Cleanup(func() { os.Chdir(orig) })
	os.Chdir(tmpDir)

	s := New()
	s.Merge([]scanner.Endpoint{{ID: "ep1", Method: "GET", Path: "/a"}})
	// Set Coverage to NaN which causes json.MarshalIndent to fail
	s.Entries[0].Coverage = math.NaN()

	err := s.Save()
	if err == nil {
		t.Fatal("expected error from MarshalIndent with NaN")
	}
}

func TestStats(t *testing.T) {
	s := New()
	s.Merge([]scanner.Endpoint{
		{ID: "ep1", Method: "GET", Path: "/a"},
		{ID: "ep2", Method: "POST", Path: "/b"},
		{ID: "ep3", Method: "PUT", Path: "/c"},
		{ID: "ep4", Method: "DELETE", Path: "/d"},
	})
	s.MarkPass("ep1")
	s.MarkDone("ep2", 80)
	s.MarkImprove("ep3", 50)

	total, pass, todo := s.Stats()
	if total != 4 {
		t.Fatalf("expected total 4, got %d", total)
	}
	if pass != 2 {
		t.Fatalf("expected pass 2, got %d", pass)
	}
	if todo != 2 {
		t.Fatalf("expected todo 2, got %d", todo)
	}
}
