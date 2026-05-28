//ff:func feature=adapter type=adapter control=sequence
//ff:what Stops the Supabase Edge Function server process with SIGINT
package adapter

import (
	"os"
	"time"
)

// Stop sends SIGINT to the Deno server process and waits for it to exit.
func (d *DenoAdapter) Stop() error {
	if d.proc == nil || d.proc.Process == nil {
		return nil
	}

	if err := d.proc.Process.Signal(os.Interrupt); err != nil {
		return nil
	}

	done := make(chan error, 1)
	go func() {
		done <- d.proc.Wait()
	}()

	select {
	case <-done:
	case <-time.After(10 * time.Second):
		d.proc.Process.Kill()
		<-done
	}

	d.proc = nil
	return nil
}
