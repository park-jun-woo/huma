//ff:func feature=adapter type=engine control=sequence
//ff:what Sends SIGINT to the Python server and waits for the .coverage file to be written
package adapter

import (
	"os"
	"time"
)

// Stop sends SIGINT to the Python server process and waits for it to exit.
// SIGINT triggers Python's KeyboardInterrupt, leading to graceful shutdown
// and coverage.py's atexit handler writing the .coverage file.
func (a *PythonAdapter) Stop() error {
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
