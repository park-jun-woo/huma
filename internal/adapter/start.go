//ff:func feature=adapter type=engine control=sequence
//ff:what Launches the server process with GOCOVERDIR environment variable set
package adapter

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Start launches the server process with GOCOVERDIR set.
func (a *GoAdapter) Start() error {
	parts := strings.Fields(a.cfg.Start)
	if len(parts) == 0 {
		return fmt.Errorf("empty start command")
	}

	cmd := exec.Command(parts[0], parts[1:]...)
	cmd.Stdout = os.Stderr
	cmd.Stderr = os.Stderr

	env := os.Environ()
	absDir, err := filepath.Abs(a.coverDir)
	if err != nil {
		return fmt.Errorf("abs cover dir: %w", err)
	}
	env = append(env, "GOCOVERDIR="+absDir)
	for k, v := range a.cfg.Env {
		env = append(env, k+"="+v)
	}
	cmd.Env = env

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("start server: %w", err)
	}
	a.proc = cmd
	return nil
}
