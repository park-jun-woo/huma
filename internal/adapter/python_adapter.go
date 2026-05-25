//ff:type feature=adapter type=adapter
//ff:what PythonAdapter manages a Python server process with coverage.py support
package adapter

import (
	"os/exec"

	"github.com/park-jun-woo/huma/internal/config"
)

const coverageFile = ".coverage"
const coverageJSON = ".huma/cov.json"

// PythonAdapter implements Adapter using coverage.py for Python servers.
type PythonAdapter struct {
	cfg     *config.ServerConfig
	baseURL string
	proc    *exec.Cmd
}
