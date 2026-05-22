//ff:func feature=adapter type=engine control=sequence
//ff:what Builds a CoverageResult from filtered coverage blocks and handler line range
package adapter

import (
	"github.com/park-jun-woo/hurlfill/internal/coverage"
)

func buildResult(filtered []coverage.Block, handlerFile string, startLine, endLine int) (*CoverageResult, error) {
	covered, total := computeLineCoverage(filtered, startLine, endLine)

	totalCount := len(total)
	coveredCount := len(covered)
	var pct float64
	if totalCount > 0 {
		pct = float64(coveredCount) / float64(totalCount) * 100
	}

	uncoveredLines, err := readUncoveredLines(handlerFile, covered, total)
	if err != nil {
		uncoveredLines = nil
	}

	return &CoverageResult{
		Covered:   coveredCount,
		Total:     totalCount,
		Percent:   pct,
		Uncovered: uncoveredLines,
	}, nil
}
