//ff:func feature=prompt type=helper control=iteration dimension=1
//ff:what Formats response fields as a prompt section string with field paths and types
package prompt

import (
	"fmt"
	"strings"

	"github.com/park-jun-woo/huma/internal/scanner"
)

func formatResponseFields(fields []scanner.ResponseField) string {
	if len(fields) == 0 {
		return ""
	}

	var b strings.Builder
	b.WriteString("\n## Response fields\n\n")
	for _, f := range fields {
		if f.Type != "" {
			b.WriteString(fmt.Sprintf("  %s — %s\n", f.Path, f.Type))
		} else {
			b.WriteString(fmt.Sprintf("  %s\n", f.Path))
		}
	}
	return b.String()
}
