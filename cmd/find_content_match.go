//ff:func feature=scan type=helper control=iteration dimension=1
//ff:what Finds a matching endpoint by comparing parsed hurl method and normalized path pattern against all endpoints
package cmd

import (
	"path/filepath"
	"strings"

	"github.com/park-jun-woo/huma/internal/runner"
	"github.com/park-jun-woo/huma/internal/scanner"
)

func findContentMatch(method, path, name string, endpoints []scanner.Endpoint, hurlDir string) string {
	normalized := normalizePathPattern(path)
	for _, ep := range endpoints {
		if !strings.EqualFold(ep.Method, method) {
			continue
		}
		if normalizePathPattern(ep.Path) != normalized {
			continue
		}
		expected := filepath.Base(runner.HurlFileName(&ep, hurlDir))
		if expected == name {
			continue
		}
		return filepath.Join(hurlDir, expected)
	}
	return ""
}
