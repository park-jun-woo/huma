//ff:func feature=verify type=helper control=sequence
//ff:what Reports whether one client branch is runtime-bound (or exempt) to an unambiguous hit source line.

package humaquest

import (
	"github.com/park-jun-woo/huma/internal/adapter"
	"github.com/park-jun-woo/huma/internal/analyzer"
	"github.com/park-jun-woo/huma/internal/config"
)

// branchBound reports whether a single branch is runtime-bound: exempt via
// unreachable.yaml, or it has a unique source Line that was hit at runtime.
// Line==0 (no coordinate) and line cohesion (>1 branch on a line) both fail
// (strict, §3.4). Ported from bak/cmd/branch_bound.go.
func branchBound(b analyzer.ResponseBranch, cov *adapter.CoverageResult, exemptions []config.UnreachableEntry, key string, cohesion map[int]int) bool {
	if config.IsExempt(exemptions, key, b.Status) {
		return true
	}
	if b.Line == 0 || cohesion[b.Line] > 1 {
		return false
	}
	return cov.IsLineCovered(b.Line)
}
