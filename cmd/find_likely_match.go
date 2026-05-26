//ff:func feature=scan type=helper control=sequence
//ff:what Finds the best matching expected hurl filepath for a mismatched file by parsing hurl content first then falling back to keyword overlap scoring
package cmd

import (
	"path/filepath"

	"github.com/park-jun-woo/huma/internal/scanner"
)

func findLikelyMatch(name string, endpoints []scanner.Endpoint, hurlDir string) string {
	hurlPath := filepath.Join(hurlDir, name)
	method, path := parseHurlRequest(hurlPath)
	if method != "" && path != "" {
		if match := findContentMatch(method, path, name, endpoints, hurlDir); match != "" {
			return match
		}
	}
	return findKeywordMatch(name, endpoints, hurlDir)
}
