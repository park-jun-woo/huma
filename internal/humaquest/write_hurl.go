//ff:func feature=gate type=helper control=sequence level=error
//ff:what 생성된 hurl 내용을 정규 경로에 기록한다(부모 디렉터리 생성 포함). 측정 전 LLM 출력이 디스크에 닿게 하는 단계 — runner.FindHurlFile/ParseHurlEntries가 경로 기반이라 기록이 선행되어야 한다.

package humaquest

import (
	"os"
	"path/filepath"
)

// writeHurl writes the sanitized generated content to the conventional .hurl path,
// creating the hurl directory if needed. This is the step that makes the LLM's
// output reach disk before measurement: Measure (FindHurlFile) and ParseHurlEntries
// are both path-based, so the file must exist first.
func writeHurl(path, content string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, []byte(content), 0o644)
}
