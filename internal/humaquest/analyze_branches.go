//ff:func feature=ratchet type=helper control=sequence
//ff:what Runs the static analyzer on the endpoint handler source to extract its response branches (ground-truth floor).

package humaquest

import (
	"github.com/park-jun-woo/huma/internal/analyzer"
	"github.com/park-jun-woo/huma/internal/scanner"
)

// analyzeBranches runs the static analyzer on the endpoint source file. Returns
// nil when there is no source, no analyzer for the language, or analysis fails.
// Ported from bak/cmd/analyze_branches.go.
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
