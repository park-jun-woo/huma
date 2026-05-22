//ff:func feature=scan type=parser control=iteration dimension=1
//ff:what Scans a single Go file for gin route registrations and extracts endpoints
package scanner

import (
	"bufio"
	"os"
	"regexp"
)

var ginRoutePattern = regexp.MustCompile(
	`\.\s*(GET|POST|PUT|PATCH|DELETE|HEAD|OPTIONS)\s*\(\s*"([^"]+)"`,
)

func scanFile(path string) ([]Endpoint, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var endpoints []Endpoint
	s := bufio.NewScanner(f)
	lineNum := 0

	for s.Scan() {
		lineNum++
		ep := parseRoute(s.Text(), path, lineNum)
		if ep == nil {
			continue
		}
		endpoints = append(endpoints, *ep)
	}

	return endpoints, s.Err()
}
