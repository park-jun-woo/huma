//ff:func feature=gate type=helper control=iteration dimension=1
//ff:what cover 루프의 per-item verdict를 submit과 같은 형식으로 출력한다: "key -> OUTCOME (state STATE)" 한 줄에 더해, Feedback(공략집)이 있으면 그대로, 없으면 Facts(규칙·위치·기대·실제)를 들여쓴다. reins cli.renderVerdict가 unexported라 동등 출력을 둔다.

package humaquest

import (
	"fmt"
	"io"
	"strings"

	"github.com/park-jun-woo/reins/pkg/quest"
)

// renderCoverVerdict prints one item's outcome the way reins' submit does (its
// renderVerdict is unexported): the verdict line, then the pre-rendered Feedback
// walkthrough when present, else the per-Fact loop.
func renderCoverVerdict(w io.Writer, key string, it *quest.Item, v quest.Verdict) {
	fmt.Fprintf(w, "%s -> %s (state %s)\n", key, v.Outcome, it.State)
	if v.Feedback != "" {
		for _, line := range strings.Split(strings.TrimRight(v.Feedback, "\n"), "\n") {
			fmt.Fprintf(w, "  %s\n", line)
		}
		return
	}
	for _, f := range v.Facts {
		fmt.Fprintf(w, "  - %s: %s expected=%q actual=%q\n", f.Rule, f.Where, f.Expected, f.Actual)
	}
}
