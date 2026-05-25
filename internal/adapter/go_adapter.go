//ff:type feature=adapter type=adapter
//ff:what GoAdapter manages a Go server process with integration test coverage support
package adapter

import (
	"os/exec"

	"github.com/park-jun-woo/huma/internal/config"
)

const coverDir = ".huma/coverdata"
const coverOut = ".huma/coverage.out"

// GoAdapter implements Adapter using Go 1.20+ integration test coverage.
type GoAdapter struct {
	cfg      *config.ServerConfig
	baseURL  string
	coverDir string
	proc     *exec.Cmd
	built    bool
}
