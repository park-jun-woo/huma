//ff:func feature=verify type=helper control=sequence
//ff:what Reports whether this is the last retry (Tries+1 >= MaxTries) with no unreachable.yaml reason — the case where OutFail must be avoided so Apply does not auto-lock DONE.

package humaquest

import (
	"github.com/park-jun-woo/huma/internal/analyzer"
	"github.com/park-jun-woo/huma/internal/scanner"
	"github.com/park-jun-woo/reins/pkg/gate"
	"github.com/park-jun-woo/reins/pkg/quest"
)

// boundaryNoReason reports the MaxTries guard (§2): on the last attempt
// (it.Tries+1 >= quest.MaxTries) an OutFail would make quest.Apply unconditionally
// lock DONE. DONE is only legitimate as a reason-backed best-effort (C-04, §3.7),
// so this returns true when we are at the boundary AND the uncovered branches are
// NOT fully justified by unreachable.yaml — telling the caller to return OutReview
// (UNVERIFIED) instead of OutFail. Before the boundary it returns false (normal
// IMPROVE retry).
func boundaryNoReason(ctx gate.Context, ep *scanner.Endpoint, uncovered []analyzer.ResponseBranch) bool {
	lastTry := ctx.Item != nil && ctx.Item.Tries+1 >= quest.MaxTries
	return lastTry && !reasonsCover(ep, uncovered)
}
