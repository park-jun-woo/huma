//ff:func feature=scan type=engine control=iteration dimension=1
//ff:what Finds the file and line of a handler definition by name across candidate source files
package scanner

import "regexp"

// findHandlerDef scans candidate files for a definition of handlerName and
// returns the first matching file and 1-based line number.
func findHandlerDef(files []string, handlerName string) (string, int) {
	defRe := regexp.MustCompile(`\b` + regexp.QuoteMeta(handlerName) + `\s*(\(|=|:)`)
	for _, f := range files {
		if line := scanForPattern(f, defRe); line > 0 {
			return f, line
		}
	}
	return "", 0
}
