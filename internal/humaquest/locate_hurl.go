//ff:func feature=gate type=helper control=sequence level=error
//ff:what 제출물(.hurl)의 경로를 정해 파싱한다. raw가 비면 cfg.HurlDir+endpoint 관례 경로를, 아니면 raw(trim)를 경로로 본다. 파일 부재면 found=false(H-01), 그 외 파싱 에러는 err로. 부수효과 없음(텍스트 읽기만).

package humaquest

import (
	"os"
	"strings"

	"github.com/park-jun-woo/huma/internal/hurlcheck"
	"github.com/park-jun-woo/huma/internal/runner"
	"github.com/park-jun-woo/huma/internal/scanner"
)

// locateHurl resolves the .hurl submission path and parses it. The raw encoding is
// a FILE PATH (the agent writes the .hurl and submits its path):
//
//   - raw empty  → derive the conventional path runner.HurlFileName(ep, hurlDir).
//   - raw set    → treat the trimmed bytes as the path verbatim.
//
// It is G5-clean: only os.Stat + a buffered text read of the .hurl. A missing file
// is reported as found=false (the Prepare caller maps it to the H-01 verdict), not
// an error; err is reserved for genuine parse/read failures on a present file.
func locateHurl(raw []byte, ep scanner.Endpoint, hurlDir string) (path string, entries []hurlcheck.HurlEntry, found bool, err error) {
	path = strings.TrimSpace(string(raw))
	if path == "" {
		path = runner.HurlFileName(&ep, hurlDir)
	}

	if _, statErr := os.Stat(path); statErr != nil {
		// Absent (or unstattable) → H-01; let Prepare emit the FAIL verdict.
		return path, nil, false, nil
	}

	entries, err = hurlcheck.ParseHurlEntries(path)
	if err != nil {
		return path, nil, true, err
	}
	return path, entries, true, nil
}
