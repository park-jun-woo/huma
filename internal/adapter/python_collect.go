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

	executedLines, missingLines, err := coverage.ParseCoveragePy(coverageJSON, handlerFile, startLine, endLine)
	if err != nil {
		return nil, fmt.Errorf("parse coverage: %w", err)
	}

	totalCount := len(executedLines) + len(missingLines)
	coveredCount := len(executedLines)

	var pct float64
	if totalCount > 0 {
		pct = float64(coveredCount) / float64(totalCount) * 100
	}

	covered := make(map[int]bool, len(executedLines))
	for _, ln := range executedLines {
		covered[ln] = true
	}
	total := make(map[int]bool, totalCount)
	for _, ln := range executedLines {
		total[ln] = true
	}
	for _, ln := range missingLines {
		total[ln] = true
	}

	uncovered, err := readUncoveredLines(handlerFile, covered, total)
	if err != nil {
		uncovered = nil
	}

	return &CoverageResult{
		Covered:      coveredCount,
		Total:        totalCount,
		Percent:      pct,
		Uncovered:    uncovered,
		CoveredLines: covered,
	}, nil
}
