//ff:func feature=analyzer type=helper control=iteration dimension=1
//ff:what Matches a single source line against PHP/Laravel response patterns and returns matched status code
package analyzer

import "strconv"

// matchPhpLine tries each PHP/Laravel regex pattern against a line.
// Returns the matched HTTP status code, or 0.
func matchPhpLine(line string) int {
	// Explicit status code patterns (controller responses).
	for _, pat := range laravelPatterns {
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

	// API Resource patterns with explicit status code.
	for _, pat := range laravelResourcePatterns {
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

	// Exception throws with mapped status codes.
	if m := laravelExceptionPattern.FindStringSubmatch(line); m != nil {
		if code, ok := laravelExceptionStatus[m[1]]; ok {
			return code
		}
	}

	// Implicit redirect 302.
	if laravelImplicitRedirect302.MatchString(line) {
		return 302
	}

	// Implicit 200 patterns.
	if laravelImplicitJson200.MatchString(line) {
		return 200
	}

	if laravelImplicitResponseJson200.MatchString(line) {
		return 200
	}

	if laravelImplicitResource200.MatchString(line) {
		return 200
	}

	if laravelImplicitCollection200.MatchString(line) {
		return 200
	}

	return 0
}
