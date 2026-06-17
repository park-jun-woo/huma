//ff:func feature=verify type=helper control=sequence
//ff:what Measures the minimum assertion-depth grade across the client branch statuses using the cached hurl entries.

package humaquest

import (
	"github.com/park-jun-woo/huma/internal/analyzer"
	"github.com/park-jun-woo/huma/internal/hurlcheck"
)

// staticAGrade measures the minimum assertion depth across the client branch
// statuses for the transparency output. It reuses the entries Prepare already
// parsed (Phase 004) rather than re-reading the .hurl file — the only change
// from bak/cmd/static_a_grade.go, which took the file path and re-parsed it.
func staticAGrade(entries []hurlcheck.HurlEntry, branches []analyzer.ResponseBranch) int {
	return hurlcheck.MinAGrade(entries, branchStatuses(branches))
}
