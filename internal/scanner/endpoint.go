//ff:type feature=scan type=model
//ff:what Endpoint represents a discovered API route with method, path, handler, and source location
package scanner

type Endpoint struct {
	ID      string `json:"id"`
	Method  string `json:"method"`
	Path    string `json:"path"`
	Handler string `json:"handler"`
	Source  string `json:"source"`
	Line    int    `json:"line"`
}
