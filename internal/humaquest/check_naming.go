//ff:func feature=gate type=helper control=sequence
//ff:what H-04 정적 검사. 해석된 .hurl 파일명이 endpoint의 네이밍 관례(runner.HurlFileName)와 일치하는지(basename 비교) 보고한다. 부수효과 없음.

package humaquest

import (
	"path/filepath"

	"github.com/park-jun-woo/huma/internal/runner"
	"github.com/park-jun-woo/huma/internal/scanner"
)

// namingOK reports H-04: whether the resolved .hurl file name matches the naming
// convention for the endpoint. Compared by base name so a conventional file under
// any directory still passes. When the path was derived (raw empty) this is always
// true; it can only fail when the agent submits a hand-named path.
func namingOK(path string, ep scanner.Endpoint, hurlDir string) bool {
	return filepath.Base(path) == filepath.Base(runner.HurlFileName(&ep, hurlDir))
}
