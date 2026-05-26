//ff:func feature=prompt type=helper control=iteration dimension=1
//ff:what Collects formatted branch lines and status labels from response branches
package prompt

import (
	"fmt"
	"strings"

	"github.com/park-jun-woo/huma/internal/analyzer"
)

func collectBranchSection(branches []analyzer.ResponseBranch) (lines string, statusList string) {
	var lb strings.Builder
	statuses := make([]string, len(branches))
	for i, br := range branches {
		lb.WriteString(formatBranchLine(br))
		statuses[i] = fmt.Sprintf("%d", br.Status)
	}
	return lb.String(), strings.Join(statuses, ", ")
}
