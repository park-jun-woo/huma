//ff:func feature=verify type=builder control=selection
//ff:what Prints the IMPROVE or UNVERIFIED prompt for a non-passing live verdict; reports whether handled
package cmd

import (
	"fmt"

	"github.com/park-jun-woo/huma/internal/adapter"
	"github.com/park-jun-woo/huma/internal/config"
	"github.com/park-jun-woo/huma/internal/prompt"
	"github.com/park-jun-woo/huma/internal/scanner"
)

// printLiveNonPass prints the prompt for a non-passing live outcome and returns
// true when it handled the outcome (caller should stop).
func printLiveNonPass(oc outcome, ep *scanner.Endpoint, hurl string, cov *adapter.CoverageResult, cfg *config.Config) bool {
	switch oc {
	case outcomeImprove:
		fmt.Print(prompt.ImprovePrompt(ep, hurl, cov))
		return true
	case outcomeUnverified:
		fmt.Print(prompt.UnverifiedPrompt(ep, cfg))
		return true
	default:
		return false
	}
}
