//ff:func feature=adapter type=engine control=sequence
//ff:what Launches the Python server process with coverage.py instrumentation
package adapter

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Start launches the Python server process with the configured start command.
// The start command typically includes `coverage run` for instrumentation.
func (a *PythonAdapter) Start() error {
	parts := strings.Fields(a.cfg.Start)
	if len(parts) == 0 {
		return fmt.Errorf("empty start command")
	}

	cmd := exec.Command(parts[0], parts[1:]...)
	cmd.Stdout = os.Stderr
	cmd.Stderr = os.Stderr

	env := os.Environ()
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
