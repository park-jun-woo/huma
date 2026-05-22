package cmd

import (
	"fmt"
	"os"

	"github.com/park-jun-woo/hurlfill/internal/config"
	"github.com/park-jun-woo/hurlfill/internal/prompt"
	"github.com/park-jun-woo/hurlfill/internal/runner"
	"github.com/park-jun-woo/hurlfill/internal/scanner"
	"github.com/park-jun-woo/hurlfill/internal/session"
	"github.com/spf13/cobra"
)

var verifyCmd = &cobra.Command{
	Use:   "verify",
	Short: "Run hurl test for current endpoint and advance if passing",
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
			total, pass, _ := sess.Stats()
			fmt.Print(prompt.AllComplete(pass, total))
			return nil
		}

		hurl := runner.FindHurlFile(ep, cfg.HurlDir)
		if hurl == "" {
			fmt.Fprintf(os.Stderr, "No .hurl file found for %s %s\n", ep.Method, ep.Path)
			fmt.Fprintf(os.Stderr, "Expected: %s\n", runner.HurlFileName(ep, cfg.HurlDir))
			os.Exit(1)
		}

		// Coverage mode: server config is present
		if cfg.Server.Build != "" {
			return verifyWithCoverage(cfg, sess, ep, hurl)
		}

		// No coverage mode — behave as before
		result, err := runner.Run(hurl, cfg.BaseURL)
		if err != nil {
			return fmt.Errorf("hurl run failed: %w", err)
		}

		if result.Pass {
			sess.MarkPass(ep.ID)
			if err := sess.Save(); err != nil {
				return err
			}
			fmt.Print(prompt.PassPrompt(ep))
		} else {
			fmt.Print(prompt.FailPrompt(ep, hurl, result.Feedback))
		}

		return nil
	},
}

func verifyWithCoverage(cfg *config.Config, sess *session.Session, ep *scanner.Endpoint, hurl string) error {
	entry := sess.CurrentEntry()
	a := newAdapterFn(cfg)

	if err := a.Build(); err != nil {
		return fmt.Errorf("build: %w", err)
	}

	result, covResult, err := adapterRunFn(a, hurl, cfg.BaseURL, ep.Source, ep.Handler)
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
		return nil
	}

	// Coverage < 100%: check if improvement stalled
	if entry.ImproveCount >= 1 && covResult.Percent <= entry.PrevCoverage {
		sess.MarkDone(ep.ID, covResult.Percent)
		if err := sess.Save(); err != nil {
			return err
		}
		fmt.Print(prompt.PassPrompt(ep))
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
