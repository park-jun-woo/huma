//ff:func feature=verify type=engine control=sequence
//ff:what Builds the assertion-depth IMPROVE verdict (C-03) when A is the limiting CRI axis (aGrade<staged): OutFail to deepen assertions, or OutReview at the MaxTries boundary without an unreachable.yaml reason.

package humaquest

import (
	"fmt"

	"github.com/park-jun-woo/huma/internal/analyzer"
	"github.com/park-jun-woo/huma/internal/scanner"
	"github.com/park-jun-woo/reins/pkg/gate"
	"github.com/park-jun-woo/reins/pkg/quest"
)

// assertionImproveVerdict builds the verdict when the A (assertion-depth) axis
// caps CRI below the staged O/D/E ceiling (Phase 007, §3.3/§6): the endpoint has
// an oracle and execution/denominator evidence but the hurl entries assert too
// shallowly (e.g. status only, A=1). Unlike improveVerdict (coverage-shaped:
// uncovered/pct/prev — unfit for an A gap because the branches ARE covered),
// this routes a covered-but-shallow endpoint back to deepen its assertions toward
// A=3 (status + body shape + invariants). It maps to quest.OutFail → a retry;
// the shared MaxTries guard (boundaryNoReason, reused not duplicated) flips the
// last attempt to quest.OutReview when no unreachable.yaml reason justifies the
// gap, so Apply never auto-locks an unjustified DONE (§2, C-04). effective =
// min(staged,A) is the capped display tier; the Fact and Feedback expose the cap
// transparency (A caps CRI to effective while the staged ceiling is cri).
func assertionImproveVerdict(ctx gate.Context, ep *scanner.Endpoint, branches []analyzer.ResponseBranch, cri, effective int) quest.Verdict {
	if boundaryNoReason(ctx, ep, branches) {
		return boundaryReviewVerdict(ep, branches, 100)
	}
	return quest.Verdict{
		Outcome:   quest.OutFail,
		RootCause: "C-03",
		Facts: []quest.Fact{{
			Rule:     "C-03",
			Where:    ep.Method + " " + ep.Path,
			Expected: "assertion depth A=3 (status + body shape + invariants)",
			Actual:   fmt.Sprintf("A=%d caps CRI to %d (staged %d) — status only, no body/invariant assertions", effective, effective, cri),
		}},
		Feedback: assertionImproveFeedback(ep, cri, effective),
	}
}
