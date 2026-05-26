//ff:func feature=analyzer type=helper control=iteration dimension=1
//ff:what Matches a single source line against Node.js response patterns and returns matched status code
package analyzer

import "strconv"

// matchNodeLine tries each Node.js regex pattern against a line, including
// NestJS exception class mapping. Returns the matched HTTP status code, or 0.
func matchNodeLine(line string) int {
	for _, pat := range nodePatterns {
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

	m := nestExceptionPattern.FindStringSubmatch(line)
	if m == nil {
		return 0
	}
	code, ok := nestExceptionStatus[m[1]]
	if !ok {
		return 0
	}
	return code
}
