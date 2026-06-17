//ff:func feature=verify type=engine control=sequence
//ff:what Reports whether all runtime-uncovered client branches have valid unreachable.yaml reasons, gating DONE for an instrumented run.

package humaquest

import (
	"github.com/park-jun-woo/huma/internal/adapter"
	"github.com/park-jun-woo/huma/internal/analyzer"
	"github.com/park-jun-woo/huma/internal/scanner"
)

// doneReasonsSatisfied reports whether every runtime-uncovered client branch has
// a valid unreachable.yaml reason — required to grant DONE (§3.7, C-04). Ported
// from bak/cmd/done_reasons_satisfied.go, delegating the exemption check to
// reasonsCover over the line-uncovered set.
func doneReasonsSatisfied(ep *scanner.Endpoint, branches []analyzer.ResponseBranch, cov *adapter.CoverageResult) bool {
	return reasonsCover(ep, uncoveredBranches(branches, cov))
}
