//ff:func feature=ratchet type=helper control=iteration dimension=1
//ff:what Splits response branches into client (<500, gated) and advisory server (>=500, shown but not gating)
package cmd

import "github.com/park-jun-woo/huma/internal/analyzer"

// splitClientBranches partitions branches into client branches (Status < 500,
// which form the PASS-gate denominator) and advisory server branches
// (Status >= 500, kept for visibility but never silently dropped). Dropping
// 5xx would shrink the denominator and become a cheese enabler (§3.6).
func splitClientBranches(branches []analyzer.ResponseBranch) (client, advisory []analyzer.ResponseBranch) {
	for _, b := range branches {
		if b.Status >= 500 {
			advisory = append(advisory, b)
		} else {
			client = append(client, b)
		}
	}
	return client, advisory
}
