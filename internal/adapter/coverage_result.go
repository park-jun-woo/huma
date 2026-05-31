//ff:type feature=adapter type=model
//ff:what CoverageResult holds coverage analysis metrics, uncovered lines, and the runtime-covered line set for a handler
package adapter

// CoverageResult holds the coverage analysis for a handler function.
type CoverageResult struct {
	Covered   int
	Total     int
	Percent   float64
	Uncovered []UncoveredLine
	// CoveredLines is the set of source line numbers actually hit at runtime
	// within the handler range. Used for branch-line binding (§3.4): a status
	// branch is runtime-bound only if its source Line is in this set.
	CoveredLines map[int]bool
}
