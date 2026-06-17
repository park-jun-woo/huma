//ff:func feature=gate type=helper control=sequence
//ff:what reins evaluateAndApply의 판정 분기를 export 프리미티브로 재현한다: def가 gate.Evaluator면 그 Evaluate(huma의 CRI 경로), 아니면 평탄 gate.Evaluate(Rules). cover가 unexported evaluateAndApply 대신 동일 권한 비대칭(게이트만 PASS)으로 판정하기 위한 글루.

package humaquest

import (
	"github.com/park-jun-woo/reins/pkg/gate"
	"github.com/park-jun-woo/reins/pkg/quest"
)

// evaluate dispatches a prepared Context to the same verdict source reins'
// evaluateAndApply uses: the gate.Evaluator path when def implements it (huma's CRI
// Evaluate), else the flat gate.Evaluate over the rule catalog. It exists because
// evaluateAndApply is unexported, so cover replicates this one branch with exported
// primitives only.
func evaluate(def gate.Definition, ctx gate.Context) quest.Verdict {
	if ev, ok := def.(gate.Evaluator); ok {
		return ev.Evaluate(ctx)
	}
	return gate.Evaluate(def.Rules(), ctx)
}
