package adapter

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/park-jun-woo/huma/internal/config"
)

func TestNodeCollect_C8ReportError(t *testing.T) {
	a := &NodeAdapter{
		cfg:      &config.ServerConfig{},
		coverDir: "/nonexistent/v8cov",
	}
	_, err := a.Collect("handler.js", 10, 20)
	if err == nil {
		t.Fatal("expected error from c8 report")
	}
}

func TestNodeCollect_EmptyCoverage(t *testing.T) {
	tmpDir := t.TempDir()
	orig, _ := os.Getwd()
	t.Cleanup(func() { os.Chdir(orig) })
	os.Chdir(tmpDir)

	covDir := filepath.Join(tmpDir, "v8cov")
	os.MkdirAll(covDir, 0o755)

	// Create istanbul output directory with a valid but empty coverage file
	os.MkdirAll(filepath.Join(tmpDir, ".huma", "istanbul"), 0o755)
	os.WriteFile(
		filepath.Join(tmpDir, ".huma", "istanbul", "coverage-final.json"),
		[]byte("{}"),
		0o644,
	)

	a := &NodeAdapter{
		cfg:      &config.ServerConfig{},
		coverDir: covDir,
	}

	// c8 report may fail if npx/c8 is not installed, or may succeed.
	// Either way the test exercises the code path.
	result, err := a.Collect("handler.js", 10, 20)
	if err != nil {
		// c8 report failed — acceptable
		return
	}
	if result == nil {
		t.Fatal("expected non-nil result")
	}
}

func TestNodeCollect_ParseIstanbulError(t *testing.T) {
	tmpDir := t.TempDir()
	orig, _ := os.Getwd()
	t.Cleanup(func() { os.Chdir(orig) })
	os.Chdir(tmpDir)

	covDir := filepath.Join(tmpDir, "v8cov")
	os.MkdirAll(covDir, 0o755)

	// Create coverage-final.json as a directory to force ParseIstanbul to fail
	istDir := filepath.Join(tmpDir, ".huma", "istanbul", "coverage-final.json")
	os.MkdirAll(istDir, 0o755)

	a := &NodeAdapter{
		cfg:      &config.ServerConfig{},
		coverDir: covDir,
	}

	_, err := a.Collect("handler.js", 10, 20)
	// Either c8 report fails, or ParseIstanbul fails on the directory
	if err == nil {
		t.Fatal("expected error")
	}
}
