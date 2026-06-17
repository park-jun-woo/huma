//ff:func feature=gate type=engine control=sequence level=error
//ff:what 생성·기록·검증된 hurl을 측정해 verdict와 측정 coverage를 돌려준다: probe.Reset→Measure(*runner.Result 포착). result.Pass==false면 정적 CRI와 무관하게 runtimeFailVerdict(강제 비-PASS, §4 구멍 차단), 아니면 evaluateGenerated. cov는 호출자의 payload 영속화에 쓰인다.

package humaquest

import (
	"github.com/park-jun-woo/huma/internal/adapter"
	"github.com/park-jun-woo/huma/internal/scanner"
	"github.com/park-jun-woo/reins/pkg/gate"
	"github.com/park-jun-woo/reins/pkg/quest"
)

// measureGenerated measures the just-written hurl and returns the verdict plus the
// measured coverage (the caller persists the latter). It captures the *runner.Result
// that the manual path discards: when result.Pass==false it forces runtimeFailVerdict
// regardless of the static CRI (closing the §4 residual hole where a runtime-failing
// hurl could static-tier-1 PASS); otherwise it runs the normal evaluateGenerated path.
func measureGenerated(def gate.Definition, probe coverageProbe, s *quest.Session, it *quest.Item, ep scanner.Endpoint) (quest.Verdict, *adapter.CoverageResult, error) {
	if err := probe.Reset(); err != nil {
		return quest.Verdict{}, nil, err
	}
	cov, result, err := probe.Measure(ep)
	if err != nil {
		return quest.Verdict{}, nil, err
	}
	if result != nil && !result.Pass {
		return runtimeFailVerdict(ep, result.Feedback), cov, nil
	}
	verdict, err := evaluateGenerated(def, s, it, cov)
	return verdict, cov, err
}
