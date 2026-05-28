//ff:type feature=adapter type=adapter
//ff:what RustAdapter manages a Rust/Actix-web server process without coverage support
package adapter

import (
	"os/exec"

	"github.com/park-jun-woo/huma/internal/config"
)

// RustAdapter implements Adapter for Rust / Actix-web servers.
// Coverage is not supported — Collect always returns nil.
type RustAdapter struct {
	cfg     *config.ServerConfig
	baseURL string
	proc    *exec.Cmd
	built   bool
}
