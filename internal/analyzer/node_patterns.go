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
	`new (BadRequest|Unauthorized|Forbidden|NotFound|Conflict|InternalServerError)Exception`,
)

// nestExceptionStatus maps NestJS exception class prefixes to HTTP status codes.
var nestExceptionStatus = map[string]int{
	"BadRequest":          400,
	"Unauthorized":        401,
	"Forbidden":           403,
	"NotFound":            404,
	"Conflict":            409,
	"InternalServerError": 500,
}
