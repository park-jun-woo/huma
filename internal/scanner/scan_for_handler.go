//ff:func feature=scan type=reader control=iteration dimension=2
//ff:what Scans one file for a handler definition whose normalized identifier matches a target key
package scanner

import (
	"bufio"
	"os"
)

// scanForHandler returns the 1-based line number in file where a handler whose
// normalized identifier equals normTarget is DEFINED, or 0 if none. It applies
// the language-specific definition patterns for the file's extension; for each
// line it normalizes the captured identifier (group 1) and compares keys, which
// absorbs camelCase <-> PascalCase (§2.3). Patterns are tried in priority order
// and the earliest line with any matching pattern wins.
func scanForHandler(file, normTarget string) int {
	pats := handlerDefPatterns(file)
	if len(pats) == 0 {
		return 0
	}
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
		text := sc.Text()
		for _, re := range pats {
			m := re.FindStringSubmatch(text)
			if m != nil && normalizeSymbol(m[1]) == normTarget {
				return lineNum
			}
		}
	}
	return 0
}
