//ff:type feature=adapter type=adapter
//ff:what PhpAdapter manages a PHP/Laravel server process without coverage support
package adapter

import (
	"os/exec"

	"github.com/park-jun-woo/huma/internal/config"
)

// PhpAdapter implements Adapter for PHP/Laravel applications.
// Coverage is not supported — Collect always returns nil.
type PhpAdapter struct {
	cfg     *config.ServerConfig
	baseURL string
	proc    *exec.Cmd
	built   bool
}
