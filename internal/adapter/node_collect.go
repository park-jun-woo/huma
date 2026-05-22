//ff:func feature=adapter type=engine control=sequence
//ff:what Runs c8 report to convert V8 coverage to istanbul JSON and analyzes handler coverage
package adapter

import (
	"fmt"
	"os/exec"

	"github.com/park-jun-woo/hurlfill/internal/coverage"
)

// Collect runs c8 report to convert V8 coverage to istanbul JSON,
// then parses the result and analyzes coverage for the handler.
func (a *NodeAdapter) Collect(handlerFile string, startLine, endLine int) (*CoverageResult, error) {
	out, err := exec.Command("npx", "c8", "report",
		"--reporter=json",
		"--temp-directory="+a.coverDir,
		"-o", istanbulOutDir,
	).CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("c8 report: %w\n%s", err, string(out))
	}

	blocks, err := coverage.ParseIstanbul(istanbulOutFile)
	if err != nil {
		return nil, fmt.Errorf("parse istanbul: %w", err)
	}

	filtered := coverage.FilterBlocks(blocks, handlerFile, startLine, endLine)
	return buildResult(filtered, handlerFile, startLine, endLine)
}
