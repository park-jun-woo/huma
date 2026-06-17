//ff:func feature=gate type=data control=iteration dimension=1
//ff:what Definition.Rules. huma의 M/E/H/S/A/C 룰북(bak/internal/rule/)을 패키지 테이블 ruleCatalog에 선언하고, gate.RuleMeta 카탈로그로 펼쳐 반환한다(`rules` 명령·자동 rulebook 전용). ERROR→LevelFail, WARNING→LevelReview. Check는 카탈로그가 가리킬 수 있도록 시그니처만 맞춘 무발동 스텁이며, 실제 판정은 Phase 005의 Evaluator가 한다.

package humaquest

import (
	"github.com/park-jun-woo/reins/pkg/gate"
	"github.com/park-jun-woo/reins/pkg/quest"
)

// ruleCatalog is huma's violation rulebook (M/E/H/S/A/C), ported verbatim from
// bak/internal/rule/ and kept 1:1 with rulebook.md. huma's ERROR maps to
// gate.LevelFail and WARNING to gate.LevelReview.
var ruleCatalog = []struct {
	id    string
	level gate.Level
	desc  string
}{
	// M. Manifest validation
	{"M-01", gate.LevelFail, "manifest.yaml not found"},
	{"M-02", gate.LevelFail, "manifest.yaml parse error (invalid YAML)"},
	{"M-03", gate.LevelFail, "apiVersion missing or unsupported"},
	{"M-04", gate.LevelFail, "metadata.name missing"},
	{"M-05", gate.LevelFail, "backend.lang missing"},
	{"M-06", gate.LevelFail, "testing.base_url missing"},
	{"M-07", gate.LevelFail, "testing.hurl_dir missing"},
	{"M-08", gate.LevelFail, "testing.server.start missing"},
	{"M-09", gate.LevelFail, "testing.server.ready missing"},
	{"M-10", gate.LevelReview, "testing.hurl_variables empty"},

	// E. Endpoint input validation
	{"E-01", gate.LevelFail, "No OpenAPI file found and --from not specified"},
	{"E-02", gate.LevelFail, "Input file not readable"},
	{"E-03", gate.LevelFail, "Input is not valid JSON/YAML"},
	{"E-04", gate.LevelFail, "Endpoint missing method field"},
	{"E-05", gate.LevelFail, "Endpoint missing path field"},
	{"E-06", gate.LevelReview, "Endpoint missing handler field"},
	{"E-07", gate.LevelReview, "Endpoint missing file field"},
	{"E-08", gate.LevelReview, "Duplicate endpoint"},
	{"E-09", gate.LevelReview, "OpenAPI auto-detect failed, falling back to endpoint list parser"},

	// H. Hurl file validation
	{"H-01", gate.LevelFail, "Hurl file not found at expected path"},
	{"H-02", gate.LevelFail, "Hurl execution failed"},
	{"H-03", gate.LevelFail, "Hurl test failed"},
	{"H-04", gate.LevelReview, "Existing hurl file name doesn't match naming convention"},
	{"H-05", gate.LevelReview, "Hurl file missing {{host}} variable"},

	// S. Session state validation
	{"S-01", gate.LevelFail, "No session found"},
	{"S-02", gate.LevelFail, "Session file corrupt"},
	{"S-03", gate.LevelReview, "Session has stale entries"},

	// A. Adapter / server validation
	{"A-01", gate.LevelFail, "Server healthcheck failed"},
	{"A-02", gate.LevelFail, "Server build command failed"},
	{"A-03", gate.LevelFail, "Server start command failed"},
	{"A-04", gate.LevelFail, "Server ready timeout"},
	{"A-05", gate.LevelFail, "Coverage data collection failed"},
	{"A-06", gate.LevelReview, "deps.ready check failed"},

	// C. Coverage verdict (cheese-resistant gate)
	{"C-01", gate.LevelFail, "No-signal verdict cannot PASS — downgraded to UNVERIFIED"},
	{"C-02", gate.LevelFail, "Denominator is monotonic — input spec cannot shrink ground-truth branches"},
	{"C-03", gate.LevelReview, "Assertion depth below required level"},
	{"C-04", gate.LevelReview, "DONE requires an unreachable.yaml reason artifact for uncovered branches"},
}

// Rules returns huma's violation-rule catalog as gate.Rules. It is the audit
// rulebook surfaced by the `rules` command, not the verdict path: the real CRI
// judging is the Phase 005 Evaluator. Each Check is a no-fire stub that satisfies
// the gate.Rule signature so the catalog can be listed.
func (humaDef) Rules() []gate.Rule {
	out := make([]gate.Rule, len(ruleCatalog))
	for i, r := range ruleCatalog {
		out[i] = gate.Rule{
			Meta:  gate.RuleMeta{ID: r.id, Level: r.level, Desc: r.desc},
			Check: func(gate.Context) (bool, quest.Fact) { return false, quest.Fact{} },
		}
	}
	return out
}
