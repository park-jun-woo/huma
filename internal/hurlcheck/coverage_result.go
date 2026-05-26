//ff:type feature=hurlcheck type=model
//ff:what ResponseCoverageResult holds the comparison between expected and tested HTTP status codes
package hurlcheck

import "github.com/park-jun-woo/huma/internal/analyzer"

// ResponseCoverageResult holds the static response coverage analysis.
type ResponseCoverageResult struct {
	Covered int
	Total   int
	Percent float64
	Missing []analyzer.ResponseBranch
}
