//ff:func feature=analyzer type=data control=selection
//ff:what Maps Echo built-in error constant names to their HTTP status codes
package analyzer

var echoErrStatus = map[string]int{
	"ErrBadRequest":                  400,
	"ErrUnauthorized":                401,
	"ErrForbidden":                   403,
	"ErrNotFound":                    404,
	"ErrMethodNotAllowed":            405,
	"ErrRequestTimeout":              408,
	"ErrConflict":                    409,
	"ErrUnsupportedMediaType":        415,
	"ErrStatusRequestEntityTooLarge": 413,
	"ErrTooManyRequests":             429,
	"ErrInternalServerError":         500,
	"ErrServiceUnavailable":          503,
}
