//ff:func feature=prompt type=builder control=sequence
//ff:what Builds the agent instruction for missing response status codes from static analysis
package prompt

import (
	"fmt"
	"strings"

	"github.com/park-jun-woo/huma/internal/hurlcheck"
	"github.com/park-jun-woo/huma/internal/scanner"
)

// ResponseImprovePrompt builds the agent instruction for missing response status codes.
func ResponseImprovePrompt(ep *scanner.Endpoint, hurlFile string, result *hurlcheck.ResponseCoverageResult) string {
	var b strings.Builder

	b.WriteString(fmt.Sprintf("# IMPROVE  %s %s\n", ep.Method, ep.Path))
	b.WriteString(fmt.Sprintf("# Response coverage: %d/%d (%.0f%%)\n", result.Covered, result.Total, result.Percent))

	if len(result.Missing) > 0 {
		b.WriteString("# MISSING:\n")
		for _, m := range result.Missing {
			b.WriteString(formatMissingBranch(m))
		}
	}

	b.WriteString("\n## Instructions\n\n")
	b.WriteString(fmt.Sprintf("1. Edit %s\n", hurlFile))
	b.WriteString("2. Add test entries for the missing status codes above\n")
	b.WriteString("3. Run `huma next`\n")

	return b.String()
}
