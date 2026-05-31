//ff:func feature=verify type=engine control=iteration dimension=1
//ff:what Checks strict branch-line binding for all client branches with unreachable.yaml exemptions
package cmd

import (
	"github.com/park-jun-woo/huma/internal/adapter"
	"github.com/park-jun-woo/huma/internal/analyzer"
	"github.com/park-jun-woo/huma/internal/config"
	"github.com/park-jun-woo/huma/internal/scanner"
)

// branchBindingOK reports whether all client branches are runtime-bound. With
// no branches there is nothing to bind (caller relies on line-percent). The
// Go-only traceHelperCall exemption is already reflected in branch.Line; the
// unreachable.yaml reason exempts branches that cannot be bound (§3.4).
func branchBindingOK(ep *scanner.Endpoint, branches []analyzer.ResponseBranch, cov *adapter.CoverageResult) bool {
	if len(branches) == 0 {
		return true
	}
	exemptions, _ := config.LoadUnreachable()
	key := ep.Method + " " + ep.Path
	cohesion := lineCohesion(branches)
	for _, b := range branches {
		if !branchBound(b, cov, exemptions, key, cohesion) {
			return false
		}
	}
	return true
}
