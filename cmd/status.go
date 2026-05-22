//ff:func feature=session type=command control=sequence
//ff:what Defines the status command that shows test progress summary
package cmd

import (
	"fmt"
	"os"

	"github.com/park-jun-woo/hurlfill/internal/config"
	"github.com/park-jun-woo/hurlfill/internal/session"
	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show progress summary",
	RunE: func(cmd *cobra.Command, args []string) error {
		_, err := config.Load()
		if err != nil {
			return fmt.Errorf("load config: %w", err)
		}

		sess, err := session.Load()
		if err != nil {
			fmt.Fprintln(os.Stderr, "No session found. Run 'hurlfill scan' first.")
			os.Exit(1)
		}

		total, done, todo := sess.Stats()
		fmt.Printf("%d endpoints\n", total)
		fmt.Printf("DONE: %3d (%5.1f%%)\n", done, pct(done, total))
		fmt.Printf("TODO: %3d (%5.1f%%)\n", todo, pct(todo, total))
		return nil
	},
}

func pct(n, total int) float64 {
	if total == 0 {
		return 0
	}
	return float64(n) / float64(total) * 100
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
