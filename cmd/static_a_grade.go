//ff:func feature=verify type=helper control=sequence
//ff:what Measures the minimum assertion-depth grade across client branch statuses for a hurl file
package cmd

import (
	"github.com/park-jun-woo/huma/internal/analyzer"
	"github.com/park-jun-woo/huma/internal/hurlcheck"
)

// staticAGrade measures the minimum assertion depth across the client branch
// statuses for the transparency output.
func staticAGrade(hurl string, branches []analyzer.ResponseBranch) int {
	entries, err := hurlcheck.ParseHurlEntries(hurl)
	if err != nil {
		return 0
	}
	return hurlcheck.MinAGrade(entries, branchStatuses(branches))
}
