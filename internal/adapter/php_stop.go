//ff:func feature=adapter type=adapter control=sequence
//ff:what Stops the PHP/Laravel server process with SIGINT
package adapter

import (
	"os"
	"time"
)

// Stop sends SIGINT to the PHP server process and waits for it to exit.
func (p *PhpAdapter) Stop() error {
	if p.proc == nil || p.proc.Process == nil {
		return nil
	}

	if err := p.proc.Process.Signal(os.Interrupt); err != nil {
		return nil
	}

	done := make(chan error, 1)
	go func() {
		done <- p.proc.Wait()
	}()

	select {
	case <-done:
	case <-time.After(10 * time.Second):
		p.proc.Process.Kill()
		<-done
	}

	p.proc = nil
	return nil
}
