//ff:func feature=scan type=engine control=iteration dimension=1
//ff:what Links endpoints with empty Source to a handler file:line found under a root directory
package scanner

// LinkSource scans root for source files and, for every endpoint that has a
// Handler name but no Source, locates the file:line where that handler is
// defined. Endpoints whose handler cannot be located are left unlinked
// (Source stays "") — per §3.2 they honestly fall to UNVERIFIED. Returns the
// number of endpoints newly linked.
func LinkSource(endpoints []Endpoint, root string) int {
	files := collectSourceFiles(root)
	linked := 0
	for i := range endpoints {
		if linkEndpoint(&endpoints[i], files) {
			linked++
		}
	}
	return linked
}
