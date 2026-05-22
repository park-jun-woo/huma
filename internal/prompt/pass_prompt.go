//ff:func feature=prompt type=builder control=sequence
//ff:what Builds the PASS status line for a completed endpoint
package prompt

import (
	"fmt"

	"github.com/park-jun-woo/hurlfill/internal/scanner"
)

// PassPrompt builds the agent instruction for a PASS then next TODO.
func PassPrompt(ep *scanner.Endpoint) string {
	return fmt.Sprintf("# PASS  %s %s\n", ep.Method, ep.Path)
}
