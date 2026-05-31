//ff:func feature=ratchet type=helper control=sequence
//ff:what Builds the response-branch denominator as the monotonic union of source analysis and OpenAPI declarations
package cmd

import (
	"github.com/park-jun-woo/huma/internal/analyzer"
	"github.com/park-jun-woo/huma/internal/scanner"
)

// responseBranches builds the response-branch denominator. The denominator is
// the union of ground-truth source branches (the authoritative floor) and
// OpenAPI declarations (additive only). Input can never shrink the
// denominator (monotonic increase, §3.1 / C-02).
func responseBranches(ep *scanner.Endpoint, lang string) ([]analyzer.ResponseBranch, BranchProvenance) {
	src := sourceBranches(ep, lang)
	decl := declaredBranches(ep)
	union := dedupByStatus(concatBranches(src, decl))
	client, _ := splitClientBranches(union)
	return client, provenanceOf(src, decl)
}
