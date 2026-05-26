//ff:func feature=ratchet type=helper control=iteration dimension=1
//ff:what Removes server error branches (Status >= 500) that cannot be triggered by client input
package cmd

import "github.com/park-jun-woo/huma/internal/analyzer"

// filterClientBranches removes server error branches (Status >= 500)
// that cannot be triggered by client input.
func filterClientBranches(branches []analyzer.ResponseBranch) []analyzer.ResponseBranch {
	if len(branches) == 0 {
		return branches
	}
	n := 0
	for _, b := range branches {
		if b.Status < 500 {
			branches[n] = b
			n++
		}
	}
	return branches[:n]
}
