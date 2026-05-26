//ff:func feature=verify type=command control=sequence
//ff:what Defines the verify command that runs hurl tests and advances on pass
package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/park-jun-woo/huma/internal/config"
	"github.com/park-jun-woo/huma/internal/prompt"
	"github.com/park-jun-woo/huma/internal/rule"
	"github.com/park-jun-woo/huma/internal/runner"
	"github.com/park-jun-woo/huma/internal/scanner"
	"github.com/park-jun-woo/huma/internal/session"
	"github.com/spf13/cobra"
)

var verifyCmd = &cobra.Command{
	Use:   "verify",
	Short: "Run hurl test for current endpoint and advance if passing",
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
			total, pass, _ := sess.Stats()
			fmt.Print(prompt.AllComplete(pass, total))
			return nil
		}

		hurl := runner.FindHurlFile(ep, cfg.HurlDir)

		if cfg.Server.Start != "" {
			// LIVE MODE: healthcheck → hurl execution → coverage
			readyURL := cfg.BaseURL + cfg.Server.Ready
			if !probeCheckFn(readyURL) {
				fmt.Print(prompt.StartPrompt(cfg))
				return nil
			}
			if hurl == "" {
				detail := fmt.Sprintf("%s %s\n  Expected: %s", ep.Method, ep.Path, runner.HurlFileName(ep, cfg.HurlDir))
				fmt.Fprintln(os.Stderr, rule.H01.Format(detail))
				os.Exit(1)
			}
			return verifyWithCoverage(cfg, sess, ep, hurl)
		}

		// STATIC MODE: no server, static analysis only
		if hurl == "" {
			detail := fmt.Sprintf("%s %s\n  Expected: %s", ep.Method, ep.Path, runner.HurlFileName(ep, cfg.HurlDir))
			fmt.Fprintln(os.Stderr, rule.H01.Format(detail))
			os.Exit(1)
		}

		// .hurl exists → check static response coverage
		respResult := checkResponseCoverageFn(ep, hurl, cfg.Scan.Lang)
		if respResult != nil && respResult.Total > 0 && len(respResult.Missing) > 0 {
			sess.MarkImprove(ep.ID, respResult.Percent)
			if err := sess.Save(); err != nil {
				return err
			}
			fmt.Print(prompt.ResponseImprovePrompt(ep, hurl, respResult))
			return nil
		}

		sess.MarkPass(ep.ID)
		if err := sess.Save(); err != nil {
			return err
		}
		fmt.Print(prompt.PassPrompt(ep))
		fmt.Println()

		return nil
	},
}

func verifyWithCoverage(cfg *config.Config, sess *session.Session, ep *scanner.Endpoint, hurl string) error {
	entry := sess.CurrentEntry()
	a := newAdapterFn(cfg)

	if err := a.Build(); err != nil {
		return fmt.Errorf("%s", rule.A02.Format(err.Error()))
	}

	result, covResult, err := adapterRunFn(a, hurl, cfg.HurlVariables, ep.Source, ep.Handler)
	if err != nil {
		return err
	}

	if !result.Pass {
		fmt.Print(prompt.FailPrompt(ep, hurl, result.Feedback))
		return nil
	}

	// Hurl passed — check coverage
	if covResult == nil || covResult.Percent == 100 || covResult.Total == 0 {
		sess.MarkPass(ep.ID)
		if err := sess.Save(); err != nil {
			return err
		}
		fmt.Print(prompt.PassPrompt(ep))
		fmt.Println()
		return nil
	}

	// Coverage < 100%: check if improvement stalled
	if entry.ImproveCount >= 1 && covResult.Percent <= entry.PrevCoverage {
		sess.MarkDone(ep.ID, covResult.Percent)
		if err := sess.Save(); err != nil {
			return err
		}
		fmt.Print(prompt.PassPrompt(ep))
		fmt.Println()
		return nil
	}

	// Coverage improved or first attempt — mark improve
	sess.MarkImprove(ep.ID, covResult.Percent)
	if err := sess.Save(); err != nil {
		return err
	}
	fmt.Print(prompt.ImprovePrompt(ep, hurl, covResult))

	return nil
}

func init() {
	rootCmd.AddCommand(verifyCmd)
}
