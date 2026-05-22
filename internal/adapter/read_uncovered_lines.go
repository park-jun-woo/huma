//ff:func feature=adapter type=reader control=iteration dimension=1
//ff:what Reads source file and returns UncoveredLine entries for lines not covered
package adapter

import (
	"bufio"
	"os"
	"strings"
)

// readUncoveredLines reads the source file and returns UncoveredLine entries
// for lines that are in `total` but not in `covered`.
func readUncoveredLines(file string, covered, total map[int]bool) ([]UncoveredLine, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var uncovered []UncoveredLine
	scanner := bufio.NewScanner(f)
	lineNum := 0
	for scanner.Scan() {
		lineNum++
		if !total[lineNum] || covered[lineNum] {
			continue
		}
		code := strings.TrimRight(scanner.Text(), " \t")
		uncovered = append(uncovered, UncoveredLine{
			File: file,
			Line: lineNum,
			Code: code,
		})
	}
	return uncovered, scanner.Err()
}
