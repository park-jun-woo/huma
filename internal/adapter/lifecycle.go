//ff:func feature=adapter type=engine control=sequence
//ff:what Executes the full adapter lifecycle: reset, start, run hurl, stop, and collect coverage
package adapter

import (
	"fmt"

	"github.com/park-jun-woo/huma/internal/runner"
	"github.com/park-jun-woo/huma/internal/source"
)

// RunWithCoverage executes the full adapter lifecycle: reset, start, wait,
// run hurl test, stop, and collect coverage data.
// Returns the hurl result, coverage result (nil if hurl failed), and any error.
func RunWithCoverage(a Adapter, hurlPath string, variables map[string]string, handlerFile, handlerName string) (*runner.Result, *CoverageResult, error) {
	if err := a.Reset(); err != nil {
		return nil, nil, fmt.Errorf("reset coverage: %w", err)
	}

	if err := a.Start(); err != nil {
		return nil, nil, fmt.Errorf("start server: %w", err)
	}

	if err := a.WaitReady(); err != nil {
		a.Stop()
		return nil, nil, fmt.Errorf("wait ready: %w", err)
	}

	result, err := runner.Run(hurlPath, variables)

	// Always stop the server to trigger coverage dump
	stopErr := a.Stop()

	if err != nil {
		return nil, nil, fmt.Errorf("hurl run: %w", err)
	}
	if stopErr != nil {
		return result, nil, fmt.Errorf("stop server: %w", stopErr)
	}

	// If hurl failed, no coverage to collect
	if !result.Pass {
		return result, nil, nil
	}

	// Read handler line range for coverage filtering
	_, startLine, endLine, err := source.ReadHandler(handlerFile, handlerName)
	if err != nil {
		// Cannot determine handler bounds — skip coverage
		return result, nil, nil
	}

	covResult, err := a.Collect(handlerFile, startLine, endLine)
	if err != nil {
		return result, nil, fmt.Errorf("collect coverage: %w", err)
	}

	return result, covResult, nil
}
