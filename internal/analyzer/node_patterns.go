//ff:func feature=analyzer type=helper control=sequence
//ff:what Provides compiled regex patterns and NestJS exception mapping for Node.js response status detection
package analyzer

import "regexp"

var nodePatterns = []*regexp.Regexp{
	// .status(201) (Express, Fastify)
	regexp.MustCompile(`\.status\((\d+)\)`),
	// .sendStatus(204) (Express)
	regexp.MustCompile(`\.sendStatus\((\d+)\)`),
	// .code(201) (Fastify)
	regexp.MustCompile(`\.code\((\d+)\)`),
	// @HttpCode(201) (NestJS decorator)
	regexp.MustCompile(`@HttpCode\((\d+)\)`),
	// new HttpException("...", 400) or new HttpError("...", 400)
	regexp.MustCompile(`new Http(?:Exception|Error)\([^,]*,\s*(\d+)\)`),
}

// nestExceptionPattern matches NestJS built-in exception classes.
var nestExceptionPattern = regexp.MustCompile(
	`new (BadRequest|Unauthorized|Forbidden|NotFound|MethodNotAllowed|` +
		`NotAcceptable|RequestTimeout|Conflict|Gone|PreconditionFailed|` +
		`PayloadTooLarge|UnsupportedMediaType|UnprocessableEntity|` +
		`InternalServerError)Exception`,
)

// nestExceptionStatus maps NestJS exception class prefixes to HTTP status codes.
var nestExceptionStatus = map[string]int{
	"BadRequest":          400,
	"Unauthorized":        401,
	"Forbidden":           403,
	"NotFound":            404,
	"MethodNotAllowed":    405,
	"NotAcceptable":       406,
	"RequestTimeout":      408,
	"Conflict":            409,
	"Gone":                410,
	"PreconditionFailed":  412,
	"PayloadTooLarge":     413,
	"UnsupportedMediaType": 415,
	"UnprocessableEntity": 422,
	"InternalServerError": 500,
}

// httpStatusEnumPattern matches NestJS HttpStatus enum references.
var httpStatusEnumPattern = regexp.MustCompile(`HttpStatus\.([A-Z_]+)`)

// httpStatusEnum maps HttpStatus enum members to HTTP status codes.
var httpStatusEnum = map[string]int{
	"OK":                     200,
	"CREATED":                201,
	"ACCEPTED":               202,
	"NO_CONTENT":             204,
	"BAD_REQUEST":            400,
	"UNAUTHORIZED":           401,
	"FORBIDDEN":              403,
	"NOT_FOUND":              404,
	"METHOD_NOT_ALLOWED":     405,
	"NOT_ACCEPTABLE":         406,
	"REQUEST_TIMEOUT":        408,
	"CONFLICT":               409,
	"GONE":                   410,
	"PRECONDITION_FAILED":    412,
	"PAYLOAD_TOO_LARGE":      413,
	"UNSUPPORTED_MEDIA_TYPE": 415,
	"UNPROCESSABLE_ENTITY":   422,
}

// apiResponsePattern matches @ApiResponse({ status: <code> }) decorators.
var apiResponsePattern = regexp.MustCompile(`@ApiResponse\s*\(\s*\{[^}]*status:\s*(\d+)`)

// apiShorthandPattern matches NestJS Swagger shorthand decorators.
var apiShorthandPattern = regexp.MustCompile(`@Api(Ok|Created|Accepted|NoContent|BadRequest|Unauthorized|Forbidden|NotFound|Conflict|UnprocessableEntity)Response`)

// apiShorthandStatus maps Swagger shorthand decorator names to HTTP status codes.
var apiShorthandStatus = map[string]int{
	"Ok":                  200,
	"Created":             201,
	"Accepted":            202,
	"NoContent":           204,
	"BadRequest":          400,
	"Unauthorized":        401,
	"Forbidden":           403,
	"NotFound":            404,
	"Conflict":            409,
	"UnprocessableEntity": 422,
}
