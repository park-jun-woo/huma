//ff:func feature=ratchet type=helper control=sequence
//ff:what Concatenates source and declared branches with source first (source Line wins on dedup)
package cmd

import "github.com/park-jun-woo/huma/internal/analyzer"

// concatBranches appends source branches before declared ones so the
// source-derived Line (authoritative) wins on a later status dedup.
func concatBranches(src, decl []analyzer.ResponseBranch) []analyzer.ResponseBranch {
	out := make([]analyzer.ResponseBranch, 0, len(src)+len(decl))
	out = append(out, src...)
	out = append(out, decl...)
	return out
}
