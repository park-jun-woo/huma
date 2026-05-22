//ff:func feature=adapter type=engine control=iteration dimension=1
//ff:what Blocks on tick/deadline channels until the server responds 200 or timeout
package adapter

import (
	"net/http"
	"time"
)

func awaitReady(client *http.Client, url string, deadline <-chan time.Time, tick <-chan time.Time) error {
	for {
		err := waitForEvent(client, url, deadline, tick)
		if err == errRetry {
			continue
		}
		return err
	}
}
