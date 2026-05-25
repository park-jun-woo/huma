package adapter

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/park-jun-woo/huma/internal/config"
)

func TestPythonCollect_MkdirAllError(t *testing.T) {
	tmpDir := t.TempDir()
	orig, _ := os.Getwd()
	t.Cleanup(func() { os.Chdir(orig) })
	os.Chdir(tmpDir)

	// Create a file at .huma to block MkdirAll
	os.WriteFile(filepath.Join(tmpDir, ".huma"), []byte("block"), 0o644)

	a := &PythonAdapter{cfg: &config.ServerConfig{}}
	_, err := a.Collect("handler.py", 1, 10)
	if err == nil {
		t.Fatal("expected error from MkdirAll")
	}
}

func TestPythonCollect_CoverageCommandError(t *testing.T) {
	tmpDir := t.TempDir()
	orig, _ := os.Getwd()
	t.Cleanup(func() { os.Chdir(orig) })
	os.Chdir(tmpDir)

	// Make PATH empty so "coverage" command is not found
	origPath := os.Getenv("PATH")
	os.Setenv("PATH", "")
	t.Cleanup(func() { os.Setenv("PATH", origPath) })

	a := &PythonAdapter{cfg: &config.ServerConfig{}}
	_, err := a.Collect("handler.py", 1, 10)
	if err == nil {
		t.Fatal("expected error from coverage json command")
	}
}

func TestPythonCollect_Success(t *testing.T) {
	tmpDir := t.TempDir()
	orig, _ := os.Getwd()
	t.Cleanup(func() { os.Chdir(orig) })
	os.Chdir(tmpDir)

	// Create a fake "coverage" script that writes valid JSON to .huma/cov.json
	binDir := filepath.Join(tmpDir, "bin")
	os.MkdirAll(binDir, 0o755)

	handlerFile := filepath.Join(tmpDir, "handler.py")
	os.WriteFile(handlerFile, []byte("def handler():\n    x = 1\n    return x\n"), 0o644)

	// Build coverage.py JSON format
	covJSON := map[string]interface{}{
		"files": map[string]interface{}{
			handlerFile: map[string]interface{}{
				"executed_lines":  []int{1, 2},
				"missing_lines":   []int{3},
				"excluded_lines":  []int{},
				"summary": map[string]interface{}{
					"covered_lines":   2,
					"num_statements":  3,
					"missing_lines":   1,
					"percent_covered": 66.67,
				},
			},
		},
	}
	covData, _ := json.Marshal(covJSON)

	// Create a fake "coverage" command that writes the JSON
	script := `#!/bin/bash
cat > "` + filepath.Join(tmpDir, ".huma", "cov.json") + `" << 'JSONEOF'
` + string(covData) + `
JSONEOF
`
	os.WriteFile(filepath.Join(binDir, "coverage"), []byte(script), 0o755)

	origPath := os.Getenv("PATH")
	os.Setenv("PATH", binDir+":"+origPath)
	t.Cleanup(func() { os.Setenv("PATH", origPath) })

	a := &PythonAdapter{cfg: &config.ServerConfig{}}
	result, err := a.Collect(handlerFile, 1, 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result == nil {
		t.Fatal("expected non-nil result")
	}
	if result.Total != 3 {
		t.Fatalf("expected 3 total, got %d", result.Total)
	}
	if len(result.Uncovered) != 1 {
		t.Fatalf("expected 1 uncovered, got %d", len(result.Uncovered))
	}
}

func TestPythonCollect_ParseError(t *testing.T) {
	tmpDir := t.TempDir()
	orig, _ := os.Getwd()
	t.Cleanup(func() { os.Chdir(orig) })
	os.Chdir(tmpDir)

	// Create a fake "coverage" that writes invalid JSON
	binDir := filepath.Join(tmpDir, "bin")
	os.MkdirAll(binDir, 0o755)

	script := `#!/bin/bash
mkdir -p "` + filepath.Join(tmpDir, ".huma") + `"
echo "INVALID JSON" > "` + filepath.Join(tmpDir, ".huma", "cov.json") + `"
`
	os.WriteFile(filepath.Join(binDir, "coverage"), []byte(script), 0o755)

	origPath := os.Getenv("PATH")
	os.Setenv("PATH", binDir+":"+origPath)
	t.Cleanup(func() { os.Setenv("PATH", origPath) })

	a := &PythonAdapter{cfg: &config.ServerConfig{}}
	_, err := a.Collect("handler.py", 1, 10)
	if err == nil {
		t.Fatal("expected parse error")
	}
}

func TestPythonCollect_NoMissingLines(t *testing.T) {
	tmpDir := t.TempDir()
	orig, _ := os.Getwd()
	t.Cleanup(func() { os.Chdir(orig) })
	os.Chdir(tmpDir)

	binDir := filepath.Join(tmpDir, "bin")
	os.MkdirAll(binDir, 0o755)

	handlerFile := filepath.Join(tmpDir, "handler.py")
	os.WriteFile(handlerFile, []byte("def handler():\n    return 1\n"), 0o644)

	// No missing lines
	covJSON := map[string]interface{}{
		"files": map[string]interface{}{
			handlerFile: map[string]interface{}{
				"executed_lines":  []int{1, 2},
				"missing_lines":   []int{},
				"excluded_lines":  []int{},
				"summary": map[string]interface{}{
					"covered_lines":   2,
					"num_statements":  2,
					"missing_lines":   0,
					"percent_covered": 100.0,
				},
			},
		},
	}
	covData, _ := json.Marshal(covJSON)

	script := `#!/bin/bash
cat > "` + filepath.Join(tmpDir, ".huma", "cov.json") + `" << 'JSONEOF'
` + string(covData) + `
JSONEOF
`
	os.WriteFile(filepath.Join(binDir, "coverage"), []byte(script), 0o755)

	origPath := os.Getenv("PATH")
	os.Setenv("PATH", binDir+":"+origPath)
	t.Cleanup(func() { os.Setenv("PATH", origPath) })

	a := &PythonAdapter{cfg: &config.ServerConfig{}}
	result, err := a.Collect(handlerFile, 1, 2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Percent != 100 {
		t.Fatalf("expected 100%%, got %f", result.Percent)
	}
}
