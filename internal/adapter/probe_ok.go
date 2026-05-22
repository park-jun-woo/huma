//ff:func feature=adapter type=helper control=sequence
//ff:what Sends a GET request and returns true if the response status is 200
package adapter

import "net/http"

func probeOK(client *http.Client, url string) bool {
	resp, err := client.Get(url)
	if err != nil {
		return false
	}
	resp.Body.Close()
	return resp.StatusCode == http.StatusOK
}
