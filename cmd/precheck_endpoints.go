//ff:func feature=scan type=engine control=iteration dimension=1
//ff:what Checks existing hurl files for all endpoints and marks session entries as PASS or IMPROVE
package cmd

import (
	"github.com/park-jun-woo/huma/internal/config"
	"github.com/park-jun-woo/huma/internal/runner"
	"github.com/park-jun-woo/huma/internal/session"
)

func precheckEndpoints(sess *session.Session, cfg *config.Config) {
	if cfg == nil {
		return
	}
	for i := range sess.Entries {
		e := &sess.Entries[i]
		ep := &e.Endpoint

		hurl := runner.FindHurlFile(ep, cfg.HurlDir)
		if hurl == "" {
			continue
		}

		result := checkResponseCoverageFn(ep, hurl, cfg.Scan.Lang)
		if result != nil && result.Total > 0 && len(result.Missing) > 0 {
			sess.MarkImprove(e.ID, result.Percent)
			continue
		}

		sess.MarkPass(e.ID)
	}
}
