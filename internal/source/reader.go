package source

import (
	"bufio"
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

	// Match func definitions: "func HandlerName(" or "func (r *Type) HandlerName("
	pattern := regexp.MustCompile(`^func\s+(\([^)]*\)\s*)?` + regexp.QuoteMeta(handlerName) + `\s*\(`)

	var lines []string
	scanner := bufio.NewScanner(f)
	lineNum := 0
	startLine := 0
	collecting := false

	for scanner.Scan() {
		lineNum++
		line := scanner.Text()

		if collecting {
			// Stop at next top-level func definition (column 0)
			if strings.HasPrefix(line, "func ") && lineNum > startLine {
				break
			}
			lines = append(lines, line)
			continue
		}

		if pattern.MatchString(line) {
			startLine = lineNum
			collecting = true
			lines = append(lines, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return "", 0, 0, fmt.Errorf("scan source file: %w", err)
	}

	if !collecting {
		return "", 0, 0, fmt.Errorf("handler %q not found in %s", handlerName, file)
	}

	// Trim trailing blank lines and comments that belong to the next function
	for len(lines) > 0 {
		trimmed := strings.TrimSpace(lines[len(lines)-1])
		if trimmed == "" || strings.HasPrefix(trimmed, "//") {
			lines = lines[:len(lines)-1]
			continue
		}
		break
	}

	endLine := startLine + len(lines) - 1
	return strings.Join(lines, "\n"), startLine, endLine, nil
}
