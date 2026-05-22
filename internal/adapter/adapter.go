package adapter

// UncoveredLine represents a single line of source code that was not covered.
type UncoveredLine struct {
	File string
	Line int
	Code string
}

// CoverageResult holds the coverage analysis for a handler function.
type CoverageResult struct {
	Covered   int
	Total     int
	Percent   float64
	Uncovered []UncoveredLine
}

// Adapter manages a server process for coverage collection.
type Adapter interface {
	Build() error
	Start() error
	WaitReady() error
	Stop() error
	Collect(handlerFile string, startLine, endLine int) (*CoverageResult, error)
	Reset() error
}
