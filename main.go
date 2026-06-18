//ff:func feature=cli type=command control=sequence
//ff:what 프로그램 진입점. reins cli.NewQuestCmd로 huma 퀘스트 CLI를 조립해 실행한다. 도메인 로직은 humaquest.Def() 하나만 끼운다.

package main

import (
	"os"

	"github.com/park-jun-woo/huma/internal/humaquest"
	"github.com/park-jun-woo/reins/pkg/cli"
	"github.com/spf13/cobra"
)

// Version is the build-time default, overridable via -ldflags "-X main.Version=...".
var Version = "v0.3.0"

func main() {
	def := humaquest.Def()
	root := cli.NewQuestCmd("huma", def, cli.Options{
		Version: Version,
		// huma's `loop` owns the session-scoped server lifecycle (Phase 006/010): it
		// brings the server up once, has an LLM generate .hurl for each TODO endpoint,
		// measures live runtime coverage, injects it into the CRI gate, retries until
		// convergence, and tears down. reins' canonical submit/loop are one-shot and
		// have no server bring-up hook, so this huma-specific command supplies it.
		// reins' static Loop option is intentionally omitted: enabling it would register
		// a second `loop` command and collide with huma's live loop.
		ExtraCommands: []*cobra.Command{humaquest.NewLoopCmd(def)},
	})
	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}
