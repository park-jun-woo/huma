//ff:func feature=scan type=command control=sequence
//ff:what Defines the scan command that reads endpoints from OpenAPI, JSON, or YAML input
package cmd

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/park-jun-woo/huma/internal/config"
	"github.com/park-jun-woo/huma/internal/rule"
	"github.com/park-jun-woo/huma/internal/scanner"
	"github.com/park-jun-woo/huma/internal/session"
	"github.com/spf13/cobra"
)

var fromFlag string
var linkSourceFlag string

var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Read endpoints from OpenAPI, JSON, or stdin and create a session",
	RunE: func(cmd *cobra.Command, args []string) error {
		if fromFlag == "" {
			found := scanner.FindOpenAPIFile()
			if found == "" {
				return fmt.Errorf("%s", rule.E01.Format("Use --from to specify."))
			}
			fromFlag = found
		}

		cfg, err := config.Load()
		if err != nil && !errors.Is(err, config.ErrNoManifest) {
			return fmt.Errorf("load config: %w", err)
		}

		hurlDir := "hurl"
		if cfg != nil {
			hurlDir = cfg.HurlDir
		}

		var endpoints []scanner.Endpoint

		info, statErr := os.Stat(fromFlag)
		if fromFlag != "-" && statErr == nil && info.IsDir() {
			endpoints, err = scanner.ParseEdgeFunctions(fromFlag)
			if err != nil {
				return fmt.Errorf("scan edge functions: %w", err)
			}
		} else {
			var data []byte
			if fromFlag == "-" {
				data, err = io.ReadAll(os.Stdin)
			} else {
				data, err = os.ReadFile(fromFlag)
			}
			if err != nil {
				return fmt.Errorf("read input: %w", err)
			}

			endpoints, err = scanner.ParseEndpoints(data)
			if err != nil {
				return fmt.Errorf("parse endpoints: %w", err)
			}
		}

		if linkSourceFlag != "" {
			linked := scanner.LinkSource(endpoints, linkSourceFlag)
			fmt.Printf("Linked %d/%d endpoints to source under %s\n", linked, len(endpoints), linkSourceFlag)
		}

		sess, err := session.Load()
		if err != nil {
			sess = session.New()
		}
		sess.Merge(endpoints)
		if err := sess.Save(); err != nil {
			return fmt.Errorf("save session: %w", err)
		}

		precheckEndpoints(sess, cfg)
		if err := sess.Save(); err != nil {
			return fmt.Errorf("save session: %w", err)
		}

		var passCount, improveCount, todoCount, unverifiedCount int
		for _, e := range sess.Entries {
			switch e.Status {
			case session.StatusPass, session.StatusDone:
				passCount++
			case session.StatusImprove:
				improveCount++
			case session.StatusUnverified:
				unverifiedCount++
			default:
				todoCount++
			}
		}

		fmt.Printf("Scanned %d endpoints\n", len(endpoints))
		fmt.Printf("  PASS:       %d (existing hurl, all responses covered)\n", passCount)
		fmt.Printf("  IMPROVE:    %d (existing hurl, missing responses)\n", improveCount)
		fmt.Printf("  UNVERIFIED: %d (no oracle: source unlinked / not yet executed)\n", unverifiedCount)
		fmt.Printf("  TODO:       %d (no hurl file)\n", todoCount)

		warnMismatchedHurlFiles(hurlDir, endpoints)

		return nil
	},
}

func init() {
	scanCmd.Flags().StringVar(&fromFlag, "from", "", "path to endpoints JSON file, or - for stdin")
	scanCmd.Flags().StringVar(&linkSourceFlag, "link-source", "", "root directory to map OpenAPI handlers to source file:line")
	rootCmd.AddCommand(scanCmd)
}
