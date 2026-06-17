//ff:type feature=adapter type=adapter
//ff:what coverageProbe의 실구현. 선택된 adapter.Adapter와 manifest cfg(HurlDir/HurlVariables/server)를 명령 스코프에 소유한다(전역 상태 없음). Up=Build 1회(계측 바이너리 컴파일), Measure=엔드포인트마다 RunWithCoverage(Reset→Start→WaitReady→hurl→Stop→Collect), Down=최종 teardown. Go 커버리지는 프로세스 종료 시에만 flush되므로 기동/덤프가 엔드포인트 단위다.

package humaquest

import (
	"github.com/park-jun-woo/huma/internal/adapter"
	"github.com/park-jun-woo/huma/internal/config"
)

// liveProbe is the real coverageProbe: it drives a selected adapter.Adapter, owning the
// handle in the command scope (no global mutable state). Up compiles the instrumented
// binary once (adapter.Build); Measure runs the full Start → run hurl → Stop → Collect
// cycle PER endpoint via adapter.RunWithCoverage (Go coverage flushes only on process
// exit, so each endpoint needs its own process lifetime); Down is the final teardown.
// cfg supplies the .hurl directory, hurl variables, and server lifecycle commands.
type liveProbe struct {
	adapter adapter.Adapter
	cfg     *config.Config
}
