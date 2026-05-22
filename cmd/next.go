//ff:func feature=ratchet type=command control=sequence
//ff:what Defines the next command that shows the next untested endpoint or verifies the current one
package cmd

import (
	"fmt"
	"os"

	"github.com/park-jun-woo/hurlfill/internal/adapter"
	"github.com/park-jun-woo/hurlfill/internal/config"
	"github.com/park-jun-woo/hurlfill/internal/prompt"
	"github.com/park-jun-woo/hurlfill/internal/runner"
	"github.com/park-jun-woo/hurlfill/internal/scanner"
	"github.com/park-jun-woo/hurlfill/internal/session"
	"github.com/spf13/cobra"
)

var nextCmd = &cobra.Command{
	Use:   "next",
	Short: "Show the next untested endpoint, or verify the current one",
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
			// TODO: no .hurl file yet
			fmt.Print(prompt.TodoPrompt(ep, cfg.HurlDir, cfg.URLVar()))
			return nil
		}

		// Coverage mode: server config is present
		if cfg.Server.Start != "" {
			return runWithCoverage(cfg, sess, ep, hurl)
		}

		// No coverage mode — behave as before (hurl pass/fail only)
		result, err := runner.Run(hurl, cfg.HurlVariables)
		if err != nil {
			return fmt.Errorf("hurl run failed: %w", err)
		}

		if result.Pass {
			sess.MarkPass(ep.ID)
			if err := sess.Save(); err != nil {
				return err
			}
			fmt.Print(prompt.PassPrompt(ep))
			fmt.Println()

			next := sess.Current()
			if next == nil {
				total, pass, _ := sess.Stats()
				fmt.Print(prompt.AllComplete(pass, total))
			} else {
				fmt.Print(prompt.TodoPrompt(next, cfg.HurlDir, cfg.URLVar()))
			}
		} else {
			fmt.Print(prompt.FailPrompt(ep, hurl, result.Feedback))
		}

		return nil
	},
}

// adapterRunFn allows tests to replace adapter.RunWithCoverage.
var adapterRunFn = adapter.RunWithCoverage

// newAdapterFn allows tests to replace adapter.NewGoAdapter.
var newAdapterFn = func(cfg *config.Config) adapter.Adapter {
	if cfg.Scan.Lang == "python" {
		return adapter.NewPythonAdapter(cfg)
	}
	return adapter.NewGoAdapter(cfg)
}

func runWithCoverage(cfg *config.Config, sess *session.Session, ep *scanner.Endpoint, hurl string) error {
	entry := sess.CurrentEntry()
	a := newAdapterFn(cfg)

	if err := a.Build(); err != nil {
		return fmt.Errorf("build: %w", err)
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

		next := sess.Current()
		if next == nil {
			total, pass, _ := sess.Stats()
			fmt.Print(prompt.AllComplete(pass, total))
		} else {
			fmt.Print(prompt.TodoPrompt(next, cfg.HurlDir, cfg.URLVar()))
		}
		return nil
	}

	// Coverage < 100%: check if improvement stalled
	if entry.ImproveCount >= 1 && covResult.Percent <= entry.PrevCoverage {
		// No improvement after retry — mark done
		sess.MarkDone(ep.ID, covResult.Percent)
		if err := sess.Save(); err != nil {
			return err
		}
		fmt.Print(prompt.PassPrompt(ep))
		fmt.Println()

		next := sess.Current()
		if next == nil {
			total, pass, _ := sess.Stats()
			fmt.Print(prompt.AllComplete(pass, total))
		} else {
			fmt.Print(prompt.TodoPrompt(next, cfg.HurlDir, cfg.URLVar()))
		}
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
	rootCmd.AddCommand(nextCmd)
}
