//ff:func feature=source type=reader control=iteration dimension=1
//ff:what Reads lines from the scanner until a top-level func definition or EOF
package source

import (
	"bufio"
	"fmt"
	"strings"
)

func readUntilNextFunc(scanner *bufio.Scanner) ([]string, error) {
	var lines []string
	lines = append(lines, scanner.Text())

	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "func ") {
			break
		}
		lines = append(lines, line)
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scan source file: %w", err)
	}
	return lines, nil
}
