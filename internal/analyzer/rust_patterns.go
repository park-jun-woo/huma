//ff:func feature=analyzer type=helper control=sequence
//ff:what Provides compiled regex patterns and factory method mapping for Actix-web response status detection
package analyzer

import "regexp"

// actixPatterns contains Actix-web patterns with a single capture group.
// Each pattern must capture either a StatusCode constant name or a numeric status.
var actixPatterns = []*regexp.Regexp{
	// HttpResponse::build(StatusCode::CREATED)
	regexp.MustCompile(`HttpResponse::build\(StatusCode::([A-Z_]+)\)`),
	// HttpResponseBuilder::new(StatusCode::CREATED)
	regexp.MustCompile(`HttpResponseBuilder::new\(StatusCode::([A-Z_]+)\)`),
}

// actixImplicitPatterns match HttpResponse factory methods with fixed status codes.
// These patterns have no capture group — the status is determined by the factory name.
var actixImplicitPatterns = []*regexp.Regexp{
	regexp.MustCompile(`HttpResponse::Ok\(`),
	regexp.MustCompile(`HttpResponse::Created\(`),
	regexp.MustCompile(`HttpResponse::Accepted\(`),
	regexp.MustCompile(`HttpResponse::NoContent\(`),
	regexp.MustCompile(`HttpResponse::MovedPermanently\(`),
	regexp.MustCompile(`HttpResponse::Found\(`),
	regexp.MustCompile(`HttpResponse::BadRequest\(`),
	regexp.MustCompile(`HttpResponse::Unauthorized\(`),
	regexp.MustCompile(`HttpResponse::Forbidden\(`),
	regexp.MustCompile(`HttpResponse::NotFound\(`),
	regexp.MustCompile(`HttpResponse::MethodNotAllowed\(`),
	regexp.MustCompile(`HttpResponse::Conflict\(`),
	regexp.MustCompile(`HttpResponse::Gone\(`),
	regexp.MustCompile(`HttpResponse::UnprocessableEntity\(`),
	regexp.MustCompile(`HttpResponse::InternalServerError\(`),
	regexp.MustCompile(`HttpResponse::ServiceUnavailable\(`),
}

// actixFactoryStatus maps HttpResponse factory method names to HTTP status codes.
var actixFactoryStatus = map[string]int{
	"Ok":                  200,
	"Created":             201,
	"Accepted":            202,
	"NoContent":           204,
	"MovedPermanently":    301,
	"Found":               302,
	"BadRequest":          400,
	"Unauthorized":        401,
	"Forbidden":           403,
	"NotFound":            404,
	"MethodNotAllowed":    405,
	"Conflict":            409,
	"Gone":                410,
	"UnprocessableEntity": 422,
	"InternalServerError": 500,
	"ServiceUnavailable":  503,
}

// actixImplicitFactory matches HttpResponse::<FactoryName>( and captures the factory name.
var actixImplicitFactory = regexp.MustCompile(`HttpResponse::([A-Z][a-zA-Z]+)\(`)

// actixErrorPattern matches actix_web::error::Error<Name>() calls.
var actixErrorPattern = regexp.MustCompile(
	`error::Error(NotFound|BadRequest|Unauthorized|Forbidden|` +
		`MethodNotAllowed|Conflict|Gone|InternalServerError|ServiceUnavailable)\(`,
)

// actixErrorStatus maps actix_web error function names to HTTP status codes.
var actixErrorStatus = map[string]int{
	"NotFound":            404,
	"BadRequest":          400,
	"Unauthorized":        401,
	"Forbidden":           403,
	"MethodNotAllowed":    405,
	"Conflict":            409,
	"Gone":                410,
	"InternalServerError": 500,
	"ServiceUnavailable":  503,
}
