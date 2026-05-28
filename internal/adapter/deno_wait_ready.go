//ff:func feature=adapter type=engine control=sequence
//ff:what Polls the ready URL until the Deno server returns 200 or times out
package adapter

import (
	"net/http"
	"time"
)

// WaitReady polls the ready URL until it returns 200 or times out (30s).
func (d *DenoAdapter) WaitReady() error {
	if d.cfg.Ready == "" {
		time.Sleep(2 * time.Second)
		return nil
	}

	client := &http.Client{Timeout: 2 * time.Second}
	deadline := time.After(30 * time.Second)
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	return awaitReady(client, d.baseURL+d.cfg.Ready, deadline, ticker.C)
}
