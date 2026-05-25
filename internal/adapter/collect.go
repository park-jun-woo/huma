//ff:func feature=adapter type=engine control=sequence
//ff:what Converts raw coverage data and analyzes handler function coverage
package adapter

import (
	"fmt"
	"os/exec"

	"github.com/park-jun-woo/huma/internal/coverage"
)

// Collect converts raw coverage data and analyzes coverage for the handler.
func (a *GoAdapter) Collect(handlerFile string, startLine, endLine int) (*CoverageResult, error) {
	out, err := exec.Command("go", "tool", "covdata", "textfmt",
		"-i="+a.coverDir,
		"-o="+coverOut,
	).CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("covdata textfmt: %w\n%s", err, string(out))
	}

	blocks, err := coverage.ParseCoverageFile(coverOut)
	if err != nil {
		return nil, fmt.Errorf("parse coverage: %w", err)
	}

	filtered := coverage.FilterBlocks(blocks, handlerFile, startLine, endLine)
	return buildResult(filtered, handlerFile, startLine, endLine)
}
