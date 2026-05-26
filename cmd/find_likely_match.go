//ff:func feature=scan type=helper control=iteration dimension=1
//ff:what Finds the most likely expected hurl filename for a mismatched existing file
package cmd

import (
	"path/filepath"
	"strings"

	"github.com/park-jun-woo/huma/internal/runner"
	"github.com/park-jun-woo/huma/internal/scanner"
)

func findLikelyMatch(name string, endpoints []scanner.Endpoint, hurlDir string) string {
	nameLower := strings.ToLower(name)
	for _, ep := range endpoints {
		expected := filepath.Base(runner.HurlFileName(&ep, hurlDir))
		if expected == name {
			continue
		}
		methodPrefix := strings.ToLower(ep.Method) + "_"
		if strings.HasPrefix(nameLower, methodPrefix) {
			return filepath.Join(hurlDir, expected)
		}
	}
	return ""
}
