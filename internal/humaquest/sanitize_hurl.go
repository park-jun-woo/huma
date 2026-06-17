//ff:func feature=gate type=parser control=sequence level=error
//ff:what LLM 출력에서 순수 .hurl만 추출한다(tsma 003 C1/C3의 hurl판). 마크다운 펜스(```)가 있으면 첫 펜스 블록 안쪽을 꺼내 언어 태그 줄을 버리고, 없으면 트림한 원문을 그대로 돌려준다. 부수효과 없음.

package humaquest

import "strings"

// sanitizeHurl extracts pure .hurl text from an LLM completion, stripping markdown
// code fences and surrounding prose. When the output contains a ``` fence it returns
// the inside of the first fenced block (dropping an optional language tag on the
// fence line); otherwise it returns the trimmed raw text. It is pure (no I/O).
func sanitizeHurl(raw string) string {
	s := strings.TrimSpace(raw)
	start := strings.Index(s, "```")
	if start < 0 {
		return s
	}
	rest := s[start+3:]
	if nl := strings.IndexByte(rest, '\n'); nl >= 0 {
		rest = rest[nl+1:]
	}
	if end := strings.Index(rest, "```"); end >= 0 {
		rest = rest[:end]
	}
	return strings.TrimSpace(rest)
}
