//ff:func feature=source type=reader control=sequence
//ff:what Scans lines to collect the handler function body from start to next top-level func
package source

import (
	"bufio"
	"os"
	"regexp"
)

func collectHandler(f *os.File, pattern *regexp.Regexp) ([]string, int, error) {
	scanner := bufio.NewScanner(f)
	startLine, err := seekToHandler(scanner, pattern)
	if err != nil {
		return nil, 0, err
	}
	if startLine == 0 {
		return nil, 0, nil
	}

	lines, err := readUntilNextFunc(scanner)
	if err != nil {
		return nil, 0, err
	}

	return lines, startLine, nil
}
