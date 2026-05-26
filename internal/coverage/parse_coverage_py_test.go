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

	executed, missing, err := ParseCoveragePy(f, "handler.py", 1, 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(executed) != 3 {
		t.Fatalf("expected 3 executed lines, got %d", len(executed))
	}
	if len(missing) != 2 {
		t.Fatalf("expected 2 missing lines, got %d", len(missing))
	}
}

func TestParseCoveragePy_FileNotFound(t *testing.T) {
	_, _, err := ParseCoveragePy("/nonexistent/cov.json", "handler.py", 1, 10)
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestParseCoveragePy_InvalidJSON(t *testing.T) {
	tmpDir := t.TempDir()
	f := filepath.Join(tmpDir, "cov.json")
	os.WriteFile(f, []byte("INVALID"), 0o644)

	_, _, err := ParseCoveragePy(f, "handler.py", 1, 10)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestParseCoveragePy_FileNotInReport(t *testing.T) {
	tmpDir := t.TempDir()
	f := filepath.Join(tmpDir, "cov.json")
	content := `{"files": {"/app/other.py": {"executed_lines": [1], "missing_lines": []}}}`
	os.WriteFile(f, []byte(content), 0o644)

	executed, missing, err := ParseCoveragePy(f, "handler.py", 1, 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if executed != nil {
		t.Fatalf("expected nil executed for missing file, got %v", executed)
	}
	if missing != nil {
		t.Fatalf("expected nil missing for missing file, got %v", missing)
	}
}

func TestParseCoveragePy_FilteredByRange(t *testing.T) {
	tmpDir := t.TempDir()
	f := filepath.Join(tmpDir, "cov.json")
	content := `{
		"files": {
			"handler.py": {
				"executed_lines": [1, 7, 20],
				"missing_lines": [2, 5, 15]
			}
		}
	}`
	os.WriteFile(f, []byte(content), 0o644)

	executed, missing, err := ParseCoveragePy(f, "handler.py", 1, 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Only lines 1 and 7 are executed within [1, 10]
	if len(executed) != 2 {
		t.Fatalf("expected 2 executed, got %d", len(executed))
	}
	// Only lines 2 and 5 are missing within [1, 10]
	if len(missing) != 2 {
		t.Fatalf("expected 2 missing, got %d", len(missing))
	}
}
