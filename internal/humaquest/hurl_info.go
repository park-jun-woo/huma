//ff:type feature=gate type=model
//ff:what Prepare가 만들어 gate.Context.Submission에 싣는 huma 제출물. .hurl 위치/파싱 결과와 정적 검사(A-grade·네이밍 H-04·{{host}} H-05) 산출을 담아 Phase005 Evaluate가 재파싱 없이 소비한다.

package humaquest

import (
	"github.com/park-jun-woo/huma/internal/hurlcheck"
	"github.com/park-jun-woo/huma/internal/scanner"
)

// hurlInfo is huma's decoded submission. Prepare resolves and statically inspects
// the .hurl file (no server, no hurl execution) and stashes everything the static
// stage produced here so Phase 005's Evaluate consumes it via
// ctx.Submission.(*hurlInfo) without re-parsing the file.
//
// Fields:
//   - Endpoint: the endpoint under test (decoded from the Item payload).
//   - HurlPath: the resolved .hurl path (see locateHurl's path rule).
//   - Entries:  parsed request blocks with measured A-grade — the cached parse
//     Evaluate reuses (e.g. MinAGrade against the real branch denominator,
//     NonVacuousStatuses for coverage).
//   - AGrade:   minimum assertion-depth grade across the hurl's own non-vacuous
//     statuses (a self-contained static signal; Evaluate may recompute against the
//     full source∪OpenAPI branch union).
//   - NamingOK: H-04 — the resolved file name matches the naming convention.
//   - HostVarOK: H-05 — every non-skip entry URL references the {{<urlVar>}} host
//     template variable.
type hurlInfo struct {
	Endpoint  scanner.Endpoint
	HurlPath  string
	Entries   []hurlcheck.HurlEntry
	AGrade    int
	NamingOK  bool
	HostVarOK bool
}
