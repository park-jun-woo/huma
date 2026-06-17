//ff:func feature=gate type=helper control=iteration dimension=1 level=error
//ff:what 엔드포인트 단락(short-circuit) 판정. .huma/unreachable.yaml의 유효 면제로 endpoint의 선언된 모든 응답 분기가 면제되면 검증할 게 없으므로 SKIP 사유를 반환한다. bak/cmd/all_exempt.go 이식. 부수효과 없음(파일 읽기만).

package humaquest

import (
	"fmt"

	"github.com/park-jun-woo/huma/internal/analyzer"
	"github.com/park-jun-woo/huma/internal/config"
	"github.com/park-jun-woo/huma/internal/scanner"
)

// isExempt ports bak/cmd/all_exempt.go: an endpoint is short-circuited (SKIPPED)
// when every one of its response branches carries a valid unreachable.yaml
// exemption — there is nothing left to verify. It reads the exemption artifact and
// the endpoint's declared (OpenAPI) response branches; both are pure disk/text
// reads, so isExempt is G5-clean.
//
// The session is intentionally not consulted: exemptions are an on-disk artifact
// (.huma/unreachable.yaml), not session state. The branch view here is the
// declared set; Phase 005's Evaluate does the authoritative source∪OpenAPI-union
// exemption/coverage judgment.
func isExempt(ep scanner.Endpoint) (bool, string) {
	branches := analyzer.ParseResponses(ep.Responses, ep.Source)
	if len(branches) == 0 {
		return false, ""
	}

	exemptions, err := config.LoadUnreachable()
	if err != nil || len(exemptions) == 0 {
		return false, ""
	}

	key := ep.Method + " " + ep.Path
	for _, b := range branches {
		if !config.IsExempt(exemptions, key, b.Status) {
			return false, ""
		}
	}
	return true, fmt.Sprintf("all %d response branches exempt via unreachable.yaml — nothing to verify for %s", len(branches), key)
}
