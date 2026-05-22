//ff:func feature=source type=reader control=iteration dimension=1
//ff:what Scans lines until finding a function matching the pattern, returns its line number
package source

import (
	"bufio"
	"fmt"
	"regexp"
)

func seekToHandler(scanner *bufio.Scanner, pattern *regexp.Regexp) (int, error) {
	lineNum := 0
	for scanner.Scan() {
		lineNum++
		if pattern.MatchString(scanner.Text()) {
			return lineNum, nil
		}
	}
	if err := scanner.Err(); err != nil {
		return 0, fmt.Errorf("scan source file: %w", err)
	}
	return 0, nil
}
