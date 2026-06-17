//ff:func feature=adapter type=helper control=sequence
//ff:what manifest cfg로부터 liveProbe를 만든다 — selectAdapter로 언어별 어댑터를 골라 cfg와 함께 묶는다.

package humaquest

import "github.com/park-jun-woo/huma/internal/config"

// newLiveProbe builds the real coverageProbe for cfg, selecting the language
// adapter via selectAdapter. It is the default probe wired by the cover command.
func newLiveProbe(cfg *config.Config) *liveProbe {
	return &liveProbe{adapter: selectAdapter(cfg), cfg: cfg}
}
