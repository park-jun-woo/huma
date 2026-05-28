//ff:func feature=adapter type=adapter control=sequence
//ff:what Starts the PHP/Laravel server process
package adapter

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Start launches the PHP/Laravel server process.
func (p *PhpAdapter) Start() error {
	parts := strings.Fields(p.cfg.Start)
	if len(parts) == 0 {
		return fmt.Errorf("empty start command")
	}

	cmd := exec.Command(parts[0], parts[1:]...)
	cmd.Stdout = os.Stderr
	cmd.Stderr = os.Stderr

	env := os.Environ()
	for k, v := range p.cfg.Env {
		env = append(env, k+"="+v)
	}
	cmd.Env = env

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("start server: %w", err)
	}
	p.proc = cmd
	return nil
}
