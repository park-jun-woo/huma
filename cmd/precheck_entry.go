//ff:func feature=scan type=engine control=sequence
//ff:what Marks one entry's initial verdict from its existing hurl file (UNVERIFIED/IMPROVE/SCAFFOLDED)
package cmd

import (
	"github.com/park-jun-woo/huma/internal/config"
	"github.com/park-jun-woo/huma/internal/runner"
	"github.com/park-jun-woo/huma/internal/session"
)

// precheckEntry assigns one entry's initial verdict based on its existing hurl
// file. No-signal endpoints stay UNVERIFIED; static pre-pass tops out at
// SCAFFOLDED (CRI 1).
func precheckEntry(sess *session.Session, e *session.Entry, cfg *config.Config) {
	ep := &e.Endpoint
	hurl := runner.FindHurlFile(ep, cfg.HurlDir)
	if hurl == "" {
		return
	}
	branches, prov := responseBranches(ep, cfg.Scan.Lang)
	if !hasOracle(ep, nil) || len(branches) == 0 {
		sess.MarkUnverified(e.ID)
		return
	}
	result := checkResponseCoverageFn(ep, hurl, cfg.Scan.Lang)
	if result != nil && result.Total > 0 && len(result.Missing) > 0 {
		sess.MarkImprove(e.ID, result.Percent)
		return
	}
	if resolveGate(cfg, 1) > 1 {
		sess.MarkUnverified(e.ID)
		return
	}
	sess.SetVerdict(e.ID, 1, staticAGrade(hurl, branches), prov.String())
	sess.MarkPass(e.ID)
}
