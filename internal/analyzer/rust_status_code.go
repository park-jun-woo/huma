//ff:func feature=analyzer type=data control=sequence
//ff:what Maps Rust StatusCode constants to HTTP status codes for Actix-web
package analyzer

// rustStatusCode maps Rust StatusCode enum constants to HTTP status codes.
var rustStatusCode = map[string]int{
	"OK":                    200,
	"CREATED":               201,
	"ACCEPTED":              202,
	"NO_CONTENT":            204,
	"MOVED_PERMANENTLY":     301,
	"FOUND":                 302,
	"BAD_REQUEST":           400,
	"UNAUTHORIZED":          401,
	"FORBIDDEN":             403,
	"NOT_FOUND":             404,
	"METHOD_NOT_ALLOWED":    405,
	"NOT_ACCEPTABLE":        406,
	"CONFLICT":              409,
	"GONE":                  410,
	"PRECONDITION_FAILED":   412,
	"PAYLOAD_TOO_LARGE":     413,
	"UNSUPPORTED_MEDIA_TYPE": 415,
	"UNPROCESSABLE_ENTITY":  422,
	"TOO_MANY_REQUESTS":     429,
	"INTERNAL_SERVER_ERROR": 500,
	"SERVICE_UNAVAILABLE":   503,
}
