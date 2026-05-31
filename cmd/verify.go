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

		// .hurl exists → compute the static cheese-resistant verdict
		oc, respResult := staticVerdict(cfg, sess, ep, hurl)
		if err := sess.Save(); err != nil {
			return err
		}
		if printStaticNonPass(oc, ep, hurl, respResult, cfg) {
			return nil
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

	// Hurl passed — compute the live cheese-resistant verdict.
	oc := liveVerdict(cfg, sess, ep, hurl, covResult, entry)
	if err := sess.Save(); err != nil {
		return err
	}
	if printLiveNonPass(oc, ep, hurl, covResult, cfg) {
		return nil
	}
	fmt.Print(prompt.PassPrompt(ep))
	fmt.Println()
	return nil
}

func init() {
	rootCmd.AddCommand(verifyCmd)
}
