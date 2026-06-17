//ff:func feature=gate type=helper control=iteration dimension=1
//ff:what H-05 정적 검사. 모든 non-skip 엔트리의 URL이 {{<urlVar>}} 호스트 템플릿 변수를 참조하는지(하드코딩 base URL 금지) 보고한다. 부수효과 없음.

package humaquest

import (
	"strings"

	"github.com/park-jun-woo/huma/internal/hurlcheck"
)

// hostVarOK reports H-05: every non-skip entry with a URL must reference the
// {{<urlVar>}} host template variable rather than a hardcoded base URL. Returns
// false for an empty entry set (a vacuous file cannot confirm the convention).
func hostVarOK(entries []hurlcheck.HurlEntry, urlVar string) bool {
	if len(entries) == 0 {
		return false
	}
	token := "{{" + urlVar + "}}"
	for _, e := range entries {
		if e.Skip || e.URL == "" {
			continue
		}
		if !strings.Contains(e.URL, token) {
			return false
		}
	}
	return true
}
