//ff:type feature=adapter type=adapter
//ff:what DotnetAdapter manages an ASP.NET Core server process without coverage support
package adapter

import (
	"os/exec"

	"github.com/park-jun-woo/huma/internal/config"
)

// DotnetAdapter implements Adapter for ASP.NET Core applications.
// Coverage is not supported — Collect always returns nil.
type DotnetAdapter struct {
	cfg     *config.ServerConfig
	baseURL string
	proc    *exec.Cmd
	built   bool
}
