//ff:func feature=verify type=helper control=sequence
//ff:what Builds the PASS verdict (OutPass) at an achieved CRI tier, with the §5 transparency line (tier, branch count, A-grade, provenance) as feedback.

package humaquest

import (
	"fmt"

	"github.com/park-jun-woo/huma/internal/analyzer"
	"github.com/park-jun-woo/reins/pkg/quest"
)

// passVerdict builds the PASS verdict for an achieved CRI tier (SCAFFOLDED/
// SMOKE/COVERED). It maps to quest.OutPass → PASS (the only outcome quest.Apply
// locks to PASS). It carries no Facts (PASS has nothing to correct); the §5
// transparency line — CRI tier, client-branch count, A-grade, denominator
// provenance — rides in Feedback so submit can print it. Ports the bak
// sess.SetVerdict + sess.MarkPass.
func passVerdict(tier, aGrade int, prov string, branches []analyzer.ResponseBranch) quest.Verdict {
	return quest.Verdict{
		Outcome: quest.OutPass,
		Feedback: fmt.Sprintf("%s  CRI %d  (client branches %d, A=%d, source: %s)",
			criLabel(tier), tier, len(branches), aGrade, prov),
	}
}
