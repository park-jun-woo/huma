//ff:func feature=analyzer type=parser control=iteration dimension=1
//ff:what Analyzes Rust source using regex to extract response status codes from Actix-web handlers
package analyzer

import (
	"fmt"
	"os"
	"strings"
)

// Analyze reads a Rust source file and extracts response branches from the
// specified line range using regex patterns for Actix-web.
func (r *RustAnalyzer) Analyze(file string, handlerName string, startLine, endLine int) ([]ResponseBranch, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", file, err)
	}

	allLines := strings.Split(string(data), "\n")

	start := startLine
	if start <= 0 {
		start = 1
	}
	end := endLine
	if end <= 0 || end > len(allLines) {
		end = len(allLines)
	}

	var branches []ResponseBranch
	seen := map[int]bool{}
	for i := start - 1; i < end; i++ {
		code := matchRustLine(allLines[i])
		if code == 0 {
			continue
		}
		lineNum := i + 1
		key := lineNum*1000 + code
		if seen[key] {
			continue
		}
		seen[key] = true
		branches = append(branches, ResponseBranch{
			Status: code,
			File:   file,
			Line:   lineNum,
			Code:   strings.TrimSpace(allLines[i]),
		})
	}

	return branches, nil
}
