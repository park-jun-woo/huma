//ff:type feature=coverage type=model
//ff:what CoveragePyReport represents the top-level structure of coverage.py JSON output
package coverage

// coveragePyReport represents the top-level structure of `coverage json` output.
type coveragePyReport struct {
	Files map[string]coveragePyFile `json:"files"`
}
