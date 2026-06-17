//ff:func feature=gate type=helper control=sequence level=error
//ff:what Reads the previous attempt's line-coverage percent from the Item payload (read-only) for the IMPROVE monotonicity check.

package humaquest

import "github.com/park-jun-woo/reins/pkg/quest"

// priorCoverage decodes the Item payload's PrevCoverage (written by Phase 006's
// cover command on the previous attempt) so Evaluate can compare current vs prior
// coverage for the stalled-vs-improving note (§2). It is strictly read-only —
// Evaluate never writes Payload. A missing/old bare-Endpoint payload decodes to 0
// (no prior signal), which is the correct neutral default.
func priorCoverage(it *quest.Item) float64 {
	if it == nil {
		return 0
	}
	var ps payloadState
	if err := it.DecodePayload(&ps); err != nil {
		return 0
	}
	return ps.PrevCoverage
}
