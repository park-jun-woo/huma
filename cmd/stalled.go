//ff:func feature=verify type=helper control=sequence
//ff:what Reports whether coverage improvement has stalled after at least one retry
package cmd

import (
	"github.com/park-jun-woo/huma/internal/adapter"
	"github.com/park-jun-woo/huma/internal/session"
)

// stalled reports whether the endpoint has been retried at least once and the
// coverage did not improve over the previous attempt.
func stalled(entry *session.Entry, cov *adapter.CoverageResult) bool {
	return entry != nil && entry.ImproveCount >= 1 && cov.Percent <= entry.PrevCoverage
}
