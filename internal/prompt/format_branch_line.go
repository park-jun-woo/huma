//ff:func feature=prompt type=helper control=sequence
//ff:what Formats a response branch as an indented display line with status code, file, and line
package prompt

import (
	"fmt"
	"path/filepath"

	"github.com/park-jun-woo/huma/internal/analyzer"
)

func formatBranchLine(br analyzer.ResponseBranch) string {
	base := filepath.Base(br.File)
	if br.Code != "" {
		return fmt.Sprintf("  %d — %s:%d  %s\n", br.Status, base, br.Line, br.Code)
	}
	return fmt.Sprintf("  %d — %s:%d\n", br.Status, base, br.Line)
}
