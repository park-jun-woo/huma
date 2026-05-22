//ff:type feature=adapter type=model
//ff:what CoverageResult holds coverage analysis metrics and uncovered lines for a handler
package adapter

// CoverageResult holds the coverage analysis for a handler function.
type CoverageResult struct {
	Covered   int
	Total     int
	Percent   float64
	Uncovered []UncoveredLine
}
