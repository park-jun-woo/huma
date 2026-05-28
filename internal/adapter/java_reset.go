//ff:func feature=adapter type=engine control=sequence
//ff:what Removes and recreates the JaCoCo coverage data directory
package adapter

import (
	"fmt"
	"os"
)

// Reset removes and recreates the JaCoCo coverage data directory.
func (a *JavaAdapter) Reset() error {
	if err := os.RemoveAll(a.jacocoDir); err != nil {
		return fmt.Errorf("remove jacoco dir: %w", err)
	}
	if err := os.MkdirAll(a.jacocoDir, 0o755); err != nil {
		return fmt.Errorf("create jacoco dir: %w", err)
	}
	return nil
}
