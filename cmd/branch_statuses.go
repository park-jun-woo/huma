//ff:func feature=verify type=helper control=iteration dimension=1
//ff:what Extracts the status codes from a slice of response branches
package cmd

import "github.com/park-jun-woo/huma/internal/analyzer"

// branchStatuses returns the status codes of the supplied branches.
func branchStatuses(branches []analyzer.ResponseBranch) []int {
	out := make([]int, 0, len(branches))
	for _, b := range branches {
		out = append(out, b.Status)
	}
	return out
}
