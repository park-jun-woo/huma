//ff:func feature=ratchet type=helper control=sequence
//ff:what Concatenates source and declared branches source-first so the authoritative source Line wins on a later status dedup.

package humaquest

import "github.com/park-jun-woo/huma/internal/analyzer"

// concatBranches appends source branches before declared ones so the
// source-derived Line (authoritative) wins on a later status dedup. Ported from
// bak/cmd/concat_branches.go.
func concatBranches(src, decl []analyzer.ResponseBranch) []analyzer.ResponseBranch {
	out := make([]analyzer.ResponseBranch, 0, len(src)+len(decl))
	out = append(out, src...)
	out = append(out, decl...)
	return out
}
