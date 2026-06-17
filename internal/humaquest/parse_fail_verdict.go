//ff:func feature=verify type=engine control=sequence level=error
//ff:what 생성된 hurl이 ParseHurlEntries(경로 기반)로 파싱되지 않을 때(잘림/구문오류)의 전용 비-PASS verdict. 측정 생략·직전 hurl 복원 후 호출된다(§2-4). OutFail(재시도) + "전체 .hurl을 한 블록으로 출력하라" 피드백; 파싱 에러는 Facts.Actual에 실어 다음 시도의 lastReason으로 흐른다.

package humaquest

import (
	"github.com/park-jun-woo/huma/internal/scanner"
	"github.com/park-jun-woo/reins/pkg/quest"
)

// parseFailVerdict builds the dedicated non-PASS verdict for a generated hurl that
// ParseHurlEntries (path-based) could not parse — truncated or syntactically broken
// output. Measurement is skipped and the prior hurl restored before this is returned
// (§2-4). It maps to OutFail (retry); the parse error rides in the Fact's Actual so
// quest.Apply persists it into Attempt.Reason for the next generation's lastReason.
func parseFailVerdict(ep scanner.Endpoint, perr error) quest.Verdict {
	key := ep.Method + " " + ep.Path
	return quest.Verdict{
		Outcome:   quest.OutFail,
		RootCause: "H-03",
		Facts: []quest.Fact{{
			Rule:     "H-03",
			Where:    key,
			Expected: "a complete, parseable .hurl file",
			Actual:   "output appears truncated/invalid hurl: " + perr.Error(),
		}},
		Feedback: "output appears truncated/invalid hurl — emit the ENTIRE .hurl in one block, no prose, no markdown fences.",
	}
}
