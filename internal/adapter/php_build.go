//ff:func feature=adapter type=adapter control=sequence
//ff:what Runs the build command for PHP adapter if configured, otherwise no-op
package adapter

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Build runs the configured build command if set. If no build command is
// configured, it is a no-op. Skips if already built.
func (p *PhpAdapter) Build() error {
	if p.built {
		return nil
	}

	parts := strings.Fields(p.cfg.Build)
	if len(parts) == 0 {
		return nil
	}

	cmd := exec.Command(parts[0], parts[1:]...)
	cmd.Stdout = os.Stderr
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("build failed: %w", err)
	}
	p.built = true
	return nil
}
