package config

import (
	"os"
	"path/filepath"
	"testing"
)

func chdir(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	orig, _ := os.Getwd()
	t.Cleanup(func() { os.Chdir(orig) })
	os.Chdir(dir)
	return dir
}

func TestLoadUnreachable_Missing(t *testing.T) {
	chdir(t)
	entries, err := LoadUnreachable()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(entries))
	}
}

func TestLoadUnreachable_DropsIncomplete(t *testing.T) {
	chdir(t)
	os.MkdirAll(".huma", 0o755)
	content := `- endpoint: POST /api/x
  status: 503
  reason: gateway timeout
  evidence: handler.go:88
- endpoint: GET /api/y
  status: 500
`
	os.WriteFile(filepath.Join(".huma", "unreachable.yaml"), []byte(content), 0o644)
	entries, err := LoadUnreachable()
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 valid entry (incomplete dropped), got %d", len(entries))
	}
	if !IsExempt(entries, "POST /api/x", 503) {
		t.Fatal("expected POST /api/x 503 to be exempt")
	}
	if IsExempt(entries, "GET /api/y", 500) {
		t.Fatal("incomplete entry must not be exempt")
	}
}
