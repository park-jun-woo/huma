//ff:func feature=gate type=command control=sequence level=error
//ff:what `huma cover` ExtraCommand를 만든다. 루트의 persistent 플래그(--session/--out)를 런타임에 상속받아 읽고, manifest를 로드해 실서버 probe(newLiveProbe)를 구성한 뒤 runCover에 넘긴다. --generate면 buildBackend(--model)로 LLM backend를 만들어 무인 생성 루프를 켜고, --max-items로 처리할 distinct endpoint 수를 캡한다. 서버 라이프사이클을 소유하는 huma 전용 명령 — 계측 바이너리를 한 번 컴파일하고, 엔드포인트마다 Start/Stop/Collect로 측정하며 전 TODO를 순회한다(Go 커버리지는 프로세스 종료 시 flush).

package humaquest

import (
	"fmt"

	"github.com/park-jun-woo/huma/internal/config"
	"github.com/park-jun-woo/reins/pkg/gate"
	"github.com/spf13/cobra"
)

// NewCoverCmd builds the `cover` ExtraCommand: huma's server-owning loop. It is
// attached via cli.Options.ExtraCommands and reuses the root's persistent --session
// and --out flags (read at run time from the inherited flag set). It loads the
// manifest, constructs the real liveProbe (the default coverage source), and hands
// orchestration to runCover. The probe seam means the real adapter is the default
// here while runCover stays unit-testable with a fake.
func NewCoverCmd(def gate.Definition) *cobra.Command {
	var generate bool
	var model string
	var maxItems int
	cmd := &cobra.Command{
		Use:   "cover",
		Short: "compile once, then measure live coverage per endpoint (start/stop each) over all TODO endpoints",
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
				return fmt.Errorf("cover needs a manifest with testing.server: %w", err)
			}
			backend, err := buildBackend(generate, model)
			if err != nil {
				return err
			}
			return runCover(def, newLiveProbe(cfg), backend, maxItems, sessionPath, outPath, cmd.OutOrStdout())
		},
	}
	cmd.Flags().BoolVar(&generate, "generate", false, "generate .hurl files with an LLM and converge them via the CRI gate (opt-in)")
	cmd.Flags().StringVar(&model, "model", "claude:sonnet", "LLM backend:model for --generate (e.g. claude:sonnet, ollama:gemma3)")
	cmd.Flags().IntVar(&maxItems, "max-items", 0, "cap on DISTINCT endpoints to process (0 = all)")
	return cmd
}
