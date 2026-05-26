//ff:func feature=ratchet type=helper control=sequence
//ff:what Runs static analyzer on endpoint source to extract response branches
package cmd

import (
	"github.com/park-jun-woo/huma/internal/analyzer"
	"github.com/park-jun-woo/huma/internal/scanner"
)

// analyzeBranches runs the static analyzer on the endpoint source file.
func analyzeBranches(ep *scanner.Endpoint, lang string) []analyzer.ResponseBranch {
	if ep.Source == "" {
		return nil
	}
	a := analyzer.NewAnalyzer(lang)
	if a == nil {
		return nil
	}
	branches, err := a.Analyze(ep.Source, ep.Handler, ep.Line, 0)
	if err != nil {
		return nil
	}
	return branches
}
