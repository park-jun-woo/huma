//ff:func feature=verify type=engine control=sequence level=error
//ff:what 생성된 hurl이 런타임 실패(result.Pass==false)했을 때의 강제 비-PASS verdict. 게이트는 런타임 실패에 H-02/H-03을 emit하지 않으므로(§2-4) cover가 직접 합성한다: OutFail(재시도) + result.Feedback를 Facts.Actual에 실어 다음 시도의 lastReason으로 흐르게 하고, Feedback에는 정적 프리앰블(ruleSystem["H-03"])+전체 출력을 담는다. 정적 CRI가 tier1 PASS를 줄 수 있는 §4 잔존 구멍을 닫는다(런타임 실패는 정적 등급과 무관하게 비-PASS).

package humaquest

import (
	"strings"

	"github.com/park-jun-woo/huma/internal/scanner"
	"github.com/park-jun-woo/reins/pkg/quest"
)

// runtimeFailVerdict synthesizes the forced non-PASS verdict for a generated hurl
// that ran but FAILED (result.Pass==false). The gate never emits H-02/H-03 for a
// runtime failure (§2-4): a failing hurl yields cov=nil and would otherwise fall to
// the static-only CRI, which can hand out a tier-1 PASS (§4 residual hole). So cover
// forces OutFail here regardless of the static CRI. The runtime output is embedded in
// the Fact's Actual (capped) so quest.Apply persists it into Attempt.Reason — that is
// what lastReason feeds back to the next generation. Feedback carries the static
// preamble plus the full output for display.
func runtimeFailVerdict(ep scanner.Endpoint, feedback string) quest.Verdict {
	key := ep.Method + " " + ep.Path
	actual := strings.TrimSpace(feedback)
	if len(actual) > 1500 {
		actual = actual[:1500] + " …(truncated)"
	}
	return quest.Verdict{
		Outcome:   quest.OutFail,
		RootCause: "H-03",
		Facts: []quest.Fact{{
			Rule:     "H-03",
			Where:    key,
			Expected: "hurl run passes (every asserted branch holds)",
			Actual:   "hurl FAILED at runtime — " + actual,
		}},
		Feedback: ruleSystem["H-03"] + "\n\n" + strings.TrimSpace(feedback),
	}
}
