//ff:func feature=adapter type=engine control=sequence
//ff:what Polls the ready URL until the Node.js server returns 200 or times out
package adapter

import (
	"net/http"
	"time"
)

// WaitReady polls the ready URL until it returns 200 or times out (30s).
func (a *NodeAdapter) WaitReady() error {
	if a.cfg.Ready == "" {
		time.Sleep(2 * time.Second)
		return nil
	}

	client := &http.Client{Timeout: 2 * time.Second}
	deadline := time.After(30 * time.Second)
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	return awaitReady(client, a.cfg.Ready, deadline, ticker.C)
}
