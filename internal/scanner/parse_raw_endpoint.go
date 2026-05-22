//ff:func feature=scan type=parser control=sequence
//ff:what Converts a raw JSON endpoint entry into an Endpoint struct with ID and juicer handler parsing
package scanner

import "strings"

func parseRawEndpoint(r rawEndpoint) *Endpoint {
	if r.Method == "" || r.Path == "" {
		return nil
	}

	handler := r.Handler
	file := r.File

	// juicer format: handler contains "file:funcName"
	if handler != "" && strings.Contains(handler, ":") {
		parts := strings.SplitN(handler, ":", 2)
		if file == "" {
			file = parts[0]
		}
		handler = parts[1]
	}

	return &Endpoint{
		ID:      makeID(r.Method, r.Path),
		Method:  r.Method,
		Path:    r.Path,
		Handler: handler,
		Source:  file,
		Line:    r.Line,
	}
}
