//ff:func feature=prompt type=builder control=sequence
//ff:what Builds the agent instruction for an endpoint that passes but has uncovered lines
package prompt

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/park-jun-woo/hurlfill/internal/adapter"
	"github.com/park-jun-woo/hurlfill/internal/scanner"
)

// ImprovePrompt builds the agent instruction for an endpoint that passes
// but has uncovered lines.
func ImprovePrompt(ep *scanner.Endpoint, hurlFile string, result *adapter.CoverageResult) string {
	var b strings.Builder

	b.WriteString(fmt.Sprintf("# IMPROVE  %s %s\n", ep.Method, ep.Path))
	b.WriteString(fmt.Sprintf("# Coverage: %.0f%% (%d/%d)\n", result.Percent, result.Covered, result.Total))

	if len(result.Uncovered) > 0 {
		b.WriteString("# UNCOVERED:\n")
		for _, u := range result.Uncovered {
			base := filepath.Base(u.File)
			b.WriteString(fmt.Sprintf("#   %s:%d  %s\n", base, u.Line, u.Code))
		}
	}

	b.WriteString("\n## Instructions\n\n")
	b.WriteString(fmt.Sprintf("1. Edit %s\n", hurlFile))
	b.WriteString("2. Add test entries for the uncovered branches above\n")
	b.WriteString("3. Run `hurlfill next`\n")

	return b.String()
}
