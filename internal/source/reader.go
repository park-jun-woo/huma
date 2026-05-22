//ff:func feature=source type=reader control=sequence
//ff:what Reads a source file and extracts the handler function body with line range
package source

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

// ReadHandler reads the source file and extracts the handler function body.
// It searches for a function matching handlerName and returns lines from the
// function definition until the next top-level func or EOF.
func ReadHandler(file string, handlerName string) (string, int, int, error) {
	f, err := os.Open(file)
	if err != nil {
		return "", 0, 0, fmt.Errorf("open source file: %w", err)
	}
	defer f.Close()

	pattern := regexp.MustCompile(`^func\s+(\([^)]*\)\s*)?` + regexp.QuoteMeta(handlerName) + `\s*\(`)

	lines, startLine, err := collectHandler(f, pattern)
	if err != nil {
		return "", 0, 0, err
	}

	if startLine == 0 {
		return "", 0, 0, fmt.Errorf("handler %q not found in %s", handlerName, file)
	}

	lines = trimTrailing(lines)
	endLine := startLine + len(lines) - 1
	return strings.Join(lines, "\n"), startLine, endLine, nil
}
