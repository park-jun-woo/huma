//ff:func feature=verify type=helper control=sequence
//ff:what Reports whether a set of uncovered branches is fully justified by valid unreachable.yaml reasons (the DONE gate, §3.7/C-04).

package humaquest

import (
	"github.com/park-jun-woo/huma/internal/analyzer"
	"github.com/park-jun-woo/huma/internal/config"
	"github.com/park-jun-woo/huma/internal/scanner"
)

// reasonsCover reports whether every uncovered branch carries a valid
// unreachable.yaml reason — the precondition for granting DONE (§3.7, C-04). An
// empty uncovered set is trivially satisfied. Factored out of
// bak/cmd/done_reasons_satisfied.go so both the live and static improve paths can
// decide the MaxTries→DONE boundary over their own uncovered set.
func reasonsCover(ep *scanner.Endpoint, uncovered []analyzer.ResponseBranch) bool {
	if len(uncovered) == 0 {
		return true
	}
	exemptions, err := config.LoadUnreachable()
	if err != nil || len(exemptions) == 0 {
		return false
	}
	return allExempt(uncovered, exemptions, ep.Method+" "+ep.Path)
}
