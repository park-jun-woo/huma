//ff:func feature=gate type=helper control=selection level=error
//ff:what Definition.Prepare. Item payload의 Endpoint를 디코드하고, 면제 엔드포인트를 OutSkip으로 단락, .hurl을 위치확인/파싱(부재시 H-01 FAIL 단락)하고, 정적 검사(A-grade·H-04·H-05)를 hurlInfo에 실어 gate.Context를 만든다. 서버 미기동·hurl 미실행(G5-clean).

package humaquest

import (
	"fmt"

	"github.com/park-jun-woo/huma/internal/config"
	"github.com/park-jun-woo/huma/internal/hurlcheck"
	"github.com/park-jun-woo/huma/internal/scanner"
	"github.com/park-jun-woo/reins/pkg/gate"
	"github.com/park-jun-woo/reins/pkg/quest"
)

// Prepare decodes the raw submission into an evaluation Context. The raw bytes are
// the .hurl file PATH (empty → the conventional path is derived; see locateHurl).
// Prepare is G5-clean: it never starts a server or runs hurl — only text parsing
// and file reads. Server build/start, hurl execution, and coverage collection are
// isolated to Phase 005/006's Evaluate. It short-circuits two ways:
//   - exempt endpoint  → OutSkip (the gate is skipped → SKIPPED).
//   - missing .hurl    → OutFail/H-01 (a missing submission is a Prepare-stage
//     failure; fail-fast here keeps Evaluate from special-casing a nil submission
//     and mirrors bak/cmd/verify's H-01 exit).
func (humaDef) Prepare(s *quest.Session, it *quest.Item, raw []byte) (gate.Context, *quest.Verdict, error) {
	var ep scanner.Endpoint
	if err := it.DecodePayload(&ep); err != nil {
		return gate.Context{}, nil, err
	}

	// 1) Exempt endpoint → skip the gate entirely.
	if exempt, why := isExempt(ep); exempt {
		return gate.Context{}, &quest.Verdict{Outcome: quest.OutSkip, Feedback: why}, nil
	}

	cfg, err := config.Load()
	if err != nil {
		// No manifest → static-mode defaults, mirroring Render.
		cfg = &config.Config{HurlDir: "hurl"}
	}

	// 2) Locate + parse the .hurl submission (no execution).
	path, entries, found, err := locateHurl(raw, ep, cfg.HurlDir)
	if err != nil {
		return gate.Context{}, nil, err
	}
	if !found {
		detail := fmt.Sprintf("%s %s\n  expected: %s", ep.Method, ep.Path, path)
		return gate.Context{}, &quest.Verdict{
			Outcome:   quest.OutFail,
			RootCause: "H-01",
			Feedback:  "[H-01] hurl file not found at expected path\n  " + detail,
		}, nil
	}

	// 3) Static checks (no side effects): A-grade, naming (H-04), {{host}} (H-05).
	hi := &hurlInfo{
		Endpoint:  ep,
		HurlPath:  path,
		Entries:   entries,
		AGrade:    hurlcheck.MinAGrade(entries, hurlcheck.NonVacuousStatusList(entries)),
		NamingOK:  namingOK(path, ep, cfg.HurlDir),
		HostVarOK: hostVarOK(entries, cfg.URLVar()),
	}

	src := readHandlerSource(ep)
	return gate.Context{Item: it, Submission: hi, Source: src}, nil, nil
}
