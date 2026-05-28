//ff:func feature=adapter type=adapter control=sequence
//ff:what Runs the build command for Deno adapter if configured, otherwise no-op
package adapter

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Build runs the configured build command if set. If no build command is
// configured, it is a no-op. Skips if already built.
func (d *DenoAdapter) Build() error {
	if d.built {
		return nil
	}

	parts := strings.Fields(d.cfg.Build)
	if len(parts) == 0 {
		return nil
	}

	cmd := exec.Command(parts[0], parts[1:]...)
	cmd.Stdout = os.Stderr
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("build failed: %w", err)
	}
	d.built = true
	return nil
}
