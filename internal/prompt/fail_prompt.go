//ff:func feature=prompt type=builder control=sequence
//ff:what Builds the agent instruction text for a FAIL endpoint with feedback
package prompt

import (
	"fmt"
	"strings"

	"github.com/park-jun-woo/hurlfill/internal/scanner"
)

// FailPrompt builds the agent instruction for a FAIL endpoint.
func FailPrompt(ep *scanner.Endpoint, hurlFile string, feedback string) string {
	var b strings.Builder

	b.WriteString(fmt.Sprintf("# FAIL  %s %s\n", ep.Method, ep.Path))
	b.WriteString(fmt.Sprintf("# File: %s\n", hurlFile))
	b.WriteString("\n")
	b.WriteString(feedback)
	b.WriteString("\n")
	b.WriteString("\n## Instructions\n\n")
	b.WriteString(fmt.Sprintf("1. Edit %s\n", hurlFile))
	b.WriteString("2. Fix the failing assertion\n")
	b.WriteString("3. Run `hurlfill next`\n")

	return b.String()
}
