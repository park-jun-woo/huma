//ff:func feature=analyzer type=helper control=iteration dimension=1
//ff:what Matches a single source line against Python response patterns and returns matched status code
package analyzer

import "strconv"

// matchPythonLine tries each Python regex pattern against a line and returns
// the matched HTTP status code. Returns 0 if no pattern matches.
func matchPythonLine(line string) int {
	for _, pat := range pythonPatterns {
		m := pat.FindStringSubmatch(line)
		if m == nil {
			continue
		}
		code, err := strconv.Atoi(m[1])
		if err != nil {
			continue
		}
		return code
	}
	return 0
}
