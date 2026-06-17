//ff:func feature=gate type=engine control=sequence level=error
//ff:what 런타임 통과한 생성 hurl을 수동 경로와 동일하게 판정한다: def.Prepare(short 존중)→측정 coverage가 있으면 ctx.Grounds["coverage"] 주입→Evaluate. 수동 coverItem의 측정-후 꼬리와 같은 게이트 권위(PASS는 게이트만)를 generate 분기가 재사용한다.

package humaquest

import (
	"github.com/park-jun-woo/huma/internal/adapter"
	"github.com/park-jun-woo/reins/pkg/gate"
	"github.com/park-jun-woo/reins/pkg/quest"
)

// evaluateGenerated runs the normal post-measurement gate path for a generated hurl
// that PASSED at runtime, mirroring coverItem's manual tail: Prepare (honoring a
// short-circuit) → inject the measured coverage as Grounds["coverage"] when present →
// Evaluate. The generate branch reuses the exact same gate authority (only the gate
// originates PASS).
func evaluateGenerated(def gate.Definition, s *quest.Session, it *quest.Item, cov *adapter.CoverageResult) (quest.Verdict, error) {
	ctx, short, err := def.Prepare(s, it, nil)
	if err != nil {
		return quest.Verdict{}, err
	}
	if short != nil {
		return *short, nil
	}
	if cov != nil {
		ground, gerr := coverageGround(cov)
		if gerr != nil {
			return quest.Verdict{}, gerr
		}
		ctx.Grounds = map[string]string{"coverage": ground}
	}
	return evaluate(def, ctx), nil
}
