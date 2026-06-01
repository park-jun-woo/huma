package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/park-jun-woo/huma/internal/scanner"
)

func TestAnalyzeBranches(t *testing.T) {
	// no source → nil
	if got := analyzeBranches(&scanner.Endpoint{}, "go"); got != nil {
		t.Errorf("no source → nil, got %v", got)
	}
	// unknown lang → nil analyzer → nil
	dir := t.TempDir()
	src := filepath.Join(dir, "h.go")
	os.WriteFile(src, []byte("package main\nimport \"net/http\"\nfunc H(c interface{}) {\n\tc.JSON(http.StatusOK, nil)\n}\n"), 0o644)
	if got := analyzeBranches(&scanner.Endpoint{Source: src, Handler: "H", Line: 3}, "cobol"); got != nil {
		t.Errorf("unknown lang → nil, got %v", got)
	}
	// valid go source → branches
	got := analyzeBranches(&scanner.Endpoint{Source: src, Handler: "H", Line: 3}, "go")
	if len(got) != 1 || got[0].Status != 200 {
		t.Errorf("expected 1 branch status 200, got %v", got)
	}
	// nonexistent source file → Analyze errors → nil
	if got := analyzeBranches(&scanner.Endpoint{Source: filepath.Join(dir, "nope.go"), Handler: "H", Line: 1}, "go"); got != nil {
		t.Errorf("missing source → nil, got %v", got)
	}
}
