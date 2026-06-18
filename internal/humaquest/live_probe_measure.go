//ff:func feature=adapter type=engine control=sequence level=error
//ff:what liveProbe.Measure — 한 엔드포인트에 대해 adapter.RunWithCoverage에 위임한다: Reset→Start→WaitReady→run hurl→Stop(커버리지 flush 유발)→Collect를 엔드포인트마다 수행한다. Go 커버리지는 프로세스 종료 시에만 flush되므로 Stop이 매번 필요하다. .hurl 부재면 cov=nil(라이브 신호 없음)로 돌려 Evaluate가 정적 CRI를 내게 한다. hurl 변수는 mergeVars(cfg.HurlVariables, extraVars)로 합쳐 캡처/mint 변수(token·픽스처)가 정적 변수를 덮어쓰며, manual·--generate 두 경로 모두 이 주입을 받는다(둘 다 Measure 경유, Phase 009).

package humaquest

import (
	"github.com/park-jun-woo/huma/internal/adapter"
	"github.com/park-jun-woo/huma/internal/runner"
	"github.com/park-jun-woo/huma/internal/scanner"
)

// Measure runs one endpoint's full coverage cycle by delegating to
// adapter.RunWithCoverage: Reset → Start → WaitReady → run the .hurl → Stop → Collect.
// The Stop is what triggers Go's coverage dump (counters flush only on process exit),
// so a fresh Start/Stop happens for every endpoint — including each IMPROVE retry.
// It returns a nil CoverageResult (no live signal) when the .hurl is missing, the hurl
// run fails, or the handler bounds can't be read — the caller then evaluates the
// static-only CRI. The instrumented binary was already compiled once by Up.
func (p *liveProbe) Measure(ep scanner.Endpoint) (*adapter.CoverageResult, *runner.Result, error) {
	hurlPath := runner.FindHurlFile(&ep, p.cfg.HurlDir)
	if hurlPath == "" {
		return nil, nil, nil
	}

	vars := mergeVars(p.cfg.HurlVariables, p.extraVars)
	result, cov, err := adapter.RunWithCoverage(p.adapter, hurlPath, vars, ep.Source, ep.Handler)
	if err != nil {
		return nil, result, err
	}
	return cov, result, nil
}
