//ff:func feature=verify type=engine control=sequence
//ff:what Grants DONE for a stalled run only if uncovered branches have unreachable.yaml reasons, else UNVERIFIED
package cmd

import (
	"github.com/park-jun-woo/huma/internal/adapter"
	"github.com/park-jun-woo/huma/internal/analyzer"
	"github.com/park-jun-woo/huma/internal/scanner"
	"github.com/park-jun-woo/huma/internal/session"
)

// stalledVerdict resolves a stalled instrumented run: DONE when every uncovered
// branch has a verifiable unreachable.yaml reason (§3.7, C-04), otherwise the
// honest UNVERIFIED.
func stalledVerdict(sess *session.Session, ep *scanner.Endpoint, branches []analyzer.ResponseBranch, cov *adapter.CoverageResult, aGrade int, prov string) outcome {
	if !doneReasonsSatisfied(ep, branches, cov) {
		sess.MarkUnverified(ep.ID)
		return outcomeUnverified
	}
	sess.SetVerdict(ep.ID, 3, aGrade, prov)
	sess.MarkDone(ep.ID, cov.Percent)
	return outcomeDone
}
