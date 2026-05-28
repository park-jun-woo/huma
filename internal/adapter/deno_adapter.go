//ff:type feature=adapter type=adapter
//ff:what DenoAdapter manages a Supabase Edge Function server process without coverage support
package adapter

import (
	"os/exec"

	"github.com/park-jun-woo/huma/internal/config"
)

// DenoAdapter implements Adapter for Deno / Supabase Edge Functions.
// Coverage is not supported — Collect always returns nil.
type DenoAdapter struct {
	cfg     *config.ServerConfig
	baseURL string
	proc    *exec.Cmd
	built   bool
}
