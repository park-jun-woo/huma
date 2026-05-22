//ff:type feature=coverage type=model
//ff:what CoveragePyFile represents per-file coverage data from coverage.py JSON output
package coverage

// coveragePyFile represents per-file coverage data from coverage.py.
type coveragePyFile struct {
	ExecutedLines []int `json:"executed_lines"`
	MissingLines  []int `json:"missing_lines"`
}
