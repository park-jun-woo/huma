//ff:func feature=verify type=helper control=iteration dimension=1
//ff:what Counts how many branches share each source line — a line shared by more than one branch cannot be resolved to a single branch.

package humaquest

import "github.com/park-jun-woo/huma/internal/analyzer"

// lineCohesion returns a map of source line → number of branches at that line.
// A line shared by more than one branch is cohesive and cannot be bound (§3.4).
// Ported from bak/cmd/line_cohesion.go.
func lineCohesion(branches []analyzer.ResponseBranch) map[int]int {
	out := make(map[int]int, len(branches))
	for _, b := range branches {
		out[b.Line]++
	}
	return out
}
