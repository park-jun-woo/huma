//ff:func feature=scan type=parser control=iteration dimension=1
//ff:what Reads a hurl file and extracts the HTTP method and normalized path from the first request line
package cmd

import (
	"bufio"
	"os"
	"regexp"
	"strings"
)

var httpMethods = []string{"GET", "POST", "PUT", "DELETE", "PATCH"}

var templateVarRe = regexp.MustCompile(`\{\{(\w+)\}\}`)

func parseHurlRequest(hurlPath string) (method, path string) {
	f, err := os.Open(hurlPath)
	if err != nil {
		return "", ""
	}
	defer f.Close()

	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		m, url, ok := matchLineMethod(line)
		if !ok {
			return "", ""
		}
		return m, stripHurlURL(url)
	}
	return "", ""
}
