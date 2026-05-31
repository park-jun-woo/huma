//ff:func feature=ratchet type=helper control=sequence
//ff:what Returns the ground-truth source response branches when the endpoint has a linked source
package cmd

import (
	"github.com/park-jun-woo/huma/internal/analyzer"
	"github.com/park-jun-woo/huma/internal/scanner"
)

// sourceBranches returns the authoritative source-derived branches (the
// denominator floor). Empty when the endpoint has no source link.
func sourceBranches(ep *scanner.Endpoint, lang string) []analyzer.ResponseBranch {
	if ep.Source == "" {
		return nil
	}
	return analyzeBranches(ep, lang)
}
