//ff:type feature=adapter type=adapter
//ff:what JavaAdapter manages a Java server process with JaCoCo coverage support
package adapter

import (
	"os/exec"

	"github.com/park-jun-woo/huma/internal/config"
)

const jacocoDir = ".huma/jacoco"
const jacocoExecFile = ".huma/jacoco/jacoco.exec"
const jacocoXMLFile = ".huma/jacoco/jacoco.xml"

// JavaAdapter implements Adapter using JaCoCo for Java coverage collection.
type JavaAdapter struct {
	cfg       *config.ServerConfig
	baseURL   string
	jacocoDir string
	proc      *exec.Cmd
	built     bool
}
