//ff:func feature=ratchet type=helper control=sequence
//ff:what Returns the ground-truth source response branches (the authoritative denominator floor) when the endpoint is source-linked.

package humaquest

import (
	"github.com/park-jun-woo/huma/internal/analyzer"
	"github.com/park-jun-woo/huma/internal/scanner"
)

// sourceBranches returns the authoritative source-derived branches (the
// denominator floor). Empty when the endpoint has no source link. Ported from
// bak/cmd/source_branches.go.
func sourceBranches(ep *scanner.Endpoint, lang string) []analyzer.ResponseBranch {
	if ep.Source == "" {
		return nil
	}
	return analyzeBranches(ep, lang)
}
