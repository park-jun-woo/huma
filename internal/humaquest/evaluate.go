//ff:func feature=gate type=engine control=sequence
//ff:what Definition.Evaluate (gate.Evaluator). huma의 4단계 CRI verdict를 한 함수로: 분모(소스∪선언 합집합·단조)·A-grade·주입된 coverage를 모아 computeCRI로 staged O/D/E 천장을 내고, 4번째 축 A는 verdictFromCRI에서 fold해 quest.Verdict에 매핑한다(전체 CRI=min(O,D,A,E)). 순수/읽기전용 — ctx.Grounds["coverage"]만 읽고 서버·디스크 쓰기 없음(Phase 006 cover가 부수효과 담당).

package humaquest

import (
	"github.com/park-jun-woo/huma/internal/config"
	"github.com/park-jun-woo/reins/pkg/gate"
	"github.com/park-jun-woo/reins/pkg/quest"
)

// Evaluate is huma's gate.Evaluator path. Because humaDef implements
// gate.Evaluator, reins dispatches here instead of the Rules() catalog
// (cli/evaluate_and_apply.go: `if ev, ok := def.(gate.Evaluator); ok {
// verdict = ev.Evaluate(ctx) }`). It ports the bak staged CRI pipeline
// (static→live→smoke|covered→stalled) into one pure, read-only function:
//
//  1. Denominator: responseBranches — the monotonic union of source branches
//     (authoritative floor) and OpenAPI declarations (additive), split to the
//     gated client set, carrying provenance (§3.1/C-02).
//  2. A-grade: minimum assertion depth over the union, reusing Prepare's entries.
//  3. Coverage: ctx.Grounds["coverage"] — present only when Phase 006's cover
//     command injected a measured run; absent on plain submit → static-only CRI.
//  4. CRI: computeCRI yields the staged O/D/E evidence ceiling; the fourth axis
//     A (assertion depth) is folded on top in verdictFromCRI, so the full
//     CRI = min(O,D,A,E) is realized at the verdict stage, not in computeCRI.
//  5. Verdict: verdictFromCRI maps the tier through the require_cri gate to
//     PASS / IMPROVE(retry) / UNVERIFIED, honoring the MaxTries→DONE boundary,
//     and routes an A-capped tier (aGrade<cri) to the assertion-depth IMPROVE.
//
// Side-effect free: it reads the manifest, unreachable.yaml, and the injected
// coverage string only — no server lifecycle, no Payload writes (Phase 006 owns
// those). Given the coverage string it is deterministic and network-free.
func (humaDef) Evaluate(ctx gate.Context) quest.Verdict {
	sub, ok := ctx.Submission.(*hurlInfo)
	if !ok || sub == nil {
		return unverifiedVerdict("submission", "a decoded *hurlInfo submission", "missing or wrongly-typed submission", "internal: no hurl submission to evaluate")
	}
	ep := &sub.Endpoint

	cfg, err := config.Load()
	if err != nil {
		// No manifest → static-mode defaults, mirroring Prepare/Render.
		cfg = &config.Config{HurlDir: "hurl"}
	}

	branches, prov := responseBranches(ep, cfg.Scan.Lang)
	cov, covPresent := decodeCoverage(ctx.Grounds["coverage"])
	aGrade := staticAGrade(sub.Entries, branches)
	prev := priorCoverage(ctx.Item)

	cri := computeCRI(ep, branches, cov, covPresent)
	return verdictFromCRI(ctx, cfg, cri, ep, sub, branches, cov, prov.String(), aGrade, prev)
}
