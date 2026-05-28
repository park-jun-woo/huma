//ff:func feature=analyzer type=helper control=iteration dimension=1
//ff:what Matches a single source line against Node.js response patterns and returns matched status code
package analyzer

import "strconv"

// matchNodeLine tries each Node.js regex pattern against a line, including
// NestJS exception class mapping, HttpStatus enum, and @ApiResponse decorators.
// Returns the matched HTTP status code, or 0.
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

	if m := nestExceptionPattern.FindStringSubmatch(line); m != nil {
		if code, ok := nestExceptionStatus[m[1]]; ok {
			return code
		}
	}

	if m := httpStatusEnumPattern.FindStringSubmatch(line); m != nil {
		if code, ok := httpStatusEnum[m[1]]; ok {
			return code
		}
	}

	if m := apiResponsePattern.FindStringSubmatch(line); m != nil {
		code, err := strconv.Atoi(m[1])
		if err == nil {
			return code
		}
	}

	if m := apiShorthandPattern.FindStringSubmatch(line); m != nil {
		if code, ok := apiShorthandStatus[m[1]]; ok {
			return code
		}
	}

	if expressRedirectImplicit.MatchString(line) {
		return 302
	}

	if expressImplicit200.MatchString(line) {
		return 200
	}

	return 0
}
