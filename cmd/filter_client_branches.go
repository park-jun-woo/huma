//ff:func feature=ratchet type=helper control=sequence
//ff:what Returns only the client (gated) branches, wrapping splitClientBranches for callers that need the gate set
package cmd

import "github.com/park-jun-woo/huma/internal/analyzer"

// filterClientBranches returns only the client (gated) branches.
func filterClientBranches(branches []analyzer.ResponseBranch) []analyzer.ResponseBranch {
	if len(branches) == 0 {
		return branches
	}
	client, _ := splitClientBranches(branches)
	return client
}
