//ff:func feature=verify type=builder control=selection
//ff:what Prints the IMPROVE or UNVERIFIED prompt for a non-passing static verdict; reports whether handled
package cmd

import (
	"fmt"

	"github.com/park-jun-woo/huma/internal/config"
	"github.com/park-jun-woo/huma/internal/hurlcheck"
	"github.com/park-jun-woo/huma/internal/prompt"
	"github.com/park-jun-woo/huma/internal/scanner"
)

// printStaticNonPass prints the prompt for a non-passing static outcome and
// returns true when it handled the outcome (caller should stop).
func printStaticNonPass(oc outcome, ep *scanner.Endpoint, hurl string, res *hurlcheck.ResponseCoverageResult, cfg *config.Config) bool {
	switch oc {
	case outcomeImprove:
		fmt.Print(prompt.ResponseImprovePrompt(ep, hurl, res))
		return true
	case outcomeUnverified:
		fmt.Print(prompt.UnverifiedPrompt(ep, cfg))
		return true
	default:
		return false
	}
}
