package cmd

import (
	"fmt"
	"os"

	"github.com/park-jun-woo/hurlfill/internal/config"
	"github.com/park-jun-woo/hurlfill/internal/prompt"
	"github.com/park-jun-woo/hurlfill/internal/session"
	"github.com/spf13/cobra"
)

var promptCmd = &cobra.Command{
	Use:   "prompt",
	Short: "Output agent prompt for current TODO endpoint (no side effects)",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("load config: %w", err)
		}

		sess, err := session.Load()
		if err != nil {
			fmt.Fprintln(os.Stderr, "No session found. Run 'hurlfill scan' first.")
			os.Exit(1)
		}

		ep := sess.Current()
		if ep == nil {
			// No TODO endpoints remain — exit 1 so `while hurlfill prompt` stops
			os.Exit(1)
		}

		fmt.Print(prompt.TodoPrompt(ep, cfg.HurlDir))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(promptCmd)
}
