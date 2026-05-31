//ff:func feature=ratchet type=helper control=sequence
//ff:what Classifies whether the denominator came from source analysis, declarations, or both
package cmd

import "github.com/park-jun-woo/huma/internal/analyzer"

// provenanceOf classifies where the denominator branches originated.
func provenanceOf(src, decl []analyzer.ResponseBranch) BranchProvenance {
	return BranchProvenance{
		HasSource:   len(src) > 0,
		HasDeclared: len(decl) > 0,
	}
}
