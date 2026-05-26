//ff:func feature=prompt type=helper control=sequence
//ff:what Formats a single missing response branch as a display line with status code and location
package prompt

import (
	"fmt"
	"path/filepath"

	"github.com/park-jun-woo/huma/internal/analyzer"
)

func formatMissingBranch(m analyzer.ResponseBranch) string {
	base := filepath.Base(m.File)
	if m.Code != "" {
		return fmt.Sprintf("#   %d — %s:%d  %s\n", m.Status, base, m.Line, m.Code)
	}
	return fmt.Sprintf("#   %d — %s:%d\n", m.Status, base, m.Line)
}
