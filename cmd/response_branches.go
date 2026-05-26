//ff:func feature=ratchet type=helper control=sequence
//ff:what Extracts response branches from endpoint responses field or source file analysis
package cmd

import (
	"github.com/park-jun-woo/huma/internal/analyzer"
	"github.com/park-jun-woo/huma/internal/scanner"
)

// responseBranches extracts response branches from the endpoint.
func responseBranches(ep *scanner.Endpoint) []analyzer.ResponseBranch {
	if len(ep.Responses) > 0 {
		branches := analyzer.ParseResponses(ep.Responses, ep.Source)
		if len(branches) > 0 {
			return branches
		}
	}

	if ep.Source == "" {
		return nil
	}

	a := &analyzer.GoAnalyzer{}
	branches, err := a.Analyze(ep.Source, ep.Handler, ep.Line, 0)
	if err != nil {
		return nil
	}

	return branches
}
