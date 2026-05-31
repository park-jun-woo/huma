//ff:func feature=ratchet type=helper control=sequence
//ff:what Returns the OpenAPI-declared response branches (additive only, never shrinks the denominator)
package cmd

import (
	"github.com/park-jun-woo/huma/internal/analyzer"
	"github.com/park-jun-woo/huma/internal/scanner"
)

// declaredBranches returns the OpenAPI-declared branches. These are additive
// only — they can extend but never reduce the source floor.
func declaredBranches(ep *scanner.Endpoint) []analyzer.ResponseBranch {
	if len(ep.Responses) == 0 {
		return nil
	}
	return analyzer.ParseResponses(ep.Responses, ep.Source)
}
