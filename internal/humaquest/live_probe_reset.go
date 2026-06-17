//ff:func feature=adapter type=engine control=sequence level=error
//ff:what liveProbe.Reset — no-op. 엔드포인트 간 커버리지 격리는 Measure가 위임하는 adapter.RunWithCoverage 내부의 Reset이 매 측정마다 수행하므로 여기선 할 일이 없다. coverageProbe 인터페이스 안정성을 위해 메서드는 남겨둔다(fake probe·오케스트레이션 불변).

package humaquest

// Reset is a no-op for liveProbe: adapter.RunWithCoverage (called by Measure) resets
// coverage internally at the start of every per-endpoint cycle, so there is nothing to
// isolate here. The method is kept to preserve the coverageProbe interface (the fake
// probe and runCover/coverItem orchestration stay unchanged).
func (p *liveProbe) Reset() error {
	return nil
}
