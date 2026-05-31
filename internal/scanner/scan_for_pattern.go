//ff:func feature=scan type=reader control=iteration dimension=1
//ff:what Scans a file line by line and returns the first 1-based line matching a pattern (0 if none)
package scanner

import (
	"bufio"
	"os"
	"regexp"
)

// scanForPattern returns the first 1-based line number in file matching re,
// or 0 if the file cannot be read or no line matches.
func scanForPattern(file string, re *regexp.Regexp) int {
	fh, err := os.Open(file)
	if err != nil {
		return 0
	}
	defer fh.Close()
	lineNum := 0
	sc := bufio.NewScanner(fh)
	sc.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	for sc.Scan() {
		lineNum++
		if re.MatchString(sc.Text()) {
			return lineNum
		}
	}
	return 0
}
