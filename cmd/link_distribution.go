//ff:func feature=scan type=helper control=iteration dimension=1
//ff:what Formats the matched-file extension breakdown for the link-source summary
package cmd

import (
	"fmt"
	"sort"
	"strings"

	"github.com/park-jun-woo/huma/internal/scanner"
)

// distribution formats the matched-file extension breakdown, e.g. "go: 142",
// noting the low-confidence fallback when the backend lang was unknown.
func distribution(r scanner.LinkResult) string {
	if len(r.ByExt) == 0 {
		if !r.LangKnown {
			return "lang=unknown, low-confidence"
		}
		return "none"
	}
	exts := make([]string, 0, len(r.ByExt))
	for e := range r.ByExt {
		exts = append(exts, e)
	}
	sort.Strings(exts)
	parts := make([]string, 0, len(exts))
	for _, e := range exts {
		parts = append(parts, fmt.Sprintf("%s: %d", strings.TrimPrefix(e, "."), r.ByExt[e]))
	}
	out := strings.Join(parts, ", ")
	if !r.LangKnown {
		out += " (lang=unknown, low-confidence)"
	}
	return out
}
