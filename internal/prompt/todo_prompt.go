//ff:func feature=prompt type=builder control=sequence
//ff:what Builds the agent instruction text for a TODO endpoint
package prompt

import (
	"fmt"
	"strings"

	"github.com/park-jun-woo/huma/internal/runner"
	"github.com/park-jun-woo/huma/internal/scanner"
	"github.com/park-jun-woo/huma/internal/source"
)

// TodoPrompt builds the agent instruction for a TODO endpoint.
func TodoPrompt(ep *scanner.Endpoint, hurlDir, urlVar string) string {
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

	b.WriteString("\n## Hurl example\n\n")
	b.WriteString(hurlExample(ep.Method, ep.Path, urlVar))
	b.WriteString("\n")

	hurlFile := runner.HurlFileName(ep, hurlDir)
	b.WriteString("\n## Instructions\n\n")
	b.WriteString(fmt.Sprintf("1. Write %s\n", hurlFile))
	b.WriteString("2. Run `huma next`\n")

	return b.String()
}
