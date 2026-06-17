//ff:func feature=gate type=helper control=sequence level=error
//ff:what 측정된 adapter.CoverageResult를 gate.Context.Grounds["coverage"]에 실을 JSON 문자열로 직렬화한다. Evaluate의 decodeCoverage와 짝(왕복) — Percent/Total/CoveredLines를 보존한다.

package humaquest

import (
	"encoding/json"

	"github.com/park-jun-woo/huma/internal/adapter"
)

// coverageGround marshals a measured CoverageResult into the JSON string the cover
// command injects as ctx.Grounds["coverage"]. It is the write side of the contract
// Evaluate reads via decodeCoverage: presence signals a live run, Total==0 → SMOKE,
// Total>0 → COVERED-eligible, and CoveredLines feeds the branch-line binding check.
func coverageGround(cov *adapter.CoverageResult) (string, error) {
	b, err := json.Marshal(cov)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
