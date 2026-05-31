//ff:func feature=verify type=helper control=sequence
//ff:what Resolves the CRI gate: explicit require_cri if set, else the reachable maximum
package cmd

import "github.com/park-jun-woo/huma/internal/config"

// resolveGate resolves the effective minimum CRI gate. An explicit
// testing.require_cri wins (a deliberate strictness choice); otherwise huma
// auto-requires the maximum CRI reachable in the current situation (§4).
func resolveGate(cfg *config.Config, reachableMax int) int {
	if cfg.RequireCRI != nil {
		return *cfg.RequireCRI
	}
	return reachableMax
}
