//ff:type feature=coverage type=model
//ff:what IstanbulFileCoverage represents coverage data for a single file in istanbul JSON format
package coverage

// istanbulFileCoverage represents coverage data for a single file.
type istanbulFileCoverage struct {
	Path         string                   `json:"path"`
	StatementMap map[string]istanbulRange `json:"statementMap"`
	S            map[string]int           `json:"s"`
}
