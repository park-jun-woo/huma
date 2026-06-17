//ff:func feature=verify type=helper control=sequence
//ff:what Builds the UNVERIFIED verdict (OutReview, RootCause C-01) — the honest non-pass for a no-signal/under-gated endpoint.

package humaquest

import "github.com/park-jun-woo/reins/pkg/quest"

// unverifiedVerdict builds the UNVERIFIED verdict: a no-signal or under-gated
// endpoint is never PASS (the §0 invariant, C-01). It maps to quest.OutReview →
// REVIEW (an oracle/human is needed), carrying a C-01 Fact so the rendered
// Attempt.Reason and the next-prompt explain the missing oracle/gate. Ports the
// bak sess.MarkUnverified outcome.
func unverifiedVerdict(where, expected, actual, feedback string) quest.Verdict {
	return quest.Verdict{
		Outcome:   quest.OutReview,
		RootCause: "C-01",
		Facts: []quest.Fact{{
			Rule:     "C-01",
			Where:    where,
			Expected: expected,
			Actual:   actual,
		}},
		Feedback: feedback,
	}
}
