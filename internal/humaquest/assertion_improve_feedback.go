//ff:func feature=verify type=builder control=iteration dimension=1
//ff:what Builds the assertion-depth IMPROVE feedback (C-03): exposes the cap (A caps CRI to N, staged M) and, if response fields are known, lists a few as concrete A=3 body-shape assertion targets.

package humaquest

import (
	"fmt"
	"strings"

	"github.com/park-jun-woo/huma/internal/scanner"
)

// assertionImproveFeedback renders the human IMPROVE hint when A is the limiting
// axis (Phase 007, §3.3/§4). It first exposes the cap transparency — A caps CRI
// to the effective tier while the staged O/D/E ceiling is higher — then asks the
// agent to deepen assertions from status-only toward A=3 (status + body shape +
// invariants). When the endpoint carries OpenAPI/source response fields it lists
// a few of their paths as concrete body-shape assertion targets; otherwise it
// gives the generic shape/invariant hint. Kept simple per §3.3 (no schema
// codegen): naming up to three field paths is enough to unstick the agent.
func assertionImproveFeedback(ep *scanner.Endpoint, cri, effective int) string {
	base := fmt.Sprintf("IMPROVE: A=%d caps CRI to %d (staged %d) — status만 검증 중, body shape/불변식 assert 추가", effective, effective, cri)
	paths := make([]string, 0, 3)
	for _, f := range ep.ResponseFields {
		if f.Path == "" {
			continue
		}
		paths = append(paths, f.Path)
		if len(paths) == 3 {
			break
		}
	}
	if len(paths) == 0 {
		return base + "."
	}
	return fmt.Sprintf("%s (assert response body fields toward A=3: %s).", base, strings.Join(paths, ", "))
}
