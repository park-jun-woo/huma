//ff:func feature=gate type=engine control=iteration dimension=1 level=error
//ff:what cover 오케스트레이션. 세션 로드→probe.Up(계측 바이너리 1회 컴파일, defer Down)→TODO 순회(각 아이템: coverItem 측정·판정→commitVerdict 래칫·Save·Export→출력)→끝. backend!=nil이면 coverItem이 generate 모드로 동작한다(별도 내부 재시도 루프 없음 — FAIL은 TODO로 남아 NextTODO가 재선택, MaxTries→DONE). maxItems>0이면 TODO를 떠난 distinct endpoint 수가 그 수에 닿을 때 멈춘다(attempt가 아니라 endpoint 카운트). 서버 기동/덤프는 엔드포인트마다 Measure 안의 RunWithCoverage가 담당한다(Go 커버리지는 프로세스 종료 시 flush). probe가 시임이라 가짜로 단위 검증 가능.

package humaquest

import (
	"io"

	"github.com/park-jun-woo/reins/pkg/gate"
	"github.com/park-jun-woo/reins/pkg/llm"
	"github.com/park-jun-woo/reins/pkg/quest"
)

// runCover is the cover command's orchestration: it loads the session, compiles the
// instrumented binary ONCE via probe.Up (deferring Down so teardown happens on
// error/panic), then loops over every TODO item — measuring coverage, injecting it,
// evaluating, persisting monotonicity (coverItem), and applying+saving+exporting the
// verdict (commitVerdict). Each Measure runs its own Start → run hurl → Stop → Collect
// cycle (Go coverage flushes only on process exit), so IMPROVE retries re-measure with
// a fresh process (a failing item stays TODO and NextTODO re-selects it until MaxTries
// locks it DONE). The probe is the injectable seam, so this loop is unit-testable with
// a fake adapter and a canned CoverageResult.
func runCover(def gate.Definition, probe coverageProbe, backend llm.Backend, maxItems int, sessionPath, outPath string, w io.Writer) error {
	s, err := loadCoverSession(sessionPath)
	if err != nil {
		return err
	}

	if err := probe.Up(); err != nil {
		return err
	}
	defer probe.Down()

	sink, err := newJSONLSink(outPath)
	if err != nil {
		return err
	}

	done := 0
	for it := s.NextTODO(); it != nil; it = s.NextTODO() {
		verdict, err := coverItem(def, probe, s, it, backend)
		if err != nil {
			return err
		}
		if err := commitVerdict(s, it, verdict, sink, sessionPath); err != nil {
			return err
		}
		renderCoverVerdict(w, it.Key, it, verdict)
		// --max-items counts DISTINCT endpoints, not attempts: a FAILed item stays
		// TODO (NextTODO re-selects it), so only count when it leaves TODO.
		if it.State == quest.TODO {
			continue
		}
		done++
		if maxItems > 0 && done >= maxItems {
			return nil
		}
	}
	return nil
}
