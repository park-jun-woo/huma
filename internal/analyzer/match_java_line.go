//ff:func feature=analyzer type=helper control=iteration dimension=1
//ff:what Matches a single source line against Java response patterns and returns matched status code
package analyzer

import "strconv"

// matchJavaLine tries each Java regex pattern against a line, including
// Spring Boot, Quarkus/JAX-RS patterns, factory methods, and exception mapping.
// Returns the matched HTTP status code, or 0.
func matchJavaLine(line string) int {
	// Spring Boot patterns with numeric status
	for _, pat := range springPatterns {
		m := pat.FindStringSubmatch(line)
		if m == nil {
			continue
		}
		// Check if it's a numeric capture
		if code, err := strconv.Atoi(m[1]); err == nil {
			return code
		}
		// Check if it's an HttpStatus enum name
		if code, ok := javaHttpStatus[m[1]]; ok {
			return code
		}
	}

	// Spring Boot factory methods
	if m := springFactoryMethod.FindStringSubmatch(line); m != nil {
		if code, ok := springFactoryStatus[m[1]]; ok {
			return code
		}
	}

	// Quarkus/JAX-RS patterns with numeric status
	for _, pat := range quarkusPatterns {
		m := pat.FindStringSubmatch(line)
		if m == nil {
			continue
		}
		if code, err := strconv.Atoi(m[1]); err == nil {
			return code
		}
		if code, ok := javaHttpStatus[m[1]]; ok {
			return code
		}
	}

	// Quarkus factory methods
	if m := quarkusFactoryMethod.FindStringSubmatch(line); m != nil {
		if code, ok := quarkusFactoryStatus[m[1]]; ok {
			return code
		}
	}

	// JAX-RS exception classes
	if m := jaxrsExceptionPattern.FindStringSubmatch(line); m != nil {
		if code, ok := jaxrsExceptionStatus[m[1]]; ok {
			return code
		}
	}

	return 0
}
