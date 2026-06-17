//ff:func feature=ratchet type=helper control=iteration dimension=1
//ff:what Deduplicates response branches by status code, keeping the first (source-preferred) occurrence.

package humaquest

import "github.com/park-jun-woo/huma/internal/analyzer"

// dedupByStatus removes duplicate branches by status code, keeping the first
// occurrence (source-preferred via concatBranches order). Ported from
// bak/cmd/dedup_by_status.go.
func dedupByStatus(branches []analyzer.ResponseBranch) []analyzer.ResponseBranch {
	seen := make(map[int]bool, len(branches))
	out := make([]analyzer.ResponseBranch, 0, len(branches))
	for _, b := range branches {
		if seen[b.Status] {
			continue
		}
		seen[b.Status] = true
		out = append(out, b)
	}
	return out
}
