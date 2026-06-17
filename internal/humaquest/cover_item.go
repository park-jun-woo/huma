//ff:func feature=gate type=engine control=sequence level=error
//ff:what 한 TODO 아이템의 측정→판정을 수행한다: payload에서 Endpoint/prior 디코드→probe.Reset(liveProbe에선 no-op)→Measure(엔드포인트마다 Start→run hurl→Stop(flush)→Collect)→def.Prepare(short-circuit 존중)→측정됐으면 ctx.Grounds["coverage"] 주입→Evaluate→(Evaluate 직후, Apply 전에) payloadState(PrevCoverage/ImproveCount) 영속화. verdict만 돌려주고 래칫/Save/Export는 호출자(runCover)가 한다. PASS 발급은 여전히 게이트 Evaluate에만 있다.

package humaquest

import (
	"github.com/park-jun-woo/reins/pkg/gate"
	"github.com/park-jun-woo/reins/pkg/llm"
	"github.com/park-jun-woo/reins/pkg/quest"
)

// coverItem measures and evaluates one item, replicating reins' unexported
// evaluateAndApply head (Prepare → short|Evaluate) with one huma addition: it injects
// the measured coverage as ctx.Grounds["coverage"] before Evaluate and persists the
// IMPROVE monotonicity (PrevCoverage/ImproveCount) into the Item payload right after
// Evaluate — Evaluate is read-only, so cover owns this write. When backend != nil
// (--generate) it delegates to generateItem, which LLM-generates the .hurl, writes it
// to disk, then measures/evaluates; a nil backend is the manual path (behavior
// unchanged). It returns the verdict; the caller applies the ratchet, saves, and
// exports (the PASS-lock tail). PASS still originates only at the gate Evaluate.
func coverItem(def gate.Definition, probe coverageProbe, s *quest.Session, it *quest.Item, backend llm.Backend) (quest.Verdict, error) {
	var ps payloadState
	if err := it.DecodePayload(&ps); err != nil {
		return quest.Verdict{}, err
	}
	ep := ps.Endpoint

	if backend != nil {
		return generateItem(def, probe, s, it, ep, ps, backend)
	}

	if err := probe.Reset(); err != nil {
		return quest.Verdict{}, err
	}
	cov, _, err := probe.Measure(ep)
	if err != nil {
		return quest.Verdict{}, err
	}

	// Build the Context via the SAME Prepare the submit path uses. Pass an empty raw
	// so locateHurl derives the conventional .hurl path (matching Measure's lookup).
	ctx, short, err := def.Prepare(s, it, nil)
	if err != nil {
		return quest.Verdict{}, err
	}
	if short != nil {
		// Exempt (SKIPPED) or missing-hurl (H-01 FAIL): honor it like submit does.
		return *short, nil
	}

	// Inject the measured coverage as a ground so Evaluate yields a live CRI.
	if cov != nil {
		ground, err := coverageGround(cov)
		if err != nil {
			return quest.Verdict{}, err
		}
		ctx.Grounds = map[string]string{"coverage": ground}
	}

	verdict := evaluate(def, ctx)

	// Persist IMPROVE monotonicity AFTER Evaluate, BEFORE the caller's Apply/Save.
	if err := persistCoverState(it, ep, cov, ps); err != nil {
		return quest.Verdict{}, err
	}
	return verdict, nil
}
