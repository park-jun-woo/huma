//ff:func feature=adapter type=helper control=sequence
//ff:what Sends a GET request with a short timeout and returns true if the server responds with 200
package adapter

import (
	"net/http"
	"time"
)

func ProbeCheck(url string) bool {
	client := &http.Client{Timeout: 2 * time.Second}
	return probeOK(client, url)
}
