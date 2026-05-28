//ff:func feature=analyzer type=helper control=iteration dimension=1
//ff:what Matches a single source line against ASP.NET Core response patterns and returns matched status code
package analyzer

import "strconv"

// matchDotnetLine tries each ASP.NET Core regex pattern against a line.
// Returns the matched HTTP status code, or 0.
func matchDotnetLine(line string) int {
	for _, pat := range dotnetNumericPatterns {
		m := pat.FindStringSubmatch(line)
		if m == nil {
			continue
		}
		if code, err := strconv.Atoi(m[1]); err == nil {
			return code
		}
	}

	// Controller helper methods: return Ok(), return NotFound(), etc.
	if m := aspnetControllerMethodPattern.FindStringSubmatch(line); m != nil {
		if code, ok := aspnetMethodStatus[m[1]]; ok {
			return code
		}
	}

	// Minimal API: Results.Ok(), TypedResults.NotFound(), etc.
	if m := aspnetMinimalMethodPattern.FindStringSubmatch(line); m != nil {
		if code, ok := aspnetMethodStatus[m[1]]; ok {
			return code
		}
	}

	// return RedirectPermanent(...) — 301
	if aspnetRedirectPermanentPattern.MatchString(line) {
		return 301
	}

	// return Redirect/RedirectToAction/etc — 302
	if aspnetRedirectPattern.MatchString(line) {
		return 302
	}

	// Results.Redirect(url) — 302
	if aspnetMinimalRedirect.MatchString(line) {
		return 302
	}

	return 0
}
