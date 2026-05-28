//ff:func feature=adapter type=adapter control=sequence
//ff:what Starts the Supabase Edge Function server process
package adapter

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Start launches the Deno / Supabase Edge Function server process.
func (d *DenoAdapter) Start() error {
	parts := strings.Fields(d.cfg.Start)
	if len(parts) == 0 {
		return fmt.Errorf("empty start command")
	}

	cmd := exec.Command(parts[0], parts[1:]...)
	cmd.Stdout = os.Stderr
	cmd.Stderr = os.Stderr

	env := os.Environ()
	for k, v := range d.cfg.Env {
		env = append(env, k+"="+v)
	}
	cmd.Env = env

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("start server: %w", err)
	}
	d.proc = cmd
	return nil
}
