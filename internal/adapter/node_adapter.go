//ff:type feature=adapter type=adapter
//ff:what NodeAdapter manages a Node.js server process with V8 coverage support
package adapter

import (
	"os/exec"

	"github.com/park-jun-woo/hurlfill/internal/config"
)

const nodeCoverDir = ".hurlfill/v8cov"
const istanbulOutDir = ".hurlfill/istanbul"
const istanbulOutFile = ".hurlfill/istanbul/coverage-final.json"

// NodeAdapter implements Adapter using Node.js V8 built-in coverage.
type NodeAdapter struct {
	cfg      *config.ServerConfig
	baseURL  string
	coverDir string
	proc     *exec.Cmd
	built    bool
}
