//ff:func feature=verify type=engine control=sequence
//ff:what Builds the IMPROVE verdict for a measured-but-below-gate run: OutFail to retry (Tries++), or OutReview at the MaxTries boundary without an unreachable.yaml reason (so DONE is never auto-locked unjustified).

package humaquest

import (
	"fmt"

	"github.com/park-jun-woo/huma/internal/analyzer"
	"github.com/park-jun-woo/huma/internal/scanner"
	"github.com/park-jun-woo/reins/pkg/gate"
	"github.com/park-jun-woo/reins/pkg/quest"
)



// improveVerdict builds the verdict when coverage was measured but is below the
// gate (some client branch uncovered). The default is quest.OutFail → a retry
// (quest.Apply does Tries++; at MaxTries with reasons it best-effort locks DONE).
// The MaxTries guard (boundaryNoReason) flips the last attempt to quest.OutReview
// when no unreachable.yaml reason justifies the gap, so Apply does not auto-lock
// DONE on an unjustified stall (§2, C-04). Facts carry the coverage percent and
// the uncovered status@line list (Phase 003 improve hint); prev feeds the
// stalled-vs-improving note. Ports bak coveredVerdict/staticVerdict IMPROVE +
// stalledVerdict.
func improveVerdict(ctx gate.Context, ep *scanner.Endpoint, branches, uncovered []analyzer.ResponseBranch, pct, prev float64) quest.Verdict {
	key := ep.Method + " " + ep.Path
	if boundaryNoReason(ctx, ep, uncovered) {
		return boundaryReviewVerdict(ep, uncovered, pct)
	}
	return quest.Verdict{
		Outcome:   quest.OutFail,
		RootCause: "C-03",
		Facts: []quest.Fact{{
			Rule:     "C-03",
			Where:    key,
			Expected: fmt.Sprintf("all %d client branch(es) covered", len(branches)),
			Actual:   fmt.Sprintf("coverage %.0f%% (prev %.0f%%), uncovered: %s", pct, prev, formatBranches(uncovered)),
		}},
		Feedback: improveFeedback(uncovered, pct, prev),
	}
}
