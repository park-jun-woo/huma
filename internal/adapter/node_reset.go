//ff:func feature=adapter type=engine control=sequence
//ff:what Removes and recreates the V8 coverage data directory
package adapter

import (
	"fmt"
	"os"
)

// Reset removes and recreates the V8 coverage data directory.
func (a *NodeAdapter) Reset() error {
	if err := os.RemoveAll(a.coverDir); err != nil {
		return fmt.Errorf("remove cover dir: %w", err)
	}
	if err := os.MkdirAll(a.coverDir, 0o755); err != nil {
		return fmt.Errorf("create cover dir: %w", err)
	}
	return nil
}
