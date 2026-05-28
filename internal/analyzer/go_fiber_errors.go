//ff:func feature=analyzer type=data control=selection
//ff:what Maps Fiber built-in error constant names to their HTTP status codes
package analyzer

var fiberErrStatus = map[string]int{
	"ErrBadRequest":            400,
	"ErrUnauthorized":          401,
	"ErrForbidden":             403,
	"ErrNotFound":              404,
	"ErrMethodNotAllowed":      405,
	"ErrNotAcceptable":         406,
	"ErrRequestTimeout":        408,
	"ErrConflict":              409,
	"ErrGone":                  410,
	"ErrPreconditionFailed":    412,
	"ErrRequestEntityTooLarge": 413,
	"ErrUnsupportedMediaType":  415,
	"ErrUnprocessableEntity":   422,
	"ErrTooManyRequests":       429,
	"ErrInternalServerError":   500,
	"ErrServiceUnavailable":    503,
}
