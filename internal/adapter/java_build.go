//ff:func feature=adapter type=engine control=sequence
//ff:what Runs the configured build command for Java, skipping if already built
package adapter

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Build runs the configured build command (e.g., mvn package -DskipTests).
// If no build command is configured, it is a no-op. Skips if already built.
func (a *JavaAdapter) Build() error {
	if a.built {
		return nil
	}

	parts := strings.Fields(a.cfg.Build)
	if len(parts) == 0 {
		return nil
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
