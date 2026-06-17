//ff:func feature=verify type=helper control=iteration dimension=1
//ff:what Renders a set of branches as a compact "status@Lline" list for Facts/feedback (e.g. "404@L88, 409@L0").

package humaquest

import (
	"fmt"
	"strings"

	"github.com/park-jun-woo/huma/internal/analyzer"
)

// formatBranches renders branches as a compact "status@Lline" list for the
// IMPROVE Facts and feedback, so the rendered Attempt.Reason (Phase 003's
// improve hint) names exactly which status branches and source lines are
// uncovered. Returns "none" for an empty set.
func formatBranches(branches []analyzer.ResponseBranch) string {
	if len(branches) == 0 {
		return "none"
	}
	parts := make([]string, 0, len(branches))
	for _, b := range branches {
		parts = append(parts, fmt.Sprintf("%d@L%d", b.Status, b.Line))
	}
	return strings.Join(parts, ", ")
}
