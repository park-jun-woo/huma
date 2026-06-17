//ff:func feature=adapter type=engine control=sequence level=error
//ff:what liveProbe.Down — 최종 정리(final cleanup)다. 정상 흐름에선 서버가 Measure 안의 RunWithCoverage가 매 엔드포인트마다 Stop으로 이미 내려가 있다. Down은 중도 실패 시 남았을 수 있는 프로세스를 안전하게 닫는 idempotent teardown(Stop은 핸들이 nil이면 no-op).

package humaquest

// Down is the final cleanup. In the normal flow the server is already stopped after
// every endpoint (RunWithCoverage Stops to flush coverage), so there is no per-endpoint
// Stop left to do here. Down still calls the adapter's idempotent Stop (no-op when the
// handle is nil) so a process left running by a mid-cycle failure is torn down safely;
// it is safe to defer even when Up failed.
func (p *liveProbe) Down() error {
	return p.adapter.Stop()
}
