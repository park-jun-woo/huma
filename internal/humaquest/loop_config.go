//ff:type feature=gate type=data
//ff:what generate 모드(cover --generate)의 LLM 코칭 구성. generateSystem은 hurl 생성 System 프롬프트(L0 생성자에게 단일 .hurl만 출력하도록 지시). ruleSystem은 일반 정적 프리앰블 맵 — 판별 피드백이 아니라 보조 코칭이다. 게이트가 emit하는 RootCause는 C-01/C-03/C-04/H-01뿐이고 두 improve verdict가 모두 C-03을 쓰므로, 매 시도의 판별 피드백은 영속된 Attempt.Reason(lastReason)에서 오고 이 맵은 그 위에 얹는 정적 프리앰블일 뿐이다. H-03은 게이트 RootCause가 아니라(런타임 실패는 게이트가 emit 안 함) cover가 result.Pass==false를 감지했을 때 result.Feedback 앞에 붙이는 프리앰블로만 쓰인다.

package humaquest

// generateSystem is the System prompt handed to the LLM backend in --generate
// mode. It is the generation stage (L0) only — it has no authority over the
// gate's PASS lock; the CRI gate decides PASS independently after measurement.
const generateSystem = "You write a single Hurl (.hurl) file that verifies ONE HTTP endpoint. " +
	"Cover every client response branch (status + body shape), and assert body fields and " +
	"invariants — not just the HTTP status. Output ONLY a valid .hurl file: no prose, no " +
	"explanation, no markdown code fences."

// ruleSystem is the generic, static per-attempt coaching preamble keyed by a
// rule id. It is auxiliary, NOT the discriminating retry feedback: the gate only
// emits RootCause C-01/C-03/C-04/H-01, and BOTH improve verdicts set C-03, so a
// RootCause lookup cannot tell "uncovered branch" from "shallow assertion". The
// real per-attempt signal is the persisted Attempt.Reason (lastReason), which this
// preamble merely sits on top of. "H-03" is not a gate RootCause (runtime failures
// are never emitted by the gate); it is only the static preamble cover prepends to
// result.Feedback when it detects result.Pass==false.
var ruleSystem = map[string]string{
	"C-03": "Your hurl is too shallow or misses branches. Add [Asserts] for body shape and " +
		"invariants toward assertion depth A=3, and add requests that exercise the uncovered " +
		"response branches reported below.",
	"H-03": "Your hurl FAILED when it was run. Fix the exact assertion/status shown below and " +
		"output a full, corrected .hurl file.",
}
