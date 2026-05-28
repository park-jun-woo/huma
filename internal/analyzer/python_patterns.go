//ff:func feature=analyzer type=helper control=sequence
//ff:what Provides compiled regex patterns for detecting response status codes in Python frameworks
package analyzer

import "regexp"

var pythonPatterns = []*regexp.Regexp{
	// status=status.HTTP_201_CREATED → must be checked before generic status=(\d+)
	regexp.MustCompile(`status=status\.HTTP_(\d+)`),
	// status=201 (Django, DRF)
	regexp.MustCompile(`status=(\d+)`),
	// status_code=400 (FastAPI)
	regexp.MustCompile(`status_code=(\d+)`),
	// abort(404) or abort(404, message="...") (Flask, Flask-RESTful)
	regexp.MustCompile(`abort\((\d+)`),
	// return jsonify(...), 201 (Flask)
	regexp.MustCompile(`return\s+jsonify\(.*\),\s*(\d+)`),
	// make_response(..., 201) (Flask)
	regexp.MustCompile(`make_response\(.*,\s*(\d+)\)`),
	// return "...", 201 or return data, 201 (Flask tuple response)
	regexp.MustCompile(`return\s+.+,\s*(\d{3})\s*$`),
	// Response(data, status=201) (Flask/Werkzeug)
	regexp.MustCompile(`Response\(.*status=(\d+)`),
	// Response(data, 201) (Flask/Werkzeug positional)
	regexp.MustCompile(`Response\([^,]+,\s*(\d+)`),
	// redirect(url, 301) — explicit redirect status (Flask)
	regexp.MustCompile(`redirect\(.+,\s*(\d+)\)`),
}

// redirect(url) without explicit status — implicit 302 (Flask)
var flaskRedirectImplicit = regexp.MustCompile(`redirect\(\s*(?:url_for|request\.|['"/])`)

// djangoResponseClassPattern matches Django HttpResponseXxx classes.
var djangoResponseClassPattern = regexp.MustCompile(
	`HttpResponse(NotFound|BadRequest|Forbidden|NotAllowed|ServerError|` +
		`Redirect|PermanentRedirect|Gone)\(`,
)

var djangoResponseClassStatus = map[string]int{
	"NotFound":          404,
	"BadRequest":        400,
	"Forbidden":         403,
	"NotAllowed":        405,
	"ServerError":       500,
	"Redirect":          302,
	"PermanentRedirect": 301,
	"Gone":              410,
}

// djangoExceptionPattern matches Django built-in exception raises.
var djangoExceptionPattern = regexp.MustCompile(
	`raise\s+(Http404|PermissionDenied|SuspiciousOperation|BadRequest|ValidationError)`,
)

var djangoExceptionStatus = map[string]int{
	"Http404":             404,
	"PermissionDenied":    403,
	"SuspiciousOperation": 400,
	"BadRequest":          400,
	"ValidationError":     400,
}

// drfExceptionPattern matches DRF exception raises.
var drfExceptionPattern = regexp.MustCompile(
	`raise\s+(NotFound|PermissionDenied|AuthenticationFailed|NotAuthenticated|` +
		`ValidationError|MethodNotAllowed|ParseError|Throttled)\(`,
)

var drfExceptionStatus = map[string]int{
	"NotFound":             404,
	"PermissionDenied":     403,
	"AuthenticationFailed": 401,
	"NotAuthenticated":     401,
	"ValidationError":      400,
	"MethodNotAllowed":     405,
	"ParseError":           400,
	"Throttled":            429,
}
