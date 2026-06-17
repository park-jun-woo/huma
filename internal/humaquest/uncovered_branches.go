//ff:func feature=verify type=helper control=iteration dimension=1
//ff:what Returns the client branches whose source line was not runtime-covered (they need a reachability reason).

package humaquest

import (
	"github.com/park-jun-woo/huma/internal/adapter"
	"github.com/park-jun-woo/huma/internal/analyzer"
)

// uncoveredBranches returns branches whose source line was not hit at runtime.
// Ported from bak/cmd/uncovered_branches.go.
func uncoveredBranches(branches []analyzer.ResponseBranch, cov *adapter.CoverageResult) []analyzer.ResponseBranch {
	var out []analyzer.ResponseBranch
	for _, b := range branches {
		if cov != nil && b.Line != 0 && cov.IsLineCovered(b.Line) {
			continue
		}
		out = append(out, b)
	}
	return out
}
