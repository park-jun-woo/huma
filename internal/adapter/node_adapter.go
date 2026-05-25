//ff:type feature=adapter type=adapter
//ff:what NodeAdapter manages a Node.js server process with V8 coverage support
package adapter

import (
	"os/exec"

	"github.com/park-jun-woo/huma/internal/config"
)

const nodeCoverDir = ".huma/v8cov"
const istanbulOutDir = ".huma/istanbul"
const istanbulOutFile = ".huma/istanbul/coverage-final.json"

// NodeAdapter implements Adapter using Node.js V8 built-in coverage.
type NodeAdapter struct {
	cfg      *config.ServerConfig
	baseURL  string
	coverDir string
	proc     *exec.Cmd
	built    bool
}
