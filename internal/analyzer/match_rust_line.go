//ff:func feature=analyzer type=helper control=iteration dimension=1
//ff:what Matches a single source line against Rust/Actix-web response patterns and returns matched status code
package analyzer

// matchRustLine tries each Actix-web regex pattern against a line, including
// HttpResponse factory methods, StatusCode-based patterns, and error function mapping.
// Returns the matched HTTP status code, or 0.
func matchRustLine(line string) int {
	// Patterns with capture group (StatusCode::XXX).
	for _, pat := range actixPatterns {
		m := pat.FindStringSubmatch(line)
		if m == nil {
			continue
		}
		if code, ok := rustStatusCode[m[1]]; ok {
			return code
		}
	}

	// Implicit factory patterns (HttpResponse::Ok(), etc.).
	if m := actixImplicitFactory.FindStringSubmatch(line); m != nil {
		if code, ok := actixFactoryStatus[m[1]]; ok {
			return code
		}
	}

	// Error function patterns (error::ErrorNotFound(), etc.).
	if m := actixErrorPattern.FindStringSubmatch(line); m != nil {
		if code, ok := actixErrorStatus[m[1]]; ok {
			return code
		}
	}

	return 0
}
