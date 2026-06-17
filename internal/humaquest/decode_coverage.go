//ff:func feature=gate type=helper control=sequence level=error
//ff:what gate.Context.Grounds["coverage"]의 JSON 문자열을 adapter.CoverageResult로 디코드한다. 빈/부재(plain submit)면 (nil,false), 주입됨이면 (cov,true).

package humaquest

import (
	"encoding/json"

	"github.com/park-jun-woo/huma/internal/adapter"
)

// decodeCoverage decodes the coverage ground injected by Phase 006's cover
// command into ctx.Grounds["coverage"] (a JSON-encoded adapter.CoverageResult).
// The boolean reports presence: false on the plain `submit` path where no
// coverage was injected (raw is "") — the caller then computes a static-only
// CRI. A non-empty but malformed value is treated as absent (no panic) so a
// broken injection degrades to the honest static ceiling rather than a false
// SMOKE/COVERED.
func decodeCoverage(raw string) (*adapter.CoverageResult, bool) {
	if raw == "" {
		return nil, false
	}
	var cov adapter.CoverageResult
	if err := json.Unmarshal([]byte(raw), &cov); err != nil {
		return nil, false
	}
	return &cov, true
}
