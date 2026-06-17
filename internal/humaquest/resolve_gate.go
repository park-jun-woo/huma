//ff:func feature=verify type=helper control=sequence
//ff:what Resolves the effective CRI gate: explicit testing.require_cri if set, else the reachable maximum (auto-ceiling).

package humaquest

import "github.com/park-jun-woo/huma/internal/config"

// resolveGate resolves the effective minimum CRI gate. An explicit
// testing.require_cri wins (a deliberate strictness choice); otherwise huma
// auto-requires the maximum CRI reachable in the current situation (§4). Ported
// from bak/cmd/resolve_gate.go.
func resolveGate(cfg *config.Config, reachableMax int) int {
	if cfg.RequireCRI != nil {
		return *cfg.RequireCRI
	}
	return reachableMax
}
