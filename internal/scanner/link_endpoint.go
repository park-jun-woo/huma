//ff:func feature=scan type=helper control=sequence
//ff:what Links a single endpoint to a handler definition file:line if it is unlinked and matchable
package scanner

// linkEndpoint sets Source/Line on a single endpoint when it is unlinked and
// its handler can be found. Returns true if a new link was established.
func linkEndpoint(ep *Endpoint, files []string) bool {
	if ep.Source != "" || ep.Handler == "" {
		return false
	}
	file, line := findHandlerDef(files, ep.Handler)
	if file == "" {
		return false
	}
	ep.Source = file
	ep.Line = line
	return true
}
