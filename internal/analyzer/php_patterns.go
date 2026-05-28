//ff:func feature=analyzer type=helper control=sequence
//ff:what Provides compiled regex patterns for detecting response status codes in PHP/Laravel source
package analyzer

import "regexp"

// laravelPatterns matches Laravel controller response patterns with explicit status codes.
var laravelPatterns = []*regexp.Regexp{
	// response()->json($data, 201)
	regexp.MustCompile(`response\(\)->json\(.+,\s*(\d+)\)`),
	// response($data, 201)
	regexp.MustCompile(`response\([^,]*,\s*(\d+)\)`),
	// Response::json($data, 201)
	regexp.MustCompile(`Response::json\(.+,\s*(\d+)\)`),
	// new JsonResponse($data, 201)
	regexp.MustCompile(`new JsonResponse\(.+,\s*(\d+)\)`),
	// new Response($data, 201)
	regexp.MustCompile(`new Response\([^,]*,\s*(\d+)\)`),
	// abort(404)
	regexp.MustCompile(`abort\((\d+)`),
	// abort_if(condition, 404)
	regexp.MustCompile(`abort_if\(.+,\s*(\d+)`),
	// abort_unless(condition, 403)
	regexp.MustCompile(`abort_unless\(.+,\s*(\d+)`),
	// redirect($url, 301) — explicit redirect status
	regexp.MustCompile(`redirect\([^,]+,\s*(\d+)\)`),
}

// laravelResourcePatterns matches Laravel API Resource patterns with explicit status codes.
var laravelResourcePatterns = []*regexp.Regexp{
	// return (new UserResource($user))->response()->setStatusCode(201)
	regexp.MustCompile(`->setStatusCode\((\d+)\)`),
}

// laravelExceptionPattern matches Laravel/Symfony exception throws.
var laravelExceptionPattern = regexp.MustCompile(
	`throw\s+new\s+(NotFoundHttpException|BadRequestHttpException|` +
		`AccessDeniedHttpException|UnauthorizedHttpException|` +
		`MethodNotAllowedHttpException|ConflictHttpException|` +
		`GoneHttpException|TooManyRequestsHttpException|` +
		`UnprocessableEntityHttpException|ServiceUnavailableHttpException)\(`,
)

// laravelExceptionStatus maps Laravel/Symfony exception class names to HTTP status codes.
var laravelExceptionStatus = map[string]int{
	"NotFoundHttpException":            404,
	"BadRequestHttpException":          400,
	"AccessDeniedHttpException":        403,
	"UnauthorizedHttpException":        401,
	"MethodNotAllowedHttpException":    405,
	"ConflictHttpException":            409,
	"GoneHttpException":                410,
	"TooManyRequestsHttpException":     429,
	"UnprocessableEntityHttpException": 422,
	"ServiceUnavailableHttpException":  503,
}

// laravelImplicitJson200 matches response()->json($data) without explicit status — implicit 200.
var laravelImplicitJson200 = regexp.MustCompile(`response\(\)->json\([^,]*\)\s*;`)

// laravelImplicitRedirect302 matches redirect() without explicit status — implicit 302.
var laravelImplicitRedirect302 = regexp.MustCompile(`return\s+redirect\(\)->`)

// laravelImplicitResource200 matches return new UserResource($user) — implicit 200.
var laravelImplicitResource200 = regexp.MustCompile(`return\s+new\s+\w+Resource\(`)

// laravelImplicitCollection200 matches return UserResource::collection($users) — implicit 200.
var laravelImplicitCollection200 = regexp.MustCompile(`return\s+\w+Resource::collection\(`)

// laravelImplicitResponseJson200 matches return response()->json($data) — implicit 200.
var laravelImplicitResponseJson200 = regexp.MustCompile(`return\s+response\(\)->json\([^,]*\)`)
