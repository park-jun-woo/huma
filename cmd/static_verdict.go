//ff:func feature=verify type=engine control=sequence
//ff:what Computes and applies the static-mode verdict (no-signal→UNVERIFIED, source∪decl gate, SCAFFOLDED ceiling)
package cmd

import (
	"github.com/park-jun-woo/huma/internal/config"
	"github.com/park-jun-woo/huma/internal/hurlcheck"
	"github.com/park-jun-woo/huma/internal/scanner"
	"github.com/park-jun-woo/huma/internal/session"
)

// staticVerdict computes and applies the verdict for static mode (no server).
// It returns the chosen outcome and the static response-coverage result.
func staticVerdict(cfg *config.Config, sess *session.Session, ep *scanner.Endpoint, hurl string) (outcome, *hurlcheck.ResponseCoverageResult) {
	branches, prov := responseBranches(ep, cfg.Scan.Lang)

	// No oracle (source unlinked) or no branch denominator ⇒ UNVERIFIED (C-01).
	if !hasOracle(ep, nil) || len(branches) == 0 {
		sess.MarkUnverified(ep.ID)
		return outcomeUnverified, nil
	}

	respResult := checkResponseCoverageFn(ep, hurl, cfg.Scan.Lang)
	if respResult != nil && respResult.Total > 0 && len(respResult.Missing) > 0 {
		sess.MarkImprove(ep.ID, respResult.Percent)
		return outcomeImprove, respResult
	}

	// All client branches covered. Static honest ceiling is SCAFFOLDED (CRI 1);
	// a higher require_cri needs live execution, so fall to UNVERIFIED.
	if 1 < resolveGate(cfg, 1) {
		sess.MarkUnverified(ep.ID)
		return outcomeUnverified, respResult
	}
	sess.SetVerdict(ep.ID, 1, staticAGrade(hurl, branches), prov.String())
	sess.MarkPass(ep.ID)
	return outcomePass, respResult
}
