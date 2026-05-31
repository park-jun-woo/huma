//ff:func feature=prompt type=builder control=sequence
//ff:what Builds the UNVERIFIED instruction prompting the user to supply a source link or server instrumentation
package prompt

import (
	"fmt"
	"strings"

	"github.com/park-jun-woo/huma/internal/config"
	"github.com/park-jun-woo/huma/internal/scanner"
)

// UnverifiedPrompt explains why an endpoint cannot be granted PASS and presents
// the forced choice to obtain an independent oracle (§3.2, §5).
func UnverifiedPrompt(ep *scanner.Endpoint, cfg *config.Config) string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("# UNVERIFIED  %s %s  CRI 0\n", ep.Method, ep.Path))
	if ep.Source == "" {
		b.WriteString("# source unlinked")
	} else {
		b.WriteString(fmt.Sprintf("# source: %s", ep.Source))
	}
	if cfg != nil && cfg.Server.Start == "" {
		b.WriteString(" AND no server instrumentation (static mode)\n")
	} else {
		b.WriteString(" AND runtime uninstrumented (Total==0)\n")
	}
	b.WriteString("# No independent oracle exists — measurement failed, so this is NOT a pass.\n")
	b.WriteString("\n## Fix (pick one)\n\n")
	b.WriteString("1. Link source: `huma scan --from <openapi> --link-source <root>`\n")
	b.WriteString("2. Instrument the server: set testing.server + a coverage-enabled build\n")
	return b.String()
}
