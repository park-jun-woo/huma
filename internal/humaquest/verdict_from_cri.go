//ff:func feature=verify type=engine control=selection
//ff:what Maps the staged O/D/E ceiling to a quest.Verdict, first folding the A axis: cri>0 & aGrade<cri routes to the assertion-depth IMPROVE(C-03); else 0→UNVERIFIED, 3→COVERED gate, 2→SMOKE/IMPROVE, 1→SCAFFOLDED gate or static IMPROVE.

package humaquest

import (
	"fmt"

	"github.com/park-jun-woo/huma/internal/adapter"
	"github.com/park-jun-woo/huma/internal/analyzer"
	"github.com/park-jun-woo/huma/internal/config"
	"github.com/park-jun-woo/huma/internal/hurlcheck"
	"github.com/park-jun-woo/huma/internal/scanner"
	"github.com/park-jun-woo/reins/pkg/gate"
	"github.com/park-jun-woo/reins/pkg/quest"
)

// verdictFromCRI maps the achieved CRI tier (computeCRI) to a quest.Verdict,
// porting the bak staged pipeline (static/smoke/covered/stalled) into one switch:
//
//	0 → UNVERIFIED (OutReview): no oracle / no denominator (C-01).
//	3 → COVERED: 100% + branch-bound; gateVerdict grants PASS unless require_cri>3.
//	2 → instrumented-but-incomplete → IMPROVE; uninstrumented green → SMOKE gate.
//	1 → static: missing branch statuses → IMPROVE; all covered → SCAFFOLDED gate.
//
// IMPROVE returns OutFail (retry) or, at the MaxTries boundary without a reason,
// OutReview (so Apply never auto-locks an unjustified DONE).
//
// A-axis fold (Phase 007, §1/§3.3): cri is only the staged O/D/E ceiling; the
// fourth CRI axis A (assertion depth, aGrade) is folded here. The display tier
// is effective = min(cri, aGrade). When A is the limiting axis (cri>0 &&
// aGrade<cri) the endpoint has an oracle and execution/denominator evidence but
// shallow assertions, so it routes to the dedicated assertion-depth IMPROVE
// (C-03) BEFORE the staged switch — this preempts a PASS@low-tier that would
// otherwise cap silently. The cri>0 guard keeps case 0 (no oracle) on the
// UNVERIFIED/C-01 path, so A=0 WITH an oracle is C-03 IMPROVE, not C-01. A
// normal low-staged PASS (aGrade==cri) is not aGrade<cri, so it is untouched.
func verdictFromCRI(ctx gate.Context, cfg *config.Config, cri int, ep *scanner.Endpoint, sub *hurlInfo, branches []analyzer.ResponseBranch, cov *adapter.CoverageResult, prov string, aGrade int, prev float64) quest.Verdict {
	effective := cri
	if aGrade < effective {
		effective = aGrade // display tier = min(staged O/D/E, A)
	}
	if cri > 0 && aGrade < cri {
		// A is the limiting axis: deepen assertions first (preempts the staged
		// switch, so dual coverage+assertion deficiency at cri==2 deepens A first).
		return assertionImproveVerdict(ctx, ep, branches, cri, effective)
	}
	switch cri {
	case 0:
		return unverifiedVerdict(
			ep.Method+" "+ep.Path,
			"an independent oracle (source link or instrumented server) and >=1 client branch",
			fmt.Sprintf("source=%q, runtime-instrumented=%t, client branches=%d", ep.Source, cov != nil && cov.Total > 0, len(branches)),
			"UNVERIFIED: no independent oracle. fix: huma scan --link-source <root> OR set testing.server (+ instrumented build).",
		)
	case 3:
		return gateVerdict(cfg, 3, ep, branches, aGrade, prov)
	case 2:
		if cov != nil && cov.Total > 0 {
			return improveVerdict(ctx, ep, branches, uncoveredBranches(branches, cov), cov.Percent, prev)
		}
		return gateVerdict(cfg, 2, ep, branches, aGrade, prov)
	default:
		respResult := hurlcheck.CheckResponseCoverage(branches, hurlcheck.NonVacuousStatusList(sub.Entries))
		if respResult != nil && respResult.Total > 0 && len(respResult.Missing) > 0 {
			return improveVerdict(ctx, ep, branches, respResult.Missing, respResult.Percent, prev)
		}
		return gateVerdict(cfg, 1, ep, branches, aGrade, prov)
	}
}
