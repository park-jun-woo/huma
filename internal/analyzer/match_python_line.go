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
	if flaskRedirectImplicit.MatchString(line) {
		return 302
	}
	// Django HttpResponseXxx classes
	if m := djangoResponseClassPattern.FindStringSubmatch(line); m != nil {
		if code, ok := djangoResponseClassStatus[m[1]]; ok {
			return code
		}
	}
	// Django exception raises
	if m := djangoExceptionPattern.FindStringSubmatch(line); m != nil {
		if code, ok := djangoExceptionStatus[m[1]]; ok {
			return code
		}
	}
	// DRF exception raises
	if m := drfExceptionPattern.FindStringSubmatch(line); m != nil {
		if code, ok := drfExceptionStatus[m[1]]; ok {
			return code
		}
	}
	return 0
}
