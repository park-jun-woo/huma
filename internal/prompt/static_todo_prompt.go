//ff:func feature=prompt type=builder control=sequence
//ff:what Builds the agent instruction text for a TODO endpoint in static mode with response branches
package prompt

import (
	"fmt"
	"strings"

	"github.com/park-jun-woo/huma/internal/analyzer"
	"github.com/park-jun-woo/huma/internal/runner"
	"github.com/park-jun-woo/huma/internal/scanner"
	"github.com/park-jun-woo/huma/internal/source"
)

// StaticTodoPrompt builds the agent instruction for a TODO endpoint in static mode.
func StaticTodoPrompt(ep *scanner.Endpoint, hurlDir, urlVar string, branches []analyzer.ResponseBranch) string {
	var b strings.Builder

	b.WriteString(fmt.Sprintf("# TODO  %s %s\n", ep.Method, ep.Path))
	b.WriteString(fmt.Sprintf("# Source: %s:%d\n", ep.Source, ep.Line))
	b.WriteString(fmt.Sprintf("# Handler: %s\n", ep.Handler))

	src, _, _, err := source.ReadHandler(ep.Source, ep.Handler)
	if err == nil && src != "" {
		b.WriteString("\n## Handler source\n\n")
		b.WriteString(src)
		b.WriteString("\n")
	}

	if len(branches) > 0 {
		b.WriteString("\n## Expected responses (static analysis)\n\n")
		lines, statusList := collectBranchSection(branches)
		b.WriteString(lines)
		b.WriteString(formatResponseFields(ep.ResponseFields))
		b.WriteString("\n## Instructions\n\n")
		hurlFile := runner.HurlFileName(ep, hurlDir)
		b.WriteString(fmt.Sprintf("1. Write %s\n", hurlFile))
		b.WriteString(fmt.Sprintf("2. Include test entries for status %s\n", statusList))
		b.WriteString("3. Run `huma next`\n")
	} else {
		b.WriteString("\n## Hurl example\n\n")
		b.WriteString(hurlExample(ep.Method, ep.Path, urlVar))
		b.WriteString("\n")

		b.WriteString(formatResponseFields(ep.ResponseFields))

		hurlFile := runner.HurlFileName(ep, hurlDir)
		b.WriteString("\n## Instructions\n\n")
		b.WriteString(fmt.Sprintf("1. Write %s\n", hurlFile))
		b.WriteString("2. Run `huma next`\n")
	}

	return b.String()
}
