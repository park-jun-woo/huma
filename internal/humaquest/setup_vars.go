//ff:func feature=gate type=engine control=selection level=error
//ff:what Phase 009 / C5 dispatch: picks the test-variable path (testing.setup capture first, else testing.auth mint) and returns the captured/minted vars; logs the captured KEYS (never values) for transparency
package humaquest

import (
	"fmt"
	"io"

	"github.com/park-jun-woo/huma/internal/adapter"
	"github.com/park-jun-woo/huma/internal/config"
)

// setupVars resolves the dynamic hurl variables to inject before the cover loop.
// Capture (2-A, testing.setup) is the primary path and is tried first; if no setup
// is configured but testing.auth is, it mints an HS256 token (2-B). Nothing
// configured -> empty map. It logs only the captured KEYS (token, building_id, …),
// never the values, so a token never lands in logs.
//
// Error policy: a capture/mint error is surfaced as a loud warning and the run
// CONTINUES token-less. This is deliberate — a token-less run still produces the
// per-endpoint diagnostic (which endpoints 401), which is exactly the §5 signal, and
// failing fast would lose that. The warning makes the missing token unmissable.
func setupVars(a adapter.Adapter, cfg *config.Config, w io.Writer) map[string]string {
	var (
		vars map[string]string
		err  error
		path string
	)
	switch {
	case cfg.Setup.Hurl != "":
		path = "capture(testing.setup)"
		vars, err = captureSetup(a, cfg)
	case cfg.Auth.SecretEnv != "":
		path = "mint(testing.auth)"
		vars, err = mintToken(cfg)
	default:
		return map[string]string{}
	}

	if err != nil {
		fmt.Fprintf(w, "warning: %s failed, continuing token-less: %v\n", path, err)
		return map[string]string{}
	}
	if len(vars) > 0 {
		fmt.Fprintf(w, "setup: %s captured vars %v\n", path, sortedKeys(vars))
	}
	return vars
}
