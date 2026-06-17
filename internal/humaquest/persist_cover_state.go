//ff:func feature=gate type=helper control=sequence level=error
//ff:what 한 측정 후 IMPROVE 단조성 상태(PrevCoverage/ImproveCount)를 Item payload에 영속화한다. Evaluate는 read-only라 cover가 이 쓰기를 소유한다(Evaluate 직후, 호출자의 Apply/Save 전). 수동·generate 두 경로가 공유한다.

package humaquest

import (
	"github.com/park-jun-woo/huma/internal/adapter"
	"github.com/park-jun-woo/huma/internal/scanner"
	"github.com/park-jun-woo/reins/pkg/quest"
)

// persistCoverState writes the IMPROVE monotonicity (PrevCoverage from this run,
// ImproveCount++) into the Item payload. Evaluate is read-only, so cover owns this
// write — done after Evaluate and before the caller's Apply/Save. Both the manual
// (coverItem) and generate (generateItem) paths share it so the payload contract is
// identical.
func persistCoverState(it *quest.Item, ep scanner.Endpoint, cov *adapter.CoverageResult, ps payloadState) error {
	return it.SetPayload(payloadState{
		Endpoint:     ep,
		PrevCoverage: coveragePercent(cov),
		ImproveCount: ps.ImproveCount + 1,
	})
}
