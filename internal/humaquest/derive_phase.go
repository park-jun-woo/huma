//ff:func feature=gate type=helper control=selection
//ff:what huma 세션 Status를 reins Item 상태로 재유도해 프롬프트 종류를 고른다: 직전 Outcome=REVIEW→unverified, =FAIL→improve, 그 외 fresh는 manifest 모드로 static/todo. Render는 이 결정만 읽고 부작용 없음.

package humaquest

import "github.com/park-jun-woo/reins/pkg/quest"

// derivePhase maps the Item's own signals to a prompt variant. It reads only the
// last Attempt Outcome and the manifest mode; it never inspects a RootCause
// field (RootCause is not persisted — see Phase003 doc §2). A FAIL implies
// Tries > 0; a fresh item (no log) picks static vs todo by server config.
func derivePhase(it *quest.Item, staticMode bool) phase {
	switch lastOutcome(it) {
	case quest.OutReview:
		return phaseUnverified
	case quest.OutFail:
		return phaseImprove
	default:
		if staticMode {
			return phaseStatic
		}
		return phaseTodo
	}
}
