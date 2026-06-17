//ff:func feature=verify type=helper control=iteration dimension=1
//ff:what Reports whether every supplied branch has a valid unreachable.yaml exemption for an endpoint key.

package humaquest

import (
	"github.com/park-jun-woo/huma/internal/analyzer"
	"github.com/park-jun-woo/huma/internal/config"
)

// allExempt reports whether every branch in the set is exempt for endpoint key.
// Ported from bak/cmd/all_exempt.go.
func allExempt(branches []analyzer.ResponseBranch, exemptions []config.UnreachableEntry, key string) bool {
	for _, b := range branches {
		if !config.IsExempt(exemptions, key, b.Status) {
			return false
		}
	}
	return true
}
