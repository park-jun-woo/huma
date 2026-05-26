//ff:func feature=prompt type=command control=sequence
//ff:what Defines the prompt command that outputs agent instruction for the current TODO endpoint
package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/park-jun-woo/huma/internal/config"
	"github.com/park-jun-woo/huma/internal/prompt"
	"github.com/park-jun-woo/huma/internal/rule"
	"github.com/park-jun-woo/huma/internal/session"
	"github.com/spf13/cobra"
)

var promptCmd = &cobra.Command{
	Use:   "prompt",
	Short: "Output agent prompt for current TODO endpoint (no side effects)",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if errors.Is(err, config.ErrNoManifest) {
			fmt.Print(prompt.SetupPrompt())
			return nil
		}
		if err != nil {
			return fmt.Errorf("load config: %w", err)
		}

		sess, err := session.Load()
		if err != nil {
			fmt.Fprintln(os.Stderr, rule.S01.Format("Run 'huma scan' first."))
			os.Exit(1)
		}

		ep := sess.Current()
		if ep == nil {
			// No TODO endpoints remain — exit 1 so `while huma prompt` stops
			os.Exit(1)
		}

		fmt.Print(prompt.TodoPrompt(ep, cfg.HurlDir, cfg.URLVar()))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(promptCmd)
}
