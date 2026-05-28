//ff:func feature=analyzer type=data control=sequence
//ff:what Maps Java HttpStatus enum constants to HTTP status codes for Spring and JAX-RS
package analyzer

// javaHttpStatus maps HttpStatus enum members to HTTP status codes.
// Used by both Spring Boot (HttpStatus.CREATED) and JAX-RS (Response.Status.CREATED).
var javaHttpStatus = map[string]int{
	"OK":                       200,
	"CREATED":                  201,
	"ACCEPTED":                 202,
	"NO_CONTENT":               204,
	"MOVED_PERMANENTLY":        301,
	"FOUND":                    302,
	"BAD_REQUEST":              400,
	"UNAUTHORIZED":             401,
	"FORBIDDEN":                403,
	"NOT_FOUND":                404,
	"METHOD_NOT_ALLOWED":       405,
	"NOT_ACCEPTABLE":           406,
	"CONFLICT":                 409,
	"GONE":                     410,
	"PRECONDITION_FAILED":      412,
	"REQUEST_ENTITY_TOO_LARGE": 413,
	"UNSUPPORTED_MEDIA_TYPE":   415,
	"UNPROCESSABLE_ENTITY":     422,
	"TOO_MANY_REQUESTS":        429,
	"INTERNAL_SERVER_ERROR":    500,
	"SERVICE_UNAVAILABLE":      503,
}
