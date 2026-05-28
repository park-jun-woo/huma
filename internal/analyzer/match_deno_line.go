//ff:func feature=analyzer type=helper control=iteration dimension=1
//ff:what Matches a single source line against Deno response patterns and returns matched status code
package analyzer

import "strconv"

// matchDenoLine tries each Deno regex pattern against a line.
// Returns the matched HTTP status code, or 0.
func matchDenoLine(line string) int {
	for _, pat := range denoPatterns {
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

	if denoImplicitRedirect302.MatchString(line) {
		return 302
	}

	if denoImplicitJson200.MatchString(line) {
		return 200
	}

	if denoImplicitResponse200.MatchString(line) {
		return 200
	}

	return 0
}
