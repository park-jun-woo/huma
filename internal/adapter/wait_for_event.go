//ff:func feature=adapter type=engine control=sequence
//ff:what Selects on deadline or tick channel, returning nil on success, errRetry to continue, or error on timeout
package adapter

import (
	"fmt"
	"net/http"
	"time"
)

var errRetry = fmt.Errorf("retry")

func waitForEvent(client *http.Client, url string, deadline <-chan time.Time, tick <-chan time.Time) error {
	select {
	case <-deadline:
		return fmt.Errorf("server not ready after 30s (url: %s)", url)
	case <-tick:
		if probeOK(client, url) {
			return nil
		}
		return errRetry
	}
}
