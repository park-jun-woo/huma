//ff:func feature=gate type=helper control=sequence level=error
//ff:what cover 명령의 --generate/--model에서 LLM backend를 만든다. generate가 아니면 nil(수동 모드 — 동작 불변), 맞으면 llm.FromFlag(model)로 backend:model 문자열을 Backend로 해석한다. backend!=nil이 곧 generate 모드 신호다.

package humaquest

import "github.com/park-jun-woo/reins/pkg/llm"

// buildBackend turns the cover command's --generate/--model flags into an llm.Backend.
// Without --generate it returns nil (manual mode — cover/submit/next behavior is
// unchanged), so a nil backend IS the generate-mode-off signal threaded through
// runCover/coverItem. With --generate it resolves the "backend:model" string via
// llm.FromFlag (default claude:sonnet).
func buildBackend(generate bool, model string) (llm.Backend, error) {
	if !generate {
		return nil, nil
	}
	return llm.FromFlag(model)
}
