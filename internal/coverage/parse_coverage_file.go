//ff:func feature=coverage type=parser control=iteration dimension=1
//ff:what Reads a Go coverage.out file and returns all coverage blocks
package coverage

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// ParseCoverageFile reads a Go coverage.out file and returns all coverage blocks.
// Format: file:startLine.startCol,endLine.endCol numStatements count
func ParseCoverageFile(path string) ([]Block, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open coverage file: %w", err)
	}
	defer f.Close()

	var blocks []Block
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "mode:") {
			continue
		}

		b, err := parseLine(line)
		if err != nil {
			continue
		}
		blocks = append(blocks, b)
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scan coverage file: %w", err)
	}
	return blocks, nil
}
