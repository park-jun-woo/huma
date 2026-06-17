//ff:func feature=verify type=builder control=sequence
//ff:what Builds the IMPROVE feedback text: which branches to cover next, plus a stalled-vs-improving note from the prev/current coverage comparison.

package humaquest

import (
	"fmt"

	"github.com/park-jun-woo/huma/internal/analyzer"
)

// improveFeedback renders the human IMPROVE hint fed back to the model. It names
// the uncovered branches and adds a monotonicity note (§2): coverage that did
// not rise over the previous attempt is flagged as stalled so the agent changes
// strategy instead of resubmitting the same test.
func improveFeedback(uncovered []analyzer.ResponseBranch, pct, prev float64) string {
	trend := fmt.Sprintf("coverage %.0f%%", pct)
	if prev > 0 {
		if pct <= prev {
			trend = fmt.Sprintf("coverage stalled at %.0f%% (was %.0f%%)", pct, prev)
		} else {
			trend = fmt.Sprintf("coverage %.0f%% (up from %.0f%%)", pct, prev)
		}
	}
	return fmt.Sprintf("IMPROVE: %s — add hurl entries that exercise and assert the uncovered branches: %s.",
		trend, formatBranches(uncovered))
}
