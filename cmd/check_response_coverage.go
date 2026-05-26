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
	branches := responseBranches(ep, lang)
	if len(branches) == 0 {
		return nil
	}

	statuses, err := hurlcheck.ParseHurlStatuses(hurlFile)
	if err != nil {
		return nil
	}

	return hurlcheck.CheckResponseCoverage(branches, statuses)
}
