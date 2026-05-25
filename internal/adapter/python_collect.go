//ff:func feature=adapter type=engine control=sequence
//ff:what Runs coverage json export and parses the result for handler coverage
package adapter

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/park-jun-woo/huma/internal/coverage"
)

// Collect runs `coverage json` to export coverage data, then parses
// the JSON output to extract uncovered lines for the given handler.
func (a *PythonAdapter) Collect(handlerFile string, startLine, endLine int) (*CoverageResult, error) {
	if err := os.MkdirAll(filepath.Dir(coverageJSON), 0o755); err != nil {
		return nil, fmt.Errorf("create output dir: %w", err)
	}

	out, err := exec.Command("coverage", "json", "-o", coverageJSON).CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("coverage json: %w\n%s", err, string(out))
	}

	missingLines, err := coverage.ParseCoveragePy(coverageJSON, handlerFile, startLine, endLine)
	if err != nil {
		return nil, fmt.Errorf("parse coverage: %w", err)
	}

	totalCount := endLine - startLine + 1
	uncoveredCount := len(missingLines)
	coveredCount := totalCount - uncoveredCount
	if coveredCount < 0 {
		coveredCount = 0
	}

	var pct float64
	if totalCount > 0 {
		pct = float64(coveredCount) / float64(totalCount) * 100
	}

	uncovered := make([]UncoveredLine, len(missingLines))
	for i, line := range missingLines {
		uncovered[i] = UncoveredLine{
			File: handlerFile,
			Line: line,
		}
	}

	return &CoverageResult{
		Covered:   coveredCount,
		Total:     totalCount,
		Percent:   pct,
		Uncovered: uncovered,
	}, nil
}
