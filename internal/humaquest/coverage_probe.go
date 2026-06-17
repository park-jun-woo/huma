//ff:type feature=gate type=model
//ff:what cover 명령의 부수효과 시임(injectable seam). 서버 라이프사이클을 한 인터페이스로 추상화한다: Up=계측 바이너리 컴파일 1회(Build), Measure=엔드포인트마다 Start→run hurl→Stop(커버리지 flush)→Collect, Down=최종 teardown. Go 커버리지는 프로세스 종료 시에만 flush되므로 기동/덤프가 엔드포인트 단위다. 시임 덕에 오케스트레이션(runCover)을 실서버 없이 가짜 probe로 단위 검증할 수 있다. 실구현은 liveProbe(어댑터+runner), 테스트는 canned CoverageResult를 돌려주는 fake.

package humaquest

import (
	"github.com/park-jun-woo/huma/internal/adapter"
	"github.com/park-jun-woo/huma/internal/runner"
	"github.com/park-jun-woo/huma/internal/scanner"
)

// coverageProbe is the side-effect seam the cover command orchestrates. It owns the
// session-scoped server lifecycle (Up brings it up once, Down tears it down) and the
// per-endpoint coverage measurement (Reset isolates, Measure runs the .hurl and
// collects coverage). Abstracting it as an interface lets runCover's orchestration —
// the loop, the Grounds["coverage"] injection, payloadState persistence, quest.Apply,
// and Export — be unit-tested with a fake that returns a canned CoverageResult, no
// real server required. The default real implementation is liveProbe.
type coverageProbe interface {
	// Up compiles the instrumented binary once for the whole loop (Build only). It
	// does not start the server — that happens per endpoint inside Measure.
	Up() error
	// Reset isolates coverage between endpoints. For liveProbe this is a no-op
	// (RunWithCoverage resets internally each measurement); fakes may use it.
	Reset() error
	// Measure runs one endpoint's full Start → run hurl → Stop (coverage flush) →
	// Collect cycle and returns the collected coverage. A nil CoverageResult means no
	// live signal (hurl missing/failed or handler bounds unknown) — the caller then
	// evaluates the static-only CRI.
	Measure(ep scanner.Endpoint) (*adapter.CoverageResult, *runner.Result, error)
	// Down is the final teardown (idempotent; safe to defer). The server is already
	// stopped per endpoint, so this only cleans up a process left by a mid-cycle error.
	Down() error
}
