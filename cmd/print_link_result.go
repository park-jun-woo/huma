//ff:func feature=scan type=command control=iteration dimension=1
//ff:what Prints the link-source result: per-language link distribution and skip reasons
package cmd

import (
	"fmt"

	"github.com/park-jun-woo/huma/internal/scanner"
)

// printLinkResult renders the --link-source summary: linked count with a
// per-extension distribution, skipped count broken into ext-mismatch/ambiguous,
// and the user-facing reason for each skipped endpoint (§2.5/§2.6).
func printLinkResult(r scanner.LinkResult, total int, root string) {
	fmt.Printf("Linked %d/%d endpoints under %s  (%s)\n", r.Linked, total, root, distribution(r))
	fmt.Printf("Skipped %d  (ext-mismatch: %d, ambiguous: %d)\n", r.Skipped, r.ExtMismatch, r.Ambiguous)
	for _, msg := range r.SkipMessages {
		fmt.Println("  " + msg)
	}
}
