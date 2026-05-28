//ff:func feature=analyzer type=helper control=sequence
//ff:what Provides compiled regex patterns and exception mapping for Quarkus/JAX-RS response status detection
package analyzer

import "regexp"

// quarkusPatterns contains Quarkus/JAX-RS patterns with a single capture group.
// Each pattern must capture either a numeric status code or a Response.Status enum name.
var quarkusPatterns = []*regexp.Regexp{
	// Response.status(201)
	regexp.MustCompile(`Response\.status\((\d+)\)`),
	// Response.status(Response.Status.CREATED)
	regexp.MustCompile(`Response\.status\(Response\.Status\.([A-Z_]+)\)`),
	// throw new WebApplicationException(404)
	regexp.MustCompile(`new WebApplicationException\((\d+)\)`),
	// throw new WebApplicationException(Response.Status.NOT_FOUND)
	regexp.MustCompile(`new WebApplicationException\(Response\.Status\.([A-Z_]+)\)`),
}

// jaxrsExceptionPattern matches JAX-RS built-in exception classes.
var jaxrsExceptionPattern = regexp.MustCompile(
	`new (NotFoundException|BadRequestException|ForbiddenException|` +
		`NotAuthorizedException|NotAllowedException|NotAcceptableException|` +
		`InternalServerErrorException|ServiceUnavailableException)\(`,
)

// jaxrsExceptionStatus maps JAX-RS exception class names to HTTP status codes.
var jaxrsExceptionStatus = map[string]int{
	"NotFoundException":            404,
	"BadRequestException":          400,
	"ForbiddenException":           403,
	"NotAuthorizedException":       401,
	"NotAllowedException":          405,
	"NotAcceptableException":       406,
	"InternalServerErrorException": 500,
	"ServiceUnavailableException":  503,
}

// quarkusFactoryMethod matches JAX-RS Response factory method calls.
var quarkusFactoryMethod = regexp.MustCompile(`Response\.(ok|created|accepted|noContent|serverError)\(`)

// quarkusFactoryStatus maps JAX-RS Response factory methods to HTTP status codes.
var quarkusFactoryStatus = map[string]int{
	"ok":          200,
	"created":     201,
	"accepted":    202,
	"noContent":   204,
	"serverError": 500,
}
