//ff:func feature=verify type=engine control=sequence
//ff:what Computes the staged O/D/E evidence ceiling (0..3) the runtime signal can support; the A (assertion) axis is folded in verdictFromCRI, not here.

package humaquest

import (
	"github.com/park-jun-woo/huma/internal/adapter"
	"github.com/park-jun-woo/huma/internal/analyzer"
	"github.com/park-jun-woo/huma/internal/scanner"
)

// computeCRI returns the staged O/D/E evidence ceiling (0..3), faithful to the
// bak pipeline (static→smoke→covered). It is three of the four CRI axes (§1):
// O oracle, D denominator, E execution. The fourth axis — A (assertion depth) —
// is folded on top of this ceiling in verdictFromCRI (which already receives
// aGrade), so the full CRI = min(O,D,A,E) is realized there, not here:
//
//	0 UNVERIFIED  no oracle OR no denominator — an axis is 0 (CV-1/2/3).
//	1 SCAFFOLDED  source-linked, no execution evidence (covAbsent, E low).
//	2 SMOKE       executed green, no runtime line binding (cov.Total==0, E=2).
//	3 COVERED     100% line coverage AND strict branch-line binding (E=3, §3.4).
//
// An instrumented run that is below 100% or fails binding has not reached
// COVERED; it falls back to SMOKE (2) as the honest measured tier — the caller's
// gate/IMPROVE logic decides PASS vs retry from there.
func computeCRI(ep *scanner.Endpoint, branches []analyzer.ResponseBranch, cov *adapter.CoverageResult, covPresent bool) int {
	if !hasOracle(ep, cov) || len(branches) == 0 {
		return 0
	}
	if !covPresent {
		return 1
	}
	if cov == nil || cov.Total == 0 {
		return 2
	}
	if cov.Percent == 100 && branchBindingOK(ep, branches, cov) {
		return 3
	}
	return 2
}
