//ff:func feature=scan type=parser control=sequence
//ff:what Matches a gin route pattern in a line and returns an Endpoint or nil
package scanner

func parseRoute(line, path string, lineNum int) *Endpoint {
	matches := ginRoutePattern.FindStringSubmatch(line)
	if matches == nil {
		return nil
	}
	method := matches[1]
	routePath := matches[2]
	handler := extractHandler(line)

	return &Endpoint{
		ID:      makeID(method, routePath),
		Method:  method,
		Path:    routePath,
		Handler: handler,
		Source:  path,
		Line:    lineNum,
	}
}
