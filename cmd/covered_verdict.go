//ff:func feature=verify type=engine control=sequence
//ff:what Applies the COVERED/DONE/IMPROVE verdict for an instrumented live run using strict branch-line binding
package cmd

import (
	"github.com/park-jun-woo/huma/internal/adapter"
	"github.com/park-jun-woo/huma/internal/analyzer"
	"github.com/park-jun-woo/huma/internal/config"
	"github.com/park-jun-woo/huma/internal/scanner"
	"github.com/park-jun-woo/huma/internal/session"
)

// coveredVerdict applies the instrumented-mode verdict. COVERED (CRI 3)
// requires 100% line coverage AND strict branch-line binding (§3.4); a stalled
// sub-100% run becomes DONE only with unreachable.yaml reasons (§3.7), else
// UNVERIFIED.
func coveredVerdict(cfg *config.Config, sess *session.Session, ep *scanner.Endpoint, branches []analyzer.ResponseBranch, cov *adapter.CoverageResult, entry *session.Entry, aGrade int, prov string) outcome {
	if cov.Percent == 100 && branchBindingOK(ep, branches, cov) {
		if 3 < resolveGate(cfg, 3) {
			sess.MarkUnverified(ep.ID)
			return outcomeUnverified
		}
		sess.SetVerdict(ep.ID, 3, aGrade, prov)
		sess.MarkPass(ep.ID)
		return outcomePass
	}
	if stalled(entry, cov) {
		return stalledVerdict(sess, ep, branches, cov, aGrade, prov)
	}
	sess.MarkImprove(ep.ID, cov.Percent)
	return outcomeImprove
}
