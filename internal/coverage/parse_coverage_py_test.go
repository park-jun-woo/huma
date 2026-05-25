package coverage

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseCoveragePy_Valid(t *testing.T) {
	tmpDir := t.TempDir()
	f := filepath.Join(tmpDir, "cov.json")
	content := `{
		"files": {
			"/app/handler.py": {
				"executed_lines": [1, 2, 3],
				"missing_lines": [4, 5]
			}
		}
	}`
	os.WriteFile(f, []byte(content), 0o644)

	lines, err := ParseCoveragePy(f, "handler.py", 1, 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(lines) != 2 {
		t.Fatalf("expected 2 missing lines, got %d", len(lines))
	}
}

func TestParseCoveragePy_FileNotFound(t *testing.T) {
	_, err := ParseCoveragePy("/nonexistent/cov.json", "handler.py", 1, 10)
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestParseCoveragePy_InvalidJSON(t *testing.T) {
	tmpDir := t.TempDir()
	f := filepath.Join(tmpDir, "cov.json")
	os.WriteFile(f, []byte("INVALID"), 0o644)

	_, err := ParseCoveragePy(f, "handler.py", 1, 10)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestParseCoveragePy_FileNotInReport(t *testing.T) {
	tmpDir := t.TempDir()
	f := filepath.Join(tmpDir, "cov.json")
	content := `{"files": {"/app/other.py": {"executed_lines": [1], "missing_lines": []}}}`
	os.WriteFile(f, []byte(content), 0o644)

	lines, err := ParseCoveragePy(f, "handler.py", 1, 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if lines != nil {
		t.Fatalf("expected nil for missing file, got %v", lines)
	}
}

func TestParseCoveragePy_FilteredByRange(t *testing.T) {
	tmpDir := t.TempDir()
	f := filepath.Join(tmpDir, "cov.json")
	content := `{
		"files": {
			"handler.py": {
				"executed_lines": [1],
				"missing_lines": [2, 5, 15]
			}
		}
	}`
	os.WriteFile(f, []byte(content), 0o644)

	lines, err := ParseCoveragePy(f, "handler.py", 1, 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Only lines 2 and 5 are within [1, 10]
	if len(lines) != 2 {
		t.Fatalf("expected 2, got %d", len(lines))
	}
}
