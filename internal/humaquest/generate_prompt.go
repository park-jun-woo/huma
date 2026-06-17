//ff:func feature=prompt type=builder control=sequence level=error
//ff:what generate 모드의 user 프롬프트를 만든다: Render(s,it)(read-only 권위 프롬프트)에, 재시도면(Tries>0) 정적 코칭 프리앰블(ruleSystem["C-03"])과 직전 시도의 판별 피드백(lastReason — Attempt.Reason)을 덧붙인다. 사람·모델 피드백 패리티의 1차 메커니즘은 lastReason이고 ruleSystem은 보조 프리앰블이다.

package humaquest

import (
	"github.com/park-jun-woo/reins/pkg/gate"
	"github.com/park-jun-woo/reins/pkg/quest"
)

// generatePrompt builds the user prompt for one generation attempt. The base is the
// definition's read-only Render (the same authoring prompt `huma next` shows). On a
// retry (it.Tries>0) it appends the generic static coaching preamble (ruleSystem
// ["C-03"]) and the discriminating feedback from the previous attempt — lastReason,
// the persisted Attempt.Reason. lastReason is the primary human/model feedback-parity
// mechanism; the ruleSystem preamble is auxiliary (§3).
func generatePrompt(def gate.Definition, s *quest.Session, it *quest.Item) (string, error) {
	base, err := def.Render(s, it)
	if err != nil {
		return "", err
	}
	if it.Tries > 0 {
		base += "\n\n## Retry — fix the previous failure\n" + ruleSystem["C-03"]
		if r := lastReason(it); r != "" {
			base += "\n\nPrevious verdict:\n" + r
		}
	}
	return base, nil
}
