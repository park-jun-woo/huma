//ff:func feature=analyzer type=helper control=sequence
//ff:what Provides compiled regex patterns for detecting response status codes in ASP.NET Core C# source
package analyzer

import "regexp"

// aspnetMethodStatus maps Controller helper method names to HTTP status codes.
var aspnetMethodStatus = map[string]int{
	"Ok":                  200,
	"Created":             201,
	"CreatedAtAction":     201,
	"CreatedAtRoute":      201,
	"Accepted":            202,
	"NoContent":           204,
	"BadRequest":          400,
	"Unauthorized":        401,
	"Forbid":              403,
	"NotFound":            404,
	"Conflict":            409,
	"UnprocessableEntity": 422,
}

// aspnetControllerMethodPattern extracts the method name from return statements.
var aspnetControllerMethodPattern = regexp.MustCompile(`return\s+(Ok|Created(?:AtAction|AtRoute)?|Accepted|NoContent|BadRequest|Unauthorized|Forbid|NotFound|Conflict|UnprocessableEntity)\(`)

// aspnetStatusCodePattern matches return StatusCode(201) or StatusCode(StatusCodes.Status201...).
var aspnetStatusCodeNumeric = regexp.MustCompile(`return\s+StatusCode\((\d+)`)
var aspnetStatusCodeEnum = regexp.MustCompile(`StatusCode\(StatusCodes\.Status(\d+)`)

// aspnetRedirectPattern matches return Redirect/RedirectToAction/etc — implicit 302.
var aspnetRedirectPattern = regexp.MustCompile(`return\s+Redirect(?:ToAction|ToRoute|ToPage)?\(`)

// aspnetRedirectPermanentPattern matches return RedirectPermanent — implicit 301.
var aspnetRedirectPermanentPattern = regexp.MustCompile(`return\s+RedirectPermanent\(`)

// aspnetMinimalPatterns matches ASP.NET Core Minimal API patterns.
var aspnetMinimalMethodPattern = regexp.MustCompile(`(?:Results|TypedResults)\.(Ok|Created|Accepted|NoContent|BadRequest|Unauthorized|Forbid|NotFound|Conflict|UnprocessableEntity)\(`)

// aspnetMinimalStatusCode matches Results.StatusCode(201).
var aspnetMinimalStatusCode = regexp.MustCompile(`Results\.StatusCode\((\d+)`)

// aspnetMinimalRedirect matches Results.Redirect(url) — implicit 302.
var aspnetMinimalRedirect = regexp.MustCompile(`Results\.Redirect\(`)

// producesPattern matches [ProducesResponseType(StatusCodes.Status201Created)].
var producesPattern = regexp.MustCompile(`\[ProducesResponseType\(.*Status(\d+)`)

// producesNumericPattern matches [ProducesResponseType(201)].
var producesNumericPattern = regexp.MustCompile(`\[ProducesResponseType\((\d+)\)`)

// producesTypePattern matches [ProducesResponseType(typeof(ErrorResponse), 400)].
var producesTypePattern = regexp.MustCompile(`\[ProducesResponseType\(typeof\([^)]+\),\s*(\d+)\)`)

// dotnetNumericPatterns is the ordered list of patterns whose first capture group
// is a numeric status code. Used by matchDotnetLine's for-loop.
var dotnetNumericPatterns = []*regexp.Regexp{
	aspnetStatusCodeNumeric,
	aspnetStatusCodeEnum,
	aspnetMinimalStatusCode,
	producesTypePattern,
	producesPattern,
	producesNumericPattern,
}
