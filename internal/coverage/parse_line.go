//ff:func feature=coverage type=parser control=sequence
//ff:what Parses a single coverage line into a Block struct
package coverage

import (
	"fmt"
	"strconv"
	"strings"
)

// parseLine parses a single coverage line.
// Format: github.com/example/handler.go:41.2,43.4 1 0
func parseLine(line string) (Block, error) {
	lastSpace := strings.LastIndex(line, " ")
	if lastSpace < 0 {
		return Block{}, fmt.Errorf("no space found")
	}
	countStr := line[lastSpace+1:]
	rest := line[:lastSpace]

	lastSpace2 := strings.LastIndex(rest, " ")
	if lastSpace2 < 0 {
		return Block{}, fmt.Errorf("no second space found")
	}
	posStr := rest[:lastSpace2]

	count, err := strconv.Atoi(countStr)
	if err != nil {
		return Block{}, fmt.Errorf("parse count: %w", err)
	}

	return parsePosition(posStr, count)
}
