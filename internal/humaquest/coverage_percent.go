//ff:func feature=gate type=helper control=sequence
//ff:what 측정된 CoverageResult의 Percent를 안전하게 꺼낸다 — nil(라이브 신호 없음)이면 0. payloadState.PrevCoverage 기록용.

package humaquest

import "github.com/park-jun-woo/huma/internal/adapter"

// coveragePercent returns the measured line-coverage percent, or 0 when no live
// signal was collected (cov nil). It feeds the PrevCoverage the cover command
// persists for the next attempt's monotonicity check.
func coveragePercent(cov *adapter.CoverageResult) float64 {
	if cov == nil {
		return 0
	}
	return cov.Percent
}
