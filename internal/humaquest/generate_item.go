//ff:func feature=gate type=engine control=sequence level=error
//ff:what generate 모드의 한 아이템 처리(C3). 면제면 SKIP 단락. 아니면 prompt 빌드(Render+lastReason)→backend.Complete→sanitizeHurl→stash(기존 보호)→정규 경로 기록→ParseHurlEntries(경로 기반) 검증(실패: 복원+parseFailVerdict, 측정 생략)→measureGenerated(측정·result.Pass==false면 강제 비-PASS)→payload 영속화→PASS면 생성물 확정·그 외 직전 hurl 복원. PASS 발급은 게이트뿐.

package humaquest

import (
	"github.com/park-jun-woo/huma/internal/config"
	"github.com/park-jun-woo/huma/internal/hurlcheck"
	"github.com/park-jun-woo/huma/internal/runner"
	"github.com/park-jun-woo/huma/internal/scanner"
	"github.com/park-jun-woo/reins/pkg/gate"
	"github.com/park-jun-woo/reins/pkg/llm"
	"github.com/park-jun-woo/reins/pkg/quest"
)

// generateItem is the --generate branch for one item. Exempt endpoints short-circuit
// to SKIP (generation cannot supply an oracle for a fully-exempt endpoint). Otherwise
// it builds the prompt (Render + last attempt's reason), calls the LLM backend,
// sanitizes the output, stashes any existing hurl, writes the generated content to the
// conventional path, then validates via ParseHurlEntries (path-based — write first).
// A parse failure restores the stash and returns parseFailVerdict (no measurement). On
// a parseable hurl it measures (measureGenerated forces a non-PASS on a runtime
// failure), persists the IMPROVE state, and keeps the generated file only on PASS —
// every non-PASS restores the prior hurl so no user asset is lost. PASS still
// originates only at the gate.
func generateItem(def gate.Definition, probe coverageProbe, s *quest.Session, it *quest.Item, ep scanner.Endpoint, ps payloadState, backend llm.Backend) (quest.Verdict, error) {
	if exempt, why := isExempt(ep); exempt {
		return quest.Verdict{Outcome: quest.OutSkip, Feedback: why}, nil
	}

	cfg, err := config.Load()
	if err != nil {
		cfg = &config.Config{HurlDir: "hurl"}
	}
	hurlPath := runner.HurlFileName(&ep, cfg.HurlDir)

	prompt, err := generatePrompt(def, s, it)
	if err != nil {
		return quest.Verdict{}, err
	}
	raw, err := backend.Complete(generateSystem, prompt)
	if err != nil {
		return quest.Verdict{}, err
	}

	stash, err := stashHurl(hurlPath)
	if err != nil {
		return quest.Verdict{}, err
	}
	if err := writeHurl(hurlPath, sanitizeHurl(raw)); err != nil {
		return quest.Verdict{}, err
	}
	if _, perr := hurlcheck.ParseHurlEntries(hurlPath); perr != nil {
		if rerr := restoreHurl(stash); rerr != nil {
			return quest.Verdict{}, rerr
		}
		return parseFailVerdict(ep, perr), nil
	}

	verdict, cov, err := measureGenerated(def, probe, s, it, ep)
	if err != nil {
		_ = restoreHurl(stash)
		return quest.Verdict{}, err
	}
	if err := persistCoverState(it, ep, cov, ps); err != nil {
		_ = restoreHurl(stash)
		return quest.Verdict{}, err
	}
	if verdict.Outcome != quest.OutPass {
		if rerr := restoreHurl(stash); rerr != nil {
			return quest.Verdict{}, rerr
		}
	}
	return verdict, nil
}
