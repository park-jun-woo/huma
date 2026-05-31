//ff:func feature=ratchet type=engine control=sequence
//ff:what Performs static response analysis on an endpoint by comparing analyzer branches with hurl statuses
package cmd

import (
	"github.com/park-jun-woo/huma/internal/hurlcheck"
	"github.com/park-jun-woo/huma/internal/scanner"
)

// checkResponseCoverageFn allows tests to replace checkResponseCoverage.
var checkResponseCoverageFn = checkResponseCoverage

// checkResponseCoverage performs static response analysis on the endpoint.
func checkResponseCoverage(ep *scanner.Endpoint, hurlFile string, lang string) *hurlcheck.ResponseCoverageResult {
	branches, _ := responseBranches(ep, lang)
	if len(branches) == 0 {
		return nil
	}

	entries, err := hurlcheck.ParseHurlEntries(hurlFile)
	if err != nil {
		return nil
	}

	// Only non-vacuous entries (real status assertion, not skipped) count
	// toward coverage — blocks vacuous green (§3.5, CV-10).
	statuses := hurlcheck.NonVacuousStatusList(entries)
	return hurlcheck.CheckResponseCoverage(branches, statuses)
}
