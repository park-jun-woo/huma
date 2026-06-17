//ff:func feature=verify type=helper control=sequence
//ff:what Compares an achieved CRI tier against the resolved gate: PASS at the tier when met, else UNVERIFIED (an explicit require_cri exceeds the reachable evidence).

package humaquest

import (
	"fmt"

	"github.com/park-jun-woo/huma/internal/analyzer"
	"github.com/park-jun-woo/huma/internal/config"
	"github.com/park-jun-woo/huma/internal/scanner"
	"github.com/park-jun-woo/reins/pkg/quest"
)

// gateVerdict applies the CRI gate to an achieved tier (the bak `tier <
// resolveGate(cfg, tier)` check). With no explicit testing.require_cri the gate
// equals the achieved tier, so PASS is granted at that tier. An explicit
// require_cri higher than the reachable evidence yields UNVERIFIED — huma refuses
// false comfort and demands a stronger oracle (§4). PASS → quest.OutPass (the
// only PASS lock); the shortfall → unverifiedVerdict (OutReview).
func gateVerdict(cfg *config.Config, tier int, ep *scanner.Endpoint, branches []analyzer.ResponseBranch, aGrade int, prov string) quest.Verdict {
	gate := resolveGate(cfg, tier)
	if tier < gate {
		return unverifiedVerdict(
			ep.Method+" "+ep.Path,
			fmt.Sprintf("CRI >= %d (testing.require_cri)", gate),
			fmt.Sprintf("achieved CRI %d (%s) — evidence below the required gate", tier, criLabel(tier)),
			fmt.Sprintf("UNVERIFIED: %s reaches CRI %d but require_cri=%d. Raise evidence (instrument the server / link source / deepen assertions) or lower require_cri.", criLabel(tier), tier, gate),
		)
	}
	return passVerdict(tier, aGrade, prov, branches)
}
