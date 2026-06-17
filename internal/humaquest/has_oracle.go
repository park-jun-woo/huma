//ff:func feature=verify type=helper control=sequence
//ff:what Reports whether an independent ground-truth oracle (source link or runtime instrumentation) exists for the endpoint.

package humaquest

import (
	"github.com/park-jun-woo/huma/internal/adapter"
	"github.com/park-jun-woo/huma/internal/scanner"
)

// hasOracle reports whether an independent ground-truth oracle exists for the
// endpoint: a linked source (for branch analysis) or live runtime
// instrumentation (Total>0). No oracle ⇒ UNVERIFIED (§3.2, CV-2/CV-3). Ported
// from bak/cmd/has_oracle.go.
func hasOracle(ep *scanner.Endpoint, cov *adapter.CoverageResult) bool {
	if ep.Source != "" {
		return true
	}
	return cov != nil && cov.Total > 0
}
