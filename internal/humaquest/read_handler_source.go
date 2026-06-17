//ff:func feature=gate type=helper control=sequence
//ff:what endpoint 핸들러의 원천 소스를 읽어 gate.Context.Source(치즈방어 규칙이 재확인하는 캐시 원천)에 싣는다. 소스 파일/핸들러 부재는 빈 문자열로 graceful. 부수효과 없음(파일 읽기만).

package humaquest

import (
	"github.com/park-jun-woo/huma/internal/scanner"
	"github.com/park-jun-woo/huma/internal/source"
)

// readHandlerSource returns the endpoint's handler body for gate.Context.Source —
// the cached ground truth cheese-defense rules re-confirm. It is best-effort: a
// missing source file or unlocatable handler yields "" rather than an error, since
// Source is advisory context, not a gate input. It is G5-clean (a single file read).
func readHandlerSource(ep scanner.Endpoint) string {
	if ep.Source == "" || ep.Handler == "" {
		return ""
	}
	src, _, _, err := source.ReadHandler(ep.Source, ep.Handler)
	if err != nil {
		return ""
	}
	return src
}
