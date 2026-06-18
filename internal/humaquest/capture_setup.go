//ff:func feature=gate type=engine control=sequence level=error
//ff:what Phase 009 / 2-A capture: brings the server up once, runs testing.setup.hurl via `hurl --json`, parses its [Captures] into a map[string]string (token + any fixtures) for injection into every endpoint's hurl run
package humaquest

import (
	"fmt"

	"github.com/park-jun-woo/huma/internal/adapter"
	"github.com/park-jun-woo/huma/internal/config"
	"github.com/park-jun-woo/huma/internal/runner"
)

// captureSetup runs the user-authored setup .hurl (testing.setup.hurl) once before
// the cover loop and returns the variables it captured (token, fixture IDs, …) as a
// flat map. It is the primary (2-A) Phase 009 path: signature-agnostic and universal —
// whatever the setup hurl declares in [Captures] comes out, so fixtures come free.
//
// It owns a self-contained server lifecycle (Build is idempotent — the per-endpoint
// Up() reuses the same built binary; Start → WaitReady → run setup → Stop) because the
// login must hit a live server. Capture happens once per loop: JWT is stateless across
// the per-endpoint restarts (Phase 006), so one token stays valid.
//
// Policy: no setup configured -> empty map, no error (caller proceeds token-less).
// Setup configured but the run/parse fails (file missing, login 401, bad JSON) -> the
// error is returned so the caller can warn loudly; a silent token-less proceed would
// 401 every protected endpoint and waste a whole generate loop.
func captureSetup(a adapter.Adapter, cfg *config.Config) (map[string]string, error) {
	if cfg.Setup.Hurl == "" {
		return map[string]string{}, nil
	}

	if err := a.Build(); err != nil {
		return nil, fmt.Errorf("setup: build server: %w", err)
	}
	if err := a.Start(); err != nil {
		return nil, fmt.Errorf("setup: start server: %w", err)
	}
	defer a.Stop()
	if err := a.WaitReady(); err != nil {
		return nil, fmt.Errorf("setup: wait ready: %w", err)
	}

	captures, err := runner.RunJSON(cfg.Setup.Hurl, cfg.HurlVariables)
	if err != nil {
		return nil, fmt.Errorf("setup capture (%s): %w", cfg.Setup.Hurl, err)
	}
	return captures, nil
}
