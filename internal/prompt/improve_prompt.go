//ff:func feature=prompt type=builder control=sequence
//ff:what Builds the agent instruction for a retry after a FAIL: echoes the prior shortfall (rendered Facts — failed rule IDs / uncovered branches) and tells the agent to edit the .hurl and rerun
package prompt

import (
	"fmt"
	"strings"

	"github.com/park-jun-woo/huma/internal/scanner"
)

// ImprovePrompt builds the agent instruction for an endpoint whose previous
// attempt fell short (a FAIL). reason is the rendered Facts string from the last
// Attempt.Reason — it already carries the failed rule IDs and the quantified
// shortfall (coverage numbers, uncovered branches), so the builder echoes it
// verbatim rather than re-reading a CoverageResult.
func ImprovePrompt(ep *scanner.Endpoint, hurlFile, reason string) string {
	var b strings.Builder

	b.WriteString(fmt.Sprintf("# IMPROVE  %s %s\n", ep.Method, ep.Path))

	if reason != "" {
		b.WriteString("# Previous attempt fell short:\n")
		for _, line := range strings.Split(strings.TrimRight(reason, "\n"), "\n") {
			b.WriteString("#   " + line + "\n")
		}
	}

	b.WriteString("\n## Instructions\n\n")
	b.WriteString(fmt.Sprintf("1. Edit %s\n", hurlFile))
	b.WriteString("2. Add test entries for the uncovered branches / missing statuses above\n")
	b.WriteString("3. Run `huma next`\n")

	return b.String()
}
