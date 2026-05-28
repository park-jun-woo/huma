//ff:func feature=analyzer type=helper control=sequence
//ff:what Provides compiled regex patterns and factory method mapping for Spring Boot response status detection
package analyzer

import "regexp"

// springPatterns contains Spring Boot patterns with a single capture group.
// Each pattern must capture either a numeric status code or an HttpStatus enum name.
var springPatterns = []*regexp.Regexp{
	// ResponseEntity.status(201)
	regexp.MustCompile(`ResponseEntity\.status\((\d+)\)`),
	// ResponseEntity.status(HttpStatus.CREATED)
	regexp.MustCompile(`ResponseEntity\.status\(HttpStatus\.([A-Z_]+)\)`),
	// new ResponseEntity<>(body, HttpStatus.CREATED)
	regexp.MustCompile(`new ResponseEntity<?>?\(.*HttpStatus\.([A-Z_]+)\)`),
	// @ResponseStatus(HttpStatus.CREATED)
	regexp.MustCompile(`@ResponseStatus\((?:value\s*=\s*)?HttpStatus\.([A-Z_]+)\)`),
	// @ResponseStatus(code = HttpStatus.NOT_FOUND)
	regexp.MustCompile(`@ResponseStatus\(code\s*=\s*HttpStatus\.([A-Z_]+)\)`),
	// throw new ResponseStatusException(HttpStatus.NOT_FOUND)
	regexp.MustCompile(`new ResponseStatusException\(HttpStatus\.([A-Z_]+)`),
}

// springFactoryMethod matches ResponseEntity factory method calls.
var springFactoryMethod = regexp.MustCompile(`ResponseEntity\.(ok|created|accepted|noContent|badRequest|notFound|unprocessableEntity)\(`)

// springFactoryStatus maps ResponseEntity factory methods to HTTP status codes.
var springFactoryStatus = map[string]int{
	"ok":                  200,
	"created":             201,
	"accepted":            202,
	"noContent":           204,
	"badRequest":          400,
	"notFound":            404,
	"unprocessableEntity": 422,
}
