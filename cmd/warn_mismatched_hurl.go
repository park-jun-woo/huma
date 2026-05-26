//ff:func feature=scan type=helper control=iteration dimension=1
//ff:what Checks hurl directory for files that do not match expected naming convention and prints warnings
package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/park-jun-woo/huma/internal/runner"
	"github.com/park-jun-woo/huma/internal/scanner"
)

func warnMismatchedHurlFiles(hurlDir string, endpoints []scanner.Endpoint) {
	entries, err := os.ReadDir(hurlDir)
	if err != nil {
		return
	}

	expectedNames := make(map[string]bool)
	for _, ep := range endpoints {
		name := filepath.Base(runner.HurlFileName(&ep, hurlDir))
		expectedNames[name] = true
	}

	type mismatch struct {
		actual   string
		expected string
	}

	var mismatches []mismatch

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".hurl") {
			continue
		}
		if !expectedNames[entry.Name()] {
			expected := findLikelyMatch(entry.Name(), endpoints, hurlDir)
			mismatches = append(mismatches, mismatch{
				actual:   filepath.Join(hurlDir, entry.Name()),
				expected: expected,
			})
		}
	}

	if len(mismatches) == 0 {
		return
	}

	fmt.Printf("\n# WARNING — %d existing .hurl files don't match naming convention\n\n", len(mismatches))
	for _, m := range mismatches {
		if m.expected != "" {
			fmt.Printf("  %-40s → expected: %s\n", m.actual, m.expected)
		} else {
			fmt.Printf("  %-40s → no matching endpoint\n", m.actual)
		}
	}
	fmt.Println("\nRename to match, or huma will treat these endpoints as TODO.")
}
