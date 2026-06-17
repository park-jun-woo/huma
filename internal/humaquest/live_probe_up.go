//ff:func feature=adapter type=engine control=sequence level=error
//ff:what liveProbe.Up — 계측 바이너리를 한 번만 컴파일한다(adapter.Build). Start/WaitReady는 여기서 하지 않는다: Go 커버리지는 프로세스 종료 시에만 카운터를 flush하므로, 기동/측정/덤프는 엔드포인트마다 Measure 안의 RunWithCoverage가 담당한다. Build 실패(A-01)는 에러로 전파해 cover가 중단·teardown하게 한다.

package humaquest

import "fmt"

// Up compiles the instrumented binary once for the whole cover loop (adapter.Build
// only). It deliberately does NOT Start/WaitReady here: Go integration coverage
// (`go build -cover` / `go tool covdata`) only flushes counters on process exit, so a
// long-lived "up once" server yields no coverage. The Start → WaitReady → run hurl →
// Stop (flush) → Collect cycle happens PER endpoint inside Measure via
// adapter.RunWithCoverage. A Build failure (A-01) is returned so the cover command
// aborts and still tears down via the deferred Down.
func (p *liveProbe) Up() error {
	if err := p.adapter.Build(); err != nil {
		return fmt.Errorf("build server: %w", err)
	}
	return nil
}
