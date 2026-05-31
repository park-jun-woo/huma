//ff:func feature=verify type=engine control=sequence
//ff:what Computes and applies the live-mode verdict (oracle gate, SMOKE/COVERED tiers, branch-line binding)
package cmd

import (
	"github.com/park-jun-woo/huma/internal/adapter"
	"github.com/park-jun-woo/huma/internal/config"
	"github.com/park-jun-woo/huma/internal/scanner"
	"github.com/park-jun-woo/huma/internal/session"
)

// liveVerdict computes and applies the verdict for live mode after a passing
// hurl run. covResult may be nil (uninstrumented server).
func liveVerdict(cfg *config.Config, sess *session.Session, ep *scanner.Endpoint, hurl string, covResult *adapter.CoverageResult, entry *session.Entry) outcome {
	branches, prov := responseBranches(ep, cfg.Scan.Lang)

	// No oracle: source unlinked AND no runtime instrumentation ⇒ UNVERIFIED.
	if !hasOracle(ep, covResult) {
		sess.MarkUnverified(ep.ID)
		return outcomeUnverified
	}

	aGrade := staticAGrade(hurl, branches)
	if covResult == nil || covResult.Total == 0 {
		return smokeVerdict(cfg, sess, ep, aGrade, prov.String())
	}
	return coveredVerdict(cfg, sess, ep, branches, covResult, entry, aGrade, prov.String())
}
