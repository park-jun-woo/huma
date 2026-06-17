//ff:func feature=gate type=helper control=sequence
//ff:what it.Log 마지막 Attempt의 Outcome(PASS/REVIEW/FAIL/…)를 반환한다. 로그가 없으면 빈 Outcome. derivePhase가 프롬프트 분기에 사용. 읽기 전용.

package humaquest

import "github.com/park-jun-woo/reins/pkg/quest"

// lastOutcome returns the Outcome of the most recent logged Attempt, or "" when
// the item has no attempts yet. Read-only.
func lastOutcome(it *quest.Item) quest.Outcome {
	if n := len(it.Log); n > 0 {
		return quest.Outcome(it.Log[n-1].Outcome)
	}
	return ""
}
