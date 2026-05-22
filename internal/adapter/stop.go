//ff:func feature=adapter type=engine control=sequence
//ff:what Sends SIGINT to the server process and waits for graceful shutdown
package adapter

import (
	"os"
	"time"
)

// Stop sends SIGINT to the server process and waits for it to exit.
// SIGINT (not SIGTERM) triggers Go's coverage dump on graceful shutdown.
func (a *GoAdapter) Stop() error {
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
