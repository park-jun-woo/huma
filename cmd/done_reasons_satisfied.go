//ff:func feature=verify type=engine control=sequence
//ff:what Reports whether all uncovered client branches have valid unreachable.yaml reasons for DONE
package cmd

import (
	"github.com/park-jun-woo/huma/internal/adapter"
	"github.com/park-jun-woo/huma/internal/analyzer"
	"github.com/park-jun-woo/huma/internal/config"
	"github.com/park-jun-woo/huma/internal/scanner"
)

// doneReasonsSatisfied reports whether every uncovered client branch has a
// valid unreachable.yaml reason — required to grant DONE (§3.7, C-04).
func doneReasonsSatisfied(ep *scanner.Endpoint, branches []analyzer.ResponseBranch, cov *adapter.CoverageResult) bool {
	uncovered := uncoveredBranches(branches, cov)
	if len(uncovered) == 0 {
		return true // nothing uncovered to justify
	}
	exemptions, err := config.LoadUnreachable()
	if err != nil || len(exemptions) == 0 {
		return false
	}
	return allExempt(uncovered, exemptions, ep.Method+" "+ep.Path)
}
