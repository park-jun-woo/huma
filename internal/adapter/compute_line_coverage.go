//ff:func feature=adapter type=engine control=iteration dimension=1
//ff:what Computes covered and total line sets from filtered coverage blocks
package adapter

import (
	"github.com/park-jun-woo/hurlfill/internal/coverage"
)

func computeLineCoverage(filtered []coverage.Block, startLine, endLine int) (map[int]bool, map[int]bool) {
	covered := make(map[int]bool)
	total := make(map[int]bool)

	for _, b := range filtered {
		markBlock(b, startLine, endLine, covered, total)
	}
	return covered, total
}
