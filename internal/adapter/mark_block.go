//ff:func feature=adapter type=engine control=iteration dimension=1
//ff:what Marks lines in a coverage block as total and optionally covered within the handler range
package adapter

import (
	"github.com/park-jun-woo/huma/internal/coverage"
)

func markBlock(b coverage.Block, startLine, endLine int, covered, total map[int]bool) {
	bStart := b.StartLine
	if bStart < startLine {
		bStart = startLine
	}
	bEnd := b.EndLine
	if bEnd > endLine {
		bEnd = endLine
	}
	for line := bStart; line <= bEnd; line++ {
		total[line] = true
		if b.Count > 0 {
			covered[line] = true
		}
	}
}
