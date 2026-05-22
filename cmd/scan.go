//ff:func feature=scan type=command control=sequence
//ff:what Defines the scan command that reads endpoints from an external JSON file or stdin
package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/park-jun-woo/hurlfill/internal/config"
	"github.com/park-jun-woo/hurlfill/internal/scanner"
	"github.com/park-jun-woo/hurlfill/internal/session"
	"github.com/spf13/cobra"
)

var fromFlag string

var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Read endpoints from JSON file or stdin and create a session",
	RunE: func(cmd *cobra.Command, args []string) error {
		if fromFlag == "" {
			return fmt.Errorf("--from is required (file path or - for stdin)")
		}

		_, err := config.Load()
		if err != nil {
			return fmt.Errorf("load config: %w", err)
		}

		var data []byte
		if fromFlag == "-" {
			data, err = io.ReadAll(os.Stdin)
		} else {
			data, err = os.ReadFile(fromFlag)
		}
		if err != nil {
			return fmt.Errorf("read input: %w", err)
		}

		endpoints, err := scanner.ParseEndpoints(data)
		if err != nil {
			return fmt.Errorf("parse endpoints: %w", err)
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
	scanCmd.Flags().StringVar(&fromFlag, "from", "", "path to endpoints JSON file, or - for stdin")
	rootCmd.AddCommand(scanCmd)
}
