//ff:func feature=verify type=engine control=sequence
//ff:what Applies the SMOKE (CRI 2) verdict for a source-linked but uninstrumented server that ran green
package cmd

import (
	"github.com/park-jun-woo/huma/internal/config"
	"github.com/park-jun-woo/huma/internal/scanner"
	"github.com/park-jun-woo/huma/internal/session"
)

// smokeVerdict applies the SMOKE tier: server executed green with no runtime
// line coverage. CRI 2 is the reachable maximum here.
func smokeVerdict(cfg *config.Config, sess *session.Session, ep *scanner.Endpoint, aGrade int, prov string) outcome {
	if 2 < resolveGate(cfg, 2) {
		sess.MarkUnverified(ep.ID)
		return outcomeUnverified
	}
	sess.SetVerdict(ep.ID, 2, aGrade, prov)
	sess.MarkPass(ep.ID)
	return outcomePass
}
