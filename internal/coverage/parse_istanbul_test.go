package coverage

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseIstanbul_Valid(t *testing.T) {
	tmpDir := t.TempDir()
	f := filepath.Join(tmpDir, "coverage-final.json")
	content := `{
		"/app/handler.js": {
			"path": "/app/handler.js",
			"statementMap": {
				"0": {"start": {"line": 1, "column": 0}, "end": {"line": 3, "column": 1}},
				"1": {"start": {"line": 5, "column": 0}, "end": {"line": 7, "column": 1}}
			},
			"s": {"0": 5, "1": 0}
		}
	}`
	os.WriteFile(f, []byte(content), 0o644)

	blocks, err := ParseIstanbul(f)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(blocks) != 2 {
		t.Fatalf("expected 2 blocks, got %d", len(blocks))
	}
}

func TestParseIstanbul_FileNotFound(t *testing.T) {
	_, err := ParseIstanbul("/nonexistent/coverage-final.json")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestParseIstanbul_InvalidJSON(t *testing.T) {
	tmpDir := t.TempDir()
	f := filepath.Join(tmpDir, "coverage-final.json")
	os.WriteFile(f, []byte("INVALID"), 0o644)

	_, err := ParseIstanbul(f)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestParseIstanbul_EmptyFile(t *testing.T) {
	tmpDir := t.TempDir()
	f := filepath.Join(tmpDir, "coverage-final.json")
	os.WriteFile(f, []byte("{}"), 0o644)

	blocks, err := ParseIstanbul(f)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(blocks) != 0 {
		t.Fatalf("expected 0 blocks, got %d", len(blocks))
	}
}
