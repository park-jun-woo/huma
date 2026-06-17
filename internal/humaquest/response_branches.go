//ff:func feature=ratchet type=helper control=sequence
//ff:what Builds the response-branch denominator as the monotonic union of source analysis and OpenAPI declarations, returning the gated client set plus provenance.

package humaquest

import (
	"github.com/park-jun-woo/huma/internal/analyzer"
	"github.com/park-jun-woo/huma/internal/scanner"
)

// responseBranches builds the response-branch denominator: the union of
// ground-truth source branches (the authoritative floor) and OpenAPI
// declarations (additive only). Input can never shrink the denominator
// (monotonic increase, §3.1 / C-02). It returns the gated client branches and
// the provenance of the union. Ported from bak/cmd/response_branches.go.
func responseBranches(ep *scanner.Endpoint, lang string) ([]analyzer.ResponseBranch, BranchProvenance) {
	src := sourceBranches(ep, lang)
	decl := declaredBranches(ep)
	union := dedupByStatus(concatBranches(src, decl))
	client, _ := splitClientBranches(union)
	return client, provenanceOf(src, decl)
}
