//ff:func feature=adapter type=adapter control=sequence
//ff:what Stops the Rust/Actix-web server process with SIGINT
package adapter

import (
	"os"
	"time"
)

// Stop sends SIGINT to the Rust server process and waits for it to exit.
func (a *RustAdapter) Stop() error {
	if a.proc == nil || a.proc.Process == nil {
		return nil
	}

	if err := a.proc.Process.Signal(os.Interrupt); err != nil {
		return nil
	}

	done := make(chan error, 1)
	go func() {
		done <- a.proc.Wait()
	}()

	select {
	case <-done:
	case <-time.After(10 * time.Second):
		a.proc.Process.Kill()
		<-done
	}

	a.proc = nil
	return nil
}
