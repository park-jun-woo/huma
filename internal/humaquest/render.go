//ff:func feature=gate type=helper control=selection
//ff:what Definition.Render. Item.Payload의 Endpoint를 디코드하고 derivePhase로 Item 상태(Tries/Log)+manifest 모드를 읽어 todo/improve/unverified/static 프롬프트를 고른다. read-only — it·s.Meta를 변경하지 않는다(상태 전이는 quest.Apply 담당).

package humaquest

import (
	"github.com/park-jun-woo/huma/internal/config"
	"github.com/park-jun-woo/huma/internal/prompt"
	"github.com/park-jun-woo/huma/internal/runner"
	"github.com/park-jun-woo/huma/internal/scanner"
	"github.com/park-jun-woo/reins/pkg/quest"
)

// Render returns the authoring prompt + verification context shown by `next`.
// It is read-only: it never mutates it or s.Meta. The prompt variant is
// re-derived from the Item's own signals (Tries/Log) plus the manifest mode —
// the old huma session Status no longer exists. Missing source/manifest is
// handled gracefully (the prompt builders skip the source section, and config
// falls back to static-mode defaults).
func (humaDef) Render(s *quest.Session, it *quest.Item) (string, error) {
	var ep scanner.Endpoint
	if err := it.DecodePayload(&ep); err != nil {
		return "", err
	}

	cfg, err := config.Load()
	if err != nil {
		// No manifest → treat as static mode with default hurl dir.
		cfg = &config.Config{HurlDir: "hurl"}
	}
	staticMode := cfg.Server.Start == ""

	switch derivePhase(it, staticMode) {
	case phaseImprove:
		hurlFile := runner.HurlFileName(&ep, cfg.HurlDir)
		return prompt.ImprovePrompt(&ep, hurlFile, lastReason(it)), nil
	case phaseUnverified:
		return prompt.UnverifiedPrompt(&ep, cfg), nil
	case phaseStatic:
		return prompt.StaticTodoPrompt(&ep, cfg.HurlDir, cfg.URLVar(), nil), nil
	default: // phaseTodo
		return prompt.TodoPrompt(&ep, cfg.HurlDir, cfg.URLVar()), nil
	}
}
