//ff:func feature=adapter type=engine control=sequence
//ff:what Runs the configured build command for Node.js, skipping if already built
package adapter

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Build runs the configured build command (e.g., npm install). Skips if already built.
func (a *NodeAdapter) Build() error {
	if a.built {
		return nil
	}

	parts := strings.Fields(a.cfg.Build)
	if len(parts) == 0 {
		return fmt.Errorf("empty build command")
	}

	cmd := exec.Command(parts[0], parts[1:]...)
	cmd.Stdout = os.Stderr
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("build failed: %w", err)
	}
	a.built = true
	return nil
}
