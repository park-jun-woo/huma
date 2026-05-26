//ff:func feature=scan type=helper control=iteration dimension=1
//ff:what Finds the best matching expected hurl filepath using keyword overlap scoring between filename and endpoint paths
package cmd

import (
	"path/filepath"
	"strings"

	"github.com/park-jun-woo/huma/internal/runner"
	"github.com/park-jun-woo/huma/internal/scanner"
)

func findKeywordMatch(name string, endpoints []scanner.Endpoint, hurlDir string) string {
	nameKeywords := extractFilenameKeywords(name)
	if len(nameKeywords) == 0 {
		return ""
	}
	nameLower := strings.ToLower(name)
	methodPrefix := strings.SplitN(nameLower, "_", 2)[0]

	bestScore := 0
	bestPath := ""
	for _, ep := range endpoints {
		if strings.ToLower(ep.Method) != methodPrefix {
			continue
		}
		expected := filepath.Base(runner.HurlFileName(&ep, hurlDir))
		if expected == name {
			continue
		}
		score := countKeywordOverlap(nameKeywords, extractPathKeywords(ep.Path))
		if score <= bestScore {
			continue
		}
		bestScore = score
		bestPath = filepath.Join(hurlDir, expected)
	}
	return bestPath
}
