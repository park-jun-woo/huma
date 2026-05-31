//ff:func feature=session type=builder control=selection
//ff:what Renders a transparent per-endpoint status line with CRI label and denominator provenance
package cmd

import (
	"fmt"

	"github.com/park-jun-woo/huma/internal/session"
)

// statusLine renders one endpoint's transparent status line with its CRI label
// and denominator provenance (§5).
func statusLine(e session.Entry) string {
	prov := provenanceLabel(e.Provenance)
	switch e.Status {
	case session.StatusUnverified:
		return fmt.Sprintf("UNVERIFIED %-6s %-32s CRI 0  (source: %s)\n             → fix: --link-source OR testing.server", e.Method, e.Path, prov)
	case session.StatusPass, session.StatusDone:
		return fmt.Sprintf("%-10s %-6s %-32s CRI %d  (source: %s, A=%d)", session.CRILabel(e.CRI), e.Method, e.Path, e.CRI, prov, e.AGrade)
	case session.StatusImprove:
		return fmt.Sprintf("IMPROVE    %-6s %-32s (%.0f%%, source: %s)", e.Method, e.Path, e.Coverage, prov)
	default:
		return fmt.Sprintf("TODO       %-6s %-32s", e.Method, e.Path)
	}
}
