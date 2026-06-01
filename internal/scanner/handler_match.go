//ff:type feature=scan type=model
//ff:what handlerMatch is one file:line where a handler definition was found
package scanner

// handlerMatch is one file:line where a handler definition was found.
type handlerMatch struct {
	File string
	Line int
}
