//ff:func feature=verify type=helper control=sequence
//ff:what Builds the MaxTries-boundary UNVERIFIED verdict (OutReview, C-04) when uncovered branches lack an unreachable.yaml reason at the last attempt.

package humaquest

import (
	"fmt"

	"github.com/park-jun-woo/huma/internal/analyzer"
	"github.com/park-jun-woo/huma/internal/scanner"
	"github.com/park-jun-woo/reins/pkg/quest"
)

// boundaryReviewVerdict is the verdict improveVerdict returns at the MaxTries
// boundary when no unreachable.yaml reason justifies the uncovered branches: a
// reason-less stall must not auto-lock DONE, so it stays UNVERIFIED (OutReview,
// C-04, §2/§3.7). Split out of improveVerdict to keep that function's control
// body small (Q4).
func boundaryReviewVerdict(ep *scanner.Endpoint, uncovered []analyzer.ResponseBranch, pct float64) quest.Verdict {
	return quest.Verdict{
		Outcome:   quest.OutReview,
		RootCause: "C-04",
		Facts: []quest.Fact{{
			Rule:     "C-04",
			Where:    ep.Method + " " + ep.Path,
			Expected: "full coverage, or an unreachable.yaml reason for every uncovered branch (to grant DONE)",
			Actual:   fmt.Sprintf("coverage %.0f%%, uncovered: %s, no reason artifact", pct, formatBranches(uncovered)),
		}},
		Feedback: fmt.Sprintf("UNVERIFIED: %d uncovered client branch(es) at the retry limit with no unreachable.yaml reason — add a verifiable reason (status + source evidence) or raise coverage. Uncovered: %s.", len(uncovered), formatBranches(uncovered)),
	}
}
