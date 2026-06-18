//ff:func feature=gate type=command control=sequence level=error
//ff:what `huma loop` ExtraCommand를 만든다. 루트의 persistent 플래그(--session/--out)를 런타임에 상속받아 읽고, manifest를 로드해 실서버 probe(newLiveProbe)를 구성한 뒤 runCover에 넘긴다. Phase 010: 라이브 생성 루프를 사용자 표면 명령 loop으로 노출 — 기본 생성 ON(`--measure-only`로 측정 전용), 기본 모델 ollama:gemma4:e4b. Phase 009: newLiveProbe 직후 setupVars로 testing.setup(캡처) 또는 testing.auth(mint)를 1회 실행해 token·픽스처를 lp.extraVars에 적재한다(Measure가 정적 변수 위에 합류). 생성 모드면 buildBackend(--model)로 LLM backend를 만들어 무인 생성 루프를 켜고, --max-items로 처리할 distinct endpoint 수를 캡한다. 서버 라이프사이클을 소유하는 huma 전용 명령 — 계측 바이너리를 한 번 컴파일하고, 엔드포인트마다 Start/Stop/Collect로 측정하며 전 TODO를 순회한다(Go 커버리지는 프로세스 종료 시 flush).

package humaquest

import (
	"fmt"

	"github.com/park-jun-woo/huma/internal/config"
	"github.com/park-jun-woo/reins/pkg/gate"
	"github.com/spf13/cobra"
)

// NewLoopCmd builds the `loop` ExtraCommand: huma's server-owning, LLM-generating
// convergence loop. It is attached via cli.Options.ExtraCommands and reuses the root's
// persistent --session and --out flags (read at run time from the inherited flag set).
// It loads the manifest, constructs the real liveProbe (the default coverage source),
// and hands orchestration to runCover. The probe seam means the real adapter is the
// default here while runCover stays unit-testable with a fake.
//
// Phase 010 makes `loop` the single live-loop surface: generation is ON by default (the
// loop's reason to exist), `--measure-only` reverts to plain coverage measurement of
// existing .hurl, and the default model is the local ollama:gemma4:e4b.
func NewLoopCmd(def gate.Definition) *cobra.Command {
	var measureOnly bool
	var model string
	var maxItems int
	cmd := &cobra.Command{
		Use:   "loop",
		Short: "generate .hurl with an LLM, then measure live coverage per endpoint and converge via the CRI gate",
		RunE: func(cmd *cobra.Command, args []string) error {
			sessionPath, err := cmd.Flags().GetString("session")
			if err != nil {
				return err
			}
			outPath, err := cmd.Flags().GetString("out")
			if err != nil {
				return err
			}
			cfg, err := config.Load()
			if err != nil {
				return fmt.Errorf("loop needs a manifest with testing.server: %w", err)
			}
			// Generation is on by default; --measure-only reverts to measuring
			// existing .hurl with no LLM.
			gen := !measureOnly
			// Phase 010 §4 C1: one-line startup log (target session / mode / model)
			// so the unattended generate loop announces its disk-writing intent.
			mode := "measure-only"
			if gen {
				mode = "generate"
			}
			fmt.Fprintf(cmd.OutOrStdout(), "loop: session=%s mode=%s model=%s\n", sessionPath, mode, model)
			backend, err := buildBackend(gen, model)
			if err != nil {
				return err
			}
			// Phase 009: resolve the dynamic test variables (token/fixtures) once,
			// before the loop, and load them into the probe's extraVars. Measure
			// merges them over cfg.HurlVariables, so both manual and generate
			// endpoint runs get {{token}} injected.
			lp := newLiveProbe(cfg)
			lp.extraVars = setupVars(lp.adapter, cfg, cmd.OutOrStdout())
			return runCover(def, lp, backend, maxItems, sessionPath, outPath, cmd.OutOrStdout())
		},
	}
	cmd.Flags().BoolVar(&measureOnly, "measure-only", false, "only measure existing .hurl files; skip LLM generation")
	cmd.Flags().StringVar(&model, "model", "ollama:gemma4:e4b", "LLM backend:model for generation (e.g. ollama:gemma4:e4b, claude:sonnet)")
	cmd.Flags().IntVar(&maxItems, "max-items", 0, "cap on DISTINCT endpoints to process (0 = all)")
	return cmd
}
