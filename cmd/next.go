//ff:func feature=ratchet type=command control=sequence
//ff:what Defines the next command that shows the next untested endpoint or verifies the current one
package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/park-jun-woo/huma/internal/adapter"
	"github.com/park-jun-woo/huma/internal/config"
	"github.com/park-jun-woo/huma/internal/prompt"
	"github.com/park-jun-woo/huma/internal/rule"
	"github.com/park-jun-woo/huma/internal/runner"
	"github.com/park-jun-woo/huma/internal/scanner"
	"github.com/park-jun-woo/huma/internal/session"
	"github.com/spf13/cobra"
)

var nextCmd = &cobra.Command{
	Use:   "next",
	Short: "Show the next untested endpoint, or verify the current one",
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
				fmt.Print(prompt.TodoPrompt(ep, cfg.HurlDir, cfg.URLVar()))
				return nil
			}
			return runWithCoverage(cfg, sess, ep, hurl)
		}

		// STATIC MODE: no server, static analysis only
		if hurl == "" {
			branches := responseBranches(ep, cfg.Scan.Lang)
			fmt.Print(prompt.StaticTodoPrompt(ep, cfg.HurlDir, cfg.URLVar(), branches))
			return nil
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

		// All status codes covered → PASS
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
			nextBranches := responseBranches(next, cfg.Scan.Lang)
			fmt.Print(prompt.StaticTodoPrompt(next, cfg.HurlDir, cfg.URLVar(), nextBranches))
		}

		return nil
	},
}

// probeCheckFn allows tests to replace adapter.ProbeCheck.
var probeCheckFn = adapter.ProbeCheck

// adapterRunFn allows tests to replace adapter.RunWithCoverage.
var adapterRunFn = adapter.RunWithCoverage

// newAdapterFn allows tests to replace adapter.NewGoAdapter.
var newAdapterFn = func(cfg *config.Config) adapter.Adapter {
	switch cfg.Scan.Lang {
	case "python", "flask", "django", "drf":
		return adapter.NewPythonAdapter(cfg)
	case "node", "javascript", "typescript", "nestjs", "express", "fastify", "hono":
		return adapter.NewNodeAdapter(cfg)
	case "deno", "edge-functions":
		return adapter.NewDenoAdapter(cfg)
	case "java", "spring", "quarkus":
		return adapter.NewJavaAdapter(cfg)
	case "dotnet", "aspnet", "csharp":
		return adapter.NewDotnetAdapter(cfg)
	case "php", "laravel":
		return adapter.NewPhpAdapter(cfg)
	case "rust", "actix":
		return adapter.NewRustAdapter(cfg)
	case "go", "fiber", "echo":
		return adapter.NewGoAdapter(cfg)
	default:
		return adapter.NewGoAdapter(cfg)
	}
}

func runWithCoverage(cfg *config.Config, sess *session.Session, ep *scanner.Endpoint, hurl string) error {
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
