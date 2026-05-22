package cmd

import (
	"fmt"

	"github.com/park-jun-woo/hurlfill/internal/config"
	"github.com/park-jun-woo/hurlfill/internal/scanner"
	"github.com/park-jun-woo/hurlfill/internal/session"
	"github.com/spf13/cobra"
)

var scanCmd = &cobra.Command{
	Use:   "scan [dir]",
	Short: "Scan source code and index all API endpoints",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		_, err := config.Load()
		if err != nil {
			return fmt.Errorf("load config: %w", err)
		}

		dir := "."
		if len(args) > 0 {
			dir = args[0]
		}

		endpoints, err := scanner.Scan(dir)
		if err != nil {
			return fmt.Errorf("scan failed: %w", err)
		}

		sess, err := session.Load()
		if err != nil {
			sess = session.New()
		}
		sess.Merge(endpoints)
		if err := sess.Save(); err != nil {
			return fmt.Errorf("save session: %w", err)
		}

		fmt.Printf("Scanned %d endpoints\n", len(endpoints))
		for _, ep := range endpoints {
			fmt.Printf("  %s %s  (%s)\n", ep.Method, ep.Path, ep.Source)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(scanCmd)
}
