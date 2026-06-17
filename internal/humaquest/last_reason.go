//ff:func feature=gate type=helper control=sequence
//ff:what it.Log 마지막 Attempt의 Reason(Facts 렌더 문자열 — 실패 룰 ID/정량 증거 포함)을 반환한다. 로그가 없으면 빈 문자열. improve 프롬프트의 직전 미달 사유 힌트로 사용. 읽기 전용.

package humaquest

import "github.com/park-jun-woo/reins/pkg/quest"

// lastReason returns the Reason (the rendered Facts string, carrying failed rule
// IDs and quantified evidence) of the most recent logged Attempt, or "" when the
// item has no attempts. This is the improve hint — RootCause is never persisted,
// so the shortfall is recovered from here. Read-only.
func lastReason(it *quest.Item) string {
	if n := len(it.Log); n > 0 {
		return it.Log[n-1].Reason
	}
	return ""
}
