//ff:func feature=scan type=helper control=sequence
//ff:what Strips the template host prefix from a hurl URL and converts template variables to colon-prefixed params
package cmd

import "strings"

func stripHurlURL(raw string) string {
	url := raw
	if idx := strings.Index(url, "}}/"); idx >= 0 {
		url = url[idx+2:]
	}
	return templateVarRe.ReplaceAllString(url, ":$1")
}
